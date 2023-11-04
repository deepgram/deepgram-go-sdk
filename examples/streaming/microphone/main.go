// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// streaming
import (
	"bufio"
	"context"
	"log"
	"os"

	microphone "github.com/deepgram-devs/deepgram-go-sdk/pkg/audio/microphone"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/live"
)

func main() {
	// init library
	microphone.Initialize()

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
		Encoding:    "linear16",
		Channels:    1,
		Sample_rate: 16000,
	}

	dgClient, err := client.NewWithDefaults(ctx, "", transcriptOptions)
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
		Microphone package
	*/
	// mic stuf
	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  16000,
	})
	if err != nil {
		log.Printf("Initialize failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start the mic
	err = mic.Start()
	if err != nil {
		log.Printf("mic.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	go func() {
		// this is a blocking call
		mic.Stream(dgClient)
	}()

	log.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close mic stream
	err = mic.Stop()
	if err != nil {
		log.Printf("mic.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}

	// teardown library
	microphone.Teardown()

	// close DG client
	dgClient.Stop()
}
