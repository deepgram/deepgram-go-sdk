// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines interfaces for the live API
package interfacesv1

import "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"

/*
Shared Structs
*/
// OpenResponse is the response from the connection opening
type OpenResponse struct {
	Type string `json:"type,omitempty"`
}

// CloseResponse is the response from the connection closing
type CloseResponse struct {
	Type string `json:"type,omitempty"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse = interfaces.DeepgramError

/*
Router definition
*/
type Router interface {
	Open(or *OpenResponse) error
	Message(byMsg []byte) error
	Binary(byMsg []byte) error
	Close(or *CloseResponse) error
	Error(er *ErrorResponse) error
}

/*
WebSocketHandler this defines the things you need to implement for your specific WS protocol
*/
type WebSocketHandler interface {
	// GetURL returns the URL for the websocket connection. This has already been processed through pkg/api/version
	GetURL(host string) (string, error)

	// ProcessMessage is the entry point for processing messages based on that client's WS protocol
	ProcessMessage(wsType int, byMsg []byte) error

	// ProcessError handles any errors that occur during the WS connection
	ProcessError(err error) error

	// Start handles any setup required specific to that client's WS protocol
	Start()

	// Finish handles any cleanup required specific to that client's WS protocol
	Finish()

	// GetCloseMsg returns the message to send when closing the connection.
	// It turns out that between clients, the close message can be different.
	GetCloseMsg() []byte
}
