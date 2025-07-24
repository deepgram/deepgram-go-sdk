// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	agentinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

// TestHistoryFlags tests the new Flags struct with History field
func TestHistoryFlags(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test_Flags_struct_creation",
			test: func(t *testing.T) {
				flags := interfaces.Flags{
					History: true,
				}

				if flags.History != true {
					t.Errorf("Expected History to be true, got %v", flags.History)
				}
			},
		},
		{
			name: "Test_Flags_JSON_marshaling",
			test: func(t *testing.T) {
				flags := interfaces.Flags{
					History: true,
				}

				data, err := json.Marshal(flags)
				if err != nil {
					t.Fatalf("Failed to marshal Flags: %v", err)
				}

				expected := `{"history":true}`
				if string(data) != expected {
					t.Errorf("Expected JSON %s, got %s", expected, string(data))
				}
			},
		},
		{
			name: "Test_Flags_JSON_unmarshaling",
			test: func(t *testing.T) {
				jsonData := `{"history":false}`
				var flags interfaces.Flags

				err := json.Unmarshal([]byte(jsonData), &flags)
				if err != nil {
					t.Fatalf("Failed to unmarshal Flags: %v", err)
				}

				if flags.History != false {
					t.Errorf("Expected History to be false, got %v", flags.History)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.test)
	}
}

// TestSettingsOptionsWithFlags tests SettingsOptions with the new Flags field
func TestSettingsOptionsWithFlags(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test_SettingsOptions_with_Flags",
			test: func(t *testing.T) {
				options := interfaces.SettingsOptions{
					Type: "Settings",
					Flags: &interfaces.Flags{
						History: true,
					},
				}

				if options.Flags == nil {
					t.Fatal("Expected Flags to be set")
				}

				if !options.Flags.History {
					t.Errorf("Expected History to be true, got %v", options.Flags.History)
				}
			},
		},
		{
			name: "Test_NewSettingsOptions_default_flags",
			test: func(t *testing.T) {
				options := interfaces.NewSettingsOptions()

				if options.Flags == nil {
					t.Fatal("Expected Flags to be set in default options")
				}

				if !options.Flags.History {
					t.Errorf("Expected default History to be true, got %v", options.Flags.History)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.test)
	}
}

// TestHistoryConversationText tests the HistoryConversationText message type
func TestHistoryConversationText(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test_HistoryConversationText_struct_creation",
			test: func(t *testing.T) {
				history := agentinterfaces.HistoryConversationText{
					Type:    "History",
					Role:    "user",
					Content: "What's the weather like today?",
				}

				if history.Type != "History" {
					t.Errorf("Expected Type to be 'History', got %s", history.Type)
				}
				if history.Role != "user" {
					t.Errorf("Expected Role to be 'user', got %s", history.Role)
				}
				if history.Content != "What's the weather like today?" {
					t.Errorf("Expected Content to match, got %s", history.Content)
				}
			},
		},
		{
			name: "Test_HistoryConversationText_JSON_marshaling",
			test: func(t *testing.T) {
				history := agentinterfaces.HistoryConversationText{
					Type:    "History",
					Role:    "assistant",
					Content: "Based on the current data, it's sunny with a temperature of 72°F (22°C).",
				}

				data, err := json.Marshal(history)
				if err != nil {
					t.Fatalf("Failed to marshal HistoryConversationText: %v", err)
				}

				var unmarshaled agentinterfaces.HistoryConversationText
				err = json.Unmarshal(data, &unmarshaled)
				if err != nil {
					t.Fatalf("Failed to unmarshal HistoryConversationText: %v", err)
				}

				if unmarshaled.Type != history.Type {
					t.Errorf("Expected Type %s, got %s", history.Type, unmarshaled.Type)
				}
				if unmarshaled.Role != history.Role {
					t.Errorf("Expected Role %s, got %s", history.Role, unmarshaled.Role)
				}
				if unmarshaled.Content != history.Content {
					t.Errorf("Expected Content %s, got %s", history.Content, unmarshaled.Content)
				}
			},
		},
		{
			name: "Test_HistoryConversationText_API_spec_compliance",
			test: func(t *testing.T) {
				// Test example from API spec
				jsonData := `{
					"type": "History",
					"role": "user",
					"content": "What's the weather like today?"
				}`

				var history agentinterfaces.HistoryConversationText
				err := json.Unmarshal([]byte(jsonData), &history)
				if err != nil {
					t.Fatalf("Failed to unmarshal API spec example: %v", err)
				}

				if history.Type != "History" {
					t.Errorf("Expected Type 'History', got %s", history.Type)
				}
				if history.Role != "user" {
					t.Errorf("Expected Role 'user', got %s", history.Role)
				}
				if history.Content != "What's the weather like today?" {
					t.Errorf("Expected Content to match API spec example")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.test)
	}
}

