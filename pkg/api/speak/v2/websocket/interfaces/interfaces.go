// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package interfacesv2 defines the handler contracts for the Deepgram Flux TTS
// (v2/speak) WebSocket API.
package interfacesv2

// FluxSpeakMessageCallback is the callback-based handler interface for Flux TTS
// events. Implement all methods and pass your implementation to NewWSUsingCallback.
type FluxSpeakMessageCallback interface {
	// Open is called when the local WebSocket connection is established (SDK event).
	Open(or *OpenResponse) error

	// Connected is called when the server sends Connected — the first Deepgram
	// server message, confirming the session is ready to receive text.
	Connected(cr *ConnectedResponse) error

	// Binary is called for each binary WebSocket frame: audio data in the format
	// specified by the connection parameters.
	Binary(byMsg []byte) error

	// SpeechStarted is called at the start of each new turn, before audio streaming begins.
	SpeechStarted(ss *SpeechStartedResponse) error

	// SpeechMetadata is called at turn boundaries, after all audio for the turn has
	// been sent, with billing and timing for the completed turn.
	SpeechMetadata(sm *SpeechMetadataResponse) error

	// SpeechInterrupted is called in reply to an Interrupt, once synthesis has
	// concluded and audio generation has stopped.
	SpeechInterrupted(si *SpeechInterruptedResponse) error

	// Flushed is called as the immediate echo on receipt of a manual Flush.
	Flushed(fl *FlushedResponse) error

	// SessionMetadata is called with the final server message before the WebSocket
	// closes, reporting cumulative session totals.
	SessionMetadata(sm *SessionMetadataResponse) error

	// ConfigureSuccess is called when a mid-session Configure message was accepted.
	ConfigureSuccess(cs *ConfigureSuccessResponse) error

	// ConfigureFailure is called when a mid-session Configure message was rejected.
	ConfigureFailure(cf *ConfigureFailureResponse) error

	// Warning is called for informational warnings; synthesis continues and the
	// connection is unaffected.
	Warning(wr *WarningResponse) error

	// FatalError is called on terminal server errors. The connection will be closed
	// by the server immediately after this event.
	FatalError(fe *FatalErrorResponse) error

	// Close is called when the WebSocket connection is closed (SDK event).
	Close(cr *CloseResponse) error

	// Error is called on connection-level transport errors.
	Error(er *ErrorResponse) error

	// UnhandledEvent is called for any server message type not recognized by the router.
	UnhandledEvent(byData []byte) error
}

// FluxSpeakMessageChan is the channel-based handler interface for Flux TTS events.
// Each getter returns a slice of channel pointers, allowing multiple subscribers.
// Implement all methods and pass your implementation to NewWSUsingChan.
type FluxSpeakMessageChan interface {
	// GetOpen returns channels to receive the local connection-opened event.
	GetOpen() []*chan *OpenResponse

	// GetConnected returns channels to receive Connected server messages.
	GetConnected() []*chan *ConnectedResponse

	// GetBinary returns channels to receive binary audio frames.
	GetBinary() []*chan *[]byte

	// GetSpeechStarted returns channels to receive SpeechStarted messages.
	GetSpeechStarted() []*chan *SpeechStartedResponse

	// GetSpeechMetadata returns channels to receive SpeechMetadata messages.
	GetSpeechMetadata() []*chan *SpeechMetadataResponse

	// GetSpeechInterrupted returns channels to receive SpeechInterrupted messages.
	GetSpeechInterrupted() []*chan *SpeechInterruptedResponse

	// GetFlushed returns channels to receive Flushed messages.
	GetFlushed() []*chan *FlushedResponse

	// GetSessionMetadata returns channels to receive SessionMetadata messages.
	GetSessionMetadata() []*chan *SessionMetadataResponse

	// GetConfigureSuccess returns channels to receive ConfigureSuccess messages.
	GetConfigureSuccess() []*chan *ConfigureSuccessResponse

	// GetConfigureFailure returns channels to receive ConfigureFailure messages.
	GetConfigureFailure() []*chan *ConfigureFailureResponse

	// GetWarning returns channels to receive Warning messages.
	GetWarning() []*chan *WarningResponse

	// GetFatalError returns channels to receive fatal server error messages.
	GetFatalError() []*chan *FatalErrorResponse

	// GetClose returns channels to receive the local connection-closed event.
	GetClose() []*chan *CloseResponse

	// GetError returns channels to receive connection-level transport errors.
	GetError() []*chan *ErrorResponse

	// GetUnhandled returns channels to receive raw bytes of unrecognized message types.
	GetUnhandled() []*chan *[]byte
}
