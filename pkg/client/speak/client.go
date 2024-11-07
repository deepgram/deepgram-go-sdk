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

// New Client Name
type RESTClient = speakv1rest.RESTClient

/*
NewRESTWithDefaults creates a new speak client with all default options

Returns:
- *Client: a new speak client

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewRESTWithDefaults() *speakv1rest.RESTClient {
	return speakv1rest.NewWithDefaults()
}

/*
New creates a new speak client with the specified options

Input parameters:
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.

Returns:
- *Client: a new speak client
*/
func NewREST(apiKey string, options *interfaces.ClientOptions) *speakv1rest.RESTClient {
	return speakv1rest.New(apiKey, options)
}

/***********************************/
// WebSocket Client
/***********************************/
const (
	WebSocketPackageVersion = speakv1ws.PackageVersion
)

// WSCallback is an alias for listenv1ws.WSCallback
type WSCallback = speakv1ws.WSCallback

// WSChannel is an alias for listenv1ws.WSChannel
type WSChannel = speakv1ws.WSChannel

/*
	Using Callbacks
*/
/*
NewWSUsingCallbackForDemo creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- options: SpeakOptions which allows overriding things like model, etc.

Returns:
- *Client: a new websocket client

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWSUsingCallbackForDemo(ctx context.Context, options *interfaces.WSSpeakOptions) (*speakv1ws.WSCallback, error) {
	return speakv1ws.NewUsingCallbackForDemo(ctx, options)
}

/*
NewWSUsingCallbackWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- options: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Returns:
- *Client: a new websocket client

Notes:
  - The callback handler is set to the default handler
*/
func NewWSUsingCallbackWithDefaults(ctx context.Context, options *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.WSCallback, error) {
	return speakv1ws.NewUsingCallbackWithDefaults(ctx, options, callback)
}

/*
NewWSUsingCallbacks creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Returns:
- *Client: a new websocket client

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewWSUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.WSCallback, error) {
	return speakv1ws.NewUsingCallback(ctx, apiKey, cOptions, sOptions, callback)
}

/*
NewWSUsingCallbackWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Returns:
- *Client: a new websocket client

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewWSUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.WSCallback, error) {
	return speakv1ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}

/*
	Using Channels
*/
/*
NewWSUsingChanForDemo creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- options: SpeakOptions which allows overriding things like model, etc.

Returns:
- *Client: a new websocket client

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWSUsingChanForDemo(ctx context.Context, options *interfaces.WSSpeakOptions) (*speakv1ws.WSChannel, error) {
	return speakv1ws.NewUsingChanForDemo(ctx, options)
}

/*
NewWSUsingChanWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- options: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Returns:
- *Client: a new websocket client

Notes:
  - The callback handler is set to the default handler
*/
func NewWSUsingChanWithDefaults(ctx context.Context, options *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageChan) (*speakv1ws.WSChannel, error) {
	return speakv1ws.NewUsingChanWithDefaults(ctx, options, callback)
}

/*
NewWSUsingChan creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Returns:
- *Client: a new websocket client

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageChan) (*speakv1ws.WSChannel, error) {
	return speakv1ws.NewUsingChan(ctx, apiKey, cOptions, sOptions, callback)
}

/*
NewWSUsingChanWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Returns:
- *Client: a new websocket client

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageChan) (*speakv1ws.WSChannel, error) {
	return speakv1ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}
