// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (
	"context"
	"sync"

	"github.com/dvonthenen/websocket"

	live "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1"
	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// Client is a struct representing the websocket client connection
type Client struct {
	cOptions *interfaces.ClientOptions
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
