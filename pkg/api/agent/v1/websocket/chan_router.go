// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
)

// NewWithDefault creates a ChanRouter with the default callback handler
func NewChanWithDefault() *ChanRouter {
	chans := NewDefaultChanHandler()
	go func() {
		err := chans.Run()
		if err != nil {
			klog.V(1).Infof("chans.Run failed. Err: %v\n", err)
		}
	}()

	return NewChanRouter(chans)
}

// New creates a ChanRouter with a user-defined channels
// gocritic:ignore
func NewChanRouter(chans interfaces.AgentMessageChan) *ChanRouter {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	router := &ChanRouter{
		debugWebsocket:               strings.EqualFold(strings.ToLower(debugStr), "true"),
		binaryChan:                   make([]*chan *[]byte, 0),
		openChan:                     make([]*chan *interfaces.OpenResponse, 0),
		welcomeResponse:              make([]*chan *interfaces.WelcomeResponse, 0),
		conversationTextResponse:     make([]*chan *interfaces.ConversationTextResponse, 0),
		userStartedSpeakingResponse:  make([]*chan *interfaces.UserStartedSpeakingResponse, 0),
		agentThinkingResponse:        make([]*chan *interfaces.AgentThinkingResponse, 0),
		functionCallRequestResponse:  make([]*chan *interfaces.FunctionCallRequestResponse, 0),
		functionCallingResponse:      make([]*chan *interfaces.FunctionCallingResponse, 0),
		agentStartedSpeakingResponse: make([]*chan *interfaces.AgentStartedSpeakingResponse, 0),
		agentAudioDoneResponse:       make([]*chan *interfaces.AgentAudioDoneResponse, 0),
		injectionRefusedResponse:     make([]*chan *interfaces.InjectionRefusedResponse, 0),
		keepAliveResponse:            make([]*chan *interfaces.KeepAlive, 0),
		settingsAppliedResponse:      make([]*chan *interfaces.SettingsAppliedResponse, 0),
		closeChan:                    make([]*chan *interfaces.CloseResponse, 0),
		errorChan:                    make([]*chan *interfaces.ErrorResponse, 0),
		unhandledChan:                make([]*chan *[]byte, 0),
	}

	if chans != nil {
		router.binaryChan = append(router.binaryChan, chans.GetBinary()...)
		router.openChan = append(router.openChan, chans.GetOpen()...)
		router.welcomeResponse = append(router.welcomeResponse, chans.GetWelcome()...)
		router.conversationTextResponse = append(router.conversationTextResponse, chans.GetConversationText()...)
		router.userStartedSpeakingResponse = append(router.userStartedSpeakingResponse, chans.GetUserStartedSpeaking()...)
		router.agentThinkingResponse = append(router.agentThinkingResponse, chans.GetAgentThinking()...)
		router.functionCallRequestResponse = append(router.functionCallRequestResponse, chans.GetFunctionCallRequest()...)
		router.functionCallingResponse = append(router.functionCallingResponse, chans.GetFunctionCalling()...)
		router.agentStartedSpeakingResponse = append(router.agentStartedSpeakingResponse, chans.GetAgentStartedSpeaking()...)
		router.agentAudioDoneResponse = append(router.agentAudioDoneResponse, chans.GetAgentAudioDone()...)
		router.closeChan = append(router.closeChan, chans.GetClose()...)
		router.errorChan = append(router.errorChan, chans.GetError()...)
		router.unhandledChan = append(router.unhandledChan, chans.GetUnhandled()...)
		router.injectionRefusedResponse = append(router.injectionRefusedResponse, chans.GetInjectionRefused()...)
		router.keepAliveResponse = append(router.keepAliveResponse, chans.GetKeepAlive()...)
		router.settingsAppliedResponse = append(router.settingsAppliedResponse, chans.GetSettingsApplied()...)
	}

	return router
}

// Open sends an OpenResponse message to the callback
func (r *ChanRouter) Open(or *interfaces.OpenResponse) error {
	byMsg, err := json.Marshal(or)
	if err != nil {
		klog.V(1).Infof("json.Marshal(or) failed. Err: %v\n", err)
		return err
	}

	action := func(data []byte) error {
		var msg interfaces.OpenResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(OpenResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.openChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeOpenResponse), byMsg, action)
}

// Close sends an CloseResponse message to the callback
func (r *ChanRouter) Close(cr *interfaces.CloseResponse) error {
	byMsg, err := json.Marshal(cr)
	if err != nil {
		klog.V(1).Infof("json.Marshal(or) failed. Err: %v\n", err)
		return err
	}

	action := func(data []byte) error {
		var msg interfaces.CloseResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(CloseResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.closeChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeCloseResponse), byMsg, action)
}

// Error sends an ErrorResponse message to the callback
func (r *ChanRouter) Error(er *interfaces.ErrorResponse) error {
	byMsg, err := json.Marshal(er)
	if err != nil {
		klog.V(1).Infof("json.Marshal(er) failed. Err: %v\n", err)
		return err
	}

	action := func(data []byte) error {
		var msg interfaces.ErrorResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(ErrorResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.errorChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeErrorResponse), byMsg, action)
}

// processGeneric generalizes the handling of all message types
func (r *ChanRouter) processGeneric(msgType string, byMsg []byte, action func(data []byte) error) error {
	klog.V(6).Infof("router.%s ENTER\n", msgType)

	r.printDebugMessages(5, msgType, byMsg)

	var err error
	if err = action(byMsg); err != nil {
		klog.V(1).Infof("callback.%s failed. Err: %v\n", msgType, err)
	} else {
		klog.V(5).Infof("callback.%s succeeded\n", msgType)
	}
	klog.V(6).Infof("router.%s LEAVE\n", msgType)

	return err
}

