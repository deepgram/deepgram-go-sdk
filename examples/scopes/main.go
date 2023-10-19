// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

func main() {
	fmt.Print("This example requires a member who already exists where \"@spam.com\" is in the email.\n")
	fmt.Printf("If that use already exists, press ENTER to continue!\n")
	fmt.Printf("Otherwise, Control+C to EXIT!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

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
		projectId = item.ProjectId
		name := item.Name
		log.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectId)
		break
	}

	// list invitations
	respGet, err := mgClient.ListMembers(ctx, projectId)
	if err != nil {
		log.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	var updateID string
	for _, item := range respGet.Members {
		memberId := item.MemberId
		email := item.Email
		log.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)
		for _, scope := range item.Scopes {
			log.Printf("Scope: %s\n", scope)
		}

		if strings.Contains(email, "@spam.com") {
			updateID = memberId
		}
	}

	if updateID == "" {
		fmt.Printf("No member found with \"@spam.com\" in the email. Exiting!\n")
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

	// list invitations
	respGet, err = mgClient.ListMembers(ctx, projectId)
	if err != nil {
		log.Printf("ListMembers failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Members {
		memberId := item.MemberId
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
		memberId := item.MemberId
		email := item.Email
		log.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)
		for _, scope := range item.Scopes {
			log.Printf("Scope: %s\n", scope)
		}
	}
}
