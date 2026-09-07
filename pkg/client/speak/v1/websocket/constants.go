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
	// MessageTypeSpeak adds text to the server's synthesis buffer.
	MessageTypeSpeak string = "Speak"

	// MessageTypeFlush flushes the server's synthesis buffer.
	MessageTypeFlush string = "Flush"

	// MessageTypeClear clears the server's synthesis buffer.
	MessageTypeClear string = "Clear"

	// MessageTypeReset is kept for compatibility but emits the supported Clear message.
	// Deprecated: Use MessageTypeClear instead.
	MessageTypeReset = MessageTypeClear

	// MessageTypeClose closes the stream
	MessageTypeClose string = "Close"
)

// internal constants for retry, waits, back-off, etc.
const (
	flushPeriod = 500 * time.Millisecond
)
