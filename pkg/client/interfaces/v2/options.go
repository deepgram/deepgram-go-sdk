// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv2

import (
	"os"
	"strings"

	"k8s.io/klog/v2"
)

func (o *FluxTranscriptionOptions) Check() error {
	// currently no op
	return nil
}

func (o *ClientOptions) Parse() error {
	// Priority-based credential resolution for authentication
	// 1. Explicit AccessToken parameter (highest priority)
	// 2. Explicit APIKey parameter
	// 3. DEEPGRAM_ACCESS_TOKEN environment variable
	// 4. DEEPGRAM_API_KEY environment variable (lowest priority)

	// Thread-safe credential assignment
	o.credentialsMutex.Lock()
	if o.AccessToken == "" {
		if v := os.Getenv("DEEPGRAM_ACCESS_TOKEN"); v != "" {
			klog.V(3).Infof("DEEPGRAM_ACCESS_TOKEN found")
			o.AccessToken = v
		}
	}

	if o.APIKey == "" {
		if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
			klog.V(3).Infof("DEEPGRAM_API_KEY found")
			o.APIKey = v
		}
	}
	o.credentialsMutex.Unlock()

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

	// checks - ensure we have some form of authentication unless self-hosted
	// Use thread-safe access to check credentials
	o.credentialsMutex.RLock()
	hasCredentials := o.AccessToken != "" || o.APIKey != ""
	o.credentialsMutex.RUnlock()

	if !o.SelfHosted && !hasCredentials {
		klog.V(1).Infof("Neither DEEPGRAM_ACCESS_TOKEN nor DEEPGRAM_API_KEY is set")
		return ErrNoAPIKey // Using existing error for backward compatibility
	}

	// shared
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_SKIP_AUTH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_SKIP_AUTH found")
		o.SkipServerAuth = strings.EqualFold(strings.ToLower(v), "true")
	}

	// prerecorded
	// currently nothing

	// speech-to-text websocket
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
