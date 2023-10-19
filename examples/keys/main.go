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

	// list keys
	respGet, err := mgClient.ListKeys(ctx, projectId)
	if err != nil {
		log.Printf("ListKeys failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.ApiKeys {
		id := item.ApiKeyId
		comment := item.Comment
		log.Printf("ListKeys() - ID: %s, Comment: %s\n", id, comment)
		break
	}

	// create key
	respCreate, err := mgClient.CreateKey(ctx, projectId, &interfaces.KeyCreateRequest{
		Comment: "My Test",
		Scopes:  []string{"onprem:products"},
	})
	if err != nil {
		log.Printf("CreateKey failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("CreateKey() - Name: %s\n", respCreate.Comment)

	// get key
	respKey, err := mgClient.GetKey(ctx, projectId, respCreate.ApiKeyId)
	if err != nil {
		log.Printf("GetKey failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetKey() - ID: %s, Comment: %s\n", respKey.Key.ApiKeyId, respKey.Key.Comment)

	// delete project
	respMessage, err := mgClient.DeleteKey(ctx, projectId, respKey.Key.ApiKeyId)
	if err != nil {
		log.Printf("DeleteKey failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("DeleteKey() - Name: %s\n", respMessage.Message)

	// list invitations
	respGet, err = mgClient.ListKeys(ctx, projectId)
	if err != nil {
		log.Printf("ListKeys failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.ApiKeys {
		id := item.ApiKeyId
		comment := item.Comment
		log.Printf("ListKeys() - ID: %s, Comment: %s\n", id, comment)
	}
}
