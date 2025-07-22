// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"reflect"
	"testing"

	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func TestAgentTags_StructCreation(t *testing.T) {
	t.Run("Test Agent struct creation with tags field", func(t *testing.T) {
		// Test creating agent with tags
		tags := []string{"tag1", "tag2", "production"}
		agent := &interfacesv1.Agent{
			Language: "en",
			Tags:     tags,
		}

		// Verify the field is set correctly
		if !reflect.DeepEqual(agent.Tags, tags) {
			t.Errorf("Expected Tags to be %v, got %v", tags, agent.Tags)
		}

		// Test creating agent with single tag
		singleTag := []string{"test"}
		agent2 := &interfacesv1.Agent{
			Language: "en",
			Tags:     singleTag,
		}

		// Verify the single tag is set correctly
		if !reflect.DeepEqual(agent2.Tags, singleTag) {
			t.Errorf("Expected Tags to be %v, got %v", singleTag, agent2.Tags)
		}
	})

	t.Run("Test Agent struct with empty tags array", func(t *testing.T) {
		// Test creating agent with empty tags array
		agent := &interfacesv1.Agent{
			Language: "en",
			Tags:     []string{},
		}

		// Verify the field is an empty array
		if len(agent.Tags) != 0 {
			t.Errorf("Expected Tags to be empty array, got %v", agent.Tags)
		}
	})

	t.Run("Test Agent struct with default tags value", func(t *testing.T) {
		// Test creating agent without explicitly setting tags
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		// Verify the field defaults to nil (Go's zero value for slice)
		if agent.Tags != nil {
			t.Errorf("Expected Tags to default to nil, got %v", agent.Tags)
		}
	})
}

func TestAgentTags_JSONMarshaling(t *testing.T) {
	t.Run("Test Agent JSON marshaling with tags populated", func(t *testing.T) {
		tags := []string{"development", "test", "agent-v1"}
		agent := &interfacesv1.Agent{
			Language: "en",
			Tags:     tags,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify tags field is present and correct
		tagsInterface, exists := result["tags"]
		if !exists {
			t.Error("Expected tags field to be present in JSON")
		}

		// Convert interface{} to []interface{} and then to []string for comparison
		tagsArray, ok := tagsInterface.([]interface{})
		if !ok {
			t.Errorf("Expected tags to be array, got %T", tagsInterface)
		}

		// Convert []interface{} to []string
		var tagsStrings []string
		for _, tag := range tagsArray {
			tagStr, ok := tag.(string)
			if !ok {
				t.Errorf("Expected tag to be string, got %T", tag)
			}
			tagsStrings = append(tagsStrings, tagStr)
		}

		if !reflect.DeepEqual(tagsStrings, tags) {
			t.Errorf("Expected tags to be %v, got %v", tags, tagsStrings)
		}
	})

	t.Run("Test Agent JSON marshaling with empty tags array", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Tags:     []string{},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// With omitempty tag, empty array should be omitted from JSON
		if _, exists := result["tags"]; exists {
			t.Errorf("Expected tags to be omitted from JSON when empty array, but it was present with value %v", result["tags"])
		}
	})

	t.Run("Test Agent JSON marshaling with nil tags", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Tags:     nil,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// With omitempty tag, nil value should be omitted from JSON
		if _, exists := result["tags"]; exists {
			t.Errorf("Expected tags to be omitted from JSON when nil, but it was present with value %v", result["tags"])
		}
	})

	t.Run("Test Agent JSON marshaling with single tag", func(t *testing.T) {
		tags := []string{"production"}
		agent := &interfacesv1.Agent{
			Language: "en",
			Tags:     tags,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify single tag is correctly marshaled
		tagsInterface, exists := result["tags"]
		if !exists {
			t.Error("Expected tags field to be present in JSON")
		}

		tagsArray, ok := tagsInterface.([]interface{})
		if !ok {
			t.Errorf("Expected tags to be array, got %T", tagsInterface)
		}

		if len(tagsArray) != 1 {
			t.Errorf("Expected tags array to have 1 element, got %d", len(tagsArray))
		}

		if tagsArray[0] != "production" {
			t.Errorf("Expected single tag to be 'production', got %v", tagsArray[0])
		}
	})
}

func TestAgentTags_JSONUnmarshaling(t *testing.T) {
	t.Run("Test Agent JSON unmarshaling with tags", func(t *testing.T) {
		jsonStr := `{"language":"en","tags":["development","test","agent-v1"]}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonStr), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		expectedTags := []string{"development", "test", "agent-v1"}
		if !reflect.DeepEqual(agent.Tags, expectedTags) {
			t.Errorf("Expected Tags to be %v, got %v", expectedTags, agent.Tags)
		}
	})

	t.Run("Test Agent JSON unmarshaling without tags field", func(t *testing.T) {
		jsonStr := `{"language":"en"}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonStr), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		// Tags should be nil when not present in JSON
		if agent.Tags != nil {
			t.Errorf("Expected Tags to be nil when not in JSON, got %v", agent.Tags)
		}
	})

	t.Run("Test Agent JSON unmarshaling with empty tags array", func(t *testing.T) {
		jsonStr := `{"language":"en","tags":[]}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonStr), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		// Tags should be empty slice when empty array in JSON
		if agent.Tags == nil {
			t.Error("Expected Tags to be empty slice, got nil")
		}
		if len(agent.Tags) != 0 {
			t.Errorf("Expected Tags to be empty slice, got %v", agent.Tags)
		}
	})
}

func TestAgentTags_SettingsOptions(t *testing.T) {
	t.Run("Test SettingsOptions with tags in Agent", func(t *testing.T) {
		options := interfacesv1.NewSettingsOptions()

		// Add tags to the agent
		options.Agent.Tags = []string{"test", "schema-validation"}

		// Verify tags are set correctly
		expectedTags := []string{"test", "schema-validation"}
		if !reflect.DeepEqual(options.Agent.Tags, expectedTags) {
			t.Errorf("Expected Agent.Tags to be %v, got %v", expectedTags, options.Agent.Tags)
		}

		// Test JSON marshaling of complete SettingsOptions
		jsonData, err := json.Marshal(options)
		if err != nil {
			t.Fatalf("JSON marshaling of SettingsOptions failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal SettingsOptions JSON: %v", err)
		}

		// Navigate to agent.tags in the JSON structure
		agentInterface, exists := result["agent"]
		if !exists {
			t.Error("Expected agent field to be present in SettingsOptions JSON")
		}

		agentObj, ok := agentInterface.(map[string]interface{})
		if !ok {
			t.Errorf("Expected agent to be object, got %T", agentInterface)
		}

		tagsInterface, exists := agentObj["tags"]
		if !exists {
			t.Error("Expected tags field to be present in agent JSON")
		}

		tagsArray, ok := tagsInterface.([]interface{})
		if !ok {
			t.Errorf("Expected tags to be array, got %T", tagsInterface)
		}

		// Verify tag values
		if len(tagsArray) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tagsArray))
		}
		if tagsArray[0] != "test" || tagsArray[1] != "schema-validation" {
			t.Errorf("Expected tags ['test', 'schema-validation'], got %v", tagsArray)
		}
	})
}
