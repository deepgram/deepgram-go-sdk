// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"

	klog "k8s.io/klog/v2"

	websocketv2api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// NewUsingCallbackForDemo creates a Flux WebSocket client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewUsingCallbackForDemo(ctx context.Context, options *clientinterfaces.FluxTranscriptionOptions) (*WSCallback, error) {
	return NewUsingCallback(ctx, "", &clientinterfaces.ClientOptionsV2{}, options, nil)
}

// NewUsingCallbackWithDefaults creates a Flux WebSocket client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewUsingCallbackWithDefaults(ctx context.Context, options *clientinterfaces.FluxTranscriptionOptions, callback msginterfaces.FluxMessageCallback) (*WSCallback, error) {
	return NewUsingCallback(ctx, "", &clientinterfaces.ClientOptionsV2{}, options, callback)
}

// NewUsingCallback creates a Flux WebSocket client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
// If callback is nil, the default stdout-printing handler is used.
func NewUsingCallback(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptionsV2, tOptions *clientinterfaces.FluxTranscriptionOptions, callback msginterfaces.FluxMessageCallback) (*WSCallback, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

// NewUsingCallbackWithCancel creates a Flux WebSocket client and lets the caller supply
// their own context cancel function (BYOC — Bring Your Own Cancel).
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
// If callback is nil, the default stdout-printing handler is used.
func NewUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptionsV2, tOptions *clientinterfaces.FluxTranscriptionOptions, callback msginterfaces.FluxMessageCallback) (*WSCallback, error) {
	klog.V(6).Infof("flux.NewUsingCallbackWithCancel() ENTER\n")

	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	if err := cOptions.Parse(); err != nil {
		klog.V(1).Infof("ClientOptions.Parse() failed. Err: %v\n", err)
		return nil, err
	}
	if err := tOptions.Check(); err != nil {
		klog.V(1).Infof("FluxTranscriptionOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	if callback == nil {
		klog.V(2).Infof("Using DefaultCallbackHandler.\n")
		callback = websocketv2api.NewDefaultCallbackHandler()
	}

	var router commoninterfaces.Router
	router = websocketv2api.NewCallbackRouter(callback)

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

	klog.V(3).Infof("flux WSCallback created\n")
	klog.V(6).Infof("flux.NewUsingCallbackWithCancel() LEAVE\n")
	return &conn, nil
}
