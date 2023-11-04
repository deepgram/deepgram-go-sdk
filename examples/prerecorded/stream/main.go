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
	filePath string = "./Bueller-Life-moves-pretty-fast.mp3"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelFull,
	})

	// context
	ctx := context.Background()

	//client
	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	// send file stream
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("os.Open(%s) failed. Err: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	res, err := dg.FromStream(
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
		log.Printf("FromStream failed. Err: %v\n", err)
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
