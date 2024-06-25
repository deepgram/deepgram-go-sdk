// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package contains the code for the Keys APIs in the Deepgram Manage API
package manage

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	managev1 "github.com/deepgram/deepgram-go-sdk/pkg/client/manage/v1"
)

const (
	PackageVersion = managev1.PackageVersion
)

// Alias
type Client = managev1.Client

/*
NewWithDefaults creates a new analyze/read client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWithDefaults() *Client {
	return managev1.NewWithDefaults()
}

/*
New creates a new analyze/read client with the specified options

Input parameters:
- [Optional] apiKey: string containing the Deepgram API key. If empty, the API key is read from the environment variable DEEPGRAM_API_KEY
- [Optional] options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func New(apiKey string, options *interfaces.ClientOptions) *Client {
	return managev1.New(apiKey, options)
}
