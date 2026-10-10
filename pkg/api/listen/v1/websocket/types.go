// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"sync"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
)

type extraCache struct {
	mu    sync.RWMutex
	extra map[string]string
}

func (c *extraCache) store(extra map[string]string) {
	if len(extra) == 0 {
		return
	}

	copied := make(map[string]string, len(extra))
	for k, v := range extra {
		copied[k] = v
	}

	c.mu.Lock()
	c.extra = copied
	c.mu.Unlock()
}

func (c *extraCache) load() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.extra) == 0 {
		return nil
	}

	copied := make(map[string]string, len(c.extra))
	for k, v := range c.extra {
		copied[k] = v
	}
	return copied
}

/*
Using Channels
*/
// DefaultCallbackHandler is a default callback handler for live transcription
// Simply prints the transcript to stdout
type DefaultChanHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool

	openChan          chan *interfaces.OpenResponse
	messageChan       chan *interfaces.MessageResponse
	metadataChan      chan *interfaces.MetadataResponse
	speechStartedChan chan *interfaces.SpeechStartedResponse
	utteranceEndChan  chan *interfaces.UtteranceEndResponse
	closeChan         chan *interfaces.CloseResponse
	errorChan         chan *interfaces.ErrorResponse
	unhandledChan     chan *[]byte
}

// ChanRouter routes events
type ChanRouter struct {
	debugWebsocket bool

	extraMetadata extraCache

	// call out to channels
	openChan          []*chan *interfaces.OpenResponse
	messageChan       []*chan *interfaces.MessageResponse
	metadataChan      []*chan *interfaces.MetadataResponse
	speechStartedChan []*chan *interfaces.SpeechStartedResponse
	utteranceEndChan  []*chan *interfaces.UtteranceEndResponse
	closeChan         []*chan *interfaces.CloseResponse
	errorChan         []*chan *interfaces.ErrorResponse
	unhandledChan     []*chan *[]byte
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
	callback       interfaces.LiveMessageCallback

	extraMetadata extraCache
}

// MessageRouter is the interface for routing messages
// Deprecated: Use CallbackRouter instead
type MessageRouter = CallbackRouter
