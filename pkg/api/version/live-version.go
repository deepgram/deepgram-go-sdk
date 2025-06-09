// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

const (
	// LiveAPIVersion current supported version
	LiveAPIVersion string = "v1"

	// LivePath is the current path for live transcription
	LivePath string = "listen"
)

/*
GetLiveAPI is a function which controls the versioning of the live transcription API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for the live transcription
*/
func GetLiveAPI(ctx context.Context, host, version, path string, options *interfaces.LiveTranscriptionOptions, args ...interface{}) (string, error) {
	return getAPIURL(ctx, "live", host, version, path, options, args...)
}
