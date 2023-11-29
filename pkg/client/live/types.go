// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the live/streaming client implementation for the Deepgram API
*/
package live

import (
	"context"
	"sync"

	"github.com/dvonthenen/websocket"

	live "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1"
	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// ClientOptions defines any options for the client
type ClientOptions struct {
	Host            string // override for the host endpoint
	ApiVersion      string // override for the version used
	Path            string // override for the endpoint path usually <version/listen>
	RedirectService bool   // allows HTTP redirects to be followed
	SkipServerAuth  bool   // keeps the client from authenticating with the server
	EnableKeepAlive bool   // enables the keep alive feature
}

// Client is a struct representing the websocket client connection
type Client struct {
	cOptions *ClientOptions
	apiKey   string
	tOptions interfaces.LiveTranscriptionOptions

	sendBuf   chan []byte
	org       context.Context
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu     sync.RWMutex
	wsconn *websocket.Conn
	retry  bool

	callback msginterface.LiveMessageCallback
	router   *live.MessageRouter
}
