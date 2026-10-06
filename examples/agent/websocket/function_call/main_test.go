// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
)

func Test_ExecuteFunctionCanceledDoesNotReturnResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results := make(chan functionCallResult, 1)
	failures := make(chan functionCallError, 1)
	done := make(chan struct{})
	go func() {
		executeFunction(ctx, msginterfaces.FunctionCall{ID: "fc_123", Arguments: `{"location":"Fremont"}`}, results, failures)
		close(done)
	}()
	<-done

	select {
	case result := <-results:
		t.Fatalf("canceled function returned a result: %+v", result)
	default:
	}
	select {
	case failure := <-failures:
		t.Fatalf("canceled function returned an error: %v", failure.err)
	default:
	}
}

func Test_SendFunctionResultPrefersQueuedCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inFlight := map[string]context.CancelFunc{"fc_123": cancel}
	cancellations := make(chan *msginterfaces.FunctionCallCancelledResponse, 1)
	cancellations <- &msginterfaces.FunctionCallCancelledResponse{
		Functions: []msginterfaces.FunctionCallCancelled{{ID: "fc_123"}},
	}

	sent := false
	resultSent, err := sendFunctionResult(cancellations, inFlight, functionCallResult{
		function: msginterfaces.FunctionCall{ID: "fc_123", Name: "get_weather"},
		content:  `{"weather":"sunny"}`,
	}, func(msginterfaces.FunctionCallResponse) error {
		sent = true
		return nil
	})
	if err != nil {
		t.Fatalf("send function result: %v", err)
	}
	if resultSent || sent {
		t.Fatal("canceled function sent a response")
	}
	if _, ok := inFlight["fc_123"]; ok {
		t.Fatal("canceled function remained in flight")
	}
	select {
	case <-ctx.Done():
	default:
		t.Fatal("canceled function context remained active")
	}
}
