// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
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

	// save 1 project
	var projectID string
	for _, item := range respList.Projects {
		projectID = item.ProjectID
		name := item.Name
		fmt.Printf("ListProjects() - Name: %s, ID: %s\n", name, projectID)
		break
	}
	fmt.Printf("\n\n\n")

	// list models
	respModels, err := mgClient.ListModels(ctx, nil)
	if err != nil {
		fmt.Printf("ListModels() failed. Err: %v\n", err)
		os.Exit(1)
	}

	modelId := ""
	if respModels == nil {
		fmt.Printf("ListModels() - No models found\n")
	} else {
		for _, item := range respModels.Stt {
			id := item.UUID
			name := item.Name
			fmt.Printf("STT - ID: %s, Scope: %s\n", id, name)
			modelId = id // save one model id
		}
		for _, item := range respModels.Tts {
			id := item.UUID
			name := item.Name
			fmt.Printf("TTS - ID: %s, Scope: %s\n", id, name)
		}
	}
	fmt.Printf("\n\n\n")

	// get model
	respModel, err := mgClient.GetModel(ctx, modelId)
	if err != nil {
		fmt.Printf("GetModel failed. Err: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("GetModel() - ID: %s, Name: %s\n", respModel.UUID, respModel.Name)
	}
	fmt.Printf("\n\n\n")

	// list project models
	respModels, err = mgClient.ListProjectModels(ctx, projectID, nil)
	if err != nil {
		fmt.Printf("ListProjectModels failed. Err: %v\n", err)
		os.Exit(1)
	}

	modelId = ""
	if respModels == nil {
		fmt.Printf("ListModels() - No models found\n")
	} else {
		for _, item := range respModels.Stt {
			id := item.UUID
			name := item.Name
			fmt.Printf("STT - ID: %s, Scope: %s\n", id, name)
			modelId = id // save one model id
		}
		for _, item := range respModels.Tts {
			id := item.UUID
			name := item.Name
			fmt.Printf("TTS - ID: %s, Scope: %s\n", id, name)
		}
	}
	fmt.Printf("\n\n\n")

	// get model
	respModel, err = mgClient.GetProjectModel(ctx, projectID, modelId)
	if err != nil {
		fmt.Printf("GetProjectModel failed. Err: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("GetProjectModel() - ID: %s, Name: %s\n", respModel.UUID, respModel.Name)
	}
}
