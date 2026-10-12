// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates the reusable agent configuration and agent variable
// management endpoints. It creates an agent configuration and a template
// variable, exercises list/get/update on both, then deletes what it created.
//
// Run:
//
//	DEEPGRAM_API_KEY=<your-key> go run main.go
package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
	apiinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault,
	})

	// context
	ctx := context.Background()

	// client
	dg := client.NewWithDefaults()
	mgClient := api.New(dg)

	// find the first project
	respProject, err := mgClient.ListProjects(ctx)
	if err != nil {
		fmt.Printf("ListProjects failed. Err: %v\n", err)
		os.Exit(1)
	}

	var projectID string
	for _, item := range respProject.Projects {
		projectID = item.ProjectID
		fmt.Printf("Using project: %s (%s)\n\n", item.Name, projectID)
		break
	}

	// create a template variable the agent config can reference
	respVariable, err := mgClient.CreateAgentVariable(ctx, projectID, &apiinterfaces.AgentVariableCreateRequest{
		Key:   "DG_EXAMPLE_ROLE",
		Value: "a helpful weather assistant",
	})
	if err != nil {
		fmt.Printf("CreateAgentVariable failed. Err: %v\n", err)
		os.Exit(1)
	}
	// the create response carries only the new UUID
	variableID := respVariable.VariableID
	fmt.Printf("CreateAgentVariable() - ID: %s\n", variableID)

	// create a reusable agent configuration.
	// Config is a JSON *string* holding the "agent" block of a Settings message.
	// Template variables are referenced bare (unquoted) and substitute a whole
	// JSON value: DG_EXAMPLE_ROLE resolves to its string value at connection time.
	respCreate, err := mgClient.CreateAgent(ctx, projectID, &apiinterfaces.AgentCreateRequest{
		Config: `{
			"listen": {"provider": {"type": "deepgram", "model": "nova-3"}},
			"think": {
				"provider": {"type": "open_ai", "model": "gpt-4o-mini"},
				"prompt": DG_EXAMPLE_ROLE
			},
			"speak": {"provider": {"type": "deepgram", "model": "aura-2-thalia-en"}}
		}`,
		Metadata: map[string]interface{}{"created_by": "go-sdk-example"},
	})
	if err != nil {
		fmt.Printf("CreateAgent failed. Err: %v\n", err)
		os.Exit(1)
	}
	// the create response carries only the new UUID
	agentID := respCreate.AgentID
	fmt.Printf("CreateAgent() - ID: %s\n", agentID)

	// list agent configurations
	respList, err := mgClient.ListAgents(ctx, projectID)
	if err != nil {
		fmt.Printf("ListAgents failed. Err: %v\n", err)
		os.Exit(1)
	}
	for _, item := range respList.Agents {
		fmt.Printf("ListAgents() - ID: %s, Metadata: %v\n", item.AgentID, item.Metadata)
	}

	// get the agent configuration (returned uninterpolated: DG_EXAMPLE_ROLE appears as-is)
	respGet, err := mgClient.GetAgent(ctx, projectID, agentID)
	if err != nil {
		fmt.Printf("GetAgent failed. Err: %v\n", err)
		os.Exit(1)
	}
	// Config comes back as the JSON *string* that was stored, uninterpolated.
	fmt.Printf("GetAgent() - ID: %s, Member: %s, Config: %s\n", respGet.AgentID, respGet.MemberID, respGet.Config)

	// update the agent's metadata (the config itself is immutable).
	// The API answers with an empty body, so read the metadata back with GetAgent.
	// Metadata values are arbitrary JSON, not just strings. This PUT REPLACES the
	// whole metadata object, so created_by is resent even though it is unchanged:
	// dropping it here would delete it. Do not "simplify" this to only the keys
	// being added.
	if err := mgClient.UpdateAgentMetadata(ctx, projectID, agentID, &apiinterfaces.AgentMetadataUpdateRequest{
		Metadata: map[string]interface{}{"created_by": "go-sdk-example", "env": "demo", "revision": 2},
	}); err != nil {
		fmt.Printf("UpdateAgentMetadata failed. Err: %v\n", err)
		os.Exit(1)
	}

	respAfterUpdate, err := mgClient.GetAgent(ctx, projectID, agentID)
	if err != nil {
		fmt.Printf("GetAgent failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("UpdateAgentMetadata() - Metadata: %v\n", respAfterUpdate.Metadata)

	// list, get, and update the template variable
	respVariables, err := mgClient.ListAgentVariables(ctx, projectID)
	if err != nil {
		fmt.Printf("ListAgentVariables failed. Err: %v\n", err)
		os.Exit(1)
	}
	for _, item := range respVariables.Variables {
		fmt.Printf("ListAgentVariables() - ID: %s, Key: %s, Value: %v\n", item.VariableID, item.Key, item.Value)
	}

	respGetVariable, err := mgClient.GetAgentVariable(ctx, projectID, variableID)
	if err != nil {
		fmt.Printf("GetAgentVariable failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("GetAgentVariable() - Key: %s, Value: %v\n", respGetVariable.Key, respGetVariable.Value)

	// PATCH also answers with an empty body, so read the new value back.
	if err := mgClient.UpdateAgentVariable(ctx, projectID, variableID, &apiinterfaces.AgentVariableUpdateRequest{
		Value: "a concise weather assistant",
	}); err != nil {
		fmt.Printf("UpdateAgentVariable failed. Err: %v\n", err)
		os.Exit(1)
	}

	respAfterVariableUpdate, err := mgClient.GetAgentVariable(ctx, projectID, variableID)
	if err != nil {
		fmt.Printf("GetAgentVariable failed. Err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("UpdateAgentVariable() - Value: %v\n", respAfterVariableUpdate.Value)

	// clean up: delete the agent configuration and the variable this example
	// created. Attempt both deletes even if the first one fails, so a failure
	// deleting the agent does not strand the variable in the project.
	cleanupFailed := false

	if err := mgClient.DeleteAgent(ctx, projectID, agentID); err != nil {
		fmt.Printf("DeleteAgent failed. Err: %v\n", err)
		cleanupFailed = true
	} else {
		fmt.Printf("DeleteAgent() - deleted %s\n", agentID)
	}

	if err := mgClient.DeleteAgentVariable(ctx, projectID, variableID); err != nil {
		fmt.Printf("DeleteAgentVariable failed. Err: %v\n", err)
		cleanupFailed = true
	} else {
		fmt.Printf("DeleteAgentVariable() - deleted %s\n", variableID)
	}

	if cleanupFailed {
		os.Exit(1)
	}
}
