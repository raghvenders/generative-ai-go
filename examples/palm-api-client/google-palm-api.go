package main

import (
	"log"
	"os"

	"github.com/raghvenders/generative-ai-go/palm2"
)

func main() {
	/*
		client, err := palm2.NewPalm()
		if err != nil {
			log.Fatalln(err, "(Suggestion: If key is unavailable, Kindly create an API Key or use an existing API key from ", restclient.HyperLink("https://makersuite.google.com/app/apikey", "makersuite.google.com"))
		}

		_ = client

	*/

	//urlQueryParams := palm2.New().WithBaseUrl("https://generativelanguage.googleapis.com/v1beta2/").WithQueryParams(map[string]any{"key": "123", "q": "apple", "min_price": 100, "max price": 1000})

	//urlPathParams := palm2.New().WithBaseUrl("https://generativelanguage.googleapis.com/v1beta2/").WithPathParams("models").WithQueryParams(map[string]any{"key": os.Args[1]})

	//fmt.Printf("%+v \n", *rb)

	if len(os.Args) == 1 {
		log.Fatalln("Unable to create PaLM2 Client as apiKey is missing")
	}

	res, err := palm2.New(os.Args[1]).ListModels()
	if err != nil {
		log.Fatalln("Response", err)
	}
	log.Printf("List Models : \n %+v", res)
}
