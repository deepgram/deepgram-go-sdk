// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"net/url"
	"strings"
	"testing"

	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// Test_PrerecordedDiarizeModel verifies that setting DiarizeModel on the
// pre-recorded (batch) transcription options serializes to the diarize_model
// query parameter on the outgoing POST /v1/listen request.
func Test_PrerecordedDiarizeModel(t *testing.T) {
	t.Run("diarize_model=v2 appears in the batch request query string", func(t *testing.T) {
		uri, err := version.GetPrerecordedAPI(
			context.Background(),
			"",
			"",
			"",
			&interfaces.PreRecordedTranscriptionOptions{
				Model:        "nova-2",
				DiarizeModel: "v2",
			})
		if err != nil {
			t.Fatalf("GetPrerecordedAPI should succeed, but got %s", err)
		}

		parsed, err := url.Parse(uri)
		if err != nil {
			t.Fatalf("failed to parse generated URI %q: %s", uri, err)
		}

		if got := parsed.Query().Get("diarize_model"); got != "v2" {
			t.Errorf("expected diarize_model=v2 in query string, got %q (full URI: %s)", got, uri)
		}
	})

	// v1 and latest are also valid per the API spec (enum: latest, v1, v2).
	for _, value := range []string{"v1", "latest"} {
		value := value
		t.Run("diarize_model="+value+" appears in the batch request query string", func(t *testing.T) {
			uri, err := version.GetPrerecordedAPI(
				context.Background(),
				"",
				"",
				"",
				&interfaces.PreRecordedTranscriptionOptions{
					DiarizeModel: value,
				})
			if err != nil {
				t.Fatalf("GetPrerecordedAPI should succeed, but got %s", err)
			}
			if !strings.Contains(uri, "diarize_model="+value) {
				t.Errorf("expected diarize_model=%s in query string, got URI: %s", value, uri)
			}
		})
	}

	// When DiarizeModel is unset it must not appear (omitempty), keeping the
	// change additive and non-breaking for existing callers.
	t.Run("diarize_model is omitted when unset", func(t *testing.T) {
		uri, err := version.GetPrerecordedAPI(
			context.Background(),
			"",
			"",
			"",
			&interfaces.PreRecordedTranscriptionOptions{
				Model: "nova-2",
			})
		if err != nil {
			t.Fatalf("GetPrerecordedAPI should succeed, but got %s", err)
		}
		if strings.Contains(uri, "diarize_model") {
			t.Errorf("did not expect diarize_model in query string when unset, got URI: %s", uri)
		}
	})
}
