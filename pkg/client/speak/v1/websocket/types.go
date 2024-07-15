// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"sync"
	"time"

	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

// external structs
type TextSource struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// internal structs
type controlMessage struct {
	Type string `json:"type"`
}

// Client is an alias for WSCallback
// Deprecated: use WSCallback instead
type Client = WSCallback

// WSCallback is a struct representing the websocket client connection using callbacks
type WSCallback struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptions
	sOptions *interfaces.WSSpeakOptions

	callback msginterface.SpeakMessageCallback
	router   *commoninterfaces.Router

	// internal constants for retry, waits, back-off, etc.
	lastDatagram *time.Time
	muFinal      sync.RWMutex
	flushCount   int64
}

// WSChannel is a struct representing the websocket client connection using channels
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptions
	sOptions *interfaces.WSSpeakOptions

	chans  []*msginterface.SpeakMessageChan
	router *commoninterfaces.Router

	// internal constants for retry, waits, back-off, etc.
	lastDatagram *time.Time
	muFinal      sync.RWMutex
	flushCount   int64
}
