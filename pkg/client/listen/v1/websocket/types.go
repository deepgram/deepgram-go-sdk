// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"sync"
	"time"

	"github.com/dvonthenen/websocket"

	live "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket"
	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// Client is a struct representing the websocket client connection
type Client struct {
	cOptions *interfaces.ClientOptions
	tOptions *interfaces.LiveTranscriptionOptions

	sendBuf   chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc

	muConn   sync.RWMutex
	wsconn   *websocket.Conn
	retry    bool
	retryCnt int64

	callback msginterface.LiveMessageCallback
	router   *live.MessageRouter

	// internal constants for retry, waits, back-off, etc.
	lastDatagram *time.Time
	muFinal      sync.RWMutex
}
