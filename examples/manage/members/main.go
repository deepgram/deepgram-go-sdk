// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/rest"
)

const (
	TestAccountName = "enter-your-email@gmail.com"
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
	respList, err := mgClient.ListProjects(ctx)
	if err != nil {
		fmt.Printf("ListProjects failed. Err: %v\n", err)
		os.Exit(1)
	}

	var projectId string
	for _, item := range respList.Projects {
		projectId = item.ProjectID
		name := item.Name
		fmt.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectId)
		break
	}

	// list members
	respGet, err := mgClient.ListMembers(ctx, projectId)
	if err != nil {
		fmt.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	var delMemberId string
	for _, item := range respGet.Members {
		memberId := item.MemberID
		email := item.Email
		fmt.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)

		if email == TestAccountName {
			delMemberId = memberId
		}
	}

	if delMemberId == "" {
		fmt.Printf("This example requires a member who already email is %s.\n", TestAccountName)
		fmt.Printf("This is required to exercise the RemoveMember function.\n")
		fmt.Printf("In the absence of this, this example will exit early.\n")
		os.Exit(0)
	}

	// remove member
	respMessage, err := mgClient.RemoveMember(ctx, projectId, delMemberId)
	if err != nil {
		fmt.Printf("RemoveMember failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("RemoveMember() - Name: %s\n", respMessage.Message)

	// list members
	respGet, err = mgClient.ListMembers(ctx, projectId)
	if err != nil {
		fmt.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Members {
		memberId := item.MemberID
		email := item.Email
		fmt.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)
	}
}
