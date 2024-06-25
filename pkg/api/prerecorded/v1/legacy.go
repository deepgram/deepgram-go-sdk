// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// *********** WARNING ***********
// This package provides the PreRecorded API
//
// Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
//
// This package is frozen and no new functionality will be added.
// *********** WARNING ***********
package legacy

import (
	restv1 "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded" //lint:ignore
)

const (
	PackageVersion = restv1.PackageVersion
)

// Alias
type Client = restv1.Client

// New creates a new Client
func New(c *client.Client) *restv1.Client {
	return restv1.New(c)
}
