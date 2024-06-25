// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the speak client implementation for the Deepgram API
*/
package speak

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	speakv1rest "github.com/deepgram/deepgram-go-sdk/pkg/client/speak/v1/rest"
	speakv1ws "github.com/deepgram/deepgram-go-sdk/pkg/client/speak/v1/websocket"
)

/***********************************/
// REST Client
/***********************************/
const (
	RESTPackageVersion = speakv1rest.PackageVersion
)

// Legacy Client Name
//
// Deprecated: This struct is deprecated. Please use RestClient struct. This will be removed in a future release.
type Client = speakv1rest.Client

// New Client Name
type RestClient = speakv1rest.Client

/*
NewWithDefaults creates a new speak client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY

Deprecated: This function is deprecated. Please use NewREST(). This will be removed in a future release.
*/
func NewWithDefaults() *speakv1rest.Client {
	return speakv1rest.NewWithDefaults()
}

/*
New creates a new speak client with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.

Deprecated: This function is deprecated. Please use NewREST(). This will be removed in a future release.
*/
func New(apiKey string, options *interfaces.ClientOptions) *speakv1rest.Client {
	return speakv1rest.New(apiKey, options)
}

/*
NewRESTWithDefaults creates a new speak client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewRESTWithDefaults() *speakv1rest.Client {
	return speakv1rest.NewWithDefaults()
}

/*
New creates a new speak client with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func NewREST(apiKey string, options *interfaces.ClientOptions) *speakv1rest.Client {
	return speakv1rest.New(apiKey, options)
}

/***********************************/
// WebSocket Client
/***********************************/
const (
	WebSocketPackageVersion = speakv1ws.PackageVersion
)

type WebSocketClient = speakv1ws.Client

/*
NewWebSocketForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWebSocketForDemo(ctx context.Context, options *interfaces.SpeakOptions) (*speakv1ws.Client, error) {
	return speakv1ws.NewWebSocketForDemo(ctx, options)
}

/*
NewStreamWithDefaults creates a new websocket connection with all default options

Notes:
  - The callback handler is set to the default handler
*/
func NewWebSocketWithDefaults(ctx context.Context, options *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.Client, error) {
	return speakv1ws.NewWebSocketWithDefaults(ctx, options, callback)
}

/*
NewStream creates a new websocket connection with the specified options

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
func NewWebSocket(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.Client, error) {
	return speakv1ws.NewWebSocket(ctx, apiKey, cOptions, sOptions, callback)
}

/*
NewWebSocketWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

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
func NewWebSocketWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.Client, error) {
	return speakv1ws.NewWebSocketWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}
