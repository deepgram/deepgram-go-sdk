// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Models API:
https://developers.deepgram.com/docs/model-metadata
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

// ListModels lists all models available
// NOTE: This is a wrapper around GetModels
//
// Args:
//
//	ctx: context
//	model: model request options
//
// Returns:
//
//	*api.ModelsResult: list of models
func (c *Client) ListModels(ctx context.Context, model *api.ModelRequest) (*api.ModelsResult, error) {
	return c.GetModels(ctx, model)
}

// GetModels lists all models available
//
// Args:
//
//	ctx: context
//	model: model request options
//
// Returns:
//
//	*api.ModelsResult: list of models
func (c *Client) GetModels(ctx context.Context, model *api.ModelRequest) (*api.ModelsResult, error) {
	klog.V(6).Infof("manage.GetModels() ENTER\n")

	if model == nil {
		model = &api.ModelRequest{}
	}

	jsonStr, err := json.Marshal(model)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetModels() LEAVE\n")
		return nil, err
	}

	var resp api.ModelsResult
	err = c.APIRequest(ctx, "GET", version.ModelsURI, bytes.NewBuffer(jsonStr), &resp)
	if err != nil {
		klog.V(1).Infof("GetModels failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetModels Succeeded\n")
	}

	klog.V(6).Infof("manage.GetModels() LEAVE\n")
	return &resp, nil
}

// GetModel gets a model by ID
//
// Args:
//
//	ctx: context
//	modelID: model ID
//
// Returns:
//
//	*api.ModelResult: specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*api.ModelResult, error) {
	klog.V(6).Infof("manage.GetModel() ENTER\n")

	var resp api.ModelResult
	err := c.APIRequest(ctx, "GET", version.ModelsByIDURI, nil, &resp, modelID)
	if err != nil {
		klog.V(1).Infof("GetModel failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetModel Succeeded\n")
	}

	klog.V(6).Infof("manage.GetModel() LEAVE\n")
	return &resp, nil
}

// ListProjectModels lists all models available
// NOTE: This is a wrapper around GetProjectModels
//
// Args:
//
//	ctx: context
//	projectID: project ID
//	model: model request options
//
// Returns:
//
//	*api.ModelsResult: list of models
func (c *Client) ListProjectModels(ctx context.Context, projectID string, model *api.ModelRequest) (*api.ModelsResult, error) {
	return c.GetProjectModels(ctx, projectID, model)
}

// GetProjectModels lists all models available
//
// Args:
//
//	ctx: context
//	projectID: project ID
//	model: model request options
//
// Returns:
//
//	*api.ModelsResult: list of models
func (c *Client) GetProjectModels(ctx context.Context, projectID string, model *api.ModelRequest) (*api.ModelsResult, error) {
	klog.V(6).Infof("manage.GetProjectModels() ENTER\n")

	if model == nil {
		model = &api.ModelRequest{}
	}

	jsonStr, err := json.Marshal(model)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetProjectModels() LEAVE\n")
		return nil, err
	}

	var resp api.ModelsResult
	err = c.APIRequest(ctx, "GET", version.ModelsProjectURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("GetProjectModels failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetProjectModels Succeeded\n")
	}

	klog.V(6).Infof("manage.GetProjectModels() LEAVE\n")
	return &resp, nil
}

// GetProjectModel gets a single model within the project by ID
//
// Args:
//
//	ctx: context
//	projectID: project ID
//	modelID: model ID
//
// Returns:
//
//	*api.ModelResult: specific model
func (c *Client) GetProjectModel(ctx context.Context, projectID, modelID string) (*api.ModelResult, error) {
	klog.V(6).Infof("manage.GetProjectModel() ENTER\n")

	var resp api.ModelResult
	err := c.APIRequest(ctx, "GET", version.ModelsProjectByIDURI, nil, &resp, projectID, modelID)
	if err != nil {
		klog.V(1).Infof("GetProjectModel failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetProjectModel Succeeded\n")
	}

	klog.V(6).Infof("manage.GetProjectModel() LEAVE\n")
	return &resp, nil
}
