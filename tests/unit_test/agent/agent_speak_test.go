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

func TestAgentSpeakBackwardCompatibility(t *testing.T) {
	t.Run("Test backward compatibility - direct field access", func(t *testing.T) {
		// This test ensures that existing code patterns still work without type assertions
		agent := &interfacesv1.Agent{
			Language: "en",
		}

		// Old pattern - direct assignment (this should still work)
		agent.Speak = interfacesv1.Speak{
			Provider: map[string]interface{}{
				"type":  "deepgram",
				"model": "aura-2-thalia-en",
			},
		}

		// Old pattern - direct field access (this should still work)
		if agent.Speak.Provider["type"] != "deepgram" {
			t.Errorf("Expected provider type 'deepgram', got %v", agent.Speak.Provider["type"])
		}
		if agent.Speak.Provider["model"] != "aura-2-thalia-en" {
			t.Errorf("Expected model 'aura-2-thalia-en', got %v", agent.Speak.Provider["model"])
		}
	})

	t.Run("Test backward compatibility - JSON marshaling", func(t *testing.T) {
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

	t.Run("Test NewSettingsOptions backward compatibility", func(t *testing.T) {
		// Test that the default options still work as expected
		options := interfacesv1.NewSettingsOptions()

		// Check that the primary speak provider is properly set
		if options.Agent.Speak.Provider == nil {
			t.Error("NewSettingsOptions should set a default Speak provider")
		}

		// Direct field access should work (no type assertions needed)
		if options.Agent.Speak.Provider["type"] != "deepgram" {
			t.Errorf("Expected default provider type 'deepgram', got %v", options.Agent.Speak.Provider["type"])
		}
		if options.Agent.Speak.Provider["model"] != "aura-2-thalia-en" {
			t.Errorf("Expected default model 'aura-2-thalia-en', got %v", options.Agent.Speak.Provider["model"])
		}
	})
}

func TestAgentSpeakFallbackProviders(t *testing.T) {
	t.Run("Test fallback providers assignment", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			// Primary provider (backward compatible)
			Speak: interfacesv1.Speak{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-thalia-en",
				},
			},
		}

		// Test assignment of fallback providers
		fallbackProviders := []interfacesv1.Speak{
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
			{
				Provider: map[string]interface{}{
					"type":  "elevenlabs",
					"model": "eleven_turbo_v2",
					"voice": "alice",
				},
				Endpoint: &interfacesv1.Endpoint{
					Url: "https://api.elevenlabs.io/v1/text-to-speech",
					Headers: map[string]string{
						"authorization": "Bearer {{ELEVENLABS_API_KEY}}",
					},
				},
			},
		}

		agent.SpeakFallback = &fallbackProviders

		// Verify primary provider still works
		if agent.Speak.Provider["type"] != "deepgram" {
			t.Errorf("Expected primary provider type 'deepgram', got %v", agent.Speak.Provider["type"])
		}

		// Verify fallback providers are set
		if agent.SpeakFallback == nil {
			t.Error("SpeakFallback should not be nil after assignment")
		}

		if len(*agent.SpeakFallback) != 2 {
			t.Errorf("Expected 2 fallback providers, got %d", len(*agent.SpeakFallback))
		}

		// Verify first fallback provider (OpenAI)
		if (*agent.SpeakFallback)[0].Provider["type"] != "open_ai" {
			t.Errorf("Expected first fallback provider type 'open_ai', got %v", (*agent.SpeakFallback)[0].Provider["type"])
		}

		// Verify second fallback provider (ElevenLabs)
		if (*agent.SpeakFallback)[1].Provider["type"] != "elevenlabs" {
			t.Errorf("Expected second fallback provider type 'elevenlabs', got %v", (*agent.SpeakFallback)[1].Provider["type"])
		}
	})

	t.Run("Test fallback providers JSON marshaling", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
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

		// Verify primary speak provider is present
		speak, ok := result["speak"].(map[string]interface{})
		if !ok {
			t.Error("speak field should be an object")
		}

		provider := speak["provider"].(map[string]interface{})
		if provider["type"] != "deepgram" {
			t.Errorf("Expected primary provider type 'deepgram', got %v", provider["type"])
		}

		// Verify fallback providers are present
		speakFallback, ok := result["speak_fallback"].([]interface{})
		if !ok {
			t.Error("speak_fallback field should be an array")
		}

		if len(speakFallback) != 1 {
			t.Errorf("Expected 1 fallback provider, got %d", len(speakFallback))
		}

		// Verify fallback provider has endpoint
		firstFallback := speakFallback[0].(map[string]interface{})
		if firstFallback["endpoint"] == nil {
			t.Error("Fallback provider should have endpoint")
		}
	})
}

