// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"fmt"
	"os"
	"strings"
	"sync"

	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
)

// NewDefaultChanHandler creates a DefaultChanHandler with channels for all Flux TTS
// events. The handler's Run() goroutine prints events to stdout; it is started
// automatically. Binary audio frames are acknowledged but discarded.
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
		binaryChan:            make(chan *[]byte),
		speechStartedChan:     make(chan *interfaces.SpeechStartedResponse),
		speechMetadataChan:    make(chan *interfaces.SpeechMetadataResponse),
		speechInterruptedChan: make(chan *interfaces.SpeechInterruptedResponse),
		flushedChan:           make(chan *interfaces.FlushedResponse),
		sessionMetadataChan:   make(chan *interfaces.SessionMetadataResponse),
		configureSuccessChan:  make(chan *interfaces.ConfigureSuccessResponse),
		configureFailureChan:  make(chan *interfaces.ConfigureFailureResponse),
		warningChan:           make(chan *interfaces.WarningResponse),
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

func (h *DefaultChanHandler) GetOpen() []*chan *interfaces.OpenResponse {
	return []*chan *interfaces.OpenResponse{&h.openChan}
}

func (h *DefaultChanHandler) GetConnected() []*chan *interfaces.ConnectedResponse {
	return []*chan *interfaces.ConnectedResponse{&h.connectedChan}
}

func (h *DefaultChanHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&h.binaryChan}
}

func (h *DefaultChanHandler) GetSpeechStarted() []*chan *interfaces.SpeechStartedResponse {
	return []*chan *interfaces.SpeechStartedResponse{&h.speechStartedChan}
}

func (h *DefaultChanHandler) GetSpeechMetadata() []*chan *interfaces.SpeechMetadataResponse {
	return []*chan *interfaces.SpeechMetadataResponse{&h.speechMetadataChan}
}

func (h *DefaultChanHandler) GetSpeechInterrupted() []*chan *interfaces.SpeechInterruptedResponse {
	return []*chan *interfaces.SpeechInterruptedResponse{&h.speechInterruptedChan}
}

func (h *DefaultChanHandler) GetFlushed() []*chan *interfaces.FlushedResponse {
	return []*chan *interfaces.FlushedResponse{&h.flushedChan}
}

func (h *DefaultChanHandler) GetSessionMetadata() []*chan *interfaces.SessionMetadataResponse {
	return []*chan *interfaces.SessionMetadataResponse{&h.sessionMetadataChan}
}

func (h *DefaultChanHandler) GetConfigureSuccess() []*chan *interfaces.ConfigureSuccessResponse {
	return []*chan *interfaces.ConfigureSuccessResponse{&h.configureSuccessChan}
}

func (h *DefaultChanHandler) GetConfigureFailure() []*chan *interfaces.ConfigureFailureResponse {
	return []*chan *interfaces.ConfigureFailureResponse{&h.configureFailureChan}
}

func (h *DefaultChanHandler) GetWarning() []*chan *interfaces.WarningResponse {
	return []*chan *interfaces.WarningResponse{&h.warningChan}
}

func (h *DefaultChanHandler) GetFatalError() []*chan *interfaces.FatalErrorResponse {
	return []*chan *interfaces.FatalErrorResponse{&h.fatalErrorChan}
}

func (h *DefaultChanHandler) GetClose() []*chan *interfaces.CloseResponse {
	return []*chan *interfaces.CloseResponse{&h.closeChan}
}

func (h *DefaultChanHandler) GetError() []*chan *interfaces.ErrorResponse {
	return []*chan *interfaces.ErrorResponse{&h.errorChan}
}

func (h *DefaultChanHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandledChan}
}

// Run starts goroutines for each channel and prints events to stdout.
// Blocks until all channels are closed.
func (h *DefaultChanHandler) Run() error { //nolint:funlen // one goroutine per event channel
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range h.openChan {
			fmt.Printf("\n\n[Open] Received\n\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cr := range h.connectedChan {
			fmt.Printf("\n[Connected] request_id=%s model=%s version=%s\n", cr.RequestID, cr.ModelName, cr.ModelVersion)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for byAudio := range h.binaryChan {
			fmt.Printf("[Binary] %d bytes of audio received\n", len(*byAudio))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for ss := range h.speechStartedChan {
			fmt.Printf("\n[SpeechStarted] speech_id=%s\n", ss.SpeechID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for sm := range h.speechMetadataChan {
			fmt.Printf("\n[SpeechMetadata] speech_id=%s duration_ms=%d billable_chars=%d\n",
				sm.SpeechID, sm.AudioDurationMs, sm.BillableCharacterCount)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for si := range h.speechInterruptedChan {
			fmt.Printf("\n[SpeechInterrupted] audio_played_ms=%d\n", si.AudioPlayedMs)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for fl := range h.flushedChan {
			fmt.Printf("\n[Flushed] speech_id=%s\n", fl.SpeechID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for sm := range h.sessionMetadataChan {
			fmt.Printf("\n[SessionMetadata] total_duration_ms=%d total_billable_chars=%d\n",
				sm.TotalAudioDurationMs, sm.TotalBillableCharacterCount)
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
			fmt.Printf("\n[ConfigureFailure] code=%s description=%s\n", cf.Code, cf.Description)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for wr := range h.warningChan {
			fmt.Printf("\n[Warning] code=%s description=%s\n", wr.Code, wr.Description)
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
			fmt.Printf("\n[Error] type=%s message=%s description=%s\n", er.ErrCode, er.ErrMsg, er.Description)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for byData := range h.unhandledChan {
			fmt.Printf("\n[UnhandledEvent] Dump:\n%s\n\n", string(*byData))
		}
	}()

	wg.Wait()
	return nil
}
