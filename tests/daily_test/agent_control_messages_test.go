// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package daily_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	websocketv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent/v1/websocket"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

type controlMessagesDailyHandler struct {
	welcome          chan *msginterfaces.WelcomeResponse
	conversationText chan *msginterfaces.ConversationTextResponse
	errors           chan *msginterfaces.ErrorResponse
	unhandled        chan *[]byte
	settingsApplied  chan *msginterfaces.SettingsAppliedResponse
}

func newControlMessagesDailyHandler() *controlMessagesDailyHandler {
	return &controlMessagesDailyHandler{
		welcome:          make(chan *msginterfaces.WelcomeResponse, 1),
		conversationText: make(chan *msginterfaces.ConversationTextResponse, 16),
		errors:           make(chan *msginterfaces.ErrorResponse, 4),
		unhandled:        make(chan *[]byte, 64),
		settingsApplied:  make(chan *msginterfaces.SettingsAppliedResponse, 1),
	}
}

func (h *controlMessagesDailyHandler) GetBinary() []*chan *[]byte { return nil }
func (h *controlMessagesDailyHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return []*chan *msginterfaces.WelcomeResponse{&h.welcome}
}
func (h *controlMessagesDailyHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&h.conversationText}
}
func (h *controlMessagesDailyHandler) GetUserStartedSpeaking() []*chan *msginterfaces.UserStartedSpeakingResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetAgentThinking() []*chan *msginterfaces.AgentThinkingResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetFunctionCallRequest() []*chan *msginterfaces.FunctionCallRequestResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetAgentStartedSpeaking() []*chan *msginterfaces.AgentStartedSpeakingResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetAgentAudioDone() []*chan *msginterfaces.AgentAudioDoneResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetClose() []*chan *msginterfaces.CloseResponse { return nil }
func (h *controlMessagesDailyHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}
func (h *controlMessagesDailyHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandled}
}
func (h *controlMessagesDailyHandler) GetInjectionRefused() []*chan *msginterfaces.InjectionRefusedResponse {
	return nil
}
func (h *controlMessagesDailyHandler) GetKeepAlive() []*chan *msginterfaces.KeepAlive { return nil }
func (h *controlMessagesDailyHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsApplied}
}

// TestDaily_AgentUpdatePromptAndInjectAgentMessage checks that the API accepts
// the SDK's UpdatePrompt and InjectAgentMessage frames. Before the fix, the
// UpdatePrompt frame failed with UNPARSABLE_CLIENT_MESSAGE and closed the session.
func TestDaily_AgentUpdatePromptAndInjectAgentMessage(t *testing.T) {
	if os.Getenv("DEEPGRAM_API_KEY") == "" && os.Getenv("DEEPGRAM_ACCESS_TOKEN") == "" {
		t.Skip("DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	options := agent.NewSettingsConfigurationOptions()
	options.Agent.Greeting = ""

	handler := newControlMessagesDailyHandler()
	conn, err := agent.NewWSUsingChan(ctx, "", &interfaces.ClientOptions{}, options, handler)
	if err != nil {
		t.Fatalf("create Agent WebSocket: %v", err)
	}
	defer conn.Stop()

	if !conn.Connect() {
		t.Fatal("connect Agent WebSocket")
	}

	// waitFor drains every routed channel so the router never blocks, and
	// returns once done reports true.
	waitFor := func(step string, done func(text *msginterfaces.ConversationTextResponse, raw []byte) bool) {
		t.Helper()
		for {
			select {
			case welcome := <-handler.welcome:
				t.Logf("request_id: %s", welcome.RequestID)
			case <-handler.settingsApplied:
				if done(nil, []byte(`{"type":"SettingsApplied"}`)) {
					return
				}
			case text := <-handler.conversationText:
				if done(text, nil) {
					return
				}
			case raw := <-handler.unhandled:
				if done(nil, *raw) {
					return
				}
			case apiErr := <-handler.errors:
				t.Fatalf("%s: Agent error: %s", step, apiErr.Description)
			case <-ctx.Done():
				t.Fatalf("%s: %v", step, ctx.Err())
			}
		}
	}
	isType := func(raw []byte, want string) bool {
		var msg struct {
			Type string `json:"type"`
		}
		return raw != nil && json.Unmarshal(raw, &msg) == nil && msg.Type == want
	}

	waitFor("wait for SettingsApplied", func(_ *msginterfaces.ConversationTextResponse, raw []byte) bool {
		return isType(raw, "SettingsApplied")
	})

	if err := conn.WriteJSON(websocketv1.UpdatePrompt{
		Type:   msginterfaces.TypeUpdatePrompt,
		Prompt: "Keep every answer under ten words.",
	}); err != nil {
		t.Fatalf("send UpdatePrompt: %v", err)
	}
	waitFor("wait for PromptUpdated", func(_ *msginterfaces.ConversationTextResponse, raw []byte) bool {
		return isType(raw, "PromptUpdated")
	})

	const message = "Say the word pineapple."
	if err := conn.WriteJSON(msginterfaces.InjectAgentMessage{
		Type:    msginterfaces.TypeInjectAgentMessage,
		Content: message,
	}); err != nil {
		t.Fatalf("send InjectAgentMessage: %v", err)
	}
	waitFor("wait for injected assistant text", func(text *msginterfaces.ConversationTextResponse, _ []byte) bool {
		return text != nil && text.Role == "assistant" && text.Content == message
	})
}
