// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"encoding/json"
	"testing"
)

func Test_DeepgramErrorJSONWireShapes(t *testing.T) {
	t.Run("WebSocket event", func(t *testing.T) {
		var event DeepgramError
		if err := json.Unmarshal([]byte(`{"type":"Error","code":"INVALID_QUERY","description":"invalid endpointing"}`), &event); err != nil {
			t.Fatalf("unmarshal error event: %v", err)
		}
		if event.Type != "Error" || event.Code != "INVALID_QUERY" || event.Description != "invalid endpointing" {
			t.Fatalf("unexpected event: %#v", event)
		}

		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal error event: %v", err)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(data, &payload); err != nil {
			t.Fatalf("unmarshal marshaled event: %v", err)
		}
		if payload["type"] != "Error" || payload["code"] != "INVALID_QUERY" {
			t.Errorf("unexpected WebSocket payload: %s", data)
		}
		if _, found := payload["err_code"]; found {
			t.Errorf("WebSocket payload must not contain err_code: %s", data)
		}
	})

	t.Run("REST response", func(t *testing.T) {
		var response DeepgramError
		if err := json.Unmarshal([]byte(`{"type":"Error","err_code":"ENGINE_ERROR","err_msg":"engine failed"}`), &response); err != nil {
			t.Fatalf("unmarshal REST error: %v", err)
		}
		if response.Type != "Error" || response.ErrCode != "ENGINE_ERROR" || response.ErrMsg != "engine failed" {
			t.Fatalf("unexpected REST error: %#v", response)
		}

		data, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("marshal REST error: %v", err)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(data, &payload); err != nil {
			t.Fatalf("unmarshal marshaled REST error: %v", err)
		}
		if payload["type"] != "Error" || payload["err_code"] != "ENGINE_ERROR" || payload["err_msg"] != "engine failed" {
			t.Errorf("unexpected REST payload: %s", data)
		}
		if _, found := payload["code"]; found {
			t.Errorf("REST payload must not contain code: %s", data)
		}
	})

	t.Run("REST response without message", func(t *testing.T) {
		var response DeepgramError
		if err := json.Unmarshal([]byte(`{"type":"Error","err_code":"ENGINE_ERROR"}`), &response); err != nil {
			t.Fatalf("unmarshal REST error: %v", err)
		}

		data, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("marshal REST error: %v", err)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(data, &payload); err != nil {
			t.Fatalf("unmarshal marshaled REST error: %v", err)
		}
		if payload["type"] != "Error" || payload["err_code"] != "ENGINE_ERROR" {
			t.Errorf("unexpected REST payload: %s", data)
		}
		if _, found := payload["code"]; found {
			t.Errorf("REST payload must not contain code: %s", data)
		}
	})
}

func Test_DeepgramWarningJSONWireShapes(t *testing.T) {
	var warning DeepgramWarning
	if err := json.Unmarshal([]byte(`{"type":"Warning","code":"RATE_LIMITED","description":"slow down"}`), &warning); err != nil {
		t.Fatalf("unmarshal warning: %v", err)
	}
	if warning.Type != "Warning" || warning.Code != "RATE_LIMITED" || warning.Description != "slow down" {
		t.Fatalf("unexpected warning: %#v", warning)
	}

	data, err := json.Marshal(warning)
	if err != nil {
		t.Fatalf("marshal warning: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal marshaled warning: %v", err)
	}
	if payload["type"] != "Warning" || payload["code"] != "RATE_LIMITED" {
		t.Errorf("unexpected WebSocket warning payload: %s", data)
	}
	if _, found := payload["warn_code"]; found {
		t.Errorf("WebSocket warning payload must not contain warn_code: %s", data)
	}

	var legacy DeepgramWarning
	if err := json.Unmarshal([]byte(`{"type":"Warning","warn_code":"LEGACY_WARNING","warn_msg":"legacy"}`), &legacy); err != nil {
		t.Fatalf("unmarshal legacy warning: %v", err)
	}
	if legacy.Type != "Warning" || legacy.WarnCode != "LEGACY_WARNING" || legacy.WarnMsg != "legacy" {
		t.Fatalf("unexpected legacy warning: %#v", legacy)
	}
	data, err = json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal legacy warning: %v", err)
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal marshaled legacy warning: %v", err)
	}
	if payload["type"] != "Warning" || payload["warn_code"] != "LEGACY_WARNING" || payload["warn_msg"] != "legacy" {
		t.Errorf("unexpected legacy warning payload: %s", data)
	}
}
