// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Keys APIs in the Deepgram Manage API
*/
package manage

import (
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	manage "github.com/deepgram/deepgram-go-sdk/pkg/client/manage"
	rest "github.com/deepgram/deepgram-go-sdk/pkg/client/rest" //lint:ignore
)

const (
	PackageVersion string = "v1.0"
)

// Alias
type Client struct {
	*manage.Client
}

// New creates a new Client from
func New(client interface{}) *Client {
	switch client := client.(type) {
	case *rest.Client:
		return &Client{
			&manage.Client{
				RESTClient: &common.RESTClient{
					Client: client,
				},
			},
		}
	case *manage.Client:
		return &Client{client}
	}
	return nil
}
