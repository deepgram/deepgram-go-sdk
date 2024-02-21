// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package handles the versioning in the API for the various clients (prerecorded, live, etc.)
package version

import "errors"

const (
	// APIProtocol default protocol
	APIProtocol string = "https"

	// WSProtocol default protocol
	WSProtocol string = "wss"
)

var (
	// ErrInvalidPath invalid path
	ErrInvalidPath = errors.New("invalid path")
)
