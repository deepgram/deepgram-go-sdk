// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"encoding/json"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
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
func NewChanRouter(chans interfaces.LiveMessageChan) *ChanRouter {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	router := &ChanRouter{
		debugWebsocket:    strings.EqualFold(strings.ToLower(debugStr), "true"),
		openChan:          make([]*chan *interfaces.OpenResponse, 0),
		messageChan:       make([]*chan *interfaces.MessageResponse, 0),
		metadataChan:      make([]*chan *interfaces.MetadataResponse, 0),
		speechStartedChan: make([]*chan *interfaces.SpeechStartedResponse, 0),
		utteranceEndChan:  make([]*chan *interfaces.UtteranceEndResponse, 0),
		closeChan:         make([]*chan *interfaces.CloseResponse, 0),
		errorChan:         make([]*chan *interfaces.ErrorResponse, 0),
		unhandledChan:     make([]*chan *[]byte, 0),
	}

	if chans != nil {
		router.openChan = append(router.openChan, chans.GetOpen()...)
		router.messageChan = append(router.messageChan, chans.GetMessage()...)
		router.metadataChan = append(router.metadataChan, chans.GetMetadata()...)
		router.speechStartedChan = append(router.speechStartedChan, chans.GetSpeechStarted()...)
		router.utteranceEndChan = append(router.utteranceEndChan, chans.GetUtteranceEnd()...)
		router.closeChan = append(router.closeChan, chans.GetClose()...)
		router.errorChan = append(router.errorChan, chans.GetError()...)
		router.unhandledChan = append(router.unhandledChan, chans.GetUnhandled()...)
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

func (r *ChanRouter) processMessage(byMsg []byte) error {
	action := func(byMsg []byte) error {
		var msg interfaces.MessageResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(MessageResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.messageChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeMessageResponse), byMsg, action)
}

func (r *ChanRouter) processMetadata(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.MetadataResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(MessageResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.metadataChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeMetadataResponse), byMsg, action)
}

func (r *ChanRouter) processSpeechStartedResponse(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.SpeechStartedResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(MessageResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.speechStartedChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeSpeechStartedResponse), byMsg, action)
}

func (r *ChanRouter) processUtteranceEndResponse(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.UtteranceEndResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(UtteranceEndResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.utteranceEndChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeUtteranceEndResponse), byMsg, action)
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
	case interfaces.TypeMessageResponse:
		err = r.processMessage(byMsg)
	case interfaces.TypeMetadataResponse:
		err = r.processMetadata(byMsg)
	case interfaces.TypeSpeechStartedResponse:
		err = r.processSpeechStartedResponse(byMsg)
	case interfaces.TypeUtteranceEndResponse:
		err = r.processUtteranceEndResponse(byMsg)
	case interfaces.TypeResponse(interfaces.TypeErrorResponse):
		err = r.processErrorResponse(byMsg)
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
	// No implementation needed on STT
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
