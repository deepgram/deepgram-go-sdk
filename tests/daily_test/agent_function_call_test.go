// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package daily_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type functionCallDailyHandler struct {
	functionCallRequests  chan *msginterfaces.FunctionCallRequestResponse
	functionCallCancelled chan *msginterfaces.FunctionCallCancelledResponse
	conversationText      chan *msginterfaces.ConversationTextResponse
	errors                chan *msginterfaces.ErrorResponse
	settingsApplied       chan *msginterfaces.SettingsAppliedResponse
}

func newFunctionCallDailyHandler() *functionCallDailyHandler {
	return &functionCallDailyHandler{
		functionCallRequests:  make(chan *msginterfaces.FunctionCallRequestResponse, 1),
		functionCallCancelled: make(chan *msginterfaces.FunctionCallCancelledResponse, 1),
		conversationText:      make(chan *msginterfaces.ConversationTextResponse, 4),
		errors:                make(chan *msginterfaces.ErrorResponse, 1),
		settingsApplied:       make(chan *msginterfaces.SettingsAppliedResponse, 1),
	}
}

func (h *functionCallDailyHandler) GetBinary() []*chan *[]byte { return nil }
func (h *functionCallDailyHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return nil
}
func (h *functionCallDailyHandler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return nil
}
func (h *functionCallDailyHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&h.conversationText}
}
func (h *functionCallDailyHandler) GetUserStartedSpeaking() []*chan *msginterfaces.UserStartedSpeakingResponse {
	return nil
}
func (h *functionCallDailyHandler) GetAgentThinking() []*chan *msginterfaces.AgentThinkingResponse {
	return nil
}
func (h *functionCallDailyHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&h.functionCallRequests}
}
func (h *functionCallDailyHandler) GetFunctionCallCancelled() []*chan *msginterfaces.FunctionCallCancelledResponse {
	return []*chan *msginterfaces.FunctionCallCancelledResponse{&h.functionCallCancelled}
}
func (h *functionCallDailyHandler) GetAgentStartedSpeaking() []*chan *msginterfaces.AgentStartedSpeakingResponse {
	return nil
}
func (h *functionCallDailyHandler) GetAgentAudioDone() []*chan *msginterfaces.AgentAudioDoneResponse {
	return nil
}
func (h *functionCallDailyHandler) GetClose() []*chan *msginterfaces.CloseResponse { return nil }
func (h *functionCallDailyHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}
func (h *functionCallDailyHandler) GetUnhandled() []*chan *[]byte { return nil }
func (h *functionCallDailyHandler) GetInjectionRefused() []*chan *msginterfaces.InjectionRefusedResponse {
	return nil
}
func (h *functionCallDailyHandler) GetKeepAlive() []*chan *msginterfaces.KeepAlive { return nil }
func (h *functionCallDailyHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsApplied}
}

func TestDaily_AgentFunctionCall(t *testing.T) {
	if os.Getenv("DEEPGRAM_API_KEY") == "" && os.Getenv("DEEPGRAM_ACCESS_TOKEN") == "" {
		t.Skip("DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN is required")
	}

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

	handler := newFunctionCallDailyHandler()
	conn, err := agent.NewWSUsingChan(ctx, "", &interfaces.ClientOptions{}, options, handler)
	if err != nil {
		t.Fatalf("create Agent WebSocket: %v", err)
	}
	defer conn.Stop()

	if !conn.Connect() {
		t.Fatal("connect Agent WebSocket")
	}

	select {
	case <-handler.settingsApplied:
	case err := <-handler.errors:
		t.Fatalf("apply settings: %s", err.Description)
	case <-ctx.Done():
		t.Fatalf("wait for SettingsApplied: %v", ctx.Err())
	}

	if err := conn.WriteJSON(msginterfaces.InjectUserMessage{
		Type:    msginterfaces.TypeInjectUserMessage,
		Content: "What is the weather in Fremont, California?",
	}); err != nil {
		t.Fatalf("send InjectUserMessage: %v", err)
	}

	responseSent := false
	for {
		select {
		case request := <-handler.functionCallRequests:
			if len(request.Functions) != 1 {
				t.Fatalf("expected one function call, got %d", len(request.Functions))
			}
			function := request.Functions[0]
			if function.Name != "get_weather" || !function.ClientSide {
				t.Fatalf("unexpected function call: %+v", function)
			}

			var arguments map[string]string
			if err := json.Unmarshal([]byte(function.Arguments), &arguments); err != nil {
				t.Fatalf("decode function arguments %q: %v", function.Arguments, err)
			}
			if arguments["location"] == "" {
				t.Fatalf("function call is missing location: %q", function.Arguments)
			}

			if err := conn.WriteJSON(msginterfaces.FunctionCallResponse{
				Type:    msginterfaces.TypeFunctionCallResponse,
				ID:      function.ID,
				Name:    function.Name,
				Content: `{"weather":"sunny"}`,
			}); err != nil {
				t.Fatalf("send FunctionCallResponse: %v", err)
			}
			responseSent = true
		case cancellation := <-handler.functionCallCancelled:
			t.Fatalf("function call was cancelled: %+v", cancellation.Functions)
		case conversation := <-handler.conversationText:
			if responseSent && conversation.Role == "assistant" {
				return
			}
		case err := <-handler.errors:
			t.Fatalf("Agent error: %s", err.Description)
		case <-ctx.Done():
			t.Fatalf("complete function-call round trip: %v", ctx.Err())
		}
	}
}
