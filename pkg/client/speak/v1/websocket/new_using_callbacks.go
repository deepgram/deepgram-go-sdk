// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the speak/streaming client implementation for the Deepgram API
package websocketv1

import (
	"context"

	klog "k8s.io/klog/v2"

	websocketv1api "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewUsingCallbackForDemo(ctx context.Context, options *clientinterfaces.WSSpeakOptions) (*WSCallback, error) {
	return NewUsingCallback(ctx, "", &clientinterfaces.ClientOptions{}, options, nil)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The callback handler is set to the default handler
*/
func NewUsingCallbackWithDefaults(ctx context.Context, options *clientinterfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*WSCallback, error) {
	return NewUsingCallback(ctx, "", &clientinterfaces.ClientOptions{}, options, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewUsingCallback(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptions, sOptions *clientinterfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*WSCallback, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}

/*
NewWithCancel creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptions, sOptions *clientinterfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*WSCallback, error) {
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
		klog.V(1).Infof("SpeakOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	if callback == nil {
		klog.V(2).Infof("Using DefaultCallbackHandler.\n")
		callback = websocketv1api.NewDefaultCallbackHandler()
	}

	// init
	var router commoninterfaces.Router
	router = websocketv1api.NewCallbackRouter(callback)

	// init
	conn := Client{
		cOptions:  cOptions,
		sOptions:  sOptions,
		callback:  callback,
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
