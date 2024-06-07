// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
	"os"
	"strconv"
	"strings"

	klog "k8s.io/klog/v2"
)

func (o *ClientOptions) Parse() error {
	// general
	if o.APIKey == "" {
		if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
			klog.V(3).Infof("DEEPGRAM_API_KEY found")
			o.APIKey = v
		}
	}
	if v := os.Getenv("DEEPGRAM_HOST"); v != "" {
		klog.V(3).Infof("DEEPGRAM_HOST found")
		o.Host = v
	}
	if v := os.Getenv("DEEPGRAM_API_VERSION"); v != "" {
		klog.V(3).Infof("DEEPGRAM_API_VERSION found")
		o.APIVersion = v
	}
	if v := os.Getenv("DEEPGRAM_API_PATH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_API_PATH found")
		o.Path = v
	}
	if v := os.Getenv("DEEPGRAM_SELF_HOSTED"); v != "" {
		klog.V(3).Infof("DEEPGRAM_SELF_HOSTED found")
		o.SelfHosted = strings.EqualFold(strings.ToLower(v), "true")
	}

	// checks
	if !o.SelfHosted && o.APIKey == "" {
		klog.V(1).Infof("DEEPGRAM_API_KEY not set")
		return ErrNoAPIKey
	}

	// shared
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_SKIP_AUTH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_SKIP_AUTH found")
		o.SkipServerAuth = strings.EqualFold(strings.ToLower(v), "true")
	}

	// prerecorded
	// currently nothing

	// websocket
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_REDIRECT"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_REDIRECT found")
		o.RedirectService = strings.EqualFold(strings.ToLower(v), "true")
	}
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_KEEP_ALIVE"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_KEEP_ALIVE found")
		o.EnableKeepAlive = strings.EqualFold(strings.ToLower(v), "true")
	}

	// these require inspecting messages, therefore you must update the InspectMessage() method
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_AUTO_FLUSH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_AUTO_FLUSH found")
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			klog.V(3).Infof("DEEPGRAM_WEBSOCKET_AUTO_FLUSH set to %d", i)
			o.AutoFlushReplyDelta = i
		}
	}

	return nil
}

func (c *ClientOptions) InspectMessage() bool {
	return c.AutoFlushReplyDelta != 0
}

func (o *PreRecordedTranscriptionOptions) Check() error {
	// checks
	// currently no op

	return nil
}

func (o *LiveTranscriptionOptions) Check() error {
	// checks
	// currently no op

	return nil
}

func (o *AnalyzeOptions) Check() error {
	// checks
	// currently no op

	return nil
}

func (o *SpeakOptions) Check() error {
	// checks
	// currently no op

	return nil
}
