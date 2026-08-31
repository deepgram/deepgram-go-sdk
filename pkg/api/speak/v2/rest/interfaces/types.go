// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv2

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

/***********************************/
// Request/Input structs
/***********************************/

// SpeakOptions are the query parameters for POST /v2/speak.
type SpeakOptions = interfaces.SpeakV2Options

/***********************************/
// response/result structs
/***********************************/

// SpeakResponse carries the per-request telemetry returned as response headers
// by POST /v2/speak, mirroring Aura's REST conventions.
type SpeakResponse struct {
	ContextType      string `json:"content_type,omitempty"`
	RequestID        string `json:"request_id,omitempty"`
	ModelUUID        string `json:"model_uuid,omitempty"`
	ModelName        string `json:"model_name,omitempty"`
	Characters       int    `json:"characters,omitempty"`
	TransferEncoding string `json:"transfer_encoding,omitempty"`
	Date             string `json:"date,omitempty"`
	Filename         string `json:"filename,omitempty"`
	// Warnings carries the generic dg-warnings response header, when present.
	Warnings string `json:"warnings,omitempty"`
}

// AsyncResponse is the JSON acknowledgement returned when a callback URL is
// supplied: the request is processed asynchronously and the audio is delivered
// to the callback URL.
type AsyncResponse struct {
	RequestID string `json:"request_id"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse = interfaces.DeepgramError
