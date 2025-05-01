// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// Import dependencies
import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/agent"
	"github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// MyHandler handles all websocket events
type MyHandler struct {
	binaryChan                   chan *[]byte
	openChan                     chan *msginterfaces.OpenResponse
	welcomeResponse              chan *msginterfaces.WelcomeResponse
	conversationTextResponse     chan *msginterfaces.ConversationTextResponse
	userStartedSpeakingResponse  chan *msginterfaces.UserStartedSpeakingResponse
	agentThinkingResponse        chan *msginterfaces.AgentThinkingResponse
	functionCallRequestResponse  chan *msginterfaces.FunctionCallRequestResponse
	agentStartedSpeakingResponse chan *msginterfaces.AgentStartedSpeakingResponse
	agentAudioDoneResponse       chan *msginterfaces.AgentAudioDoneResponse
	closeChan                    chan *msginterfaces.CloseResponse
	errorChan                    chan *msginterfaces.ErrorResponse
	unhandledChan                chan *[]byte
	injectionRefusedResponse     chan *msginterfaces.InjectionRefusedResponse
	keepAliveResponse            chan *msginterfaces.KeepAlive
	settingsAppliedResponse      chan *msginterfaces.SettingsAppliedResponse
	chatLogFile                  *os.File
}

// Channel getter methods to implement AgentMessageChan interface
func (dch MyHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&dch.binaryChan}
}

func (dch MyHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return []*chan *msginterfaces.OpenResponse{&dch.openChan}
}

func (dch MyHandler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return []*chan *msginterfaces.WelcomeResponse{&dch.welcomeResponse}
}

func (dch MyHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&dch.conversationTextResponse}
}

func (dch MyHandler) GetUserStartedSpeaking() []*chan *msginterfaces.UserStartedSpeakingResponse {
	return []*chan *msginterfaces.UserStartedSpeakingResponse{&dch.userStartedSpeakingResponse}
}

func (dch MyHandler) GetAgentThinking() []*chan *msginterfaces.AgentThinkingResponse {
	return []*chan *msginterfaces.AgentThinkingResponse{&dch.agentThinkingResponse}
}

func (dch MyHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&dch.functionCallRequestResponse}
}

func (dch MyHandler) GetAgentStartedSpeaking() []*chan *msginterfaces.AgentStartedSpeakingResponse {
	return []*chan *msginterfaces.AgentStartedSpeakingResponse{&dch.agentStartedSpeakingResponse}
}

func (dch MyHandler) GetAgentAudioDone() []*chan *msginterfaces.AgentAudioDoneResponse {
	return []*chan *msginterfaces.AgentAudioDoneResponse{&dch.agentAudioDoneResponse}
}

func (dch MyHandler) GetClose() []*chan *msginterfaces.CloseResponse {
	return []*chan *msginterfaces.CloseResponse{&dch.closeChan}
}

func (dch MyHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&dch.errorChan}
}

func (dch MyHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&dch.unhandledChan}
}

func (dch MyHandler) GetInjectionRefused() []*chan *msginterfaces.InjectionRefusedResponse {
	return []*chan *msginterfaces.InjectionRefusedResponse{&dch.injectionRefusedResponse}
}

func (dch MyHandler) GetKeepAlive() []*chan *msginterfaces.KeepAlive {
	return []*chan *msginterfaces.KeepAlive{&dch.keepAliveResponse}
}

func (dch MyHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&dch.settingsAppliedResponse}
}

// 2. Initialize the Voice Agent
func NewMyHandler() *MyHandler {
	// Create chat log file
	chatLogFile, err := os.OpenFile("chatlog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to create chat log file: %v\n", err)
		return nil
	}

	handler := &MyHandler{
		binaryChan:                   make(chan *[]byte),
		openChan:                     make(chan *msginterfaces.OpenResponse),
		welcomeResponse:              make(chan *msginterfaces.WelcomeResponse),
		conversationTextResponse:     make(chan *msginterfaces.ConversationTextResponse),
		userStartedSpeakingResponse:  make(chan *msginterfaces.UserStartedSpeakingResponse),
		agentThinkingResponse:        make(chan *msginterfaces.AgentThinkingResponse),
		functionCallRequestResponse:  make(chan *msginterfaces.FunctionCallRequestResponse),
		agentStartedSpeakingResponse: make(chan *msginterfaces.AgentStartedSpeakingResponse),
		agentAudioDoneResponse:       make(chan *msginterfaces.AgentAudioDoneResponse),
		closeChan:                    make(chan *msginterfaces.CloseResponse),
		errorChan:                    make(chan *msginterfaces.ErrorResponse),
		unhandledChan:                make(chan *[]byte),
		injectionRefusedResponse:     make(chan *msginterfaces.InjectionRefusedResponse),
		keepAliveResponse:            make(chan *msginterfaces.KeepAlive),
		settingsAppliedResponse:      make(chan *msginterfaces.SettingsAppliedResponse),
		chatLogFile:                  chatLogFile,
	}

	go func() {
		handler.Run()
	}()

	return handler
}

// 3. Configure the Agent
func configureAgent() *interfaces.ClientOptions {
	// Initialize library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelFull,
	})

	// Set client options
	return &interfaces.ClientOptions{
		EnableKeepAlive: true,
	}
}

