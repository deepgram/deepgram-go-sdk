// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"strings"

	klog "k8s.io/klog/v2"

	websocketv1api "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewUsingChanForDemo(ctx context.Context, options *clientinterfaces.SettingsConfigurationOptions) (*WSChannel, error) {
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptions{}, options, nil)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The chans handler is set to the default handler which just prints all messages to the console
*/
func NewUsingChanWithDefaults(ctx context.Context, options *clientinterfaces.SettingsConfigurationOptions, chans msginterfaces.AgentMessageChan) (*WSChannel, error) { // gocritic:ignore
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptions{}, options, chans)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: SettingsConfigurationOptions which allows overriding things like language, model, etc.
- chans: AgentMessageChan which is a chans that allows you to perform actions based on the transcription
*/
func NewUsingChan(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptions, tOptions *clientinterfaces.SettingsConfigurationOptions, chans msginterfaces.AgentMessageChan) (*WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: SettingsConfigurationOptions which allows overriding things like language, model, etc.
- chans: AgentMessageChan which is a chans that allows you to perform actions based on the transcription
*/
func NewUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptions, tOptions *clientinterfaces.SettingsConfigurationOptions, chans msginterfaces.AgentMessageChan) (*WSChannel, error) {
	klog.V(6).Infof("agent.New() ENTER\n")

	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	if len(tOptions.Agent.Listen.Keyterm) > 0 && !strings.HasPrefix(tOptions.Agent.Listen.Model, "nova-3") {
		klog.V(1).Info("Keyterms are only supported with nova-3 models.")
		return nil, nil
	}
	err := cOptions.Parse()
	if err != nil {
		klog.V(1).Infof("ClientOptions.Parse() failed. Err: %v\n", err)
		return nil, err
	}
	err = tOptions.Check()
	if err != nil {
		klog.V(1).Infof("TranscribeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	if chans == nil {
		klog.V(2).Infof("Using DefaultCallbackHandler.\n")
		chans = websocketv1api.NewDefaultChanHandler()
	}

	// init
	var router commoninterfaces.Router
	router = websocketv1api.NewChanRouter(chans)

	conn := WSChannel{
		cOptions:  cOptions,
		tOptions:  tOptions,
		chans:     make([]*msginterfaces.AgentMessageChan, 0),
		router:    &router,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}

	var handler commoninterfaces.WebSocketHandler
	handler = &conn
	conn.WSClient = common.NewWS(ctx, ctxCancel, apiKey, cOptions, &handler, &router)

	klog.V(3).Infof("NewDeepGramWSClient Succeeded\n")
	klog.V(6).Infof("agent.New() LEAVE\n")

	return &conn, nil
}
