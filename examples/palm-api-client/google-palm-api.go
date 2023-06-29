package main

import (
	"log"

	"github.com/raghvenders/generative-ai-go/palm2"
)

func main() {
	client, err := palm2.NewPalm()
	if err != nil {

		log.Fatalln(err, " \n !! Kindly create a API Key or use an existing API key from the makersuite.google.com !!")
	}

	_ = client
}
