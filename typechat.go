package typechat

import (
	"fmt"
	"strings"
)

type JsonTranslator[T any] interface {
	CreateRequestPrompt(request string) string
	CreateRepairPrompt(validationError string) string
	Translate(request string) (*T, error)
	Validator() JsonValidator[T]
	Model() LanguageModel
}

type baseJsonTranslator[T any] struct {
	model         LanguageModel
	validator     JsonValidator[T]
	attemptRepair bool
	stripNulls    bool
}

func NewJsonTranslator[T any](model LanguageModel, schema string, typeName string) JsonTranslator[T] {
	return &baseJsonTranslator[T]{
		model:         model,
		validator:     NewJsonValidator[T](schema, typeName),
		attemptRepair: true,
	}
}

func (t *baseJsonTranslator[T]) CreateRequestPrompt(request string) string {
	return fmt.Sprintf("You are a service that translates user requests into JSON objects of type \"%s\" according to the following Golang definitions:\n"+
		"```\n%s```\n"+
		"The following is a user request:\n"+
		"\"\"\"\n%s\n\"\"\"\n"+
		"The following is the user request translated into a JSON object with 1 spaces of indentation and no properties with the value undefined:\n",
		t.validator.GetTypeName(), t.validator.GetSchema(), request)
}

func (t *baseJsonTranslator[T]) CreateRepairPrompt(validationError string) string {
	return fmt.Sprintf("The JSON object is invalid for the following reason:\n"+
		"\"\"\"\n%s\n\"\"\"\n"+
		"The following is a revised JSON object:\n", validationError)
}

func (t *baseJsonTranslator[T]) Translate(request string) (*T, error) {
	prompt := t.CreateRequestPrompt(request)
	attemptRepair := t.attemptRepair

	for {
		resp, err := t.model.complete(prompt)
		if err != nil {
			return nil, err
		}

		startIndex := strings.Index(resp, "{")
		endIndex := strings.LastIndex(resp, "}")

		if !(startIndex >= 0 && endIndex > startIndex) {
			return nil, fmt.Errorf("Response is not JSON:\n%s", resp)
		}

		jsonText := resp[startIndex : endIndex+1]

		result, err := t.validator.Validate(jsonText)

		if err == nil {
			return result, nil
		}

		if !attemptRepair {
			return nil, fmt.Errorf("JSON validation failed: %v\n%s", err, jsonText)
		}

		prompt += fmt.Sprintf("%s\n%s", jsonText, t.CreateRepairPrompt(err.Error()))
		attemptRepair = false
	}
}

func (t *baseJsonTranslator[T]) Validator() JsonValidator[T] {
	return t.validator
}

func (t *baseJsonTranslator[T]) Model() LanguageModel {
	return t.model
}
