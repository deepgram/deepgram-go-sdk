// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"sync"
	"sync/atomic"

	msginterface "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// finishState tracks the two milestones a graceful close waits on: the final
// SessionMetadata event and the peer's socket closure. Both clients embed a
// pointer to it, and the event observers mark the milestones as they arrive.
type finishState struct {
	sessionMetaDone chan struct{}
	peerCloseDone   chan struct{}
	sessionMetaOnce sync.Once
	peerCloseOnce   sync.Once
}

func newFinishState() *finishState {
	return &finishState{
		sessionMetaDone: make(chan struct{}),
		peerCloseDone:   make(chan struct{}),
	}
}

func (f *finishState) markSessionMetadata() {
	f.sessionMetaOnce.Do(func() { close(f.sessionMetaDone) })
}

func (f *finishState) markPeerClose() {
	f.peerCloseOnce.Do(func() { close(f.peerCloseDone) })
}

// wsHandlerCore is the subset of the internal WebSocketHandler contract both
// clients implement directly.
type wsHandlerCore interface {
	GetURL(host string) (string, error)
	ProcessMessage(wsType int, byMsg []byte) error
	ProcessError(err error) error
	Start()
	GetCloseMsg() []byte
}

// wsHandlerShim adapts a client to the internal WebSocketHandler contract. It
// exists so the handler's no-op Finish() hook does not occupy the clients' public
// Finish name, which is the user-facing graceful close (Finish(ctx)).
type wsHandlerShim struct {
	core wsHandlerCore
}

func (s *wsHandlerShim) GetURL(host string) (string, error) { return s.core.GetURL(host) }
func (s *wsHandlerShim) ProcessMessage(wsType int, byMsg []byte) error {
	return s.core.ProcessMessage(wsType, byMsg)
}
func (s *wsHandlerShim) ProcessError(err error) error { return s.core.ProcessError(err) }
func (s *wsHandlerShim) Start()                       { s.core.Start() }
func (s *wsHandlerShim) Finish()                      {}
func (s *wsHandlerShim) GetCloseMsg() []byte          { return s.core.GetCloseMsg() }

// WSCallback is a Flux TTS WebSocket client that delivers server events via a
// FluxSpeakMessageCallback.
type WSCallback struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptionsV2
	sOptions *interfaces.SpeakV2WSOptions

	callback msginterface.FluxSpeakMessageCallback
	router   *commoninterfaces.Router

	finishMu   sync.Mutex
	finish     *finishState
	pingActive atomic.Bool
}

// currentFinishState returns the finishState for the session in progress.
func (c *WSCallback) currentFinishState() *finishState {
	c.finishMu.Lock()
	defer c.finishMu.Unlock()
	return c.finish
}

// resetFinishStateIfClosed installs a fresh finishState when the previous
// session already ended (its peer-close milestone fired), so Finish works on a
// reconnected session. A live session's milestones are never discarded.
func (c *WSCallback) resetFinishStateIfClosed() {
	c.finishMu.Lock()
	defer c.finishMu.Unlock()
	select {
	case <-c.finish.peerCloseDone:
		c.finish = newFinishState()
	default:
	}
}

func (c *WSCallback) markSessionMetadataEvent() { c.currentFinishState().markSessionMetadata() }
func (c *WSCallback) markPeerCloseEvent()       { c.currentFinishState().markPeerClose() }

// WSChannel is a Flux TTS WebSocket client that delivers server events via Go channels.
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptionsV2
	sOptions *interfaces.SpeakV2WSOptions

	chans  []*msginterface.FluxSpeakMessageChan
	router *commoninterfaces.Router

	finishMu   sync.Mutex
	finish     *finishState
	pingActive atomic.Bool
}

// currentFinishState returns the finishState for the session in progress.
func (c *WSChannel) currentFinishState() *finishState {
	c.finishMu.Lock()
	defer c.finishMu.Unlock()
	return c.finish
}

// resetFinishStateIfClosed installs a fresh finishState when the previous
// session already ended (its peer-close milestone fired), so Finish works on a
// reconnected session. A live session's milestones are never discarded.
func (c *WSChannel) resetFinishStateIfClosed() {
	c.finishMu.Lock()
	defer c.finishMu.Unlock()
	select {
	case <-c.finish.peerCloseDone:
		c.finish = newFinishState()
	default:
	}
}

func (c *WSChannel) markSessionMetadataEvent() { c.currentFinishState().markSessionMetadata() }
func (c *WSChannel) markPeerCloseEvent()       { c.currentFinishState().markPeerClose() }
