// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Reporting retrieves grouped usage and billing data for the first project the
// configured API key can access.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

func main() {
	ctx := context.Background()
	management := api.New(client.NewWithDefaults())

	projectID := os.Getenv("DEEPGRAM_PROJECT_ID")
	if projectID == "" {
		log.Fatal("DEEPGRAM_PROJECT_ID is required")
	}
	usage, err := management.GetUsageBreakdown(ctx, projectID, &interfaces.UsageBreakdownRequest{
		Grouping: interfaces.UsageBreakdownGroupingEndpoint,
	})
	if err != nil {
		log.Fatal(err)
	}

	billing, err := management.GetBillingBreakdown(ctx, projectID, &interfaces.BillingBreakdownRequest{
		Grouping: []interfaces.BillingBreakdownGrouping{
			interfaces.BillingBreakdownGroupingLineItem,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	result := struct {
		Usage   interface{} `json:"usage"`
		Billing interface{} `json:"billing"`
	}{Usage: usage, Billing: billing}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
