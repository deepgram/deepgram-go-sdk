// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (
	"errors"
	"time"
	// gabs "github.com/Jeffail/gabs/v2"
)

// internal constants for retry, waits, back-off, etc.
const (
	flushPeriod = 500 * time.Millisecond

	pingPeriod = 5 * time.Second

	defaultDelayBetweenRetry int64 = 2
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
	// MessageTypeFinalize flushes the transcription from the server
	MessageTypeFinalize string = "Finalize"

	// MessageTypeCloseStream closes the stream
	MessageTypeCloseStream string = "CloseStream"
)

// errors
var (
	// ErrInvalidConnection connection is not valid
	ErrInvalidConnection = errors.New("connection is not valid")
)
