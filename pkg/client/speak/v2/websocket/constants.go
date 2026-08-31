// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import "time"

const (
	PackageVersion string = "v2.0"
)

const (
	DefaultConnectRetry int64 = 3
)

const (
	// MessageTypeSpeak sends text to be synthesized into the active turn
	MessageTypeSpeak string = "Speak"

	// MessageTypeFlush ends the active turn
	MessageTypeFlush string = "Flush"

	// MessageTypeInterrupt reports that the user barged in
	MessageTypeInterrupt string = "Interrupt"

	// MessageTypeConfigure sends a mid-session configuration update
	MessageTypeConfigure string = "Configure"

	// MessageTypeClose gracefully closes the session
	MessageTypeClose string = "Close"
)

const (
	pingPeriod = 60 * time.Second
)
