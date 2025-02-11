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
	filePath string = "./Bueller-Life-moves-pretty-fast.mp3"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelVerbose, // LogLevelStandard / LogLevelFull / LogLevelTrace / LogLevelVerbose
	})

	// Go context
	ctx := context.Background()

	// set the Transcription options
	options := &interfaces.PreRecordedTranscriptionOptions{
		Model:       "nova-3",
		Keyterms:    []string{"Bueller"},
		Punctuate:   true,
		Paragraphs:  true,
		SmartFormat: true,
		Language:    "en-US",
		Utterances:  true,
	}

	// create a Deepgram client
	c := client.NewREST("", &interfaces.ClientOptions{
		Host: "https://api.deepgram.com",
	})
	dg := api.New(c)

	// example on how to send a custom header
	// need to import (
	//	 "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	// )
	//
	// headers := make(map[string][]string, 0)
	// headers["MY-CUSTOM-HEADER"] = []string{"CUSTOM"}
	// ctx = cfginterfaces.WithCustomHeaders(ctx, headers)
	//
	// example on how to send a custom parameter
	// params := make(map[string][]string, 0)
	// params["utterances"] = []string{"true"}
	// ctx = cfginterfaces.WithCustomParameters(ctx, params)

	// send/process file to Deepgram
	res, err := dg.FromFile(ctx, filePath, options)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			fmt.Printf("DEEPGRAM ERROR:\n%s:\n%s\n", e.DeepgramError.ErrCode, e.DeepgramError.ErrMsg)
		}
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

	// dump example VTT
	vtt, err := res.ToWebVTT()
	if err != nil {
		fmt.Printf("ToWebVTT failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\n\nVTT:\n%s\n\n\n", vtt)

	// dump example SRT
	srt, err := res.ToSRT()
	if err != nil {
		fmt.Printf("ToSRT failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\n\nSRT:\n%s\n\n\n", srt)
}
