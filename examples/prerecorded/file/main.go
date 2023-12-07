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

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/prerecorded/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
)

const (
	filePath string = "./Bueller-Life-moves-pretty-fast.mp3"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelTrace, // LogLevelStandard / LogLevelFull / LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// set the Transcription options
	options := interfaces.PreRecordedTranscriptionOptions{
		Punctuate:  true,
		Diarize:    true,
		Language:   "en-US",
		Utterances: true,
	}

	// create a Deepgram client
	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	// send/process file to Deepgram
	res, err := dg.FromFile(ctx, filePath, options)
	if err != nil {
		log.Printf("FromStream failed. Err: %v\n", err)
		os.Exit(1)
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	// make the JSON pretty
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		log.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("\n\nResult:\n%s\n\n", prettyJson)

	// dump example VTT
	vtt, err := res.ToWebVTT()
	if err != nil {
		log.Printf("ToWebVTT failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("\n\n\nVTT:\n%s\n\n\n", vtt)

	// dump example SRT
	srt, err := res.ToSRT()
	if err != nil {
		log.Printf("ToSRT failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("\n\n\nSRT:\n%s\n\n\n", srt)
}
