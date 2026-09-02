// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package speak provides the entry points for the Deepgram Flux TTS (v2/speak) APIs:
// the batch REST transport (POST /v2/speak) and the streaming WebSocket transport
// (wss://api.deepgram.com/v2/speak).
package speak

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
	speakv2rest "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/rest"
	speakv2ws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/websocket"
)

const (
	RESTPackageVersion      = speakv2rest.PackageVersion
	WebSocketPackageVersion = speakv2ws.PackageVersion
)

// RESTClient is an alias for the Flux TTS batch (REST) transport client.
type RESTClient = speakv2rest.Client

// WSCallback is an alias for the callback-based Flux TTS WebSocket client.
type WSCallback = speakv2ws.WSCallback

// WSChannel is an alias for the channel-based Flux TTS WebSocket client.
type WSChannel = speakv2ws.WSChannel

// -----------------------------------------------------------------------
// REST factory functions
// -----------------------------------------------------------------------

// NewRESTWithDefaults creates a Flux TTS batch client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
// An error (never a nil client) is returned when configuration fails, e.g. when
// no credentials are available.
func NewRESTWithDefaults() (*RESTClient, error) {
	return speakv2rest.NewWithDefaults()
}

// NewREST creates a Flux TTS batch client with the specified options.
// An error (never a nil client) is returned when configuration fails.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
func NewREST(apiKey string, options *interfacesv1.ClientOptions) (*RESTClient, error) {
	return speakv2rest.New(apiKey, options)
}

// -----------------------------------------------------------------------
// Callback-based WebSocket factory functions
// -----------------------------------------------------------------------

// NewWSUsingCallbackForDemo creates a Flux TTS client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingCallbackForDemo(ctx context.Context, sOptions *interfaces.SpeakV2WSOptions) (*WSCallback, error) {
	return speakv2ws.NewUsingCallbackForDemo(ctx, sOptions)
}

// NewWSUsingCallbackWithDefaults creates a Flux TTS client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingCallbackWithDefaults(ctx context.Context, sOptions *interfaces.SpeakV2WSOptions, callback msginterfaces.FluxSpeakMessageCallback) (*WSCallback, error) {
	return speakv2ws.NewUsingCallbackWithDefaults(ctx, sOptions, callback)
}

// NewWSUsingCallback creates a Flux TTS client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptionsV2, sOptions *interfaces.SpeakV2WSOptions, callback msginterfaces.FluxSpeakMessageCallback) (*WSCallback, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return speakv2ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}

// NewWSUsingCallbackWithCancel creates a Flux TTS client and lets the caller supply their own
// context cancel function (BYOC — Bring Your Own Cancel).
func NewWSUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptionsV2, sOptions *interfaces.SpeakV2WSOptions, callback msginterfaces.FluxSpeakMessageCallback) (*WSCallback, error) {
	return speakv2ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}

// -----------------------------------------------------------------------
// Channel-based WebSocket factory functions
// -----------------------------------------------------------------------

// NewWSUsingChanForDemo creates a Flux TTS client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingChanForDemo(ctx context.Context, sOptions *interfaces.SpeakV2WSOptions) (*WSChannel, error) {
	return speakv2ws.NewUsingChanForDemo(ctx, sOptions)
}

// NewWSUsingChanWithDefaults creates a Flux TTS client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingChanWithDefaults(ctx context.Context, sOptions *interfaces.SpeakV2WSOptions, chans msginterfaces.FluxSpeakMessageChan) (*WSChannel, error) {
	return speakv2ws.NewUsingChanWithDefaults(ctx, sOptions, chans)
}

// NewWSUsingChan creates a Flux TTS client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptionsV2, sOptions *interfaces.SpeakV2WSOptions, chans msginterfaces.FluxSpeakMessageChan) (*WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return speakv2ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, chans)
}

// NewWSUsingChanWithCancel creates a Flux TTS client and lets the caller supply their own
// context cancel function (BYOC — Bring Your Own Cancel).
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptionsV2, sOptions *interfaces.SpeakV2WSOptions, chans msginterfaces.FluxSpeakMessageChan) (*WSChannel, error) {
	return speakv2ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, chans)
}
