// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the code for the Balances APIs in the Deepgram Manage API

Please see:
https://developers.deepgram.com/reference/get-all-balances
https://developers.deepgram.com/reference/get-balance
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

// ListBalances lists all balances for a project
func (c *ManageClient) ListBalances(ctx context.Context, projectId string) (*api.BalancesResult, error) {
	klog.V(6).Infof("manage.ListBalances() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.BalancesURI, nil, projectId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListBalances() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.ListBalances() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.BalancesResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.ListBalances() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.ListBalances() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ListBalances Succeeded\n")
	klog.V(6).Infof("manage.ListBalances() LEAVE\n")
	return &resp, nil
}

// GetBalance gets a balance for a project
func (c *ManageClient) GetBalance(ctx context.Context, projectId string, balanceId string) (*api.BalanceResult, error) {
	klog.V(6).Infof("manage.GetBalance() ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.BalancesByIdURI, nil, projectId, balanceId)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetBalance() LEAVE\n")
		return nil, err
	}
	klog.V(4).Infof("Calling %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("manage.GetBalance() LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.BalanceResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("manage.GetBalance() LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("manage.GetBalance() LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("GetBalance Succeeded\n")
	klog.V(6).Infof("manage.GetBalance() LEAVE\n")
	return &resp, nil
}
