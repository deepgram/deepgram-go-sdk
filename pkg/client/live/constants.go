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
	pingPeriod = 5 * time.Second

	connectionRetryInfinite  int64 = 0
	defaultConnectRetry      int64 = 3
	defaultDelayBetweenRetry int64 = 2

	ChunkSize        = 1024 * 2
	TerminationSleep = 100 * time.Millisecond
)

// errors
var (
	// ErrInvalidConnection connection is not valid
	ErrInvalidConnection = errors.New("connection is not valid")
)