// Setup Event Handlers
func (dch MyHandler) Run() error {
	wgReceivers := sync.WaitGroup{}

	// Handle binary data
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		counter := 0
		lastBytesReceived := time.Now().Add(-7 * time.Second)

		for br := range dch.binaryChan {
			fmt.Printf("\n\n[Binary Data Received]\n")
			fmt.Printf("Size: %d bytes\n", len(*br))

			if lastBytesReceived.Add(5 * time.Second).Before(time.Now()) {
				counter = counter + 1
				file, err := os.OpenFile(fmt.Sprintf("output_%d.wav", counter), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
				if err != nil {
					fmt.Printf("Failed to open file. Err: %v\n", err)
					continue
				}
				// Add WAV header
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

	// Handle conversation text
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ctr := range dch.conversationTextResponse {
			fmt.Printf("\n\n[ConversationTextResponse]\n")
			fmt.Printf("%s: %s\n\n", ctr.Role, ctr.Content)

			// Write to chat log
			if err := dch.writeToChatLog(ctr.Role, ctr.Content); err != nil {
				fmt.Printf("Failed to write to chat log: %v\n", err)
			}

			// If this is an agent response, we can consider the conversation turn complete
			if ctr.Role == "assistant" {
				fmt.Printf("Agent has responded, waiting for next user input...\n")
			}
		}
	}()

	// Handle keep alive responses
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for range dch.keepAliveResponse {
			fmt.Printf("\n\n[KeepAliveResponse]\n")
			// Write to chat log
			if err := dch.writeToChatLog("system", "Keep alive received"); err != nil {
				fmt.Printf("Failed to write to chat log: %v\n", err)
			}
		}
	}()

	// Handle agent audio done
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for range dch.agentAudioDoneResponse {
			fmt.Printf("\n\n[AgentAudioDoneResponse]\n")
			fmt.Printf("Audio processing completed\n")

			// Write to chat log
			if err := dch.writeToChatLog("system", "Audio processing completed"); err != nil {
				fmt.Printf("Failed to write to chat log: %v\n", err)
			}
		}
	}()

	// Handle other events
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.openChan {
			fmt.Printf("\n\n[OpenResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.welcomeResponse {
			fmt.Printf("\n\n[WelcomeResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.userStartedSpeakingResponse {
			fmt.Printf("\n\n[UserStartedSpeakingResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.agentThinkingResponse {
			fmt.Printf("\n\n[AgentThinkingResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.functionCallRequestResponse {
			fmt.Printf("\n\n[FunctionCallRequestResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.agentStartedSpeakingResponse {
			fmt.Printf("\n\n[AgentStartedSpeakingResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.settingsAppliedResponse {
			fmt.Printf("\n\n[SettingsAppliedResponse]\n\n")
		}
	}()

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.closeChan {
			fmt.Printf("\n\n[CloseResponse]\n\n")
		}
	}()

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

	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for byData := range dch.unhandledChan {
			fmt.Printf("\n[UnhandledEvent]\n")
			fmt.Printf("Raw message: %s\n", string(*byData))
		}
	}()

	// Wait for all receivers to finish
	wgReceivers.Wait()
	return nil
}

// Helper function to write to chat log
func (dch *MyHandler) writeToChatLog(role, content string) error {
	if dch.chatLogFile == nil {
		return fmt.Errorf("chat log file not initialized")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, role, content)

	_, err := dch.chatLogFile.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("failed to write to chat log: %v", err)
	}

	return nil
}

// Main function
func main() {
	// Print instructions
	fmt.Print("\n\nPress ENTER to exit!\n\n")

	// Initialize context
	ctx := context.Background()

	// Configure agent
	cOptions := configureAgent()

	// Set transcription options
	tOptions := client.NewSettingsConfigurationOptions()
	tOptions.Audio.Input.Encoding = "linear16"
	tOptions.Audio.Input.SampleRate = 24000
	tOptions.Agent.Think.Provider.Type = "open_ai"
	tOptions.Agent.Think.Provider.Model = "gpt-4o-mini"
	tOptions.Agent.Think.Prompt = "You are a helpful AI assistant."
	tOptions.Agent.Listen.Provider.Type = "deepgram"
	tOptions.Agent.Listen.Provider.Model = "nova-3"
	tOptions.Agent.Speak.Provider.Type = "deepgram"
	tOptions.Agent.Speak.Provider.Model = "aura-2-thalia-en"
	tOptions.Agent.Language = "en"
	tOptions.Agent.Greeting = "Hello! How can I help you today?"

	// Create handler
	fmt.Printf("Creating new Deepgram WebSocket client...\n")
	handler := NewMyHandler()
	if handler == nil {
		fmt.Printf("Failed to create handler\n")
		return
	}
	defer handler.chatLogFile.Close()

	// Create client
	callback := msginterfaces.AgentMessageChan(*handler)
	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, tOptions, callback)
	if err != nil {
		fmt.Printf("ERROR creating LiveTranscription connection:\n- Error: %v\n- Type: %T\n", err, err)
		return
	}

	// Connect to Deepgram
	fmt.Printf("Attempting to connect to Deepgram WebSocket...\n")
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Printf("WebSocket connection failed - check your API key and network connection\n")
		os.Exit(1)
	}
	fmt.Printf("Successfully connected to Deepgram WebSocket\n")

	// Stream audio from URL
	audioURL := "https://dpgr.am/spacewalk.wav"
	httpClient := new(http.Client)
	resp, err := httpClient.Get(audioURL)
	if err != nil {
		fmt.Printf("Failed to fetch audio from URL. Err: %v\n", err)
		return
	}
	fmt.Printf("Stream is up and running %s\n", reflect.TypeOf(resp))
	buf := bufio.NewReader(resp.Body)
	go func() {
		fmt.Printf("Starting to stream audio from URL...\n")
		defer resp.Body.Close()
		err = dgClient.Stream(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("Failed to stream audio. Err: %v\n", err)
			return
		}
		fmt.Printf("Finished streaming audio from URL\n")
	}()

	// Wait for user input to exit
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// Cleanup
	fmt.Printf("Stopping Agent...\n")
	dgClient.Stop()
	fmt.Printf("\n\nProgram exiting...\n")
}
