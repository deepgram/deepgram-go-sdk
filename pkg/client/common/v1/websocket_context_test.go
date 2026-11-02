// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

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

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type contextCloseRouter struct {
	closeOnce sync.Once
	closed    chan struct{}
	closes    atomic.Int32
}

func (r *contextCloseRouter) Open(*commoninterfaces.OpenResponse) error { return nil }
func (r *contextCloseRouter) Message([]byte) error                      { return nil }
func (r *contextCloseRouter) Binary([]byte) error                       { return nil }
func (r *contextCloseRouter) Error(*commoninterfaces.ErrorResponse) error {
	return nil
}
func (r *contextCloseRouter) Close(*commoninterfaces.CloseResponse) error {
	r.closes.Add(1)
	r.closeOnce.Do(func() { close(r.closed) })
	return nil
}

type contextCloseHandler struct {
	url        string
	finishOnce sync.Once
	finished   chan struct{}
	finishes   atomic.Int32
}

func (h *contextCloseHandler) GetURL(string) (string, error) { return h.url, nil }
func (h *contextCloseHandler) ProcessMessage(int, []byte) error {
	return nil
}
func (h *contextCloseHandler) ProcessError(error) error { return nil }
func (h *contextCloseHandler) Start()                   {}
func (h *contextCloseHandler) Finish() {
	h.finishes.Add(1)
	h.finishOnce.Do(func() { close(h.finished) })
}
func (h *contextCloseHandler) GetCloseMsg() []byte { return nil }

type blockingContextRouter struct {
	openStarted  chan struct{}
	openRelease  chan struct{}
	closeStarted chan struct{}
	closeRelease chan struct{}
}

func (r *blockingContextRouter) Open(*commoninterfaces.OpenResponse) error {
	close(r.openStarted)
	<-r.openRelease
	return nil
}
func (r *blockingContextRouter) Message([]byte) error { return nil }
func (r *blockingContextRouter) Binary([]byte) error  { return nil }
func (r *blockingContextRouter) Error(*commoninterfaces.ErrorResponse) error {
	return nil
}
func (r *blockingContextRouter) Close(*commoninterfaces.CloseResponse) error {
	close(r.closeStarted)
	<-r.closeRelease
	return nil
}

func Test_WSClientContextCancellationClosesActiveConnection(t *testing.T) {
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
	handler := &contextCloseHandler{
		url:      "ws" + strings.TrimPrefix(server.URL, "http"),
		finished: make(chan struct{}),
	}
	router := &contextCloseRouter{closed: make(chan struct{})}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect")
	}
	cancel()

	for name, signal := range map[string]<-chan struct{}{
		"router close":   router.closed,
		"handler finish": handler.finished,
		"server close":   serverClosed,
	} {
		select {
		case <-signal:
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for %s", name)
		}
	}

	if got := router.closes.Load(); got != 1 {
		t.Errorf("Close events = %d, want 1", got)
	}
	if got := handler.finishes.Load(); got != 1 {
		t.Errorf("Finish calls = %d, want 1", got)
	}
}

func Test_WSClientStaleContextDoesNotCloseReconnectedSession(t *testing.T) {
	secondMessageReceived := make(chan struct{}, 1)
	secondConnectionClosed := make(chan struct{}, 1)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	var connections atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connection := connections.Add(1)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if connection == 2 {
					secondConnectionClosed <- struct{}{}
				}
				return
			}
			if connection == 2 {
				secondMessageReceived <- struct{}{}
			}
		}
	}))
	defer server.Close()

	firstCtx, cancelFirst := context.WithCancel(context.Background())
	defer cancelFirst()
	handler := &contextCloseHandler{
		url:      "ws" + strings.TrimPrefix(server.URL, "http"),
		finished: make(chan struct{}),
	}
	router := &contextCloseRouter{closed: make(chan struct{})}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(firstCtx, cancelFirst, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	firstWS, response, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("connect first session: %v", err)
	}
	if response != nil {
		response.Body.Close()
	}
	defer firstWS.Close()

	firstSessionDone := make(chan struct{})
	client.muConn.Lock()
	client.wsconn = firstWS
	client.sessionDone = firstSessionDone
	client.terminalCloseStarted = false
	client.muConn.Unlock()
	firstWatcherDone := client.closeOnContextDone(firstCtx, firstWS, firstSessionDone)

	secondCtx, cancelSecond := context.WithCancel(context.Background())
	defer cancelSecond()
	secondWS, response, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("reconnect second session: %v", err)
	}
	if response != nil {
		response.Body.Close()
	}
	defer secondWS.Close()

	secondSessionDone := make(chan struct{})
	client.muConn.Lock()
	client.wsconn = secondWS
	client.sessionDone = secondSessionDone
	client.terminalCloseStarted = false
	client.muConn.Unlock()
	secondWatcherDone := client.closeOnContextDone(secondCtx, secondWS, secondSessionDone)

	cancelFirst()
	select {
	case <-firstWatcherDone:
	case <-time.After(time.Second):
		t.Fatal("stale context watcher did not complete its identity check")
	}

	if err := secondWS.WriteMessage(websocket.TextMessage, []byte("still connected")); err != nil {
		t.Fatalf("stale context closed the reconnected session: %v", err)
	}
	select {
	case <-secondMessageReceived:
	case <-time.After(time.Second):
		t.Fatal("reconnected session did not receive a message")
	}

	cancelSecond()
	select {
	case <-secondConnectionClosed:
	case <-time.After(time.Second):
		t.Fatal("active context did not close the reconnected session")
	}
	select {
	case <-secondWatcherDone:
	case <-time.After(time.Second):
		t.Fatal("active context watcher did not finish")
	}
}

