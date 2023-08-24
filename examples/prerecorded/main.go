package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/deepgram-devs/deepgram-go-sdk/deepgram"
)

func main() {
	credentials := "DEEPGRAM_API_KEY"
	dg := deepgram.NewClient(credentials)

	filePath := "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"
	var res interface{}
	var err error

	if isURL(filePath) {
		res, err = dg.PreRecordedFromURL(
			deepgram.UrlSource{Url: filePath},
			deepgram.PreRecordedTranscriptionOptions{
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

		source := deepgram.ReadStreamSource{Stream: file, Mimetype: "YOUR_FILE_MIME_TYPE"}

		res, err = dg.PreRecordedFromStream(
			source,
			deepgram.PreRecordedTranscriptionOptions{
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
