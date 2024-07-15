// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
)

// These are the message types that can be received from the live API
type TypeResponse commoninterfaces.TypeResponse

// These are the message types that can be received from the text-to-speech streaming API
const (
	// message types
	TypeOpenResponse                  = commoninterfaces.TypeOpenResponse
	TypeMetadataResponse TypeResponse = "Metadata"
	TypeFlushedResponse  TypeResponse = "Flushed"
	TypeClearedResponse  TypeResponse = "Cleared"
	TypeCloseResponse                 = commoninterfaces.TypeCloseResponse

	// "Error" type
	TypeWarningResponse   TypeResponse = "Warning"
	TypeErrorResponse                  = commoninterfaces.TypeErrorResponse
	TypeUnhandledResponse TypeResponse = "Unhandled"
)
