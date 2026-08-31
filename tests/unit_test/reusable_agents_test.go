// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
	apiinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	manageclient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

// The response bodies in this file are the shapes the live API returned when the
// endpoints were probed on 2026-09-23. They are deliberately verbatim: the
// published reference and openapi.yaml disagree with production on the id field
// names, on whether the list endpoints wrap their array, and on whether config
// comes back as a JSON string or a parsed object.
const (
	testProjectID  = "9db43a4f-4f24-4c0c-a2fc-1c0f19f77869"
	testAgentID    = "3a1c22de-9d20-4b23-b4a3-04a30b71d0b7"
	testVariableID = "f9d254d3-b0e8-4b19-8ff2-30a67b2b3c98"
	testMemberID   = "43d2a494-1f0e-4a93-a0a9-6a7e55a6a0f1"

	// The exact config string a create request sends and a read echoes back.
	testConfigString = `{"language":"en","think":{"prompt":DG_PROMPT}}`

	testVariableValue = "a helpful assistant"

	// The bare JSON array both list endpoints return when a project is empty.
	emptyJSONArray = "[]"
)

// recordedRequest captures what the SDK actually sent for wire-level assertions.
type recordedRequest struct {
	Method string
	Path   string
	Body   []byte
}

// newManageTestServer starts an httptest server that records the request and
// returns the given JSON body, plus a manage API client pointed at it. An empty
// responseJSON reproduces the 200-with-no-body the update and delete endpoints
// actually answer with.
func newManageTestServer(t *testing.T, responseJSON string) (*recordedRequest, *api.Client, func()) {
	t.Helper()

	recorded := &recordedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorded.Method = r.Method
		recorded.Path = r.URL.Path
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		recorded.Body = body
		if responseJSON == "" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(responseJSON)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))

	dg := manageclient.New("testkey", &interfaces.ClientOptions{Host: srv.URL})
	return recorded, api.New(dg), srv.Close
}

// agentReadBody is the body GET /agents/{id} returns, and also one element of
// the bare array GET /agents returns.
func agentReadBody(t *testing.T) string {
	t.Helper()
	body, err := json.Marshal(map[string]interface{}{
		"agent_uuid":  testAgentID,
		"member_id":   testMemberID,
		"api_version": 1,
		"config":      testConfigString,
		"metadata":    map[string]string{"env": "prod"},
	})
	if err != nil {
		t.Fatalf("failed to build the agent fixture: %v", err)
	}
	return string(body)
}

