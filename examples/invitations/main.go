// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"log"
	"os"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

func main() {
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
	respGet, err := mgClient.ListInvitations(ctx, projectId)
	if err != nil {
		log.Printf("ListInvitations failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Invites {
		id := item.Email
		scope := item.Scope
		log.Printf("ListInvitations() - ID: %s, Scope: %s\n", id, scope)
	}

	// send invite
	respMessage, err := mgClient.SendInvitation(ctx, projectId, &interfaces.InvitationRequest{
		Email: "spam@spam.com",
		Scope: "member",
	})
	if err != nil {
		log.Printf("SendInvitation failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("SendInvitation() - Name: %s\n", respMessage.Message)

	// list invitations
	respGet, err = mgClient.ListInvitations(ctx, projectId)
	if err != nil {
		log.Printf("ListInvitations failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Invites {
		id := item.Email
		scope := item.Scope
		log.Printf("ListInvitations() - ID: %s, Scope: %s\n", id, scope)
	}

	// delete project
	respMessage, err = mgClient.DeleteInvitation(ctx, projectId, "spam@spam.com")
	if err != nil {
		log.Printf("DeleteInvitation failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("DeleteInvitation() - Name: %s\n", respMessage.Message)

	// list invitations
	respGet, err = mgClient.ListInvitations(ctx, projectId)
	if err != nil {
		log.Printf("ListInvitations failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.Invites {
		id := item.Email
		scope := item.Scope
		log.Printf("ListInvitations() - ID: %s, Scope: %s\n", id, scope)
	}

	// TODO: There isnt an API call to add a member to a project. So will leave this commented out as an example
	// Leave Project
	// respMessage, err = mgClient.LeaveProject(ctx, projectId)
	// if err != nil {
	// 	log.Printf("LeaveProject failed. Err: %v\n", err)
	// 	os.Exit(1)
	// }
	// log.Printf("LeaveProject() - Name: %s\n", respMessage.Message)
}
