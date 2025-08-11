// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
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
	// MessageTypeFlush flushes the audio from the server
	MessageTypeSpeak string = "Speak"

	// MessageTypeFlush flushes the audio from the server
	MessageTypeFlush string = "Flush"

	// MessageTypeClear clears the audio from the server
	MessageTypeClear string = "Clear"

	// MessageTypeReset resets the text buffer
	MessageTypeReset string = "Reset"

	// MessageTypeClose closes the stream
	MessageTypeClose string = "Close"
)

// internal constants for retry, waits, back-off, etc.
const (
	flushPeriod = 500 * time.Millisecond
)
