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

	analyze "github.com/deepgram/deepgram-go-sdk/pkg/api/analyze/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/analyze"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

const (
	filePath string = "./conversation.txt"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelStandard, // LogLevelStandard / LogLevelFull / LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// set the anaylze options
	rOptions := &interfaces.AnalyzeOptions{
		Language:  "en",
		Summarize: true,
	}

	// create a Deepgram client
	c := client.NewWithDefaults()
	dg := analyze.New(c)

	// send/process file to Deepgram
	res, err := dg.FromFile(ctx, filePath, rOptions)
	if err != nil {
		fmt.Printf("FromFile failed. Err: %v\n", err)
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
