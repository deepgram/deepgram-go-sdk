// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	clientws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent/v1/websocket"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const (
	agentTestAPIKey                = "test-api-key"
	customFromThinkProviderPayload = `{"type":"__customFromThinkProvider","content":{"event":"ready","turn":1}}`
)

var _ msginterfaces.CustomFromThinkProviderChan = (*agentws.DefaultChanHandler)(nil)

type customThinkProviderHandler struct {
	*legacyFunctionCallHandler
	customFromThinkProvider chan *msginterfaces.CustomFromThinkProviderResponse
}

func newCustomThinkProviderHandler() *customThinkProviderHandler {
	return &customThinkProviderHandler{
		legacyFunctionCallHandler: newLegacyFunctionCallHandler(),
		customFromThinkProvider:   make(chan *msginterfaces.CustomFromThinkProviderResponse, 1),
	}
}

func (h *customThinkProviderHandler) GetCustomFromThinkProvider() []*chan *msginterfaces.CustomFromThinkProviderResponse {
	return []*chan *msginterfaces.CustomFromThinkProviderResponse{&h.customFromThinkProvider}
}

func Test_CustomToThinkProviderMarshaling(t *testing.T) {
	request := clientws.CustomToThinkProvider{
		Type:    msginterfaces.TypeCustomToThinkProvider,
		Content: json.RawMessage(`{"event":"start","values":[1,true,null]}`),
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal CustomToThinkProvider: %v", err)
	}

	const expected = `{"type":"__customToThinkProvider","content":{"event":"start","values":[1,true,null]}}`
	if string(data) != expected {
		t.Errorf("unexpected CustomToThinkProvider payload:\nexpected: %s\nactual:   %s", expected, data)
	}

	data, err = json.Marshal(msginterfaces.CustomToThinkProvider{
		Type:    msginterfaces.TypeCustomToThinkProvider,
		Content: nil,
	})
	if err != nil {
		t.Fatalf("marshal CustomToThinkProvider null content: %v", err)
	}
	if string(data) != `{"type":"__customToThinkProvider","content":null}` {
		t.Errorf("null content must be preserved as a valid JSON value: %s", data)
	}
}

func Test_CustomToThinkProviderValidation(t *testing.T) {
	conn := &clientws.WSChannel{}
	if err := conn.SendCustomToThinkProvider(json.RawMessage(`not-json`)); err == nil {
		t.Fatal("invalid custom Think provider content must be rejected before writing")
	}
}

func newAgentWireTestServer(t *testing.T) (host string, frames chan string, shutdown func()) {
	t.Helper()

	frames = make(chan string, 16)
	upgrader := websocket.Upgrader{}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("test server upgrade failed: %s", err)
			return
		}
		defer conn.Close()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if msgType == websocket.TextMessage {
				frames <- string(msg)
			}
		}
	}))

	return strings.TrimPrefix(srv.URL, "https://"), frames, srv.Close
}

func newAgentWireClient(t *testing.T, host string) *agent.WSChannel {
	t.Helper()
	t.Setenv("DEEPGRAM_HOST", "")

	conn, err := agent.NewWSUsingChan(
		context.Background(),
		agentTestAPIKey,
		&interfaces.ClientOptions{Host: host, SkipServerAuth: true},
		agent.NewSettingsConfigurationOptions(),
		newLegacyFunctionCallHandler(),
	)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %v", err)
	}
	if !conn.Connect() {
		t.Fatal("Connect to the test server failed")
	}
	t.Cleanup(conn.Stop)

	return conn
}

func awaitAgentWireFrame(t *testing.T, frames <-chan string) string {
	t.Helper()

	select {
	case frame := <-frames:
		return frame
	case <-time.After(5 * time.Second):
		t.Fatal("server never received a text frame")
		return ""
	}
}

