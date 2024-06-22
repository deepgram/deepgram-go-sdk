// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Projects API:
https://developers.deepgram.com/reference/get-projects
https://developers.deepgram.com/reference/get-project
https://developers.deepgram.com/reference/update-project
https://developers.deepgram.com/reference/delete-project
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

// ListProjects lists all projects for a user
func (c *Client) ListProjects(ctx context.Context) (*api.ProjectsResult, error) {
	klog.V(6).Infof("manage.ListProjects() ENTER\n")

	var resp api.ProjectsResult
	err := c.APIRequest(ctx, "GET", version.ProjectsURI, nil, &resp)
	if err != nil {
		klog.V(1).Infof("ListProjects failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListProjects Succeeded\n")
	}

	klog.V(6).Infof("manage.ListProjects() LEAVE\n")
	return &resp, nil
}

// GetProject gets a project by ID
func (c *Client) GetProject(ctx context.Context, projectID string) (*api.ProjectResult, error) {
	klog.V(6).Infof("manage.GetProject() ENTER\n")

	var resp api.ProjectResult
	err := c.APIRequest(ctx, "GET", version.ProjectsByIDURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("GetProject failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetProject Succeeded\n")
	}

	klog.V(6).Infof("manage.GetProject() LEAVE\n")
	return &resp, nil
}

// UpdateProject updates a project
func (c *Client) UpdateProject(ctx context.Context, projectID string, proj *api.ProjectUpdateRequest) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.UpdateProject() ENTER\n")

	jsonStr, err := json.Marshal(proj)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		return nil, err
	}

	var resp api.MessageResult
	err = c.APIRequest(ctx, "PATCH", version.ProjectsByIDURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("UpdateProject failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("UpdateProject Succeeded\n")
	}

	klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
	return &resp, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(ctx context.Context, projectID string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.DeleteProject() ENTER\n")

	var resp api.MessageResult
	err := c.APIRequest(ctx, "DELETE", version.InvitationsByIDURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("DeleteProject failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("DeleteProject Succeeded\n")
	}

	klog.V(6).Infof("manage.DeleteProject() LEAVE\n")
	return &resp, nil
}
