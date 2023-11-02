// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Members APIs in the Deepgram Manage API

Please see:
https://developers.deepgram.com/reference/get-members
https://developers.deepgram.com/reference/remove-member
*/
package manage

import (
	"context"
	"net/http"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
)

// ListMembers lists all members for a project
func (c *ManageClient) ListMembers(ctx context.Context, projectId string) (*api.MembersResult, error) {
	klog.V(6).Infof("manage.ListMembers() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.MembersURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListMembers() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListMembers() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MembersResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.ListMembers() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.ListMembers() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ListMembers Succeeded\n")
	klog.V(6).Infof("manage.ListMembers() LEAVE\n")
	return &resp, nil
}

// RemoveMember removes a member from a project
func (c *ManageClient) RemoveMember(ctx context.Context, projectId string, memberId string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.RemoveMember() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.MembersByIdURI, nil, projectId, memberId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.RemoveMember() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "DELETE", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.RemoveMember() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.RemoveMember() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.RemoveMember() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("RemoveMember Succeeded\n")
	klog.V(6).Infof("manage.RemoveMember() LEAVE\n")
	return &resp, nil
}
