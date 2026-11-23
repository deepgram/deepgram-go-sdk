// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type reconnectTestHandler struct {
	url      string
	starts   atomic.Int32
	finishes atomic.Int32
	onStart  func()
}

func (h *reconnectTestHandler) GetURL(string) (string, error) { return h.url, nil }
func (h *reconnectTestHandler) ProcessMessage(int, []byte) error {
	return nil
}
func (h *reconnectTestHandler) ProcessError(error) error { return nil }
func (h *reconnectTestHandler) Start() {
	h.starts.Add(1)
	if h.onStart != nil {
		h.onStart()
	}
}
func (h *reconnectTestHandler) Finish()             { h.finishes.Add(1) }
func (h *reconnectTestHandler) GetCloseMsg() []byte { return nil }

type reconnectTestRouter struct {
	opens  atomic.Int32
	closes atomic.Int32
}

func (r *reconnectTestRouter) Open(*commoninterfaces.OpenResponse) error {
	r.opens.Add(1)
	return nil
}
func (r *reconnectTestRouter) Message([]byte) error { return nil }
func (r *reconnectTestRouter) Binary([]byte) error  { return nil }
func (r *reconnectTestRouter) Close(*commoninterfaces.CloseResponse) error {
	r.closes.Add(1)
	return nil
}
func (r *reconnectTestRouter) Error(*commoninterfaces.ErrorResponse) error { return nil }

func Test_WSClientAutomaticallyReconnectsAfterUnexpectedTransportFailure(t *testing.T) {
	var connections atomic.Int32
	firstConnected := make(chan struct{})
	headers := make(chan http.Header, 1)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		if connections.Add(1) == 1 {
			close(firstConnected)
			// Closing the transport without a WebSocket close frame simulates an
			// unexpected network failure.
			_ = conn.UnderlyingConn().Close()
			return
		}
		headers <- r.Header.Clone()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	events := make(chan interfaces.ReconnectEvent, 3)
	handler := &reconnectTestHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	router := &reconnectTestRouter{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{
		WSHeaderProcessor: func(headers http.Header) {
			headers.Set("X-Reconnect-Test", "restored")
		},
		ReconnectPolicy: &interfaces.ReconnectPolicy{
			MaxRetries:   2,
			InitialDelay: time.Millisecond,
			MaxDelay:     time.Millisecond,
			OnReconnect: func(event interfaces.ReconnectEvent) {
				events <- event
			},
		},
	}, &handlerInterface, &routerInterface)
	t.Cleanup(client.Stop)

	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect")
	}
	select {
	case <-firstConnected:
	case <-time.After(time.Second):
		t.Fatal("first connection was not established")
	}

	wantEvents := []interfaces.ReconnectEventType{interfaces.ReconnectAttempt, interfaces.ReconnectSuccess}
	for _, want := range wantEvents {
		select {
		case event := <-events:
			if event.Type != want {
				t.Fatalf("event type = %q, want %q", event.Type, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for %q event", want)
		}
	}
	select {
	case got := <-headers:
		if got.Get("Authorization") != "token test-api-key" {
			t.Errorf("Authorization = %q, want token authentication", got.Get("Authorization"))
		}
		if got.Get("X-Reconnect-Test") != "restored" {
			t.Errorf("custom header = %q, want restored", got.Get("X-Reconnect-Test"))
		}
	case <-time.After(time.Second):
		t.Fatal("reconnected request did not reach the server")
	}
	if got := connections.Load(); got != 2 {
		t.Errorf("connections = %d, want 2", got)
	}
	if got := handler.starts.Load(); got != 2 {
		t.Errorf("handler starts = %d, want 2", got)
	}
	if got := handler.finishes.Load(); got != 0 {
		t.Errorf("handler finishes = %d before Stop, want 0", got)
	}
}

func Test_WSClientAutomaticReconnectResendsRequiredProtocolSetup(t *testing.T) {
	var connections atomic.Int32
	setupReceived := make(chan int, 2)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		connection := int(connections.Add(1))
		if _, _, err := conn.ReadMessage(); err != nil {
			t.Errorf("connection %d did not receive required setup: %v", connection, err)
			return
		}
		setupReceived <- connection
		if connection == 1 {
			_ = conn.UnderlyingConn().Close()
			return
		}
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	handler := &reconnectTestHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	router := &reconnectTestRouter{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{
		ReconnectPolicy: &interfaces.ReconnectPolicy{
			InitialDelay: time.Millisecond,
			MaxDelay:     time.Millisecond,
		},
	}, &handlerInterface, &routerInterface)
	t.Cleanup(client.Stop)
	handler.onStart = func() {
		if err := client.WriteJSON(map[string]string{"type": "Settings"}); err != nil {
			t.Errorf("send required protocol setup: %v", err)
		}
	}

	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect")
	}
	for want := 1; want <= 2; want++ {
		select {
		case got := <-setupReceived:
			if got != want {
				t.Fatalf("setup connection = %d, want %d", got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for setup on connection %d", want)
		}
	}
}

func Test_WSClientAutomaticReconnectDoesNotWriteOrReconnectAfterCancellation(t *testing.T) {
	var connections atomic.Int32
	firstConnected := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		connections.Add(1)
		close(firstConnected)
		_ = conn.UnderlyingConn().Close()
	}))
	t.Cleanup(server.Close)

	handler := &reconnectTestHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	router := &reconnectTestRouter{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	ctx, cancel := context.WithCancel(context.Background())
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{
		ReconnectPolicy: &interfaces.ReconnectPolicy{
			InitialDelay: 200 * time.Millisecond,
			MaxDelay:     200 * time.Millisecond,
		},
	}, &handlerInterface, &routerInterface)

	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect")
	}
	select {
	case <-firstConnected:
	case <-time.After(time.Second):
		t.Fatal("first connection was not established")
	}
	deadline := time.Now().Add(time.Second)
	for {
		client.muConn.RLock()
		reconnecting := client.reconnecting
		client.muConn.RUnlock()
		if reconnecting {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("client did not enter automatic reconnect")
		}
		time.Sleep(time.Millisecond)
	}
	if err := client.WriteBinary([]byte("must not be replayed")); !errors.Is(err, ErrInvalidConnection) {
		t.Fatalf("WriteBinary error = %v, want ErrInvalidConnection", err)
	}
	cancel()
	time.Sleep(300 * time.Millisecond)

	if got := connections.Load(); got != 1 {
		t.Errorf("connections after cancellation = %d, want 1", got)
	}
	if got := router.closes.Load(); got != 1 {
		t.Errorf("terminal close events = %d, want 1", got)
	}
	if got := handler.finishes.Load(); got != 1 {
		t.Errorf("handler finishes = %d, want 1", got)
	}
}

