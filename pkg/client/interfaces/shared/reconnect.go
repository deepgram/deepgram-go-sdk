// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import "time"

// ReconnectEventType identifies a lifecycle event produced by an automatic
// WebSocket reconnection policy.
type ReconnectEventType string

const (
	// ReconnectAttempt is emitted immediately before each connection attempt.
	ReconnectAttempt ReconnectEventType = "attempt"
	// ReconnectSuccess is emitted after a replacement connection is ready.
	ReconnectSuccess ReconnectEventType = "success"
	// ReconnectFailure is emitted after all attempts have failed.
	ReconnectFailure ReconnectEventType = "failure"
)

// ReconnectEvent reports the state of an automatic reconnection attempt.
// Attempt starts at one. Error is populated for terminal failures.
type ReconnectEvent struct {
	Type    ReconnectEventType
	Attempt int
	Error   error
}

// ReconnectPolicy enables automatic recovery from unexpected WebSocket
// transport failures. A nil policy preserves the SDK's existing behavior:
// callers must explicitly call AttemptReconnect.
//
// MaxRetries is the number of replacement connection attempts; zero uses the
// default of three. InitialDelay is waited before the first attempt; zero uses
// one second. Later delays double up to MaxDelay, which defaults to thirty
// seconds. Jitter is the proportional random variation applied to each delay;
// zero uses the default of 0.2.
//
// OnReconnect runs synchronously in the recovery goroutine and should return
// promptly. The SDK never replays caller-provided audio, control, or
// application frames.
type ReconnectPolicy struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Jitter       float64
	OnReconnect  func(ReconnectEvent)
}
