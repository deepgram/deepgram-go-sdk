// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func TestSettingsMipOptOut_StructCreation(t *testing.T) {
	t.Run("Test SettingsOptions struct creation with mip_opt_out field", func(t *testing.T) {
		// Test creating settings with mip_opt_out set to true
		settings := &interfacesv1.SettingsOptions{
			Type:      "Settings",
			MipOptOut: true,
		}

		// Verify the field is set correctly
		if settings.MipOptOut != true {
			t.Errorf("Expected MipOptOut to be true, got %v", settings.MipOptOut)
		}

		// Test creating settings with mip_opt_out set to false
		settings2 := &interfacesv1.SettingsOptions{
			Type:      "Settings",
			MipOptOut: false,
		}

		// Verify the field is set correctly
		if settings2.MipOptOut != false {
			t.Errorf("Expected MipOptOut to be false, got %v", settings2.MipOptOut)
		}
	})

	t.Run("Test SettingsOptions struct with default mip_opt_out value", func(t *testing.T) {
		// Test creating settings without explicitly setting mip_opt_out
		settings := &interfacesv1.SettingsOptions{
			Type: "Settings",
		}

		// Verify the field defaults to false (Go's zero value for bool)
		if settings.MipOptOut != false {
			t.Errorf("Expected MipOptOut to default to false, got %v", settings.MipOptOut)
		}
	})
}

