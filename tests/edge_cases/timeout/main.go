// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
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

	// use the default callback handler which just dumps all messages to the screen
	dgClient, err := client.NewWebSocketWithDefaults(ctx, tOptions, nil)
	if err != nil {
		fmt.Println("ERROR creating LiveClient connection:", err)
		return
	}

	// info
	fmt.Printf("\n\nThis should timeout in roughly 12 seconds from now with a CloseResponse message.\n")
	fmt.Printf("\n\nYou will first see an OpenResponse message followed by CloseResponse in 12 seconds.\n")

	// connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Println("Client.Connect failed")
		os.Exit(1)
	}

	// wait for user input to exit
	fmt.Printf("Press ENTER to exit...\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close client
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
}