func TestAgentSpeakJSONUnmarshaling(t *testing.T) {
	t.Run("Test unmarshaling backward compatible JSON", func(t *testing.T) {
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

		// Direct field access should work (no type assertions needed)
		if agent.Speak.Provider["type"] != "deepgram" {
			t.Errorf("Expected provider type 'deepgram', got %v", agent.Speak.Provider["type"])
		}
		if agent.Speak.Provider["model"] != "aura-2-thalia-en" {
			t.Errorf("Expected model 'aura-2-thalia-en', got %v", agent.Speak.Provider["model"])
		}
	})

	t.Run("Test unmarshaling JSON with fallback providers", func(t *testing.T) {
		jsonData := `{
			"language": "en",
			"speak": {
				"provider": {
					"type": "deepgram",
					"model": "aura-2-thalia-en"
				}
			},
			"speak_fallback": [
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

		// Verify primary provider
		if agent.Speak.Provider["type"] != "deepgram" {
			t.Errorf("Expected primary provider type 'deepgram', got %v", agent.Speak.Provider["type"])
		}

		// Verify fallback providers
		if agent.SpeakFallback == nil {
			t.Error("SpeakFallback should not be nil after unmarshaling")
		}

		if len(*agent.SpeakFallback) != 1 {
			t.Errorf("Expected 1 fallback provider, got %d", len(*agent.SpeakFallback))
		}

		// Check fallback provider
		fallback := (*agent.SpeakFallback)[0]
		if fallback.Provider["type"] != "open_ai" {
			t.Errorf("Expected fallback provider type 'open_ai', got %v", fallback.Provider["type"])
		}

		if fallback.Endpoint == nil {
			t.Error("Fallback provider should have endpoint")
		}
	})
}

func TestAgentSpeakEdgeCases(t *testing.T) {
	t.Run("Test nil fallback providers", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Speak: interfacesv1.Speak{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-thalia-en",
				},
			},
			SpeakFallback: nil,
		}

		// Primary provider should still work
		if agent.Speak.Provider["type"] != "deepgram" {
			t.Error("Primary provider should work even with nil fallback")
		}

		if agent.SpeakFallback != nil {
			t.Error("SpeakFallback should be nil when not set")
		}
	})

	t.Run("Test empty fallback providers", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Speak: interfacesv1.Speak{
				Provider: map[string]interface{}{
					"type":  "deepgram",
					"model": "aura-2-thalia-en",
				},
			},
			SpeakFallback: &[]interfacesv1.Speak{},
		}

		// Primary provider should still work
		if agent.Speak.Provider["type"] != "deepgram" {
			t.Error("Primary provider should work with empty fallback")
		}

		if len(*agent.SpeakFallback) != 0 {
			t.Errorf("Expected empty fallback slice, got length %d", len(*agent.SpeakFallback))
		}
	})

	t.Run("Test primary provider with endpoint", func(t *testing.T) {
		agent := &interfacesv1.Agent{
			Language: "en",
			Speak: interfacesv1.Speak{
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
			},
		}

		// Direct field access should work
		if agent.Speak.Provider["type"] != "open_ai" {
			t.Errorf("Expected provider type 'open_ai', got %v", agent.Speak.Provider["type"])
		}

		if agent.Speak.Endpoint == nil {
			t.Error("Speak should have endpoint")
		}

		if agent.Speak.Endpoint.Url != "https://api.openai.com/v1/audio/speech" {
			t.Errorf("Expected URL 'https://api.openai.com/v1/audio/speech', got %v", agent.Speak.Endpoint.Url)
		}
	})
}
