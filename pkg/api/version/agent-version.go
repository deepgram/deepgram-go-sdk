// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"
)

const (
	// APIVersion current supported version
	AgentAPIVersion string = "v1"

	// AgentPath is the current path for agent API
	AgentPath string = "agent/converse"
)

/*
GetAgentAPI is a function which controls the versioning of the agent API and provides mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for the agent API
*/
func GetAgentAPI(ctx context.Context, host, version, path string /*options *interfaces.SettingsConfigurationOptions,*/, args ...interface{}) (string, error) {
	if version == "" {
		version = AgentAPIVersion
	}
	if path == "" {
		path = AgentPath
	}
	return getAPIURL(ctx, "agent", host, version, path, nil, args...)
}
