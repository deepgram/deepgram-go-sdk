// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/manage"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelTrace, // LogLevelStandard / LogLevelFull / LogLevelTrace
	})

	// context
	ctx := context.Background()

	//client
	dg := client.NewWithDefaults()
	mgClient := api.New(dg)

	// list projects
	respProject, err := mgClient.ListProjects(ctx)
	if err != nil {
		fmt.Printf("ListProjects failed. Err: %v\n", err)
		os.Exit(1)
	}

	var projectID string
	for _, item := range respProject.Projects {
		projectID = item.ProjectID
		name := item.Name
		fmt.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectID)
		break
	}

	// list balances
	respList, err := mgClient.ListBalances(ctx, projectID)
	if err != nil {
		fmt.Printf("ListBalances failed. Err: %v\n", err)
		os.Exit(1)
	}

	var amount float64
	var id string
	for _, item := range respList.Balances {
		id = item.BalanceID
		amount = item.Amount
		fmt.Printf("ListBalances() - ID: %s, Amount: %f\n", id, amount)
	}

	// get first balance
	respGet, err := mgClient.GetBalance(ctx, projectID, id)
	if err != nil {
		fmt.Printf("GetBalance failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("GetBalance() - ID: %s, Amount: %f\n", id, respGet.Amount)
}
