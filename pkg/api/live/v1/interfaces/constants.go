// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package legacy

import (
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
)

const (
	TypeOpenResponse          = interfacesv1.TypeOpenResponse
	TypeMessageResponse       = interfacesv1.TypeMessageResponse
	TypeMetadataResponse      = interfacesv1.TypeMetadataResponse
	TypeUtteranceEndResponse  = interfacesv1.TypeUtteranceEndResponse
	TypeSpeechStartedResponse = interfacesv1.TypeSpeechStartedResponse
	TypeFinalizeResponse      = interfacesv1.TypeFinalizeResponse
	TypeCloseStreamResponse   = interfacesv1.TypeCloseStreamResponse
	TypeCloseResponse         = interfacesv1.TypeCloseResponse
	TypeErrorResponse         = interfacesv1.TypeErrorResponse
)
