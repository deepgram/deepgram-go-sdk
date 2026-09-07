// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"encoding/json"
	"testing"
)

func TestClearControlMessage(t *testing.T) {
	if MessageTypeReset != MessageTypeClear {
		t.Fatalf("MessageTypeReset = %q, want %q", MessageTypeReset, MessageTypeClear)
	}

	data, err := json.Marshal(controlMessage{Type: MessageTypeClear})
	if err != nil {
		t.Fatalf("marshal clear control message: %v", err)
	}

	if got, want := string(data), `{"type":"Clear"}`; got != want {
		t.Errorf("clear payload = %s, want %s", got, want)
	}
}
