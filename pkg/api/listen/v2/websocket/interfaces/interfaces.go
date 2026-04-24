// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package interfacesv2 defines the handler contracts for the Deepgram Flux
// (v2/listen) WebSocket API.
package interfacesv2

// FluxMessageCallback is the callback-based handler interface for Flux events.
// Implement all methods and pass your implementation to NewWSUsingCallback.
//
// Usage:
//
//	type MyHandler struct{}
//
//	func (h *MyHandler) Open(or *OpenResponse) error           { return nil }
//	func (h *MyHandler) Connected(cr *ConnectedResponse) error { fmt.Println("connected:", cr.RequestID); return nil }
//	func (h *MyHandler) TurnInfo(tr *TurnInfoResponse) error {
//	    switch tr.EventType {
//	    case TurnEventEndOfTurn:
//	        fmt.Println("Final:", tr.Transcript)
//	    case TurnEventUpdate:
//	        fmt.Println("Interim:", tr.Transcript)
//	    }
//	    return nil
//	}
//	func (h *MyHandler) ConfigureSuccess(cs *ConfigureSuccessResponse) error { return nil }
//	func (h *MyHandler) ConfigureFailure(cf *ConfigureFailureResponse) error { return nil }
//	func (h *MyHandler) FatalError(fe *FatalErrorResponse) error             { return nil }
//	func (h *MyHandler) Close(cr *CloseResponse) error                       { return nil }
//	func (h *MyHandler) Error(er *ErrorResponse) error                       { return nil }
//	func (h *MyHandler) UnhandledEvent(byData []byte) error                  { return nil }
type FluxMessageCallback interface {
	// Open is called when the local WebSocket connection is established (SDK event).
	Open(or *OpenResponse) error

	// Connected is called when the server sends ListenV2Connected — the first Deepgram
	// server message, confirming the session is ready to receive audio.
	Connected(cr *ConnectedResponse) error

	// TurnInfo is called for all turn lifecycle events. Inspect tr.EventType:
	//   TurnEventStartOfTurn    — new speech turn detected
	//   TurnEventUpdate         — interim transcript update
	//   TurnEventEagerEndOfTurn — early end-of-turn signal
	//   TurnEventTurnResumed    — turn continued after EagerEndOfTurn
	//   TurnEventEndOfTurn      — final transcript for the turn
	TurnInfo(tr *TurnInfoResponse) error

	// ConfigureSuccess is called when a mid-session Configure message was accepted.
	ConfigureSuccess(cs *ConfigureSuccessResponse) error

	// ConfigureFailure is called when a mid-session Configure message was rejected.
	ConfigureFailure(cf *ConfigureFailureResponse) error

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

// FluxMessageChan is the channel-based handler interface for Flux events.
// Each getter returns a slice of channel pointers, allowing multiple subscribers.
// Implement all methods and pass your implementation to NewWSUsingChan.
//
// Usage:
//
//	type MyHandler struct {
//	    connectedChan chan *ConnectedResponse
//	    turnInfoChan  chan *TurnInfoResponse
//	    closeChan     chan *CloseResponse
//	    errorChan     chan *ErrorResponse
//	    // ... etc
//	}
//
//	func (h *MyHandler) GetOpen() []*chan *OpenResponse              { return nil }
//	func (h *MyHandler) GetConnected() []*chan *ConnectedResponse    { return []*chan *ConnectedResponse{&h.connectedChan} }
//	func (h *MyHandler) GetTurnInfo() []*chan *TurnInfoResponse      { return []*chan *TurnInfoResponse{&h.turnInfoChan} }
//	func (h *MyHandler) GetConfigureSuccess() []*chan *ConfigureSuccessResponse { return nil }
//	func (h *MyHandler) GetConfigureFailure() []*chan *ConfigureFailureResponse { return nil }
//	func (h *MyHandler) GetFatalError() []*chan *FatalErrorResponse  { return nil }
//	func (h *MyHandler) GetClose() []*chan *CloseResponse            { return []*chan *CloseResponse{&h.closeChan} }
//	func (h *MyHandler) GetError() []*chan *ErrorResponse            { return []*chan *ErrorResponse{&h.errorChan} }
//	func (h *MyHandler) GetUnhandled() []*chan *[]byte               { return nil }
type FluxMessageChan interface {
	// GetOpen returns channels to receive the local connection-opened event.
	GetOpen() []*chan *OpenResponse

	// GetConnected returns channels to receive ListenV2Connected server messages.
	GetConnected() []*chan *ConnectedResponse

	// GetTurnInfo returns channels to receive ListenV2TurnInfo messages.
	// Inspect the received TurnInfoResponse.EventType field to distinguish event types.
	GetTurnInfo() []*chan *TurnInfoResponse

	// GetConfigureSuccess returns channels to receive ListenV2ConfigureSuccess messages.
	GetConfigureSuccess() []*chan *ConfigureSuccessResponse

	// GetConfigureFailure returns channels to receive ListenV2ConfigureFailure messages.
	GetConfigureFailure() []*chan *ConfigureFailureResponse

	// GetFatalError returns channels to receive ListenV2FatalError messages.
	GetFatalError() []*chan *FatalErrorResponse

	// GetClose returns channels to receive the local connection-closed event.
	GetClose() []*chan *CloseResponse

	// GetError returns channels to receive connection-level transport errors.
	GetError() []*chan *ErrorResponse

	// GetUnhandled returns channels to receive raw bytes of unrecognized message types.
	GetUnhandled() []*chan *[]byte
}
