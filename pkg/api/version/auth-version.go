// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"
)

const (
	// version
	AuthAPIVersion string = "v1"

	// grant token
	GrantTokenURI string = "auth/grant"
)

/*
GetAuthAPI is a function which controls the versioning of the auth API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for auth
*/
func GetAuthAPI(ctx context.Context, host, version, path string, vals interface{}, args ...interface{}) (string, error) {
	return getAPIURL(ctx, "auth", host, version, path, nil, args...)
}
