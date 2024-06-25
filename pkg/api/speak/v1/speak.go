// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// *********** WARNING ***********
// This package provides the API for Speak on a REST interface
//
// Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
//
// This package is frozen and no new functionality will be added.
// *********** WARNING ***********
package legacy

import (
	speakv1 "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/rest"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

const (
	PackageVersion = speakv1.PackageVersion
)

// Alias
type Client = speakv1.Client

// New creates a new Client
func New(c *client.Client) *Client {
	return speakv1.New(c)
}
