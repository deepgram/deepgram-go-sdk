// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package replay

import (
	"os"
	"sync"

	wav "github.com/youpy/go-wav"
)

// ReplayOpts defines options for this device
type ReplayOptions struct {
	FullFilename string
}

// Client is a replay device. In this case, an audio stream.
type Client struct {
	options ReplayOptions

	// wav
	file    *os.File
	decoder *wav.Reader

	// operational stuff
	stopChan chan struct{}
	mute     sync.Mutex
	muted    bool
}
