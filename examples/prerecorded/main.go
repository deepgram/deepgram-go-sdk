// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"

	prerecorded "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/prerecorded/v1"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

func main() {
	// context
	ctx := context.Background()

	//client
	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	filePath := "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"
	var res interface{}
	var err error

	// send stream
	if isURL(filePath) {
		res, err = dg.FromURL(
			ctx,
			filePath,
			interfaces.PreRecordedTranscriptionOptions{
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

		res, err = dg.FromStream(
			ctx,
			file,
			interfaces.PreRecordedTranscriptionOptions{
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

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("RecognitionResult json.Marshal failed. Err: %v\n", err)
		return
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		log.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		return
	}
	log.Printf("\n\nResult:\n%s\n\n", prettyJson)
}

// Function to check if a string is a valid URL
func isURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}
