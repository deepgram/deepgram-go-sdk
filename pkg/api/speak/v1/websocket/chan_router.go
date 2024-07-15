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

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
)

// NewWithDefault creates a ChanRouter with the default callback handler
func NewChanRouterWithDefault() *ChanRouter {
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
func NewChanRouter(chans interfaces.SpeakMessageChan) *ChanRouter {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	router := &ChanRouter{
		debugWebsocket: strings.EqualFold(strings.ToLower(debugStr), "true"),
		binaryChan:     make([]*chan *[]byte, 0),
		openChan:       make([]*chan *interfaces.OpenResponse, 0),
		metadataChan:   make([]*chan *interfaces.MetadataResponse, 0),
		flushedChan:    make([]*chan *interfaces.FlushedResponse, 0),
		clearedChan:    make([]*chan *interfaces.ClearedResponse, 0),
		closeChan:      make([]*chan *interfaces.CloseResponse, 0),
		warningChan:    make([]*chan *interfaces.WarningResponse, 0),
		errorChan:      make([]*chan *interfaces.ErrorResponse, 0),
		unhandledChan:  make([]*chan *[]byte, 0),
	}

	if chans != nil {
		router.binaryChan = append(router.binaryChan, chans.GetBinary()...)
		router.openChan = append(router.openChan, chans.GetOpen()...)
		router.metadataChan = append(router.metadataChan, chans.GetMetadata()...)
		router.flushedChan = append(router.flushedChan, chans.GetFlush()...)
		router.clearedChan = append(router.clearedChan, chans.GetClear()...)
		router.closeChan = append(router.closeChan, chans.GetClose()...)
		router.warningChan = append(router.warningChan, chans.GetWarning()...)
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

func (r *ChanRouter) processMetadata(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.MetadataResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(MetadataResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.metadataChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeMetadataResponse), byMsg, action)
}

func (r *ChanRouter) processFlushed(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.FlushedResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(FlushedResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.flushedChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeFlushedResponse), byMsg, action)
}

func (r *ChanRouter) processCleared(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.ClearedResponse
		if err := json.Unmarshal(byMsg, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(ClearedResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.clearedChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeClearedResponse), byMsg, action)
}

func (r *ChanRouter) processWarningResponse(byMsg []byte) error {
	action := func(data []byte) error {
		var msg interfaces.WarningResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			klog.V(1).Infof("json.Unmarshal(WarningResponse) failed. Err: %v\n", err)
			return err
		}

		for _, ch := range r.warningChan {
			*ch <- &msg
		}
		return nil
	}

	return r.processGeneric(string(interfaces.TypeWarningResponse), byMsg, action)
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
	case interfaces.TypeMetadataResponse:
		err = r.processMetadata(byMsg)
	case interfaces.TypeFlushedResponse:
		err = r.processFlushed(byMsg)
	case interfaces.TypeClearedResponse:
		err = r.processCleared(byMsg)
	case interfaces.TypeWarningResponse:
		err = r.processWarningResponse(byMsg)
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
