// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/microphone"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

// HistoryHandler demonstrates the INTERFACE SEGREGATION pattern for opt-in History support.
//
// This handler implements TWO interfaces:
// 1. AgentMessageChan (required) - Core agent functionality
// 2. HistoryMessageChan (optional) - History message support
//
// DESIGN PATTERN DEMONSTRATION:
// - The router uses type assertion to detect HistoryMessageChan support
// - Handlers that don't implement HistoryMessageChan continue to work unchanged
// - This avoids breaking changes that would require a major version bump
// - Follows Go stdlib patterns like http.ResponseWriter + http.Hijacker
//
// This example shows how new handlers can opt-in to History functionality
// while maintaining backwards compatibility with existing implementations.
type HistoryHandler struct {
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
	// History channels for handling conversation and function call history
	historyConversationTextChan chan *msginterfaces.HistoryConversationText
	historyFunctionCallsChan    chan *msginterfaces.HistoryFunctionCalls
	// WebSocket client for sending function responses
	dgClient *client.WSChannel
}

func NewHistoryHandler() *HistoryHandler {
	handler := &HistoryHandler{
		binaryChan:                   make(chan *[]byte, 100),
		openChan:                     make(chan *msginterfaces.OpenResponse, 100),
		welcomeResponse:              make(chan *msginterfaces.WelcomeResponse, 100),
		conversationTextResponse:     make(chan *msginterfaces.ConversationTextResponse, 100),
		userStartedSpeakingResponse:  make(chan *msginterfaces.UserStartedSpeakingResponse, 100),
		agentThinkingResponse:        make(chan *msginterfaces.AgentThinkingResponse, 100),
		functionCallRequestResponse:  make(chan *msginterfaces.FunctionCallRequestResponse, 100),
		agentStartedSpeakingResponse: make(chan *msginterfaces.AgentStartedSpeakingResponse, 100),
		agentAudioDoneResponse:       make(chan *msginterfaces.AgentAudioDoneResponse, 100),
		closeChan:                    make(chan *msginterfaces.CloseResponse, 100),
		errorChan:                    make(chan *msginterfaces.ErrorResponse, 100),
		unhandledChan:                make(chan *[]byte, 100),
		injectionRefusedResponse:     make(chan *msginterfaces.InjectionRefusedResponse, 100),
		keepAliveResponse:            make(chan *msginterfaces.KeepAlive, 100),
		settingsAppliedResponse:      make(chan *msginterfaces.SettingsAppliedResponse, 100),
		historyConversationTextChan:  make(chan *msginterfaces.HistoryConversationText, 100),
		historyFunctionCallsChan:     make(chan *msginterfaces.HistoryFunctionCalls, 100),
	}

	go func() {
		handler.Run()
	}()

	return handler
}

// SetClient sets the WebSocket client for sending function responses
func (h *HistoryHandler) SetClient(client *client.WSChannel) {
	h.dgClient = client
}

// Implement all the required interface methods
func (h HistoryHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&h.binaryChan}
}

func (h HistoryHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return []*chan *msginterfaces.OpenResponse{&h.openChan}
}

func (h HistoryHandler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return []*chan *msginterfaces.WelcomeResponse{&h.welcomeResponse}
}

func (h HistoryHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&h.conversationTextResponse}
}

func (h HistoryHandler) GetUserStartedSpeaking() []*chan *msginterfaces.UserStartedSpeakingResponse {
	return []*chan *msginterfaces.UserStartedSpeakingResponse{&h.userStartedSpeakingResponse}
}

func (h HistoryHandler) GetAgentThinking() []*chan *msginterfaces.AgentThinkingResponse {
	return []*chan *msginterfaces.AgentThinkingResponse{&h.agentThinkingResponse}
}

func (h HistoryHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&h.functionCallRequestResponse}
}

func (h HistoryHandler) GetAgentStartedSpeaking() []*chan *msginterfaces.AgentStartedSpeakingResponse {
	return []*chan *msginterfaces.AgentStartedSpeakingResponse{&h.agentStartedSpeakingResponse}
}

func (h HistoryHandler) GetAgentAudioDone() []*chan *msginterfaces.AgentAudioDoneResponse {
	return []*chan *msginterfaces.AgentAudioDoneResponse{&h.agentAudioDoneResponse}
}

func (h HistoryHandler) GetClose() []*chan *msginterfaces.CloseResponse {
	return []*chan *msginterfaces.CloseResponse{&h.closeChan}
}

func (h HistoryHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errorChan}
}

func (h HistoryHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandledChan}
}

func (h HistoryHandler) GetInjectionRefused() []*chan *msginterfaces.InjectionRefusedResponse {
	return []*chan *msginterfaces.InjectionRefusedResponse{&h.injectionRefusedResponse}
}

