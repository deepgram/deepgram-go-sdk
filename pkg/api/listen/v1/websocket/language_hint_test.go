// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func Test_LanguageHint(t *testing.T) {
	t.Run("serializes repeated streaming query parameters", func(t *testing.T) {
		endpoint, err := version.GetLiveAPI(context.Background(), "", "", "", &clientinterfaces.LiveTranscriptionOptions{
			Language:     "multi",
			LanguageHint: []string{"en", "es"},
		})
		if err != nil {
			t.Fatalf("build live API URL: %v", err)
		}
		parsed, err := url.Parse(endpoint)
		if err != nil {
			t.Fatalf("parse live API URL: %v", err)
		}
		if got, want := parsed.Query()["language_hint"], []string{"en", "es"}; !reflect.DeepEqual(got, want) {
			t.Errorf("language_hint = %q, want %q", got, want)
		}
	})

	t.Run("decodes applied hints", func(t *testing.T) {
		var response interfaces.MessageResponse
		if err := json.Unmarshal([]byte(`{"type":"Results","channel":{"alternatives":[{"languages_hinted":["en","es"]}]}}`), &response); err != nil {
			t.Fatalf("decode streaming response: %v", err)
		}
		if got, want := response.Channel.Alternatives[0].LanguagesHinted, []string{"en", "es"}; !reflect.DeepEqual(got, want) {
			t.Errorf("LanguagesHinted = %q, want %q", got, want)
		}
	})

	for _, test := range []struct {
		name    string
		payload string
		wantNil bool
	}{
		{name: "empty", payload: `{"languages_hinted":[]}`},
		{name: "absent", payload: `{}`, wantNil: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var alternative interfaces.Alternative
			if err := json.Unmarshal([]byte(test.payload), &alternative); err != nil {
				t.Fatalf("decode alternative: %v", err)
			}
			if (alternative.LanguagesHinted == nil) != test.wantNil {
				t.Errorf("LanguagesHinted nil = %t, want %t", alternative.LanguagesHinted == nil, test.wantNil)
			}
		})
	}
}
