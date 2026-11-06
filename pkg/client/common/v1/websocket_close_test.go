// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type closeCountingRouter struct {
	mu     sync.Mutex
	closes int
}

func (r *closeCountingRouter) Open(*commoninterfaces.OpenResponse) error { return nil }
func (r *closeCountingRouter) Message([]byte) error                      { return nil }
func (r *closeCountingRouter) Binary([]byte) error                       { return nil }
func (r *closeCountingRouter) Error(*commoninterfaces.ErrorResponse) error {
	return nil
}
func (r *closeCountingRouter) Close(*commoninterfaces.CloseResponse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closes++
	return nil
}

type closeCountingHandler struct {
	mu       sync.Mutex
	finishes int
	urlCalls int
}

func (h *closeCountingHandler) GetURL(string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.urlCalls++
	return "", errors.New("unexpected reconnect")
}
func (h *closeCountingHandler) ProcessMessage(int, []byte) error {
	return nil
}
func (h *closeCountingHandler) ProcessError(error) error { return nil }
func (h *closeCountingHandler) Start()                   {}
func (h *closeCountingHandler) GetCloseMsg() []byte      { return nil }
func (h *closeCountingHandler) Finish() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.finishes++
}

func Test_WSClientTerminalCloseRunsOnce(t *testing.T) {
	router := &closeCountingRouter{}
	handler := &closeCountingHandler{}
	var routerInterface commoninterfaces.Router = router
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	client := WSClient{
		router:          &routerInterface,
		processMessages: &handlerInterface,
	}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.closeWsWithState(true, false, false)
		}()
	}
	wg.Wait()

	if router.closes != 1 {
		t.Fatalf("router Close calls = %d, want 1", router.closes)
	}
	if handler.finishes != 1 {
		t.Fatalf("handler Finish calls = %d, want 1", handler.finishes)
	}
}

func Test_WriteControlWithoutConnectionDoesNotReconnect(t *testing.T) {
	handler := &closeCountingHandler{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	client := WSClient{
		retry:           true,
		processMessages: &handlerInterface,
	}

	err := client.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
	if !errors.Is(err, ErrInvalidConnection) {
		t.Fatalf("WriteControl error = %v, want ErrInvalidConnection", err)
	}
	if handler.urlCalls != 0 {
		t.Errorf("WriteControl attempted %d reconnects, want 0", handler.urlCalls)
	}
}

func Test_WriteBinaryWithoutConnectionAttemptsReconnect(t *testing.T) {
	handler := &closeCountingHandler{}
	router := &closeCountingRouter{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(context.Background(), func() {}, "", &interfaces.ClientOptions{SelfHosted: true}, &handlerInterface, &routerInterface)
	if client == nil {
		t.Fatal("NewWS returned nil")
	}
	client.retry = true
	client.retryCnt = 1

	err := client.WriteBinary([]byte{1})
	if !errors.Is(err, ErrInvalidConnection) {
		t.Fatalf("WriteBinary error = %v, want ErrInvalidConnection", err)
	}
	if handler.urlCalls != 1 {
		t.Errorf("WriteBinary attempted %d reconnects, want 1", handler.urlCalls)
	}
}

func Test_InternalConnectWaitsForCancelledSession(t *testing.T) {
	handler := &closeCountingHandler{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	oldCtx, cancelOld := context.WithCancel(context.Background())
	sessionDone := make(chan struct{})
	client := NewWS(oldCtx, cancelOld, "", &interfaces.ClientOptions{SelfHosted: true}, &handlerInterface, nil)
	if client == nil {
		t.Fatal("NewWS returned nil")
	}
	client.wsconn = &websocket.Conn{}
	client.sessionDone = sessionDone
	cancelOld()

	newCtx := context.Background()
	done := make(chan struct{})
	go func() {
		_ = client.internalConnectWithCancel(newCtx, func() {}, 1, true, nil)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("reconnect returned before the cancelled session closed")
	case <-time.After(50 * time.Millisecond):
	}
	client.muConn.RLock()
	currentCtx := client.ctx
	client.muConn.RUnlock()
	if currentCtx != oldCtx {
		t.Fatal("reconnect replaced the cancelled session context before it closed")
	}

	client.muConn.Lock()
	client.wsconn = nil
	close(sessionDone)
	client.muConn.Unlock()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("reconnect did not continue after the cancelled session closed")
	}
	client.muConn.RLock()
	currentCtx = client.ctx
	client.muConn.RUnlock()
	if currentCtx != newCtx {
		t.Fatal("reconnect did not install the replacement context")
	}
}

func Test_InternalConnectKeepsActiveSessionContext(t *testing.T) {
	handler := &closeCountingHandler{}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	oldCtx, cancelOld := context.WithCancel(context.Background())
	client := NewWS(oldCtx, cancelOld, "", &interfaces.ClientOptions{SelfHosted: true}, &handlerInterface, nil)
	if client == nil {
		t.Fatal("NewWS returned nil")
	}
	client.wsconn = &websocket.Conn{}
	client.sessionDone = make(chan struct{})

	newCtx := context.Background()
	if got := client.internalConnectWithCancel(newCtx, func() {}, 1, true, nil); got == nil {
		t.Fatal("reconnect did not reuse the active transport")
	}
	client.muConn.RLock()
	currentCtx := client.ctx
	client.muConn.RUnlock()
	if currentCtx != oldCtx {
		t.Fatal("reconnect replaced the active session context")
	}
}
