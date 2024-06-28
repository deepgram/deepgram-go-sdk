// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/manage"
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

	var projectID string
	for _, item := range respList.Projects {
		projectID = item.ProjectID
		name := item.Name
		fmt.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectID)
		break
	}

	// list members
	respGet, err := mgClient.ListMembers(ctx, projectID)
	if err != nil {
		fmt.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	var updateID string
	for _, item := range respGet.Members {
		memberID := item.MemberID
		email := item.Email
		fmt.Printf("ListMembers() - ID: %s, Scope: %s\n", memberID, email)
		for _, scope := range item.Scopes {
			fmt.Printf("Scope: %s\n", scope)
		}

		if email == TestAccountName {
			updateID = memberID
		}
	}

	if updateID == "" {
		fmt.Printf("This example requires a member who already email is %s.\n", TestAccountName)
		fmt.Printf("This is required to exercise the RemoveMember function.\n")
		fmt.Printf("In the absence of this, this example will exit early.\n")
		os.Exit(0)
	}

	// get scope
	respScope, err := mgClient.GetMemberScopes(ctx, projectID, updateID)
	if err != nil {
		fmt.Printf("GetMemberScopes failed. Err: %v\n", err)
		os.Exit(1)
	}

	scope := respScope.Scopes
	for _, item := range scope {
		fmt.Printf("GetMemberScopes() - Scope: %s\n", item)
	}

	// update scope
	respMessage, err := mgClient.UpdateMemberScopes(ctx, projectID, updateID, &interfaces.ScopeUpdateRequest{
		Scope: "admin",
	})
	if err != nil {
		fmt.Printf("UpdateMemberScopes failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("UpdateMemberScopes() - Name: %s\n", respMessage.Message)

	// list members
	respGet, err = mgClient.ListMembers(ctx, projectID)
	if err != nil {
		fmt.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Members {
		memberID := item.MemberID
		email := item.Email
		fmt.Printf("ListMembers() - ID: %s, Scope: %s\n", memberID, email)
		for _, scope := range item.Scopes {
			fmt.Printf("Scope: %s\n", scope)
		}
	}

	// update scope
	respMessage, err = mgClient.UpdateMemberScopes(ctx, projectID, updateID, &interfaces.ScopeUpdateRequest{
		Scope: "member",
	})
	if err != nil {
		fmt.Printf("UpdateMemberScopes failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("UpdateMemberScopes() - Name: %s\n", respMessage.Message)

	// list members
	respGet, err = mgClient.ListMembers(ctx, projectID)
	if err != nil {
		fmt.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Members {
		memberID := item.MemberID
		email := item.Email
		fmt.Printf("ListMembers() - ID: %s, Scope: %s\n", memberID, email)
		for _, scope := range item.Scopes {
			fmt.Printf("Scope: %s\n", scope)
		}
	}
}
