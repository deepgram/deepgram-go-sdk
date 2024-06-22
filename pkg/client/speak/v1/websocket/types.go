// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"sync"

	"github.com/dvonthenen/websocket"

	speak "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket"
	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

// Client is a struct representing the websocket client connection
type Client struct {
	cOptions *interfaces.ClientOptions
	sOptions *interfaces.SpeakOptions

	sendBuf   chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc

	muConn   sync.RWMutex
	wsconn   *websocket.Conn
	retry    bool
	retryCnt int64

	callback msginterface.SpeakMessageCallback
	router   *speak.MessageRouter
}
