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
	"strings"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/pkg/audio/microphone"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/live"
)

type MyCallback struct{}

func (c MyCallback) Message(mr *api.MessageResponse) error {
	// handle the message
	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		return nil
	}
	log.Printf("\n%s\n", sentence)
	return nil
}
func (c MyCallback) Metadata(md *api.MetadataResponse) error {
	// handle the metadata
	log.Printf("\nMetadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	log.Printf("Metadata.Channels: %d\n", md.Channels)
	log.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(md.Created))
	return nil
}
func (c MyCallback) Error(er *api.ErrorResponse) error {
	// handle the error
	log.Printf("\nError.Type: %s\n", er.Type)
	log.Printf("Error.Message: %s\n", er.Message)
	log.Printf("Error.Description: %s\n\n", er.Description)
	return nil
}

func main() {
	// init library
	microphone.Initialize()

	/*
		DG Streaming API
	*/
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault, // LogLevelDefault, LogLevelFull, LogLevelTrace
	})

	// context
	ctx := context.Background()

	// options
	options := interfaces.LiveTranscriptionOptions{
		Language:   "en-US",
		Punctuate:  true,
		Encoding:   "linear16",
		Channels:   1,
		SampleRate: 16000,
	}

	// callback
	callback := MyCallback{}

	dgClient, err := client.NewWithDefaults(ctx, options, callback)
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

	log.Printf("Program exiting...\n")
	// time.Sleep(120 * time.Second)

}
