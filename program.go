package typechat

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
)

var (
	//go:embed program_schema.tpl
	programSchemaText string
)

const Steps = "@step"
const Func = "@func"
const Args = "@args"
const Ref = "@ref"

type Program struct {
	Steps []*FuncCall `json:"@steps"`
}

type FuncCall struct {
	Func string       `json:"@func"`
	Args []Expression `json:"@args,omitempty"`
}

// Expression is a int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string | Program | ResultReference
type Expression any
type Result any

type ResultReference struct {
	Ref int `json:"@ref"`
}

type ProgramTranslator struct {
	JsonTranslator[Program]
}

func NewProgramTranslator(model LanguageModel, schema string) JsonTranslator[Program] {
	return &ProgramTranslator{
		&baseJsonTranslator[Program]{
			model:     model,
			validator: NewProgramValidator(schema),
		},
	}
}

func (t *ProgramTranslator) CreateRequestPrompt(request string) string {
	return fmt.Sprintf("You are a service that translates user requests into programs represented as JSON using the following TypeScript definitions:\n"+
		"```\n%s```\n"+
		"The programs can call functions from the API defined in the following TypeScript definitions:\n"+
		"```\n%s```\n"+
		"The following is a user request:\n"+
		"```\n%s\n```\n"+
		"The following is the user request translated into a JSON program object with 2 spaces of indentation and no properties with the value undefined:\n",
		programSchemaText, t.Validator().GetSchema(), request)
}

func (t *ProgramTranslator) CreateRepairPrompt(validationError string) string {
	return fmt.Sprintf("The JSON program object is invalid for the following reason:\n"+
		"\"\"\"\n%s\n\"\"\""+
		"The following is a revised JSON program object:\n",
		validationError)
}

func (t *ProgramTranslator) Translate(request string) (*Program, error) {
	prompt := t.CreateRequestPrompt(request)

	resp, err := t.JsonTranslator.Model().complete(prompt)
	if err != nil {
		return nil, err
	}

	startIndex := strings.Index(resp, "{")
	endIndex := strings.LastIndex(resp, "}")

	if !(startIndex >= 0 && endIndex > startIndex) {
		return nil, errors.New(fmt.Sprintf("Response is not JSON:\n%s", resp))
	}

	jsonText := resp[startIndex : endIndex+1]
	program, err := t.Validator().Validate(jsonText)
	if err == nil {
		return program, nil
	}

	prompt += fmt.Sprintf("%s\n%s", jsonText, t.CreateRepairPrompt(err.Error()))

	return nil, nil
}

func (t *ProgramTranslator) Validator() JsonValidator[Program] {
	return t.JsonTranslator.Validator()
}

type OnCallFunc func(fn string, args []Expression) (Result, error)

func EvaluateJsonProgram(program *Program, onCall OnCallFunc) (Result, error) {
	var results []Result

	for _, step := range program.Steps {
		result, err := evaluate(step, onCall, results)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if len(results) > 0 {
		return results[len(results)-1], nil
	}
	return nil, nil
}

func evaluate(funcCall *FuncCall, onCall OnCallFunc, results []Result) (Result, error) {
	var expressions []Expression

	for i := range funcCall.Args {
		switch funcCall.Args[i].(type) {
		case map[string]any:
			m := funcCall.Args[i].(map[string]any)
			if _, ok := m[Func]; ok {
				result, err := onCall(m[Func].(string), evaluateArray(m[Args].([]any), onCall))
				if err != nil {
					return nil, err
				}
				expressions = append(expressions, result)

			} else if _, ok := m[Ref]; ok {
				expressions = append(expressions, results[int(m[Ref].(float64))])
			}
		case int:
			expressions = append(expressions, funcCall.Args[i])
		case float64:
			expressions = append(expressions, funcCall.Args[i])
		}
	}

	result, err := onCall(funcCall.Func, expressions)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func evaluateArray(args []any, onCall OnCallFunc) []Expression {
	var expressions []Expression
	for _, arg := range args {
		switch arg.(type) {
		case map[string]any:
			m := arg.(map[string]any)
			if _, ok := m[Func]; ok {
				result, err := onCall(m[Func].(string), evaluateArray(m[Args].([]any), onCall))
				if err != nil {
					return nil
				}
				expressions = append(expressions, result)
			}
		case int:
			expressions = append(expressions, arg)
		case float64:
			expressions = append(expressions, arg)
		}
	}

	return expressions
}
