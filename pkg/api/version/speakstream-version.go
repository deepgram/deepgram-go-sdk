// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

const (
	// SpeakStreamAPIVersion current supported version
	SpeakStreamAPIVersion string = "v1"

	// SpeakStreamPath is the current path
	SpeakStreamPath string = "speak-stream"
)

/*
GetSpeakStreamAPI is a function which controls the versioning of the text-to-speech API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for the text-to-speech
*/
func GetSpeakStreamAPI(ctx context.Context, host, version, path string, options *interfaces.WSSpeakOptions, args ...interface{}) (string, error) {
	return getAPIURL(ctx, "speak-stream", host, version, path, options, args...)
}
