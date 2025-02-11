// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

func main() {
	// init library
	client.InitWithDefault()

	// Go context
	ctx := context.Background()

	// set the Transcription options
	tOptions := &interfaces.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}

	// create a Deepgram client
	cOptions := &interfaces.ClientOptions{
		EnableKeepAlive: true,
	}

	// use the default callback handler which just dumps all messages to the screen
	dgClient, err := client.NewWebSocketUsingChan(ctx, "", cOptions, tOptions, nil)
	if err != nil {
		fmt.Println("ERROR creating LiveClient connection:", err)
		return
	}

	// connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Println("Client.Connect failed")
		os.Exit(1)
	}

	// wait for user input to exit
	fmt.Printf("This demonstrates using KeepAlives. Press ENTER to exit...\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close client
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
}
