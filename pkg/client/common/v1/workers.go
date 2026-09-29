// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"context"
	"sync"
)

// Workers tracks the background goroutines (keepalive, auto flush) that a client starts
// for one connection. Starting a new set first stops the previous one, so reconnecting
// does not leave duplicate goroutines running against the same client.
type Workers struct {
	mu     sync.Mutex
	cancel context.CancelFunc
}

// Start stops the previous set of workers and returns the context for the next set.
// The returned context is canceled by Stop, by the next call to Start, or when parent is done.
func (w *Workers) Start(parent context.Context) context.Context {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cancel != nil {
		w.cancel()
	}
	ctx, cancel := context.WithCancel(parent)
	w.cancel = cancel
	return ctx
}

// Stop cancels the current set of workers. It is safe to call more than once.
func (w *Workers) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
}
