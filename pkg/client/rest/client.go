// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// *********** WARNING ***********
// This package provides a generic reusable REST client
//
// Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
//
// This package is frozen and no new functionality will be added.
// *********** WARNING ***********
package rest

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	restv1 "github.com/deepgram/deepgram-go-sdk/pkg/client/rest/v1"
)

const (
	PackageVersion = restv1.PackageVersion
)

// Alias
type Client = restv1.Client

// NewWithDefaults creates a REST client with default options
func NewWithDefaults() *Client {
	return New(&interfaces.ClientOptions{})
}

// New REST client
func New(options *interfaces.ClientOptions) *Client {
	return restv1.New(options)
}
