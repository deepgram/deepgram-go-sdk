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
	"sync"
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

func cleanupReconnectClient(client reconnectingFluxClient) {
	client.Stop()
	switch client := client.(type) {
	case *WSCallback:
		client.stopWorkers()
	case *WSChannel:
		client.stopWorkers()
	}
}

// idleReader yields no audio, so Stream keeps polling its context until the client stops.
type idleReader struct{}

func (idleReader) Read([]byte) (int, error) {
	time.Sleep(time.Millisecond)
	return 0, nil
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
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
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
		_ = client.Stream(idleReader{})
	}()

	time.Sleep(20 * time.Millisecond)
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

func Test_ReconnectReusesActivePingWorker(t *testing.T) {
	originalPingPeriod := pingPeriod
	pingPeriod = 20 * time.Millisecond
	t.Cleanup(func() { pingPeriod = originalPingPeriod })

	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewUsingCallback(context.Background(), "test-api-key", &interfaces.ClientOptionsV2{
		Host:            "ws" + strings.TrimPrefix(server.URL, "http"),
		EnableKeepAlive: true,
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	t.Cleanup(func() {
		client.Stop()
		client.stopWorkers()
	})

	pingStarted := make(chan struct{})
	releasePing := make(chan struct{})
	var blockOnce sync.Once
	client.beforePing = func() {
		blockOnce.Do(func() {
			close(pingStarted)
			<-releasePing
		})
	}

	if !client.Connect() {
		t.Fatal("connect")
	}
	select {
	case <-pingStarted:
	case <-time.After(time.Second):
		t.Fatal("ping did not start")
	}

	reconnected := make(chan bool, 1)
	go func() { reconnected <- client.AttemptReconnect(context.Background(), 3) }()
	select {
	case ok := <-reconnected:
		if !ok {
			t.Fatal("reconnect failed")
		}
	case <-time.After(time.Second):
		t.Fatal("healthy reconnect did not reuse the active connection")
	}

	close(releasePing)
}

func Test_StartWaitsForCancelledPingWorker(t *testing.T) {
	originalPingPeriod := pingPeriod
	pingPeriod = 20 * time.Millisecond
	t.Cleanup(func() { pingPeriod = originalPingPeriod })

	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewUsingCallback(context.Background(), "test-api-key", &interfaces.ClientOptionsV2{
		Host:            "ws" + strings.TrimPrefix(server.URL, "http"),
		EnableKeepAlive: true,
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	t.Cleanup(func() {
		client.Stop()
		client.stopWorkers()
	})

	pingStarted := make(chan struct{})
	releasePing := make(chan struct{})
	var blockOnce sync.Once
	client.beforePing = func() {
		blockOnce.Do(func() {
			close(pingStarted)
			<-releasePing
		})
	}

	if !client.Connect() {
		t.Fatal("connect")
	}
	select {
	case <-pingStarted:
	case <-time.After(time.Second):
		t.Fatal("ping did not start")
	}

	client.Finish()
	started := make(chan struct{})
	go func() {
		client.Start()
		close(started)
	}()
	select {
	case <-started:
		t.Fatal("Start returned while the cancelled ping was in flight")
	case <-time.After(100 * time.Millisecond):
	}

	close(releasePing)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("Start did not wait for the cancelled ping worker")
	}
}

func Test_WriteBinaryReconnectWaitsForInFlightPing(t *testing.T) {
	originalPingPeriod := pingPeriod
	pingPeriod = 20 * time.Millisecond
	t.Cleanup(func() { pingPeriod = originalPingPeriod })

	closeFirst := make(chan struct{})
	secondBinary := make(chan struct{}, 1)
	var connections atomic.Int32
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
			<-closeFirst
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			return
		}

		messageType, _, err := conn.ReadMessage()
		if err == nil && messageType == websocket.BinaryMessage {
			secondBinary <- struct{}{}
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewUsingCallback(context.Background(), "test-api-key", &interfaces.ClientOptionsV2{
		Host:            "ws" + strings.TrimPrefix(server.URL, "http"),
		EnableKeepAlive: true,
		WSHeaderProcessor: func(headers http.Header) {
			headers.Del("Host")
		},
	}, &interfaces.FluxTranscriptionOptions{Model: "flux-general-en"}, nil)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	t.Cleanup(func() {
		client.Stop()
		client.stopWorkers()
	})

	pingStarted := make(chan struct{})
	releasePing := make(chan struct{})
	var blockOnce sync.Once
	client.beforePing = func() {
		blockOnce.Do(func() {
			close(pingStarted)
			<-releasePing
		})
	}

	if !client.Connect() {
		t.Fatal("connect")
	}
	select {
	case <-pingStarted:
	case <-time.After(time.Second):
		t.Fatal("ping did not start")
	}

	close(closeFirst)
	deadline := time.Now().Add(time.Second)
	for client.WSClient.HasActiveConnection() && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if client.WSClient.HasActiveConnection() {
		t.Fatal("peer close did not clear the active transport")
	}

	written := make(chan error, 1)
	go func() { written <- client.WriteBinary([]byte{1}) }()
	select {
	case <-written:
		t.Fatal("WriteBinary reconnected while the old ping was in flight")
	case <-time.After(100 * time.Millisecond):
	}

	close(releasePing)
	select {
	case err := <-written:
		if err != nil {
			t.Fatalf("WriteBinary: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("WriteBinary did not reconnect after the old ping exited")
	}
	select {
	case <-secondBinary:
	case <-time.After(time.Second):
		t.Fatal("replacement transport did not receive binary audio")
	}
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
	t.Cleanup(func() { cleanupReconnectClient(client) })

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
