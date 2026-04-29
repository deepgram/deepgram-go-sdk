// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import "time"

const (
	PackageVersion string = "v2.0"
)

const (
	DefaultConnectRetry int64 = 3

	ChunkSize        = 1024 * 2
	TerminationSleep = 100 * time.Millisecond
)

const (
	// MessageTypeCloseStream terminates the session
	MessageTypeCloseStream string = "CloseStream"

	// MessageTypeConfigure sends a mid-session configuration update
	MessageTypeConfigure string = "Configure"
)

const (
	pingPeriod = 60 * time.Second
)
