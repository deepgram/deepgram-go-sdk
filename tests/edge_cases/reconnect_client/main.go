// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// streaming
import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/pkg/audio/microphone"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

// Implement your own callback
type MyCallback struct {
	sb *strings.Builder
}

func (c MyCallback) Message(mr *api.MessageResponse) error {
	// handle the message
	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	if len(mr.Channel.Alternatives) == 0 || sentence == "" {
		return nil
	}

	if mr.IsFinal {
		c.sb.WriteString(sentence)
		c.sb.WriteString(" ")

		if mr.SpeechFinal {
			fmt.Printf("[------- Is Final]: %s\n", c.sb.String())
			c.sb.Reset()
		}
	} else {
		fmt.Printf("[Interim Result]: %s\n", sentence)
	}

	return nil
}

func (c MyCallback) Open(ocr *api.OpenResponse) error {
	// handle the open
	fmt.Printf("\n[Open] Received\n")
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

func (c MyCallback) SpeechStarted(ssr *api.SpeechStartedResponse) error {
	fmt.Printf("\n[SpeechStarted] Received\n")
	return nil
}

func (c MyCallback) UtteranceEnd(ur *api.UtteranceEndResponse) error {
	utterance := strings.TrimSpace(c.sb.String())
	if len(utterance) > 0 {
		fmt.Printf("[------- UtteranceEnd]: %s\n", utterance)
		c.sb.Reset()
	} else {
		fmt.Printf("\n[UtteranceEnd] Received\n")
	}
	return nil
}

func (c MyCallback) Close(ocr *api.CloseResponse) error {
	// handle the close
	fmt.Printf("\n[Close] Received\n")
	return nil
}

func (c MyCallback) Error(er *api.ErrorResponse) error {
	// handle the error
	fmt.Printf("\n[Error] Received\n")
	fmt.Printf("Error.Type: %s\n", er.Type)
	fmt.Printf("Error.ErrCode: %s\n", er.ErrCode)
	fmt.Printf("Error.Description: %s\n\n", er.Description)
	return nil
}

func (c MyCallback) UnhandledEvent(byData []byte) error {
	// handle the unhandled event
	fmt.Printf("\n[UnhandledEvent] Received\n")
	fmt.Printf("UnhandledEvent: %s\n\n", string(byData))
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
		LogLevel: client.LogLevelTrace, // LogLevelDefault, LogLevelFull, LogLevelDebug, LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// client options
	cOptions := &interfaces.ClientOptions{
		EnableKeepAlive: true,
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

	// implement your own callback
	callback := MyCallback{
		sb: &strings.Builder{},
	}

	// create a Deepgram client
	dgClient, err := client.NewWebSocketUsingCallback(ctx, "", cOptions, tOptions, callback)
	if err != nil {
		fmt.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(5 * time.Second)
		}

		// connect the websocket to Deepgram
		bConnected := dgClient.AttemptReconnect(context.Background(), 3)
		if !bConnected {
			fmt.Println("Client.AttemptReconnect failed")
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
			err := mic.Stream(dgClient)
			if err != nil {
				fmt.Printf("mic.Stream non-fatal error. Err: %v\n", err)
			}
		}()

		// sleep for 10 seconds
		time.Sleep(10 * time.Second)

		// close and repeat
		err = mic.Stop()
		if err != nil {
			fmt.Printf("mic.Stop non-fatal error. Err: %v\n", err)
		}
		dgClient.Stop()
	}

	// teardown library
	microphone.Teardown()

	fmt.Printf("\n\nProgram exiting...\n")
}
