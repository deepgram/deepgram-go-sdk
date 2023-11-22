// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (
	"errors"
	"time"
	// gabs "github.com/Jeffail/gabs/v2"
)

const (
	pingPeriod = 5 * time.Second

	connectionRetryInfinite  int64 = 0
	defaultConnectRetry      int64 = 3
	defaultDelayBetweenRetry int64 = 2

	CHUNK_SIZE        = 1024 * 2
	TERMINATION_SLEEP = 100 * time.Millisecond
)

var (
	// ErrInvalidConnection connection is not valid
	ErrInvalidConnection = errors.New("connection is not valid")
)