func TestSettingsMipOptOut_JSONMarshaling(t *testing.T) {
	t.Run("Test SettingsOptions JSON marshaling with mip_opt_out set to true", func(t *testing.T) {
		// Create settings with mip_opt_out set to true
		settings := &interfacesv1.SettingsOptions{
			Type:      "Settings",
			MipOptOut: true,
			Audio: interfacesv1.Audio{
				Input: &interfacesv1.Input{
					Encoding:   "linear16",
					SampleRate: 16000,
				},
			},
			Agent: interfacesv1.Agent{
				Language: "en",
			},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(settings)
		if err != nil {
			t.Fatalf("Failed to marshal settings to JSON: %v", err)
		}

		// Parse back to verify
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify mip_opt_out field is present and correct
		if result["mip_opt_out"] != true {
			t.Errorf("Expected mip_opt_out to be true, got %v", result["mip_opt_out"])
		}
	})

	t.Run("Test SettingsOptions JSON marshaling with mip_opt_out set to false", func(t *testing.T) {
		// Create settings with mip_opt_out set to false
		settings := &interfacesv1.SettingsOptions{
			Type:      "Settings",
			MipOptOut: false,
			Audio: interfacesv1.Audio{
				Input: &interfacesv1.Input{
					Encoding:   "linear16",
					SampleRate: 16000,
				},
			},
			Agent: interfacesv1.Agent{
				Language: "en",
			},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(settings)
		if err != nil {
			t.Fatalf("Failed to marshal settings to JSON: %v", err)
		}

		// Parse back to verify
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// With omitempty tag, mip_opt_out should be omitted when false
		if _, exists := result["mip_opt_out"]; exists {
			t.Errorf("Expected mip_opt_out to be omitted from JSON when false, but it was present with value %v", result["mip_opt_out"])
		}
	})

	t.Run("Test SettingsOptions JSON marshaling with default mip_opt_out value", func(t *testing.T) {
		// Create settings without explicitly setting mip_opt_out
		settings := &interfacesv1.SettingsOptions{
			Type: "Settings",
			Audio: interfacesv1.Audio{
				Input: &interfacesv1.Input{
					Encoding:   "linear16",
					SampleRate: 16000,
				},
			},
			Agent: interfacesv1.Agent{
				Language: "en",
			},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(settings)
		if err != nil {
			t.Fatalf("Failed to marshal settings to JSON: %v", err)
		}

		// Parse back to verify
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// With omitempty tag and default false value, mip_opt_out should be omitted
		if _, exists := result["mip_opt_out"]; exists {
			t.Errorf("Expected mip_opt_out to be omitted from JSON when using default value, but it was present with value %v", result["mip_opt_out"])
		}
	})
}

func TestSettingsMipOptOut_JSONUnmarshaling(t *testing.T) {
	t.Run("Test SettingsOptions JSON unmarshaling with mip_opt_out set to true", func(t *testing.T) {
		// JSON with mip_opt_out set to true
		jsonData := `{
			"type": "Settings",
			"mip_opt_out": true,
			"audio": {
				"input": {
					"encoding": "linear16",
					"sample_rate": 16000
				}
			},
			"agent": {
				"language": "en"
			}
		}`

		var settings interfacesv1.SettingsOptions
		err := json.Unmarshal([]byte(jsonData), &settings)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify the field is set correctly
		if settings.MipOptOut != true {
			t.Errorf("Expected MipOptOut to be true, got %v", settings.MipOptOut)
		}
	})

	t.Run("Test SettingsOptions JSON unmarshaling with mip_opt_out set to false", func(t *testing.T) {
		// JSON with mip_opt_out set to false
		jsonData := `{
			"type": "Settings",
			"mip_opt_out": false,
			"audio": {
				"input": {
					"encoding": "linear16",
					"sample_rate": 16000
				}
			},
			"agent": {
				"language": "en"
			}
		}`

		var settings interfacesv1.SettingsOptions
		err := json.Unmarshal([]byte(jsonData), &settings)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify the field is set correctly
		if settings.MipOptOut != false {
			t.Errorf("Expected MipOptOut to be false, got %v", settings.MipOptOut)
		}
	})

	t.Run("Test SettingsOptions JSON unmarshaling without mip_opt_out field", func(t *testing.T) {
		// JSON without mip_opt_out field
		jsonData := `{
			"type": "Settings",
			"audio": {
				"input": {
					"encoding": "linear16",
					"sample_rate": 16000
				}
			},
			"agent": {
				"language": "en"
			}
		}`

		var settings interfacesv1.SettingsOptions
		err := json.Unmarshal([]byte(jsonData), &settings)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify the field defaults to false
		if settings.MipOptOut != false {
			t.Errorf("Expected MipOptOut to default to false, got %v", settings.MipOptOut)
		}
	})
}

func TestSettingsMipOptOut_NewSettingsOptions(t *testing.T) {
	t.Run("Test NewSettingsOptions creates SettingsOptions with correct defaults", func(t *testing.T) {
		// Create new settings using the constructor
		settings := interfacesv1.NewSettingsOptions()

		// Verify basic structure
		if settings.Type != "Settings" {
			t.Errorf("Expected Type to be 'Settings', got %v", settings.Type)
		}

		// Verify audio defaults
		if settings.Audio.Input == nil {
			t.Error("Expected Audio.Input to be initialized")
		}
		if settings.Audio.Output == nil {
			t.Error("Expected Audio.Output to be initialized")
		}

		// Verify agent defaults
		if settings.Agent.Language != "en" {
			t.Errorf("Expected Agent.Language to be 'en', got %v", settings.Agent.Language)
		}

		// Marshal to JSON to verify structure
		jsonData, err := json.Marshal(settings)
		if err != nil {
			t.Fatalf("Failed to marshal settings to JSON: %v", err)
		}

		// Parse back to verify
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify that mip_opt_out defaults to false and is omitted when false (due to omitempty)
		if _, exists := result["mip_opt_out"]; exists {
			t.Errorf("Expected mip_opt_out to be omitted from JSON when false, but it was present")
		}
	})

	t.Run("Test setting mip_opt_out on NewSettingsOptions", func(t *testing.T) {
		// Create new settings and set mip_opt_out
		settings := interfacesv1.NewSettingsOptions()
		settings.MipOptOut = true

		// Marshal to JSON to verify
		jsonData, err := json.Marshal(settings)
		if err != nil {
			t.Fatalf("Failed to marshal settings to JSON: %v", err)
		}

		// Parse back to verify
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify mip_opt_out is present when true
		if result["mip_opt_out"] != true {
			t.Errorf("Expected mip_opt_out to be true, got %v", result["mip_opt_out"])
		}
	})

	t.Run("Test mip_opt_out with all other fields", func(t *testing.T) {
		// Create settings with all fields including mip_opt_out
		settings := &interfacesv1.SettingsOptions{
			Type:         "Settings",
			Experimental: true,
			MipOptOut:    true,
			Audio: interfacesv1.Audio{
				Input: &interfacesv1.Input{
					Encoding:   "linear16",
					SampleRate: 24000,
				},
				Output: &interfacesv1.Output{
					Encoding:   "mp3",
					SampleRate: 24000,
					Bitrate:    48000,
					Container:  "none",
				},
			},
			Agent: interfacesv1.Agent{
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
				Greeting: "Hello, I'm your AI assistant.",
			},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(settings)
		if err != nil {
			t.Fatalf("Failed to marshal settings to JSON: %v", err)
		}

		// Parse back to verify
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Verify mip_opt_out is present when true
		if result["mip_opt_out"] != true {
			t.Errorf("Expected mip_opt_out to be true, got %v", result["mip_opt_out"])
		}

		// Verify it's at the top level, not inside agent
		if agent, ok := result["agent"].(map[string]interface{}); ok {
			if _, exists := agent["mip_opt_out"]; exists {
				t.Error("mip_opt_out should not be in agent, it should be at top level")
			}
		}
	})
}
