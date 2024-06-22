// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Keys API:
https://developers.deepgram.com/reference/list-keys
https://developers.deepgram.com/reference/get-key
https://developers.deepgram.com/reference/create-key
https://developers.deepgram.com/reference/delete-key
*/
package manage

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
)

// ListKeys lists all keys for a project
func (c *Client) ListKeys(ctx context.Context, projectID string) (*api.KeysResult, error) {
	klog.V(6).Infof("manage.ListKeys() ENTER\n")

	var resp api.KeysResult
	err := c.APIRequest(ctx, "GET", version.KeysURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListKeys failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListKeys Succeeded\n")
	}

	klog.V(6).Infof("manage.ListKeys() LEAVE\n")
	return &resp, nil
}

// GetKey gets a key for a project
func (c *Client) GetKey(ctx context.Context, projectID, keyID string) (*api.KeyResult, error) {
	klog.V(6).Infof("manage.GetKey() ENTER\n")

	var resp api.KeyResult
	err := c.APIRequest(ctx, "GET", version.KeysByIDURI, nil, &resp, projectID, keyID)
	if err != nil {
		klog.V(1).Infof("GetKey failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetKey Succeeded\n")
	}

	klog.V(6).Infof("manage.GetKey() LEAVE\n")
	return &resp, nil
}

// CreateKey creates a key for a project
func (c *Client) CreateKey(ctx context.Context, projectID string, key *api.KeyCreateRequest) (*api.APIKey, error) {
	klog.V(6).Infof("manage.CreateKey() ENTER\n")

	var expirationStr string
	if !key.ExpirationDate.IsZero() {
		expirationStr = key.ExpirationDate.Format(time.RFC3339)
	}

	type InternalKeyCreateRequest struct {
		Comment        string   `json:"comment"`
		Scopes         []string `json:"scopes"`
		ExpirationDate string   `json:"expiration_date,omitempty"`
		TimeToLive     int      `json:"time_to_live_in_seconds,omitempty"`
		// Tags           []string `json:"tags"`
	}
	internalKey := InternalKeyCreateRequest{
		Comment:        key.Comment,
		Scopes:         key.Scopes,
		ExpirationDate: expirationStr,
		TimeToLive:     key.TimeToLive,
		// Tags:           key.Tags,
	}

	jsonStr, err := json.Marshal(internalKey)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.CreateKey() LEAVE\n")
		return nil, err
	}

	var resp api.APIKey
	err = c.APIRequest(ctx, "POST", version.KeysURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("CreateKey failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("CreateKey Succeeded\n")
	}

	klog.V(6).Infof("manage.CreateKey() LEAVE\n")
	return &resp, nil
}

// DeleteKey deletes a key for a project
func (c *Client) DeleteKey(ctx context.Context, projectID, keyID string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.DeleteKey() ENTER\n")

	var resp api.MessageResult
	err := c.APIRequest(ctx, "DELETE", version.KeysByIDURI, nil, &resp, projectID, keyID)
	if err != nil {
		klog.V(1).Infof("DeleteKey failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("DeleteKey Succeeded\n")
	}

	klog.V(6).Infof("manage.DeleteKey() LEAVE\n")
	return &resp, nil
}
