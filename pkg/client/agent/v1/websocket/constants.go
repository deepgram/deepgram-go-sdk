// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
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

var (
	// ErrInvalidUpdateListen is returned when UpdateListen has no provider.
	ErrInvalidUpdateListen = errors.New("agent UpdateListen requires a listen provider")
	// ErrInvalidUpdateThink is returned when UpdateThink has no provider.
	ErrInvalidUpdateThink = errors.New("agent UpdateThink requires a think provider")
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
)

// internal constants for retry, waits, back-off, etc.
const (
	flushPeriod = 500 * time.Millisecond
	pingPeriod  = 5 * time.Second
)
