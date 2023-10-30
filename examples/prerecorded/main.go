package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/prerecorded/v1"
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

	prClient := api.New(dg)

	filePath := "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"
	var res interface{}
	var err error

	if isURL(filePath) {
		res, err = prClient.PreRecordedFromURL(
			api.UrlSource{Url: filePath},
			api.PreRecordedTranscriptionOptions{
				Punctuate:  true,
				Diarize:    true,
				Language:   "en-US",
				Utterances: true,
			},
		)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			log.Panicf("error opening file %s: %v", filePath, err)
		}
		defer file.Close()

		source := api.ReadStreamSource{Stream: file, Mimetype: "YOUR_FILE_MIME_TYPE"}

		res, err = prClient.PreRecordedFromStream(
			source,
			api.PreRecordedTranscriptionOptions{
				Punctuate:  true,
				Diarize:    true,
				Language:   "en-US",
				Utterances: true,
			},
		)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
	}

	jsonStr, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	log.Printf("%s", jsonStr)
}

// Function to check if a string is a valid URL
func isURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}
