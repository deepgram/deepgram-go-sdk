// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Balances API:
https://developers.deepgram.com/reference/get-all-balances
https://developers.deepgram.com/reference/get-balance
*/
package manage

import (
	"context"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
)

// ListBalances lists all balances for a project
func (c *Client) ListBalances(ctx context.Context, projectID string) (*api.BalancesResult, error) {
	klog.V(6).Infof("manage.ListBalances() ENTER\n")

	var resp api.BalancesResult
	err := c.APIRequest(ctx, "GET", version.BalancesURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListBalances failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListBalances Succeeded\n")
	}

	klog.V(6).Infof("manage.ListBalances() LEAVE\n")
	return &resp, nil
}

// GetBalance gets a balance for a project
func (c *Client) GetBalance(ctx context.Context, projectID, balanceID string) (*api.BalanceResult, error) {
	klog.V(6).Infof("manage.GetBalance() ENTER\n")

	var resp api.BalanceResult
	err := c.APIRequest(ctx, "GET", version.BalancesByIDURI, nil, &resp, projectID, balanceID)
	if err != nil {
		klog.V(1).Infof("GetBalance failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetBalance Succeeded\n")
	}

	klog.V(6).Infof("manage.GetBalance() LEAVE\n")
	return &resp, nil
}
