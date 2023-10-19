// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"log"
	"os"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

func main() {
	// context
	ctx := context.Background()

	//client
	dg := client.NewWithDefaults()
	mgClient := api.New(dg)

	// list projects
	respProject, err := mgClient.ListProjects(ctx)
	if err != nil {
		log.Printf("ListProjects failed. Err: %v\n", err)
		os.Exit(1)
	}

	var projectId string
	for _, item := range respProject.Projects {
		projectId = item.ProjectId
		name := item.Name
		log.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectId)
		break
	}

	// list balances
	respList, err := mgClient.ListBalances(ctx, projectId)
	if err != nil {
		log.Printf("ListBalances failed. Err: %v\n", err)
		os.Exit(1)
	}

	var amount float64
	var id string
	for _, item := range respList.Balances {
		id = item.BalanceId
		amount = item.Amount
		log.Printf("ListBalances() - ID: %s, Amount: %f\n", id, amount)
	}

	// get first balance
	respGet, err := mgClient.GetBalance(ctx, projectId, id)
	if err != nil {
		log.Printf("GetBalance failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetBalance() - ID: %s, Amount: %f\n", id, respGet.Amount)
}
