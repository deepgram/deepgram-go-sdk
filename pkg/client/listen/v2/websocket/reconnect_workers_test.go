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
	"time"

	"github.com/dvonthenen/websocket"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func Test_ReconnectDoesNotDuplicatePing(t *testing.T) {
	originalPingPeriod := pingPeriod
	pingPeriod = 50 * time.Millisecond
	t.Cleanup(func() { pingPeriod = originalPingPeriod })

	var connections, pings atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		if connections.Add(1) == 1 {
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(time.Second),
			)
			return
		}

		conn.SetPingHandler(func(string) error {
			pings.Add(1)
			return nil
		})
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	options := &interfaces.ClientOptionsV2{
		Host:            "ws" + strings.TrimPrefix(server.URL, "http"),
		EnableKeepAlive: true,
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}
	client, err := NewUsingCallback(context.Background(), "test-api-key", options, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	t.Cleanup(client.Stop)

	if !client.Connect() {
		t.Fatal("connect")
	}
	time.Sleep(100 * time.Millisecond)
	if !client.AttemptReconnect(context.Background(), 3) {
		t.Fatal("reconnect")
	}
	if got := connections.Load(); got != 2 {
		t.Fatalf("connections = %d, want 2", got)
	}

	time.Sleep(pingPeriod + 25*time.Millisecond)
	if got := pings.Load(); got != 1 {
		t.Errorf("pings after reconnect = %d, want 1", got)
	}
}
