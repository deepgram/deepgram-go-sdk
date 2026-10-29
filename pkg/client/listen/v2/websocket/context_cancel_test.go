// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

type contextCloseCallback struct {
	closeOnce sync.Once
	closed    chan struct{}
	closes    atomic.Int32
}

func (c *contextCloseCallback) Open(*msginterfaces.OpenResponse) error { return nil }
func (c *contextCloseCallback) Connected(*msginterfaces.ConnectedResponse) error {
	return nil
}
func (c *contextCloseCallback) TurnInfo(*msginterfaces.TurnInfoResponse) error { return nil }
func (c *contextCloseCallback) ConfigureSuccess(*msginterfaces.ConfigureSuccessResponse) error {
	return nil
}
func (c *contextCloseCallback) ConfigureFailure(*msginterfaces.ConfigureFailureResponse) error {
	return nil
}
func (c *contextCloseCallback) FatalError(*msginterfaces.FatalErrorResponse) error { return nil }
func (c *contextCloseCallback) Close(*msginterfaces.CloseResponse) error {
	c.closes.Add(1)
	c.closeOnce.Do(func() { close(c.closed) })
	return nil
}
func (c *contextCloseCallback) Error(*msginterfaces.ErrorResponse) error { return nil }
func (c *contextCloseCallback) UnhandledEvent([]byte) error              { return nil }

func Test_ContextCancellationClosesFluxConnection(t *testing.T) {
	serverClosed := make(chan struct{}, 1)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		_, _, _ = conn.ReadMessage()
		serverClosed <- struct{}{}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	handler := &contextCloseCallback{closed: make(chan struct{})}
	client, err := NewUsingCallback(ctx, "test-api-key",
		&interfaces.ClientOptionsV2{Host: "ws" + strings.TrimPrefix(server.URL, "http")},
		&interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, handler)
	if err != nil {
		t.Fatalf("NewUsingCallback failed: %v", err)
	}
	if !client.Connect() {
		t.Fatal("connect")
	}
	cancel()

	for name, signal := range map[string]<-chan struct{}{
		"callback close": handler.closed,
		"server close":   serverClosed,
	} {
		select {
		case <-signal:
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for %s", name)
		}
	}

	if got := handler.closes.Load(); got != 1 {
		t.Errorf("Close events = %d, want 1", got)
	}
}
