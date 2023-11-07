// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Implementation for a replay device. In this case, replays an audio file to stream into a listener
package replay

import (
	"errors"
)

const (
	defaultBytesToRead int = 2048
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
