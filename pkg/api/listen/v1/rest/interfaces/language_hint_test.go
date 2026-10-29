// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func Test_LanguageHint(t *testing.T) {
	t.Run("serializes repeated prerecorded query parameters", func(t *testing.T) {
		endpoint, err := version.GetPrerecordedAPI(context.Background(), "", "", "", &interfaces.PreRecordedTranscriptionOptions{
			Language:     "multi",
			LanguageHint: []string{"en", "es"},
		})
		if err != nil {
			t.Fatalf("build prerecorded API URL: %v", err)
		}
		parsed, err := url.Parse(endpoint)
		if err != nil {
			t.Fatalf("parse endpoint: %v", err)
		}
		if got, want := parsed.Query()["language_hint"], []string{"en", "es"}; !reflect.DeepEqual(got, want) {
			t.Errorf("language_hint = %q, want %q", got, want)
		}
	})

	t.Run("decodes applied hints and utterance languages", func(t *testing.T) {
		var response PreRecordedResponse
		if err := json.Unmarshal([]byte(`{"results":{"channels":[{"alternatives":[{"languages_hinted":["en","es"]}]}],"utterances":[{"languages":["en","es"]}]}}`), &response); err != nil {
			t.Fatalf("decode prerecorded response: %v", err)
		}
		if got, want := response.Results.Channels[0].Alternatives[0].LanguagesHinted, []string{"en", "es"}; !reflect.DeepEqual(got, want) {
			t.Errorf("LanguagesHinted = %q, want %q", got, want)
		}
		if got, want := response.Results.Utterances[0].Languages, []string{"en", "es"}; !reflect.DeepEqual(got, want) {
			t.Errorf("Utterance.Languages = %q, want %q", got, want)
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
			var alternative Alternative
			if err := json.Unmarshal([]byte(test.payload), &alternative); err != nil {
				t.Fatalf("decode alternative: %v", err)
			}
			if (alternative.LanguagesHinted == nil) != test.wantNil {
				t.Errorf("LanguagesHinted nil = %t, want %t", alternative.LanguagesHinted == nil, test.wantNil)
			}
		})
	}
}
