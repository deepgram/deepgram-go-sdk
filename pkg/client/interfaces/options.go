// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
	"os"
	"strings"

	klog "k8s.io/klog/v2"
)

func (o *ClientOptions) Parse() error {
	// general
	if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
		klog.V(3).Infof("DEEPGRAM_API_KEY found")
		o.ApiKey = v
	}
	if !o.OnPrem && o.ApiKey == "" {
		klog.V(1).Infof("DEEPGRAM_API_KEY not set")
		return ErrNoApiKey
	}
	if v := os.Getenv("DEEPGRAM_HOST"); v != "" {
		klog.V(3).Infof("DEEPGRAM_HOST found")
		o.Host = v
	}
	if v := os.Getenv("DEEPGRAM_API_VERSION"); v != "" {
		klog.V(3).Infof("DEEPGRAM_API_VERSION found")
		o.ApiVersion = v
	}
	if v := os.Getenv("DEEPGRAM_API_PATH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_API_PATH found")
		o.Path = v
	}
	if v := os.Getenv("DEEPGRAM_ON_PREM"); v != "" {
		klog.V(3).Infof("DEEPGRAM_ON_PREM found")
		o.OnPrem = strings.EqualFold(strings.ToLower(v), "true")
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

	return nil
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
