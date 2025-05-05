// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Token APIs in the Deepgram auth API
*/
package auth

import (
	auth "github.com/deepgram/deepgram-go-sdk/v2/pkg/client/auth"
	common "github.com/deepgram/deepgram-go-sdk/v2/pkg/client/common/v1"
	rest "github.com/deepgram/deepgram-go-sdk/v2/pkg/client/rest" //lint:ignore
)

const (
	PackageVersion string = "v1.0"
)

// Alias
type Client struct {
	*auth.Client
}

// New creates a new Client from
func New(client interface{}) *Client {
	switch client := client.(type) {
	case *rest.Client:
		return &Client{
			&auth.Client{
				RESTClient: &common.RESTClient{
					Client: client,
				},
			},
		}
	case *auth.Client:
		return &Client{client}
	}
	return nil
}
