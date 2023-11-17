// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Projects APIs in the Deepgram Manage API

Please see:
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
	"net/http"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// ListProjects lists all projects for a user
func (c *ManageClient) ListProjects(ctx context.Context) (*api.ProjectsResult, error) {
	klog.V(6).Infof("manage.ListProjects() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.ProjectsURI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListProjects() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListProjects() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.ProjectsResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.ListProjects() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.ListProjects() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ListProjects Succeeded\n")
	klog.V(6).Infof("manage.ListProjects() LEAVE\n")
	return &resp, nil
}

// GetProject gets a project by ID
func (c *ManageClient) GetProject(ctx context.Context, projectId string) (*api.ProjectResult, error) {
	klog.V(6).Infof("manage.GetProject() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.ProjectsByIdURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetProject() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetProject() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.ProjectResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetProject() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetProject() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetProject Succeeded\n")
	klog.V(6).Infof("manage.GetProject() LEAVE\n")
	return &resp, nil
}

// UpdateProject updates a project
func (c *ManageClient) UpdateProject(ctx context.Context, projectId string, proj *api.ProjectUpdateRequest) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.UpdateProject() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.ProjectsByIdURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	jsonStr, err := json.Marshal(proj)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", URI, bytes.NewBuffer(jsonStr))
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("UpdateProject Succeeded\n")
	klog.V(6).Infof("manage.UpdateProject() LEAVE\n")
	return &resp, nil
}

// DeleteProject deletes a project
func (c *ManageClient) DeleteProject(ctx context.Context, projectId string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.DeleteProject() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.ProjectsByIdURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteProject() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "DELETE", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteProject() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.DeleteProject() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteProject() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("DeleteProject Succeeded\n")
	klog.V(6).Infof("manage.DeleteProject() LEAVE\n")
	return &resp, nil
}
