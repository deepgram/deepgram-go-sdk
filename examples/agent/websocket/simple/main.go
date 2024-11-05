// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// streaming
import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/pkg/audio/microphone"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/agent"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

type MyHandler struct {
	binaryChan                   chan *[]byte
	openChan                     chan *msginterfaces.OpenResponse
	welcomeResponse              chan *msginterfaces.WelcomeResponse
	conversationTextResponse     chan *msginterfaces.ConversationTextResponse
	userStartedSpeakingResponse  chan *msginterfaces.UserStartedSpeakingResponse
	agentThinkingResponse        chan *msginterfaces.AgentThinkingResponse
	functionCallRequestResponse  chan *msginterfaces.FunctionCallRequestResponse
	functionCallingResponse      chan *msginterfaces.FunctionCallingResponse
	agentStartedSpeakingResponse chan *msginterfaces.AgentStartedSpeakingResponse
	agentAudioDoneResponse       chan *msginterfaces.AgentAudioDoneResponse
	endOfThoughtResponse         chan *msginterfaces.EndOfThoughtResponse
	closeChan                    chan *msginterfaces.CloseResponse
	errorChan                    chan *msginterfaces.ErrorResponse
	unhandledChan                chan *[]byte
}

func NewMyHandler() *MyHandler {
	handler := &MyHandler{
		binaryChan:                   make(chan *[]byte),
		openChan:                     make(chan *msginterfaces.OpenResponse),
		welcomeResponse:              make(chan *msginterfaces.WelcomeResponse),
		conversationTextResponse:     make(chan *msginterfaces.ConversationTextResponse),
		userStartedSpeakingResponse:  make(chan *msginterfaces.UserStartedSpeakingResponse),
		agentThinkingResponse:        make(chan *msginterfaces.AgentThinkingResponse),
		functionCallRequestResponse:  make(chan *msginterfaces.FunctionCallRequestResponse),
		functionCallingResponse:      make(chan *msginterfaces.FunctionCallingResponse),
		agentStartedSpeakingResponse: make(chan *msginterfaces.AgentStartedSpeakingResponse),
		agentAudioDoneResponse:       make(chan *msginterfaces.AgentAudioDoneResponse),
		endOfThoughtResponse:         make(chan *msginterfaces.EndOfThoughtResponse),
		closeChan:                    make(chan *msginterfaces.CloseResponse),
		errorChan:                    make(chan *msginterfaces.ErrorResponse),
		unhandledChan:                make(chan *[]byte),
	}

	go func() {
		handler.Run()
	}()

	return handler
}

// GetBinary returns the binary channels
func (dch MyHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&dch.binaryChan}
}

// GetOpen returns the open channels
func (dch MyHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return []*chan *msginterfaces.OpenResponse{&dch.openChan}
}

// GetWelcomeResponse returns the welcome response channels
func (dch MyHandler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return []*chan *msginterfaces.WelcomeResponse{&dch.welcomeResponse}
}

// GetConversationTextResponse returns the conversation text response channels
func (dch MyHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&dch.conversationTextResponse}
}

// GetUserStartedSpeakingResponse returns the user started speaking response channels
func (dch MyHandler) GetUserStartedSpeaking() []*chan *msginterfaces.UserStartedSpeakingResponse {
	return []*chan *msginterfaces.UserStartedSpeakingResponse{&dch.userStartedSpeakingResponse}
}

// GetAgentThinkingResponse returns the agent thinking response channels
func (dch MyHandler) GetAgentThinking() []*chan *msginterfaces.AgentThinkingResponse {
	return []*chan *msginterfaces.AgentThinkingResponse{&dch.agentThinkingResponse}
}

// GetFunctionCallRequestResponse returns the function call request response channels
func (dch MyHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&dch.functionCallRequestResponse}
}

// GetFunctionCallingResponse returns the function calling response channels
func (dch MyHandler) GetFunctionCalling() []*chan *msginterfaces.FunctionCallingResponse {
	return []*chan *msginterfaces.FunctionCallingResponse{&dch.functionCallingResponse}
}

// GetAgentStartedSpeakingResponse returns the agent started speaking response channels
func (dch MyHandler) GetAgentStartedSpeaking() []*chan *msginterfaces.AgentStartedSpeakingResponse {
	return []*chan *msginterfaces.AgentStartedSpeakingResponse{&dch.agentStartedSpeakingResponse}
}

// GetAgentAudioDoneResponse returns the agent audio done response channels
func (dch MyHandler) GetAgentAudioDone() []*chan *msginterfaces.AgentAudioDoneResponse {
	return []*chan *msginterfaces.AgentAudioDoneResponse{&dch.agentAudioDoneResponse}
}

