// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type customThinkProviderHandler struct {
	*agentws.DefaultChanHandler
	customFromThinkProvider chan *msginterfaces.CustomFromThinkProviderResponse
	settingsApplied         chan *msginterfaces.SettingsAppliedResponse
	errors                  chan *msginterfaces.ErrorResponse
}

func newCustomThinkProviderHandler() *customThinkProviderHandler {
	return &customThinkProviderHandler{
		DefaultChanHandler:      agentws.NewDefaultChanHandler(),
		customFromThinkProvider: make(chan *msginterfaces.CustomFromThinkProviderResponse, 1),
		settingsApplied:         make(chan *msginterfaces.SettingsAppliedResponse, 1),
		errors:                  make(chan *msginterfaces.ErrorResponse, 1),
	}
}

func (h *customThinkProviderHandler) GetCustomFromThinkProvider() []*chan *msginterfaces.CustomFromThinkProviderResponse {
	return []*chan *msginterfaces.CustomFromThinkProviderResponse{&h.customFromThinkProvider}
}

func (h *customThinkProviderHandler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settingsApplied}
}

func (h *customThinkProviderHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	endpoint := os.Getenv("CUSTOM_THINK_PROVIDER_URL")
	if !strings.HasPrefix(endpoint, "wss://") {
		return fmt.Errorf("set CUSTOM_THINK_PROVIDER_URL to a wss:// custom Think provider endpoint")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	options := agent.NewSettingsConfigurationOptions()
	options.Agent.Think.Endpoint = &interfacesv1.Endpoint{Url: endpoint}
	handler := newCustomThinkProviderHandler()

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
		return fmt.Errorf("apply settings: %s", response.Description)
	case <-ctx.Done():
		return fmt.Errorf("wait for SettingsApplied: %w", ctx.Err())
	}

	if err := conn.SendCustomToThinkProvider(json.RawMessage(`{"event":"status"}`)); err != nil {
		return fmt.Errorf("send custom Think provider message: %w", err)
	}

	select {
	case response := <-handler.customFromThinkProvider:
		fmt.Printf("Custom Think provider response: %s\n", response.Content)
		return nil
	case response := <-handler.errors:
		return fmt.Errorf("agent error: %s", response.Description)
	case <-ctx.Done():
		return fmt.Errorf("wait for custom Think provider response: %w", ctx.Err())
	}
}
