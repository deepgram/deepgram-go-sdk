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
	"strings"
	"sync"

	klog "k8s.io/klog/v2"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/pkg/audio/microphone"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

type MyHandler struct {
	openChan          chan *msginterfaces.OpenResponse
	messageChan       chan *msginterfaces.MessageResponse
	metadataChan      chan *msginterfaces.MetadataResponse
	speechStartedChan chan *msginterfaces.SpeechStartedResponse
	utteranceEndChan  chan *msginterfaces.UtteranceEndResponse
	closeChan         chan *msginterfaces.CloseResponse
	errorChan         chan *msginterfaces.ErrorResponse
	unhandledChan     chan *[]byte
}

func NewMyHandler() *MyHandler {
	handler := &MyHandler{
		openChan:          make(chan *msginterfaces.OpenResponse),
		messageChan:       make(chan *msginterfaces.MessageResponse),
		metadataChan:      make(chan *msginterfaces.MetadataResponse),
		speechStartedChan: make(chan *msginterfaces.SpeechStartedResponse),
		utteranceEndChan:  make(chan *msginterfaces.UtteranceEndResponse),
		closeChan:         make(chan *msginterfaces.CloseResponse),
		errorChan:         make(chan *msginterfaces.ErrorResponse),
		unhandledChan:     make(chan *[]byte),
	}

	go func() {
		handler.Run()
	}()

	return handler
}

// GetOpen returns the open channels
func (dch MyHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return []*chan *msginterfaces.OpenResponse{&dch.openChan}
}

// GetMessage returns the message channels
func (dch MyHandler) GetMessage() []*chan *msginterfaces.MessageResponse {
	return []*chan *msginterfaces.MessageResponse{&dch.messageChan}
}

// GetMetadata returns the metadata channels
func (dch MyHandler) GetMetadata() []*chan *msginterfaces.MetadataResponse {
	return []*chan *msginterfaces.MetadataResponse{&dch.metadataChan}
}

// GetSpeechStarted returns the speech started channels
func (dch MyHandler) GetSpeechStarted() []*chan *msginterfaces.SpeechStartedResponse {
	return []*chan *msginterfaces.SpeechStartedResponse{&dch.speechStartedChan}
}

// GetUtteranceEnd returns the utterance end channels
func (dch MyHandler) GetUtteranceEnd() []*chan *msginterfaces.UtteranceEndResponse {
	return []*chan *msginterfaces.UtteranceEndResponse{&dch.utteranceEndChan}
}

// GetClose returns the close channels
func (dch MyHandler) GetClose() []*chan *msginterfaces.CloseResponse {
	return []*chan *msginterfaces.CloseResponse{&dch.closeChan}
}

// GetError returns the error channels
func (dch MyHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&dch.errorChan}
}

// GetUnhandled returns the unhandled event channels
func (dch MyHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&dch.unhandledChan}
}

// Open is the callback for when the connection opens
// golintci: funlen
func (dch MyHandler) Run() error {
	wgReceivers := sync.WaitGroup{}

	// open channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.openChan {
			fmt.Printf("\n\n[OpenResponse]\n\n")
		}
	}()

	// message channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for mr := range dch.messageChan {
			sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

			if len(mr.Channel.Alternatives) == 0 || sentence == "" {
				klog.V(7).Infof("DEEPGRAM - no transcript")
				continue
			}

			if mr.IsFinal {
				fmt.Printf("\n[MessageResponse] (Final) %s\n", sentence)
			} else {
				fmt.Printf("\n[MessageResponse] (Interim) %s\n", sentence)
			}
		}
	}()

	// metadata channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for mr := range dch.metadataChan {
			fmt.Printf("\n\nMetadata.RequestID: %s\n", strings.TrimSpace(mr.RequestID))
			fmt.Printf("Metadata.Channels: %d\n", mr.Channels)
			fmt.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(mr.Created))
		}
	}()

	// speech started channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.speechStartedChan {
			fmt.Printf("\n[SpeechStarted]\n")
		}
	}()

	// utterance end channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.utteranceEndChan {
			fmt.Printf("\n[UtteranceEnd]\n")
		}
	}()

	// close channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.closeChan {
			fmt.Printf("\n\n[CloseResponse]\n\n")
		}
	}()

	// error channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for er := range dch.errorChan {
			fmt.Printf("\n[ErrorResponse]\n")
			fmt.Printf("\nError.Type: %s\n", er.ErrCode)
			fmt.Printf("Error.Message: %s\n", er.ErrMsg)
			fmt.Printf("Error.Description: %s\n\n", er.Description)
			fmt.Printf("Error.Variant: %s\n\n", er.Variant)
		}
	}()

	// unhandled event channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for byData := range dch.unhandledChan {
			fmt.Printf("\n[UnhandledEvent]")
			fmt.Printf("Dump:\n%s\n\n", string(*byData))
		}
	}()

	// wait for all receivers to finish
	wgReceivers.Wait()

	return nil
}

func main() {
	// init library
	microphone.Initialize()

	// print instructions
	fmt.Print("\n\nPress ENTER to exit!\n\n")

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
		EnableKeepAlive: true,
	}

	// set the Transcription options
	tOptions := &interfaces.LiveTranscriptionOptions{
		Model:       "nova-3",
		Keyterms:    []string{"deepgram"},
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

	// example on how to send a custom parameter
	// params := make(map[string][]string, 0)
	// params["dictation"] = []string{"true"}
	// ctx = interfaces.WithCustomParameters(ctx, params)

	// implement your own callback
	var callback msginterfaces.LiveMessageChan
	callback = *NewMyHandler()

	// create a Deepgram client
	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, tOptions, callback)
	if err != nil {
		fmt.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
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

	// wait for user input to exit
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

	fmt.Printf("\n\nProgram exiting...\n")
	// time.Sleep(120 * time.Second)
}
