// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Implementation of a common functions used in portaudio
package common

import (
	"github.com/gordonklaus/portaudio"
	klog "k8s.io/klog/v2"
)

var (
	bHasBeenInitialized = false
)

// Initialize inits the library. This handles OS level init of the library.
func Initialize() {
	if bHasBeenInitialized {
		klog.V(3).Infof("portaudio has already been initialized\n")
		return
	}

	err := portaudio.Initialize()
	if err != nil {
		klog.V(1).Infof("portaudio.Initialize failed. Err: %v\n", err)
		return
	}

	klog.V(4).Infof("portaudio.Initialize succeeded\n")
	bHasBeenInitialized = true
}

// Teardown cleans up the library. This handles OS level cleanup.
func Teardown() {
	if !bHasBeenInitialized {
		klog.V(3).Infof("portaudio is not initialized\n")
		return
	}

	err := portaudio.Terminate()
	if err != nil {
		klog.V(1).Infof("portaudio.Terminate failed. Err: %v\n", err)
		return
	}

	klog.V(4).Infof("portaudio.Terminate succeeded\n")
	bHasBeenInitialized = false
}
