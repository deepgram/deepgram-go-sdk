// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Usage APIs in the Deepgram Manage API

Please see:
https://developers.deepgram.com/reference/get-all-requests
https://developers.deepgram.com/reference/get-request
https://developers.deepgram.com/reference/summarize-usage
https://developers.deepgram.com/reference/get-fields
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

// ListRequests lists all requests for a project
func (c *ManageClient) ListRequests(ctx context.Context, projectId string, use *api.UsageListRequest) (*api.UsageListResult, error) {
	klog.V(6).Infof("manage.ListRequests() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.UsageRequestURI, use, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListRequests() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListRequests() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.UsageListResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.ListRequests() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.ListRequests() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ListRequests Succeeded\n")
	klog.V(6).Infof("manage.ListRequests() LEAVE\n")
	return &resp, nil
}

// GetRequest gets a request by ID
func (c *ManageClient) GetRequest(ctx context.Context, projectId string, requestId string) (*api.UsageRequestResult, error) {
	klog.V(6).Infof("manage.GetRequest() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.UsageRequestByIdURI, nil, projectId, requestId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetRequest() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetRequest() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.UsageRequestResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetRequest() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetRequest() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetRequest Succeeded\n")
	klog.V(6).Infof("manage.GetRequest() LEAVE\n")
	return &resp, nil
}

// GetFields gets a list of fields for a project
func (c *ManageClient) GetFields(ctx context.Context, projectId string, use *api.UsageListRequest) (*api.UsageFieldResult, error) {
	klog.V(6).Infof("manage.GetFields() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.UsageFieldsURI, use, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetFields() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetFields() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.UsageFieldResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetFields() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetFields() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetFields Succeeded\n")
	klog.V(6).Infof("manage.GetFields() LEAVE\n")
	return &resp, nil
}

// GetUsage gets a usage by ID
func (c *ManageClient) GetUsage(ctx context.Context, projectId string, use *api.UsageRequest) (*api.UsageResult, error) {
	klog.V(6).Infof("manage.GetUsage() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.UsageURI, use, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetUsage() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetUsage() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.UsageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetUsage() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetUsage() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetUsage Succeeded\n")
	klog.V(6).Infof("manage.GetUsage() LEAVE\n")
	return &resp, nil
}