func (r *ChanRouter) processWelcome(byMsg []byte) error {
	action := func(byMsg []byte) error {
		var msg interfaces.WelcomeResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(WelcomeResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.welcomeResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeWelcomeResponse), byMsg, action)
}

func (r *ChanRouter) processConversationText(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.ConversationTextResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(ConversationTextResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.conversationTextResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeConversationTextResponse), byMsg, action)
}

func (r *ChanRouter) processUserStartedSpeaking(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.UserStartedSpeakingResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(UserStartedSpeakingResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.userStartedSpeakingResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeUserStartedSpeakingResponse), byMsg, action)
}

func (r *ChanRouter) processAgentThinking(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.AgentThinkingResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(AgentThinkingResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.agentThinkingResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeAgentThinkingResponse), byMsg, action)
}

func (r *ChanRouter) processFunctionCallRequest(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.FunctionCallRequestResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(FunctionCallRequestResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.functionCallRequestResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeFunctionCallRequestResponse), byMsg, action)
}

func (r *ChanRouter) processFunctionCalling(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.FunctionCallingResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(FunctionCallingResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.functionCallingResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeFunctionCallingResponse), byMsg, action)
}

func (r *ChanRouter) processAgentStartedSpeaking(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.AgentStartedSpeakingResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(AgentStartedSpeakingResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.agentStartedSpeakingResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeAgentStartedSpeakingResponse), byMsg, action)
}

func (r *ChanRouter) processAgentAudioDone(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.AgentAudioDoneResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(AgentAudioDoneResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.agentAudioDoneResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeAgentAudioDoneResponse), byMsg, action)
}

func (r *ChanRouter) processErrorResponse(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.ErrorResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(MessageResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.errorChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeErrorResponse), byMsg, action)
}

func (r *ChanRouter) processInjectionRefused(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.InjectionRefusedResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(InjectionRefusedResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.injectionRefusedResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeInjectionRefusedResponse), byMsg, action)
}

func (r *ChanRouter) processKeepAlive(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.KeepAlive
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(KeepAlive) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.keepAliveResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeKeepAlive), byMsg, action)
}

func (r *ChanRouter) processSettingsApplied(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.SettingsAppliedResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(SettingsAppliedResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.settingsAppliedResponse {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeSettingsAppliedResponse), byMsg, action)
}

// Message handles platform messages and routes them appropriately based on the MessageType
func (r *ChanRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("router.Message ENTER\n")

	if r.debugWebsocket {
		klog.V(5).Infof("Raw Message:\n%s\n", string(byMsg))
	}

	var mt interfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(6).Infof("router.Message LEAVE\n")
		return err
	}

	var err error
	switch interfaces.TypeResponse(mt.Type) {
	case interfaces.TypeWelcomeResponse:
		err = r.processWelcome(byMsg)
	case interfaces.TypeConversationTextResponse:
		err = r.processConversationText(byMsg)
	case interfaces.TypeUserStartedSpeakingResponse:
		err = r.processUserStartedSpeaking(byMsg)
	case interfaces.TypeAgentThinkingResponse:
		err = r.processAgentThinking(byMsg)
	case interfaces.TypeFunctionCallRequestResponse:
		err = r.processFunctionCallRequest(byMsg)
	case interfaces.TypeFunctionCallingResponse:
		err = r.processFunctionCalling(byMsg)
	case interfaces.TypeAgentStartedSpeakingResponse:
		err = r.processAgentStartedSpeaking(byMsg)
	case interfaces.TypeAgentAudioDoneResponse:
		err = r.processAgentAudioDone(byMsg)
	case interfaces.TypeResponse(interfaces.TypeErrorResponse):
		err = r.processErrorResponse(byMsg)
	case interfaces.TypeInjectionRefusedResponse:
		err = r.processInjectionRefused(byMsg)
	case interfaces.TypeKeepAlive:
		err = r.processKeepAlive(byMsg)
	case interfaces.TypeSettingsAppliedResponse:
		err = r.processSettingsApplied(byMsg)
	default:
		err = r.UnhandledMessage(byMsg)
	}

	if err == nil {
		klog.V(6).Infof("MessageType(%s) after - Result: succeeded\n", mt.Type)
	} else {
		klog.V(5).Infof("MessageType(%s) after - Result: %v\n", mt.Type, err)
	}
	klog.V(6).Infof("router.Message LEAVE\n")
	return err
}

// Binary handles platform messages and routes them appropriately based on the MessageType
func (r *ChanRouter) Binary(byMsg []byte) error {
	klog.V(6).Infof("router.Binary ENTER\n")

	klog.V(5).Infof("Binary Message:\n%s...\n", hex.EncodeToString(byMsg[:20]))
	for _, ch := range r.binaryChan {
		*ch <- &byMsg
	}

	klog.V(6).Infof("router.Binary LEAVE\n")
	return nil
}

// UnhandledMessage logs and handles any unexpected message types
func (r *ChanRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("router.UnhandledMessage ENTER\n")
	r.printDebugMessages(3, "UnhandledMessage", byMsg)

	for _, ch := range r.unhandledChan {
		*ch <- &byMsg
	}

	klog.V(1).Infof("Unknown Event was received\n")
	klog.V(6).Infof("router.UnhandledMessage LEAVE\n")
	return ErrInvalidMessageType
}

// printDebugMessages formats and logs debugging messages
func (r *ChanRouter) printDebugMessages(level klog.Level, function string, byMsg []byte) {
	prettyJSON, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Format failed. Err: %v\n", err)
		return
	}
	klog.V(level).Infof("\n\n-----------------------------------------------\n")
	klog.V(level).Infof("%s RAW:\n%s\n", function, prettyJSON)
	klog.V(level).Infof("-----------------------------------------------\n\n\n")
}
