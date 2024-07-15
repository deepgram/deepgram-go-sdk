// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// streaming
import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	microphone "github.com/deepgram/deepgram-go-sdk/pkg/audio/microphone"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

func main() {
	// init library
	microphone.Initialize()

	/*
		DG Streaming API
	*/
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault, // LogLevelDefault, LogLevelFull, LogLevelDebug, LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// client options
	cOptions := &interfaces.ClientOptions{
		EnableKeepAlive:     true,
		AutoFlushReplyDelta: 2000,
	}

	// set the Transcription options
	tOptions := &interfaces.LiveTranscriptionOptions{
		Model:       "nova-2",
		Language:    "en-US",
		Punctuate:   true,
		Encoding:    "linear16",
		Channels:    1,
		SampleRate:  16000,
		SmartFormat: true,
		VadEvents:   true,
		// To get UtteranceEnd, the following must be set:
		InterimResults: true,
		UtteranceEndMs: "1000",
		// End of UtteranceEnd settings
	}

	// create a Deepgram client
	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, tOptions, nil)
	if err != nil {
		fmt.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// instructions
	fmt.Print("\n\nPress ENTER to exit!\n\n")

	// connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Println("Client.Connect failed")
		os.Exit(1)
	}

	go func() {
		// feed the microphone stream to the Deepgram client (this is a blocking call)
		for {
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
			err = mic.Stream(dgClient)
			if err != nil {
				fmt.Printf("mic.Stream failed. Err: %v\n", err)
				time.Sleep(2 * time.Second)
			}

			// close mic stream
			mic.Stop()

			if err == nil {
				break
			}
		}
	}()

	// wait for user input to exit
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// teardown library
	microphone.Teardown()

	// close DG client
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
	// time.Sleep(120 * time.Second)

}
