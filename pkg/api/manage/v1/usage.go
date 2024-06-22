// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Usage API:
https://developers.deepgram.com/reference/get-all-requests
https://developers.deepgram.com/reference/get-request
https://developers.deepgram.com/reference/summarize-usage
https://developers.deepgram.com/reference/get-fields
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

// ListRequests lists all requests for a project
func (c *Client) ListRequests(ctx context.Context, projectID string, use *api.UsageListRequest) (*api.UsageListResult, error) {
	klog.V(6).Infof("manage.ListRequests() ENTER\n")

	jsonStr, err := json.Marshal(use)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		return nil, err
	}

	var resp api.UsageListResult
	err = c.APIRequest(ctx, "GET", version.UsageRequestURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListRequests failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListRequests Succeeded\n")
	}

	klog.V(6).Infof("manage.ListRequests() LEAVE\n")
	return &resp, nil
}

// GetRequest gets a request by ID
func (c *Client) GetRequest(ctx context.Context, projectID, requestID string) (*api.UsageRequestResult, error) {
	klog.V(6).Infof("manage.GetRequest() ENTER\n")

	var resp api.UsageRequestResult
	err := c.APIRequest(ctx, "GET", version.UsageRequestByIDURI, nil, &resp, projectID, requestID)
	if err != nil {
		klog.V(1).Infof("GetRequest failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetRequest Succeeded\n")
	}

	klog.V(6).Infof("manage.GetRequest() LEAVE\n")
	return &resp, nil
}

// GetFields gets a list of fields for a project
func (c *Client) GetFields(ctx context.Context, projectID string, use *api.UsageListRequest) (*api.UsageFieldResult, error) {
	klog.V(6).Infof("manage.GetFields() ENTER\n")

	jsonStr, err := json.Marshal(use)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		return nil, err
	}

	var resp api.UsageFieldResult
	err = c.APIRequest(ctx, "GET", version.UsageFieldsURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("GetFields failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetFields Succeeded\n")
	}

	klog.V(6).Infof("manage.GetFields() LEAVE\n")
	return &resp, nil
}

// GetUsage gets a usage by ID
func (c *Client) GetUsage(ctx context.Context, projectID string, use *api.UsageRequest) (*api.UsageResult, error) {
	klog.V(6).Infof("manage.GetUsage() ENTER\n")

	jsonStr, err := json.Marshal(use)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		return nil, err
	}

	var resp api.UsageResult
	err = c.APIRequest(ctx, "GET", version.UsageURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("GetUsage failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetUsage Succeeded\n")
	}

	klog.V(6).Infof("manage.GetUsage() LEAVE\n")
	return &resp, nil
}
