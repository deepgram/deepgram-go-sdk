// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"sync"
	"testing"

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
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
}

func (h *closeCountingHandler) GetURL(string) (string, error) { return "", nil }
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
