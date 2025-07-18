// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func TestAgentMipOptOut_StructCreation(t *testing.T) {
	t.Run("Test Agent struct creation with mip_opt_out field", func(t *testing.T) {
		// Test creating agent with mip_opt_out set to true
		agent := &interfacesv1.Agent{
			Language:  "en",
			MipOptOut: true,
		}

		// Verify the field is set correctly
		if agent.MipOptOut != true {
			t.Errorf("Expected MipOptOut to be true, got %v", agent.MipOptOut)
		}

		// Test creating agent with mip_opt_out set to false
		agent2 := &interfacesv1.Agent{
			Language:  "en",
			MipOptOut: false,
		}

		// Verify the field is set correctly
		if agent2.MipOptOut != false {
			t.Errorf("Expected MipOptOut to be false, got %v", agent2.MipOptOut)
		}
	})

	t.Run("Test Agent struct with default mip_opt_out value", func(t *testing.T) {
		// Test creating agent without explicitly setting mip_opt_out
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		// Verify the field defaults to false (Go's zero value for bool)
		if agent.MipOptOut != false {
			t.Errorf("Expected MipOptOut to default to false, got %v", agent.MipOptOut)
		}
	})
}

func TestAgentMipOptOut_JSONMarshaling(t *testing.T) {
	t.Run("Test Agent JSON marshaling with mip_opt_out set to true", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language:  "en",
			MipOptOut: true,
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

		// Verify mip_opt_out field is present and correct
		if result["mip_opt_out"] != true {
			t.Errorf("Expected mip_opt_out to be true, got %v", result["mip_opt_out"])
		}
	})

	t.Run("Test Agent JSON marshaling with mip_opt_out set to false", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language:  "en",
			MipOptOut: false,
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

		// With omitempty tag, false value should be omitted from JSON
		if _, exists := result["mip_opt_out"]; exists {
			t.Errorf("Expected mip_opt_out to be omitted from JSON when false, but it was present with value %v", result["mip_opt_out"])
		}
	})

	t.Run("Test Agent JSON marshaling with default mip_opt_out value", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
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

		// With omitempty tag, default false value should be omitted from JSON
		if _, exists := result["mip_opt_out"]; exists {
			t.Errorf("Expected mip_opt_out to be omitted from JSON when using default value, but it was present with value %v", result["mip_opt_out"])
		}
	})
}

