// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// streaming
import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/pkg/audio/microphone"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/live"
)

// Implement your own callback
type MyCallback struct{}

func (c MyCallback) Message(mr *api.MessageResponse) error {
	// handle the message
	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		return nil
	}
	fmt.Printf("\n%s\n", sentence)
	return nil
}

func (c MyCallback) Metadata(md *api.MetadataResponse) error {
	// handle the metadata
	fmt.Printf("\n[Metadata] Received\n")
	fmt.Printf("Metadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	fmt.Printf("Metadata.Channels: %d\n", md.Channels)
	fmt.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(md.Created))
	return nil
}

func (c MyCallback) UtteranceEnd(ur *api.UtteranceEndResponse) error {
	fmt.Printf("\n[UtteranceEnd] Received\n")
	return nil
}

func (c MyCallback) Error(er *api.ErrorResponse) error {
	// handle the error
	fmt.Printf("\n[Error] Received\n")
	fmt.Printf("Error.Type: %s\n", er.Type)
	fmt.Printf("Error.Message: %s\n", er.Message)
	fmt.Printf("Error.Description: %s\n\n", er.Description)
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
		LogLevel: client.LogLevelDebug, // LogLevelDefault, LogLevelFull, LogLevelDebug, LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// set the Transcription options
	options := interfaces.LiveTranscriptionOptions{
		Model:      "nova-2",
		Language:   "en-US",
		Punctuate:  true,
		Encoding:   "linear16",
		Channels:   1,
		SampleRate: 16000,
		// To get UtteranceEnd, the following must be set:
		// InterimResults: true,
		// UtteranceEndMs: "1000",
	}

	// implement your own callback
	callback := MyCallback{}

	// create a Deepgram client
	dgClient, err := client.NewWithDefaults(ctx, options, callback)
	if err != nil {
		fmt.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// connect the websocket to Deepgram
	wsconn := dgClient.Connect()
	if wsconn == nil {
		fmt.Println("Client.Connect failed")
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
		fmt.Printf("Initialize failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start the mic
	err = mic.Start()
	if err != nil {
		fmt.Printf("mic.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	go func() {
		// feed the microphone stream to the Deepgram client (this is a blocking call)
		mic.Stream(dgClient)
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close mic stream
	err = mic.Stop()
	if err != nil {
		fmt.Printf("mic.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}

	// teardown library
	microphone.Teardown()

	// close DG client
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
	// time.Sleep(120 * time.Second)

}
