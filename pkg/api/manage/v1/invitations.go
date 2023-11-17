// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Invitations APIs in the Deepgram Manage API

Please see:
https://developers.deepgram.com/reference/list-invites
https://developers.deepgram.com/reference/send-invites
https://developers.deepgram.com/reference/delete-invite
https://developers.deepgram.com/reference/leave-project
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

// ListInvitations lists all invitations for a project
func (c *ManageClient) ListInvitations(ctx context.Context, projectId string) (*api.InvitationsResult, error) {
	klog.V(6).Infof("manage.ListInvitations() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.InvitationsURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListInvitations() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListInvitations() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.InvitationsResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.ListInvitations() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.ListInvitations() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ListInvitations Succeeded\n")
	klog.V(6).Infof("manage.ListInvitations() LEAVE\n")
	return &resp, nil
}

// SendInvitation sends an invitation to a project
func (c *ManageClient) SendInvitation(ctx context.Context, projectId string, invite *api.InvitationRequest) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.SendInvitation() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.InvitationsURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.SendInvitation() ENTER\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	jsonStr, err := json.Marshal(invite)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.SendInvitation() ENTER\n")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", URI, bytes.NewBuffer(jsonStr))
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.SendInvitation() ENTER\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.SendInvitation() ENTER\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.SendInvitation() ENTER\n")
		return nil, err
	}

	klog.V(3).Infof("SendInvitation Succeeded\n")
	klog.V(6).Infof("manage.SendInvitation() ENTER\n")
	return &resp, nil
}

// DeleteInvitation deletes an invitation to a project
func (c *ManageClient) DeleteInvitation(ctx context.Context, projectId string, email string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.DeleteInvitation() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.InvitationsByIdURI, nil, projectId, email)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteInvitation() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "DELETE", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteInvitation() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.DeleteInvitation() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteInvitation() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("DeleteInvitation Succeeded\n")
	klog.V(6).Infof("manage.DeleteInvitation() LEAVE\n")
	return &resp, nil
}

// LeaveProject leaves a project
func (c *ManageClient) LeaveProject(ctx context.Context, projectId string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.LeaveProject() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.InvitationsLeaveURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.LeaveProject() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "DELETE", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.LeaveProject() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.LeaveProject() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.LeaveProject() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("LeaveProject Succeeded\n")
	klog.V(6).Infof("manage.LeaveProject() LEAVE\n")
	return &resp, nil
}
