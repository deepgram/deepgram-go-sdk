// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

// These are the message types that can be received from the live API
const (
	// message types
	TypeMessageResponse       string = "Results"
	TypeMetadataResponse      string = "Metadata"
	TypeUtteranceEndResponse  string = "UtteranceEnd"
	TypeSpeechStartedResponse string = "SpeechStarted"

	// Error type
	TypeErrorResponse string = "Error"
)
