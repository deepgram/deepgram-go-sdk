// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (
	"context"
	"sync"

	"github.com/dvonthenen/websocket"

	live "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1"
	msginterface "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
)

// ClientOptions defines any options for the client
type ClientOptions struct {
	Host            string
	ApiVersion      string
	Path            string
	RedirectService bool
	SkipServerAuth  bool
}

// Client return websocket client connection
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
