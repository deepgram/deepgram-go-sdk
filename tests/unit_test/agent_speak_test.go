// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"reflect"
	"testing"

	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func TestAgentSpeakSingleProvider(t *testing.T) {
	t.Run("Test single Speak provider assignment", func(t *testing.T) {
		// Create agent with single speak provider
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		// Test assignment of single Speak struct
		singleSpeak := interfacesv1.Speak{
			Provider: map[string]interface{}{
				"type":  "deepgram",
				"model": "aura-2-thalia-en",
			},
		}

		agent.Speak = singleSpeak

		// Verify the assignment worked
		if agent.Speak == nil {
			t.Error("Agent.Speak should not be nil after assignment")
		}

		// Type assertion to verify it's still a Speak struct
		speak, ok := agent.Speak.(interfacesv1.Speak)
		if !ok {
			t.Error("Agent.Speak should be of type Speak")
		}

		// Verify the provider data
		if speak.Provider["type"] != "deepgram" {
			t.Errorf("Expected provider type 'deepgram', got %v", speak.Provider["type"])
		}
		if speak.Provider["model"] != "aura-2-thalia-en" {
			t.Errorf("Expected model 'aura-2-thalia-en', got %v", speak.Provider["model"])
		}
	})

	t.Run("Test single Speak provider JSON marshaling", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Speak: interfacesv1.Speak{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-zeus-en",
				},
			},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(agent)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}

		// Verify JSON structure
		expectedJSON := `{"language":"en","listen":{"provider":null},"think":{"provider":null},"speak":{"provider":{"model":"aura-2-zeus-en","type":"deepgram"}}}`

		// Parse both JSONs to compare structured data
		var expected, actual map[string]interface{}
		if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
			t.Fatalf("Failed to unmarshal expected JSON: %v", err)
		}
		if err := json.Unmarshal(jsonData, &actual); err != nil {
			t.Fatalf("Failed to unmarshal actual JSON: %v", err)
		}

		// Compare speak sections
		expectedSpeak := expected["speak"].(map[string]interface{})
		actualSpeak := actual["speak"].(map[string]interface{})

		if !reflect.DeepEqual(expectedSpeak["provider"], actualSpeak["provider"]) {
			t.Errorf("JSON marshaling mismatch.\nExpected: %v\nActual: %v", expectedSpeak["provider"], actualSpeak["provider"])
		}
	})
}

func TestAgentSpeakMultipleProviders(t *testing.T) {
	t.Run("Test multiple Speak providers assignment", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		// Test assignment of multiple Speak providers
		multipleSpeak := []interfacesv1.Speak{
			{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-thalia-en",
				},
			},
			{
				Provider: map[string]interface{}{
					"type":  "open_ai",
					"model": "tts-1",
					"voice": "shimmer",
				},
				Endpoint: &interfacesv1.Endpoint{
					Url: "https://api.openai.com/v1/audio/speech",
					Headers: map[string]string{
						"authorization": "Bearer {{OPENAI_API_KEY}}",
					},
				},
			},
		}

		agent.Speak = multipleSpeak

		// Verify the assignment worked
		if agent.Speak == nil {
			t.Error("Agent.Speak should not be nil after assignment")
		}

		// Type assertion to verify it's a slice of Speak
		speakSlice, ok := agent.Speak.([]interfacesv1.Speak)
		if !ok {
			t.Error("Agent.Speak should be of type []Speak")
		}

		// Verify we have 2 providers
		if len(speakSlice) != 2 {
			t.Errorf("Expected 2 providers, got %d", len(speakSlice))
		}

		// Verify first provider (Deepgram)
		if speakSlice[0].Provider["type"] != "deepgram" {
			t.Errorf("Expected first provider type 'deepgram', got %v", speakSlice[0].Provider["type"])
		}

		// Verify second provider (OpenAI)
		if speakSlice[1].Provider["type"] != "open_ai" {
			t.Errorf("Expected second provider type 'open_ai', got %v", speakSlice[1].Provider["type"])
		}
		if speakSlice[1].Endpoint == nil {
			t.Error("Expected second provider to have endpoint")
		}
	})

	t.Run("Test multiple Speak providers JSON marshaling", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Speak: []interfacesv1.Speak{
				{
					Provider: map[string]interface{}{
						"type":  "deepgram",
						"model": "aura-2-thalia-en",
					},
				},
				{
					Provider: map[string]interface{}{
						"type":  "open_ai",
						"model": "tts-1",
						"voice": "shimmer",
					},
					Endpoint: &interfacesv1.Endpoint{
						Url: "https://api.openai.com/v1/audio/speech",
						Headers: map[string]string{
							"authorization": "Bearer {{OPENAI_API_KEY}}",
						},
					},
				},
			},
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

		// Verify speak is an array
		speak, ok := result["speak"].([]interface{})
		if !ok {
			t.Error("speak field should be an array")
		}

		if len(speak) != 2 {
			t.Errorf("Expected 2 speak providers, got %d", len(speak))
		}

		// Verify first provider
		firstProvider := speak[0].(map[string]interface{})
		provider1 := firstProvider["provider"].(map[string]interface{})
		if provider1["type"] != "deepgram" {
			t.Errorf("Expected first provider type 'deepgram', got %v", provider1["type"])
		}

		// Verify second provider has endpoint
		secondProvider := speak[1].(map[string]interface{})
		if secondProvider["endpoint"] == nil {
			t.Error("Second provider should have endpoint")
		}
	})
}