// TestHistoryFunctionCalls tests the HistoryFunctionCalls message type
func TestHistoryFunctionCalls(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test_HistoryFunctionCalls_struct_creation",
			test: func(t *testing.T) {
				functionCall := agentinterfaces.FunctionCall{
					ID:         "fc_12345678-90ab-cdef-1234-567890abcdef",
					Name:       "check_order_status",
					ClientSide: true,
					Arguments:  `{"order_id": "ORD-123456"}`,
					Response:   "Order #123456 status: Shipped - Expected delivery date: 2024-03-15",
				}

				history := agentinterfaces.HistoryFunctionCalls{
					Type:          "History",
					FunctionCalls: []agentinterfaces.FunctionCall{functionCall},
				}

				if history.Type != "History" {
					t.Errorf("Expected Type to be 'History', got %s", history.Type)
				}
				if len(history.FunctionCalls) != 1 {
					t.Errorf("Expected 1 function call, got %d", len(history.FunctionCalls))
				}

				fc := history.FunctionCalls[0]
				if fc.ID != "fc_12345678-90ab-cdef-1234-567890abcdef" {
					t.Errorf("Expected ID to match, got %s", fc.ID)
				}
				if fc.Name != "check_order_status" {
					t.Errorf("Expected Name to be 'check_order_status', got %s", fc.Name)
				}
				if !fc.ClientSide {
					t.Errorf("Expected ClientSide to be true, got %v", fc.ClientSide)
				}
			},
		},
		{
			name: "Test_HistoryFunctionCalls_JSON_marshaling",
			test: func(t *testing.T) {
				history := agentinterfaces.HistoryFunctionCalls{
					Type: "History",
					FunctionCalls: []agentinterfaces.FunctionCall{
						{
							ID:         "fc_test",
							Name:       "test_function",
							ClientSide: false,
							Arguments:  `{"param": "value"}`,
							Response:   "Success",
						},
					},
				}

				data, err := json.Marshal(history)
				if err != nil {
					t.Fatalf("Failed to marshal HistoryFunctionCalls: %v", err)
				}

				var unmarshaled agentinterfaces.HistoryFunctionCalls
				err = json.Unmarshal(data, &unmarshaled)
				if err != nil {
					t.Fatalf("Failed to unmarshal HistoryFunctionCalls: %v", err)
				}

				if unmarshaled.Type != history.Type {
					t.Errorf("Expected Type %s, got %s", history.Type, unmarshaled.Type)
				}
				if len(unmarshaled.FunctionCalls) != len(history.FunctionCalls) {
					t.Errorf("Expected %d function calls, got %d", len(history.FunctionCalls), len(unmarshaled.FunctionCalls))
				}

				fc := unmarshaled.FunctionCalls[0]
				if fc.ID != "fc_test" {
					t.Errorf("Expected ID 'fc_test', got %s", fc.ID)
				}
				if fc.ClientSide != false {
					t.Errorf("Expected ClientSide false, got %v", fc.ClientSide)
				}
			},
		},
		{
			name: "Test_HistoryFunctionCalls_API_spec_compliance",
			test: func(t *testing.T) {
				// Test example from API spec
				jsonData := `{
					"type": "History",
					"function_calls": [
						{
							"id": "fc_12345678-90ab-cdef-1234-567890abcdef",
							"name": "check_order_status",
							"client_side": true,
							"arguments": "{\"order_id\": \"ORD-123456\"}",
							"response": "Order #123456 status: Shipped - Expected delivery date: 2024-03-15"
						}
					]
				}`

				var history agentinterfaces.HistoryFunctionCalls
				err := json.Unmarshal([]byte(jsonData), &history)
				if err != nil {
					t.Fatalf("Failed to unmarshal API spec example: %v", err)
				}

				if history.Type != "History" {
					t.Errorf("Expected Type 'History', got %s", history.Type)
				}
				if len(history.FunctionCalls) != 1 {
					t.Errorf("Expected 1 function call, got %d", len(history.FunctionCalls))
				}

				fc := history.FunctionCalls[0]
				if fc.ID != "fc_12345678-90ab-cdef-1234-567890abcdef" {
					t.Errorf("Expected ID from API spec, got %s", fc.ID)
				}
				if fc.Name != "check_order_status" {
					t.Errorf("Expected Name 'check_order_status', got %s", fc.Name)
				}
				if !fc.ClientSide {
					t.Errorf("Expected ClientSide true, got %v", fc.ClientSide)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.test)
	}
}

