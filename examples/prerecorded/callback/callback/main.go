// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	prettyjson "github.com/hokaccha/go-prettyjson"

	prerecorded "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/prerecorded/v1"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

const (
	url string = "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"
)

func main() {
	// init library
	client.InitWithDefault()

	// context
	ctx := context.Background()

	//client
	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	// send stream
	res, err := dg.FromURL(
		ctx,
		url,
		interfaces.PreRecordedTranscriptionOptions{
			Punctuate:  true,
			Diarize:    true,
			Language:   "en-US",
			Utterances: true,
			Callback:   "https://136.22.14.123:3000/v1/callback",
		},
	)
	if err != nil {
		log.Printf("FromURL failed. Err: %v\n", err)
		os.Exit(1)
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		log.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("\n\nResult:\n%s\n\n", prettyJson)
}
