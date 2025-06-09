// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
)

// These are the message types that can be received from the live API
type TypeResponse commoninterfaces.TypeResponse

const (
	// message types
	TypeOpenResponse                       = commoninterfaces.TypeOpenResponse
	TypeMessageResponse       TypeResponse = "Results"
	TypeMetadataResponse      TypeResponse = "Metadata"
	TypeUtteranceEndResponse  TypeResponse = "UtteranceEnd"
	TypeSpeechStartedResponse TypeResponse = "SpeechStarted"
	TypeFinalizeResponse      TypeResponse = "Finalize"
	TypeCloseResponse                      = commoninterfaces.TypeCloseResponse
	TypeCloseStreamResponse   TypeResponse = "CloseStream"
	TypeErrorResponse                      = commoninterfaces.TypeErrorResponse
)
