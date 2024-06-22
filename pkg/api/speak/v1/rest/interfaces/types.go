// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/***********************************/
// Request/Input structs
/***********************************/
type SpeakOptions interfaces.SpeakOptions

/***********************************/
// response/result structs
/***********************************/
type SpeakResponse struct {
	ContextType      string `json:"content_type,omitempty"`
	RequestID        string `json:"request_id,omitempty"`
	ModelUUID        string `json:"model_uuid,omitempty"`
	Characters       int    `json:"characters,omitempty"`
	ModelName        string `json:"model_name,omitempty"`
	TransferEncoding string `json:"transfer_encoding,omitempty"`
	Date             string `json:"date,omitempty"`
	Filename         string `json:"filename,omitempty"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
