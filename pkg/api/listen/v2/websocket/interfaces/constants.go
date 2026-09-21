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
	TypeFatalError        TypeResponse = "Error"   // server sends {"type":"Error"} for fatal errors
	TypeWarningResponse   TypeResponse = "Warning" // non-fatal informational message
)

// Warning codes for WarningResponse.Code. The set is open: the server may add
// new values over time.
const (
	// WarningCodeForceEndTurnNoActiveTurn — a ForceEndTurn message arrived when no
	// turn was active (before StartOfTurn or after EndOfTurn). The message is
	// ignored: no EndOfTurn is emitted and turn_index does not advance. Timing
	// races between an external end-of-turn signal and the server's StartOfTurn
	// are normal, so this is informational, not an error.
	WarningCodeForceEndTurnNoActiveTurn = "FORCE_END_TURN_NO_ACTIVE_TURN"
)

// TurnEvent values for TurnInfoResponse.EventType
const (
	TurnEventStartOfTurn    = "StartOfTurn"
	TurnEventUpdate         = "Update"
	TurnEventEagerEndOfTurn = "EagerEndOfTurn"
	TurnEventTurnResumed    = "TurnResumed"
	TurnEventEndOfTurn      = "EndOfTurn"
)

// TurnTrigger values for TurnInfoResponse.Trigger on "EndOfTurn" events.
// The field is an open string: the server may add new values over time.
const (
	// TurnTriggerModel — the turn ended via Flux's native end-of-turn detection.
	TurnTriggerModel = "model"
	// TurnTriggerManual — the turn ended because the client sent ForceEndTurn.
	TurnTriggerManual = "manual"
	// TurnTriggerTimeout — the turn ended because eot_timeout_ms elapsed.
	TurnTriggerTimeout = "timeout"
)
