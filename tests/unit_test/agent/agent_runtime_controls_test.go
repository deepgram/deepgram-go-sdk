// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
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
	agentclient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
	clientws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent/v1/websocket"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const (
	listenUpdatedPayload = `{"type":"ListenUpdated"}`
	thinkUpdatedPayload  = `{"type":"ThinkUpdated"}`
	forceEndTurnWarning  = `{"type":"Warning","code":"FORCE_END_TURN_UNSUPPORTED","description":"ForceEndTurn requires a Deepgram V2 listen provider"}`
)

type runtimeControlHandler struct {
	*legacyFunctionCallHandler
	listenUpdated chan *msginterfaces.ListenUpdatedResponse
	thinkUpdated  chan *msginterfaces.ThinkUpdatedResponse
	warnings      chan *msginterfaces.WarningResponse
}

func newRuntimeControlHandler() *runtimeControlHandler {
	return &runtimeControlHandler{
		legacyFunctionCallHandler: newLegacyFunctionCallHandler(),
		listenUpdated:             make(chan *msginterfaces.ListenUpdatedResponse, 1),
		thinkUpdated:              make(chan *msginterfaces.ThinkUpdatedResponse, 1),
		warnings:                  make(chan *msginterfaces.WarningResponse, 1),
	}
}

func (h *runtimeControlHandler) GetListenUpdated() []*chan *msginterfaces.ListenUpdatedResponse {
	return []*chan *msginterfaces.ListenUpdatedResponse{&h.listenUpdated}
}

func (h *runtimeControlHandler) GetThinkUpdated() []*chan *msginterfaces.ThinkUpdatedResponse {
	return []*chan *msginterfaces.ThinkUpdatedResponse{&h.thinkUpdated}
}

func (h *runtimeControlHandler) GetWarning() []*chan *msginterfaces.WarningResponse {
	return []*chan *msginterfaces.WarningResponse{&h.warnings}
}

