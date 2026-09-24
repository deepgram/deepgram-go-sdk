// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Agents API (reusable voice-agent configurations):
https://developers.deepgram.com/docs/reusable-agent-configurations

Save an agent's core settings (the "agent" block of a Settings message) once and
get back an agent_uuid to pass in place of the full agent object on every
Voice Agent connection.
*/
package manage

import (
	"bytes"
	"context"
	"encoding/json"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
)

// ListAgents lists all agent configurations for a project. Configurations are
// returned in their uninterpolated form — template variable placeholders
// (DG_<VARIABLE_NAME>) appear as-is rather than with their substituted values.
func (c *Client) ListAgents(ctx context.Context, projectID string) (*api.AgentsResult, error) {
	klog.V(6).Infof("manage.ListAgents() ENTER\n")

	var resp api.AgentsResult
	err := c.APIRequest(ctx, "GET", version.AgentsURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListAgents failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListAgents Succeeded\n")
	}

	klog.V(6).Infof("manage.ListAgents() LEAVE\n")
	return &resp, err
}

// GetAgent gets an agent configuration for a project, in its uninterpolated
// form. The returned Config is the JSON string that was stored, echoed back
// verbatim rather than as a parsed object.
func (c *Client) GetAgent(ctx context.Context, projectID, agentID string) (*api.AgentResult, error) {
	klog.V(6).Infof("manage.GetAgent() ENTER\n")

	var resp api.AgentResult
	err := c.APIRequest(ctx, "GET", version.AgentsByIDURI, nil, &resp, projectID, agentID)
	if err != nil {
		klog.V(1).Infof("GetAgent failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetAgent Succeeded\n")
	}

	klog.V(6).Infof("manage.GetAgent() LEAVE\n")
	return &resp, err
}

// CreateAgent creates a reusable agent configuration for a project.
// agent.Config is a JSON string holding the "agent" block of a Settings
// message; template variables are referenced bare (unquoted) inside it, so a
// config that uses them is not parseable JSON until the server substitutes
// their values. The response carries the new agent_uuid and nothing else, so
// only AgentID is populated on the returned result; it can be passed in place
// of the full agent object in future Settings messages.
func (c *Client) CreateAgent(ctx context.Context, projectID string, agent *api.AgentCreateRequest) (*api.AgentResult, error) {
	klog.V(6).Infof("manage.CreateAgent() ENTER\n")

	jsonStr, err := json.Marshal(agent)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.CreateAgent() LEAVE\n")
		return nil, err
	}

	var resp api.AgentResult
	err = c.APIRequest(ctx, "POST", version.AgentsURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("CreateAgent failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("CreateAgent Succeeded\n")
	}

	klog.V(6).Infof("manage.CreateAgent() LEAVE\n")
	return &resp, err
}

// UpdateAgentMetadata updates the metadata associated with an agent
// configuration. The config itself is immutable — to change the configuration,
// delete the existing agent and create a new one.
//
// This REPLACES the entire metadata object: any key absent from update.Metadata
// is deleted, with no error and no warning. Call GetAgent first and resend every
// key you mean to keep.
//
// The API answers with an empty body, so success is signaled by a nil error;
// call GetAgent to read the stored metadata back.
func (c *Client) UpdateAgentMetadata(ctx context.Context, projectID, agentID string, update *api.AgentMetadataUpdateRequest) error {
	klog.V(6).Infof("manage.UpdateAgentMetadata() ENTER\n")

	jsonStr, err := json.Marshal(update)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateAgentMetadata() LEAVE\n")
		return err
	}

	err = c.APIRequest(ctx, "PUT", version.AgentsByIDURI, bytes.NewBuffer(jsonStr), nil, projectID, agentID)
	if err != nil {
		klog.V(1).Infof("UpdateAgentMetadata failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("UpdateAgentMetadata Succeeded\n")
	}

	klog.V(6).Infof("manage.UpdateAgentMetadata() LEAVE\n")
	return err
}

// DeleteAgent deletes an agent configuration for a project.
// Deleting an agent configuration can cause a production outage if your service
// references this agent UUID — migrate all active sessions to a new
// configuration before deleting.
// The API answers with an empty body, so success is signaled by a nil error.
func (c *Client) DeleteAgent(ctx context.Context, projectID, agentID string) error {
	klog.V(6).Infof("manage.DeleteAgent() ENTER\n")

	err := c.APIRequest(ctx, "DELETE", version.AgentsByIDURI, nil, nil, projectID, agentID)
	if err != nil {
		klog.V(1).Infof("DeleteAgent failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("DeleteAgent Succeeded\n")
	}

	klog.V(6).Infof("manage.DeleteAgent() LEAVE\n")
	return err
}
