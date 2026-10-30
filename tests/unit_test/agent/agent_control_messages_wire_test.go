// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	websocketv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent/v1/websocket"
)

func assertWireFields(t *testing.T, v interface{}, want map[string]string, absent ...string) {
	t.Helper()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal %T: %v", v, err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal %s: %v", data, err)
	}

	for key, value := range want {
		raw, ok := got[key]
		if !ok {
			t.Errorf("%T JSON %s is missing %q", v, data, key)
			continue
		}
		if raw != value {
			t.Errorf("%T JSON %s has %q = %v, want %q", v, data, key, raw, value)
		}
	}
	for _, key := range absent {
		if _, ok := got[key]; ok {
			t.Errorf("%T JSON %s must not contain %q", v, data, key)
		}
	}
}

func Test_UpdatePromptWireShape(t *testing.T) {
	const prompt = "Always answer in French."

	t.Run("interfaces type", func(t *testing.T) {
		assertWireFields(t, msginterfaces.UpdatePrompt{
			Type:   msginterfaces.TypeUpdatePrompt,
			Prompt: prompt,
		}, map[string]string{"type": "UpdatePrompt", "prompt": prompt}, "instructions")
	})

	t.Run("client type", func(t *testing.T) {
		assertWireFields(t, websocketv1.UpdatePrompt{
			Type:   msginterfaces.TypeUpdatePrompt,
			Prompt: prompt,
		}, map[string]string{"type": "UpdatePrompt", "prompt": prompt}, "instructions")
	})

	// The API closes the session when prompt is missing but accepts an empty string.
	t.Run("empty prompt is still sent", func(t *testing.T) {
		assertWireFields(t, websocketv1.UpdatePrompt{
			Type: msginterfaces.TypeUpdatePrompt,
		}, map[string]string{"type": "UpdatePrompt", "prompt": ""}, "instructions")
	})
}

func Test_InjectAgentMessageWireShape(t *testing.T) {
	const message = "Say the word pineapple."

	t.Run("interfaces type", func(t *testing.T) {
		assertWireFields(t, msginterfaces.InjectAgentMessage{
			Type:    msginterfaces.TypeInjectAgentMessage,
			Content: message,
		}, map[string]string{"type": "InjectAgentMessage", "message": message}, "content")
	})

	t.Run("client type", func(t *testing.T) {
		assertWireFields(t, websocketv1.InjectAgentMessage{
			Type:    msginterfaces.TypeInjectAgentMessage,
			Content: message,
		}, map[string]string{"type": "InjectAgentMessage", "message": message}, "content")
	})

	// The API closes the session when message is missing but accepts an empty string.
	t.Run("empty message is still sent", func(t *testing.T) {
		assertWireFields(t, websocketv1.InjectAgentMessage{
			Type: msginterfaces.TypeInjectAgentMessage,
		}, map[string]string{"type": "InjectAgentMessage", "message": ""}, "content")
	})
}