func (dch MyHandler) GetEndOfThought() []*chan *msginterfaces.EndOfThoughtResponse {
	return []*chan *msginterfaces.EndOfThoughtResponse{&dch.endOfThoughtResponse}
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

	// binary channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		counter := 0
		lastBytesReceived := time.Now().Add(-7 * time.Second)

		for br := range dch.binaryChan {
			fmt.Printf("\n\n[Binary Data]\n\n")
			fmt.Printf("Size: %d\n\n", len(*br))

			if lastBytesReceived.Add(5 * time.Second).Before(time.Now()) {
				counter = counter + 1
				file, err := os.OpenFile(fmt.Sprintf("output_%d.wav", counter), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
				if err != nil {
					fmt.Printf("Failed to open file. Err: %v\n", err)
					continue
				}
				// Add a wav audio container header to the file if you want to play the audio
				// using a media player like VLC, Media Player, or Apple Music
				header := []byte{
					0x52, 0x49, 0x46, 0x46, // "RIFF"
					0x00, 0x00, 0x00, 0x00, // Placeholder for file size
					0x57, 0x41, 0x56, 0x45, // "WAVE"
					0x66, 0x6d, 0x74, 0x20, // "fmt "
					0x10, 0x00, 0x00, 0x00, // Chunk size (16)
					0x01, 0x00, // Audio format (1 for PCM)
					0x01, 0x00, // Number of channels (1)
					0x80, 0x3e, 0x00, 0x00, // Sample rate (16000)
					0x00, 0x7d, 0x00, 0x00, // Byte rate (16000 * 2)
					0x02, 0x00, // Block align (2)
					0x10, 0x00, // Bits per sample (16)
					0x64, 0x61, 0x74, 0x61, // "data"
					0x00, 0x00, 0x00, 0x00, // Placeholder for data size
				}

				_, err = file.Write(header)
				if err != nil {
					fmt.Printf("Failed to write header to file. Err: %v\n", err)
					continue
				}
				file.Close()
			}

			fmt.Printf("Dumping to WAV file\n")
			file, err := os.OpenFile(fmt.Sprintf("output_%d.wav", counter), os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				fmt.Printf("Failed to open file. Err: %v\n", err)
				continue
			}

			_, err = file.Write(*br)
			file.Close()

			if err != nil {
				fmt.Printf("Failed to write to file. Err: %v\n", err)
				continue
			}

			lastBytesReceived = time.Now()
		}
	}()

	// open channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.openChan {
			fmt.Printf("\n\n[OpenResponse]\n\n")
		}
	}()

	// welcome response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.welcomeResponse {
			fmt.Printf("\n\n[WelcomeResponse]\n\n")
		}
	}()

	// conversation text response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ctr := range dch.conversationTextResponse {
			fmt.Printf("\n\n[ConversationTextResponse]\n")
			fmt.Printf("%s: %s\n\n", ctr.Role, ctr.Content)
		}
	}()

	// user started speaking response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.userStartedSpeakingResponse {
			fmt.Printf("\n\n[UserStartedSpeakingResponse]\n\n")
		}
	}()

	// agent thinking response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.agentThinkingResponse {
			fmt.Printf("\n\n[AgentThinkingResponse]\n\n")
		}
	}()

	// function call request response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.functionCallRequestResponse {
			fmt.Printf("\n\n[FunctionCallRequestResponse]\n\n")
		}
	}()

	// function calling response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.functionCallingResponse {
			fmt.Printf("\n\n[FunctionCallingResponse]\n\n")
		}
	}()

	// agent started speaking response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.agentStartedSpeakingResponse {
			fmt.Printf("\n\n[AgentStartedSpeakingResponse]\n\n")
		}
	}()

	// agent audio done response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.agentAudioDoneResponse {
			fmt.Printf("\n\n[AgentAudioDoneResponse]\n\n")
		}
	}()

	// end of thought response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.endOfThoughtResponse {
			fmt.Printf("\n\n[EndOfThoughtResponse]\n\n")
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
	tOptions := client.NewSettingsConfigurationOptions()
	tOptions.Agent.Think.Provider.Type = "open_ai"
	tOptions.Agent.Think.Model = "gpt-4o-mini"
	tOptions.Agent.Think.Instructions = "You are a helpful AI assistant."

	// implement your own callback
	var callback msginterfaces.AgentMessageChan
	callback = *NewMyHandler()

	// create a Deepgram client
	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, tOptions, callback)
	if err != nil {
		fmt.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// connect the websocket to Deepgram
	fmt.Printf("Starting Agent...\n")
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
	fmt.Printf("Starting Microphone...\n")
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
	fmt.Printf("Stopping Microphone...\n")
	err = mic.Stop()
	if err != nil {
		fmt.Printf("mic.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}

	// teardown library
	microphone.Teardown()

	// close DG client
	fmt.Printf("Stopping Agent...\n")
	dgClient.Stop()

	fmt.Printf("\n\nProgram exiting...\n")
}
