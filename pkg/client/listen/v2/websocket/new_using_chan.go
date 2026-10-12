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

// NewUsingChanForDemo creates a Flux WebSocket client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewUsingChanForDemo(ctx context.Context, options *clientinterfaces.FluxTranscriptionOptions) (*WSChannel, error) {
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptionsV2{}, options, nil)
}

// NewUsingChanWithDefaults creates a Flux WebSocket client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewUsingChanWithDefaults(ctx context.Context, options *clientinterfaces.FluxTranscriptionOptions, chans msginterfaces.FluxMessageChan) (*WSChannel, error) { //nolint:gocritic
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptionsV2{}, options, chans)
}

// NewUsingChan creates a Flux WebSocket client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
// If chans is nil, the default stdout-printing channel handler is used.
func NewUsingChan(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptionsV2, tOptions *clientinterfaces.FluxTranscriptionOptions, chans msginterfaces.FluxMessageChan) (*WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

// NewUsingChanWithCancel creates a Flux WebSocket client and lets the caller supply
// their own context cancel function (BYOC — Bring Your Own Cancel).
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
// If chans is nil, the default stdout-printing channel handler is used.
func NewUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptionsV2, tOptions *clientinterfaces.FluxTranscriptionOptions, chans msginterfaces.FluxMessageChan) (*WSChannel, error) { //nolint:gocritic
	klog.V(6).Infof("flux.NewUsingChanWithCancel() ENTER\n")

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

	if chans == nil {
		klog.V(2).Infof("Using DefaultChanHandler.\n")
		chans = websocketv2api.NewDefaultChanHandler()
	}

	var router commoninterfaces.Router
	router = websocketv2api.NewChanRouter(chans)

	conn := WSChannel{
		cOptions:  cOptions,
		tOptions:  tOptions,
		chans:     make([]*msginterfaces.FluxMessageChan, 0),
		router:    &router,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}

	var handler commoninterfaces.WebSocketHandler
	handler = &conn
	conn.WSClient = common.NewWS(ctx, ctxCancel, apiKey, cOptions, &handler, &router)

	klog.V(3).Infof("flux WSChannel created\n")
	klog.V(6).Infof("flux.NewUsingChanWithCancel() LEAVE\n")
	return &conn, nil
}
