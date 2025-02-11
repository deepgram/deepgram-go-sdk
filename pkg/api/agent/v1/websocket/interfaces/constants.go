// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
)

// These are the message types that can be received from the live API
type TypeResponse commoninterfaces.TypeResponse

// client message types
const (
	TypeSettingsConfiguration = "SettingsConfiguration"
	TypeUpdateInstructions    = "UpdateInstructions"
	TypeUpdateSpeak           = "UpdateSpeak"
	TypeInjectAgentMessage    = "InjectAgentMessage"
	TypeFunctionCallResponse  = "FunctionCallResponse"
	TypeKeepAlive             = "KeepAlive"
	TypeClose                 = "Close"
)

// server message types
const (
	// message types
	TypeOpenResponse                 = commoninterfaces.TypeOpenResponse
	TypeWelcomeResponse              = "Welcome"
	TypeConversationTextResponse     = "ConversationText"
	TypeUserStartedSpeakingResponse  = "UserStartedSpeaking"
	TypeAgentThinkingResponse        = "AgentThinking"
	TypeFunctionCallRequestResponse  = "FunctionCallRequest"
	TypeFunctionCallingResponse      = "FunctionCalling"
	TypeAgentStartedSpeakingResponse = "AgentStartedSpeaking"
	TypeAgentAudioDoneResponse       = "AgentAudioDone"
	TypeCloseResponse                = commoninterfaces.TypeCloseResponse
	TypeErrorResponse                = commoninterfaces.TypeErrorResponse
	TypeInjectionRefusedResponse     = "InjectionRefused"
	TypeSettingsAppliedResponse      = "SettingsApplied"
)
