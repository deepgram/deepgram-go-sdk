// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"os"
	"strconv"
	"strings"

	klog "k8s.io/klog/v2"
)

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

	// these require inspecting messages, therefore you must update the InspectMessage() method
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_REPLY_AUTO_FLUSH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_REPLY_AUTO_FLUSH found")
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			klog.V(3).Infof("DEEPGRAM_WEBSOCKET_REPLY_AUTO_FLUSH set to %d", i)
			o.AutoFlushReplyDelta = i
		}
	}

	// text-to-speech websocket
	// these require inspecting messages, therefore you must update the InspectMessage() method
	if v := os.Getenv("DEEPGRAM_WEBSOCKET_SPEAK_AUTO_FLUSH"); v != "" {
		klog.V(3).Infof("DEEPGRAM_WEBSOCKET_SPEAK_AUTO_FLUSH found")
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			klog.V(3).Infof("DEEPGRAM_WEBSOCKET_SPEAK_AUTO_FLUSH set to %d", i)
			o.AutoFlushSpeakDelta = i
		}
	}

	return nil
}

// InspectListenMessage returns true if the Listen message should be inspected
func (o *ClientOptions) InspectListenMessage() bool {
	return o.AutoFlushReplyDelta != 0
}

// InspectSpeakMessage returns true if the Speak message should be inspected
func (o *ClientOptions) InspectSpeakMessage() bool {
	return o.AutoFlushSpeakDelta != 0
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

func (o *WSSpeakOptions) Check() error {
	// checks
	// currently no op

	return nil
}

func NewSettingsOptions() *SettingsOptions {
	options := &SettingsOptions{
		Type: TypeSettings,
		Audio: Audio{
			Input: &Input{
				Encoding:   "linear16",
				SampleRate: 16000,
			},
			Output: &Output{
				Encoding:   "linear16",
				SampleRate: 16000,
				Container:  "none",
			},
		},
		Agent: Agent{
			Language: "en",
			Listen: Listen{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "nova-3",
				},
			},
			Think: Think{
				Provider: map[string]interface{}{
					"type":  "open_ai",
					"model": "gpt-4o-mini",
				},
			},
		},
	}

	// Set the default speak provider using the new interface{} approach
	options.Agent.Speak = Speak{
		Provider: map[string]interface{}{
			"type":  "deepgram",
			"model": "aura-2-thalia-en",
		},
	}

	return options
}
func (o *SettingsOptions) Check() error {
	// checks
	// currently no op

	return nil
}
