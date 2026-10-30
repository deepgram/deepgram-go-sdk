// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This text-only example requires a Think Live Early Access project.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type thinkLiveHandler struct {
	*agentws.DefaultChanHandler
	settingsApplied     chan *msginterfaces.SettingsAppliedResponse
	conversationTexts   chan *msginterfaces.ConversationTextResponse
	frustrationDetected chan *msginterfaces.FrustrationDetectedResponse
	errors              chan *msginterfaces.ErrorResponse
}

func newThinkLiveHandler() *thinkLiveHandler {
	return &thinkLiveHandler{
		DefaultChanHandler:  agentws.NewDefaultChanHandler(),
		settingsApplied:     make(chan *msginterfaces.SettingsAppliedResponse, 1),
		conversationTexts:   make(chan *msginterfaces.ConversationTextResponse, 4),
		frustrationDetected: make(chan *msginterfaces.FrustrationDetectedResponse, 1),
		errors:              make(chan *msginterfaces.ErrorResponse, 1),
	}
}

func (h *thinkLiveHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsApplied}
}

func (h *thinkLiveHandler) GetConversationText() []*chan *msginterfaces.ConversationTextResponse {
	return []*chan *msginterfaces.ConversationTextResponse{&h.conversationTexts}
}

func (h *thinkLiveHandler) GetFrustrationDetected() []*chan *msginterfaces.FrustrationDetectedResponse {
	return []*chan *msginterfaces.FrustrationDetectedResponse{&h.frustrationDetected}
}

func (h *thinkLiveHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sensitivity := 0.2
	presenceCheckMaxCount := 2
	presenceCheckIntervalMS := 5000
	frustrationMaxCount := 2
	options := agent.NewSettingsConfigurationOptions()
	options.Agent.Listen.Provider = map[string]interface{}{
		"type":    "deepgram",
		"model":   "flux-general-en",
		"version": "v2",
	}
	options.Agent.Think.Prompt = "You are Stella, a concise customer support agent."
	options.Agent.Speak.Provider = map[string]interface{}{
		"type":    "deepgram",
		"model":   "flux-kit-en",
		"version": "v2",
	}
	options.Agent.ThinkLive = interfacesv1.NewThinkLiveSettings()
	options.Agent.ThinkLive.Provider.AgentBackchannelSensitivity = &sensitivity
	options.Agent.ThinkLive.Provider.PresenceCheckMaxCount = &presenceCheckMaxCount
	options.Agent.ThinkLive.Provider.PresenceCheckIntervalMS = &presenceCheckIntervalMS
	options.Agent.ThinkLive.Provider.FrustrationMaxCount = &frustrationMaxCount
	options.Agent.ThinkLive.Provider.RestrictedTerms = []string{"human", "transfer"}

	handler := newThinkLiveHandler()
	conn, err := agent.NewWSUsingChanWithDefaults(ctx, options, handler)
	if err != nil {
		return fmt.Errorf("create Agent WebSocket: %w", err)
	}
	defer conn.Stop()

	if !conn.Connect() {
		return fmt.Errorf("connect Agent WebSocket")
	}

	select {
	case <-handler.settingsApplied:
	case response := <-handler.errors:
		return fmt.Errorf("apply Think Live settings: %s", response.Description)
	case <-ctx.Done():
		return fmt.Errorf("wait for SettingsApplied: %w", ctx.Err())
	}

	if err := conn.WriteJSON(msginterfaces.InjectUserMessage{
		Type:    msginterfaces.TypeInjectUserMessage,
		Content: "Please tell me one short fact about the first all-female spacewalk.",
	}); err != nil {
		return fmt.Errorf("send InjectUserMessage: %w", err)
	}

	for {
		select {
		case event := <-handler.frustrationDetected:
			fmt.Printf("Frustration detected: count=%d escalated=%t\n", event.Count, event.Escalated)
		case event := <-handler.conversationTexts:
			behavior := event.ConversationalBehavior
			if behavior == "" {
				behavior = "default"
			}
			fmt.Printf("%s (%s): %s\n", event.Role, behavior, event.Content)
			if event.Role == "assistant" && behavior == "default" {
				return nil
			}
		case response := <-handler.errors:
			return fmt.Errorf("agent error: %s", response.Description)
		case <-ctx.Done():
			return fmt.Errorf("complete Think Live turn: %w", ctx.Err())
		}
	}
}
