// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"
	"time"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	clientws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent/v1/websocket"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const functionCallRequestPayload = `{
	"type": "FunctionCallRequest",
	"functions": [{
		"id": "fc_12345678-90ab-cdef-1234-567890abcdef",
		"name": "get_weather",
		"arguments": "{\"location\":\"Fremont, CA 94539\"}",
		"client_side": true,
		"thought_signature": "abc123"
	}]
}`

const functionCallCancelledPayload = `{
	"type": "FunctionCallCancelled",
	"functions": [{
		"id": "fc_12345678-90ab-cdef-1234-567890abcdef",
		"name": "get_weather"
	}]
}`

const (
	functionCallID               = "fc_12345678-90ab-cdef-1234-567890abcdef"
	functionCallName             = "get_weather"
	functionCallThoughtSignature = "abc123"
)

type legacyFunctionCallHandler struct {
	functionCallRequests chan *msginterfaces.FunctionCallRequestResponse
	unhandled            chan *[]byte
}

func newLegacyFunctionCallHandler() *legacyFunctionCallHandler {
	return &legacyFunctionCallHandler{
		functionCallRequests: make(chan *msginterfaces.FunctionCallRequestResponse, 1),
		unhandled:            make(chan *[]byte, 1),
	}
}

type functionCallHandler struct {
	*legacyFunctionCallHandler
	functionCallCancelled chan *msginterfaces.FunctionCallCancelledResponse
}

func newFunctionCallHandler() *functionCallHandler {
	return &functionCallHandler{
		legacyFunctionCallHandler: newLegacyFunctionCallHandler(),
		functionCallCancelled:     make(chan *msginterfaces.FunctionCallCancelledResponse, 1),
	}
}

func (h *legacyFunctionCallHandler) GetBinary() []*chan *[]byte { return nil }
func (h *legacyFunctionCallHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetUserStartedSpeaking() []*chan *msginterfaces.UserStartedSpeakingResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetAgentThinking() []*chan *msginterfaces.AgentThinkingResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return []*chan *msginterfaces.FunctionCallRequestResponse{&h.functionCallRequests}
}
func (h *functionCallHandler) GetFunctionCallCancelled() []*chan *msginterfaces.FunctionCallCancelledResponse {
	return []*chan *msginterfaces.FunctionCallCancelledResponse{&h.functionCallCancelled}
}
func (h *legacyFunctionCallHandler) GetAgentStartedSpeaking() []*chan *msginterfaces.AgentStartedSpeakingResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetAgentAudioDone() []*chan *msginterfaces.AgentAudioDoneResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetClose() []*chan *msginterfaces.CloseResponse { return nil }
func (h *legacyFunctionCallHandler) GetError() []*chan *msginterfaces.ErrorResponse { return nil }
func (h *legacyFunctionCallHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandled}
}
func (h *legacyFunctionCallHandler) GetInjectionRefused() []*chan *msginterfaces.InjectionRefusedResponse {
	return nil
}
func (h *legacyFunctionCallHandler) GetKeepAlive() []*chan *msginterfaces.KeepAlive { return nil }
func (h *legacyFunctionCallHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return nil
}

func Test_FunctionCallRequestRouting(t *testing.T) {
	handler := newFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(functionCallRequestPayload)); err != nil {
		t.Fatalf("FunctionCallRequest must route successfully: %v", err)
	}

	select {
	case request := <-handler.functionCallRequests:
		if request.Type != msginterfaces.TypeFunctionCallRequestResponse {
			t.Errorf("expected type %q, got %q", msginterfaces.TypeFunctionCallRequestResponse, request.Type)
		}
		if len(request.Functions) != 1 {
			t.Fatalf("expected 1 function call, got %d", len(request.Functions))
		}

		function := request.Functions[0]
		if function.ID != functionCallID {
			t.Errorf("unexpected function ID: %q", function.ID)
		}
		if function.Name != functionCallName {
			t.Errorf("unexpected function name: %q", function.Name)
		}
		if function.Arguments != `{"location":"Fremont, CA 94539"}` {
			t.Errorf("unexpected function arguments: %q", function.Arguments)
		}
		if !function.ClientSide {
			t.Error("expected client-side function call")
		}
		if function.ThoughtSignature != functionCallThoughtSignature {
			t.Errorf("unexpected thought signature: %q", function.ThoughtSignature)
		}
	case <-time.After(time.Second):
		t.Fatal("FunctionCallRequest was not delivered to the handler")
	}
}

func Test_FunctionCallResponseMarshaling(t *testing.T) {
	response := clientws.FunctionCallResponse{
		Type:             msginterfaces.TypeFunctionCallResponse,
		ID:               functionCallID,
		Name:             functionCallName,
		Content:          `{"temperature_c":21}`,
		ThoughtSignature: functionCallThoughtSignature,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal FunctionCallResponse: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal FunctionCallResponse: %v", err)
	}
	if payload["type"] != response.Type || payload["id"] != response.ID || payload["name"] != response.Name || payload["content"] != response.Content || payload["thought_signature"] != response.ThoughtSignature {
		t.Errorf("unexpected FunctionCallResponse payload: %s", data)
	}
	if _, found := payload["function_call_id"]; found {
		t.Errorf("obsolete function_call_id field present in payload: %s", data)
	}
	if _, found := payload["output"]; found {
		t.Errorf("obsolete output field present in payload: %s", data)
	}

	emptyContent, err := json.Marshal(clientws.FunctionCallResponse{Name: functionCallName})
	if err != nil {
		t.Fatalf("marshal FunctionCallResponse with empty content: %v", err)
	}
	if err := json.Unmarshal(emptyContent, &payload); err != nil {
		t.Fatalf("unmarshal FunctionCallResponse with empty content: %v", err)
	}
	if _, found := payload["content"]; !found {
		t.Errorf("empty FunctionCallResponse content must be present: %s", emptyContent)
	}
}

func Test_FunctionCallCancelledRouting(t *testing.T) {
	handler := newFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(functionCallCancelledPayload)); err != nil {
		t.Fatalf("FunctionCallCancelled must route successfully: %v", err)
	}

	select {
	case cancellation := <-handler.functionCallCancelled:
		if cancellation.Type != msginterfaces.TypeFunctionCallCancelledResponse {
			t.Errorf("expected type %q, got %q", msginterfaces.TypeFunctionCallCancelledResponse, cancellation.Type)
		}
		if len(cancellation.Functions) != 1 {
			t.Fatalf("expected 1 canceled function call, got %d", len(cancellation.Functions))
		}
		if cancellation.Functions[0].ID != functionCallID {
			t.Errorf("unexpected canceled function ID: %q", cancellation.Functions[0].ID)
		}
		if cancellation.Functions[0].Name != functionCallName {
			t.Errorf("unexpected canceled function name: %q", cancellation.Functions[0].Name)
		}
	case <-time.After(time.Second):
		t.Fatal("FunctionCallCancelled was not delivered to the handler")
	}
}

func Test_FunctionCallCancelledFallsBackToUnhandled(t *testing.T) {
	handler := newLegacyFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(functionCallCancelledPayload)); err != nil {
		t.Fatalf("FunctionCallCancelled must fall back successfully: %v", err)
	}

	select {
	case message := <-handler.unhandled:
		if string(*message) != functionCallCancelledPayload {
			t.Errorf("unexpected unhandled message: %s", *message)
		}
	case <-time.After(time.Second):
		t.Fatal("FunctionCallCancelled was not delivered to the unhandled handler")
	}
}

func Test_FunctionDefinitionMarshaling(t *testing.T) {
	clientSideFunction := interfacesv1.Functions{
		Name:          functionCallName,
		Description:   "Get the current weather for a location.",
		DeferUntilEOT: true,
		Parameters: interfacesv1.Parameters{
			Type: "object",
			PropertySchemas: map[string]interface{}{
				"location": map[string]string{
					"type":        "string",
					"description": "The city or location to get weather for.",
				},
			},
			Required: []string{"location"},
		},
	}

	data, err := json.Marshal(clientSideFunction)
	if err != nil {
		t.Fatalf("marshal function definition: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal function definition: %v", err)
	}
	parameters, ok := payload["parameters"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing function parameters: %s", data)
	}
	properties, ok := parameters["properties"].(map[string]interface{})
	if !ok || properties["location"] == nil {
		t.Fatalf("missing location property: %s", data)
	}
	if _, found := payload["endpoint"]; found {
		t.Errorf("client-side function must omit endpoint: %s", data)
	}
	if payload["defer_until_eot"] != true {
		t.Errorf("expected defer_until_eot=true, got %v", payload["defer_until_eot"])
	}

	serverSideFunction := interfacesv1.Functions{
		Name: functionCallName,
		Endpoint: interfacesv1.Endpoint{
			Url:    "https://example.com/weather",
			Method: "POST",
		},
	}
	data, err = json.Marshal(serverSideFunction)
	if err != nil {
		t.Fatalf("marshal server-side function definition: %v", err)
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal server-side function definition: %v", err)
	}
	endpoint, ok := payload["endpoint"].(map[string]interface{})
	if !ok || endpoint["url"] != "https://example.com/weather" || endpoint["method"] != "POST" {
		t.Errorf("unexpected server-side endpoint: %s", data)
	}
}

func Test_LegacyFunctionPropertiesMarshaling(t *testing.T) {
	data, err := json.Marshal(interfacesv1.Parameters{
		Type: "object",
		Properties: interfacesv1.Properties{
			Item: interfacesv1.Item{
				Type:        "string",
				Description: "The city or location to get weather for.",
			},
		},
		Required: []string{"location"},
	})
	if err != nil {
		t.Fatalf("marshal legacy function parameters: %v", err)
	}

	const expected = `{"type":"object","properties":{"item":{"type":"string","description":"The city or location to get weather for."}},"required":["location"]}`
	if string(data) != expected {
		t.Errorf("legacy parameters changed:\nexpected: %s\nactual:   %s", expected, data)
	}
}

func Test_FunctionCallClientSideFalseSurvivesMarshal(t *testing.T) {
	data, err := json.Marshal(msginterfaces.FunctionCall{ClientSide: false})
	if err != nil {
		t.Fatalf("marshal FunctionCall: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal FunctionCall: %v", err)
	}
	if payload["client_side"] != false {
		t.Errorf("expected client_side=false, got %v", payload["client_side"])
	}
}
