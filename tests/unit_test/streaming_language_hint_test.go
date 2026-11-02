// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	listenws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// liveQuery builds the Listen v1 WebSocket URL for the given options and
// returns its parsed query string.
func liveQuery(t *testing.T, options *interfaces.LiveTranscriptionOptions) url.Values {
	t.Helper()

	uri, err := version.GetLiveAPI(context.Background(), "", "", "", options)
	if err != nil {
		t.Fatalf("GetLiveAPI should succeed, but got %s", err)
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		t.Fatalf("failed to parse generated URI %q: %s", uri, err)
	}
	return parsed.Query()
}

// Test_LiveLanguageHint verifies that LanguageHint on the Nova streaming options
// serializes to repeated language_hint query parameters, and that the
// languages_hinted response field decodes from Results messages.
func Test_LiveLanguageHint(t *testing.T) {
	t.Run("single hint appears in the streaming query string", func(t *testing.T) {
		q := liveQuery(t, &interfaces.LiveTranscriptionOptions{
			Model:        "nova-3",
			Language:     "multi",
			LanguageHint: []string{"es"},
		})

		if got := q.Get("language"); got != "multi" {
			t.Errorf("expected language=multi, got %q", got)
		}
		if got := q["language_hint"]; !reflect.DeepEqual(got, []string{"es"}) {
			t.Errorf("expected language_hint=[es], got %v", got)
		}
	})

	t.Run("multiple hints are sent as repeated keys", func(t *testing.T) {
		q := liveQuery(t, &interfaces.LiveTranscriptionOptions{
			Model:        "nova-3",
			Language:     "multi",
			LanguageHint: []string{"en", "es"},
		})

		if got := q["language_hint"]; !reflect.DeepEqual(got, []string{"en", "es"}) {
			t.Errorf("expected language_hint=[en es] as repeated keys, got %v", got)
		}
	})

	t.Run("no hint omits the query parameter", func(t *testing.T) {
		q := liveQuery(t, &interfaces.LiveTranscriptionOptions{
			Model:    "nova-3",
			Language: "multi",
		})

		if _, ok := q["language_hint"]; ok {
			t.Errorf("expected no language_hint parameter, got %v", q["language_hint"])
		}
	})

	t.Run("decodes languages_hinted, languages and word language", func(t *testing.T) {
		var response listenws.MessageResponse
		if err := json.Unmarshal([]byte(`{
			"type":"Results",
			"is_final":true,
			"channel":{"alternatives":[{
				"transcript":"hola mundo",
				"confidence":0.98,
				"languages":["es","en"],
				"languages_hinted":["en","es"],
				"words":[{"word":"hola","start":0.1,"end":0.4,"confidence":0.99,"language":"es"}]
			}]}
		}`), &response); err != nil {
			t.Fatalf("decode live message: %v", err)
		}

		if len(response.Channel.Alternatives) != 1 {
			t.Fatalf("expected one alternative, got %d", len(response.Channel.Alternatives))
		}
		alternative := response.Channel.Alternatives[0]
		if !reflect.DeepEqual(alternative.LanguagesHinted, []string{"en", "es"}) {
			t.Errorf("expected languages_hinted=[en es], got %v", alternative.LanguagesHinted)
		}
		if !reflect.DeepEqual(alternative.Languages, []string{"es", "en"}) {
			t.Errorf("expected languages=[es en], got %v", alternative.Languages)
		}
		if len(alternative.Words) != 1 || alternative.Words[0].Language != "es" {
			t.Errorf("expected word language es, got %#v", alternative.Words)
		}
	})

	t.Run("empty languages_hinted decodes to a non-nil empty slice", func(t *testing.T) {
		var response listenws.MessageResponse
		if err := json.Unmarshal([]byte(`{
			"type":"Results",
			"channel":{"alternatives":[{"transcript":"hello","languages_hinted":[]}]}
		}`), &response); err != nil {
			t.Fatalf("decode live message: %v", err)
		}

		hinted := response.Channel.Alternatives[0].LanguagesHinted
		if hinted == nil || len(hinted) != 0 {
			t.Errorf("expected non-nil empty languages_hinted, got %#v", hinted)
		}
	})
}
