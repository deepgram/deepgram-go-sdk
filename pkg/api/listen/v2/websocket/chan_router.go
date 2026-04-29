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

// NewChanWithDefault creates a ChanRouter backed by the default channel handler.
func NewChanWithDefault() *ChanRouter {
	chans := NewDefaultChanHandler()
	go func() {
		if err := chans.Run(); err != nil {
			klog.V(1).Infof("chans.Run failed. Err: %v\n", err)
		}
	}()
	return NewChanRouter(chans)
}

// NewChanRouter creates a ChanRouter that dispatches messages to the provided FluxMessageChan.
func NewChanRouter(chans interfaces.FluxMessageChan) *ChanRouter { //nolint:gocritic
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	router := &ChanRouter{
		debugWebsocket:       strings.EqualFold(strings.ToLower(debugStr), "true"),
		openChan:             make([]*chan *interfaces.OpenResponse, 0),
		connectedChan:        make([]*chan *interfaces.ConnectedResponse, 0),
		turnInfoChan:         make([]*chan *interfaces.TurnInfoResponse, 0),
		configureSuccessChan: make([]*chan *interfaces.ConfigureSuccessResponse, 0),
		configureFailureChan: make([]*chan *interfaces.ConfigureFailureResponse, 0),
		fatalErrorChan:       make([]*chan *interfaces.FatalErrorResponse, 0),
		closeChan:            make([]*chan *interfaces.CloseResponse, 0),
		errorChan:            make([]*chan *interfaces.ErrorResponse, 0),
		unhandledChan:        make([]*chan *[]byte, 0),
	}

	if chans != nil {
		router.openChan = append(router.openChan, chans.GetOpen()...)
		router.connectedChan = append(router.connectedChan, chans.GetConnected()...)
		router.turnInfoChan = append(router.turnInfoChan, chans.GetTurnInfo()...)
		router.configureSuccessChan = append(router.configureSuccessChan, chans.GetConfigureSuccess()...)
		router.configureFailureChan = append(router.configureFailureChan, chans.GetConfigureFailure()...)
		router.fatalErrorChan = append(router.fatalErrorChan, chans.GetFatalError()...)
		router.closeChan = append(router.closeChan, chans.GetClose()...)
		router.errorChan = append(router.errorChan, chans.GetError()...)
		router.unhandledChan = append(router.unhandledChan, chans.GetUnhandled()...)
	}

	return router
}

// Open sends the connection-opened event to registered open channels.
func (r *ChanRouter) Open(or *interfaces.OpenResponse) error {
	byMsg, err := json.Marshal(or)
	if err != nil {
		return err
	}
	action := func(data []byte) error {
		var msg interfaces.OpenResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.openChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeOpenResponse), byMsg, action)
}

// Close sends the connection-closed event to registered close channels.
func (r *ChanRouter) Close(cr *interfaces.CloseResponse) error {
	byMsg, err := json.Marshal(cr)
	if err != nil {
		return err
	}
	action := func(data []byte) error {
		var msg interfaces.CloseResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.closeChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeCloseResponse), byMsg, action)
}

// Error sends a connection-level error to registered error channels.
func (r *ChanRouter) Error(er *interfaces.ErrorResponse) error {
	byMsg, err := json.Marshal(er)
	if err != nil {
		return err
	}
	action := func(data []byte) error {
		var msg interfaces.ErrorResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.errorChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeErrorResponse), byMsg, action)
}

// Binary handles binary frames (no-op for STT).
func (r *ChanRouter) Binary(byMsg []byte) error {
	return nil
}

// Message decodes the "type" field and routes the payload to the appropriate channels.
func (r *ChanRouter) Message(byMsg []byte) error {
	klog.V(6).Infof("flux.ChanRouter.Message ENTER\n")

	if r.debugWebsocket {
		klog.V(5).Infof("Raw Message:\n%s\n", string(byMsg))
	}

	var mt interfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(6).Infof("flux.ChanRouter.Message LEAVE\n")
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
	klog.V(6).Infof("flux.ChanRouter.Message LEAVE\n")
	return err
}

// UnhandledMessage routes an unrecognized message to the unhandled channels.
func (r *ChanRouter) UnhandledMessage(byMsg []byte) error {
	klog.V(6).Infof("flux.ChanRouter.UnhandledMessage ENTER\n")
	r.printDebugMessages(3, "UnhandledMessage", byMsg)
	klog.V(1).Infof("Unknown Flux event received\n")
	klog.V(6).Infof("flux.ChanRouter.UnhandledMessage LEAVE\n")

	action := func(data []byte) error {
		for _, ch := range r.unhandledChan {
			*ch <- &data
		}
		return nil
	}
	err := r.processGeneric("UnhandledMessage", byMsg, action)
	if err != nil {
		return err
	}
	return ErrInvalidMessageType
}

func (r *ChanRouter) processConnected(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.ConnectedResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.connectedChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeConnectedResponse), byMsg, action)
}

func (r *ChanRouter) processTurnInfo(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.TurnInfoResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.turnInfoChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeTurnInfoResponse), byMsg, action)
}

func (r *ChanRouter) processConfigureSuccess(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.ConfigureSuccessResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.configureSuccessChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeConfigureSuccess), byMsg, action)
}

func (r *ChanRouter) processConfigureFailure(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.ConfigureFailureResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.configureFailureChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeConfigureFailure), byMsg, action)
}

func (r *ChanRouter) processFatalError(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.FatalErrorResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		for _, ch := range r.fatalErrorChan {
			*ch <- &msg
		}
		return nil
	}
	return r.processGeneric(string(interfaces.TypeFatalError), byMsg, action)
}

func (r *ChanRouter) processGeneric(msgType string, byMsg []byte, action func(data []byte) error) error {
	klog.V(6).Infof("router.%s ENTER\n", msgType)
	r.printDebugMessages(5, msgType, byMsg)

	var err error
	if err = action(byMsg); err != nil {
		klog.V(1).Infof("chan.%s failed. Err: %v\n", msgType, err)
	} else {
		klog.V(5).Infof("chan.%s succeeded\n", msgType)
	}
	klog.V(6).Infof("router.%s LEAVE\n", msgType)
	return err
}

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
