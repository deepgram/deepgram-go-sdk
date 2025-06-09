// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"strings"

	klog "k8s.io/klog/v2"

	websocketv1api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewUsingCallbackForDemo(ctx context.Context, options *clientinterfaces.LiveTranscriptionOptions) (*WSCallback, error) {
	return New(ctx, "", &clientinterfaces.ClientOptions{}, options, nil)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewUsingCallbackWithDefaults(ctx context.Context, options *clientinterfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*WSCallback, error) {
	return NewUsingCallback(ctx, "", &clientinterfaces.ClientOptions{}, options, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription
*/
func NewUsingCallback(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptions, tOptions *clientinterfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*WSCallback, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription
*/
func NewUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptions, tOptions *clientinterfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*WSCallback, error) {
	klog.V(6).Infof("New() ENTER\n")

	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	if len(tOptions.Keyterm) > 0 && !strings.HasPrefix(tOptions.Model, "nova-3") {
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

	if callback == nil {
		klog.V(2).Infof("Using DefaultCallbackHandler.\n")
		callback = websocketv1api.NewDefaultCallbackHandler()
	}

	// init
	var router commoninterfaces.Router
	router = websocketv1api.NewCallbackRouter(callback)

	conn := WSCallback{
		cOptions:  cOptions,
		tOptions:  tOptions,
		callback:  callback,
		router:    &router,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}

	var handler commoninterfaces.WebSocketHandler
	handler = &conn
	conn.WSClient = common.NewWS(ctx, ctxCancel, apiKey, cOptions, &handler, &router)

	klog.V(3).Infof("NewDeepGramWSClient Succeeded\n")
	klog.V(6).Infof("New() LEAVE\n")

	return &conn, nil
}

/***********************************/
// Deprecated functions
/***********************************/
/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY

Deprecated: Use NewUsingCallbackForDemo instead
*/
func NewForDemo(ctx context.Context, options *clientinterfaces.LiveTranscriptionOptions) (*WSCallback, error) {
	return NewUsingCallbackForDemo(ctx, options)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console

Deprecated: Use NewUsingCallbackWithDefaults instead
*/
func NewWithDefaults(ctx context.Context, options *clientinterfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*WSCallback, error) {
	return NewUsingCallback(ctx, "", &clientinterfaces.ClientOptions{}, options, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Deprecated: Use NewUsingCallback instead
*/
func New(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptions, tOptions *clientinterfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*WSCallback, error) {
	return NewUsingCallback(ctx, apiKey, cOptions, tOptions, callback)
}