// TestAgentContext tests the Agent struct with Context field
func TestAgentContext(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test_Agent_with_Context",
			test: func(t *testing.T) {
				conversationText := interfaces.HistoryConversationText{
					Type:    "History",
					Role:    "user",
					Content: "Hello",
				}

				functionCalls := interfaces.HistoryFunctionCalls{
					Type: "History",
					FunctionCalls: []interfaces.FunctionCall{
						{
							ID:         "fc_test",
							Name:       "test_func",
							ClientSide: true,
							Arguments:  "{}",
							Response:   "OK",
						},
					},
				}

				agent := interfaces.Agent{
					Language: "en",
					Context: &interfaces.Context{
						Messages: []interfaces.ContextMessage{
							conversationText,
							functionCalls,
						},
					},
				}

				if agent.Context == nil {
					t.Fatal("Expected Context to be set")
				}
				if len(agent.Context.Messages) != 2 {
					t.Errorf("Expected 2 context messages, got %d", len(agent.Context.Messages))
				}

				// Verify the interface works
				msg1 := agent.Context.Messages[0]
				if msg1.GetType() != "History" {
					t.Errorf("Expected first message type 'History', got %s", msg1.GetType())
				}

				msg2 := agent.Context.Messages[1]
				if msg2.GetType() != "History" {
					t.Errorf("Expected second message type 'History', got %s", msg2.GetType())
				}
			},
		},
		{
			name: "Test_Agent_Context_JSON_marshaling",
			test: func(t *testing.T) {
				agent := interfaces.Agent{
					Language: "en",
					Context: &interfaces.Context{
						Messages: []interfaces.ContextMessage{
							interfaces.HistoryConversationText{
								Type:    "History",
								Role:    "user",
								Content: "Test message",
							},
						},
					},
				}

				data, err := json.Marshal(agent)
				if err != nil {
					t.Fatalf("Failed to marshal Agent with Context: %v", err)
				}

				// Verify the JSON contains expected fields
				jsonStr := string(data)
				if !contains(jsonStr, "context") {
					t.Error("Expected JSON to contain 'context' field")
				}
				if !contains(jsonStr, "messages") {
					t.Error("Expected JSON to contain 'messages' field")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.test)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || contains(s[1:], substr) || s[:len(substr)] == substr)
}
