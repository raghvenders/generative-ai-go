package main

import (
	"log"

	"github.com/raghvenders/generative-ai-go/palm2"
)

func main() {
	res, err := palm2.NewPalm()
	if err != nil {
		log.Fatalf("Error Response %v", err)
	}

	models, err := res.ListModels()
	if err != nil {
		log.Fatalf("Error Response %v", err)
	}
	log.Printf("List Models : \n %+v", models)
}
