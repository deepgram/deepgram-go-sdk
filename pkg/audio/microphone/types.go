// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package microphone

import (
	"sync"

	"github.com/gordonklaus/portaudio"
)

// AudioConfig init config for library
type AudioConfig struct {
	InputChannels int
	SamplingRate  float32
}

// Microphone is a representation of a microphone
type Microphone struct {
	// microphone
	stream *portaudio.Stream

	// buffer
	intBuf []int16

	// operational
	stopChan chan struct{}
	mute     sync.Mutex
	muted    bool
}
