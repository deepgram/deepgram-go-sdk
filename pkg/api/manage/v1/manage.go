// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package manage

import (
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/rest"
)

// ManageClient is the client for the Deepgram Manage API
type ManageClient struct {
	*client.Client
}

// New creates a new ManageClient
func New(client *client.Client) *ManageClient {
	return &ManageClient{client}
}
