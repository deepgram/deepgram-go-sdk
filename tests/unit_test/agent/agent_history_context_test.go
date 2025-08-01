// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"strings"
	"testing"

	agentinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

// Test constants for maintainability
const (
	testUserRole             = "user"
	testAssistantRole        = "assistant"
	testWeatherQuery         = "What's the weather like today?"
	testAssistantResponse    = "Based on the current data, it's sunny with a temperature of 72°F (22°C). The forecast shows clear skies throughout the day."
	testFunctionCallID       = "fc_12345678-90ab-cdef-1234-567890abcdef"
	testFunctionName         = "check_order_status"
	testFunctionArguments    = `{"order_id": "ORD-123456"}`
	testFunctionResponse     = "Order #123456 status: Shipped - Expected delivery date: 2024-03-15"
	testFunctionCallIDSimple = "fc_test"
	testFunctionNameSimple   = "test_function"
	testSimpleArguments      = `{"param": "value"}`
	testSimpleResponse       = "Success"
	testSimpleContent        = "test"
	testHelloContent         = "Hello"
	testTestMessage          = "Test message"
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
			name: "Test_Flags_JSON_marshaling_false",
			test: func(t *testing.T) {
				flags := interfaces.Flags{
					History: false,
				}

				data, err := json.Marshal(flags)
				if err != nil {
					t.Fatalf("Failed to marshal Flags: %v", err)
				}

				expected := `{"history":false}`
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
					Type:    agentinterfaces.TypeHistoryConversationText,
					Role:    testUserRole,
					Content: testWeatherQuery,
				}

				if history.Type != agentinterfaces.TypeHistoryConversationText {
					t.Errorf("Expected Type to be '%s', got %s", agentinterfaces.TypeHistoryConversationText, history.Type)
				}
				if history.Role != testUserRole {
					t.Errorf("Expected Role to be '%s', got %s", testUserRole, history.Role)
				}
				if history.Content != testWeatherQuery {
					t.Errorf("Expected Content to match, got %s", history.Content)
				}
			},
		},
		{
			name: "Test_HistoryConversationText_JSON_marshaling",
			test: func(t *testing.T) {
				history := agentinterfaces.HistoryConversationText{
					Type:    agentinterfaces.TypeHistoryConversationText,
					Role:    testAssistantRole,
					Content: testAssistantResponse,
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
			name: "Test_HistoryConversationText_empty_type_included",
			test: func(t *testing.T) {
				// Test that empty Type field is included in JSON (not omitted)
				history := agentinterfaces.HistoryConversationText{
					Type:    "", // Empty type should still be included
					Role:    testUserRole,
					Content: testSimpleContent,
				}

				data, err := json.Marshal(history)
				if err != nil {
					t.Fatalf("Failed to marshal HistoryConversationText: %v", err)
				}

				jsonStr := string(data)
				if !strings.Contains(jsonStr, "\"type\":\"\"") {
					t.Errorf("Expected empty Type field to be included in JSON, got: %s", jsonStr)
				}
			},
		},
		{
			name: "Test_HistoryConversationText_API_spec_compliance",
			test: func(t *testing.T) {
				// Test example from API spec
				jsonData := `{
					"type": "` + agentinterfaces.TypeHistoryConversationText + `",
					"role": "` + testUserRole + `",
					"content": "` + testWeatherQuery + `"
				}`

				var history agentinterfaces.HistoryConversationText
				err := json.Unmarshal([]byte(jsonData), &history)
				if err != nil {
					t.Fatalf("Failed to unmarshal API spec example: %v", err)
				}

				if history.Type != agentinterfaces.TypeHistoryConversationText {
					t.Errorf("Expected Type '%s', got %s", agentinterfaces.TypeHistoryConversationText, history.Type)
				}
				if history.Role != testUserRole {
					t.Errorf("Expected Role '%s', got %s", testUserRole, history.Role)
				}
				if history.Content != testWeatherQuery {
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
					ID:         testFunctionCallID,
					Name:       testFunctionName,
					ClientSide: true,
					Arguments:  testFunctionArguments,
					Response:   testFunctionResponse,
				}

				history := agentinterfaces.HistoryFunctionCalls{
					Type:          agentinterfaces.TypeHistoryFunctionCalls,
					FunctionCalls: []agentinterfaces.FunctionCall{functionCall},
				}

				if history.Type != agentinterfaces.TypeHistoryFunctionCalls {
					t.Errorf("Expected Type to be '%s', got %s", agentinterfaces.TypeHistoryFunctionCalls, history.Type)
				}
				if len(history.FunctionCalls) != 1 {
					t.Errorf("Expected 1 function call, got %d", len(history.FunctionCalls))
				}

				fc := history.FunctionCalls[0]
				if fc.ID != testFunctionCallID {
					t.Errorf("Expected ID to match, got %s", fc.ID)
				}
				if fc.Name != testFunctionName {
					t.Errorf("Expected Name to be '%s', got %s", testFunctionName, fc.Name)
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
					Type: agentinterfaces.TypeHistoryFunctionCalls,
					FunctionCalls: []agentinterfaces.FunctionCall{
						{
							ID:         testFunctionCallIDSimple,
							Name:       testFunctionNameSimple,
							ClientSide: false,
							Arguments:  testSimpleArguments,
							Response:   testSimpleResponse,
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
				if fc.ID != testFunctionCallIDSimple {
					t.Errorf("Expected ID '%s', got %s", testFunctionCallIDSimple, fc.ID)
				}
				if fc.ClientSide != false {
					t.Errorf("Expected ClientSide false, got %v", fc.ClientSide)
				}
			},
		},
		{
			name: "Test_HistoryFunctionCalls_empty_type_included",
			test: func(t *testing.T) {
				// Test that empty Type field is included in JSON (not omitted)
				history := agentinterfaces.HistoryFunctionCalls{
					Type: "", // Empty type should still be included
					FunctionCalls: []agentinterfaces.FunctionCall{
						{
							ID:   testSimpleContent,
							Name: testFunctionNameSimple,
						},
					},
				}

				data, err := json.Marshal(history)
				if err != nil {
					t.Fatalf("Failed to marshal HistoryFunctionCalls: %v", err)
				}

				jsonStr := string(data)
				if !strings.Contains(jsonStr, "\"type\":\"\"") {
					t.Errorf("Expected empty Type field to be included in JSON, got: %s", jsonStr)
				}
			},
		},
		{
			name: "Test_HistoryFunctionCalls_API_spec_compliance",
			test: func(t *testing.T) {
				// Test example from API spec
				jsonData := `{
					"type": "` + agentinterfaces.TypeHistoryFunctionCalls + `",
					"function_calls": [
						{
							"id": "` + testFunctionCallID + `",
							"name": "` + testFunctionName + `",
							"client_side": true,
							"arguments": "{\"order_id\": \"ORD-123456\"}",
							"response": "` + testFunctionResponse + `"
						}
					]
				}`

				var history agentinterfaces.HistoryFunctionCalls
				err := json.Unmarshal([]byte(jsonData), &history)
				if err != nil {
					t.Fatalf("Failed to unmarshal API spec example: %v", err)
				}

				if history.Type != agentinterfaces.TypeHistoryFunctionCalls {
					t.Errorf("Expected Type '%s', got %s", agentinterfaces.TypeHistoryFunctionCalls, history.Type)
				}
				if len(history.FunctionCalls) != 1 {
					t.Errorf("Expected 1 function call, got %d", len(history.FunctionCalls))
				}

				fc := history.FunctionCalls[0]
				if fc.ID != testFunctionCallID {
					t.Errorf("Expected ID from API spec, got %s", fc.ID)
				}
				if fc.Name != testFunctionName {
					t.Errorf("Expected Name '%s', got %s", testFunctionName, fc.Name)
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
					Type:    agentinterfaces.TypeHistoryConversationText,
					Role:    testUserRole,
					Content: testHelloContent,
				}

				functionCalls := interfaces.HistoryFunctionCalls{
					Type: agentinterfaces.TypeHistoryFunctionCalls,
					FunctionCalls: []interfaces.FunctionCall{
						{
							ID:         testFunctionCallIDSimple,
							Name:       testFunctionNameSimple,
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
				if msg1.GetType() != agentinterfaces.TypeHistoryConversationText {
					t.Errorf("Expected first message type '%s', got %s", agentinterfaces.TypeHistoryConversationText, msg1.GetType())
				}

				msg2 := agent.Context.Messages[1]
				if msg2.GetType() != agentinterfaces.TypeHistoryFunctionCalls {
					t.Errorf("Expected second message type '%s', got %s", agentinterfaces.TypeHistoryFunctionCalls, msg2.GetType())
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
								Type:    agentinterfaces.TypeHistoryConversationText,
								Role:    testUserRole,
								Content: testTestMessage,
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
				if !strings.Contains(jsonStr, "context") {
					t.Error("Expected JSON to contain 'context' field")
				}
				if !strings.Contains(jsonStr, "messages") {
					t.Error("Expected JSON to contain 'messages' field")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.test)
	}
}
