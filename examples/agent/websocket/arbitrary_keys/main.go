// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

// Import dependencies
import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v2/pkg/api/agent/v1/websocket/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v2/pkg/client/agent"
	"github.com/deepgram/deepgram-go-sdk/v2/pkg/client/interfaces"
)

// MyHandler handles all websocket events
type MyHandler struct {
	binaryChan                   chan *[]byte
	openChan                     chan *msginterfaces.OpenResponse
	welcomeResponse              chan *msginterfaces.WelcomeResponse
	conversationTextResponse     chan *msginterfaces.ConversationTextResponse
	userStartedSpeakingResponse  chan *msginterfaces.UserStartedSpeakingResponse
	agentThinkingResponse        chan *msginterfaces.AgentThinkingResponse
	agentStartedSpeakingResponse chan *msginterfaces.AgentStartedSpeakingResponse
	agentAudioDoneResponse       chan *msginterfaces.AgentAudioDoneResponse
	closeChan                    chan *msginterfaces.CloseResponse
	errorChan                    chan *msginterfaces.ErrorResponse
	unhandledChan                chan *[]byte
	injectionRefusedResponse     chan *msginterfaces.InjectionRefusedResponse
	keepAliveResponse            chan *msginterfaces.KeepAlive
	settingsAppliedResponse      chan *msginterfaces.SettingsAppliedResponse
	functionCallRequestResponse  chan *msginterfaces.FunctionCallRequestResponse
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

func (dch MyHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&dch.functionCallRequestResponse}
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
		agentStartedSpeakingResponse: make(chan *msginterfaces.AgentStartedSpeakingResponse),
		agentAudioDoneResponse:       make(chan *msginterfaces.AgentAudioDoneResponse),
		closeChan:                    make(chan *msginterfaces.CloseResponse),
		errorChan:                    make(chan *msginterfaces.ErrorResponse),
		unhandledChan:                make(chan *[]byte),
		injectionRefusedResponse:     make(chan *msginterfaces.InjectionRefusedResponse),
		keepAliveResponse:            make(chan *msginterfaces.KeepAlive),
		settingsAppliedResponse:      make(chan *msginterfaces.SettingsAppliedResponse),
		functionCallRequestResponse:  make(chan *msginterfaces.FunctionCallRequestResponse),
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
		LogLevel: client.LogLevelVerbose,
	})

	// Set client options
	return &interfaces.ClientOptions{
		EnableKeepAlive: true,
	}
}

// Setup Event Handlers
func (dch MyHandler) Run() error {
	wgReceivers := sync.WaitGroup{}

	// Handle other events
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.openChan {
			fmt.Printf("\n\n[OpenResponse]\n\n")
		}
	}()

	// welcome channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.welcomeResponse {
			fmt.Printf("\n\n[WelcomeResponse]\n\n")
		}
	}()

	// settings applied channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.settingsAppliedResponse {
			fmt.Printf("\n\n[SettingsAppliedResponse]\n\n")
		}
	}()

	// close channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for closeResp := range dch.closeChan {
			fmt.Printf("\n\n[CloseResponse]\n")
			fmt.Printf("Close response received\n")
			fmt.Printf("Close response type: %+v\n", closeResp)
			fmt.Printf("\n")
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
			fmt.Printf("\n[UnhandledEvent]\n")
			fmt.Printf("Raw message: %s\n", string(*byData))
		}
	}()

	// Handle function call request
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for range dch.functionCallRequestResponse {
			fmt.Printf("\n\n[FunctionCallRequestResponse]\n\n")
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
	fmt.Printf("Program starting\n")
	// Print instructions
	fmt.Print("\n\nPress ENTER to exit!\n\n")

	// Initialize context
	ctx := context.Background()
	fmt.Printf("Context initialized\n")

	// Configure agent
	cOptions := configureAgent()
	fmt.Printf("Agent configured\n")

	// Set transcription options
	tOptions := client.NewSettingsConfigurationOptions()
	tOptions.Audio.Input.Encoding = "linear16"
	tOptions.Audio.Input.SampleRate = 48000
	tOptions.Agent.Think.Provider.Type = "open_ai"
	tOptions.Agent.Think.Provider.Model = "gpt-4o-mini"
	tOptions.Agent.Think.Prompt = "You are a helpful AI assistant."
	tOptions.Agent.Listen.Provider.Type = "deepgram"
	tOptions.Agent.Listen.Provider.Model = "nova-3"
	tOptions.Agent.Speak.Provider["ype"] = "deepgram"
	tOptions.Agent.Speak.Provider["model"] = "aura-2-thalia-en"
	tOptions.Agent.Speak.Provider["arbitrary_key"] = "test"
	tOptions.Agent.Language = "en"
	tOptions.Agent.Greeting = "Hello! How can I help you today?"
	fmt.Printf("Transcription options set\n")

	// Create handler
	fmt.Printf("Creating new Deepgram WebSocket client...\n")
	handler := NewMyHandler()
	if handler == nil {
		fmt.Printf("Failed to create handler\n")
		return
	}
	fmt.Printf("Handler created\n")
	defer handler.chatLogFile.Close()

	// Create client
	callback := msginterfaces.AgentMessageChan(*handler)
	fmt.Printf("Callback created\n")
	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, tOptions, callback)
	if err != nil {
		fmt.Printf("ERROR creating LiveTranscription connection:\n- Error: %v\n- Type: %T\n", err, err)
		return
	}
	fmt.Printf("Deepgram client created\n")

	// Connect to Deepgram
	fmt.Printf("Attempting to connect to Deepgram WebSocket...\n")
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Printf("WebSocket connection failed - check your API key and network connection\n")
		os.Exit(1)
	}
	fmt.Printf("Successfully connected to Deepgram WebSocket\n")
	fmt.Printf("Starting cleanup sequence...\n")
	fmt.Printf("Calling dgClient.Stop()\n")
	dgClient.Stop()
	fmt.Printf("dgClient.Stop() completed\n")
	fmt.Printf("\n\nProgram exiting...\n")
}
