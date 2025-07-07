// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	websocketv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent/v1/websocket"
)

func TestInjectUserMessage_StructCreation(t *testing.T) {
	t.Run("Test InjectUserMessage struct creation", func(t *testing.T) {
		testContent := "Hello, how can you help me today?"

		// Test creating the message struct
		msg := msginterfaces.InjectUserMessage{
			Type:    msginterfaces.TypeInjectUserMessage,
			Content: testContent,
		}

		// Verify the struct fields are set correctly
		if msg.Type != msginterfaces.TypeInjectUserMessage {
			t.Errorf("Expected Type to be %s, got %s", msginterfaces.TypeInjectUserMessage, msg.Type)
		}

		if msg.Content != testContent {
			t.Errorf("Expected Content to be %s, got %s", testContent, msg.Content)
		}
	})
}

func TestInjectUserMessage_JSONMarshaling(t *testing.T) {
	t.Run("Test InjectUserMessage JSON marshaling", func(t *testing.T) {
		testContent := "What services do you offer?"

		// Create the message
		msg := msginterfaces.InjectUserMessage{
			Type:    msginterfaces.TypeInjectUserMessage,
			Content: testContent,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(msg)
		if err != nil {
			t.Errorf("Failed to marshal InjectUserMessage to JSON: %v", err)
		}

		// Verify JSON structure
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON: %v", err)
		}

		// Check the JSON contains the expected fields
		if result["type"] != msginterfaces.TypeInjectUserMessage {
			t.Errorf("Expected JSON type field to be %s, got %v", msginterfaces.TypeInjectUserMessage, result["type"])
		}

		if result["content"] != testContent {
			t.Errorf("Expected JSON content field to be %s, got %v", testContent, result["content"])
		}
	})
}

func TestInjectUserMessage_JSONUnmarshaling(t *testing.T) {
	t.Run("Test InjectUserMessage JSON unmarshaling", func(t *testing.T) {
		// Create test JSON string using the constant
		testJSON := `{
			"type": "` + msginterfaces.TypeInjectUserMessage + `",
			"content": "Can you explain how this works?"
		}`

		// Unmarshal from JSON
		var msg msginterfaces.InjectUserMessage
		err := json.Unmarshal([]byte(testJSON), &msg)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON to InjectUserMessage: %v", err)
		}

		// Verify the struct was populated correctly
		if msg.Type != msginterfaces.TypeInjectUserMessage {
			t.Errorf("Expected Type to be %s, got %s", msginterfaces.TypeInjectUserMessage, msg.Type)
		}

		if msg.Content != "Can you explain how this works?" {
			t.Errorf("Expected Content to be 'Can you explain how this works?', got %s", msg.Content)
		}
	})
}

func TestInjectUserMessage_TypeAlias(t *testing.T) {
	t.Run("Test InjectUserMessage type alias works correctly", func(t *testing.T) {
		testContent := "I need assistance with my account"

		// Test using the client type alias
		clientMsg := websocketv1.InjectUserMessage{
			Type:    msginterfaces.TypeInjectUserMessage,
			Content: testContent,
		}

		// Verify the alias works the same as the interface type
		if clientMsg.Type != msginterfaces.TypeInjectUserMessage {
			t.Errorf("Expected Type to be %s, got %s", msginterfaces.TypeInjectUserMessage, clientMsg.Type)
		}

		if clientMsg.Content != testContent {
			t.Errorf("Expected Content to be %s, got %s", testContent, clientMsg.Content)
		}

		// Test JSON marshaling with the alias
		jsonData, err := json.Marshal(clientMsg)
		if err != nil {
			t.Errorf("Failed to marshal client type alias to JSON: %v", err)
		}

		// Verify it produces the same JSON structure
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON: %v", err)
		}

		if result["type"] != msginterfaces.TypeInjectUserMessage {
			t.Errorf("Expected JSON type field to be %s, got %v", msginterfaces.TypeInjectUserMessage, result["type"])
		}
	})
}

func TestInjectUserMessage_Constants(t *testing.T) {
	t.Run("Test InjectUserMessage constants are properly defined", func(t *testing.T) {
		// Verify the constant is defined and has the expected value
		expectedType := "InjectUserMessage"
		if msginterfaces.TypeInjectUserMessage != expectedType {
			t.Errorf("Expected TypeInjectUserMessage constant to be %s, got %s", expectedType, msginterfaces.TypeInjectUserMessage)
		}
	})
}

func TestInjectUserMessage_EmptyContent(t *testing.T) {
	t.Run("Test InjectUserMessage with empty content", func(t *testing.T) {
		// Test creating message with empty content
		msg := msginterfaces.InjectUserMessage{
			Type:    msginterfaces.TypeInjectUserMessage,
			Content: "",
		}

		// Should still marshal to valid JSON
		jsonData, err := json.Marshal(msg)
		if err != nil {
			t.Errorf("Failed to marshal InjectUserMessage with empty content: %v", err)
		}

		// Verify JSON structure
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON: %v", err)
		}

		// With omitempty tag, empty content field should be omitted from JSON entirely
		if _, exists := result["content"]; exists {
			t.Errorf("Expected content key to be omitted from JSON when empty, but it was present")
		}

		// Verify we can unmarshal back to the struct correctly
		var reconstructed msginterfaces.InjectUserMessage
		err = json.Unmarshal(jsonData, &reconstructed)
		if err != nil {
			t.Errorf("Failed to unmarshal back to struct: %v", err)
		}

		if reconstructed.Type != msginterfaces.TypeInjectUserMessage {
			t.Errorf("Expected Type to be %s, got %s", msginterfaces.TypeInjectUserMessage, reconstructed.Type)
		}

		if reconstructed.Content != "" {
			t.Errorf("Expected Content to be empty string, got %s", reconstructed.Content)
		}
	})
}

func TestInjectUserMessage_SpecialCharacters(t *testing.T) {
	t.Run("Test InjectUserMessage with special characters", func(t *testing.T) {
		// Test with various special characters and unicode
		testContent := "Hello! How do you handle émojis 🚀 and special chars like @#$%?"

		msg := msginterfaces.InjectUserMessage{
			Type:    msginterfaces.TypeInjectUserMessage,
			Content: testContent,
		}

		// Marshal and unmarshal to ensure proper handling
		jsonData, err := json.Marshal(msg)
		if err != nil {
			t.Errorf("Failed to marshal InjectUserMessage with special characters: %v", err)
		}

		var result msginterfaces.InjectUserMessage
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON with special characters: %v", err)
		}

		if result.Content != testContent {
			t.Errorf("Special characters not preserved. Expected %s, got %s", testContent, result.Content)
		}
	})
}
