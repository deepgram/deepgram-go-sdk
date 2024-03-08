// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	prettyjson "github.com/hokaccha/go-prettyjson"

	speak "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

const (
	textToSpeech string = "Hello, World!"
	filePath     string = "./test.mp3"
)

func main() {
	// init library
	client.InitWithDefault()

	// Go context
	ctx := context.Background()

	// set the Transcription options
	options := interfaces.SpeakOptions{
		Model: "aura-asteria-en",
	}

	// create a Deepgram client
	c := client.NewWithDefaults()
	dg := speak.New(c)

	// send/process file to Deepgram
	res, err := dg.ToSave(ctx, filePath, textToSpeech, options)
	if err != nil {
		fmt.Printf("FromStream failed. Err: %v\n", err)
		os.Exit(1)
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	// make the JSON pretty
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\nResult:\n%s\n\n", prettyJson)
}
