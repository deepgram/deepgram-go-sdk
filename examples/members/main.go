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

	var delMemberId string
	for _, item := range respGet.Members {
		memberId := item.MemberId
		email := item.Email
		log.Printf("ListMembers() - ID: %s, Scope: %s\n", memberId, email)

		if strings.Contains(email, "@spam.com") {
			delMemberId = memberId
		}
	}

	if delMemberId == "" {
		fmt.Printf("Unable to find member with email containing \"@spam.com\". Exiting!\n")
		os.Exit(0)
	}

	// remove member
	respMessage, err := mgClient.RemoveMember(ctx, projectId, delMemberId)
	if err != nil {
		log.Printf("RemoveMember failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("RemoveMember() - Name: %s\n", respMessage.Message)

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
	}
}
