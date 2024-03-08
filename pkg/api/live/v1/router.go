// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the live API for Deepgram
package live

import (
	"encoding/json"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

// MessageRouter is helper struct that routes events
type MessageRouter struct {
	callback       interfaces.LiveMessageCallback
	debugWebsocket bool
}

// NewWithDefault creates a default MessageRouter
func NewWithDefault() *MessageRouter {
	return New(NewDefaultCallbackHandler())
}

// New creates a MessageRouter with user defined callback
func New(callback interfaces.LiveMessageCallback) *MessageRouter {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG_WEBSOCKET"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG_WEBSOCKET found")
		debugStr = v
	}

	return &MessageRouter{
		callback:       callback,
		debugWebsocket: (strings.EqualFold(strings.ToLower(debugStr), "true")),
	}
}

// OpenResponse handles the OpenResponse message
func (r *MessageRouter) OpenHelper(or *interfaces.OpenResponse) error {
	obj, err := json.Marshal(or)
	if err != nil {
		klog.V(1).Infof("Open json.Marshal failed. Err: %v\n", err)
		return err
	}

	return r.OpenResponse(obj)
}

// OpenResponse handles the OpenResponse message
func (r *MessageRouter) OpenResponse(byMsg []byte) error {
	klog.V(6).Infof("router.OpenResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "OpenResponse", byMsg)

	var or interfaces.OpenResponse
	err := json.Unmarshal(byMsg, &or)
	if err != nil {
		klog.V(1).Infof("OpenResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.OpenResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Open(&or)
		if err != nil {
			klog.V(1).Infof("callback.OpenResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.OpenResponse succeeded\n")
		}
		klog.V(6).Infof("router.OpenResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.OpenResponse LEAVE\n")

	return nil
}

// Message handles platform messages
func (r *MessageRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("router.Message ENTER\n")

	if r.debugWebsocket {
		klog.V(5).Infof("Raw Message:\n%s\n", string(byMsg))
	}

	// what is the high level message here?
	var mt MessageType
	err := json.Unmarshal(byMsg, &mt)
	if err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(6).Infof("router.Message LEAVE\n")
		return err
	}

	klog.V(6).Infof("MessageType(%s) before\n", mt.Type)

	switch mt.Type {
	case interfaces.TypeMessageResponse:
		err = r.MessageResponse(byMsg)
	case interfaces.TypeMetadataResponse:
		err = r.MetadataResponse(byMsg)
	case interfaces.TypeSpeechStartedResponse:
		err = r.SpeechStartedResponse(byMsg)
	case interfaces.TypeUtteranceEndResponse:
		err = r.UtteranceEndResponse(byMsg)
	case interfaces.TypeErrorResponse:
		err = r.ErrorResponse(byMsg)
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

// MessageResponse handles the MessageResponse message
func (r *MessageRouter) MessageResponse(byMsg []byte) error {
	klog.V(6).Infof("router.MessageResponse ENTER\n")

	var mr interfaces.MessageResponse
	err := json.Unmarshal(byMsg, &mr)
	if err != nil {
		klog.V(1).Infof("MessageResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.MessageResponse LEAVE\n")
		return err
	}

	// this is too chatty is sentence is " " reduce the frequency of the log
	if len(mr.Channel.Alternatives) > 0 && len(strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)) == 0 {
		r.printDebugMessages(7, "MessageResponse", byMsg)
	} else {
		r.printDebugMessages(5, "MessageResponse", byMsg)
	}

	if r.callback != nil {
		err := r.callback.Message(&mr)
		if err != nil {
			klog.V(1).Infof("callback.MessageResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.MessageResponse succeeded\n")
		}
		klog.V(6).Infof("router.MessageResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.MessageResponse LEAVE\n")

	return ErrUserCallbackNotDefined
}

func (r *MessageRouter) MetadataResponse(byMsg []byte) error {
	klog.V(6).Infof("router.MetadataResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "MetadataResponse", byMsg)

	var md interfaces.MetadataResponse
	err := json.Unmarshal(byMsg, &md)
	if err != nil {
		klog.V(1).Infof("MetadataResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.MetadataResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Metadata(&md)
		if err != nil {
			klog.V(1).Infof("callback.MetadataResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.MetadataResponse succeeded\n")
		}
		klog.V(6).Infof("router.MetadataResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.MetadataResponse ENTER\n")

	return nil
}

func (r *MessageRouter) SpeechStartedResponse(byMsg []byte) error {
	klog.V(6).Infof("router.SpeechStartedResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "SpeechStartedResponse", byMsg)

	var ssr interfaces.SpeechStartedResponse
	err := json.Unmarshal(byMsg, &ssr)
	if err != nil {
		klog.V(1).Infof("SpeechStartedResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.SpeechStartedResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.SpeechStarted(&ssr)
		if err != nil {
			klog.V(1).Infof("callback.SpeechStartedResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.SpeechStartedResponse succeeded\n")
		}
		klog.V(6).Infof("router.SpeechStartedResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.SpeechStartedResponse ENTER\n")

	return nil
}

func (r *MessageRouter) UtteranceEndResponse(byMsg []byte) error {
	klog.V(6).Infof("router.UtteranceEndResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "UtteranceEndResponse", byMsg)

	var ur interfaces.UtteranceEndResponse
	err := json.Unmarshal(byMsg, &ur)
	if err != nil {
		klog.V(1).Infof("UtteranceEndResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.UtteranceEndResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.UtteranceEnd(&ur)
		if err != nil {
			klog.V(1).Infof("callback.UtteranceEndResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.UtteranceEndResponse succeeded\n")
		}
		klog.V(6).Infof("router.UtteranceEndResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.UtteranceEndResponse ENTER\n")

	return nil
}

// CloseHelper handles the CloseResponse message
func (r *MessageRouter) CloseHelper(cr *interfaces.CloseResponse) error {
	obj, err := json.Marshal(cr)
	if err != nil {
		klog.V(1).Infof("close json.Marshal failed. Err: %v\n", err)
		return err
	}

	return r.CloseResponse(obj)
}

// CloseResponse handles the CloseResponse message
func (r *MessageRouter) CloseResponse(byMsg []byte) error {
	klog.V(6).Infof("router.CloseResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "CloseResponse", byMsg)

	var cr interfaces.CloseResponse
	err := json.Unmarshal(byMsg, &cr)
	if err != nil {
		klog.V(1).Infof("CloseResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.CloseResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Close(&cr)
		if err != nil {
			klog.V(1).Infof("callback.CloseResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.CloseResponse succeeded\n")
		}
		klog.V(6).Infof("router.CloseResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.CloseResponse ENTER\n")

	return nil
}

func (r *MessageRouter) ErrorResponse(byMsg []byte) error {
	klog.V(6).Infof("router.ErrorResponse ENTER\n")

	// trace debugging
	r.printDebugMessages(5, "ErrorResponse", byMsg)

	var er interfaces.ErrorResponse
	err := json.Unmarshal(byMsg, &er)
	if err != nil {
		klog.V(1).Infof("ErrorResponse json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("router.ErrorResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Error(&er)
		if err != nil {
			klog.V(1).Infof("callback.ErrorResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.ErrorResponse succeeded\n")
		}
		klog.V(6).Infof("router.ErrorResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("User callback is undefined\n")
	klog.V(6).Infof("router.ErrorResponse ENTER\n")

	return nil
}

// UnhandledMessage handles the UnhandledMessage message
func (r *MessageRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("router.UnhandledMessage ENTER\n")

	// trace debugging
	r.printDebugMessages(3, "UnhandledMessage", byMsg)

	if r.callback != nil {
		err := r.callback.UnhandledEvent(byMsg)
		if err != nil {
			klog.V(1).Infof("callback.ErrorResponse failed. Err: %v\n", err)
		} else {
			klog.V(5).Infof("callback.ErrorResponse succeeded\n")
		}
		klog.V(6).Infof("router.ErrorResponse LEAVE\n")
		return err
	}

	klog.V(1).Infof("Unknown Event was received\n")
	klog.V(6).Infof("router.UnhandledMessage LEAVE\n")
	return ErrInvalidMessageType
}

func (r *MessageRouter) printDebugMessages(level klog.Level, function string, byMsg []byte) {
	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Format failed. Err: %v\n", err)
	}

	klog.V(level).Infof("\n\n-----------------------------------------------\n")
	klog.V(level).Infof("%s RAW:\n%s\n", function, prettyJson)
	klog.V(level).Infof("-----------------------------------------------\n\n\n")
}