// Test_AgentsAPI exercises the reusable agent configuration endpoints end to end
// against a mock server: request method, path, body serialization, and response parsing.
func Test_AgentsAPI(t *testing.T) {
	t.Run("ListAgents parses the bare top-level JSON array", func(t *testing.T) {
		// The live endpoint returns an array, NOT {"agents":[...]}. A model that
		// expects the wrapper decodes this body into an empty list without erroring.
		recorded, mgClient, closeSrv := newManageTestServer(t, "["+agentReadBody(t)+"]")
		defer closeSrv()

		resp, err := mgClient.ListAgents(context.Background(), testProjectID)
		if err != nil {
			t.Fatalf("ListAgents failed: %v", err)
		}

		if recorded.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agents"; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		if len(resp.Agents) != 1 {
			t.Fatalf("expected 1 agent from the bare array, got %d", len(resp.Agents))
		}
		agent := resp.Agents[0]
		if agent.AgentID != testAgentID {
			t.Errorf("expected agent_uuid %s, got %q", testAgentID, agent.AgentID)
		}
		if agent.MemberID != testMemberID {
			t.Errorf("expected member_id %s, got %q", testMemberID, agent.MemberID)
		}
		if agent.APIVersion != 1 {
			t.Errorf("expected api_version 1, got %d", agent.APIVersion)
		}
		if agent.Config != testConfigString {
			t.Errorf("expected config string %q, got %q", testConfigString, agent.Config)
		}
		if agent.Metadata["env"] != "prod" {
			t.Errorf("expected metadata.env prod, got %v", agent.Metadata["env"])
		}
	})

	t.Run("ListAgents does not accept the wrapped object shape the spec documents", func(t *testing.T) {
		// Proves the decoder really reads an array: the documented-but-wrong
		// {"agents":[...]} body must not silently produce a populated list.
		_, mgClient, closeSrv := newManageTestServer(t, `{"agents":[`+agentReadBody(t)+`]}`)
		defer closeSrv()

		resp, err := mgClient.ListAgents(context.Background(), testProjectID)
		if err == nil && len(resp.Agents) > 0 {
			t.Fatalf(`a wrapped {"agents":[...]} body must not decode as a list, got %d agents`, len(resp.Agents))
		}
	})

	t.Run("CreateAgent sends config as a JSON string and parses the uuid-only response", func(t *testing.T) {
		// The live create response carries the new uuid and nothing else.
		recorded, mgClient, closeSrv := newManageTestServer(t, `{"agent_uuid":"`+testAgentID+`"}`)
		defer closeSrv()

		resp, err := mgClient.CreateAgent(context.Background(), testProjectID, &apiinterfaces.AgentCreateRequest{
			Config:   testConfigString,
			Metadata: map[string]interface{}{"team": "devrel"},
		})
		if err != nil {
			t.Fatalf("CreateAgent failed: %v", err)
		}

		if recorded.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agents"; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}

		// The config field must serialize as a JSON *string*, not a nested object.
		var body map[string]interface{}
		if err := json.Unmarshal(recorded.Body, &body); err != nil {
			t.Fatalf("request body is not valid JSON: %v", err)
		}
		configStr, ok := body["config"].(string)
		if !ok {
			t.Fatalf("expected config to be a JSON string, got %T", body["config"])
		}
		if configStr != testConfigString {
			t.Errorf("unexpected config string: %s", configStr)
		}
		if _, present := body["api_version"]; present {
			t.Errorf("api_version should be omitted when unset, body: %s", recorded.Body)
		}

		if resp.AgentID != testAgentID {
			t.Errorf("expected agent_uuid %s, got %q", testAgentID, resp.AgentID)
		}
	})

	t.Run("GetAgent decodes the string config the server echoes back", func(t *testing.T) {
		// A map[string]interface{} config field does not merely come back empty
		// here: json.Unmarshal fails outright putting a string into a map, so this
		// subtest fails loudly against the pre-fix model.
		recorded, mgClient, closeSrv := newManageTestServer(t, agentReadBody(t))
		defer closeSrv()

		resp, err := mgClient.GetAgent(context.Background(), testProjectID, testAgentID)
		if err != nil {
			t.Fatalf("GetAgent failed: %v", err)
		}
		if recorded.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agents/" + testAgentID; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		if resp.AgentID != testAgentID {
			t.Errorf("expected agent_uuid %s, got %q", testAgentID, resp.AgentID)
		}
		if resp.MemberID != testMemberID {
			t.Errorf("expected member_id %s, got %q", testMemberID, resp.MemberID)
		}
		if resp.APIVersion != 1 {
			t.Errorf("expected api_version 1, got %d", resp.APIVersion)
		}
		if resp.Config != testConfigString {
			t.Errorf("expected config %q, got %q", testConfigString, resp.Config)
		}
	})

	t.Run("the agent config round trip is symmetric", func(t *testing.T) {
		// Request and response both carry config as a string, so re-marshaling a
		// read result reproduces exactly the wire shape it was decoded from.
		_, mgClient, closeSrv := newManageTestServer(t, agentReadBody(t))
		defer closeSrv()

		resp, err := mgClient.GetAgent(context.Background(), testProjectID, testAgentID)
		if err != nil {
			t.Fatalf("GetAgent failed: %v", err)
		}

		remarshaled, err := json.Marshal(resp.AgentConfiguration)
		if err != nil {
			t.Fatalf("re-marshal failed: %v", err)
		}
		var got, want map[string]interface{}
		if err := json.Unmarshal(remarshaled, &got); err != nil {
			t.Fatalf("re-marshaled body is not valid JSON: %v", err)
		}
		if err := json.Unmarshal([]byte(agentReadBody(t)), &want); err != nil {
			t.Fatalf("fixture is not valid JSON: %v", err)
		}
		for _, key := range []string{"agent_uuid", "member_id", "api_version", "config"} {
			if got[key] != want[key] {
				t.Errorf("round trip changed %s: sent %v, got %v", key, want[key], got[key])
			}
		}
		for _, phantom := range []string{"created_at", "updated_at", "agent_id"} {
			if _, present := got[phantom]; present {
				t.Errorf("%s is not part of the wire contract and must not be emitted", phantom)
			}
		}
	})

	t.Run("UpdateAgentMetadata sends PUT and tolerates the empty response body", func(t *testing.T) {
		// The live endpoint answers 200 with no body. Decoding an empty reader
		// fails with io.EOF, so this catches any model that still expects a body.
		recorded, mgClient, closeSrv := newManageTestServer(t, "")
		defer closeSrv()

		err := mgClient.UpdateAgentMetadata(context.Background(), testProjectID, testAgentID,
			&apiinterfaces.AgentMetadataUpdateRequest{Metadata: map[string]interface{}{"env": "staging"}})
		if err != nil {
			t.Fatalf("UpdateAgentMetadata failed: %v", err)
		}
		if recorded.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agents/" + testAgentID; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		if want := `{"metadata":{"env":"staging"}}`; string(recorded.Body) != want {
			t.Errorf("expected body %s, got %s", want, recorded.Body)
		}
	})

	t.Run("DeleteAgent sends DELETE and tolerates the empty response body", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, "")
		defer closeSrv()

		if err := mgClient.DeleteAgent(context.Background(), testProjectID, testAgentID); err != nil {
			t.Fatalf("DeleteAgent failed: %v", err)
		}
		if recorded.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agents/" + testAgentID; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
	})

	t.Run("errors from the server propagate to the caller", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{"err_code":"BAD_REQUEST","err_msg":"This project already has a variable with that name"}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer srv.Close()

		mgClient := api.New(manageclient.New("testkey", &interfaces.ClientOptions{Host: srv.URL}))

		if _, err := mgClient.CreateAgent(context.Background(), testProjectID,
			&apiinterfaces.AgentCreateRequest{Config: "not json"}); err == nil {
			t.Fatal("expected an error from a 400 response, got nil")
		}
		// The empty-body endpoints must still surface a failure even though they
		// never decode anything on the happy path.
		if err := mgClient.DeleteAgent(context.Background(), testProjectID, testAgentID); err == nil {
			t.Fatal("expected an error from a 400 response, got nil")
		}
		if err := mgClient.UpdateAgentMetadata(context.Background(), testProjectID, testAgentID,
			&apiinterfaces.AgentMetadataUpdateRequest{Metadata: map[string]interface{}{"a": "b"}}); err == nil {
			t.Fatal("expected an error from a 400 response, got nil")
		}
	})
}

