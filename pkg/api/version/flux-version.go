// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

const (
	// FluxAPIVersion is the default API version for the Flux endpoint
	FluxAPIVersion string = "v2"

	// FluxPath is the path for the Flux listen endpoint
	FluxPath string = "listen"
)

// GetFluxAPI builds the WebSocket URL for the Deepgram Flux (v2/listen) endpoint.
// It produces: wss://host/v2/listen?<FluxTranscriptionOptions>
//
// Overrides:
//   - host: defaults to api.deepgram.com if empty
//   - version: defaults to "v2" if empty (can be overridden via ClientOptions.APIVersion)
//   - path: defaults to "listen" if empty (can be overridden via ClientOptions.Path)
func GetFluxAPI(ctx context.Context, host, version, path string, options *interfaces.FluxTranscriptionOptions, args ...interface{}) (string, error) {
	if version == "" {
		version = FluxAPIVersion
	}
	if path == "" {
		path = FluxPath
	}
	return getAPIURL(ctx, APITypeLive, host, version, path, options, args...)
}
