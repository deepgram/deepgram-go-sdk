// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package speak

import (
	"context"
	"sync"
	"time"

	"github.com/dvonthenen/websocket"

	speak "github.com/deepgram/deepgram-go-sdk/pkg/api/speak-stream/v1"
	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/speak-stream/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// Client is a struct representing the websocket client connection
type StreamClient struct {
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

	// internal constants for retry, waits, back-off, etc.
	lastDatagram *time.Time
	muFinal      sync.RWMutex
}
