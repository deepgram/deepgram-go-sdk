// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package legacy

import (
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
)

/***********************************/
// Request/Input structs
/***********************************/
type LiveOptions = interfacesv1.LiveTranscriptionOptions

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
type MessageType = interfacesv1.MessageType

/***********************************/
// shared/common structs
/***********************************/
// Word is a single word in a transcript
type Word = interfacesv1.Word

// Alternative is a single alternative in a transcript
type Alternative = interfacesv1.Alternative

// Channel is a single channel in a transcript
type Channel = interfacesv1.Channel

// ModelInfo is the model information for a transcript
type ModelInfo = interfacesv1.ModelInfo

// Metadata is the metadata for a transcript
type Metadata = interfacesv1.Metadata

/***********************************/
// Request/Input structs
/***********************************/
type LiveTranscriptionOptions interfacesv1.LiveTranscriptionOptions

/***********************************/
// Results from Live Transcription
/***********************************/
type OpenResponse = interfacesv1.OpenResponse
type MessageResponse = interfacesv1.MessageResponse
type MetadataResponse = interfacesv1.MetadataResponse
type UtteranceEndResponse = interfacesv1.UtteranceEndResponse
type SpeechStartedResponse = interfacesv1.SpeechStartedResponse
type CloseResponse = interfacesv1.CloseResponse
type ErrorResponse = interfacesv1.ErrorResponse
