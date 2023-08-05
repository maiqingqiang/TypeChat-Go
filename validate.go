package typechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type JsonValidator[T any] interface {
	GetSchema() string
	GetTypeName() string
	CreateModuleTextFromJson(jsonObject *T) (string, error)
	Validate(jsonText string) (*T, error)
}

type baseJsonValidator[T any] struct {
	schema     string
	typeName   string
	stripNulls bool
}

func (v *baseJsonValidator[T]) GetSchema() string {
	return v.schema
}

func (v *baseJsonValidator[T]) GetTypeName() string {
	return v.typeName
}

func (v *baseJsonValidator[T]) CreateModuleTextFromJson(jsonObject *T) (string, error) {
	marshal, err := json.Marshal(jsonObject)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("package main\n\nconst json = `%s`", marshal), nil
}

func (v *baseJsonValidator[T]) Validate(jsonText string) (*T, error) {
	var result *T
	err := json.Unmarshal([]byte(jsonText), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewJsonValidator[T any](schema string, typeName string) JsonValidator[T] {
	return &baseJsonValidator[T]{
		schema:     schema,
		typeName:   typeName,
		stripNulls: true,
	}
}

type ProgramValidator struct {
	JsonValidator[Program]
}

const ModuleText = `package main

func program(api API) Result {
	%s
}
`

func NewProgramValidator(schema string) JsonValidator[Program] {
	return &ProgramValidator{
		NewJsonValidator[Program](schema, "Program"),
	}
}

func (v *ProgramValidator) CreateModuleTextFromJson(program *Program) (string, error) {
	stepsLen := len(program.Steps)

	if !(stepsLen > 0 && program.Steps[0].Func != "") {
		return "", errors.New("struct is not a valid program")
	}

	currentStep := 0
	funcBody := ""
	for currentStep < stepsLen {
		if (stepsLen - 1) == currentStep {
			funcBody += fmt.Sprintf("return %s", v.exprToString(program.Steps[currentStep]))
		} else {
			funcBody += fmt.Sprintf("step%d := %s \n", currentStep+1, v.exprToString(program.Steps[currentStep]))
		}

		currentStep++
	}

	return fmt.Sprintf(ModuleText, funcBody), nil
}

func (v *ProgramValidator) exprToString(expr *FuncCall) string {
	return v.objectToString(expr)
}

func (v *ProgramValidator) objectToString(expr *FuncCall) string {
	fn := expr.Func

	if len(expr.Args) > 0 {
		return fmt.Sprintf("api.%s(%s)", fn, v.arrayToString(expr.Args))
	} else {
		return fmt.Sprintf("api.%s()", fn)
	}
}

func (v *ProgramValidator) arrayToString(args []Expression) string {
	var list []string
	for i := range args {
		switch args[i].(type) {
		case map[string]any:
			m := args[i].(map[string]any)

			var a []Expression

			if _, ok := m[Args]; ok {
				for _, item := range m[Args].([]any) {
					a = append(a, item)
				}

				list = append(list, v.objectToString(&FuncCall{
					Func: m[Func].(string),
					Args: a,
				}))
			} else if _, ok := m[Ref]; ok {

			}
		case int:
			list = append(list, fmt.Sprintf("%v", args[i]))
		case float64:
			list = append(list, fmt.Sprintf("%v", args[i]))
		}
	}
	return strings.Join(list, ",")
}
