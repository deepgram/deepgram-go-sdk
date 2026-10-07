// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// After the server closes a connection and the caller reconnects, the client must send one
// KeepAlive per period, not one per connection it has had.
func Test_ReconnectDoesNotDuplicateKeepAlive(t *testing.T) {
	if testing.Short() {
		t.Skip("waits for a full keepalive period")
	}

	var connections, keepAlives atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		if connections.Add(1) == 1 {
			// Close the first connection gracefully, as a server does when it ends a session.
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(time.Second),
			)
			return
		}

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.TextMessage && strings.Contains(string(data), MessageTypeKeepAlive) {
				keepAlives.Add(1)
			}
		}
	}))
	t.Cleanup(server.Close)

	options := &interfaces.ClientOptions{
		Host:            "ws" + strings.TrimPrefix(server.URL, "http"),
		EnableKeepAlive: true,
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}
	client, err := NewUsingCallback(context.Background(), "test-api-key", options, &interfaces.LiveTranscriptionOptions{}, nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	t.Cleanup(client.Stop)

	if !client.Connect() {
		t.Fatal("connect")
	}
	// Give the client time to observe the server's close before it reconnects.
	time.Sleep(500 * time.Millisecond)
	if !client.AttemptReconnect(context.Background(), 3) {
		t.Fatal("reconnect")
	}
	if got := connections.Load(); got != 2 {
		t.Fatalf("connections = %d, want 2", got)
	}

	// The first keepalive of a single worker arrives one period after the reconnect; a second
	// worker from the first connection would add another on roughly the same tick.
	time.Sleep(pingPeriod + 700*time.Millisecond)
	if got := keepAlives.Load(); got != 1 {
		t.Errorf("KeepAlive messages after one period = %d, want 1", got)
	}
}
