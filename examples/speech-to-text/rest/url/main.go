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

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

const (
	url string = "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"
)

func main() {
	// init library
	client.InitWithDefault()

	// context
	ctx := context.Background()

	// send stream to Deepgram
	options := &interfaces.PreRecordedTranscriptionOptions{
		Model:      "nova-3",
		Keyterms:   []string{"deepgram"},
		Punctuate:  true,
		Diarize:    true,
		Language:   "en-US",
		Utterances: true,
		Redact:     []string{"pci", "ssn"},
	}

	// create a Deepgram client
	c := client.NewRESTWithDefaults()
	dg := api.New(c)

	// send the URL to Deepgram
	res, err := dg.FromURL(ctx, url, options)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			fmt.Printf("DEEPGRAM ERROR:\n%s:\n%s\n", e.DeepgramError.ErrCode, e.DeepgramError.ErrMsg)
		}
		fmt.Printf("FromURL failed. Err: %v\n", err)
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
