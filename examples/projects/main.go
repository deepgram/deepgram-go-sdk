package main

import (
	"log"
	"os"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

func main() {
	var deepgramApiKey string
	if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
		log.Println("DEEPGRAM_API_KEY found")
		deepgramApiKey = v
	} else {
		log.Fatal("DEEPGRAM_API_KEY not found")
		os.Exit(1)
	}

	dg := client.New(deepgramApiKey)

	mgClient := api.New(dg)

	resp, err := mgClient.ListProjects()
	if err != nil {
		log.Printf("ListProjects failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, project := range resp.Projects {
		log.Printf("Name: %s\n", project.Name)
	}
}
