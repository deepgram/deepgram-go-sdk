// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package prerecorded

import (
	"errors"
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrInvalidURIExtension couldn't find a period to indicate a file extension
	ErrInvalidURIExtension = errors.New("couldn't find a period to indicate a file extension")
)
