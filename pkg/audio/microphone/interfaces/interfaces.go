// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Microphone defines a interface for a Microphone implementation
package interfaces

import "io"

// Microphone defines a interface for a Microphone implementation
type Microphone interface {
	Start() error
	Read() ([]int16, error)
	Stream(w io.Writer) error
	Mute()
	Unmute()
	Stop() error
}
