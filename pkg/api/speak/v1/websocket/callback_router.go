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

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
)

// NewWithDefault creates a CallbackRouter with the default callback handler
func NewCallbackWithDefault() *CallbackRouter {
	var callback interfaces.SpeakMessageCallback
	handler := NewDefaultCallbackHandler()
	callback = handler
	return NewCallbackRouter(callback)
}

// New creates a CallbackRouter with a user-defined callback
func NewCallbackRouter(callback interfaces.SpeakMessageCallback) *CallbackRouter {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}
	return &CallbackRouter{
		callback:       callback,
		debugWebsocket: strings.EqualFold(strings.ToLower(debugStr), "true"),
	}
}

// OpenHelper handles the OpenResponse message
func (r *CallbackRouter) Open(or *interfaces.OpenResponse) error {
	return r.callback.Open(or)
}

// CloseHelper handles the OpenResponse message
func (r *CallbackRouter) Close(or *interfaces.CloseResponse) error {
	return r.callback.Close(or)
}

// ErrorHelper handles the ErrorResponse message
func (r *CallbackRouter) Error(er *interfaces.ErrorResponse) error {
	return r.callback.Error(er)
}

// processMessage generalizes the handling of all message types
func (r *CallbackRouter) processGeneric(msgType string, byMsg []byte, action func(data *interface{}) error, data interface{}) error {
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

func (r *CallbackRouter) processFlushed(byMsg []byte) error {
	var msg interfaces.FlushedResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Flush(&msg)
	}

	return r.processGeneric(string(interfaces.TypeFlushedResponse), byMsg, action, msg)
}

func (r *CallbackRouter) processCleared(byMsg []byte) error {
	var msg interfaces.ClearedResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Clear(&msg)
	}

	return r.processGeneric(string(interfaces.TypeClearedResponse), byMsg, action, msg)
}

func (r *CallbackRouter) processMetadata(byMsg []byte) error {
	var msg interfaces.MetadataResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Metadata(&msg)
	}

	return r.processGeneric(string(interfaces.TypeMetadataResponse), byMsg, action, msg)
}

func (r *CallbackRouter) processWarningResponse(byMsg []byte) error {
	var msg interfaces.WarningResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Warning(&msg)
	}

	return r.processGeneric(string(interfaces.TypeWarningResponse), byMsg, action, msg)
}

func (r *CallbackRouter) processErrorResponse(byMsg []byte) error {
	var msg interfaces.ErrorResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}

	action := func(data *interface{}) error {
		return r.callback.Error(&msg)
	}

	return r.processGeneric(string(interfaces.TypeErrorResponse), byMsg, action, msg)
}

// Message handles platform messages and routes them appropriately based on the MessageType
func (r *CallbackRouter) Message(byMsg []byte) error {
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
	case interfaces.TypeFlushedResponse:
		err = r.processFlushed(byMsg)
	case interfaces.TypeClearedResponse:
		err = r.processCleared(byMsg)
	case interfaces.TypeMetadataResponse:
		err = r.processMetadata(byMsg)
	case interfaces.TypeWarningResponse:
		err = r.processWarningResponse(byMsg)
	case interfaces.TypeResponse(interfaces.TypeErrorResponse):
		err = r.processErrorResponse(byMsg)
	default:
		err = r.UnhandledMessage(byMsg)
		klog.V(1).Infof("Message type %s is unhandled\n", mt.Type)
	}

	if err == nil {
		klog.V(6).Infof("MessageType(%s) after - Result: succeeded\n", mt.Type)
	} else {
		klog.V(5).Infof("MessageType(%s) after - Result: %v\n", mt.Type, err)
	}
	klog.V(6).Infof("router.Message LEAVE\n")
	return err
}

// Binary handles binary messages
func (r *CallbackRouter) Binary(byMsg []byte) error {
	klog.V(6).Infof("router.Binary ENTER\n")

	err := r.callback.Binary(byMsg)
	if err != nil {
		klog.V(1).Infof("callback.Binary failed. Err: %v\n", err)
	} else {
		klog.V(5).Infof("callback.Binary succeeded\n")
	}

	klog.V(6).Infof("router.Binary LEAVE\n")
	return err
}

// UnhandledMessage logs and handles any unexpected message types
func (r *CallbackRouter) UnhandledMessage(byMsg []byte) error {
	action := func(data *interface{}) error {
		return r.callback.UnhandledEvent(byMsg)
	}

	err := r.processGeneric(string(interfaces.TypeUnhandledResponse), byMsg, action, byMsg)
	if err != nil {
		klog.V(1).Infof("callback.UnhandledEvent failed. Err: %v\n", err)
	}

	return ErrInvalidMessageType
}

// printDebugMessages formats and logs debugging messages
func (r *CallbackRouter) printDebugMessages(level klog.Level, function string, byMsg []byte) {
	prettyJSON, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Format failed. Err: %v\n", err)
		return
	}
	klog.V(level).Infof("\n\n-----------------------------------------------\n")
	klog.V(level).Infof("%s RAW:\n%s\n", function, prettyJSON)
	klog.V(level).Infof("-----------------------------------------------\n\n\n")
}
