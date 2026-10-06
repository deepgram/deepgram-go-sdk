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