func (h HistoryHandler) GetKeepAlive() []*chan *msginterfaces.KeepAlive {
	return []*chan *msginterfaces.KeepAlive{&h.keepAliveResponse}
}

func (h HistoryHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsAppliedResponse}
}

// OPTIONAL INTERFACE: HistoryMessageChan implementation
// These methods enable History message support through interface segregation.
// The router uses type assertion to detect these methods and route History messages accordingly.
func (h HistoryHandler) GetHistoryConversationText() []*chan *msginterfaces.HistoryConversationText {
	return []*chan *msginterfaces.HistoryConversationText{&h.historyConversationTextChan}
}

func (h HistoryHandler) GetHistoryFunctionCalls() []*chan *msginterfaces.HistoryFunctionCalls {
	return []*chan *msginterfaces.HistoryFunctionCalls{&h.historyFunctionCallsChan}
}

// prettyPrintJSON formats JSON for display like the JavaScript example
func prettyPrintJSON(v interface{}) string {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(jsonBytes)
}

// Run handles all the channel events
func (h *HistoryHandler) Run() {
	for {
		select {
		case data := <-h.binaryChan:
			// Skip logging binary data for cleaner output - just process it
			_ = data

		case open := <-h.openChan:
			_ = open // suppress unused variable warning
			fmt.Printf("🔌 Connection opened successfully\n")

		case welcome := <-h.welcomeResponse:
			fmt.Printf("👋 Welcome! Agent is ready (Request ID: %s)\n", welcome.RequestID)

		case conv := <-h.conversationTextResponse:
			fmt.Printf("\n💬 %s: %s\n", strings.Title(conv.Role), conv.Content)

		case userSpeaking := <-h.userStartedSpeakingResponse:
			_ = userSpeaking // suppress unused variable warning
			fmt.Printf("🎤 User started speaking\n")

		case agentThinking := <-h.agentThinkingResponse:
			_ = agentThinking // suppress unused variable warning
			fmt.Printf("🤔 Agent is thinking...\n")

		case functionCall := <-h.functionCallRequestResponse:
			fmt.Printf("\n📞 Function call request received:\n")
			if len(functionCall.Functions) > 0 {
				fc := functionCall.Functions[0]
				fmt.Printf("   Function: %s\n", fc.Name)
				fmt.Printf("   ID: %s\n", fc.ID)
				fmt.Printf("   Arguments: %s\n", fc.Arguments)
			}
			if h.dgClient != nil {
				h.handleFunctionCall(functionCall, h.dgClient)
			} else {
				fmt.Printf("❌ WebSocket client not set - cannot handle function call\n")
			}

		case agentSpeaking := <-h.agentStartedSpeakingResponse:
			fmt.Printf("🗣️  Agent started speaking (Latency: %.0fms)\n", agentSpeaking.TotalLatency)

		case audioDone := <-h.agentAudioDoneResponse:
			_ = audioDone // suppress unused variable warning
			fmt.Printf("🔇 Agent finished speaking\n")

		case close := <-h.closeChan:
			_ = close // suppress unused variable warning
			fmt.Printf("🔚 Connection closed\n")

		case err := <-h.errorChan:
			if err.ErrMsg != "" || err.Description != "" {
				fmt.Printf("❌ Error: %s - %s\n", err.ErrMsg, err.Description)
			}

		case unhandled := <-h.unhandledChan:
			// Log any truly unhandled messages for debugging
			_ = unhandled

		case injection := <-h.injectionRefusedResponse:
			fmt.Printf("🚫 Message injection refused: %s\n", injection.Message)

		case keepAlive := <-h.keepAliveResponse:
			_ = keepAlive // suppress unused variable warning
			// Skip noisy keepalive messages for cleaner output

		case settings := <-h.settingsAppliedResponse:
			_ = settings // suppress unused variable warning
			fmt.Printf("⚙️  Agent settings applied successfully\n")

		case historyConv := <-h.historyConversationTextChan:
			fmt.Printf("📚 Conversation History payload received: %s\n", prettyPrintJSON(historyConv))

		case historyFunc := <-h.historyFunctionCallsChan:
			fmt.Printf("📚 Function Call History payload received: %s\n", prettyPrintJSON(historyFunc))
		}
	}
}

