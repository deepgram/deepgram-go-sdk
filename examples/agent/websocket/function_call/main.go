// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	errors               chan *msginterfaces.ErrorResponse
	settingsApplied      chan *msginterfaces.SettingsAppliedResponse
}

func newFunctionCallHandler() *functionCallHandler {
	return &functionCallHandler{
		DefaultChanHandler:   agentws.NewDefaultChanHandler(),
		functionCallRequests: make(chan *msginterfaces.FunctionCallRequestResponse, 1),
		errors:               make(chan *msginterfaces.ErrorResponse, 1),
		settingsApplied:      make(chan *msginterfaces.SettingsAppliedResponse, 1),
	}
}

func (h *functionCallHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&h.functionCallRequests}
}

func (h *functionCallHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}

func (h *functionCallHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsApplied}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	functions := []interfacesv1.Functions{
		{
			Name:        "get_weather",
			Description: "Return a fixed weather result for the requested location.",
			Parameters: interfacesv1.Parameters{
				Type: "object",
				Properties: map[string]interface{}{
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
		log.Fatal(err)
	}
	defer conn.Stop()

	if !conn.Connect() {
		log.Fatal("connect Agent WebSocket")
	}

	select {
	case <-handler.settingsApplied:
	case err := <-handler.errors:
		log.Fatal(err.Description)
	case <-ctx.Done():
		log.Fatal("wait for SettingsApplied")
	}

	if err := conn.WriteJSON(msginterfaces.InjectUserMessage{
		Type:    msginterfaces.TypeInjectUserMessage,
		Content: "What is the weather in Fremont, California?",
	}); err != nil {
		log.Fatal(err)
	}

	select {
	case request := <-handler.functionCallRequests:
		for _, function := range request.Functions {
			var arguments map[string]string
			if err := json.Unmarshal([]byte(function.Arguments), &arguments); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Calling %s with %v\n", function.Name, arguments)

			if err := conn.WriteJSON(msginterfaces.FunctionCallResponse{
				Type:    msginterfaces.TypeFunctionCallResponse,
				ID:      function.ID,
				Name:    function.Name,
				Content: `{"weather":"sunny"}`,
			}); err != nil {
				log.Fatal(err)
			}
		}
	case <-ctx.Done():
		log.Fatal(ctx.Err())
	}
}