// Test_AgentVariablesAPI exercises the agent template variable endpoints end to
// end against a mock server.
func Test_AgentVariablesAPI(t *testing.T) {
	variableReadBody := `{
		"agent_variable_uuid": "` + testVariableID + `",
		"member_id": "` + testMemberID + `",
		"api_version": 1,
		"key": "DG_ROLE",
		"value": "` + testVariableValue + `",
		"is_sensitive": false
	}`

	t.Run("ListAgentVariables parses the bare top-level JSON array", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, "["+variableReadBody+"]")
		defer closeSrv()

		resp, err := mgClient.ListAgentVariables(context.Background(), testProjectID)
		if err != nil {
			t.Fatalf("ListAgentVariables failed: %v", err)
		}
		if recorded.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agent-variables"; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		if len(resp.Variables) != 1 {
			t.Fatalf("expected 1 variable from the bare array, got %d", len(resp.Variables))
		}
		variable := resp.Variables[0]
		if variable.VariableID != testVariableID {
			t.Errorf("expected agent_variable_uuid %s, got %q", testVariableID, variable.VariableID)
		}
		if variable.MemberID != testMemberID {
			t.Errorf("expected member_id %s, got %q", testMemberID, variable.MemberID)
		}
		if variable.APIVersion != 1 {
			t.Errorf("expected api_version 1, got %d", variable.APIVersion)
		}
		if variable.Key != "DG_ROLE" {
			t.Errorf("expected key DG_ROLE, got %s", variable.Key)
		}
		if variable.Value != testVariableValue {
			t.Errorf("unexpected value: %v", variable.Value)
		}
		if variable.IsSensitive {
			t.Error("expected is_sensitive false")
		}
	})

	t.Run("ListAgentVariables does not accept the wrapped object shape the spec documents", func(t *testing.T) {
		_, mgClient, closeSrv := newManageTestServer(t, `{"variables":[`+variableReadBody+`]}`)
		defer closeSrv()

		resp, err := mgClient.ListAgentVariables(context.Background(), testProjectID)
		if err == nil && len(resp.Variables) > 0 {
			t.Fatalf(`a wrapped {"variables":[...]} body must not decode as a list, got %d`, len(resp.Variables))
		}
	})

	t.Run("CreateAgentVariable serializes string values and parses the uuid-only response", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, `{"agent_variable_uuid":"`+testVariableID+`"}`)
		defer closeSrv()

		resp, err := mgClient.CreateAgentVariable(context.Background(), testProjectID,
			&apiinterfaces.AgentVariableCreateRequest{Key: "DG_ROLE", Value: testVariableValue})
		if err != nil {
			t.Fatalf("CreateAgentVariable failed: %v", err)
		}
		if recorded.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agent-variables"; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		// is_sensitive is optional on the wire (the API defaults it to false) but
		// the SDK always sends it so the request states that intent explicitly.
		if want := `{"key":"DG_ROLE","value":"` + testVariableValue + `","is_sensitive":false}`; string(recorded.Body) != want {
			t.Errorf("expected body %s, got %s", want, recorded.Body)
		}
		if resp.VariableID != testVariableID {
			t.Errorf("expected agent_variable_uuid %s, got %q", testVariableID, resp.VariableID)
		}
	})

	t.Run("CreateAgentVariable serializes non-string JSON values", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, `{"agent_variable_uuid":"`+testVariableID+`"}`)
		defer closeSrv()

		if _, err := mgClient.CreateAgentVariable(context.Background(), testProjectID,
			&apiinterfaces.AgentVariableCreateRequest{
				Key:   "DG_LIMITS",
				Value: map[string]interface{}{"max_turns": 10},
			}); err != nil {
			t.Fatalf("CreateAgentVariable failed: %v", err)
		}
		if want := `{"key":"DG_LIMITS","value":{"max_turns":10},"is_sensitive":false}`; string(recorded.Body) != want {
			t.Errorf("expected body %s, got %s", want, recorded.Body)
		}
	})

	t.Run("GetAgentVariable hits the by-id path and reads every field", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, variableReadBody)
		defer closeSrv()

		resp, err := mgClient.GetAgentVariable(context.Background(), testProjectID, testVariableID)
		if err != nil {
			t.Fatalf("GetAgentVariable failed: %v", err)
		}
		if recorded.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agent-variables/" + testVariableID; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		if resp.VariableID != testVariableID {
			t.Errorf("expected agent_variable_uuid %s, got %q", testVariableID, resp.VariableID)
		}
		if resp.Value != testVariableValue {
			t.Errorf("unexpected value: %v", resp.Value)
		}

		// created_at / updated_at are not part of the contract and must not be
		// re-emitted; is_sensitive false must survive the round trip.
		remarshaled, err := json.Marshal(resp.AgentVariable)
		if err != nil {
			t.Fatalf("re-marshal failed: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal(remarshaled, &got); err != nil {
			t.Fatalf("re-marshaled body is not valid JSON: %v", err)
		}
		for _, phantom := range []string{"created_at", "updated_at", "variable_id"} {
			if _, present := got[phantom]; present {
				t.Errorf("%s is not part of the wire contract and must not be emitted", phantom)
			}
		}
		sensitive, present := got["is_sensitive"]
		if !present {
			t.Error("is_sensitive must survive a round trip even when false")
		} else if sensitive != false {
			t.Errorf("expected is_sensitive false, got %v", sensitive)
		}
	})

	t.Run("UpdateAgentVariable sends PATCH and tolerates the empty response body", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, "")
		defer closeSrv()

		err := mgClient.UpdateAgentVariable(context.Background(), testProjectID, testVariableID,
			&apiinterfaces.AgentVariableUpdateRequest{Value: 123})
		if err != nil {
			t.Fatalf("UpdateAgentVariable failed: %v", err)
		}
		if recorded.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agent-variables/" + testVariableID; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
		if want := `{"value":123}`; string(recorded.Body) != want {
			t.Errorf("expected body %s, got %s", want, recorded.Body)
		}
	})

	t.Run("DeleteAgentVariable sends DELETE and tolerates the empty response body", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, "")
		defer closeSrv()

		if err := mgClient.DeleteAgentVariable(context.Background(), testProjectID, testVariableID); err != nil {
			t.Fatalf("DeleteAgentVariable failed: %v", err)
		}
		if recorded.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", recorded.Method)
		}
		if want := "/v1/projects/" + testProjectID + "/agent-variables/" + testVariableID; recorded.Path != want {
			t.Errorf("expected path %s, got %s", want, recorded.Path)
		}
	})
}

