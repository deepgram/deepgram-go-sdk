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
)

const customFromThinkProviderPayload = `{"type":"__customFromThinkProvider","content":{"event":"ready","turn":1}}`

var _ msginterfaces.CustomFromThinkProviderChan = (*agentws.DefaultChanHandler)(nil)

type customThinkProviderHandler struct {
	*legacyFunctionCallHandler
	customFromThinkProvider chan *msginterfaces.CustomFromThinkProviderResponse
}

func newCustomThinkProviderHandler() *customThinkProviderHandler {
	return &customThinkProviderHandler{
		legacyFunctionCallHandler: newLegacyFunctionCallHandler(),
		customFromThinkProvider:   make(chan *msginterfaces.CustomFromThinkProviderResponse, 1),
	}
}

func (h *customThinkProviderHandler) GetCustomFromThinkProvider() []*chan *msginterfaces.CustomFromThinkProviderResponse {
	return []*chan *msginterfaces.CustomFromThinkProviderResponse{&h.customFromThinkProvider}
}

func Test_CustomToThinkProviderMarshaling(t *testing.T) {
	request := clientws.CustomToThinkProvider{
		Type:    msginterfaces.TypeCustomToThinkProvider,
		Content: json.RawMessage(`{"event":"start","values":[1,true,null]}`),
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal CustomToThinkProvider: %v", err)
	}

	const expected = `{"type":"__customToThinkProvider","content":{"event":"start","values":[1,true,null]}}`
	if string(data) != expected {
		t.Errorf("unexpected CustomToThinkProvider payload:\nexpected: %s\nactual:   %s", expected, data)
	}

	data, err = json.Marshal(msginterfaces.CustomToThinkProvider{
		Type:    msginterfaces.TypeCustomToThinkProvider,
		Content: nil,
	})
	if err != nil {
		t.Fatalf("marshal CustomToThinkProvider null content: %v", err)
	}
	if string(data) != `{"type":"__customToThinkProvider","content":null}` {
		t.Errorf("null content must be preserved as a valid JSON value: %s", data)
	}
}

func Test_CustomToThinkProviderValidation(t *testing.T) {
	conn := &clientws.WSChannel{}
	if err := conn.SendCustomToThinkProvider(json.RawMessage(`not-json`)); err == nil {
		t.Fatal("invalid custom Think provider content must be rejected before writing")
	}
}

func Test_CustomFromThinkProviderRouting(t *testing.T) {
	handler := newCustomThinkProviderHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(customFromThinkProviderPayload)); err != nil {
		t.Fatalf("CustomFromThinkProvider must route successfully: %v", err)
	}

	select {
	case response := <-handler.customFromThinkProvider:
		if response.Type != msginterfaces.TypeCustomFromThinkProviderResponse {
			t.Errorf("expected type %q, got %q", msginterfaces.TypeCustomFromThinkProviderResponse, response.Type)
		}
		if string(response.Content) != `{"event":"ready","turn":1}` {
			t.Errorf("unexpected custom provider content: %s", response.Content)
		}
	case <-time.After(time.Second):
		t.Fatal("CustomFromThinkProvider was not delivered to the handler")
	}
}

func Test_CustomFromThinkProviderJSONStringContent(t *testing.T) {
	handler := newCustomThinkProviderHandler()
	router := agentws.NewChanRouter(handler)

	const payload = `{"type":"__customFromThinkProvider","content":"provider returned plain text"}`
	if err := router.Message([]byte(payload)); err != nil {
		t.Fatalf("CustomFromThinkProvider string content must route successfully: %v", err)
	}

	select {
	case response := <-handler.customFromThinkProvider:
		if string(response.Content) != `"provider returned plain text"` {
			t.Errorf("non-JSON provider content must remain a JSON string: %s", response.Content)
		}
	case <-time.After(time.Second):
		t.Fatal("CustomFromThinkProvider string content was not delivered to the handler")
	}
}

func Test_CustomFromThinkProviderFallsBackToUnhandled(t *testing.T) {
	handler := newLegacyFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(customFromThinkProviderPayload)); err != nil {
		t.Fatalf("CustomFromThinkProvider must fall back successfully: %v", err)
	}

	select {
	case message := <-handler.unhandled:
		if string(*message) != customFromThinkProviderPayload {
			t.Errorf("unexpected unhandled message: %s", *message)
		}
	case <-time.After(time.Second):
		t.Fatal("CustomFromThinkProvider was not delivered to the unhandled handler")
	}
}
