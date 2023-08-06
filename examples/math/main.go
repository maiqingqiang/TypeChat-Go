package main

import (
	"fmt"
	"github.com/maiqingqiang/typechat-go"
	"github.com/spf13/cast"
	"log"
	"os"
)

func main() {
	model, err := typechat.NewLanguageModel()
	if err != nil {
		log.Fatal(err)
	}

	schema, err := os.ReadFile("math_schema.go")
	if err != nil {
		log.Fatalf("os.ReadFile Error: %v\n", err)
	}

	translator := typechat.NewProgramTranslator(model, string(schema))

	_ = typechat.ProcessRequests("➕➖✖️➗🟰> ", os.Args[1], func(request string) error {
		response, err := translator.Translate(request)
		if err != nil {
			log.Fatalf("translator.Translate Error: %v\n", err)
		}

		programStr, err := translator.Validator().CreateModuleTextFromJson(response)
		if err != nil {
			log.Fatalf("CreateModuleTextFromJson Error: %v\n", err)
		}
		log.Println(programStr)

		log.Println(fmt.Sprintf("Running program:"))
		result, err := typechat.EvaluateJsonProgram(response, handleCall)
		if err != nil {
			log.Fatalf("EvaluateJsonProgram Error: %v\n", err)
		}
		log.Println(fmt.Sprintf("Result: %d", result))

		return nil
	})

}

func handleCall(fn string, args []typechat.Expression) (typechat.Result, error) {
	switch fn {
	case "Add":
		return cast.ToInt(args[0]) + cast.ToInt(args[1]), nil
	case "Sub":
		return cast.ToInt(args[0]) - cast.ToInt(args[1]), nil
	case "Mul":
		return cast.ToInt(args[0]) * cast.ToInt(args[1]), nil
	case "Div":
		return cast.ToInt(args[0]) / cast.ToInt(args[1]), nil
	case "Neg":
		return -cast.ToInt(args[0]), nil
	case "Id":
		return cast.ToInt(args[0]), nil
	}

	return 0, nil
}
