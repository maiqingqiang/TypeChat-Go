package main

import (
	"github.com/maiqingqiang/typechat"
	"log"
	"os"
)

func main() {
	model, err := typechat.NewLanguageModel()
	if err != nil {
		log.Fatal(err)
	}

	schema, err := os.ReadFile("calendar_actions_schema.go")
	if err != nil {
		log.Fatalf("os.ReadFile Error: %v\n", err)
	}

	translator := typechat.NewJsonTranslator[CalendarActions](model, string(schema), "CalendarActions")

	_ = typechat.ProcessRequests("📅> ", os.Args[1], func(request string) error {
		calendarActions, err := translator.Translate(request)
		if err != nil {
			log.Fatalf("translator.Translate Error: %v\n", err)
		}

		log.Printf("%+v", calendarActions)

		for i := range calendarActions.Actions {
			if calendarActions.Actions[i].UnknownAction != nil {
				log.Fatalf("I didn't understand the following:\n%s", calendarActions.Actions[i].UnknownAction.Text)
			}
		}

		return nil
	})
}
