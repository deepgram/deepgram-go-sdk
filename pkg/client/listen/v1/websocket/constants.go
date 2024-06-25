// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"errors"
	"time"
)

const (
	PackageVersion string = "v1.0"
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
	// MessageTypeKeepAlive keep the connection alive
	MessageTypeKeepAlive string = "KeepAlive"

	// MessageTypeFinalize flushes the transcription from the server
	MessageTypeFinalize string = "Finalize"
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
	flushPeriod = 500 * time.Millisecond
	pingPeriod  = 5 * time.Second

	defaultDelayBetweenRetry int64 = 2
)
