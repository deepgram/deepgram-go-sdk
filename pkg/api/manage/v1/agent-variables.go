// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
Agent Variables API (template variables for reusable agent configurations):
https://developers.deepgram.com/docs/reusable-agent-configurations

Template variables follow the DG_<VARIABLE_NAME> naming format and can
substitute any JSON value in an agent configuration.
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

// ListAgentVariables lists all agent template variables for a project
func (c *Client) ListAgentVariables(ctx context.Context, projectID string) (*api.AgentVariablesResult, error) {
	klog.V(6).Infof("manage.ListAgentVariables() ENTER\n")

	var resp api.AgentVariablesResult
	err := c.APIRequest(ctx, "GET", version.AgentVariablesURI, nil, &resp, projectID)
	if err != nil {
		klog.V(1).Infof("ListAgentVariables failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("ListAgentVariables Succeeded\n")
	}

	klog.V(6).Infof("manage.ListAgentVariables() LEAVE\n")
	return &resp, err
}

// GetAgentVariable gets an agent template variable for a project
func (c *Client) GetAgentVariable(ctx context.Context, projectID, variableID string) (*api.AgentVariableResult, error) {
	klog.V(6).Infof("manage.GetAgentVariable() ENTER\n")

	var resp api.AgentVariableResult
	err := c.APIRequest(ctx, "GET", version.AgentVariablesByIDURI, nil, &resp, projectID, variableID)
	if err != nil {
		klog.V(1).Infof("GetAgentVariable failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("GetAgentVariable Succeeded\n")
	}

	klog.V(6).Infof("manage.GetAgentVariable() LEAVE\n")
	return &resp, err
}

// CreateAgentVariable creates an agent template variable for a project.
// variable.Key follows the DG_<VARIABLE_NAME> naming format; variable.Value can
// be any valid JSON type (string, number, boolean, object, or array).
// The request always carries is_sensitive. The field is optional on the wire
// and the API defaults it to false, which is also the only value it accepts
// today, so callers leave IsSensitive at its zero value.
// The response carries the new agent_variable_uuid and nothing else, so only
// VariableID is populated on the returned result.
func (c *Client) CreateAgentVariable(ctx context.Context, projectID string, variable *api.AgentVariableCreateRequest) (*api.AgentVariableResult, error) {
	klog.V(6).Infof("manage.CreateAgentVariable() ENTER\n")

	jsonStr, err := json.Marshal(variable)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.CreateAgentVariable() LEAVE\n")
		return nil, err
	}

	var resp api.AgentVariableResult
	err = c.APIRequest(ctx, "POST", version.AgentVariablesURI, bytes.NewBuffer(jsonStr), &resp, projectID)
	if err != nil {
		klog.V(1).Infof("CreateAgentVariable failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("CreateAgentVariable Succeeded\n")
	}

	klog.V(6).Infof("manage.CreateAgentVariable() LEAVE\n")
	return &resp, err
}

// UpdateAgentVariable updates the value of an existing agent template variable.
// The API answers with an empty body, so success is signaled by a nil error;
// call GetAgentVariable to read the stored value back.
func (c *Client) UpdateAgentVariable(ctx context.Context, projectID, variableID string, update *api.AgentVariableUpdateRequest) error {
	klog.V(6).Infof("manage.UpdateAgentVariable() ENTER\n")

	jsonStr, err := json.Marshal(update)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("manage.UpdateAgentVariable() LEAVE\n")
		return err
	}

	err = c.APIRequest(ctx, "PATCH", version.AgentVariablesByIDURI, bytes.NewBuffer(jsonStr), nil, projectID, variableID)
	if err != nil {
		klog.V(1).Infof("UpdateAgentVariable failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("UpdateAgentVariable Succeeded\n")
	}

	klog.V(6).Infof("manage.UpdateAgentVariable() LEAVE\n")
	return err
}

// DeleteAgentVariable deletes an agent template variable for a project.
// The API answers with an empty body, so success is signaled by a nil error.
func (c *Client) DeleteAgentVariable(ctx context.Context, projectID, variableID string) error {
	klog.V(6).Infof("manage.DeleteAgentVariable() ENTER\n")

	err := c.APIRequest(ctx, "DELETE", version.AgentVariablesByIDURI, nil, nil, projectID, variableID)
	if err != nil {
		klog.V(1).Infof("DeleteAgentVariable failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("DeleteAgentVariable Succeeded\n")
	}

	klog.V(6).Infof("manage.DeleteAgentVariable() LEAVE\n")
	return err
}
