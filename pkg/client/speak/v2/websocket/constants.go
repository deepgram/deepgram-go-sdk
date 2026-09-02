// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"errors"
	"regexp"
	"time"
)

const (
	PackageVersion string = "v2.0"
)

// ErrGracefulCloseIncomplete is returned by Finish when the server closes the
// connection before delivering the final SessionMetadata event.
var ErrGracefulCloseIncomplete = errors.New("connection closed before the final SessionMetadata was received")

// wsErrorRegexp decomposes gorilla-style websocket transport errors into
// code/number/description for typed error responses.
var wsErrorRegexp = regexp.MustCompile(`websocket: ([a-z]+) (\d+) .+: (.+)`)

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
	// DefaultKeepAlivePeriod is the default interval between WebSocket
	// protocol-level keepalive pings when ClientOptions.EnableKeepAlive is set.
	// The /v2/speak server closes idle sessions after 60 seconds (NET-0004), so
	// pings are sent at half that deadline to reach the server comfortably before
	// it. Override with ClientOptions.KeepAlivePeriod.
	DefaultKeepAlivePeriod = 30 * time.Second

	// ServerIdleTimeout is the server-side /v2/speak idle deadline: sessions idle
	// for this long are closed with NET-0004. Exported so tests can assert the
	// keepalive interval stays well inside it.
	ServerIdleTimeout = 60 * time.Second
)
