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

	// list invitations
	respGet, err := mgClient.ListInvitations(ctx, projectID)
	if err != nil {
		fmt.Printf("ListInvitations failed. Err: %v\n", err)
		os.Exit(1)
	}

	if len(respGet.Invites) == 0 {
		fmt.Printf("ListInvitations() - No invitations found\n")
	} else {
		for _, item := range respGet.Invites {
			id := item.Email
			scope := item.Scope
			fmt.Printf("ListInvitations() - ID: %s, Scope: %s\n", id, scope)
		}
	}

	// send invite
	respMessage, err := mgClient.SendInvitation(ctx, projectID, &interfaces.InvitationRequest{
		Email: "spam@spam.com",
		Scope: "member",
	})
	if err != nil {
		fmt.Printf("SendInvitation failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SendInvitation() - Result: %s\n", respMessage.Message)

	// list invitations
	respGet, err = mgClient.ListInvitations(ctx, projectID)
	if err != nil {
		fmt.Printf("ListInvitations failed. Err: %v\n", err)
		os.Exit(1)
	}

	if len(respGet.Invites) == 0 {
		fmt.Printf("ListInvitations() - No invitations found\n")
	} else {
		for _, item := range respGet.Invites {
			id := item.Email
			scope := item.Scope
			fmt.Printf("ListInvitations() - ID: %s, Scope: %s\n", id, scope)
		}
	}

	// delete invitation
	respMessage, err = mgClient.DeleteInvitation(ctx, projectID, "spam@spam.com")
	if err != nil {
		fmt.Printf("DeleteInvitation failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DeleteInvitation() - Result: %s\n", respMessage.Message)

	// list invitations
	respGet, err = mgClient.ListInvitations(ctx, projectID)
	if err != nil {
		fmt.Printf("ListInvitations failed. Err: %v\n", err)
		os.Exit(1)
	}

	if len(respGet.Invites) == 0 {
		fmt.Printf("ListInvitations() - No invitations found\n")
	} else {
		for _, item := range respGet.Invites {
			id := item.Email
			scope := item.Scope
			fmt.Printf("ListInvitations() - ID: %s, Scope: %s\n", id, scope)
		}
	}

	// There isnt an API call to add a member to a project. So will leave this commented out as an example
	// Leave Project
	// respMessage, err = mgClient.LeaveProject(ctx, projectID)
	// if err != nil {
	// 	fmt.Printf("LeaveProject failed. Err: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("LeaveProject() - Name: %s\n", respMessage.Message)
}
