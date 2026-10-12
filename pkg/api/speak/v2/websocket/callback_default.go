// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
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

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
)

// NewDefaultCallbackHandler creates a DefaultCallbackHandler that prints all Flux TTS
// events to stdout. Binary audio frames are acknowledged but discarded.
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

// printDebug renders an event as pretty JSON to the debug log. It returns true
// when debug output was produced and the caller should skip the stdout line.
func (h DefaultCallbackHandler) printDebug(name string, v interface{}) bool {
	if !h.debugWebsocket {
		return false
	}
	data, err := json.Marshal(v)
	if err != nil {
		return false
	}
	prettyJSON, err := prettyjson.Format(data)
	if err != nil {
		return false
	}
	klog.V(2).Infof("\n\n%s Object:\n%s\n\n", name, prettyJSON)
	return true
}

func (h DefaultCallbackHandler) Open(or *interfaces.OpenResponse) error {
	if h.printDebug("Open", or) {
		return nil
	}
	fmt.Printf("\n\n[Open] Received\n\n")
	return nil
}

func (h DefaultCallbackHandler) Connected(cr *interfaces.ConnectedResponse) error {
	if h.printDebug("Connected", cr) {
		return nil
	}
	fmt.Printf("\n[Connected] request_id=%s model=%s version=%s\n", cr.RequestID, cr.ModelName, cr.ModelVersion)
	return nil
}

func (h DefaultCallbackHandler) Binary(byMsg []byte) error {
	fmt.Printf("[Binary] %d bytes of audio received\n", len(byMsg))
	return nil
}

func (h DefaultCallbackHandler) SpeechStarted(ss *interfaces.SpeechStartedResponse) error {
	if h.printDebug("SpeechStarted", ss) {
		return nil
	}
	fmt.Printf("\n[SpeechStarted] speech_id=%s\n", ss.SpeechID)
	return nil
}

func (h DefaultCallbackHandler) SpeechMetadata(sm *interfaces.SpeechMetadataResponse) error {
	if h.printDebug("SpeechMetadata", sm) {
		return nil
	}
	fmt.Printf("\n[SpeechMetadata] speech_id=%s duration_ms=%d billable_chars=%d\n",
		sm.SpeechID, sm.AudioDurationMs, sm.BillableCharacterCount)
	return nil
}

func (h DefaultCallbackHandler) SpeechInterrupted(si *interfaces.SpeechInterruptedResponse) error {
	if h.printDebug("SpeechInterrupted", si) {
		return nil
	}
	fmt.Printf("\n[SpeechInterrupted] audio_played_ms=%d\n", si.AudioPlayedMs)
	return nil
}

func (h DefaultCallbackHandler) Flushed(fl *interfaces.FlushedResponse) error {
	if h.printDebug("Flushed", fl) {
		return nil
	}
	fmt.Printf("\n[Flushed] speech_id=%s\n", fl.SpeechID)
	return nil
}

func (h DefaultCallbackHandler) SessionMetadata(sm *interfaces.SessionMetadataResponse) error {
	if h.printDebug("SessionMetadata", sm) {
		return nil
	}
	fmt.Printf("\n[SessionMetadata] total_duration_ms=%d total_billable_chars=%d\n",
		sm.TotalAudioDurationMs, sm.TotalBillableCharacterCount)
	return nil
}

func (h DefaultCallbackHandler) ConfigureSuccess(cs *interfaces.ConfigureSuccessResponse) error {
	if h.printDebug("ConfigureSuccess", cs) {
		return nil
	}
	fmt.Printf("\n[ConfigureSuccess] Received\n")
	return nil
}

func (h DefaultCallbackHandler) ConfigureFailure(cf *interfaces.ConfigureFailureResponse) error {
	if h.printDebug("ConfigureFailure", cf) {
		return nil
	}
	fmt.Printf("\n[ConfigureFailure] code=%s description=%s\n", cf.Code, cf.Description)
	return nil
}

func (h DefaultCallbackHandler) Warning(wr *interfaces.WarningResponse) error {
	if h.printDebug("Warning", wr) {
		return nil
	}
	fmt.Printf("\n[Warning] code=%s description=%s\n", wr.Code, wr.Description)
	return nil
}

func (h DefaultCallbackHandler) FatalError(fe *interfaces.FatalErrorResponse) error {
	if h.printDebug("FatalError", fe) {
		return nil
	}
	fmt.Printf("\n[FatalError] code=%s description=%s\n", fe.Code, fe.Description)
	return nil
}

func (h DefaultCallbackHandler) Close(cr *interfaces.CloseResponse) error {
	if h.printDebug("Close", cr) {
		return nil
	}
	fmt.Printf("\n\n[Close] Received\n\n")
	return nil
}

func (h DefaultCallbackHandler) Error(er *interfaces.ErrorResponse) error {
	if h.printDebug("Error", er) {
		return nil
	}
	fmt.Printf("\n[Error] type=%s message=%s description=%s\n", er.ErrCode, er.ErrMsg, er.Description)
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
