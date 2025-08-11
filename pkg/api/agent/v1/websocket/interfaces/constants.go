// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
)

// These are the message types that can be received from the live API
type TypeResponse commoninterfaces.TypeResponse

// client message types
const (
	TypeSettings                = "Settings"
	TypeUpdatePrompt            = "UpdatePrompt"
	TypeUpdateSpeak             = "UpdateSpeak"
	TypeInjectAgentMessage      = "InjectAgentMessage"
	TypeInjectUserMessage       = "InjectUserMessage"
	TypeFunctionCallResponse    = "FunctionCallResponse"
	TypeHistoryConversationText = "History"
	TypeHistoryFunctionCalls    = "History"
	TypeKeepAlive               = "KeepAlive"
	TypeClose                   = "Close"
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
	TypeAgentStartedSpeakingResponse = "AgentStartedSpeaking"
	TypeAgentAudioDoneResponse       = "AgentAudioDone"
	TypeCloseResponse                = commoninterfaces.TypeCloseResponse
	TypeErrorResponse                = commoninterfaces.TypeErrorResponse
	TypeInjectionRefusedResponse     = "InjectionRefused"
	TypeSettingsAppliedResponse      = "SettingsApplied"
)
