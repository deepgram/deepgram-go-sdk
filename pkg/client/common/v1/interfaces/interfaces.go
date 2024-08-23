// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines interfaces for the live API
package interfacesv1

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
