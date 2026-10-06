// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type functionCallHandler struct {
	*agentws.DefaultChanHandler
	functionCallRequests chan *msginterfaces.FunctionCallRequestResponse
	functionCallCanceled chan *msginterfaces.FunctionCallCancelledResponse
	conversationText     chan *msginterfaces.ConversationTextResponse
	errors               chan *msginterfaces.ErrorResponse
	settingsApplied      chan *msginterfaces.SettingsAppliedResponse
}

func newFunctionCallHandler() *functionCallHandler {
	return &functionCallHandler{
		DefaultChanHandler:   agentws.NewDefaultChanHandler(),
		functionCallRequests: make(chan *msginterfaces.FunctionCallRequestResponse, 1),
		functionCallCanceled: make(chan *msginterfaces.FunctionCallCancelledResponse, 1),
		conversationText:     make(chan *msginterfaces.ConversationTextResponse, 4),
		errors:               make(chan *msginterfaces.ErrorResponse, 1),
		settingsApplied:      make(chan *msginterfaces.SettingsAppliedResponse, 1),
	}
}

type functionCallResult struct {
	function msginterfaces.FunctionCall
	content  string
}

type functionCallError struct {
	id  string
	err error
}

func (h *functionCallHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&h.conversationText}
}

func (h *functionCallHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&h.functionCallRequests}
}

func (h *functionCallHandler) GetFunctionCallCancelled() []*chan *msginterfaces.FunctionCallCancelledResponse {
	return []*chan *msginterfaces.FunctionCallCancelledResponse{&h.functionCallCanceled}
}

func (h *functionCallHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}

func (h *functionCallHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsApplied}
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func executeFunction(ctx context.Context, function msginterfaces.FunctionCall, results chan<- functionCallResult, failures chan<- functionCallError) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	var arguments map[string]string
	if err := json.Unmarshal([]byte(function.Arguments), &arguments); err != nil {
		select {
		case failures <- functionCallError{id: function.ID, err: err}:
		case <-ctx.Done():
		}
		return
	}
	fmt.Printf("Calling %s with %v\n", function.Name, arguments)

	select {
	case results <- functionCallResult{function: function, content: `{"weather":"sunny"}`}:
	case <-ctx.Done():
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	functions := []interfacesv1.Functions{
		{
			Name:        "get_weather",
			Description: "Return a fixed weather result for the requested location.",
			Parameters: interfacesv1.Parameters{
				Type: "object",
				PropertySchemas: map[string]interface{}{
					"location": map[string]string{
						"type": "string",
					},
				},
				Required: []string{"location"},
			},
		},
	}
	options := agent.NewSettingsConfigurationOptions()
	options.Agent.Think.Prompt = "Always call get_weather when asked about weather."
	options.Agent.Think.Functions = &functions
	options.Agent.Greeting = ""

	handler := newFunctionCallHandler()
	conn, err := agent.NewWSUsingChan(ctx, "", &interfaces.ClientOptions{}, options, handler)
	if err != nil {
		return fmt.Errorf("create Agent WebSocket: %w", err)
	}
	defer conn.Stop()

	if !conn.Connect() {
		return fmt.Errorf("connect Agent WebSocket")
	}

	select {
	case <-handler.settingsApplied:
	case err := <-handler.errors:
		return fmt.Errorf("apply settings: %s", err.Description)
	case <-ctx.Done():
		return fmt.Errorf("wait for SettingsApplied: %w", ctx.Err())
	}

	if err := conn.WriteJSON(msginterfaces.InjectUserMessage{
		Type:    msginterfaces.TypeInjectUserMessage,
		Content: "What is the weather in Fremont, California?",
	}); err != nil {
		return fmt.Errorf("send InjectUserMessage: %w", err)
	}

	responseSent := false
	results := make(chan functionCallResult)
	failures := make(chan functionCallError)
	inFlight := make(map[string]context.CancelFunc)
	defer func() {
		for _, cancel := range inFlight {
			cancel()
		}
	}()

	for {
		select {
		case request := <-handler.functionCallRequests:
			for _, function := range request.Functions {
				if !function.ClientSide {
					continue
				}

				if cancel, ok := inFlight[function.ID]; ok {
					cancel()
				}
				functionCtx, cancel := context.WithCancel(ctx)
				inFlight[function.ID] = cancel
				go executeFunction(functionCtx, function, results, failures)
			}
		case cancellation := <-handler.functionCallCanceled:
			for _, function := range cancellation.Functions {
				if cancel, ok := inFlight[function.ID]; ok {
					cancel()
					delete(inFlight, function.ID)
				}
			}
		case result := <-results:
			cancel, ok := inFlight[result.function.ID]
			if !ok {
				continue
			}
			cancel()
			delete(inFlight, result.function.ID)

			if err := conn.WriteJSON(msginterfaces.FunctionCallResponse{
				Type:             msginterfaces.TypeFunctionCallResponse,
				ID:               result.function.ID,
				Name:             result.function.Name,
				Content:          result.content,
				ThoughtSignature: result.function.ThoughtSignature,
			}); err != nil {
				return fmt.Errorf("send FunctionCallResponse: %w", err)
			}
			responseSent = true
		case failure := <-failures:
			if cancel, ok := inFlight[failure.id]; ok {
				cancel()
				delete(inFlight, failure.id)
				return fmt.Errorf("decode function arguments for %q: %w", failure.id, failure.err)
			}
		case conversation := <-handler.conversationText:
			if responseSent && conversation.Role == "assistant" {
				fmt.Printf("Agent: %s\n", conversation.Content)
				return nil
			}
		case err := <-handler.errors:
			return fmt.Errorf("agent error: %s", err.Description)
		case <-ctx.Done():
			return fmt.Errorf("complete function call: %w", ctx.Err())
		}
	}
}
