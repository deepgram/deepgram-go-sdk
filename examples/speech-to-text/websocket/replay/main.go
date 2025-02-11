// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
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

	replay "github.com/deepgram/deepgram-go-sdk/pkg/audio/replay"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

func main() {
	/*
		DG Streaming API
	*/
	// init library
	client.InitWithDefault()

	// Go context
	ctx := context.Background()

	// set the Transcription options
	options := &interfaces.LiveTranscriptionOptions{
		Model:      "nova-3",
		Keyterms:   []string{"deepgram"},
		Language:   "en-US",
		Punctuate:  true,
		Encoding:   "mulaw",
		Channels:   1,
		SampleRate: 8000,
	}

	// create a Deepgram client
	dgClient, err := client.NewWSUsingChanForDemo(ctx, options)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
		log.Println("Client.Connect failed")
		os.Exit(1)
	}

	/*
		Replay wav into Live stream
	*/
	play, err := replay.New(replay.Options{
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
		// feed the WAV stream to the Deepgram client (this is a blocking call)
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

	log.Printf("Program exiting...\n")
}
