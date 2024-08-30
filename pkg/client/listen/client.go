// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the prerecorded client implementation for the Deepgram API
*/
package listen

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	listenv1rest "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"
	listenv1ws "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/websocket"
)

/***********************************/
// REST Client
/***********************************/
const (
	RESTPackageVersion = listenv1rest.PackageVersion
)

// RestClient is an alias for listenv1rest.Client
type RESTClient = listenv1rest.RESTClient

/*
NewRESTWithDefaults creates a new analyze/read client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewRESTWithDefaults() *listenv1rest.RESTClient {
	return listenv1rest.NewWithDefaults()
}

/*
NewREST creates a new prerecorded client with the specified options

Input parameters:
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func NewREST(apiKey string, options *interfaces.ClientOptions) *listenv1rest.RESTClient {
	return listenv1rest.New(apiKey, options)
}

/***********************************/
// WebSocket / Streaming / Live
/***********************************/
const (
	WebSocketPackageVersion = listenv1ws.PackageVersion
)

// WSCallback is an alias for listenv1ws.WSCallback
type WSCallback = listenv1ws.WSCallback

// WSChannel is an alias for listenv1ws.WSChannel
type WSChannel = listenv1ws.WSChannel

/*
	Using Callbacks
*/
/*
NewWSUsingCallbackForDemo creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWSUsingCallbackForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSCallback, error) {
	return listenv1ws.NewUsingCallbackForDemo(ctx, options)
}

// NewWebSocketUsingCallbackForDemo is an alias for NewWSUsingCallbackForDemo
// TODO: Deprecate this function later
func NewWebSocketUsingCallbackForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSCallback, error) {
	return NewWSUsingCallbackForDemo(ctx, options)
}

/*
NewWSUsingCallbackWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewWSUsingCallbackWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return listenv1ws.NewUsingCallbackWithDefaults(ctx, tOptions, callback)
}

// NewWebSocketUsingCallbackWithDefaults is an alias for NewWSUsingCallbackWithDefaults
// TODO: Deprecate this function later
func NewWebSocketUsingCallbackWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return NewWSUsingCallbackWithDefaults(ctx, tOptions, callback)
}

/*
NewWSUsingCallback creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewWSUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv1ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

// NewWebSocketUsingCallback is an alias for NewWSUsingCallback
// TODO: Deprecate this function later
func NewWebSocketUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return NewWSUsingCallback(ctx, apiKey, cOptions, tOptions, callback)
}

/*
NewWSUsingCallbackWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewWSUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return listenv1ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

// NewWebSocketUsingCallbackWithCancel is an alias for NewWSUsingCallbackWithCancel
// TODO: Deprecate this function later
func NewWebSocketUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return NewWSUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

/*
	Using Channels
*/
/*
NewWSUsingChanForDemo creates a new websocket connection for demo purposes only

Input parameters:
- ctx: context.Context object
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWSUsingChanForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSChannel, error) {
	return listenv1ws.NewUsingChanForDemo(ctx, options)
}

// NewWebSocketUsingChanForDemo is an alias for NewWSUsingChanForDemo
// TODO: Deprecate this function later
func NewWebSocketUsingChanForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSChannel, error) {
	return NewWSUsingChanForDemo(ctx, options)
}

/*
NewWebSocketUsingChanWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The chans handler is set to the default handler which just prints all messages to the console
*/
func NewWSUsingChanWithDefaults(ctx context.Context, options *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error) {
	return listenv1ws.NewUsingChanWithDefaults(ctx, options, chans)
}

// NewWebSocketUsingChanWithDefaults is an alias for NewWSUsingChanWithDefaults
// TODO: Deprecate this function later
func NewWebSocketUsingChanWithDefaults(ctx context.Context, options *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error) {
	return NewWSUsingChanWithDefaults(ctx, options, chans)
}

/*
NewWSUsingChan creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- chans: LiveMessageCallback which is a chans that allows you to perform actions based on the transcription
*/
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv1ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

