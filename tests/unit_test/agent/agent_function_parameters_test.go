// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"strings"
	"testing"

	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

// TestParametersProperties tests OpenAI-style function schemas
func TestParametersProperties(t *testing.T) {
	t.Run("Test_Properties_JSON_marshaling", func(t *testing.T) {
		// Create a function with OpenAI-style parameter structure
		function := interfacesv1.Functions{
			Name:        "get_weather",
			Description: "Get current weather information for a specific location",
			Parameters: interfacesv1.Parameters{
				Type: "object",
				Properties: map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city and state/country, e.g. San Francisco, CA",
					},
					"unit": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"celsius", "fahrenheit"},
						"description": "The temperature unit to use",
					},
				},
				Required: []string{"location"},
			},
			Endpoint: &interfacesv1.Endpoint{
				Url:    "",
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
		}

		// Marshal to JSON
		jsonBytes, err := json.Marshal(function)
		if err != nil {
			t.Fatalf("Failed to marshal function: %v", err)
		}

		// Verify the JSON contains the expected OpenAI-style structure
		jsonStr := string(jsonBytes)

		// Check that it contains proper parameter structure
		if !strings.Contains(jsonStr, `"name":"get_weather"`) {
			t.Errorf("JSON should contain function name: %s", jsonStr)
		}

		if !strings.Contains(jsonStr, `"type":"object"`) {
			t.Errorf("JSON should contain parameters type: %s", jsonStr)
		}

		if !strings.Contains(jsonStr, `"location"`) {
			t.Errorf("JSON should contain location parameter: %s", jsonStr)
		}

		if !strings.Contains(jsonStr, `"unit"`) {
			t.Errorf("JSON should contain unit parameter: %s", jsonStr)
		}

		if !strings.Contains(jsonStr, `"required":["location"]`) {
			t.Errorf("JSON should contain required array: %s", jsonStr)
		}

		// Verify it can be unmarshaled back
		var unmarshaledFunction interfacesv1.Functions
		err = json.Unmarshal(jsonBytes, &unmarshaledFunction)
		if err != nil {
			t.Fatalf("Failed to unmarshal function: %v", err)
		}

		// Verify the unmarshaled data
		if unmarshaledFunction.Name != "get_weather" {
			t.Errorf("Expected name 'get_weather', got %s", unmarshaledFunction.Name)
		}

		if unmarshaledFunction.Parameters.Type != "object" {
			t.Errorf("Expected parameters type 'object', got %s", unmarshaledFunction.Parameters.Type)
		}

		if len(unmarshaledFunction.Parameters.Required) != 1 || unmarshaledFunction.Parameters.Required[0] != "location" {
			t.Errorf("Expected required ['location'], got %v", unmarshaledFunction.Parameters.Required)
		}

		// Check Properties was preserved
		if unmarshaledFunction.Parameters.Properties == nil {
			t.Error("Properties should not be nil after unmarshaling")
		}

		locationProp, exists := unmarshaledFunction.Parameters.Properties["location"]
		if !exists {
			t.Error("location property should exist in Properties")
		}

		locationMap, ok := locationProp.(map[string]interface{})
		if !ok {
			t.Error("location property should be a map")
		}

		if locationMap["type"] != "string" {
			t.Errorf("Expected location type 'string', got %v", locationMap["type"])
		}
	})
}
