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

	var id string
	var name string
	for _, item := range respList.Projects {
		id = item.ProjectId
		name = item.Name
		log.Printf("ListProjects() - Name: %s, ID: %s\n", name, id)
		break
	}

	// if deleteId == "" {
	// 	log.Printf("This example requires a project who already exists where \"DELETE-ME\" is in the name.\n")
	// 	log.Printf("This is required to exercise the UpdateProject and DeleteProject function.\n")
	// 	log.Printf("In the absence of this, this example will exit early.\n")
	// 	os.Exit(0)
	// }

	// get first project
	respGet, err := mgClient.GetProject(ctx, id)
	if err != nil {
		log.Printf("GetProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetProject() - Name: %s\n", respGet.Name)

	// update project
	respMessage, err := mgClient.UpdateProject(ctx, id, &interfaces.ProjectUpdateRequest{
		Name: "My TEST RENAME Example",
	})
	if err != nil {
		log.Printf("UpdateProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("UpdateProject() - Name: %s\n", respMessage.Message)

	// get project
	respGet, err = mgClient.GetProject(ctx, id)
	if err != nil {
		log.Printf("GetProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetProject() - Name: %s\n", respGet.Name)

	// update project
	respMessage, err = mgClient.UpdateProject(ctx, id, &interfaces.ProjectUpdateRequest{
		Name: name,
	})
	if err != nil {
		log.Printf("UpdateProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("UpdateProject() - Name: %s\n", respMessage.Message)

	// get project
	respGet, err = mgClient.GetProject(ctx, id)
	if err != nil {
		log.Printf("GetProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetProject() - Name: %s\n", respGet.Name)

	// delete project
	// respMessage, err = mgClient.DeleteProject(ctx, deleteId)
	// if err != nil {
	// 	log.Printf("DeleteProject failed. Err: %v\n", err)
	// 	os.Exit(1)
	// }
	// log.Printf("DeleteProject() - Name: %s\n", respMessage.Message)
}
