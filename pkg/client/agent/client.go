// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the agent client implementation for the Deepgram API
*/
package agent

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
	listenv1ws "github.com/deepgram/deepgram-go-sdk/pkg/client/agent/v1/websocket"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/***********************************/
// WebSocket Agent
/***********************************/
const (
	WebSocketPackageVersion = listenv1ws.PackageVersion
)

// WSChannel is an alias for listenv1ws.WSChannel
type WSChannel = listenv1ws.WSChannel

// options
func NewSettingsConfigurationOptions() *interfaces.SettingsConfigurationOptions {
	return interfaces.NewSettingsConfigurationOptions()
}

/*
	Using Channels
*/
/*
NewWSUsingChanForDemo creates a new websocket connection for demo purposes only

Input parameters:
- ctx: context.Context object
- tOptions: SettingsConfigurationOptions which allows overriding things like language, model, etc.

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWSUsingChanForDemo(ctx context.Context, options *interfaces.SettingsConfigurationOptions) (*listenv1ws.WSChannel, error) {
	return listenv1ws.NewUsingChanForDemo(ctx, options)
}

/*
NewWebSocketUsingChanWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- tOptions: SettingsConfigurationOptions which allows overriding things like language, model, etc.

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The chans handler is set to the default handler which just prints all messages to the console
*/
func NewWSUsingChanWithDefaults(ctx context.Context, options *interfaces.SettingsConfigurationOptions, chans msginterfaces.AgentMessageChan) (*listenv1ws.WSChannel, error) {
	return listenv1ws.NewUsingChanWithDefaults(ctx, options, chans)
}

/*
NewWSUsingChan creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: SettingsConfigurationOptions which allows overriding things like language, model, etc.
- chans: AgentMessageChan which is a chans that allows you to perform actions based on the transcription
*/
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.SettingsConfigurationOptions, chans msginterfaces.AgentMessageChan) (*listenv1ws.WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv1ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

/*
NewWSUsingChanWithCancel creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: SettingsConfigurationOptions which allows overriding things like language, model, etc.
- chans: AgentMessageChan which is a chans that allows you to perform actions based on the transcription
*/
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.SettingsConfigurationOptions, chans msginterfaces.AgentMessageChan) (*listenv1ws.WSChannel, error) {
	return listenv1ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}
