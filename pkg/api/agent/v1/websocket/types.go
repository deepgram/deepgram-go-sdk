// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
)

/*
Using Channels
*/
// DefaultCallbackHandler is a default callback handler for live transcription
// Simply prints the transcript to stdout
type DefaultChanHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool

	binaryChan                   chan *[]byte
	openChan                     chan *interfaces.OpenResponse
	welcomeResponse              chan *interfaces.WelcomeResponse
	conversationTextResponse     chan *interfaces.ConversationTextResponse
	userStartedSpeakingResponse  chan *interfaces.UserStartedSpeakingResponse
	agentThinkingResponse        chan *interfaces.AgentThinkingResponse
	functionCallRequestResponse  chan *interfaces.FunctionCallRequestResponse
	functionCallingResponse      chan *interfaces.FunctionCallingResponse
	agentStartedSpeakingResponse chan *interfaces.AgentStartedSpeakingResponse
	agentAudioDoneResponse       chan *interfaces.AgentAudioDoneResponse
	injectionRefusedResponse     chan *interfaces.InjectionRefusedResponse
	keepAliveResponse            chan *interfaces.KeepAlive
	settingsAppliedResponse      chan *interfaces.SettingsAppliedResponse
	closeChan                    chan *interfaces.CloseResponse
	errorChan                    chan *interfaces.ErrorResponse
	unhandledChan                chan *[]byte
}

// ChanRouter routes events
type ChanRouter struct {
	debugWebsocket bool

	// call out to channels
	binaryChan                   []*chan *[]byte
	openChan                     []*chan *interfaces.OpenResponse
	welcomeResponse              []*chan *interfaces.WelcomeResponse
	conversationTextResponse     []*chan *interfaces.ConversationTextResponse
	userStartedSpeakingResponse  []*chan *interfaces.UserStartedSpeakingResponse
	agentThinkingResponse        []*chan *interfaces.AgentThinkingResponse
	functionCallRequestResponse  []*chan *interfaces.FunctionCallRequestResponse
	functionCallingResponse      []*chan *interfaces.FunctionCallingResponse
	agentStartedSpeakingResponse []*chan *interfaces.AgentStartedSpeakingResponse
	agentAudioDoneResponse       []*chan *interfaces.AgentAudioDoneResponse
	injectionRefusedResponse     []*chan *interfaces.InjectionRefusedResponse
	keepAliveResponse            []*chan *interfaces.KeepAlive
	settingsAppliedResponse      []*chan *interfaces.SettingsAppliedResponse
	closeChan                    []*chan *interfaces.CloseResponse
	errorChan                    []*chan *interfaces.ErrorResponse
	unhandledChan                []*chan *[]byte
}
