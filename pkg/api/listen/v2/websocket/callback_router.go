// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"encoding/json"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
)

// NewCallbackWithDefault creates a CallbackRouter backed by the default stdout handler.
func NewCallbackWithDefault() *CallbackRouter {
	return NewCallbackRouter(NewDefaultCallbackHandler())
}

// NewCallbackRouter creates a CallbackRouter with the provided FluxMessageCallback.
func NewCallbackRouter(callback interfaces.FluxMessageCallback) *CallbackRouter {
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

// Open forwards the connection-opened event to the callback.
func (r *CallbackRouter) Open(or *interfaces.OpenResponse) error {
	return r.callback.Open(or)
}

// Close forwards the connection-closed event to the callback.
func (r *CallbackRouter) Close(cr *interfaces.CloseResponse) error {
	return r.callback.Close(cr)
}

// Error forwards a connection-level error to the callback.
func (r *CallbackRouter) Error(er *interfaces.ErrorResponse) error {
	return r.callback.Error(er)
}

// Binary handles binary frames (no-op for STT).
func (r *CallbackRouter) Binary(byMsg []byte) error {
	return nil
}

// Message decodes the "type" field of a server JSON message and routes to the appropriate
// callback method.
func (r *CallbackRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("flux.CallbackRouter.Message ENTER\n")

	if r.debugWebsocket {
		klog.V(5).Infof("Raw Message:\n%s\n", string(byMsg))
	}

	var mt interfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(6).Infof("flux.CallbackRouter.Message LEAVE\n")
		return err
	}

	var err error
	switch interfaces.TypeResponse(mt.Type) {
	case interfaces.TypeConnectedResponse:
		err = r.processConnected(byMsg)
	case interfaces.TypeTurnInfoResponse:
		err = r.processTurnInfo(byMsg)
	case interfaces.TypeConfigureSuccess:
		err = r.processConfigureSuccess(byMsg)
	case interfaces.TypeConfigureFailure:
		err = r.processConfigureFailure(byMsg)
	case interfaces.TypeFatalError:
		err = r.processFatalError(byMsg)
	default:
		err = r.UnhandledMessage(byMsg)
	}

	if err == nil {
		klog.V(6).Infof("MessageType(%s) succeeded\n", mt.Type)
	} else {
		klog.V(5).Infof("MessageType(%s) failed: %v\n", mt.Type, err)
	}
	klog.V(6).Infof("flux.CallbackRouter.Message LEAVE\n")
	return err
}

// UnhandledMessage logs and forwards any unrecognized message to the callback.
func (r *CallbackRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("flux.CallbackRouter.UnhandledMessage ENTER\n")
	r.printDebugMessages(3, "UnhandledMessage", byMsg)
	klog.V(1).Infof("Unknown Flux event received\n")
	klog.V(6).Infof("flux.CallbackRouter.UnhandledMessage LEAVE\n")
	err := r.callback.UnhandledEvent(byMsg)
	if err != nil {
		return err
	}
	return ErrInvalidMessageType
}

func (r *CallbackRouter) processConnected(byMsg []byte) error {
	var msg interfaces.ConnectedResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}
	return r.processGeneric("ConnectedResponse", byMsg, func(data *interface{}) error {
		return r.callback.Connected(&msg)
	}, msg)
}

func (r *CallbackRouter) processTurnInfo(byMsg []byte) error {
	var msg interfaces.TurnInfoResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}
	return r.processGeneric("TurnInfoResponse", byMsg, func(data *interface{}) error {
		return r.callback.TurnInfo(&msg)
	}, msg)
}

func (r *CallbackRouter) processConfigureSuccess(byMsg []byte) error {
	var msg interfaces.ConfigureSuccessResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}
	return r.processGeneric("ConfigureSuccessResponse", byMsg, func(data *interface{}) error {
		return r.callback.ConfigureSuccess(&msg)
	}, msg)
}

func (r *CallbackRouter) processConfigureFailure(byMsg []byte) error {
	var msg interfaces.ConfigureFailureResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}
	return r.processGeneric("ConfigureFailureResponse", byMsg, func(data *interface{}) error {
		return r.callback.ConfigureFailure(&msg)
	}, msg)
}

func (r *CallbackRouter) processFatalError(byMsg []byte) error {
	var msg interfaces.FatalErrorResponse
	if err := json.Unmarshal(byMsg, &msg); err != nil {
		return err
	}
	return r.processGeneric("FatalErrorResponse", byMsg, func(data *interface{}) error {
		return r.callback.FatalError(&msg)
	}, msg)
}

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