func Test_SendCustomToThinkProviderWirePayloads(t *testing.T) {
	host, frames, shutdown := newAgentWireTestServer(t)
	defer shutdown()

	conn := newAgentWireClient(t, host)
	awaitAgentWireFrame(t, frames) // Settings are sent as part of Connect().
	for _, test := range []struct {
		name     string
		content  json.RawMessage
		expected string
	}{
		{
			name:     "object",
			content:  json.RawMessage(`{"event":"ready","turn":1}`),
			expected: `{"type":"__customToThinkProvider","content":{"event":"ready","turn":1}}`,
		},
		{
			name:     "string",
			content:  json.RawMessage(`"provider returned plain text"`),
			expected: `{"type":"__customToThinkProvider","content":"provider returned plain text"}`,
		},
		{
			name:     "null",
			content:  json.RawMessage(`null`),
			expected: `{"type":"__customToThinkProvider","content":null}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := conn.SendCustomToThinkProvider(test.content); err != nil {
				t.Fatalf("SendCustomToThinkProvider failed: %v", err)
			}
			if frame := awaitAgentWireFrame(t, frames); frame != test.expected {
				t.Errorf("unexpected custom Think provider frame:\nexpected: %s\nactual:   %s", test.expected, frame)
			}
		})
	}
}

func captureAgentSettingsOutput(t *testing.T, fn func()) string {
	t.Helper()

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	if flag.Lookup("v") == nil {
		klog.InitFlags(nil)
	}
	originalStdout := os.Stdout
	originalVerbosity := flag.Lookup("v").Value.String()
	var logs bytes.Buffer

	os.Stdout = stdoutWriter
	klog.SetOutput(&logs)
	if err := flag.Set("v", "6"); err != nil {
		t.Fatalf("set klog verbosity: %v", err)
	}
	t.Cleanup(func() {
		klog.Flush()
		klog.SetOutput(os.Stderr)
		_ = flag.Set("v", originalVerbosity)
		os.Stdout = originalStdout
		_ = stdoutWriter.Close()
		_ = stdoutReader.Close()
	})

	fn()
	klog.Flush()
	if err := stdoutWriter.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}
	stdout, err := io.ReadAll(stdoutReader)
	if err != nil {
		t.Fatalf("read stdout pipe: %v", err)
	}

	return string(stdout) + logs.String()
}

func Test_AgentStartDoesNotLogSettingsEndpoints(t *testing.T) {
	const endpointToken = "https://think.example.test/custom/endpoint-token"
	const headerToken = "header-secret-token"

	host, frames, shutdown := newAgentWireTestServer(t)
	defer shutdown()
	t.Setenv("DEEPGRAM_HOST", "")
	settings := agent.NewSettingsConfigurationOptions()
	settings.Agent.Think.Endpoint = &interfacesv1.Endpoint{
		Url: endpointToken,
		Headers: map[string]string{
			"Authorization": headerToken,
		},
	}
	conn, err := agent.NewWSUsingChan(
		context.Background(),
		agentTestAPIKey,
		&interfaces.ClientOptions{Host: host, SkipServerAuth: true},
		settings,
		newLegacyFunctionCallHandler(),
	)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %v", err)
	}
	connected := false
	output := captureAgentSettingsOutput(t, func() {
		connected = conn.Connect()
	})
	if !connected {
		t.Fatal("Connect to the test server failed")
	}
	t.Cleanup(conn.Stop)

	if strings.Contains(output, endpointToken) || strings.Contains(output, headerToken) {
		t.Fatalf("settings endpoint or headers leaked to output: %s", output)
	}

	var settingsFrame interfacesv1.SettingsOptions
	if err := json.Unmarshal([]byte(awaitAgentWireFrame(t, frames)), &settingsFrame); err != nil {
		t.Fatalf("unmarshal settings frame: %v", err)
	}
	if settingsFrame.Type != settings.Type ||
		settingsFrame.Agent.Think.Endpoint == nil ||
		settingsFrame.Agent.Think.Endpoint.Url != endpointToken ||
		settingsFrame.Agent.Think.Endpoint.Headers["Authorization"] != headerToken {
		t.Fatalf("settings payload changed: %v", settingsFrame)
	}
}

func Test_CustomFromThinkProviderRouting(t *testing.T) {
	handler := newCustomThinkProviderHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(customFromThinkProviderPayload)); err != nil {
		t.Fatalf("CustomFromThinkProvider must route successfully: %v", err)
	}

	select {
	case response := <-handler.customFromThinkProvider:
		if response.Type != msginterfaces.TypeCustomFromThinkProviderResponse {
			t.Errorf("expected type %q, got %q", msginterfaces.TypeCustomFromThinkProviderResponse, response.Type)
		}
		if string(response.Content) != `{"event":"ready","turn":1}` {
			t.Errorf("unexpected custom provider content: %s", response.Content)
		}
	case <-time.After(time.Second):
		t.Fatal("CustomFromThinkProvider was not delivered to the handler")
	}
}

func Test_CustomFromThinkProviderJSONStringContent(t *testing.T) {
	handler := newCustomThinkProviderHandler()
	router := agentws.NewChanRouter(handler)

	const payload = `{"type":"__customFromThinkProvider","content":"provider returned plain text"}`
	if err := router.Message([]byte(payload)); err != nil {
		t.Fatalf("CustomFromThinkProvider string content must route successfully: %v", err)
	}

	select {
	case response := <-handler.customFromThinkProvider:
		if string(response.Content) != `"provider returned plain text"` {
			t.Errorf("non-JSON provider content must remain a JSON string: %s", response.Content)
		}
	case <-time.After(time.Second):
		t.Fatal("CustomFromThinkProvider string content was not delivered to the handler")
	}
}

func Test_CustomFromThinkProviderFallsBackToUnhandled(t *testing.T) {
	handler := newLegacyFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(customFromThinkProviderPayload)); err != nil {
		t.Fatalf("CustomFromThinkProvider must fall back successfully: %v", err)
	}

	select {
	case message := <-handler.unhandled:
		if string(*message) != customFromThinkProviderPayload {
			t.Errorf("unexpected unhandled message: %s", *message)
		}
	case <-time.After(time.Second):
		t.Fatal("CustomFromThinkProvider was not delivered to the unhandled handler")
	}
}
