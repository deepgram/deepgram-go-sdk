// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
)

// NewDefaultChanHandler creates a DefaultChanHandler with buffered channels for all Flux events.
// The handler's Run() goroutine prints events to stdout; it is started automatically.
func NewDefaultChanHandler() *DefaultChanHandler {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}
	var debugExtStr string
	if v := os.Getenv("DEEPGRAM_DEBUG_VERBOSE"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG_VERBOSE found")
		debugExtStr = v
	}
	handler := &DefaultChanHandler{
		debugWebsocket:        strings.EqualFold(debugStr, "true"),
		debugWebsocketVerbose: strings.EqualFold(debugExtStr, "true"),
		openChan:              make(chan *interfaces.OpenResponse),
		connectedChan:         make(chan *interfaces.ConnectedResponse),
		turnInfoChan:          make(chan *interfaces.TurnInfoResponse),
		configureSuccessChan:  make(chan *interfaces.ConfigureSuccessResponse),
		configureFailureChan:  make(chan *interfaces.ConfigureFailureResponse),
		fatalErrorChan:        make(chan *interfaces.FatalErrorResponse),
		closeChan:             make(chan *interfaces.CloseResponse),
		errorChan:             make(chan *interfaces.ErrorResponse),
		unhandledChan:         make(chan *[]byte),
	}

	go func() {
		if err := handler.Run(); err != nil {
			klog.V(1).Infof("DefaultChanHandler.Run failed. Err: %v\n", err)
		}
	}()

	return handler
}

func (h DefaultChanHandler) GetOpen() []*chan *interfaces.OpenResponse {
	return []*chan *interfaces.OpenResponse{&h.openChan}
}

func (h DefaultChanHandler) GetConnected() []*chan *interfaces.ConnectedResponse {
	return []*chan *interfaces.ConnectedResponse{&h.connectedChan}
}

func (h DefaultChanHandler) GetTurnInfo() []*chan *interfaces.TurnInfoResponse {
	return []*chan *interfaces.TurnInfoResponse{&h.turnInfoChan}
}

func (h DefaultChanHandler) GetConfigureSuccess() []*chan *interfaces.ConfigureSuccessResponse {
	return []*chan *interfaces.ConfigureSuccessResponse{&h.configureSuccessChan}
}

func (h DefaultChanHandler) GetConfigureFailure() []*chan *interfaces.ConfigureFailureResponse {
	return []*chan *interfaces.ConfigureFailureResponse{&h.configureFailureChan}
}

func (h DefaultChanHandler) GetFatalError() []*chan *interfaces.FatalErrorResponse {
	return []*chan *interfaces.FatalErrorResponse{&h.fatalErrorChan}
}

func (h DefaultChanHandler) GetClose() []*chan *interfaces.CloseResponse {
	return []*chan *interfaces.CloseResponse{&h.closeChan}
}

func (h DefaultChanHandler) GetError() []*chan *interfaces.ErrorResponse {
	return []*chan *interfaces.ErrorResponse{&h.errorChan}
}

func (h DefaultChanHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandledChan}
}

// Run starts goroutines for each channel and prints events to stdout.
// Blocks until all channels are closed.
func (h DefaultChanHandler) Run() error { //nolint:funlen,gocyclo
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for or := range h.openChan {
			if h.debugWebsocket {
				data, _ := json.Marshal(or)
				if pretty, err := prettyjson.Format(data); err == nil {
					klog.V(2).Infof("\n\nOpen Object:\n%s\n\n", pretty)
				}
			}
			fmt.Printf("\n\n[Open] Received\n\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cr := range h.connectedChan {
			if h.debugWebsocket {
				data, _ := json.Marshal(cr)
				if pretty, err := prettyjson.Format(data); err == nil {
					klog.V(2).Infof("\n\nConnected Object:\n%s\n\n", pretty)
				}
				continue
			}
			fmt.Printf("\n[Connected] request_id=%s sequence_id=%d\n", cr.RequestID, cr.SequenceID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for tr := range h.turnInfoChan {
			if h.debugWebsocket {
				data, _ := json.Marshal(tr)
				if pretty, err := prettyjson.Format(data); err == nil {
					klog.V(2).Infof("\n\nTurnInfo Object:\n%s\n\n", pretty)
				}
				continue
			}
			switch tr.EventType {
			case interfaces.TurnEventEndOfTurn, interfaces.TurnEventEagerEndOfTurn:
				fmt.Printf("\n[TurnInfo] event_type=%s turn_index=%d\n  transcript: %s\n",
					tr.EventType, tr.TurnIndex, strings.TrimSpace(tr.Transcript))
			default:
				fmt.Printf("\n[TurnInfo] event_type=%s turn_index=%d\n", tr.EventType, tr.TurnIndex)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range h.configureSuccessChan {
			fmt.Printf("\n[ConfigureSuccess] Received\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cf := range h.configureFailureChan {
			fmt.Printf("\n[ConfigureFailure] request_id=%s sequence_id=%d\n", cf.RequestID, cf.SequenceID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for fe := range h.fatalErrorChan {
			fmt.Printf("\n[FatalError] code=%s description=%s\n", fe.Code, fe.Description)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range h.closeChan {
			fmt.Printf("\n\n[Close] Received\n\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for er := range h.errorChan {
			if h.debugWebsocket {
				data, _ := json.Marshal(er)
				if pretty, err := prettyjson.Format(data); err == nil {
					klog.V(2).Infof("\n\nError Object:\n%s\n\n", pretty)
				}
				continue
			}
			fmt.Printf("\n[Error] type=%s message=%s description=%s\n", er.ErrCode, er.ErrMsg, er.Description)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for byData := range h.unhandledChan {
			if pretty, err := prettyjson.Format(*byData); err == nil {
				klog.V(2).Infof("\n\nUnhandled Object:\n%s\n\n", pretty)
			}
			fmt.Printf("\n[UnhandledEvent] Dump:\n%s\n\n", string(*byData))
		}
	}()

	wg.Wait()
	return nil
}
