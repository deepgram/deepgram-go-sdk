// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Members API:
https://developers.deepgram.com/reference/get-members
https://developers.deepgram.com/reference/remove-member
*/
package manage

import (
	"context"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
)

// ListMembers lists all members for a project
func (c *Client) ListMembers(ctx context.Context, projectID string) (*api.MembersResult, error) {
	klog.V(6).Infof("manage.ListMembers() ENTER\n")

	var resp api.MembersResult
	err := c.APIRequest(ctx, "GET", version.MembersURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListMembers failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListMembers Succeeded\n")
	}

	klog.V(6).Infof("manage.ListMembers() LEAVE\n")
	return &resp, nil
}

// RemoveMember removes a member from a project
func (c *Client) RemoveMember(ctx context.Context, projectID, memberID string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.RemoveMember() ENTER\n")

	var resp api.MessageResult
	err := c.APIRequest(ctx, "DELETE", version.MembersByIDURI, nil, &resp, projectID, memberID)
	if err != nil {
		klog.V(1).Infof("RemoveMember failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("RemoveMember Succeeded\n")
	}

	klog.V(6).Infof("manage.RemoveMember() LEAVE\n")
	return &resp, nil
}
