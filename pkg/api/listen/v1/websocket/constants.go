// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import "errors"

const (
	PackageVersion string = "v1.0"
)

var (
	// ErrInvalidMessageType invalid message type
	ErrInvalidMessageType = errors.New("invalid message type")

	// ErrUserCallbackNotDefined user callback not defined or invalid
	ErrUserCallbackNotDefined = errors.New("user callback not defined or invalid")

	// ErrUserChanNotDefined user chan not defined or invalid
	ErrUserChanNotDefined = errors.New("user chan not defined or invalid")
)
