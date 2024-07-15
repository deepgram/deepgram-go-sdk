// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
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
)

const (
	// MessageTypeKeepAlive keep the connection alive
	MessageTypeKeepAlive string = "KeepAlive"

	// MessageTypeFinalize flushes the transcription from the server
	MessageTypeFinalize string = "Finalize"
)

// internal constants for retry, waits, back-off, etc.
const (
	flushPeriod = 500 * time.Millisecond
	pingPeriod  = 5 * time.Second
)