// Test_AgentListsAreBareArrays pins the list representation itself. Both list
// endpoints answer with a bare JSON array at the top level, so the list types
// unmarshal from and marshal back to an array rather than a wrapper object,
// while keeping the .Agents / .Variables accessors the rest of the package uses.
func Test_AgentListsAreBareArrays(t *testing.T) {
	t.Run("agents", func(t *testing.T) {
		var list apiinterfaces.AgentsList
		if err := json.Unmarshal([]byte(`[{"agent_uuid":"`+testAgentID+`","config":"{}"}]`), &list); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if len(list.Agents) != 1 || list.Agents[0].AgentID != testAgentID {
			t.Fatalf("unexpected list: %+v", list.Agents)
		}
		out, err := json.Marshal(list)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		if len(out) == 0 || out[0] != '[' {
			t.Errorf("expected a bare JSON array, got %s", out)
		}
	})

	t.Run("agents empty", func(t *testing.T) {
		out, err := json.Marshal(apiinterfaces.AgentsList{})
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		if string(out) != emptyJSONArray {
			t.Errorf("expected %s, got %s", emptyJSONArray, out)
		}
	})

	t.Run("variables", func(t *testing.T) {
		var list apiinterfaces.AgentVariablesList
		if err := json.Unmarshal([]byte(`[{"agent_variable_uuid":"`+testVariableID+`","key":"DG_K"}]`), &list); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if len(list.Variables) != 1 || list.Variables[0].VariableID != testVariableID {
			t.Fatalf("unexpected list: %+v", list.Variables)
		}
		out, err := json.Marshal(list)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		if len(out) == 0 || out[0] != '[' {
			t.Errorf("expected a bare JSON array, got %s", out)
		}
	})

	t.Run("variables empty", func(t *testing.T) {
		out, err := json.Marshal(apiinterfaces.AgentVariablesList{})
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		if string(out) != emptyJSONArray {
			t.Errorf("expected %s, got %s", emptyJSONArray, out)
		}
	})
}

