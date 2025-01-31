// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/***********************************/
// Request/Input structs
/***********************************/
type SettingsConfigurationOptions interfaces.SettingsConfigurationOptions

// UpdateInstructions is the request to update the Agent instructions
type UpdateInstructions struct {
	Type         string `json:"type,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

// UpdateSpeak is the request to update model for speaking
type UpdateSpeak struct {
	Type  string `json:"type,omitempty"`
	Model string `json:"model,omitempty"`
}

// InjectAgentMessage is the request to inject a message into the Agent
type InjectAgentMessage struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// FunctionCallResponse is the response from a function call
type FunctionCallResponse struct {
	Type           string `json:"type,omitempty"`
	FunctionCallID string `json:"function_call_id,omitempty"`
	Output         string `json:"output,omitempty"`
}

// KeepAlive is the request to keep the connection alive
type KeepAlive struct {
	Type string `json:"type,omitempty"`
}

// Close terminates the connection
type Close struct {
	Type string `json:"type,omitempty"`
}

/***********************************/
// MessageType is the header to bootstrap you way unmarshalling other messages
/***********************************/
/*
	Example:
	{
		"type": "message",
		"message": {
			...
		}
	}
*/
type MessageType struct {
	Type string `json:"type"`
}

/***********************************/
// shared/common structs
/***********************************/
// None is a placeholder

/***********************************/
// Results from Agent/Server
/***********************************/
// OpenResponse is the response from opening the connection
type OpenResponse = commoninterfaces.OpenResponse

// WelcomeResponse is the response from the welcome message
type WelcomeResponse struct {
	Type      string `json:"type,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// ConversationTextResponse is the response from the conversation text
type ConversationTextResponse struct {
	Type    string `json:"type,omitempty"`
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// UserStartedSpeakingResponse is the response from the user starting to speak
type UserStartedSpeakingResponse struct {
	Type string `json:"type,omitempty"`
}

// AgentThinkingResponse is the response from the Agent thinking
type AgentThinkingResponse struct {
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
}

// FunctionCallRequestResponse is the response from a function call request
type FunctionCallRequestResponse struct {
	Type           string            `json:"type,omitempty"`
	FunctionName   string            `json:"function_name,omitempty"`
	FunctionCallID string            `json:"function_call_id,omitempty"`
	Input          map[string]string `json:"input,omitempty"` // TODO: this is still undefined
}

// FunctionCallingResponse is the response from a function calling
type FunctionCallingResponse struct {
	Type   string            `json:"type,omitempty"`
	Output map[string]string `json:"output,omitempty"` // TODO: this is still undefined
}

// AgentStartedSpeakingResponse is the response from the Agent starting to speak
type AgentStartedSpeakingResponse struct {
	Type         string  `json:"type,omitempty"`
	TotalLatency float64 `json:"total_latency,omitempty"`
	TtsLatency   float64 `json:"tts_latency,omitempty"`
	TttLatency   float64 `json:"ttt_latency,omitempty"`
}

// AgentAudioDoneResponse is the response from the Agent audio done
type AgentAudioDoneResponse struct {
	Type string `json:"type,omitempty"`
}

// CloseResponse is the response from closing the connection
type CloseResponse = commoninterfaces.CloseResponse

// ErrorResponse is the Deepgram specific response error
type ErrorResponse = interfaces.DeepgramError

// InjectionRefusedResponse is the response when an agent message injection is refused
type InjectionRefusedResponse struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// SettingsAppliedResponse is the response confirming settings were applied
type SettingsAppliedResponse struct {
	Type string `json:"type,omitempty"`
}
