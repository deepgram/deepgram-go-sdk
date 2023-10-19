// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package manage

import (
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

type ManageClient struct {
	*client.Client
}

func New(client *client.Client) *ManageClient {
	return &ManageClient{client}
}
