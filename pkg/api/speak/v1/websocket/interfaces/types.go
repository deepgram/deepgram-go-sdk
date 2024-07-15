// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/***********************************/
// Request/Input structs
/***********************************/
type SpeakOptions interfaces.SpeakOptions

/***********************************/
// MessageType is the header to bootstrap you way unmarshalling other messages
/***********************************/
/*
	Example:
	{
		"type": "message",
		"message": {
			...
		}
	}
*/
type MessageType struct {
	Type string `json:"type"`
}

// OpenResponse is the response from opening the connection
type OpenResponse = commoninterfaces.OpenResponse

// MetadataResponse is the response from the text-to-speech request which contains metadata about the request
type MetadataResponse struct {
	Type      string `json:"type,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// FlushedResponse is the response which indicates that the server has flushed the buffer and is ready to return audio
type FlushedResponse struct {
	Type       string `json:"type,omitempty"`
	SequenceID int    `json:"sequence_id,omitempty"`
}

// ClearedResponse is the response which indicates that the server has cleared the buffer
type ClearedResponse struct {
	Type       string `json:"type,omitempty"`
	SequenceID int    `json:"sequence_id,omitempty"`
}

// CloseResponse is the response from closing the connection
type CloseResponse = commoninterfaces.CloseResponse

// WarningResponse is the Deepgram specific response warning
type WarningResponse = interfaces.DeepgramWarning

// ErrorResponse is the Deepgram specific response error
type ErrorResponse = interfaces.DeepgramError
