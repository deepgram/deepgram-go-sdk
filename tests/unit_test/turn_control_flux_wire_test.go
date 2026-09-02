// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listenv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2"
)

// newFluxTestServer starts a TLS WebSocket server that records every text frame it
// receives and answers each ForceEndTurn frame with a FORCE_END_TURN_NO_ACTIVE_TURN
// Warning, mimicking the live /v2/listen behavior when no turn is active.
// It returns the host (ip:port) to dial and a channel of received text frames.
func newFluxTestServer(t *testing.T) (host string, frames chan string, shutdown func()) {
	t.Helper()

	frames = make(chan string, 16)
	upgrader := websocket.Upgrader{}

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("test server upgrade failed: %s", err)
			return
		}
		defer conn.Close()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return // client closed the connection
			}
			if msgType != websocket.TextMessage {
				continue
			}
			frames <- string(msg)
			if strings.Contains(string(msg), `"ForceEndTurn"`) {
				warning := `{"type":"Warning","code":"FORCE_END_TURN_NO_ACTIVE_TURN","description":"ForceEndTurn received with no active turn"}`
				if err := conn.WriteMessage(websocket.TextMessage, []byte(warning)); err != nil {
					return
				}
			}
		}
	}))

	return strings.TrimPrefix(srv.URL, "https://"), frames, srv.Close
}

func fluxWireClientOptions(host string) *interfaces.ClientOptionsV2 {
	return &interfaces.ClientOptionsV2{
		Host:           host,
		SkipServerAuth: true, // the test server uses a self-signed certificate
	}
}

func fluxWireTranscriptionOptions() *interfaces.FluxTranscriptionOptions {
	return &interfaces.FluxTranscriptionOptions{
		Model:      "flux-general-en",
		Encoding:   "linear16",
		SampleRate: 16000,
	}
}

// awaitForceEndTurnFrame asserts exactly one well-formed ForceEndTurn text frame
// reaches the server: the first frame matches the documented wire format and no
// stray second frame follows it.
func awaitForceEndTurnFrame(t *testing.T, frames chan string) {
	t.Helper()

	select {
	case frame := <-frames:
		if frame != forceEndTurnWireFormat {
			t.Errorf("expected frame %s, got %s", forceEndTurnWireFormat, frame)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server never received the ForceEndTurn text frame")
	}

	select {
	case frame := <-frames:
		t.Errorf("expected exactly one frame, got a second one: %s", frame)
	case <-time.After(200 * time.Millisecond):
	}
}

// wireCallbackHandler is a FluxMessageCallback + FluxWarningCallback that records
// warnings delivered through the live client.
type wireCallbackHandler struct {
	baseFluxCallback
	warnings chan *msginterfaces.WarningResponse
}

func (c *wireCallbackHandler) Warning(wr *msginterfaces.WarningResponse) error {
	c.warnings <- wr
	return nil
}

// Test_FluxForceEndTurnWireCallback drives the public callback client end to end
// against a mock WebSocket server (Greg's S1 on PR #350): ForceEndTurn() emits
// exactly one JSON text frame, the resulting server Warning routes back through the
// live client as a non-fatal event, and write errors propagate to the caller.
func Test_FluxForceEndTurnWireCallback(t *testing.T) {
	host, frames, shutdown := newFluxTestServer(t)
	defer shutdown()

	handler := &wireCallbackHandler{warnings: make(chan *msginterfaces.WarningResponse, 4)}
	dgClient, err := listenv2.NewWSUsingCallback(
		context.Background(), MockAPIKey, fluxWireClientOptions(host), fluxWireTranscriptionOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingCallback failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect to the test server failed")
	}

	if err := dgClient.ForceEndTurn(); err != nil {
		t.Fatalf("ForceEndTurn failed: %s", err)
	}
	awaitForceEndTurnFrame(t, frames)

	// The no-active-turn Warning the server sent back must arrive as a routine event.
	select {
	case wr := <-handler.warnings:
		if wr.Code != msginterfaces.WarningCodeForceEndTurnNoActiveTurn {
			t.Errorf("expected warning code %q, got %q", msginterfaces.WarningCodeForceEndTurnNoActiveTurn, wr.Code)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the server's Warning never reached the callback handler")
	}

	dgClient.Stop()
	if err := dgClient.ForceEndTurn(); err == nil {
		t.Error("ForceEndTurn after Stop must propagate a write error, got nil")
	}
}

// wireChanHandler is a FluxMessageChan + FluxWarningChan that records warnings
// delivered through the live client.
type wireChanHandler struct {
	baseFluxChan
	warningChan chan *msginterfaces.WarningResponse
}

func (h *wireChanHandler) GetWarning() []*chan *msginterfaces.WarningResponse {
	return []*chan *msginterfaces.WarningResponse{&h.warningChan}
}

// Test_FluxForceEndTurnWireChannel is the channel-client counterpart of
// Test_FluxForceEndTurnWireCallback.
func Test_FluxForceEndTurnWireChannel(t *testing.T) {
	host, frames, shutdown := newFluxTestServer(t)
	defer shutdown()

	handler := &wireChanHandler{
		baseFluxChan: *newBaseFluxChan(),
		warningChan:  make(chan *msginterfaces.WarningResponse, 4),
	}
	dgClient, err := listenv2.NewWSUsingChan(
		context.Background(), MockAPIKey, fluxWireClientOptions(host), fluxWireTranscriptionOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect to the test server failed")
	}

	if err := dgClient.ForceEndTurn(); err != nil {
		t.Fatalf("ForceEndTurn failed: %s", err)
	}
	awaitForceEndTurnFrame(t, frames)

	select {
	case wr := <-handler.warningChan:
		if wr.Code != msginterfaces.WarningCodeForceEndTurnNoActiveTurn {
			t.Errorf("expected warning code %q, got %q", msginterfaces.WarningCodeForceEndTurnNoActiveTurn, wr.Code)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the server's Warning never reached the warning channel")
	}

	dgClient.Stop()
	if err := dgClient.ForceEndTurn(); err == nil {
		t.Error("ForceEndTurn after Stop must propagate a write error, got nil")
	}
}
