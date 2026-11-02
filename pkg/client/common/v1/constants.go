// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"errors"
	"regexp"
	"time"
)

// external constants
const (
	DefaultConnectRetry int64 = 3

	ChunkSize        = 1024 * 2
	TerminationSleep = 100 * time.Millisecond

	// socket errors
	FatalReadSocketErr  string = "read: can't assign requested address"
	FatalWriteSocketErr string = "write: broken pipe"
	UseOfClosedSocket   string = "use of closed network connection"
	UnknownDeepgramErr  string = "unknown deepgram error"

	// socket successful close error
	SuccessfulSocketErr string = "close 1000"
)

// errors
var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrInvalidConnection connection is not valid
	ErrInvalidConnection = errors.New("connection is not valid")

	// ErrFatalPanicRecovered fatal panic recovered
	ErrFatalPanicRecovered = errors.New("fatal panic - attempt to recover")
)

var transportErrorCodeRegexp = regexp.MustCompile(`^[A-Z][A-Z0-9]*(?:[-_][A-Z0-9]+)+\b`)

// TransportErrorCode returns the service code at the start of a WebSocket close reason.
func TransportErrorCode(reason string) string {
	return transportErrorCodeRegexp.FindString(reason)
}

// internal constants for retry, waits, back-off, etc.
const (
	defaultDelayBetweenRetry int64 = 2
)
