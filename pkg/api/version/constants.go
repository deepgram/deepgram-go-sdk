// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package handles the versioning in the API for the various clients (prerecorded, live, etc.)
package version

import "errors"

const (
	// APIProtocol default protocol
	HTTPProtocol string = "https"

	// WSProtocol default protocol
	WSProtocol string = "wss"

	// APIPathListen
	APIPathListen string = "listen"

	// APITypeAgent
	APITypeAgent string = "agent"

	// APITypeLive
	APITypeLive string = "live"

	// APITypeSpeakStream
	APITypeSpeakStream string = "speak-stream"
)

const (
	// DefaultAPIVersion is the current supported default version for APIs
	DefaultAPIVersion string = "v1"
)

var (
	// ErrInvalidPath invalid path
	ErrInvalidPath = errors.New("invalid path")

	// APIPathMap maps the API types to their default paths
	APIPathMap = map[string]string{
		"agent":        "agent",
		"analyze":      "read",
		"prerecorded":  APIPathListen,
		"speak":        "speak",
		"speak-stream": "speak",
		"manage":       "",
		"live":         APIPathListen,
	}
)
