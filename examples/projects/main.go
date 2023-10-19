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
	fmt.Print("This example requires a project who already exists where \"TEST\" is in the name.\n")
	fmt.Printf("If that project already exists, press ENTER to continue!\n")
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

	var deleteId string
	var deleteName string
	for _, item := range respList.Projects {
		id := item.ProjectId
		name := item.Name
		log.Printf("ListProjects() - Name: %s, ID: %s\n", name, id)

		if strings.Contains(name, "TEST") {
			deleteId = id
			deleteName = name
		}
	}

	if deleteId == "" {
		fmt.Printf("No project found with \"TEST\" in the name. Exiting!\n")
		os.Exit(0)
	}

	// get first project
	respGet, err := mgClient.GetProject(ctx, deleteId)
	if err != nil {
		log.Printf("GetProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetProject() - Name: %s\n", respGet.Name)

	// update project
	respMessage, err := mgClient.UpdateProject(ctx, deleteId, &interfaces.ProjectUpdateRequest{
		Name: "My TEST RENAME Example",
	})
	if err != nil {
		log.Printf("UpdateProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("UpdateProject() - Name: %s\n", respMessage.Message)

	// get project
	respGet, err = mgClient.GetProject(ctx, deleteId)
	if err != nil {
		log.Printf("GetProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetProject() - Name: %s\n", respGet.Name)

	// update project
	respMessage, err = mgClient.UpdateProject(ctx, deleteId, &interfaces.ProjectUpdateRequest{
		Name: deleteName,
	})
	if err != nil {
		log.Printf("UpdateProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("UpdateProject() - Name: %s\n", respMessage.Message)

	// get project
	respGet, err = mgClient.GetProject(ctx, deleteId)
	if err != nil {
		log.Printf("GetProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetProject() - Name: %s\n", respGet.Name)

	// delete project
	respMessage, err = mgClient.DeleteProject(ctx, deleteId)
	if err != nil {
		log.Printf("DeleteProject failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("DeleteProject() - Name: %s\n", respMessage.Message)
}