func TestAgentMipOptOut_JSONUnmarshaling(t *testing.T) {
	t.Run("Test Agent JSON unmarshaling with mip_opt_out set to true", func(t *testing.T) {
		jsonData := `{
			"language": "en",
			"mip_opt_out": true
		}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonData), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		// Verify the field is set correctly
		if agent.MipOptOut != true {
			t.Errorf("Expected MipOptOut to be true, got %v", agent.MipOptOut)
		}
	})

	t.Run("Test Agent JSON unmarshaling with mip_opt_out set to false", func(t *testing.T) {
		jsonData := `{
			"language": "en",
			"mip_opt_out": false
		}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonData), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		// Verify the field is set correctly
		if agent.MipOptOut != false {
			t.Errorf("Expected MipOptOut to be false, got %v", agent.MipOptOut)
		}
	})

	t.Run("Test Agent JSON unmarshaling without mip_opt_out field", func(t *testing.T) {
		jsonData := `{
			"language": "en"
		}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonData), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		// Verify the field defaults to false when not present in JSON
		if agent.MipOptOut != false {
			t.Errorf("Expected MipOptOut to default to false, got %v", agent.MipOptOut)
		}
	})
}

func TestAgentMipOptOut_BackwardCompatibility(t *testing.T) {
	t.Run("Test backward compatibility with existing Agent usage", func(t *testing.T) {
		// This test ensures existing code patterns still work
		agent := &interfacesv1.Agent{
			Language: "en",
			Listen: interfacesv1.Listen{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "nova-3",
				},
			},
			Think: interfacesv1.Think{
				Provider: map[string]interface{}{
					"type":  "open_ai",
					"model": "gpt-4o-mini",
				},
			},
			Speak: interfacesv1.Speak{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-thalia-en",
				},
			},
			Greeting: "Hello!",
		}

		// Verify that existing functionality still works
		if agent.Language != "en" {
			t.Errorf("Expected language to be 'en', got %v", agent.Language)
		}

		if agent.Greeting != "Hello!" {
			t.Errorf("Expected greeting to be 'Hello!', got %v", agent.Greeting)
		}

		// Verify that mip_opt_out defaults to false
		if agent.MipOptOut != false {
			t.Errorf("Expected MipOptOut to default to false, got %v", agent.MipOptOut)
		}

		// Test JSON marshaling with existing fields
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify existing fields are present
		if result["language"] != "en" {
			t.Errorf("Expected language to be 'en', got %v", result["language"])
		}

		if result["greeting"] != "Hello!" {
			t.Errorf("Expected greeting to be 'Hello!', got %v", result["greeting"])
		}

		// Verify mip_opt_out is omitted when false (default)
		if _, exists := result["mip_opt_out"]; exists {
			t.Errorf("Expected mip_opt_out to be omitted from JSON when false, but it was present")
		}
	})
}

func TestAgentMipOptOut_SettingsOptions(t *testing.T) {
	t.Run("Test mip_opt_out field in SettingsOptions", func(t *testing.T) {
		// Test creating SettingsOptions with mip_opt_out
		options := &interfacesv1.SettingsOptions{
			Type: "Settings",
			Agent: interfacesv1.Agent{
				Language:  "en",
				MipOptOut: true,
			},
		}

		// Verify the field is set correctly
		if options.Agent.MipOptOut != true {
			t.Errorf("Expected Agent.MipOptOut to be true, got %v", options.Agent.MipOptOut)
		}

		// Test JSON marshaling
		jsonData, err := json.Marshal(options)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify nested structure
		agent, ok := result["agent"].(map[string]interface{})
		if !ok {
			t.Error("agent field should be an object")
		}

		if agent["mip_opt_out"] != true {
			t.Errorf("Expected agent.mip_opt_out to be true, got %v", agent["mip_opt_out"])
		}
	})

	t.Run("Test NewSettingsOptions with default mip_opt_out", func(t *testing.T) {
		// Test that NewSettingsOptions creates proper defaults
		options := interfacesv1.NewSettingsOptions()

		// Verify the field defaults to false
		if options.Agent.MipOptOut != false {
			t.Errorf("Expected Agent.MipOptOut to default to false, got %v", options.Agent.MipOptOut)
		}
	})
}

func TestAgentMipOptOut_EdgeCases(t *testing.T) {
	t.Run("Test mip_opt_out with all other fields", func(t *testing.T) {
		// Test creating agent with all fields including mip_opt_out
		agent := &interfacesv1.Agent{
			Language: "en",
			Listen: interfacesv1.Listen{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "nova-3",
				},
			},
			Think: interfacesv1.Think{
				Provider: map[string]interface{}{
					"type":  "open_ai",
					"model": "gpt-4o-mini",
				},
			},
			Speak: interfacesv1.Speak{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-thalia-en",
				},
			},
			SpeakFallback: &[]interfacesv1.Speak{
				{
					Provider: map[string]interface{}{
						"type":  "open_ai",
						"model": "tts-1",
					},
				},
			},
			Greeting:  "Hello!",
			MipOptOut: true,
		}

		// Verify all fields are set correctly
		if agent.Language != "en" {
			t.Errorf("Expected language to be 'en', got %v", agent.Language)
		}

		if agent.MipOptOut != true {
			t.Errorf("Expected MipOptOut to be true, got %v", agent.MipOptOut)
		}

		// Test JSON marshaling with all fields
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Parse JSON to verify structure
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify mip_opt_out is present when true
		if result["mip_opt_out"] != true {
			t.Errorf("Expected mip_opt_out to be true, got %v", result["mip_opt_out"])
		}

		// Verify other fields are still present
		if result["language"] != "en" {
			t.Errorf("Expected language to be 'en', got %v", result["language"])
		}

		if result["greeting"] != "Hello!" {
			t.Errorf("Expected greeting to be 'Hello!', got %v", result["greeting"])
		}
	})
}
