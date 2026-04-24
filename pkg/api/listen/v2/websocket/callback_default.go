// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
)

// NewDefaultCallbackHandler creates a DefaultCallbackHandler that prints all Flux events to stdout.
func NewDefaultCallbackHandler() *DefaultCallbackHandler {
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
	return &DefaultCallbackHandler{
		debugWebsocket:        strings.EqualFold(debugStr, "true"),
		debugWebsocketVerbose: strings.EqualFold(debugExtStr, "true"),
	}
}

func (h DefaultCallbackHandler) Open(or *interfaces.OpenResponse) error {
	if h.debugWebsocket {
		data, err := json.Marshal(or)
		if err != nil {
			return err
		}
		prettyJSON, err := prettyjson.Format(data)
		if err != nil {
			return err
		}
		klog.V(2).Infof("\n\nOpen Object:\n%s\n\n", prettyJSON)
		return nil
	}
	fmt.Printf("\n\n[Open] Received\n\n")
	return nil
}

func (h DefaultCallbackHandler) Connected(cr *interfaces.ConnectedResponse) error {
	if h.debugWebsocket {
		data, err := json.Marshal(cr)
		if err != nil {
			return err
		}
		prettyJSON, err := prettyjson.Format(data)
		if err != nil {
			return err
		}
		klog.V(2).Infof("\n\nConnected Object:\n%s\n\n", prettyJSON)
		return nil
	}
	fmt.Printf("\n[Connected] request_id=%s sequence_id=%d\n", cr.RequestID, cr.SequenceID)
	return nil
}

func (h DefaultCallbackHandler) TurnInfo(tr *interfaces.TurnInfoResponse) error {
	if h.debugWebsocket {
		data, err := json.Marshal(tr)
		if err != nil {
			return err
		}
		prettyJSON, err := prettyjson.Format(data)
		if err != nil {
			return err
		}
		klog.V(2).Infof("\n\nTurnInfo Object:\n%s\n\n", prettyJSON)
		return nil
	}

	switch tr.EventType {
	case interfaces.TurnEventEndOfTurn:
		fmt.Printf("\n[TurnInfo] event_type=%s turn_index=%d\n  transcript: %s\n", tr.EventType, tr.TurnIndex, strings.TrimSpace(tr.Transcript))
	case interfaces.TurnEventEagerEndOfTurn:
		fmt.Printf("\n[TurnInfo] event_type=%s turn_index=%d\n  transcript: %s\n", tr.EventType, tr.TurnIndex, strings.TrimSpace(tr.Transcript))
	default:
		fmt.Printf("\n[TurnInfo] event_type=%s turn_index=%d\n", tr.EventType, tr.TurnIndex)
		if h.debugWebsocketVerbose && len(tr.Transcript) > 0 {
			fmt.Printf("  transcript: %s\n", strings.TrimSpace(tr.Transcript))
		}
	}
	return nil
}

func (h DefaultCallbackHandler) ConfigureSuccess(cs *interfaces.ConfigureSuccessResponse) error {
	fmt.Printf("\n[ConfigureSuccess] Received\n")
	return nil
}

func (h DefaultCallbackHandler) ConfigureFailure(cf *interfaces.ConfigureFailureResponse) error {
	fmt.Printf("\n[ConfigureFailure] request_id=%s sequence_id=%d\n", cf.RequestID, cf.SequenceID)
	return nil
}

func (h DefaultCallbackHandler) FatalError(fe *interfaces.FatalErrorResponse) error {
	fmt.Printf("\n[FatalError] code=%s description=%s\n", fe.Code, fe.Description)
	return nil
}

func (h DefaultCallbackHandler) Close(cr *interfaces.CloseResponse) error {
	if h.debugWebsocket {
		data, err := json.Marshal(cr)
		if err != nil {
			return err
		}
		prettyJSON, err := prettyjson.Format(data)
		if err != nil {
			return err
		}
		klog.V(2).Infof("\n\nClose Object:\n%s\n\n", prettyJSON)
		return nil
	}
	fmt.Printf("\n\n[Close] Received\n\n")
	return nil
}

func (h DefaultCallbackHandler) Error(er *interfaces.ErrorResponse) error {
	if h.debugWebsocket {
		data, err := json.Marshal(er)
		if err != nil {
			return err
		}
		prettyJSON, err := prettyjson.Format(data)
		if err != nil {
			return err
		}
		klog.V(2).Infof("\n\nError Object:\n%s\n\n", prettyJSON)
		return nil
	}
	fmt.Printf("\n[Error]\n")
	fmt.Printf("Error.Type: %s\n", er.ErrCode)
	fmt.Printf("Error.Message: %s\n", er.ErrMsg)
	fmt.Printf("Error.Description: %s\n\n", er.Description)
	return nil
}

func (h DefaultCallbackHandler) UnhandledEvent(byData []byte) error {
	if h.debugWebsocket {
		prettyJSON, err := prettyjson.Format(byData)
		if err != nil {
			klog.V(2).Infof("\n\nRaw Data:\n%s\n\n", string(byData))
		} else {
			klog.V(2).Infof("\n\nUnhandled Object:\n%s\n\n", prettyJSON)
		}
		return nil
	}
	fmt.Printf("\n[UnhandledEvent] Dump:\n%s\n\n", string(byData))
	return nil
}