func Test_WSClientAutomaticReconnectDoesNotReconnectAfterStop(t *testing.T) {
	var connections atomic.Int32
	firstConnected := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		connections.Add(1)
		close(firstConnected)
		_ = conn.UnderlyingConn().Close()
	}))
	t.Cleanup(server.Close)

	handler := &reconnectTestHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	router := &reconnectTestRouter{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{
		ReconnectPolicy: &interfaces.ReconnectPolicy{
			InitialDelay: 200 * time.Millisecond,
			MaxDelay:     200 * time.Millisecond,
		},
	}, &handlerInterface, &routerInterface)

	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect")
	}
	select {
	case <-firstConnected:
	case <-time.After(time.Second):
		t.Fatal("first connection was not established")
	}
	deadline := time.Now().Add(time.Second)
	for {
		client.muConn.RLock()
		reconnecting := client.reconnecting
		client.muConn.RUnlock()
		if reconnecting {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("client did not enter automatic reconnect")
		}
		time.Sleep(time.Millisecond)
	}

	client.Stop()
	time.Sleep(300 * time.Millisecond)
	if got := connections.Load(); got != 1 {
		t.Errorf("connections after Stop = %d, want 1", got)
	}
	if got := router.closes.Load(); got != 1 {
		t.Errorf("terminal close events = %d, want 1", got)
	}
	if got := handler.finishes.Load(); got != 1 {
		t.Errorf("handler finishes = %d, want 1", got)
	}
}

func Test_WSClientAutomaticReconnectReportsTerminalFailure(t *testing.T) {
	var connections atomic.Int32
	firstConnected := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if connections.Add(1) > 1 {
			http.Error(w, "temporary outage", http.StatusServiceUnavailable)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		close(firstConnected)
		_ = conn.UnderlyingConn().Close()
	}))
	t.Cleanup(server.Close)

	events := make(chan interfaces.ReconnectEvent, 4)
	handler := &reconnectTestHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	router := &reconnectTestRouter{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{
		ReconnectPolicy: &interfaces.ReconnectPolicy{
			MaxRetries:   2,
			InitialDelay: time.Millisecond,
			MaxDelay:     time.Millisecond,
			OnReconnect: func(event interfaces.ReconnectEvent) {
				events <- event
			},
		},
	}, &handlerInterface, &routerInterface)

	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect")
	}
	select {
	case <-firstConnected:
	case <-time.After(time.Second):
		t.Fatal("first connection was not established")
	}

	var terminal interfaces.ReconnectEvent
	for i := 0; i < 3; i++ {
		select {
		case event := <-events:
			if event.Type == interfaces.ReconnectFailure {
				terminal = event
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for reconnect lifecycle event")
		}
	}
	if terminal.Type != interfaces.ReconnectFailure {
		t.Fatal("did not receive terminal reconnect failure")
	}
	if terminal.Attempt != 2 {
		t.Errorf("terminal attempt = %d, want 2", terminal.Attempt)
	}
	if terminal.Error == nil {
		t.Error("terminal failure did not retain the transport error")
	}
	deadline := time.Now().Add(time.Second)
	for router.closes.Load() != 1 || handler.finishes.Load() != 1 {
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if got := router.closes.Load(); got != 1 {
		t.Errorf("terminal close events = %d, want 1", got)
	}
	if got := handler.finishes.Load(); got != 1 {
		t.Errorf("handler finishes = %d, want 1", got)
	}
}
