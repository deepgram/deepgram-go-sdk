// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"encoding/json"
	"testing"
)

func Test_ErrorResponsePreservesRESTWireShape(t *testing.T) {
	var response ErrorResponse
	if err := json.Unmarshal([]byte(`{"type":"Error","err_code":"ENGINE_ERROR","err_msg":"engine failed"}`), &response); err != nil {
		t.Fatalf("unmarshal REST ErrorResponse: %v", err)
	}
	if response.Type != "Error" || response.ErrCode != "ENGINE_ERROR" || response.ErrMsg != "engine failed" {
		t.Fatalf("unexpected ErrorResponse: %#v", response)
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal REST ErrorResponse: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal marshaled REST ErrorResponse: %v", err)
	}
	if payload["type"] != "Error" || payload["err_code"] != "ENGINE_ERROR" {
		t.Errorf("unexpected REST payload: %s", data)
	}
	if _, found := payload["code"]; found {
		t.Errorf("REST payload must not contain code: %s", data)
	}
}
