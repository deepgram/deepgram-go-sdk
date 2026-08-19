// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"reflect"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v1/websocket/interfaces"
)

// Test_SpeakWebSocketMetadataResponse verifies that every field Deepgram sends
// on a speak v1 Metadata event survives unmarshalling. Prior to this the struct
// only declared type and request_id, so the three model_* fields were silently
// discarded by encoding/json.
func Test_SpeakWebSocketMetadataResponse(t *testing.T) {
	// Captured from a live /v1/speak socket.
	const payload = `{"type":"Metadata","request_id":"01a01a7f-f079-74c3-bfc5-ac08f549baae",` +
		`"model_name":"aura-asteria-en","model_version":"2024-11-19.0",` +
		`"model_uuid":"ecb76e9d-f2db-4127-8060-79b05590d22f","additional_model_uuids":[]}`

	t.Run("all fields on the wire land on the struct", func(t *testing.T) {
		var md msginterfaces.MetadataResponse
		if err := json.Unmarshal([]byte(payload), &md); err != nil {
			t.Fatalf("failed to unmarshal Metadata payload: %s", err)
		}

		expected := msginterfaces.MetadataResponse{
			Type:         "Metadata",
			RequestID:    "01a01a7f-f079-74c3-bfc5-ac08f549baae",
			ModelName:    "aura-asteria-en",
			ModelVersion: "2024-11-19.0",
			ModelUUID:    "ecb76e9d-f2db-4127-8060-79b05590d22f",
		}
		// additional_model_uuids is [] on the wire, which unmarshals to an
		// empty non-nil slice, so compare it separately from the scalars.
		actual := md
		actual.AdditionalModelUUIDs = nil

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("unmarshalled Metadata mismatch\n got: %+v\nwant: %+v", actual, expected)
		}
		if md.AdditionalModelUUIDs == nil || len(md.AdditionalModelUUIDs) != 0 {
			t.Errorf("expected AdditionalModelUUIDs to be an empty slice, got %#v", md.AdditionalModelUUIDs)
		}
	})

	t.Run("additional_model_uuids is populated when present", func(t *testing.T) {
		const withUUIDs = `{"type":"Metadata","request_id":"req","model_name":"aura-2-thalia-en",` +
			`"model_version":"1.0","model_uuid":"uuid-a","additional_model_uuids":["uuid-b","uuid-c"]}`

		var md msginterfaces.MetadataResponse
		if err := json.Unmarshal([]byte(withUUIDs), &md); err != nil {
			t.Fatalf("failed to unmarshal Metadata payload: %s", err)
		}
		if want := []string{"uuid-b", "uuid-c"}; !reflect.DeepEqual(md.AdditionalModelUUIDs, want) {
			t.Errorf("expected AdditionalModelUUIDs %v, got %v", want, md.AdditionalModelUUIDs)
		}
	})

	// The required fields must be emitted even when zero-valued so a re-marshalled
	// struct stays contract-valid; additional_model_uuids is optional and stays omitted.
	t.Run("required fields are emitted when zero-valued", func(t *testing.T) {
		data, err := json.Marshal(msginterfaces.MetadataResponse{})
		if err != nil {
			t.Fatalf("failed to marshal Metadata: %s", err)
		}

		var got map[string]any
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("failed to unmarshal marshalled Metadata: %s", err)
		}

		for _, key := range []string{"type", "request_id", "model_name", "model_version", "model_uuid"} {
			if _, ok := got[key]; !ok {
				t.Errorf("expected required key %q in marshalled Metadata, got %s", key, data)
			}
		}
		if _, ok := got["additional_model_uuids"]; ok {
			t.Errorf("did not expect optional key additional_model_uuids when unset, got %s", data)
		}
	})
}

// Test_SpeakWebSocketAckSequenceID verifies that sequence_id survives a
// re-marshal of the Flushed and Cleared acks. Deepgram numbers the first ack in
// a session sequence_id: 0, which omitempty used to drop from the output.
func Test_SpeakWebSocketAckSequenceID(t *testing.T) {
	t.Run("FlushedResponse keeps sequence_id 0", func(t *testing.T) {
		data, err := json.Marshal(msginterfaces.FlushedResponse{Type: "Flushed", SequenceID: 0})
		if err != nil {
			t.Fatalf("failed to marshal FlushedResponse: %s", err)
		}
		if want := `{"type":"Flushed","sequence_id":0}`; string(data) != want {
			t.Errorf("expected %s, got %s", want, data)
		}
	})

	t.Run("ClearedResponse keeps sequence_id 0", func(t *testing.T) {
		// Captured from a live /v1/speak socket.
		const payload = `{"type":"Cleared","sequence_id":0}`

		var cr msginterfaces.ClearedResponse
		if err := json.Unmarshal([]byte(payload), &cr); err != nil {
			t.Fatalf("failed to unmarshal Cleared payload: %s", err)
		}
		if cr.Type != "Cleared" || cr.SequenceID != 0 {
			t.Errorf("unexpected ClearedResponse: %+v", cr)
		}

		data, err := json.Marshal(cr)
		if err != nil {
			t.Fatalf("failed to marshal ClearedResponse: %s", err)
		}
		if string(data) != payload {
			t.Errorf("expected round-trip to preserve %s, got %s", payload, data)
		}
	})
}
