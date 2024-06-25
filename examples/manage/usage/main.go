// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	prettyjson "github.com/hokaccha/go-prettyjson"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/manage"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault, // LogLevelDefault, LogLevelFull, LogLevelDebug, LogLevelTrace
	})

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

	// list requests
	respRequestGet, err := mgClient.ListRequests(ctx, projectID, &interfaces.UsageListRequest{})
	if err != nil {
		fmt.Printf("ListRequests failed. Err: %v\n", err)
		os.Exit(1)
	}

	// dump
	data, err := json.Marshal(respRequestGet)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	// make the JSON pretty
	prettyJSON, err := prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\nResult:\n%s\n\n", prettyJSON)

	var requestID string
	for _, item := range respRequestGet.Requests {
		requestID = item.RequestID
		break
	}

	// get request
	respRequest, err := mgClient.GetRequest(ctx, projectID, requestID)
	if err != nil {
		fmt.Printf("GetRequest failed. Err: %v\n", err)
		os.Exit(1)
	}

	// dump
	data, err = json.Marshal(respRequest)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	// make the JSON pretty
	prettyJSON, err = prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\nResult:\n%s\n\n", prettyJSON)

	// get fields
	respFieldsGet, err := mgClient.GetFields(ctx, projectID, &interfaces.UsageListRequest{})
	if err != nil {
		fmt.Printf("GetFields failed. Err: %v\n", err)
		os.Exit(1)
	}

	// dump
	data, err = json.Marshal(respFieldsGet)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	// make the JSON pretty
	prettyJSON, err = prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\nResult:\n%s\n\n", prettyJSON)

	// get usage
	respUsageGet, err := mgClient.GetUsage(ctx, projectID, &interfaces.UsageRequest{})
	if err != nil {
		fmt.Printf("GetUsage failed. Err: %v\n", err)
		os.Exit(1)
	}

	// dump
	data, err = json.Marshal(respUsageGet)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	// make the JSON pretty
	prettyJSON, err = prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n\nResult:\n%s\n\n", prettyJSON)
}
