// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// *********** WARNING ***********
// This package provides the live client implementation for the Deepgram API
//
// Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
//
// This package is frozen and no new functionality will be added.
// *********** WARNING ***********
package live

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces" //lint:ignore
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	listenv1ws "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/websocket"
)

const (
	PackageVersion = listenv1ws.PackageVersion
)

/***********************************/
// LiveClient
/***********************************/
/*
Client is an alias for the listenv1ws.Client

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
type Client = listenv1ws.Client

/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
func NewForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*Client, error) {
	return listenv1ws.New(ctx, "", &interfaces.ClientOptions{}, options, nil)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Input parameters:
- ctx: context.Context object
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
func NewWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	return listenv1ws.New(ctx, "", &interfaces.ClientOptions{}, tOptions, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
func New(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return listenv1ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

/*
NewWithCancel creates a new websocket connection but has facilities to BYOC (Bring Your Own Cancel)

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

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
func NewWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	return listenv1ws.NewUsingCallbackWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}
