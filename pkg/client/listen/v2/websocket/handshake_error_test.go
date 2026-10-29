// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

type handshakeErrorCallback struct {
	errors chan *msginterfaces.ErrorResponse
}

func (c *handshakeErrorCallback) Open(*msginterfaces.OpenResponse) error { return nil }
func (c *handshakeErrorCallback) Connected(*msginterfaces.ConnectedResponse) error {
	return nil
}
func (c *handshakeErrorCallback) TurnInfo(*msginterfaces.TurnInfoResponse) error { return nil }
func (c *handshakeErrorCallback) ConfigureSuccess(*msginterfaces.ConfigureSuccessResponse) error {
	return nil
}
func (c *handshakeErrorCallback) ConfigureFailure(*msginterfaces.ConfigureFailureResponse) error {
	return nil
}
func (c *handshakeErrorCallback) FatalError(*msginterfaces.FatalErrorResponse) error { return nil }
func (c *handshakeErrorCallback) Close(*msginterfaces.CloseResponse) error           { return nil }
func (c *handshakeErrorCallback) Error(response *msginterfaces.ErrorResponse) error {
	c.errors <- response
	return nil
}
func (c *handshakeErrorCallback) UnhandledEvent([]byte) error { return nil }

func Test_HandshakeErrorUsesSharedTransport(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"err_code":"INVALID_QUERY","err_msg":"invalid endpointing","description":"endpointing must be an integer"}`))
	}))
	defer server.Close()

	handler := &handshakeErrorCallback{errors: make(chan *msginterfaces.ErrorResponse, 1)}
	client, err := NewUsingCallback(
		context.Background(), "test-api-key",
		&interfaces.ClientOptionsV2{Host: "ws" + strings.TrimPrefix(server.URL, "http")},
		&interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, handler)
	if err != nil {
		t.Fatalf("NewUsingCallback failed: %v", err)
	}
	if client.Connect() {
		t.Fatal("Connect succeeded after a 400 response")
	}
	if got := attempts.Load(); got != 1 {
		t.Errorf("handshake attempts = %d, want 1", got)
	}

	select {
	case response := <-handler.errors:
		if response.ErrCode != "INVALID_QUERY" {
			t.Errorf("ErrCode = %q, want %q", response.ErrCode, "INVALID_QUERY")
		}
		if response.Variant != "400 Bad Request" {
			t.Errorf("Variant = %q, want %q", response.Variant, "400 Bad Request")
		}
		if response.Description != "endpointing must be an integer" {
			t.Errorf("Description = %q, want %q", response.Description, "endpointing must be an integer")
		}
	default:
		t.Fatal("expected a Flux error callback after the failed handshake")
	}
}
