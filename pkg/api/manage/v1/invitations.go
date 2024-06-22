// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Invitations API
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

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
)

// ListInvitations lists all invitations for a project
func (c *Client) ListInvitations(ctx context.Context, projectID string) (*api.InvitationsResult, error) {
	klog.V(6).Infof("manage.ListInvitations() ENTER\n")

	var resp api.InvitationsResult
	err := c.APIRequest(ctx, "GET", version.InvitationsURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListInvitations failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListInvitations Succeeded\n")
	}

	klog.V(6).Infof("manage.ListInvitations() LEAVE\n")
	return &resp, err
}

// SendInvitation sends an invitation to a project
func (c *Client) SendInvitation(ctx context.Context, projectID string, invite *api.InvitationRequest) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.SendInvitation() ENTER\n")

	jsonStr, err := json.Marshal(invite)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		return nil, err
	}

	var resp api.MessageResult
	err = c.APIRequest(ctx, "POST", version.InvitationsURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("SendInvitation failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("SendInvitation Succeeded\n")
	}

	klog.V(6).Infof("manage.SendInvitation() ENTER\n")
	return &resp, err
}

// DeleteInvitation deletes an invitation to a project
func (c *Client) DeleteInvitation(ctx context.Context, projectID, email string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.DeleteInvitation() ENTER\n")

	var resp api.MessageResult
	err := c.APIRequest(ctx, "DELETE", version.InvitationsByIDURI, nil, &resp, projectID, email)
	if err != nil {
		klog.V(1).Infof("DeleteInvitation failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("DeleteInvitation Succeeded\n")
	}

	klog.V(6).Infof("manage.DeleteInvitation() LEAVE\n")
	return &resp, err
}

// LeaveProject leaves a project
func (c *Client) LeaveProject(ctx context.Context, projectID string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.LeaveProject() ENTER\n")

	var resp api.MessageResult
	err := c.APIRequest(ctx, "DELETE", version.InvitationsLeaveURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("LeaveProject failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("LeaveProject Succeeded\n")
	}

	klog.V(6).Infof("manage.LeaveProject() LEAVE\n")
	return &resp, err
}