func Test_AgentRuntimeControlSerialization(t *testing.T) {
	listen := &interfacesv1.Listen{Provider: map[string]interface{}{
		"type":          "deepgram",
		"version":       "v2",
		"model":         "flux-general-en",
		"eot_threshold": 1.0,
	}}
	think := &interfacesv1.Think{Provider: map[string]interface{}{
		"type":  "open_ai",
		"model": "gpt-4o-mini",
	}, Prompt: "Answer concisely."}

	tests := []struct {
		name string
		msg  interface{}
		want string
	}{
		{
			name: "UpdateListen",
			msg:  clientws.UpdateListen{Type: msginterfaces.TypeUpdateListen, Listen: listen},
			want: `{"type":"UpdateListen","listen":{"provider":{"eot_threshold":1,"model":"flux-general-en","type":"deepgram","version":"v2"}}}`,
		},
		{
			name: "UpdateThink",
			msg:  clientws.UpdateThink{Type: msginterfaces.TypeUpdateThink, Think: think},
			want: `{"type":"UpdateThink","think":{"provider":{"model":"gpt-4o-mini","type":"open_ai"},"prompt":"Answer concisely."}}`,
		},
		{
			name: "ForceEndTurn",
			msg:  clientws.ForceEndTurn{Type: msginterfaces.TypeForceEndTurn},
			want: `{"type":"ForceEndTurn"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.msg)
			if err != nil {
				t.Fatalf("marshal %s: %v", tt.name, err)
			}
			if string(got) != tt.want {
				t.Errorf("%s wire format mismatch: got %s, want %s", tt.name, got, tt.want)
			}
		})
	}
}

func Test_AgentRuntimeControlRouting(t *testing.T) {
	handler := newRuntimeControlHandler()
	router := agentws.NewChanRouter(handler)

	if err := router.Message([]byte(listenUpdatedPayload)); err != nil {
		t.Fatalf("ListenUpdated routing failed: %v", err)
	}
	if err := router.Message([]byte(thinkUpdatedPayload)); err != nil {
		t.Fatalf("ThinkUpdated routing failed: %v", err)
	}
	if err := router.Message([]byte(forceEndTurnWarning)); err != nil {
		t.Fatalf("Warning routing failed: %v", err)
	}

	select {
	case event := <-handler.listenUpdated:
		if event.Type != msginterfaces.TypeListenUpdatedResponse {
			t.Errorf("unexpected ListenUpdated type: %q", event.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("ListenUpdated was not delivered")
	}

	select {
	case event := <-handler.thinkUpdated:
		if event.Type != msginterfaces.TypeThinkUpdatedResponse {
			t.Errorf("unexpected ThinkUpdated type: %q", event.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("ThinkUpdated was not delivered")
	}

	select {
	case event := <-handler.warnings:
		if event.Code != "FORCE_END_TURN_UNSUPPORTED" {
			t.Errorf("unexpected warning code: %q", event.Code)
		}
	case <-time.After(time.Second):
		t.Fatal("Warning was not delivered")
	}
}

func Test_AgentRuntimeControlFallsBackToUnhandled(t *testing.T) {
	handler := newLegacyFunctionCallHandler()
	router := agentws.NewChanRouter(handler)

	for _, payload := range []string{listenUpdatedPayload, thinkUpdatedPayload, forceEndTurnWarning} {
		if err := router.Message([]byte(payload)); err != nil {
			t.Fatalf("legacy handler must safely fall back for %s: %v", payload, err)
		}
		select {
		case message := <-handler.unhandled:
			if string(*message) != payload {
				t.Errorf("unexpected fallback message: %s", *message)
			}
		case <-time.After(time.Second):
			t.Fatalf("legacy handler did not receive %s", payload)
		}
	}
}

func Test_AgentRuntimeControlValidation(t *testing.T) {
	var client clientws.WSChannel

	if err := client.UpdateListen(nil); !errors.Is(err, clientws.ErrInvalidUpdateListen) {
		t.Errorf("UpdateListen(nil) error = %v, want ErrInvalidUpdateListen", err)
	}
	if err := client.UpdateListen(&interfacesv1.Listen{}); !errors.Is(err, clientws.ErrInvalidUpdateListen) {
		t.Errorf("UpdateListen without provider error = %v, want ErrInvalidUpdateListen", err)
	}
	if err := client.UpdateThink(nil); !errors.Is(err, clientws.ErrInvalidUpdateThink) {
		t.Errorf("UpdateThink(nil) error = %v, want ErrInvalidUpdateThink", err)
	}
	if err := client.UpdateThink(&interfacesv1.Think{}); !errors.Is(err, clientws.ErrInvalidUpdateThink) {
		t.Errorf("UpdateThink without provider error = %v, want ErrInvalidUpdateThink", err)
	}
}

func newAgentRuntimeControlTestServer(t *testing.T) (string, chan string, func()) {
	t.Helper()

	frames := make(chan string, 8)
	upgrader := websocket.Upgrader{}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("test server upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType != websocket.TextMessage {
				continue
			}

			frame := string(message)
			switch {
			case strings.Contains(frame, `"type":"UpdateListen"`):
				frames <- frame
				_ = conn.WriteMessage(websocket.TextMessage, []byte(listenUpdatedPayload))
			case strings.Contains(frame, `"type":"UpdateThink"`):
				frames <- frame
				_ = conn.WriteMessage(websocket.TextMessage, []byte(thinkUpdatedPayload))
			case strings.Contains(frame, `"type":"ForceEndTurn"`):
				frames <- frame
				_ = conn.WriteMessage(websocket.TextMessage, []byte(forceEndTurnWarning))
			}
		}
	}))

	return strings.TrimPrefix(srv.URL, "https://"), frames, srv.Close
}

func newAgentRuntimeControlClient(t *testing.T, host string, handler *runtimeControlHandler) *clientws.WSChannel {
	t.Helper()

	client, err := agentclient.NewWSUsingChan(
		context.Background(),
		"",
		&interfacesv1.ClientOptions{Host: host, SelfHosted: true, SkipServerAuth: true},
		&interfacesv1.SettingsOptions{},
		handler,
	)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %v", err)
	}
	if !client.Connect() {
		t.Fatal("Connect to the test server failed")
	}

	return client
}

func awaitAgentRuntimeControlFrame(t *testing.T, frames chan string, want string) {
	t.Helper()

	select {
	case frame := <-frames:
		if frame != want {
			t.Errorf("runtime control frame = %s, want %s", frame, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("server did not receive runtime control frame %s", want)
	}
}

func Test_AgentRuntimeControlsWireChannel(t *testing.T) {
	host, frames, shutdown := newAgentRuntimeControlTestServer(t)
	defer shutdown()

	handler := newRuntimeControlHandler()
	client := newAgentRuntimeControlClient(t, host, handler)

	listen := &interfacesv1.Listen{Provider: map[string]interface{}{
		"type":          "deepgram",
		"version":       "v2",
		"model":         "flux-general-en",
		"eot_threshold": 1.0,
	}}
	if err := client.UpdateListen(listen); err != nil {
		t.Fatalf("UpdateListen failed: %v", err)
	}
	awaitAgentRuntimeControlFrame(t, frames, `{"type":"UpdateListen","listen":{"provider":{"eot_threshold":1,"model":"flux-general-en","type":"deepgram","version":"v2"}}}`)

	select {
	case event := <-handler.listenUpdated:
		if event.Type != msginterfaces.TypeListenUpdatedResponse {
			t.Errorf("ListenUpdated type = %q, want %q", event.Type, msginterfaces.TypeListenUpdatedResponse)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ListenUpdated was not routed through the connected client")
	}

	think := &interfacesv1.Think{Provider: map[string]interface{}{
		"type":  "open_ai",
		"model": "gpt-4o-mini",
	}, Prompt: "Answer concisely."}
	if err := client.UpdateThink(think); err != nil {
		t.Fatalf("UpdateThink failed: %v", err)
	}
	awaitAgentRuntimeControlFrame(t, frames, `{"type":"UpdateThink","think":{"provider":{"model":"gpt-4o-mini","type":"open_ai"},"prompt":"Answer concisely."}}`)

	select {
	case event := <-handler.thinkUpdated:
		if event.Type != msginterfaces.TypeThinkUpdatedResponse {
			t.Errorf("ThinkUpdated type = %q, want %q", event.Type, msginterfaces.TypeThinkUpdatedResponse)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ThinkUpdated was not routed through the connected client")
	}

	if err := client.ForceEndTurn(); err != nil {
		t.Fatalf("ForceEndTurn failed: %v", err)
	}
	awaitAgentRuntimeControlFrame(t, frames, `{"type":"ForceEndTurn"}`)

	select {
	case event := <-handler.warnings:
		if event.Code != "FORCE_END_TURN_UNSUPPORTED" {
			t.Errorf("warning code = %q, want FORCE_END_TURN_UNSUPPORTED", event.Code)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Warning was not routed through the connected client")
	}

	client.Stop()
	if err := client.UpdateListen(listen); err == nil {
		t.Error("UpdateListen after Stop returned nil")
	}
	if err := client.UpdateThink(think); err == nil {
		t.Error("UpdateThink after Stop returned nil")
	}
	if err := client.ForceEndTurn(); err == nil {
		t.Error("ForceEndTurn after Stop returned nil")
	}
}

type runtimeControlMarshalFailure string

func (v runtimeControlMarshalFailure) MarshalJSON() ([]byte, error) {
	return nil, errors.New(string(v))
}

func Test_AgentUpdateThinkDoesNotLogPayload(t *testing.T) {
	klog.InitFlags(nil)
	verbosity := flag.Lookup("v")
	if verbosity == nil {
		t.Fatal("klog verbosity flag is not initialized")
	}
	previousVerbosity := verbosity.Value.String()
	if err := flag.Set("v", "6"); err != nil {
		t.Fatalf("set klog verbosity: %v", err)
	}

	var logs bytes.Buffer
	klog.LogToStderr(false)
	klog.SetOutput(&logs)
	t.Cleanup(func() {
		klog.Flush()
		klog.SetOutput(os.Stderr)
		klog.LogToStderr(true)
		if err := flag.Set("v", previousVerbosity); err != nil {
			t.Errorf("restore klog verbosity: %v", err)
		}
	})

	host, _, shutdown := newAgentRuntimeControlTestServer(t)
	defer shutdown()

	client := newAgentRuntimeControlClient(t, host, newRuntimeControlHandler())
	defer client.Stop()

	const sentinel = "agent-runtime-controls-authorization-header-sentinel"
	if err := client.UpdateThink(&interfacesv1.Think{
		Provider: map[string]interface{}{"type": "open_ai"},
		Endpoint: &interfacesv1.Endpoint{Headers: map[string]string{"authorization": sentinel}},
	}); err != nil {
		t.Fatalf("UpdateThink with endpoint headers failed: %v", err)
	}
	if err := client.UpdateThink(&interfacesv1.Think{
		Provider: map[string]interface{}{"authorization": runtimeControlMarshalFailure(sentinel)},
	}); err == nil {
		t.Fatal("UpdateThink with an unmarshalable payload returned nil")
	}

	klog.Flush()
	if strings.Contains(logs.String(), sentinel) {
		t.Fatal("UpdateThink leaked its payload to logs")
	}
}
