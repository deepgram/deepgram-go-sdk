// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"log"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
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
		log.Printf("ListProjects failed. Err: %v\n", err)
		os.Exit(1)
	}

	var projectId string
	for _, item := range respList.Projects {
		projectId = item.ProjectID
		name := item.Name
		log.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectId)
		break
	}

	// list members
	respGet, err := mgClient.ListMembers(ctx, projectId)
	if err != nil {
		log.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	var updateID string
	for _, item := range respGet.Members {
		memberId := item.MemberID
		email := item.Email
		log.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)
		for _, scope := range item.Scopes {
			log.Printf("Scope: %s\n", scope)
		}

		if email == TestAccountName {
			updateID = memberId
		}
	}

	if updateID == "" {
		log.Printf("This example requires a member who already email is %s.\n", TestAccountName)
		log.Printf("This is required to exercise the RemoveMember function.\n")
		log.Printf("In the absence of this, this example will exit early.\n")
		os.Exit(0)
	}

	// get scope
	respScope, err := mgClient.GetMemberScopes(ctx, projectId, updateID)
	if err != nil {
		log.Printf("GetMemberScopes failed. Err: %v\n", err)
		os.Exit(1)
	}

	scope := respScope.Scopes
	for _, item := range scope {
		log.Printf("GetMemberScopes() - Scope: %s\n", item)
	}

	// update scope
	respMessage, err := mgClient.UpdateMemberScopes(ctx, projectId, updateID, &interfaces.ScopeUpdateRequest{
		Scope: "admin",
	})
	if err != nil {
		log.Printf("UpdateMemberScopes failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("UpdateMemberScopes() - Name: %s\n", respMessage.Message)

	// list members
	respGet, err = mgClient.ListMembers(ctx, projectId)
	if err != nil {
		log.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Members {
		memberId := item.MemberID
		email := item.Email
		log.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)
		for _, scope := range item.Scopes {
			log.Printf("Scope: %s\n", scope)
		}
	}

	// update scope
	respMessage, err = mgClient.UpdateMemberScopes(ctx, projectId, updateID, &interfaces.ScopeUpdateRequest{
		Scope: "member",
	})
	if err != nil {
		log.Printf("UpdateMemberScopes failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("UpdateMemberScopes() - Name: %s\n", respMessage.Message)

	// list members
	respGet, err = mgClient.ListMembers(ctx, projectId)
	if err != nil {
		log.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Members {
		memberId := item.MemberID
		email := item.Email
		log.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)
		for _, scope := range item.Scopes {
			log.Printf("Scope: %s\n", scope)
		}
	}
}
