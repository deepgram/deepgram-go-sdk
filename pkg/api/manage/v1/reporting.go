// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Reporting API project usage and billing endpoints.
package manage

import (
	"context"
	"net/http"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
)

// GetUsageBreakdown retrieves usage grouped by one reporting dimension.
func (c *Client) GetUsageBreakdown(ctx context.Context, projectID string, request *api.UsageBreakdownRequest) (*api.UsageBreakdownResult, error) {
	if request == nil {
		request = &api.UsageBreakdownRequest{}
	}
	var response api.UsageBreakdownResult
	if err := c.APIRequestWithQuery(ctx, http.MethodGet, version.UsageBreakdownURI, request, nil, &response, projectID); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetBillingBreakdown retrieves billing grouped by one or more reporting dimensions.
func (c *Client) GetBillingBreakdown(ctx context.Context, projectID string, request *api.BillingBreakdownRequest) (*api.BillingBreakdownResult, error) {
	if request == nil {
		request = &api.BillingBreakdownRequest{}
	}
	var response api.BillingBreakdownResult
	if err := c.APIRequestWithQuery(ctx, http.MethodGet, version.BillingBreakdownURI, request, nil, &response, projectID); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetBillingFields lists values that can filter a billing breakdown.
func (c *Client) GetBillingFields(ctx context.Context, projectID string, request *api.BillingFieldsRequest) (*api.BillingFieldsResult, error) {
	if request == nil {
		request = &api.BillingFieldsRequest{}
	}
	var response api.BillingFieldsResult
	if err := c.APIRequestWithQuery(ctx, http.MethodGet, version.BillingFieldsURI, request, nil, &response, projectID); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListPurchases lists a project's order transactions.
func (c *Client) ListPurchases(ctx context.Context, projectID string, request *api.PurchasesRequest) (*api.PurchasesResult, error) {
	if request == nil {
		request = &api.PurchasesRequest{}
	}
	var response api.PurchasesResult
	if err := c.APIRequestWithQuery(ctx, http.MethodGet, version.PurchasesURI, request, nil, &response, projectID); err != nil {
		return nil, err
	}
	return &response, nil
}
