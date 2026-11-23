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

// After an automatic reconnect, the client must send one KeepAlive per period,
// not one per connection it has had.
func Test_AutomaticReconnectDoesNotDuplicateKeepAlive(t *testing.T) {
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
			// Drop the transport without a WebSocket close frame.
			_ = conn.UnderlyingConn().Close()
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
		ReconnectPolicy: &interfaces.ReconnectPolicy{
			MaxRetries:   1,
			InitialDelay: time.Millisecond,
			MaxDelay:     time.Millisecond,
		},
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
	deadline := time.Now().Add(time.Second)
	for connections.Load() != 2 {
		if time.Now().After(deadline) {
			t.Fatalf("connections = %d, want 2", connections.Load())
		}
		time.Sleep(time.Millisecond)
	}

	// The first keepalive of a single worker arrives one period after the reconnect; a second
	// worker from the first connection would add another on roughly the same tick.
	time.Sleep(pingPeriod + 700*time.Millisecond)
	if got := keepAlives.Load(); got != 1 {
		t.Errorf("KeepAlive messages after one period = %d, want 1", got)
	}
}
