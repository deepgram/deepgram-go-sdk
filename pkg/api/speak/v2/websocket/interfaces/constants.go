// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv2

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
)

// TypeResponse is the type of server-sent message, decoded from the "type" JSON field.
type TypeResponse commoninterfaces.TypeResponse

const (
	// Shared with common WebSocket infrastructure
	TypeOpenResponse  = commoninterfaces.TypeOpenResponse  // "Open"
	TypeCloseResponse = commoninterfaces.TypeCloseResponse // "Close"
	TypeErrorResponse = commoninterfaces.TypeErrorResponse // "Error"

	// Flux TTS server messages (wire "type" field values)
	TypeConnectedResponse         TypeResponse = "Connected"
	TypeSpeechStartedResponse     TypeResponse = "SpeechStarted"
	TypeSpeechMetadataResponse    TypeResponse = "SpeechMetadata"
	TypeSpeechInterruptedResponse TypeResponse = "SpeechInterrupted"
	TypeFlushedResponse           TypeResponse = "Flushed"
	TypeSessionMetadataResponse   TypeResponse = "SessionMetadata"
	TypeConfigureSuccess          TypeResponse = "ConfigureSuccess"
	TypeConfigureFailure          TypeResponse = "ConfigureFailure"
	TypeWarningResponse           TypeResponse = "Warning"
	TypeFatalError                TypeResponse = "Error" // server sends {"type":"Error"} for fatal errors
)

// PlaybackOffsetTypeTimeMs is the only supported playback-offset unit for
// Interrupt messages: milliseconds of session audio played.
const PlaybackOffsetTypeTimeMs = "time_ms"