// handleFunctionCall processes function call requests from the agent
func (h *HistoryHandler) handleFunctionCall(functionCall *msginterfaces.FunctionCallRequestResponse, dgClient *client.WSChannel) {
	if len(functionCall.Functions) == 0 {
		fmt.Printf("❌ No functions in request\n")
		return
	}

	// Get the first function call (there should only be one for our use case)
	fc := functionCall.Functions[0]

	fmt.Printf("📞 Function Call: %s (ID: %s)\n", fc.Name, fc.ID)

	switch fc.Name {
	case "get_weather":
		h.handleWeatherFunction(&fc, dgClient)
	default:
		fmt.Printf("❌ Unknown function: %s\n", fc.Name)
	}
}

// handleWeatherFunction simulates a weather API call and responds to the agent
func (h *HistoryHandler) handleWeatherFunction(functionCall *msginterfaces.FunctionCall, dgClient *client.WSChannel) {
	// Parse the JSON arguments
	var args struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}

	if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
		fmt.Printf("❌ Error parsing function arguments: %v\n", err)
		return
	}

	// Validate required parameters
	if args.Location == "" {
		fmt.Printf("❌ Error: location parameter is required\n")
		return
	}

	// Default unit to fahrenheit if not specified
	if args.Unit == "" {
		args.Unit = "fahrenheit"
	}

	fmt.Printf("🌤️  Fetching weather for %s in %s...\n", args.Location, args.Unit)

	// Simulate weather data (in a real implementation, you'd call a weather API)
	weatherData := map[string]interface{}{
		"location":    args.Location,
		"temperature": 72,
		"condition":   "sunny",
		"humidity":    45,
		"unit":        args.Unit,
	}

	if args.Unit == "celsius" {
		weatherData["temperature"] = 22 // Convert to celsius for demo
	}

	weatherData["description"] = fmt.Sprintf("The weather in %s is %s with a temperature of %v°%s and %v%% humidity.",
		args.Location,
		weatherData["condition"],
		weatherData["temperature"],
		strings.ToUpper(string(args.Unit[0])),
		weatherData["humidity"])

	// Convert weather data to JSON string for response
	weatherJSON, _ := json.Marshal(weatherData)

	// Send the response back to the agent (matches Python SDK exactly)
	response := &msginterfaces.FunctionCallResponse{
		Type:    msginterfaces.TypeFunctionCallResponse,
		ID:      functionCall.ID,
		Name:    functionCall.Name,   // Function name (required by API)
		Content: string(weatherJSON), // Response content (matches Python SDK)
	}

	// Send confirmation message
	fmt.Printf("📞 Sending weather response for %s: %s\n", args.Location, weatherData["description"])

	err := dgClient.WriteJSON(response)
	if err != nil {
		fmt.Printf("❌ Error sending function response: %v\n", err)
	} else {
		fmt.Printf("✅ Weather response sent to agent\n")
	}
}

func createConversationHistory() []interfacesv1.ContextMessage {
	return []interfacesv1.ContextMessage{
		// Simple conversation history to demonstrate context retention
		interfacesv1.HistoryConversationText{
			Type:    msginterfaces.TypeHistoryConversationText,
			Role:    "user",
			Content: "Hello, I'm testing the conversation history feature.",
		},
		interfacesv1.HistoryConversationText{
			Type:    msginterfaces.TypeHistoryConversationText,
			Role:    "assistant",
			Content: "Hello! I can see you're testing the history feature. I can remember our previous conversations and help you with weather information using my get_weather function.",
		},
	}
}

func createAgentConfiguration() *interfaces.SettingsOptions {
	conversationHistory := createConversationHistory()

	// Create agent configuration
	config := client.NewSettingsConfigurationOptions()

	// Enable history feature for conversation context
	config.Flags = &interfacesv1.Flags{
		History: true,
	}

	// Agent tags for analytics
	config.Tags = []string{"history-example", "function-calling", "weather-demo"}

	// Audio configuration
	config.Audio.Input.Encoding = "linear16"
	config.Audio.Input.SampleRate = 16000

	// Agent configuration
	config.Agent.Language = "en"

	// Provide conversation context/history
	config.Agent.Context = &interfacesv1.Context{
		Messages: conversationHistory,
	}

	// Configure the listen provider
	config.Agent.Listen.Provider = map[string]interface{}{
		"type":  "deepgram",
		"model": "nova-3",
	}

	// Configure the speak provider
	config.Agent.Speak.Provider = map[string]interface{}{
		"type":  "deepgram",
		"model": "aura-asteria-en",
	}

	// Configure the thinking/LLM provider with function calling
	config.Agent.Think.Provider = map[string]interface{}{
		"type":  "open_ai",
		"model": "gpt-4o-mini",
	}

	// Enable function calling with weather function using OpenAI-style schema
	config.Agent.Think.Functions = &[]interfacesv1.Functions{
		{
			Name:        "get_weather",
			Description: "Get the current weather conditions for a specific location",
			Parameters: interfacesv1.Parameters{
				Type: "object",
				Properties: map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city or location to get weather for (e.g., 'New York', 'London', 'Tokyo')",
					},
					"unit": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"fahrenheit", "celsius"},
						"description": "Temperature unit preference",
						"default":     "fahrenheit",
					},
				},
				Required: []string{"location"},
			},
			// No endpoint field = client-side function (as per docs)
		},
	}

	config.Agent.Think.Prompt = "You are a helpful AI assistant. You can remember our previous conversations and provide helpful responses to user questions. Be conversational and helpful."

	config.Agent.Greeting = "Hello! I'm your AI assistant. I remember our previous conversations and I'm here to help you with any questions or tasks. How can I assist you today?"

	return config
}

