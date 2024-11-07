// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package listen

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	listenv1rest "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"
	listenv1ws "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/websocket"
)

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
