// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"

	klog "k8s.io/klog/v2"

	websocketv1api "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewUsingChanForDemo(ctx context.Context, options *clientinterfaces.WSSpeakOptions) (*WSChannel, error) {
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptions{}, options, nil)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The chans handler is set to the default handler which just prints all messages to the console
*/
func NewUsingChanWithDefaults(ctx context.Context, options *clientinterfaces.WSSpeakOptions, chans msginterfaces.SpeakMessageChan) (*WSChannel, error) { // gocritic:ignore
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptions{}, options, chans)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- chans: LiveMessageCallback which is a chans that allows you to perform actions based on the transcription
*/
func NewUsingChan(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptions, sOptions *clientinterfaces.WSSpeakOptions, chans msginterfaces.SpeakMessageChan) (*WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, chans)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- chans: LiveMessageCallback which is a chans that allows you to perform actions based on the transcription
*/
func NewUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptions, sOptions *clientinterfaces.WSSpeakOptions, chans msginterfaces.SpeakMessageChan) (*WSChannel, error) {
	klog.V(6).Infof("speak.New() ENTER\n")

	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	err := cOptions.Parse()
	if err != nil {
		klog.V(1).Infof("ClientOptions.Parse() failed. Err: %v\n", err)
		return nil, err
	}
	err = sOptions.Check()
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
		sOptions:  sOptions,
		chans:     make([]*msginterfaces.SpeakMessageChan, 0),
		router:    &router,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}

	var handler commoninterfaces.WebSocketHandler
	handler = &conn
	conn.WSClient = common.NewWS(ctx, ctxCancel, apiKey, cOptions, &handler, &router)

	klog.V(3).Infof("NewDeepGramWSClient Succeeded\n")
	klog.V(6).Infof("speak.New() LEAVE\n")

	return &conn, nil
}
