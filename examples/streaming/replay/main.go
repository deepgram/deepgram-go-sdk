// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// streaming
import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	replay "github.com/deepgram-devs/deepgram-go-sdk/pkg/audio/replay"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/live"
)

func main() {
	/*
		DG Streaming API
	*/
	// init library
	client.InitWithDefault()

	// context
	ctx := context.Background()

	// options
	transcriptOptions := interfaces.LiveTranscriptionOptions{
		Language:    "en-US",
		Punctuate:   true,
		Encoding:    "mulaw",
		Channels:    1,
		Sample_rate: 8000,
	}

	dgClient, err := client.NewForDemo(ctx, transcriptOptions)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// call connect!
	wsconn := dgClient.Connect()
	if wsconn == nil {
		log.Println("Client.Connect failed")
		os.Exit(1)
	}

	/*
		Replay wav into Live stream
	*/
	play, err := replay.New(replay.ReplayOptions{
		FullFilename: "testing.wav",
	})
	if err != nil {
		fmt.Printf("replay.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start replay
	err = play.Start()
	if err != nil {
		fmt.Printf("replay.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	go func() {
		// this is a blocking call
		play.Stream(dgClient)
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close stream
	err = play.Stop()
	if err != nil {
		fmt.Printf("replay.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}

	// close client
	dgClient.Stop()
}
