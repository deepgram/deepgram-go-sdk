// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
)

/*
Using Channels
*/
// DefaultCallbackHandler is a default callback handler for live transcription
// Simply prints the transcript to stdout
type DefaultChanHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool

	binaryChan    chan *[]byte
	openChan      chan *interfaces.OpenResponse
	metadataChan  chan *interfaces.MetadataResponse
	flushedChan   chan *interfaces.FlushedResponse
	clearedChan   chan *interfaces.ClearedResponse
	closeChan     chan *interfaces.CloseResponse
	warningChan   chan *interfaces.WarningResponse
	errorChan     chan *interfaces.ErrorResponse
	unhandledChan chan *[]byte
}

// ChanRouter routes events
type ChanRouter struct {
	debugWebsocket bool

	// call out to channels
	binaryChan    []*chan *[]byte
	openChan      []*chan *interfaces.OpenResponse
	metadataChan  []*chan *interfaces.MetadataResponse
	flushedChan   []*chan *interfaces.FlushedResponse
	clearedChan   []*chan *interfaces.ClearedResponse
	closeChan     []*chan *interfaces.CloseResponse
	warningChan   []*chan *interfaces.WarningResponse
	errorChan     []*chan *interfaces.ErrorResponse
	unhandledChan []*chan *[]byte
}

/*
Using Callbacks
*/
// DefaultCallbackHandler is a default callback handler for live transcription
// Simply prints the transcript to stdout
type DefaultCallbackHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool
}

// CallbackRouter routes events
type CallbackRouter struct {
	debugWebsocket bool
	callback       interfaces.SpeakMessageCallback
}

// MessageRouter is the interface for routing messages
// Deprecated: Use CallbackRouter instead
type MessageRouter = CallbackRouter
