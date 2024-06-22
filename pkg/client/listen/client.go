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

type RestClient = listenv1rest.Client

/*
NewRESTWithDefaults creates a new analyze/read client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewRESTWithDefaults() *listenv1rest.Client {
	return listenv1rest.NewWithDefaults()
}

/*
NewREST creates a new prerecorded client with the specified options

Input parameters:
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func NewREST(apiKey string, options *interfaces.ClientOptions) *listenv1rest.Client {
	return listenv1rest.New(apiKey, options)
}

/***********************************/
// LiveClient
/***********************************/
const (
	WebSocketPackageVersion = listenv1ws.PackageVersion
)

type WebSocketClient = listenv1ws.Client

/*
NewWebSocketForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWebSocketForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.Client, error) {
	return listenv1ws.New(ctx, "", &interfaces.ClientOptions{}, options, nil)
}

/*
NewWebSocketWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewWebSocketWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.Client, error) {
	return listenv1ws.New(ctx, "", &interfaces.ClientOptions{}, tOptions, callback)
}

/*
NewWebSocket creates a new websocket connection with the specified options

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
func NewWebSocket(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.Client, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv1ws.NewWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

/*
NewWebSocketWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

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
func NewWebSocketWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.Client, error) {
	return listenv1ws.NewWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}
