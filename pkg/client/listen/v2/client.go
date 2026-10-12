// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package listen provides the entry points for the Deepgram Flux (v2/listen) WebSocket API.
package listen

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listenv2ws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2/websocket"
)

const (
	WebSocketPackageVersion = listenv2ws.PackageVersion
)

// WSCallback is an alias for the callback-based Flux WebSocket client.
type WSCallback = listenv2ws.WSCallback

// WSChannel is an alias for the channel-based Flux WebSocket client.
type WSChannel = listenv2ws.WSChannel

// -----------------------------------------------------------------------
// Callback-based factory functions
// -----------------------------------------------------------------------

// NewWSUsingCallbackForDemo creates a Flux client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingCallbackForDemo(ctx context.Context, tOptions *interfaces.FluxTranscriptionOptions) (*WSCallback, error) {
	return listenv2ws.NewUsingCallbackForDemo(ctx, tOptions)
}

// NewWSUsingCallbackWithDefaults creates a Flux client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingCallbackWithDefaults(ctx context.Context, tOptions *interfaces.FluxTranscriptionOptions, callback msginterfaces.FluxMessageCallback) (*WSCallback, error) {
	return listenv2ws.NewUsingCallbackWithDefaults(ctx, tOptions, callback)
}

// NewWSUsingCallback creates a Flux client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptionsV2, tOptions *interfaces.FluxTranscriptionOptions, callback msginterfaces.FluxMessageCallback) (*WSCallback, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv2ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

// NewWSUsingCallbackWithCancel creates a Flux client and lets the caller supply their own
// context cancel function (BYOC — Bring Your Own Cancel).
func NewWSUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptionsV2, tOptions *interfaces.FluxTranscriptionOptions, callback msginterfaces.FluxMessageCallback) (*WSCallback, error) {
	return listenv2ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

// -----------------------------------------------------------------------
// Channel-based factory functions
// -----------------------------------------------------------------------

// NewWSUsingChanForDemo creates a Flux client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingChanForDemo(ctx context.Context, tOptions *interfaces.FluxTranscriptionOptions) (*WSChannel, error) {
	return listenv2ws.NewUsingChanForDemo(ctx, tOptions)
}

// NewWSUsingChanWithDefaults creates a Flux client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingChanWithDefaults(ctx context.Context, tOptions *interfaces.FluxTranscriptionOptions, chans msginterfaces.FluxMessageChan) (*WSChannel, error) {
	return listenv2ws.NewUsingChanWithDefaults(ctx, tOptions, chans)
}

// NewWSUsingChan creates a Flux client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptionsV2, tOptions *interfaces.FluxTranscriptionOptions, chans msginterfaces.FluxMessageChan) (*WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv2ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

// NewWSUsingChanWithCancel creates a Flux client and lets the caller supply their own
// context cancel function (BYOC — Bring Your Own Cancel).
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptionsV2, tOptions *interfaces.FluxTranscriptionOptions, chans msginterfaces.FluxMessageChan) (*WSChannel, error) {
	return listenv2ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}