func Test_WSClientContextCancellationClosesTransportBeforeBlockingHandlers(t *testing.T) {
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
	handler := &contextCloseHandler{
		url:      "ws" + strings.TrimPrefix(server.URL, "http"),
		finished: make(chan struct{}),
	}
	router := &blockingContextRouter{
		openStarted:  make(chan struct{}),
		openRelease:  make(chan struct{}),
		closeStarted: make(chan struct{}),
		closeRelease: make(chan struct{}),
	}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

	connected := make(chan bool, 1)
	go func() { connected <- client.ConnectWithCancel(ctx, cancel, 1) }()
	select {
	case <-router.openStarted:
	case <-time.After(time.Second):
		t.Fatal("Open handler was not called")
	}

	close(router.openRelease)
	select {
	case ok := <-connected:
		if !ok {
			t.Fatal("ConnectWithCancel failed")
		}
	case <-time.After(time.Second):
		t.Fatal("ConnectWithCancel did not return")
	}

	cancel()
	select {
	case <-router.closeStarted:
	case <-time.After(time.Second):
		t.Fatal("Close handler was not called")
	}

	select {
	case <-serverClosed:
	case <-time.After(time.Second):
		t.Fatal("context cancellation did not close the transport while Close was blocked")
	}

	close(router.closeRelease)
	select {
	case <-handler.finished:
	case <-time.After(time.Second):
		t.Fatal("handler Finish was not called after Close returned")
	}
}

func Test_WSClientReportsSuccessfulHandshakeBeforeImmediatePeerClose(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(time.Second),
		)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := &contextCloseHandler{
		url:      "ws" + strings.TrimPrefix(server.URL, "http"),
		finished: make(chan struct{}),
	}
	router := &blockingContextRouter{
		openStarted:  make(chan struct{}),
		openRelease:  make(chan struct{}),
		closeStarted: make(chan struct{}),
		closeRelease: make(chan struct{}),
	}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

	connected := make(chan bool, 1)
	go func() { connected <- client.ConnectWithCancel(ctx, cancel, 1) }()
	for name, signal := range map[string]<-chan struct{}{
		"Open handler":  router.openStarted,
		"Close handler": router.closeStarted,
	} {
		select {
		case <-signal:
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for %s", name)
		}
	}

	close(router.closeRelease)
	close(router.openRelease)
	select {
	case ok := <-connected:
		if !ok {
			t.Error("ConnectWithCancel returned false after a successful handshake")
		}
	case <-time.After(time.Second):
		t.Fatal("ConnectWithCancel did not return after releasing handlers")
	}
}

func Test_WSClientContextCancellationDuringOpenReturnsFailure(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()
		_, _, _ = conn.ReadMessage()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	handler := &contextCloseHandler{
		url:      "ws" + strings.TrimPrefix(server.URL, "http"),
		finished: make(chan struct{}),
	}
	router := &blockingContextRouter{
		openStarted:  make(chan struct{}),
		openRelease:  make(chan struct{}),
		closeStarted: make(chan struct{}),
		closeRelease: make(chan struct{}),
	}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

	connected := make(chan bool, 1)
	go func() { connected <- client.ConnectWithCancel(ctx, cancel, 1) }()
	select {
	case <-router.openStarted:
	case <-time.After(time.Second):
		t.Fatal("Open handler was not called")
	}

	cancel()
	close(router.openRelease)
	select {
	case ok := <-connected:
		if ok {
			t.Error("ConnectWithCancel returned true after context cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("ConnectWithCancel did not return after Open was released")
	}

	select {
	case <-router.closeStarted:
	case <-time.After(time.Second):
		t.Fatal("Close handler was not called")
	}
	close(router.closeRelease)
}