// Test_AgentMetadataAcceptsArbitraryJSONValues pins the metadata contract. The
// API stores and returns whatever JSON a client wrote, so metadata set by the
// console or another SDK can hold numbers, booleans, and nested objects. A
// map[string]string field does not merely drop those values: it fails the
// decode of the entire response with "cannot unmarshal number into Go struct
// field AgentResult.metadata of type string", so GetAgent returns nothing at all.
func Test_AgentMetadataAcceptsArbitraryJSONValues(t *testing.T) {
	const mixedMetadataBody = `{"agent_uuid":"` + testAgentID + `","member_id":"` + testMemberID + `",` +
		`"api_version":1,"config":"{}",` +
		`"metadata":{"created_by":"console","count":7,"flag":true,"nested":{"a":1}}}`

	assertMixedMetadata := func(t *testing.T, metadata map[string]interface{}) {
		t.Helper()
		if got, ok := metadata["created_by"].(string); !ok || got != "console" {
			t.Errorf("expected metadata.created_by %q, got %#v", "console", metadata["created_by"])
		}
		if got, ok := metadata["count"].(float64); !ok || got != 7 {
			t.Errorf("expected metadata.count 7, got %#v", metadata["count"])
		}
		if got, ok := metadata["flag"].(bool); !ok || !got {
			t.Errorf("expected metadata.flag true, got %#v", metadata["flag"])
		}
		nested, ok := metadata["nested"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected metadata.nested to decode as an object, got %#v", metadata["nested"])
		}
		if got, ok := nested["a"].(float64); !ok || got != 1 {
			t.Errorf("expected metadata.nested.a 1, got %#v", nested["a"])
		}
	}

	t.Run("GetAgent decodes non-string metadata values", func(t *testing.T) {
		_, mgClient, closeSrv := newManageTestServer(t, mixedMetadataBody)
		defer closeSrv()

		resp, err := mgClient.GetAgent(context.Background(), testProjectID, testAgentID)
		if err != nil {
			t.Fatalf("GetAgent failed on non-string metadata: %v", err)
		}
		// The rest of the response must survive too: a metadata type error
		// takes the whole object down, not just the one field.
		if resp.AgentID != testAgentID {
			t.Errorf("expected agent_uuid %s, got %q", testAgentID, resp.AgentID)
		}
		if resp.Config != "{}" {
			t.Errorf("expected config %q, got %q", "{}", resp.Config)
		}
		assertMixedMetadata(t, resp.Metadata)
	})

	t.Run("ListAgents decodes non-string metadata values", func(t *testing.T) {
		_, mgClient, closeSrv := newManageTestServer(t, "["+mixedMetadataBody+"]")
		defer closeSrv()

		resp, err := mgClient.ListAgents(context.Background(), testProjectID)
		if err != nil {
			t.Fatalf("ListAgents failed on non-string metadata: %v", err)
		}
		if len(resp.Agents) != 1 {
			t.Fatalf("expected 1 agent, got %d", len(resp.Agents))
		}
		assertMixedMetadata(t, resp.Agents[0].Metadata)
	})

	t.Run("non-string metadata survives a re-marshal", func(t *testing.T) {
		_, mgClient, closeSrv := newManageTestServer(t, mixedMetadataBody)
		defer closeSrv()

		resp, err := mgClient.GetAgent(context.Background(), testProjectID, testAgentID)
		if err != nil {
			t.Fatalf("GetAgent failed on non-string metadata: %v", err)
		}
		out, err := json.Marshal(resp.Metadata)
		if err != nil {
			t.Fatalf("re-marshal failed: %v", err)
		}
		want := `{"count":7,"created_by":"console","flag":true,"nested":{"a":1}}`
		if string(out) != want {
			t.Errorf("expected re-marshaled metadata %s, got %s", want, out)
		}
	})

	t.Run("CreateAgent sends non-string metadata values to the wire", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, `{"agent_uuid":"`+testAgentID+`"}`)
		defer closeSrv()

		if _, err := mgClient.CreateAgent(context.Background(), testProjectID,
			&apiinterfaces.AgentCreateRequest{
				Config:   "{}",
				Metadata: map[string]interface{}{"count": 7, "flag": true},
			}); err != nil {
			t.Fatalf("CreateAgent failed: %v", err)
		}
		if want := `{"config":"{}","metadata":{"count":7,"flag":true}}`; string(recorded.Body) != want {
			t.Errorf("expected body %s, got %s", want, recorded.Body)
		}
	})

	t.Run("UpdateAgentMetadata sends non-string metadata values to the wire", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, "")
		defer closeSrv()

		if err := mgClient.UpdateAgentMetadata(context.Background(), testProjectID, testAgentID,
			&apiinterfaces.AgentMetadataUpdateRequest{
				Metadata: map[string]interface{}{"count": 7, "nested": map[string]interface{}{"a": 1}},
			}); err != nil {
			t.Fatalf("UpdateAgentMetadata failed: %v", err)
		}
		if want := `{"metadata":{"count":7,"nested":{"a":1}}}`; string(recorded.Body) != want {
			t.Errorf("expected body %s, got %s", want, recorded.Body)
		}
	})
}

