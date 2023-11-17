// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Keys APIs in the Deepgram Manage API

Please see:
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
	"net/http"
	"time"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// ListKeys lists all keys for a project
func (c *ManageClient) ListKeys(ctx context.Context, projectId string) (*api.KeysResult, error) {
	klog.V(6).Infof("manage.ListKeys() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.KeysURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListKeys() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListKeys() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.KeysResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.ListKeys() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.ListKeys() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ListKeys Succeeded\n")
	klog.V(6).Infof("manage.ListKeys() LEAVE\n")
	return &resp, nil
}

// GetKey gets a key for a project
func (c *ManageClient) GetKey(ctx context.Context, projectId string, keyId string) (*api.KeyResult, error) {
	klog.V(6).Infof("manage.GetKey() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.KeysByIdURI, nil, projectId, keyId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetKey() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetKey() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.KeyResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetKey() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetKey() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetKey Succeeded\n")
	klog.V(6).Infof("manage.GetKey() LEAVE\n")
	return &resp, nil
}

// CreateKey creates a key for a project
func (c *ManageClient) CreateKey(ctx context.Context, projectId string, key *api.KeyCreateRequest) (*api.APIKey, error) {
	klog.V(6).Infof("manage.CreateKey() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	var expirationStr string
	if !key.ExpirationDate.IsZero() {
		expirationStr = key.ExpirationDate.Format(time.RFC3339)
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.KeysURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.CreateKey() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	type InternalKeyCreateRequest struct {
		Comment        string   `json:"comment"`
		Scopes         []string `json:"scopes"`
		ExpirationDate string   `json:"expiration_date,omitempty"`
		TimeToLive     int      `json:"time_to_live,omitempty"`
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

	req, err := http.NewRequestWithContext(ctx, "POST", URI, bytes.NewBuffer(jsonStr))
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.CreateKey() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.APIKey
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.CreateKey() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.CreateKey() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("CreateKey Succeeded\n")
	klog.V(6).Infof("manage.CreateKey() LEAVE\n")
	return &resp, nil
}

// DeleteKey deletes a key for a project
func (c *ManageClient) DeleteKey(ctx context.Context, projectId string, keyId string) (*api.MessageResult, error) {
	klog.V(6).Infof("manage.DeleteKey() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.KeysByIdURI, nil, projectId, keyId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteKey() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "DELETE", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteKey() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.DeleteKey() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.DeleteKey() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("DeleteKey Succeeded\n")
	klog.V(6).Infof("manage.DeleteKey() LEAVE\n")
	return &resp, nil
}
