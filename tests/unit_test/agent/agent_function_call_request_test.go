// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
)

// Test_AgentFunctionCallRequestResponse verifies that the functions array the
// Agent API sends on a FunctionCallRequest event lands on the struct instead of
// being silently discarded.
func Test_AgentFunctionCallRequestResponse(t *testing.T) {
	// Captured from a live /v1/agent/converse socket.
	const payload = `{"type":"FunctionCallRequest","functions":[` +
		`{"id":"call_0X6nawtijkWevoU1XGuMKSQx","name":"set_count","arguments":"{\"count\":7}","client_side":true}]}`

	var msg msginterfaces.FunctionCallRequestResponse
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		t.Fatalf("failed to unmarshal FunctionCallRequest payload: %s", err)
	}

	if msg.Type != "FunctionCallRequest" {
		t.Errorf("expected type FunctionCallRequest, got %q", msg.Type)
	}
	if len(msg.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d: %+v", len(msg.Functions), msg)
	}

	fn := msg.Functions[0]
	if fn.ID != "call_0X6nawtijkWevoU1XGuMKSQx" || fn.Name != "set_count" || !fn.ClientSide {
		t.Errorf("unexpected function: %+v", fn)
	}

	var args struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal([]byte(fn.Arguments), &args); err != nil {
		t.Fatalf("arguments is not a JSON-encoded object: %s", err)
	}
	if args.Count != 7 {
		t.Errorf("expected count 7, got %d", args.Count)
	}

	// Re-marshaling must keep the functions so a proxy can forward the event.
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal FunctionCallRequest: %s", err)
	}
	var round map[string]any
	if err := json.Unmarshal(data, &round); err != nil {
		t.Fatalf("failed to unmarshal marshaled FunctionCallRequest: %s", err)
	}
	if _, ok := round["functions"]; !ok {
		t.Errorf("expected functions key after round trip, got %s", data)
	}
}
