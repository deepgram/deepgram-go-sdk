// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func Test_LiveDetectEntities(t *testing.T) {
	t.Run("serializes the streaming query parameter", func(t *testing.T) {
		endpoint, err := version.GetLiveAPI(
			context.Background(),
			"",
			"",
			"",
			&clientinterfaces.LiveTranscriptionOptions{DetectEntities: true},
		)
		if err != nil {
			t.Fatalf("build live API URL: %v", err)
		}

		parsed, err := url.Parse(endpoint)
		if err != nil {
			t.Fatalf("parse live API URL: %v", err)
		}
		if got := parsed.Query().Get("detect_entities"); got != "true" {
			t.Fatalf("expected detect_entities=true, got %q", got)
		}
	})

	t.Run("decodes top-level entities", func(t *testing.T) {
		var response interfaces.MessageResponse
		if err := json.Unmarshal([]byte(`{
			"type":"Results",
			"entities":[{
				"label":"PERSON",
				"value":"Ada Lovelace",
				"raw_value":"ada lovelace",
				"confidence":0.98,
				"start_word":2,
				"end_word":4
			}]
		}`), &response); err != nil {
			t.Fatalf("decode live message: %v", err)
		}

		if len(response.Entities) != 1 {
			t.Fatalf("expected one entity, got %d", len(response.Entities))
		}
		entity := response.Entities[0]
		if entity.Label != "PERSON" || entity.Value != "Ada Lovelace" || entity.RawValue != "ada lovelace" {
			t.Fatalf("unexpected entity: %#v", entity)
		}
		if entity.Confidence != 0.98 || entity.StartWord != 2 || entity.EndWord != 4 {
			t.Fatalf("unexpected entity metadata: %#v", entity)
		}
	})
}