// NewWebSocketUsingChan is an alias for NewWSUsingChan
// TODO: Deprecate this function later
func NewWebSocketUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error) {
	return NewWSUsingChan(ctx, apiKey, cOptions, tOptions, chans)
}

/*
NewWSUsingChanWithCancel creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- chans: LiveMessageCallback which is a chans that allows you to perform actions based on the transcription
*/
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error) {
	return listenv1ws.NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

// NewWebSocketUsingChanWithCancel is an alias for NewWSUsingChanWithCancel
// TODO: Deprecate this function later
func NewWebSocketUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error) {
	return NewWSUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, chans)
}

/***********************************/
// Deprecated (THESE WILL STILL WORK,
// BUT WILL BE REMOVED IN A FUTURE RELEASE)
/***********************************/
/*
NewWebSocketForDemo creates a new websocket connection with all default options

Please see NewWebSocketUsingCallbackForDemo for more information.

TODO: Deprecate this function later
*/
func NewWebSocketForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.Client, error) {
	return NewWSUsingCallbackForDemo(ctx, options)
}

/*
NewWebSocketWithDefaults creates a new websocket connection with all default options

Please see NewWebSocketUsingCallbackWithDefaults for more information.

TODO: Deprecate this function later
*/
func NewWebSocketWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.Client, error) {
	return NewWSUsingCallbackWithDefaults(ctx, tOptions, callback)
}

/*
NewWebSocket creates a new websocket connection with the specified options

Please see NewWebSocketUsingCallback for more information.

TODO: Deprecate this function later
*/
func NewWebSocket(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.Client, error) {
	return NewWSUsingCallback(ctx, apiKey, cOptions, tOptions, callback)
}

/*
NewWebSocketWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

Please see NewWebSocketUsingCallbackWithCancel for more information.

TODO: Deprecate this function later
*/
func NewWebSocketWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.Client, error) {
	return NewWSUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

/***********************************/
// REST Client
/***********************************/
// PreRecordedClient is an alias for listenv1rest.Client
//
// Deprecated: This package is deprecated. Use RestClient instead. This will be removed in a future release.
type PreRecordedClient = listenv1rest.RESTClient

// NewPreRecordedWithDefaults is an alias for NewRESTWithDefaults
//
// Deprecated: This package is deprecated. Use NewRESTWithDefaults instead. This will be removed in a future release.
func NewPreRecordedWithDefaults() *listenv1rest.RESTClient {
	return NewRESTWithDefaults()
}

// NewPreRecorded is an alias for NewREST
//
// Deprecated: This package is deprecated. Use NewREST instead. This will be removed in a future release.
func NewPreRecorded(apiKey string, options *interfaces.ClientOptions) *listenv1rest.RESTClient {
	return NewREST(apiKey, options)
}

/***********************************/
// WebSocket / Streaming / Live
/***********************************/
// LiveClient is an alias for listenv1rest.Client
//
// Deprecated: This alias is deprecated. Use WSCallback instead. This will be removed in a future release.
type LiveClient = listenv1ws.Client

/*
	Older "Live" functions
*/
// NewLiveForDemo is an alias for NewWebSocketForDemo
//
// Deprecated: This package is deprecated. Use NewWebSocketForDemo instead. This will be removed in a future release.
func NewLiveForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSCallback, error) {
	return NewWebSocketForDemo(ctx, options)
}

// NewLiveWithDefaults is an alias for NewWebSocketWithDefaults
//
// Deprecated: This package is deprecated. Use NewWebSocketWithDefaults instead. This will be removed in a future release.
func NewLiveWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return NewWebSocketWithDefaults(ctx, tOptions, callback)
}

// NewLive is an alias for NewWebSocket
//
// Deprecated: This package is deprecated. Use NewWebSocket instead. This will be removed in a future release.
func NewLive(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return NewWebSocket(ctx, apiKey, cOptions, tOptions, callback)
}

// NewLiveWithCancel is an alias for NewWebSocketWithCancel
//
// Deprecated: This package is deprecated. Use NewWebSocketWithCancel instead. This will be removed in a future release.
func NewLiveWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error) {
	return NewWebSocketWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}
