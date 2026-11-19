// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	klog "k8s.io/klog/v2"
)

var (
	handlerLogFlagsOnce sync.Once
	handlerLogMu        sync.Mutex
)

func Test_DefaultChanHandlerDoesNotLogPayloads(t *testing.T) {
	const secret = "agent-default-handler-secret"

	handler := newTestDefaultChanHandler()
	done := make(chan error)
	output := captureHandlerOutput(t, func() {
		go func() {
			done <- handler.Run()
		}()

		handler.conversationTextResponse <- &interfaces.ConversationTextResponse{
			Type:    "ConversationText",
			Content: secret,
		}
		unhandled := []byte(`{"type":"FutureEvent","token":"` + secret + `"}`)
		handler.unhandledChan <- &unhandled
		handler.functionCallCancelledResponse <- &interfaces.FunctionCallCancelledResponse{
			Functions: []interfaces.FunctionCallCancelled{{Name: secret}},
		}
		handler.errorChan <- &interfaces.ErrorResponse{
			ErrCode:     "TEST_ERROR",
			ErrMsg:      secret,
			Description: secret,
			Variant:     secret,
		}
		closeTestDefaultChanHandler(handler)
		if err := <-done; err != nil {
			t.Fatalf("run default channel handler: %v", err)
		}

		NewChanRouter(nil).printDebugMessages(2, "ConversationText", []byte(`{"content":"`+secret+`"}`))
	})
	if strings.Contains(output, secret) {
		t.Fatal("default channel handler leaked a payload to log or stdout output")
	}
}

func newTestDefaultChanHandler() *DefaultChanHandler {
	return &DefaultChanHandler{
		debugWebsocket:                true,
		binaryChan:                    make(chan *[]byte),
		openChan:                      make(chan *interfaces.OpenResponse),
		welcomeResponse:               make(chan *interfaces.WelcomeResponse),
		conversationTextResponse:      make(chan *interfaces.ConversationTextResponse),
		userStartedSpeakingResponse:   make(chan *interfaces.UserStartedSpeakingResponse),
		agentThinkingResponse:         make(chan *interfaces.AgentThinkingResponse),
		frustrationDetectedResponse:   make(chan *interfaces.FrustrationDetectedResponse),
		functionCallRequestResponse:   make(chan *interfaces.FunctionCallRequestResponse),
		functionCallCancelledResponse: make(chan *interfaces.FunctionCallCancelledResponse),
		agentStartedSpeakingResponse:  make(chan *interfaces.AgentStartedSpeakingResponse),
		agentAudioDoneResponse:        make(chan *interfaces.AgentAudioDoneResponse),
		injectionRefusedResponse:      make(chan *interfaces.InjectionRefusedResponse),
		keepAliveResponse:             make(chan *interfaces.KeepAlive),
		settingsAppliedResponse:       make(chan *interfaces.SettingsAppliedResponse),
		listenUpdatedResponse:         make(chan *interfaces.ListenUpdatedResponse),
		thinkUpdatedResponse:          make(chan *interfaces.ThinkUpdatedResponse),
		thinkLiveUpdatedResponse:      make(chan *interfaces.ThinkLiveUpdatedResponse),
		warningResponse:               make(chan *interfaces.WarningResponse),
		closeChan:                     make(chan *interfaces.CloseResponse),
		errorChan:                     make(chan *interfaces.ErrorResponse),
		unhandledChan:                 make(chan *[]byte),
	}
}

func closeTestDefaultChanHandler(handler *DefaultChanHandler) {
	close(handler.binaryChan)
	close(handler.openChan)
	close(handler.welcomeResponse)
	close(handler.conversationTextResponse)
	close(handler.userStartedSpeakingResponse)
	close(handler.agentThinkingResponse)
	close(handler.frustrationDetectedResponse)
	close(handler.functionCallRequestResponse)
	close(handler.functionCallCancelledResponse)
	close(handler.agentStartedSpeakingResponse)
	close(handler.agentAudioDoneResponse)
	close(handler.injectionRefusedResponse)
	close(handler.keepAliveResponse)
	close(handler.settingsAppliedResponse)
	close(handler.listenUpdatedResponse)
	close(handler.thinkUpdatedResponse)
	close(handler.thinkLiveUpdatedResponse)
	close(handler.warningResponse)
	close(handler.closeChan)
	close(handler.errorChan)
	close(handler.unhandledChan)
}

func captureHandlerOutput(t *testing.T, action func()) string {
	t.Helper()

	handlerLogMu.Lock()
	defer handlerLogMu.Unlock()
	handlerLogFlagsOnce.Do(func() {
		if flag.Lookup("v") == nil {
			klog.InitFlags(nil)
		}
	})

	verbosity := flag.Lookup("v")
	logToStderr := flag.Lookup("logtostderr")
	if verbosity == nil || logToStderr == nil {
		t.Fatal("klog flags are not registered")
	}
	previousVerbosity := verbosity.Value.String()
	previousLogToStderr := logToStderr.Value.String()
	previousStdout := os.Stdout
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	var stdout bytes.Buffer
	stdoutDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&stdout, stdoutReader)
		close(stdoutDone)
	}()
	var logs bytes.Buffer

	if err := flag.Set("v", "6"); err != nil {
		t.Fatalf("set klog verbosity: %v", err)
	}
	if err := flag.Set("logtostderr", "false"); err != nil {
		t.Fatalf("disable klog stderr output: %v", err)
	}
	os.Stdout = stdoutWriter
	klog.SetOutput(&logs)
	t.Cleanup(func() {
		klog.Flush()
		klog.SetOutput(os.Stderr)
		os.Stdout = previousStdout
		_ = flag.Set("v", previousVerbosity)
		_ = flag.Set("logtostderr", previousLogToStderr)
	})

	action()
	klog.Flush()
	if err := stdoutWriter.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}
	<-stdoutDone
	if err := stdoutReader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}

	return logs.String() + stdout.String()
}
