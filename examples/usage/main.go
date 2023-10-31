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

	// list requests
	respRequestGet, err := mgClient.ListRequests(ctx, projectId, &interfaces.UsageListRequest{})
	if err != nil {
		log.Printf("ListRequests failed. Err: %v\n", err)
		os.Exit(1)
	}

	var requestId string
	for _, item := range respRequestGet.Requests {
		requestId = item.RequestId
		path := item.Path
		log.Printf("ListRequests() - ID: %s, Path: %s\n", requestId, path)
	}

	// get request
	respRequest, err := mgClient.GetRequest(ctx, projectId, requestId)
	if err != nil {
		log.Printf("GetRequest failed. Err: %v\n", err)
		os.Exit(1)
	}
	log.Printf("GetRequest() - ID: %s, Path: %s\n", respRequest.RequestId, respRequest.Path)

	// get fields
	respFieldsGet, err := mgClient.GetFields(ctx, projectId, &interfaces.UsageListRequest{})
	if err != nil {
		log.Printf("GetFields failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, model := range respFieldsGet.Models {
		modelId := model.ModelID
		name := model.Name
		log.Printf("GetFields() - ID: %s, Name: %s\n", modelId, name)
	}
	for _, method := range respFieldsGet.ProcessingMethods {
		log.Printf("GetFields() - method: %s\n", method)
	}

	// get usage
	respUsageGet, err := mgClient.GetUsage(ctx, projectId, &interfaces.UsageRequest{})
	if err != nil {
		log.Printf("GetUsage failed. Err: %v\n", err)
		os.Exit(1)
	}

	for _, item := range respUsageGet.Results {
		numOfRequests := item.Requests
		log.Printf("GetRequest() - %d Calls/%s\n", numOfRequests, respUsageGet.Resolution.Units)
	}
}
