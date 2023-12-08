// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
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
	"net/http"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// GetMemberScopes gets the scopes for a member
func (c *ManageClient) GetMemberScopes(ctx context.Context, projectId string, memberId string) (*api.ScopeResult, error) {
	klog.V(6).Infof("manage.GetMemberScopes() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.MembersScopeByIdURI, nil, projectId, memberId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetMemberScopes() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetMemberScopes() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.ScopeResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetMemberScopes() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetMemberScopes() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetMemberScopes Succeeded\n")
	klog.V(6).Infof("manage.GetMemberScopes() LEAVE\n")
	return &resp, nil
}

// UpdateMemberScopes updates the scopes for a member
func (c *ManageClient) UpdateMemberScopes(ctx context.Context, projectId string, memberId string, scope *api.ScopeUpdateRequest) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.UpdateMemberScopes() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.MembersScopeByIdURI, nil, projectId, memberId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	jsonStr, err := json.Marshal(scope)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", URI, bytes.NewBuffer(jsonStr))
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("UpdateMemberScopes Succeeded\n")
	klog.V(6).Infof("manage.UpdateMemberScopes() LEAVE\n")
	return &resp, nil
}