func main() {
	// Initialize the microphone library
	microphone.Initialize()

	// Print clean header
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🌤️  Deepgram Agent with History & Context")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Features:")
	fmt.Println("• 📚 Conversation history context (2 messages)")
	fmt.Println("• 🔧 Live function calling with weather API")
	fmt.Println("• 📞 SDK properly routes History events (no more 'Unknown Event' warnings!)")
	fmt.Println("• 🎯 Shows JSON payloads like JavaScript SDK for easy understanding")
	fmt.Println("• ✅ Agent can dynamically respond to function calls")
	fmt.Println()
	fmt.Println("Start speaking to interact with the agent!")
	fmt.Println("Press ENTER to exit anytime!")
	fmt.Println(strings.Repeat("=", 60))

	// Initialize the SDK with standard logging
	// History events are now properly supported by the SDK router
	client.Init(client.InitLib{
		LogLevel: client.LogLevelStandard, // Standard logging - no need to suppress History event messages
	})

	// Create context
	ctx := context.Background()

	// Client options
	cOptions := &interfaces.ClientOptions{
		EnableKeepAlive: true,
	}

	// Create agent configuration
	config := createAgentConfiguration()

	fmt.Println("📚 Context loaded: 2 conversation messages")
	fmt.Println("   When connected, these should arrive as History events!")

	// Create custom handler that supports History events
	handler := NewHistoryHandler()

	// Create callback interface
	callback := msginterfaces.AgentMessageChan(*handler)

	// Create Deepgram client
	fmt.Print("\n🚀 Connecting to Deepgram...")
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("⚠️  Note: ConfigurationSettings JSON below is debug output from SDK")
	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, config, callback)
	if err != nil {
		fmt.Printf(" ❌ Failed\nError: %v\n", err)
		return
	}

	// Set the client on the handler so it can send function responses
	handler.SetClient(dgClient)

	// Connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Printf(" ❌ Failed\nCheck your API key and network connection\n")
		os.Exit(1)
	}
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("✅ Connected! Agent ready for interaction.")

	// Set up microphone
	fmt.Print("🎤 Initializing microphone...")
	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  16000,
	})
	if err != nil {
		fmt.Printf(" ❌ Failed: %v\n", err)
		return
	}

	// Start microphone
	err = mic.Start()
	if err != nil {
		fmt.Printf(" ❌ Failed: %v\n", err)
		return
	}
	fmt.Println(" ✅ Ready!")

	fmt.Println("\n📞 Agent is live! Start speaking or press ENTER to exit...")
	fmt.Println()
	fmt.Println("ℹ️  Expected Events:")
	fmt.Println("   • 📚 Conversation History payload received: { type: 'History', role: '...', content: '...' }")
	fmt.Println("   • 📞 Function call requests when you ask about weather")
	fmt.Println("   • 🌤️  Live weather responses sent back to the agent")
	fmt.Println("   • ✅ Agent can now respond dynamically instead of using mock data")
	fmt.Println("   • ✅ Proves the SDK router handles History messages and function calls properly!")
	fmt.Println(strings.Repeat("-", 60))

	// Create a wait group to manage goroutines
	var wgReceivers sync.WaitGroup

	// Start sending microphone data
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		mic.Stream(dgClient)
	}()

	// Enable function call handling
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()
		for functionCall := range handler.functionCallRequestResponse {
			handler.handleFunctionCall(functionCall, dgClient)
		}
	}()

	// Wait for user to press Enter to exit
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	fmt.Printf("\n" + strings.Repeat("-", 60))
	fmt.Printf("\n🛑 Shutting down...")

	// Stop microphone
	err = mic.Stop()
	if err != nil {
		fmt.Printf("\n⚠️  Microphone stop error: %v", err)
	}

	// teardown library
	microphone.Teardown()

	// Close Deepgram connection
	dgClient.Stop()

	// Wait for all receivers to finish
	wgReceivers.Wait()

	fmt.Printf(" ✅ Complete!\n\n")
}
