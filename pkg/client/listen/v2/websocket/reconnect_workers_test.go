// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

type reconnectingFluxClient interface {
	Connect() bool
	AttemptReconnect(context.Context, int64) bool
	Stream(io.Reader) error
	Stop()
}

// audioReader continuously yields audio so Stream exercises WriteBinary during reconnect.
type audioReader struct{}

func (audioReader) Read(p []byte) (int, error) {
	time.Sleep(time.Millisecond)
	p[0] = 0
	return 1, nil
}

// The stream tests exercise Stream reading the client context while AttemptReconnect
// replaces it. They detect a regression only under the race detector:
// go test -race ./pkg/client/listen/v2/websocket
func Test_StreamDuringReconnectIsRaceFree(t *testing.T) {
	testStreamDuringReconnect(t, func(options *interfaces.ClientOptionsV2) (reconnectingFluxClient, error) {
		return NewUsingCallback(context.Background(), "test-api-key", options, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	})
}

func Test_StreamDuringReconnectIsRaceFreeChannel(t *testing.T) {
	testStreamDuringReconnect(t, func(options *interfaces.ClientOptionsV2) (reconnectingFluxClient, error) {
		return NewUsingChan(context.Background(), "test-api-key", options, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	})
}

func testStreamDuringReconnect(t *testing.T, newClient func(*interfaces.ClientOptionsV2) (reconnectingFluxClient, error)) {
	t.Helper()
	audioReceived := make(chan struct{}, 1)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.BinaryMessage {
				select {
				case audioReceived <- struct{}{}:
				default:
				}
			}
		}
	}))
	t.Cleanup(server.Close)

	options := &interfaces.ClientOptionsV2{
		Host: "ws" + strings.TrimPrefix(server.URL, "http"),
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}
	client, err := newClient(options)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	if !client.Connect() {
		t.Fatal("connect")
	}
	streamDone := make(chan struct{})
	go func() {
		defer close(streamDone)
		_ = client.Stream(audioReader{})
	}()

	select {
	case <-audioReceived:
	case <-time.After(time.Second):
		t.Fatal("stream did not write audio")
	}
	if !client.AttemptReconnect(context.Background(), 3) {
		t.Fatal("reconnect")
	}
	time.Sleep(20 * time.Millisecond)

	client.Stop()
	select {
	case <-streamDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Stream did not return after Stop")
	}
}

func Test_ReconnectDoesNotDuplicatePing(t *testing.T) {
	testReconnectDoesNotDuplicatePing(t, func(options *interfaces.ClientOptionsV2) (reconnectingFluxClient, error) {
		return NewUsingCallback(context.Background(), "test-api-key", options, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	})
}

func Test_ReconnectDoesNotDuplicatePingChannel(t *testing.T) {
	testReconnectDoesNotDuplicatePing(t, func(options *interfaces.ClientOptionsV2) (reconnectingFluxClient, error) {
		return NewUsingChan(context.Background(), "test-api-key", options, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	})
}

func testReconnectDoesNotDuplicatePing(t *testing.T, newClient func(*interfaces.ClientOptionsV2) (reconnectingFluxClient, error)) {
	t.Helper()
	originalPingPeriod := pingPeriod
	pingPeriod = 50 * time.Millisecond
	t.Cleanup(func() { pingPeriod = originalPingPeriod })

	var connections, pings atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connection := connections.Add(1)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		if connection == 1 {
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
	client, err := newClient(options)
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

	// Wait for the first ping instead of assuming the ticker fires on time.
	deadline := time.Now().Add(2 * time.Second)
	for pings.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if pings.Load() == 0 {
		t.Fatal("no ping received after reconnect")
	}
	// One ticker sends at most periods+1 pings in that many periods; it drops
	// late ticks rather than bursting. A duplicate worker sends about twice as many.
	const periods = 4
	pings.Store(0)
	time.Sleep(periods * pingPeriod)
	if got := pings.Load(); got > periods+1 {
		t.Errorf("pings in %d ping periods = %d, want at most %d", periods, got, periods+1)
	}
}