// Test_AgentMetadataUpdateSendsEveryKey pins the SDK half of a destructive
// contract. PUT .../agents/{id} REPLACES the whole metadata object server-side:
// a key omitted from the request is deleted, silently, with a 200 (verified
// against production 2026-09-24). The SDK therefore must put exactly the map the
// caller handed it on the wire — dropping, filtering, or merging a key here
// would delete stored metadata the caller never asked to remove.
func Test_AgentMetadataUpdateSendsEveryKey(t *testing.T) {
	metadata := map[string]interface{}{
		"created_by": "go-sdk",
		"env":        "production",
		"revision":   2,
		"enabled":    true,
		"empty":      "",
		"zero":       0,
		"disabled":   false,
		"nested":     map[string]interface{}{"a": 1},
	}

	recorded, mgClient, closeSrv := newManageTestServer(t, "")
	defer closeSrv()

	if err := mgClient.UpdateAgentMetadata(context.Background(), testProjectID, testAgentID,
		&apiinterfaces.AgentMetadataUpdateRequest{Metadata: metadata}); err != nil {
		t.Fatalf("UpdateAgentMetadata failed: %v", err)
	}

	// Byte-exact: Go sorts map keys when marshaling, so this is deterministic.
	// Zero values are included on purpose — "" / 0 / false are meaningful
	// metadata, and an omitempty anywhere on this path would delete them.
	want := `{"metadata":{"created_by":"go-sdk","disabled":false,"empty":"","enabled":true,` +
		`"env":"production","nested":{"a":1},"revision":2,"zero":0}}`
	if string(recorded.Body) != want {
		t.Errorf("expected body %s, got %s", want, recorded.Body)
	}

	// And structurally: every key the caller passed reaches the wire, and the
	// SDK adds none of its own.
	var sent struct {
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := json.Unmarshal(recorded.Body, &sent); err != nil {
		t.Fatalf("failed to decode the recorded request body: %v", err)
	}
	if len(sent.Metadata) != len(metadata) {
		t.Errorf("expected %d metadata keys on the wire, got %d: %v",
			len(metadata), len(sent.Metadata), sent.Metadata)
	}
	for key := range metadata {
		if _, present := sent.Metadata[key]; !present {
			t.Errorf("metadata key %q was dropped on the way to the wire", key)
		}
	}
}

// Test_AgentVariableZeroValuesReachTheWire pins the serialization of values whose
// Go zero value is meaningful. A variable legitimately holds false, 0, or "", and
// an omitempty on Value would silently drop each of them, turning "set this
// variable to false" into "send no value at all".
func Test_AgentVariableZeroValuesReachTheWire(t *testing.T) {
	cases := []struct {
		name       string
		value      interface{}
		wantCreate string
		wantUpdate string
	}{
		{"false", false, `{"key":"DG_K","value":false,"is_sensitive":false}`, `{"value":false}`},
		{"zero", 0, `{"key":"DG_K","value":0,"is_sensitive":false}`, `{"value":0}`},
		{"empty string", "", `{"key":"DG_K","value":"","is_sensitive":false}`, `{"value":""}`},
		{"null", nil, `{"key":"DG_K","value":null,"is_sensitive":false}`, `{"value":null}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorded, mgClient, closeSrv := newManageTestServer(t, `{"agent_variable_uuid":"`+testVariableID+`"}`)
			defer closeSrv()
			if _, err := mgClient.CreateAgentVariable(context.Background(), testProjectID,
				&apiinterfaces.AgentVariableCreateRequest{Key: "DG_K", Value: tc.value}); err != nil {
				t.Fatalf("CreateAgentVariable failed: %v", err)
			}
			if string(recorded.Body) != tc.wantCreate {
				t.Errorf("create: expected body %s, got %s", tc.wantCreate, recorded.Body)
			}

			recordedUpd, mgClientUpd, closeSrvUpd := newManageTestServer(t, "")
			defer closeSrvUpd()
			if err := mgClientUpd.UpdateAgentVariable(context.Background(), testProjectID, testVariableID,
				&apiinterfaces.AgentVariableUpdateRequest{Value: tc.value}); err != nil {
				t.Fatalf("UpdateAgentVariable failed: %v", err)
			}
			if string(recordedUpd.Body) != tc.wantUpdate {
				t.Errorf("update: expected body %s, got %s", tc.wantUpdate, recordedUpd.Body)
			}
		})
	}

	t.Run("a decoded false value survives a re-marshal", func(t *testing.T) {
		var variable apiinterfaces.AgentVariable
		if err := json.Unmarshal([]byte(`{"agent_variable_uuid":"`+testVariableID+`","key":"DG_K","value":false,"is_sensitive":false}`), &variable); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		out, err := json.Marshal(variable)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatalf("re-marshaled body is not valid JSON: %v", err)
		}
		value, present := got["value"]
		if !present {
			t.Fatalf("value false was dropped on re-marshal: %s", out)
		}
		if value != false {
			t.Errorf("expected value false, got %v", value)
		}
	})
}

// Test_AgentRequestMarshalFailures proves the json.Marshal error paths return the
// error to the caller and never issue an HTTP request with a half-built body.
func Test_AgentRequestMarshalFailures(t *testing.T) {
	// A channel cannot be marshaled to JSON, so it drives every Marshal error path.
	unmarshalable := make(chan int)

	t.Run("CreateAgentVariable returns the marshal error", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, `{}`)
		defer closeSrv()

		resp, err := mgClient.CreateAgentVariable(context.Background(), testProjectID,
			&apiinterfaces.AgentVariableCreateRequest{Key: "DG_K", Value: unmarshalable})
		if err == nil {
			t.Fatal("expected a marshal error, got nil")
		}
		if resp != nil {
			t.Errorf("expected a nil result alongside the error, got %+v", resp)
		}
		if recorded.Method != "" {
			t.Errorf("no HTTP request should have been sent, but got %s %s", recorded.Method, recorded.Path)
		}
	})

	t.Run("UpdateAgentVariable returns the marshal error", func(t *testing.T) {
		recorded, mgClient, closeSrv := newManageTestServer(t, `{}`)
		defer closeSrv()

		err := mgClient.UpdateAgentVariable(context.Background(), testProjectID, testVariableID,
			&apiinterfaces.AgentVariableUpdateRequest{Value: unmarshalable})
		if err == nil {
			t.Fatal("expected a marshal error, got nil")
		}
		if recorded.Method != "" {
			t.Errorf("no HTTP request should have been sent, but got %s %s", recorded.Method, recorded.Path)
		}
	})
}
