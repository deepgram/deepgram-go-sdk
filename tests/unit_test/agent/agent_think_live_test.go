// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const frustrationDetectedPayload = `{
	"type": "FrustrationDetected",
	"count": 2,
	"escalated": true
}`

const thinkLiveConversationPayload = `{
	"type": "ConversationText",
	"role": "assistant",
	"conversational_behavior": "filler",
	"content": "Let me check that for you."
}`

type thinkLiveHandler struct {
	*legacyFunctionCallHandler
	conversationTexts   chan *msginterfaces.ConversationTextResponse
	frustrationDetected chan *msginterfaces.FrustrationDetectedResponse
}

func Test_ThinkLiveRestrictedTermsMarshaling(t *testing.T) {
	for _, test := range []struct {
		name  string
		terms []string
		want  string
	}{
		{name: "nil omits", want: `"type":"deepgram"`},
		{name: "empty preserves", terms: []string{}, want: `"restricted_terms":[]`},
		{name: "values preserve", terms: []string{"human"}, want: `"restricted_terms":["human"]`},
	} {
		t.Run(test.name, func(t *testing.T) {
			provider := interfacesv1.NewThinkLiveSettings().Provider
			provider.RestrictedTerms = test.terms
			data, err := json.Marshal(provider)
			if err != nil {
				t.Fatalf("marshal provider: %v", err)
			}
			if !strings.Contains(string(data), test.want) {
				t.Errorf("payload %s does not contain %s", data, test.want)
			}
		})
	}
}

func newThinkLiveHandler() *thinkLiveHandler {
	return &thinkLiveHandler{
		legacyFunctionCallHandler: newLegacyFunctionCallHandler(),
		conversationTexts:         make(chan *msginterfaces.ConversationTextResponse, 1),
		frustrationDetected:       make(chan *msginterfaces.FrustrationDetectedResponse, 1),
	}
}

func (h *thinkLiveHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&h.conversationTexts}
}

func (h *thinkLiveHandler) GetFrustrationDetected() []*chan *msginterfaces.FrustrationDetectedResponse {
	return []*chan *msginterfaces.FrustrationDetectedResponse{&h.frustrationDetected}
}

func Test_ThinkLiveSettingsMarshaling(t *testing.T) {
	sensitivity := 0.0
	presenceCheckMaxCount := 0
	presenceCheckIntervalMS := 5000
	frustrationMaxCount := 0

	settings := interfacesv1.NewThinkLiveSettings()
	settings.Provider.AgentBackchannelSensitivity = &sensitivity
	settings.Provider.PresenceCheckMaxCount = &presenceCheckMaxCount
	settings.Provider.PresenceCheckIntervalMS = &presenceCheckIntervalMS
	settings.Provider.FrustrationMaxCount = &frustrationMaxCount
	settings.Provider.RestrictedTerms = []string{"human", "transfer"}

	data, err := json.Marshal(struct {
		Agent interfacesv1.Agent `json:"agent"`
	}{
		Agent: interfacesv1.Agent{ThinkLive: settings},
	})
	if err != nil {
		t.Fatalf("marshal Think Live settings envelope: %v", err)
	}

	const expected = `{"provider":{"type":"deepgram","version":"v1","agent_backchannel_sensitivity":0,"presence_check_max_count":0,"presence_check_interval_ms":5000,"frustration_max_count":0,"restricted_terms":["human","transfer"]}}`
	var payload struct {
		Agent map[string]json.RawMessage `json:"agent"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal Think Live settings envelope: %v", err)
	}
	if string(payload.Agent["think_live"]) != expected {
		t.Errorf("unexpected Think Live settings payload:\nwant %s\n got %s", expected, payload.Agent["think_live"])
	}
}

func Test_ThinkLiveEventsRouting(t *testing.T) {
	handler := newThinkLiveHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(frustrationDetectedPayload)); err != nil {
		t.Fatalf("route FrustrationDetected: %v", err)
	}
	if err := router.Message([]byte(thinkLiveConversationPayload)); err != nil {
		t.Fatalf("route ConversationText: %v", err)
	}

	select {
	case event := <-handler.frustrationDetected:
		if event.Type != msginterfaces.TypeFrustrationDetectedResponse || event.Count != 2 || !event.Escalated {
			t.Errorf("unexpected FrustrationDetected event: %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("FrustrationDetected was not delivered to the handler")
	}

	select {
	case event := <-handler.conversationTexts:
		if event.ConversationalBehavior != "filler" {
			t.Errorf("unexpected conversational behavior: %q", event.ConversationalBehavior)
		}
	case <-time.After(time.Second):
		t.Fatal("ConversationText was not delivered to the handler")
	}
}

func Test_ThinkLiveFrustrationFallsBackToUnhandled(t *testing.T) {
	handler := newLegacyFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(frustrationDetectedPayload)); err != nil {
		t.Fatalf("FrustrationDetected must fall back successfully: %v", err)
	}

	select {
	case message := <-handler.unhandled:
		if string(*message) != frustrationDetectedPayload {
			t.Errorf("unexpected unhandled message: %s", *message)
		}
	case <-time.After(time.Second):
		t.Fatal("FrustrationDetected was not delivered to the unhandled handler")
	}
}
