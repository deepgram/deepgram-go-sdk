// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// *********** WARNING ***********
// This package provides the prerecorded client implementation for the Deepgram API
//
// Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
//
// This package is frozen and no new functionality will be added.
// *********** WARNING ***********
package prerecorded

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	listenv1rest "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"
)

/***********************************/
// PreRecordedClient
/***********************************/
const (
	PackageVersion = listenv1rest.PackageVersion
)

type Client = listenv1rest.Client

/*
NewWithDefaults creates a new analyze/read client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
func NewWithDefaults() *Client {
	return listenv1rest.NewWithDefaults()
}

/*
New creates a new prerecorded client with the specified options

Input parameters:
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.

Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
*/
func New(apiKey string, options *interfaces.ClientOptions) *Client {
	return listenv1rest.New(apiKey, options)
}
