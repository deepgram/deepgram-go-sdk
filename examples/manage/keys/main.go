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

	// list keys
	respGet, err := mgClient.ListKeys(ctx, projectID)
	if err != nil {
		fmt.Printf("ListKeys failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.APIKeys {
		id := item.APIKey.APIKeyID
		comment := item.APIKey.Comment
		fmt.Printf("ListKeys() - ID: %s, Comment: %s\n", id, comment)
	}

	// create key
	respCreate, err := mgClient.CreateKey(ctx, projectID, &interfaces.KeyCreateRequest{
		Comment: "My Test",
		Scopes:  []string{"onprem:products"},
	})
	if err != nil {
		fmt.Printf("CreateKey failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("CreateKey() - Name: %s\n", respCreate.Comment)

	// list keys
	respGet, err = mgClient.ListKeys(ctx, projectID)
	if err != nil {
		fmt.Printf("ListKeys failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.APIKeys {
		id := item.APIKey.APIKeyID
		comment := item.APIKey.Comment
		fmt.Printf("ListKeys() - ID: %s, Comment: %s\n", id, comment)
	}

	// get key
	respKey, err := mgClient.GetKey(ctx, projectID, respCreate.APIKeyID)
	if err != nil {
		fmt.Printf("GetKey failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("GetKey() - ID: %s, Comment: %s\n", respKey.APIKey.APIKeyID, respKey.APIKey.Comment)

	// delete key
	respMessage, err := mgClient.DeleteKey(ctx, projectID, respKey.APIKey.APIKeyID)
	if err != nil {
		fmt.Printf("DeleteKey failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DeleteKey() - Name: %s\n", respMessage.Message)

	// list keys
	respGet, err = mgClient.ListKeys(ctx, projectID)
	if err != nil {
		fmt.Printf("ListKeys failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respGet.APIKeys {
		id := item.APIKey.APIKeyID
		comment := item.APIKey.Comment
		fmt.Printf("ListKeys() - ID: %s, Comment: %s\n", id, comment)
	}
}
