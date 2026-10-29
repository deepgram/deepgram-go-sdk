// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
	"github.com/dvonthenen/websocket"
)

func Test_ClearControlMessage(t *testing.T) {
	if MessageTypeReset != MessageTypeClear {
		t.Fatalf("MessageTypeReset = %q, want %q", MessageTypeReset, MessageTypeClear)
	}

	t.Run("callback Clear", func(t *testing.T) {
		testCallbackControlMessage(t, func(client *WSCallback) error { return client.Clear() })
	})
	t.Run("callback Reset", func(t *testing.T) {
		testCallbackControlMessage(t, func(client *WSCallback) error { return client.Reset() })
	})
	t.Run("channel Clear", func(t *testing.T) {
		testChannelControlMessage(t, func(client *WSChannel) error { return client.Clear() })
	})
	t.Run("channel Reset", func(t *testing.T) {
		testChannelControlMessage(t, func(client *WSChannel) error { return client.Reset() })
	})
}

func Test_SendControlJSON(t *testing.T) {
	payload := json.RawMessage("{\n  \"type\": \"Custom\",\n  \"value\": [1, true]\n}")

	t.Run("callback", func(t *testing.T) {
		host, received := newControlMessageServer(t)
		client, err := NewUsingCallback(context.Background(), "test-api-key", testClientOptions(host), testSpeakOptions(), nil)
		if err != nil {
			t.Fatalf("create callback client: %v", err)
		}
		t.Cleanup(client.Stop)
		if !client.Connect() {
			t.Fatal("connect callback client")
		}
		if err := client.SendControlJSON(payload); err != nil {
			t.Fatalf("send raw control message: %v", err)
		}
		expectControlMessage(t, received, string(payload))
	})

	t.Run("channel", func(t *testing.T) {
		host, received := newControlMessageServer(t)
		client, err := NewUsingChan(context.Background(), "test-api-key", testClientOptions(host), testSpeakOptions(), nil)
		if err != nil {
			t.Fatalf("create channel client: %v", err)
		}
		t.Cleanup(client.Stop)
		if !client.Connect() {
			t.Fatal("connect channel client")
		}
		if err := client.SendControlJSON(payload); err != nil {
			t.Fatalf("send raw control message: %v", err)
		}
		expectControlMessage(t, received, string(payload))
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var client WSCallback
		if err := client.SendControlJSON(json.RawMessage(`{"type":`)); !errors.Is(err, common.ErrInvalidJSON) {
			t.Errorf("SendControlJSON() error = %v, want ErrInvalidJSON", err)
		}
	})
}

func testCallbackControlMessage(t *testing.T, call func(*WSCallback) error) {
	t.Helper()
	host, received := newControlMessageServer(t)
	client, err := NewUsingCallback(context.Background(), "test-api-key", testClientOptions(host), testSpeakOptions(), nil)
	if err != nil {
		t.Fatalf("create callback client: %v", err)
	}
	t.Cleanup(client.Stop)
	if !client.Connect() {
		t.Fatal("connect callback client")
	}
	if err := call(client); err != nil {
		t.Fatalf("send control message: %v", err)
	}
	expectClearControlMessage(t, received)
}

func testChannelControlMessage(t *testing.T, call func(*WSChannel) error) {
	t.Helper()
	host, received := newControlMessageServer(t)
	client, err := NewUsingChan(context.Background(), "test-api-key", testClientOptions(host), testSpeakOptions(), nil)
	if err != nil {
		t.Fatalf("create channel client: %v", err)
	}
	t.Cleanup(client.Stop)
	if !client.Connect() {
		t.Fatal("connect channel client")
	}
	if err := call(client); err != nil {
		t.Fatalf("send control message: %v", err)
	}
	expectClearControlMessage(t, received)
}

func newControlMessageServer(t *testing.T) (string, <-chan string) {
	t.Helper()
	received := make(chan string, 2)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.TextMessage {
				received <- string(data)
			}
		}
	}))
	t.Cleanup(server.Close)

	return "ws" + strings.TrimPrefix(server.URL, "http"), received
}

func testClientOptions(host string) *interfaces.ClientOptions {
	return &interfaces.ClientOptions{
		Host: host,
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}
}

func testSpeakOptions() *interfaces.WSSpeakOptions {
	return &interfaces.WSSpeakOptions{Model: "aura-asteria-en"}
}

func expectClearControlMessage(t *testing.T, received <-chan string) {
	t.Helper()
	expectControlMessage(t, received, `{"type":"Clear"}`)
}

func expectControlMessage(t *testing.T, received <-chan string, want string) {
	t.Helper()
	select {
	case got := <-received:
		if got != want {
			t.Errorf("control message = %s, want %s", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for control message")
	}
}
