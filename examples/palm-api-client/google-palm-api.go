package main

import (
	"log"

	"github.com/raghvenders/generative-ai-go/palm2"
)

func main() {
	palm2, err := palm2.NewPalm()
	if err != nil {
		log.Fatalf("Unable to Create Palm2 Client : %v", err)
	}

	models, err := palm2.ListModels()
	if err != nil {
		log.Fatalf("Error Response %v", err)
	}
	log.Printf("List Models : %+v \n", models)

	res, err := palm2.GenerateText("Hello")
	if err != nil {
		log.Fatalf("Error Response %v", err)
	}
	log.Printf("Response : %+v \n", res)
}
