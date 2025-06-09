// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Token API:
https://developers.deepgram.com/reference/token-based-auth-api/grant-token
*/
package auth

import (
	"context"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
)

// GrantToken generates a JWT.
func (c *Client) GrantToken(ctx context.Context) (*api.GrantToken, error) {
	klog.V(6).Infof("auth.GrantToken() ENTER\n")

	var resp api.GrantToken
	err := c.APIRequest(ctx, "POST", version.GrantTokenURI, nil, &resp)
	if err != nil {
		klog.V(1).Infof("GrantToken failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GrantToken Succeeded\n")
	}

	klog.V(6).Infof("auth.GrantToken() LEAVE\n")
	return &resp, nil
}
