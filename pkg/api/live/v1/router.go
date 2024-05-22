// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT
package live

import (
	"encoding/json"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

// MessageRouter routes events
type MessageRouter struct {
	callback       interfaces.LiveMessageCallback
	debugWebsocket bool
}

// NewWithDefault creates a MessageRouter with the default callback handler
func NewWithDefault() *MessageRouter {
	return New(NewDefaultCallbackHandler())
}

// New creates a MessageRouter with a user-defined callback
func New(callback interfaces.LiveMessageCallback) *MessageRouter {
	debugStr := os.Getenv("DEEPGRAM_DEBUG_WEBSOCKET")
	return &MessageRouter{
		callback:       callback,
		debugWebsocket: strings.EqualFold(strings.ToLower(debugStr), "true"),
	}
}

// OpenHelper handles the OpenResponse message
func (r *MessageRouter) OpenHelper(or *interfaces.OpenResponse) error {
	return r.callback.Open(or)
}

// OpenResponse handles the OpenResponse message
func (r *MessageRouter) CloseHelper(or *interfaces.CloseResponse) error {
	return r.callback.Close(or)
}

// ErrorResponse handles the OpenResponse message
func (r *MessageRouter) ErrorHelper(er *interfaces.ErrorResponse) error {
	return r.callback.Error(er)
}

// processMessage generalizes the handling of all message types
func (r *MessageRouter) processGeneric(msgType string, byMsg []byte, action func(data *interface{}) error, data interface{}) error {
	klog.V(6).Infof("router.%s ENTER\n", msgType)

	r.printDebugMessages(5, msgType, byMsg)

	var err error
	if err = action(&data); err != nil {
		klog.V(1).Infof("callback.%s failed. Err: %v\n", msgType, err)
	} else {
		klog.V(5).Infof("callback.%s succeeded\n", msgType)
	}
	klog.V(6).Infof("router.%s LEAVE\n", msgType)

	return err
}

func (r *MessageRouter) processMessage(byMsg []byte) error {
	var msg interfaces.MessageResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Message(&msg)
	}

	return r.processGeneric("MessageResponse", byMsg, action, msg)
}

func (r *MessageRouter) processMetadata(byMsg []byte) error {
	var msg interfaces.MetadataResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Metadata(&msg)
	}

	return r.processGeneric("MetadataResponse", byMsg, action, msg)
}

func (r *MessageRouter) processSpeechStartedResponse(byMsg []byte) error {
	var msg interfaces.SpeechStartedResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.SpeechStarted(&msg)
	}

	return r.processGeneric("SpeechStartedResponse", byMsg, action, msg)
}

func (r *MessageRouter) processUtteranceEndResponse(byMsg []byte) error {
	var msg interfaces.UtteranceEndResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.UtteranceEnd(&msg)
	}

	return r.processGeneric("UtteranceEndResponse", byMsg, action, msg)
}

func (r *MessageRouter) processErrorResponse(byMsg []byte) error {
	var msg interfaces.ErrorResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Error(&msg)
	}

	return r.processGeneric("ErrorResponse", byMsg, action, msg)
}

// Message handles platform messages and routes them appropriately based on the MessageType
func (r *MessageRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("router.Message ENTER\n")

	if r.debugWebsocket {
		klog.V(5).Infof("Raw Message:\n%s\n", string(byMsg))
	}

	var mt MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(6).Infof("router.Message LEAVE\n")
		return err
	}

	var err error
	switch mt.Type {
	case interfaces.TypeMessageResponse:
		err = r.processMessage(byMsg)
	case interfaces.TypeMetadataResponse:
		err = r.processMetadata(byMsg)
	case interfaces.TypeSpeechStartedResponse:
		err = r.processSpeechStartedResponse(byMsg)
	case interfaces.TypeUtteranceEndResponse:
		err = r.processUtteranceEndResponse(byMsg)
	case interfaces.TypeErrorResponse:
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

// UnhandledMessage logs and handles any unexpected message types
func (r *MessageRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("router.UnhandledMessage ENTER\n")
	r.printDebugMessages(3, "UnhandledMessage", byMsg)
	klog.V(1).Infof("Unknown Event was received\n")
	klog.V(6).Infof("router.UnhandledMessage LEAVE\n")
	return ErrInvalidMessageType
}

// printDebugMessages formats and logs debugging messages
func (r *MessageRouter) printDebugMessages(level klog.Level, function string, byMsg []byte) {
	prettyJSON, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Format failed. Err: %v\n", err)
		return
	}
	klog.V(level).Infof("\n\n-----------------------------------------------\n")
	klog.V(level).Infof("%s RAW:\n%s\n", function, prettyJSON)
	klog.V(level).Infof("-----------------------------------------------\n\n\n")
}
