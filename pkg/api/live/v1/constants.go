// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import "errors"

var (
	// ErrInvalidMessageType invalid message type
	ErrInvalidMessageType = errors.New("invalid message type")

	// ErrUserCallbackNotDefined user callback object not defined
	ErrUserCallbackNotDefined = errors.New("user callback object not defined")
)