func TestAgentSpeakJSONUnmarshaling(t *testing.T) {
	t.Run("Test unmarshaling single provider JSON", func(t *testing.T) {
		jsonData := `{
			"language": "en",
			"speak": {
				"provider": {
					"type": "deepgram",
					"model": "aura-2-thalia-en"
				}
			}
		}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonData), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		if agent.Speak == nil {
			t.Error("Agent.Speak should not be nil after unmarshaling")
		}

		// The JSON unmarshaling will create a map[string]interface{} for the speak field
		// We need to handle this in our application logic
		speakMap, ok := agent.Speak.(map[string]interface{})
		if !ok {
			t.Error("Unmarshaled speak should be map[string]interface{}")
		}

		provider := speakMap["provider"].(map[string]interface{})
		if provider["type"] != "deepgram" {
			t.Errorf("Expected provider type 'deepgram', got %v", provider["type"])
		}
	})

	t.Run("Test unmarshaling multiple providers JSON", func(t *testing.T) {
		jsonData := `{
			"language": "en",
			"speak": [
				{
					"provider": {
						"type": "deepgram",
						"model": "aura-2-thalia-en"
					}
				},
				{
					"provider": {
						"type": "open_ai",
						"model": "tts-1",
						"voice": "shimmer"
					},
					"endpoint": {
						"url": "https://api.openai.com/v1/audio/speech",
						"headers": {
							"authorization": "Bearer {{OPENAI_API_KEY}}"
						}
					}
				}
			]
		}`

		var agent interfacesv1.Agent
		err := json.Unmarshal([]byte(jsonData), &agent)
		if err != nil {
			t.Fatalf("JSON unmarshaling failed: %v", err)
		}

		if agent.Speak == nil {
			t.Error("Agent.Speak should not be nil after unmarshaling")
		}

		// The JSON unmarshaling will create a []interface{} for the speak field
		speakSlice, ok := agent.Speak.([]interface{})
		if !ok {
			t.Error("Unmarshaled speak should be []interface{}")
		}

		if len(speakSlice) != 2 {
			t.Errorf("Expected 2 providers, got %d", len(speakSlice))
		}

		// Check first provider
		firstProvider := speakSlice[0].(map[string]interface{})
		provider1 := firstProvider["provider"].(map[string]interface{})
		if provider1["type"] != "deepgram" {
			t.Errorf("Expected first provider type 'deepgram', got %v", provider1["type"])
		}

		// Check second provider
		secondProvider := speakSlice[1].(map[string]interface{})
		if secondProvider["endpoint"] == nil {
			t.Error("Second provider should have endpoint")
		}
	})
}

func TestAgentSpeakEdgeCases(t *testing.T) {
	t.Run("Test nil assignment", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		agent.Speak = nil

		if agent.Speak != nil {
			t.Error("Agent.Speak should be nil after nil assignment")
		}
	})

	t.Run("Test empty slice assignment", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		agent.Speak = []interfacesv1.Speak{}

		speakSlice, ok := agent.Speak.([]interfacesv1.Speak)
		if !ok {
			t.Error("Agent.Speak should be of type []Speak")
		}

		if len(speakSlice) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(speakSlice))
		}
	})

	t.Run("Test single provider with endpoint", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		agent.Speak = interfacesv1.Speak{
			Provider: map[string]interface{}{
				"type":  "open_ai",
				"model": "tts-1",
				"voice": "alloy",
			},
			Endpoint: &interfacesv1.Endpoint{
				Url: "https://api.openai.com/v1/audio/speech",
				Headers: map[string]string{
					"authorization": "Bearer test-key",
				},
			},
		}

		speak, ok := agent.Speak.(interfacesv1.Speak)
		if !ok {
			t.Error("Agent.Speak should be of type Speak")
		}

		if speak.Endpoint == nil {
			t.Error("Speak should have endpoint")
		}

		if speak.Endpoint.Url != "https://api.openai.com/v1/audio/speech" {
			t.Errorf("Expected URL 'https://api.openai.com/v1/audio/speech', got %v", speak.Endpoint.Url)
		}
	})
}

func TestAgentSpeakBackwardCompatibility(t *testing.T) {
	t.Run("Test backward compatibility with existing code patterns", func(t *testing.T) {
		// This test ensures that existing code that assigns Speak struct still works
		agent := &interfacesv1.Agent{}

		// Old pattern - direct assignment
		agent.Speak = interfacesv1.Speak{
			Provider: map[string]interface{}{
				"type":  "deepgram",
				"model": "aura-2-thalia-en",
			},
		}

		// Verify it still works
		speak, ok := agent.Speak.(interfacesv1.Speak)
		if !ok {
			t.Error("Backward compatibility test failed: Agent.Speak should accept Speak struct")
		}

		if speak.Provider["type"] != "deepgram" {
			t.Error("Backward compatibility test failed: Provider data should be preserved")
		}
	})

	t.Run("Test NewSettingsOptions backward compatibility", func(t *testing.T) {
		// Test that the default options still work as expected
		options := interfacesv1.NewSettingsOptions()

		if options.Agent.Speak == nil {
			t.Error("NewSettingsOptions should set a default Speak provider")
		}

		// Check that it's a single Speak struct (backward compatible)
		speak, ok := options.Agent.Speak.(interfacesv1.Speak)
		if !ok {
			t.Error("Default Speak should be a single Speak struct for backward compatibility")
		}

		if speak.Provider["type"] != "deepgram" {
			t.Error("Default Speak provider should be deepgram")
		}
	})
}
