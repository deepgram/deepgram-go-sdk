// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

var (
	logCaptureFlagsOnce sync.Once
	logCaptureMu        sync.Mutex
)

func Test_StartDoesNotLogConfigurationSettings(t *testing.T) {
	const secret = "agent-settings-log-secret"

	settings := interfaces.NewSettingsConfigurationOptions()
	settings.Agent.Think.Prompt = secret
	client := &WSChannel{tOptions: settings}
	var payload []byte
	var sendErr error
	output := captureLogAndStdout(t, func() {
		sendErr = client.sendConfigurationSettings(func(settings interface{}) error {
			payload, sendErr = json.Marshal(settings)
			return sendErr
		})
	})
	if sendErr != nil {
		t.Fatalf("send configuration settings: %v", sendErr)
	}
	if !bytes.Contains(payload, []byte(secret)) {
		t.Fatal("settings payload was not passed to the websocket writer")
	}
	if strings.Contains(output, secret) {
		t.Fatal("settings payload leaked to log or stdout output")
	}
}

func captureLogAndStdout(t *testing.T, action func()) string {
	t.Helper()

	logCaptureMu.Lock()
	defer logCaptureMu.Unlock()
	logCaptureFlagsOnce.Do(func() {
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
