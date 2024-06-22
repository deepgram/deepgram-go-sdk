// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	prettyjson "github.com/hokaccha/go-prettyjson"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/prerecorded/v1" //lint:ignore
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded" //lint:ignore
)

const (
	filePath string = "./CallCenterPhoneCall.mp3"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault, // LogLevelStandard / LogLevelFull / LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// set the Transcription options
	options := &interfaces.PreRecordedTranscriptionOptions{
		Model:       "nova-2",
		Punctuate:   true,
		Language:    "en-US",
		SmartFormat: true,
		Intents:     true,
	}

	// create a Deepgram client
	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	// send/process file to Deepgram
	res, err := dg.FromFile(ctx, filePath, options)
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
	prettyJSON, err := prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\nResult:\n%s\n\n", prettyJSON)
}
