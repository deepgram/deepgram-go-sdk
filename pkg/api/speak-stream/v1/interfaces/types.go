// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
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

// SpeakStreamResponse is the response from the text-to-speech request
type SpeakStreamResponse struct {
	ContentType string `json:"content_type,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
	ModelUUID   string `json:"model_uuid,omitempty"`
	ModelName   string `json:"model_name,omitempty"`
	Date        string `json:"date,omitempty"`
}

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

// OpenResponse is the response from the connection opening
type OpenResponse struct {
	Type string `json:"type,omitempty"`
}

// CloseResponse is the response from the connection closing
type CloseResponse struct {
	Type string `json:"type,omitempty"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
