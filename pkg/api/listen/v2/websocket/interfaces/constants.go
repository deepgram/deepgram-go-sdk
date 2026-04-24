// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
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

	// Flux-specific server messages (wire "type" field values)
	TypeConnectedResponse TypeResponse = "Connected"
	TypeTurnInfoResponse  TypeResponse = "TurnInfo"
	TypeConfigureSuccess  TypeResponse = "ConfigureSuccess"
	TypeConfigureFailure  TypeResponse = "ConfigureFailure"
	TypeFatalError        TypeResponse = "Error" // server sends {"type":"Error"} for fatal errors
)

// TurnEvent values for TurnInfoResponse.EventType
const (
	TurnEventStartOfTurn    = "StartOfTurn"
	TurnEventUpdate         = "Update"
	TurnEventEagerEndOfTurn = "EagerEndOfTurn"
	TurnEventTurnResumed    = "TurnResumed"
	TurnEventEndOfTurn      = "EndOfTurn"
)
