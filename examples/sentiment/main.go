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

	schema, err := os.ReadFile("sentiment_schema.go")
	if err != nil {
		log.Fatalf("os.ReadFile Error: %v\n", err)
	}

	translator := typechat.NewJsonTranslator[SentimentResponse](model, string(schema), "SentimentResponse")

	_ = typechat.ProcessRequests("😀> ", os.Args[1], func(request string) error {
		response, err := translator.Translate(request)
		if err != nil {
			log.Fatalf("translator.Translate Error: %v\n", err)
		}

		log.Printf("The sentiment is %s\n", Sentiment(response.Sentiment))
		return nil
	})
}
