// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Scopes API:
https://developers.deepgram.com/reference/get-member-scopes
https://developers.deepgram.com/reference/update-scope
*/
package manage

import (
	"bytes"
	"context"
	"encoding/json"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
)

// GetMemberScopes gets the scopes for a member
func (c *Client) GetMemberScopes(ctx context.Context, projectID, memberID string) (*api.ScopeResult, error) {
	klog.V(6).Infof("manage.GetMemberScopes() ENTER\n")

	var resp api.ScopeResult
	err := c.APIRequest(ctx, "GET", version.MembersScopeByIDURI, nil, &resp, projectID, memberID)
	if err != nil {
		klog.V(1).Infof("GetMemberScopes failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetMemberScopes Succeeded\n")
	}

	klog.V(6).Infof("manage.GetMemberScopes() LEAVE\n")
	return &resp, nil
}

// UpdateMemberScopes updates the scopes for a member
func (c *Client) UpdateMemberScopes(ctx context.Context, projectID, memberID string, scope *api.ScopeUpdateRequest) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.UpdateMemberScopes() ENTER\n")

	jsonStr, err := json.Marshal(scope)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		return nil, err
	}

	var resp api.MessageResult
	err = c.APIRequest(ctx, "PUT", version.MembersScopeByIDURI, bytes.NewBuffer(jsonStr), &resp, projectID, memberID)
	if err != nil {
		klog.V(1).Infof("UpdateMemberScopes failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("UpdateMemberScopes Succeeded\n")
	}

	klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
	return &resp, nil
}
