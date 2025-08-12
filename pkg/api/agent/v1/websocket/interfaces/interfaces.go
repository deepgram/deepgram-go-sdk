// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines interfaces for the live API
package interfacesv1

// AgentMessageChan is the core interface for receiving agent message notifications.
// This interface should remain unchanged to maintain backwards compatibility.
type AgentMessageChan interface {
	GetBinary() []*chan *[]byte
	GetOpen() []*chan *OpenResponse
	GetWelcome() []*chan *WelcomeResponse
	GetConversationText() []*chan *ConversationTextResponse
	GetUserStartedSpeaking() []*chan *UserStartedSpeakingResponse
	GetAgentThinking() []*chan *AgentThinkingResponse
	GetFunctionCallRequest() []*chan *FunctionCallRequestResponse
	GetAgentStartedSpeaking() []*chan *AgentStartedSpeakingResponse
	GetAgentAudioDone() []*chan *AgentAudioDoneResponse
	GetClose() []*chan *CloseResponse
	GetError() []*chan *ErrorResponse
	GetUnhandled() []*chan *[]byte
	GetInjectionRefused() []*chan *InjectionRefusedResponse
	GetKeepAlive() []*chan *KeepAlive
	GetSettingsApplied() []*chan *SettingsAppliedResponse
}

// HistoryMessageChan is an optional interface for receiving History message notifications.
// Implement this interface in addition to AgentMessageChan to be non-breaking.
type HistoryMessageChan interface {
	GetHistoryConversationText() []*chan *HistoryConversationText
	GetHistoryFunctionCalls() []*chan *HistoryFunctionCalls
}
