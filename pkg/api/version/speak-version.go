// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

/*
GetSpeakAPI is a function which controls the versioning of the live transcription API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for the live transcription
*/
func GetSpeakAPI(ctx context.Context, host, version, path string, options *interfaces.SpeakOptions, args ...interface{}) (string, error) {
	return getAPIURL(ctx, "speak", host, version, path, options, args...)
}

const (
	// SpeakV2APIVersion is the default API version for the Flux TTS (v2/speak) endpoints
	SpeakV2APIVersion string = "v2"

	// SpeakV2Path is the path for the Flux TTS endpoints
	SpeakV2Path string = "speak"
)

// GetSpeakV2API builds the REST URL for the Deepgram Flux TTS batch endpoint.
// It produces: https://host/v2/speak?<SpeakV2Options>
//
// Overrides:
//   - host: defaults to api.deepgram.com if empty
//   - version: defaults to "v2" if empty (can be overridden via ClientOptions.APIVersion)
//   - path: defaults to "speak" if empty (can be overridden via ClientOptions.Path)
func GetSpeakV2API(ctx context.Context, host, version, path string, options *interfaces.SpeakV2Options, args ...interface{}) (string, error) {
	if version == "" {
		version = SpeakV2APIVersion
	}
	if path == "" {
		path = SpeakV2Path
	}
	return getAPIURL(ctx, "speak", host, version, path, options, args...)
}

// GetSpeakV2StreamAPI builds the WebSocket URL for the Deepgram Flux TTS streaming endpoint.
// It produces: wss://host/v2/speak?<SpeakV2WSOptions>
//
// Overrides:
//   - host: defaults to api.deepgram.com if empty
//   - version: defaults to "v2" if empty (can be overridden via ClientOptions.APIVersion)
//   - path: defaults to "speak" if empty (can be overridden via ClientOptions.Path)
func GetSpeakV2StreamAPI(ctx context.Context, host, version, path string, options *interfaces.SpeakV2WSOptions, args ...interface{}) (string, error) {
	if version == "" {
		version = SpeakV2APIVersion
	}
	if path == "" {
		path = SpeakV2Path
	}
	return getAPIURL(ctx, APITypeSpeakStream, host, version, path, options, args...)
}
