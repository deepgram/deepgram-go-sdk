// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	prerecorded "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v1/rest"
)

func Test_PreRecordedLanguageHintWire(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Query()["language_hint"], []string{"en", "es"}; !reflect.DeepEqual(got, want) {
			t.Errorf("language_hint = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"metadata":{},"results":{"channels":[]}}`))
	}))
	t.Cleanup(server.Close)

	dg := prerecorded.New(client.New("test-api-key", &interfaces.ClientOptions{Host: server.URL}))
	if _, err := dg.FromURL(context.Background(), "https://example.com/audio.wav", &interfaces.PreRecordedTranscriptionOptions{
		Model:        "nova-3",
		Language:     "multi",
		LanguageHint: []string{"en", "es"},
	}); err != nil {
		t.Fatalf("FromURL failed: %v", err)
	}
}
