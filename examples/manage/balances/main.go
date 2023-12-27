// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
)

func main() {
	// init library
	client.InitWithDefault()

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

	var projectId string
	for _, item := range respProject.Projects {
		projectId = item.ProjectID
		name := item.Name
		fmt.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectId)
		break
	}

	// list balances
	respList, err := mgClient.ListBalances(ctx, projectId)
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
	respGet, err := mgClient.GetBalance(ctx, projectId, id)
	if err != nil {
		fmt.Printf("GetBalance failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("GetBalance() - ID: %s, Amount: %f\n", id, respGet.Amount)
}
