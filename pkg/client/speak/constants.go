// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package speak

import (
	"errors"
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

const (
	// MessageTypeFlush flushes the audio from the server
	MessageTypeFlush string = "Flush"

	// MessageTypeReset resets the text buffer
	MessageTypeReset string = "Reset"

	// MessageTypeClose closes the stream
	MessageTypeClose string = "Close"
)

// errors
var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrInvalidConnection connection is not valid
	ErrInvalidConnection = errors.New("connection is not valid")
)

// internal constants for retry, waits, back-off, etc.
const (
	defaultDelayBetweenRetry int64 = 2
)
