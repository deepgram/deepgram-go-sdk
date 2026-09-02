// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
)

// finishMarker is how observers report graceful-close milestones. The clients
// implement it against their CURRENT session's finishState, so milestones from a
// previous session (before a reconnect) cannot satisfy a later Finish.
type finishMarker interface {
	markSessionMetadataEvent()
	markPeerCloseEvent()
}

// callbackFinishObserver wraps a user FluxSpeakMessageCallback, forwarding every
// event unchanged while marking the graceful-close milestones (SessionMetadata and
// peer closure) on the client's finishState. The user callback runs first, so when
// Finish(ctx) returns the handler has already seen every event up to and including
// SessionMetadata.
type callbackFinishObserver struct {
	next   msginterfaces.FluxSpeakMessageCallback
	marker finishMarker
}

func (o *callbackFinishObserver) Open(or *msginterfaces.OpenResponse) error {
	return o.next.Open(or)
}

func (o *callbackFinishObserver) Connected(cr *msginterfaces.ConnectedResponse) error {
	return o.next.Connected(cr)
}

func (o *callbackFinishObserver) Binary(byMsg []byte) error {
	return o.next.Binary(byMsg)
}

func (o *callbackFinishObserver) SpeechStarted(ss *msginterfaces.SpeechStartedResponse) error {
	return o.next.SpeechStarted(ss)
}

func (o *callbackFinishObserver) SpeechMetadata(sm *msginterfaces.SpeechMetadataResponse) error {
	return o.next.SpeechMetadata(sm)
}

func (o *callbackFinishObserver) SpeechInterrupted(si *msginterfaces.SpeechInterruptedResponse) error {
	return o.next.SpeechInterrupted(si)
}

func (o *callbackFinishObserver) Flushed(fl *msginterfaces.FlushedResponse) error {
	return o.next.Flushed(fl)
}

func (o *callbackFinishObserver) SessionMetadata(sm *msginterfaces.SessionMetadataResponse) error {
	err := o.next.SessionMetadata(sm)
	o.marker.markSessionMetadataEvent()
	return err
}

func (o *callbackFinishObserver) ConfigureSuccess(cs *msginterfaces.ConfigureSuccessResponse) error {
	return o.next.ConfigureSuccess(cs)
}

func (o *callbackFinishObserver) ConfigureFailure(cf *msginterfaces.ConfigureFailureResponse) error {
	return o.next.ConfigureFailure(cf)
}

func (o *callbackFinishObserver) Warning(wr *msginterfaces.WarningResponse) error {
	return o.next.Warning(wr)
}

func (o *callbackFinishObserver) FatalError(fe *msginterfaces.FatalErrorResponse) error {
	return o.next.FatalError(fe)
}

func (o *callbackFinishObserver) Close(cr *msginterfaces.CloseResponse) error {
	err := o.next.Close(cr)
	o.marker.markPeerCloseEvent()
	return err
}

func (o *callbackFinishObserver) Error(er *msginterfaces.ErrorResponse) error {
	return o.next.Error(er)
}

func (o *callbackFinishObserver) UnhandledEvent(byData []byte) error {
	return o.next.UnhandledEvent(byData)
}

// chanFinishObserver wraps a user FluxSpeakMessageChan, exposing the user's
// channels unchanged while adding internal SessionMetadata and Close channels the
// client uses to mark the graceful-close milestones on its finishState.
type chanFinishObserver struct {
	next          msginterfaces.FluxSpeakMessageChan
	sessionMetaCh chan *msginterfaces.SessionMetadataResponse
	closeCh       chan *msginterfaces.CloseResponse
}

func newChanFinishObserver(next msginterfaces.FluxSpeakMessageChan) *chanFinishObserver {
	return &chanFinishObserver{
		next: next,
		// buffered so the router's blocking fan-out never stalls on the internal
		// observer, even after its consumer goroutine has exited
		sessionMetaCh: make(chan *msginterfaces.SessionMetadataResponse, 4),
		closeCh:       make(chan *msginterfaces.CloseResponse, 4),
	}
}

// run consumes the internal channels and marks the client's current-session
// milestones until the client context is canceled.
func (o *chanFinishObserver) run(ctx context.Context, marker finishMarker) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-o.sessionMetaCh:
				marker.markSessionMetadataEvent()
			case <-o.closeCh:
				marker.markPeerCloseEvent()
			}
		}
	}()
}

func (o *chanFinishObserver) GetOpen() []*chan *msginterfaces.OpenResponse {
	return o.next.GetOpen()
}

func (o *chanFinishObserver) GetConnected() []*chan *msginterfaces.ConnectedResponse {
	return o.next.GetConnected()
}

func (o *chanFinishObserver) GetBinary() []*chan *[]byte {
	return o.next.GetBinary()
}

func (o *chanFinishObserver) GetSpeechStarted() []*chan *msginterfaces.SpeechStartedResponse {
	return o.next.GetSpeechStarted()
}

func (o *chanFinishObserver) GetSpeechMetadata() []*chan *msginterfaces.SpeechMetadataResponse {
	return o.next.GetSpeechMetadata()
}

func (o *chanFinishObserver) GetSpeechInterrupted() []*chan *msginterfaces.SpeechInterruptedResponse {
	return o.next.GetSpeechInterrupted()
}

func (o *chanFinishObserver) GetFlushed() []*chan *msginterfaces.FlushedResponse {
	return o.next.GetFlushed()
}

func (o *chanFinishObserver) GetSessionMetadata() []*chan *msginterfaces.SessionMetadataResponse {
	user := o.next.GetSessionMetadata()
	out := make([]*chan *msginterfaces.SessionMetadataResponse, 0, len(user)+1)
	out = append(out, user...)
	return append(out, &o.sessionMetaCh)
}

func (o *chanFinishObserver) GetConfigureSuccess() []*chan *msginterfaces.ConfigureSuccessResponse {
	return o.next.GetConfigureSuccess()
}

func (o *chanFinishObserver) GetConfigureFailure() []*chan *msginterfaces.ConfigureFailureResponse {
	return o.next.GetConfigureFailure()
}

func (o *chanFinishObserver) GetWarning() []*chan *msginterfaces.WarningResponse {
	return o.next.GetWarning()
}

func (o *chanFinishObserver) GetFatalError() []*chan *msginterfaces.FatalErrorResponse {
	return o.next.GetFatalError()
}

func (o *chanFinishObserver) GetClose() []*chan *msginterfaces.CloseResponse {
	user := o.next.GetClose()
	out := make([]*chan *msginterfaces.CloseResponse, 0, len(user)+1)
	out = append(out, user...)
	return append(out, &o.closeCh)
}

func (o *chanFinishObserver) GetError() []*chan *msginterfaces.ErrorResponse {
	return o.next.GetError()
}

func (o *chanFinishObserver) GetUnhandled() []*chan *[]byte {
	return o.next.GetUnhandled()
}
