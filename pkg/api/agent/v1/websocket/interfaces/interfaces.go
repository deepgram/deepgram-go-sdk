// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines interfaces for the live API
package interfacesv1

/*
Chan Interfaces
*/
// AgentMessageChan is a callback used to receive notifcations for platforms messages
type AgentMessageChan interface {
	GetBinary() []*chan *[]byte
	GetOpen() []*chan *OpenResponse
	GetWelcome() []*chan *WelcomeResponse
	GetConversationText() []*chan *ConversationTextResponse
	GetUserStartedSpeaking() []*chan *UserStartedSpeakingResponse
	GetAgentThinking() []*chan *AgentThinkingResponse
	GetFunctionCallRequest() []*chan *FunctionCallRequestResponse
	GetFunctionCalling() []*chan *FunctionCallingResponse
	GetAgentStartedSpeaking() []*chan *AgentStartedSpeakingResponse
	GetAgentAudioDone() []*chan *AgentAudioDoneResponse
	GetClose() []*chan *CloseResponse
	GetError() []*chan *ErrorResponse
	GetUnhandled() []*chan *[]byte
	GetInjectionRefused() []*chan *InjectionRefusedResponse
	GetKeepAlive() []*chan *KeepAlive
	GetSettingsApplied() []*chan *SettingsAppliedResponse
}
