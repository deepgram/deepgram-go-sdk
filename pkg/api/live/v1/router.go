// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (
	"encoding/json"
	"errors"
	"log"

	prettyjson "github.com/hokaccha/go-prettyjson"

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

// MessageRouter is helper struct that routes events
type MessageRouter struct {
	callback interfaces.LiveMessageCallback
}

// NewWithDefault creates a default MessageRouter
func NewWithDefault() *MessageRouter {
	return New(NewDefaultCallbackHandler())
}

// New creates a MessageRouter with user defined callback
func New(callback interfaces.LiveMessageCallback) *MessageRouter {
	return &MessageRouter{
		callback: callback,
	}
}

// Message handles platform messages
func (r *MessageRouter) Message(byMsg []byte) error {
	// log.Printf("MessageRouter::Message ENTER\n")

	// what is the high level message here?
	var mt MessageType
	err := json.Unmarshal(byMsg, &mt)
	if err != nil {
		log.Printf("MessageRouter json.Unmarshal(MessageType) failed. Err: %v\n", err)
		log.Printf("MessageRouter LEAVE\n")
		return err
	}

	switch mt.Type {
	case interfaces.TypeErrorResponse:
		return r.HandleError(byMsg)
	case interfaces.TypeMessageResponse:
		return r.MessageResponse(byMsg)
	default:
		return r.UnhandledMessage(byMsg)
	}

	return nil
}

// MessageResponse handles the MessageResponse message
func (r *MessageRouter) MessageResponse(byMsg []byte) error {
	// log.Printf("MessageResponse ENTER\n")

	// trace debugging
	// r.printDebugMessages("MessageResponse", byMsg)

	var mr interfaces.MessageResponse
	err := json.Unmarshal(byMsg, &mr)
	if err != nil {
		log.Printf("MessageResponse json.Unmarshal failed. Err: %v\n", err)
		log.Printf("MessageResponse LEAVE\n")
		return err
	}

	if r.callback != nil {
		err := r.callback.Message(&mr)
		if err != nil {
			log.Printf("callback.MessageResponse failed. Err: %v\n", err)
			// } else {
			// 	log.Printf("callback.MessageResponse succeeded\n")
		}
		// log.Printf("MessageResponse LEAVE\n")
		return err
	}

	// log.Printf("User callback is undefined\n")
	// log.Printf("MessageResponse LEAVE\n")
	return ErrUserCallbackNotDefined
}

// HandleError handles error messages
func (r *MessageRouter) HandleError(byMsg []byte) error {
	log.Printf("HandleError ENTER\n")

	// trace debugging
	r.printDebugMessages("HandleError", byMsg)

	var er interfaces.ErrorResponse
	err := json.Unmarshal(byMsg, &er)
	if err != nil {
		log.Printf("HandleError json.Unmarshal failed. Err: %v\n", err)
		log.Printf("HandleError LEAVE\n")
		return err
	}

	b, err := json.MarshalIndent(er, "", "    ")
	if err != nil {
		log.Printf("HandleError MarshalIndent failed. Err: %v\n", err)
		log.Printf("HandleError LEAVE\n")
		return err
	}

	log.Printf("\n\nError: %s\n\n", string(b))
	log.Printf("HandleError LEAVE\n")
	return errors.New(string(b))
}

// UnhandledMessage handles the UnhandledMessage message
func (r *MessageRouter) UnhandledMessage(byMsg []byte) error {
	log.Printf("UnhandledMessage ENTER\n")

	// trace debugging
	r.printDebugMessages("UnhandledMessage", byMsg)

	log.Printf("User callback is undefined\n")
	log.Printf("UnhandledMessage LEAVE\n")
	return ErrInvalidMessageType
}

func (r *MessageRouter) printDebugMessages(function string, byMsg []byte) {
	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		log.Printf("prettyjson.Marshal failed. Err: %v\n", err)
	}

	log.Printf("\n\n-----------------------------------------------\n")
	log.Printf("%s RAW:\n%s\n", function, prettyJson)
	log.Printf("-----------------------------------------------\n\n\n")
}
