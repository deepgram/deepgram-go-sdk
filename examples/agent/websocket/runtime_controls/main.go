// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// runtime_controls updates a Flux end-of-turn setting after SettingsApplied and
// waits for the server's ListenUpdated acknowledgement. It sends no audio.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type handler struct {
	*agentws.DefaultChanHandler
	welcome       chan *msginterfaces.WelcomeResponse
	settings      chan *msginterfaces.SettingsAppliedResponse
	listenUpdated chan *msginterfaces.ListenUpdatedResponse
}

func newHandler() *handler {
	return &handler{
		DefaultChanHandler: agentws.NewDefaultChanHandler(),
		welcome:            make(chan *msginterfaces.WelcomeResponse, 1),
		settings:           make(chan *msginterfaces.SettingsAppliedResponse, 1),
		listenUpdated:      make(chan *msginterfaces.ListenUpdatedResponse, 1),
	}
}

func (h *handler) GetWelcome() []*chan *msginterfaces.WelcomeResponse {
	return []*chan *msginterfaces.WelcomeResponse{&h.welcome}
}

func (h *handler) GetSettingsApplied() []*chan *msginterfaces.SettingsAppliedResponse {
	return []*chan *msginterfaces.SettingsAppliedResponse{&h.settings}
}

func (h *handler) GetListenUpdated() []*chan *msginterfaces.ListenUpdatedResponse {
	return []*chan *msginterfaces.ListenUpdatedResponse{&h.listenUpdated}
}

func main() {
	ctx := context.Background()
	options := agent.NewSettingsConfigurationOptions()
	options.Agent.Listen.Provider = map[string]interface{}{
		"type":    "deepgram",
		"version": "v2",
		"model":   "flux-general-en",
	}

	h := newHandler()
	conn, err := agent.NewWSUsingChanWithDefaults(ctx, options, h)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Stop()

	if !conn.Connect() {
		log.Fatal("connect failed")
	}

	select {
	case <-h.welcome:
	case <-time.After(10 * time.Second):
		log.Fatal("timed out waiting for Welcome")
	}

	select {
	case <-h.settings:
	case <-time.After(10 * time.Second):
		log.Fatal("timed out waiting for SettingsApplied")
	}

	if err := conn.UpdateListen(&interfacesv1.Listen{Provider: map[string]interface{}{
		"type":          "deepgram",
		"version":       "v2",
		"model":         "flux-general-en",
		"eot_threshold": 0.8,
	}}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-h.listenUpdated:
		fmt.Println("Listen update applied")
	case <-time.After(10 * time.Second):
		log.Fatal("timed out waiting for ListenUpdated")
	}
}
