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

	speakws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	speak "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
)

type speakRawControlHandler struct {
	cleared chan *speakws.ClearedResponse
	errors  chan *speakws.ErrorResponse
}

func (h *speakRawControlHandler) Open(*speakws.OpenResponse) error         { return nil }
func (h *speakRawControlHandler) Metadata(*speakws.MetadataResponse) error { return nil }
func (h *speakRawControlHandler) Flush(*speakws.FlushedResponse) error     { return nil }
func (h *speakRawControlHandler) Clear(response *speakws.ClearedResponse) error {
	h.cleared <- response
	return nil
}
func (h *speakRawControlHandler) Close(*speakws.CloseResponse) error     { return nil }
func (h *speakRawControlHandler) Warning(*speakws.WarningResponse) error { return nil }
func (h *speakRawControlHandler) Error(response *speakws.ErrorResponse) error {
	h.errors <- response
	return nil
}
func (h *speakRawControlHandler) UnhandledEvent([]byte) error { return nil }
func (h *speakRawControlHandler) Binary([]byte) error         { return nil }

func TestDaily_SpeakRawControlJSON(t *testing.T) {
	if os.Getenv("DEEPGRAM_API_KEY") == "" && os.Getenv("DEEPGRAM_ACCESS_TOKEN") == "" {
		t.Skip("DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	handler := &speakRawControlHandler{
		cleared: make(chan *speakws.ClearedResponse, 1),
		errors:  make(chan *speakws.ErrorResponse, 1),
	}
	client, err := speak.NewWSUsingCallback(ctx, "", &interfaces.ClientOptions{}, &interfaces.WSSpeakOptions{
		Model: "aura-asteria-en",
	}, handler)
	if err != nil {
		t.Fatalf("create Speak WebSocket: %v", err)
	}
	t.Cleanup(client.Stop)

	if !client.Connect() {
		t.Fatal("connect Speak WebSocket")
	}
	payload := json.RawMessage("{\n  \"type\": \"Clear\"\n}")
	if err := client.SendControlJSON(payload); err != nil {
		t.Fatalf("send raw Clear control: %v", err)
	}

	select {
	case response := <-handler.cleared:
		if response.Type != "Cleared" {
			t.Errorf("Cleared type = %q, want %q", response.Type, "Cleared")
		}
	case response := <-handler.errors:
		t.Fatalf("Speak error: %s", response.Description)
	case <-ctx.Done():
		t.Fatalf("wait for Cleared response: %v", ctx.Err())
	}
}
