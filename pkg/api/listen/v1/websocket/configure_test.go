// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"encoding/json"
	"testing"
	"time"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
)

type configureCallback struct {
	successes []*interfaces.ConfigureSuccessResponse
	failures  []*interfaces.ConfigureFailureResponse
	unhandled [][]byte
}

func (c *configureCallback) Open(*interfaces.OpenResponse) error { return nil }
func (c *configureCallback) Message(*interfaces.MessageResponse) error {
	return nil
}
func (c *configureCallback) Metadata(*interfaces.MetadataResponse) error { return nil }
func (c *configureCallback) SpeechStarted(*interfaces.SpeechStartedResponse) error {
	return nil
}
func (c *configureCallback) UtteranceEnd(*interfaces.UtteranceEndResponse) error { return nil }
func (c *configureCallback) Close(*interfaces.CloseResponse) error               { return nil }
func (c *configureCallback) Error(*interfaces.ErrorResponse) error               { return nil }
func (c *configureCallback) UnhandledEvent(data []byte) error {
	c.unhandled = append(c.unhandled, data)
	return nil
}
func (c *configureCallback) ConfigureSuccess(msg *interfaces.ConfigureSuccessResponse) error {
	c.successes = append(c.successes, msg)
	return nil
}
func (c *configureCallback) ConfigureFailure(msg *interfaces.ConfigureFailureResponse) error {
	c.failures = append(c.failures, msg)
	return nil
}

type legacyConfigureCallback struct {
	unhandled [][]byte
}

func (c *legacyConfigureCallback) Open(*interfaces.OpenResponse) error { return nil }
func (c *legacyConfigureCallback) Message(*interfaces.MessageResponse) error {
	return nil
}
func (c *legacyConfigureCallback) Metadata(*interfaces.MetadataResponse) error { return nil }
func (c *legacyConfigureCallback) SpeechStarted(*interfaces.SpeechStartedResponse) error {
	return nil
}
func (c *legacyConfigureCallback) UtteranceEnd(*interfaces.UtteranceEndResponse) error { return nil }
func (c *legacyConfigureCallback) Close(*interfaces.CloseResponse) error               { return nil }
func (c *legacyConfigureCallback) Error(*interfaces.ErrorResponse) error               { return nil }
func (c *legacyConfigureCallback) UnhandledEvent(data []byte) error {
	c.unhandled = append(c.unhandled, data)
	return nil
}

type configureChannels struct {
	success chan *interfaces.ConfigureSuccessResponse
	failure chan *interfaces.ConfigureFailureResponse
}

func newConfigureChannels() *configureChannels {
	return &configureChannels{
		success: make(chan *interfaces.ConfigureSuccessResponse, 1),
		failure: make(chan *interfaces.ConfigureFailureResponse, 1),
	}
}

func (c *configureChannels) GetOpen() []*chan *interfaces.OpenResponse         { return nil }
func (c *configureChannels) GetMessage() []*chan *interfaces.MessageResponse   { return nil }
func (c *configureChannels) GetMetadata() []*chan *interfaces.MetadataResponse { return nil }
func (c *configureChannels) GetSpeechStarted() []*chan *interfaces.SpeechStartedResponse {
	return nil
}
func (c *configureChannels) GetUtteranceEnd() []*chan *interfaces.UtteranceEndResponse {
	return nil
}
func (c *configureChannels) GetConfigureSuccess() []*chan *interfaces.ConfigureSuccessResponse {
	return []*chan *interfaces.ConfigureSuccessResponse{&c.success}
}
func (c *configureChannels) GetConfigureFailure() []*chan *interfaces.ConfigureFailureResponse {
	return []*chan *interfaces.ConfigureFailureResponse{&c.failure}
}
func (c *configureChannels) GetClose() []*chan *interfaces.CloseResponse { return nil }
func (c *configureChannels) GetError() []*chan *interfaces.ErrorResponse { return nil }
func (c *configureChannels) GetUnhandled() []*chan *[]byte               { return nil }

type legacyConfigureChannels struct {
	unhandled chan *[]byte
}

func newLegacyConfigureChannels() *legacyConfigureChannels {
	return &legacyConfigureChannels{unhandled: make(chan *[]byte, 1)}
}

func (c *legacyConfigureChannels) GetOpen() []*chan *interfaces.OpenResponse         { return nil }
func (c *legacyConfigureChannels) GetMessage() []*chan *interfaces.MessageResponse   { return nil }
func (c *legacyConfigureChannels) GetMetadata() []*chan *interfaces.MetadataResponse { return nil }
func (c *legacyConfigureChannels) GetSpeechStarted() []*chan *interfaces.SpeechStartedResponse {
	return nil
}
func (c *legacyConfigureChannels) GetUtteranceEnd() []*chan *interfaces.UtteranceEndResponse {
	return nil
}
func (c *legacyConfigureChannels) GetClose() []*chan *interfaces.CloseResponse { return nil }
func (c *legacyConfigureChannels) GetError() []*chan *interfaces.ErrorResponse { return nil }
func (c *legacyConfigureChannels) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&c.unhandled}
}

func TestConfigureMessageSerialization(t *testing.T) {
	t.Run("omitted keyterms", func(t *testing.T) {
		data, err := json.Marshal(interfaces.ConfigureMessage{Type: "Configure"})
		if err != nil {
			t.Fatalf("marshal ConfigureMessage: %v", err)
		}
		if got, want := string(data), `{"type":"Configure"}`; got != want {
			t.Errorf("ConfigureMessage = %s, want %s", got, want)
		}
	})

	t.Run("explicit empty keyterms clears the value", func(t *testing.T) {
		empty := []string{}
		data, err := json.Marshal(interfaces.ConfigureMessage{Type: "Configure", Keyterms: &empty})
		if err != nil {
			t.Fatalf("marshal ConfigureMessage: %v", err)
		}
		if got, want := string(data), `{"type":"Configure","keyterms":[]}`; got != want {
			t.Errorf("ConfigureMessage = %s, want %s", got, want)
		}
	})

	t.Run("nullable keyterms are retained", func(t *testing.T) {
		var keyterms []string
		data, err := json.Marshal(interfaces.ConfigureMessage{Type: "Configure", Keyterms: &keyterms})
		if err != nil {
			t.Fatalf("marshal ConfigureMessage: %v", err)
		}
		if got, want := string(data), `{"type":"Configure","keyterms":null}`; got != want {
			t.Errorf("ConfigureMessage = %s, want %s", got, want)
		}
	})

	t.Run("keyterms and false feature values are retained", func(t *testing.T) {
		keyterms := []string{"Deepgram"}
		data, err := json.Marshal(interfaces.ConfigureMessage{
			Type:     "Configure",
			Keyterms: &keyterms,
			Features: map[string]bool{"punctuate": false},
		})
		if err != nil {
			t.Fatalf("marshal ConfigureMessage: %v", err)
		}
		if got, want := string(data), `{"type":"Configure","keyterms":["Deepgram"],"features":{"punctuate":false}}`; got != want {
			t.Errorf("ConfigureMessage = %s, want %s", got, want)
		}
	})
}

func TestConfigureResponseDecoding(t *testing.T) {
	var success interfaces.ConfigureSuccessResponse
	if err := json.Unmarshal([]byte(`{"type":"ConfigureSuccess","request_id":"request-1","keyterms":["Deepgram"]}`), &success); err != nil {
		t.Fatalf("decode ConfigureSuccess: %v", err)
	}
	if success.Type != "ConfigureSuccess" || success.RequestID != "request-1" || len(success.Keyterms) != 1 || success.Keyterms[0] != "Deepgram" {
		t.Errorf("unexpected ConfigureSuccess: %#v", success)
	}
	success = interfaces.ConfigureSuccessResponse{}
	if err := json.Unmarshal([]byte(`{"type":"ConfigureSuccess","request_id":"request-2"}`), &success); err != nil {
		t.Fatalf("decode ConfigureSuccess without keyterms: %v", err)
	}
	if success.Keyterms != nil {
		t.Errorf("ConfigureSuccess keyterms = %#v, want nil when omitted", success.Keyterms)
	}

	var failure interfaces.ConfigureFailureResponse
	if err := json.Unmarshal([]byte(`{"type":"ConfigureFailure","request_id":"request-1","code":"INVALID_KEYTERM","description":"invalid keyterm"}`), &failure); err != nil {
		t.Fatalf("decode ConfigureFailure: %v", err)
	}
	if failure.Type != "ConfigureFailure" || failure.RequestID != "request-1" || failure.Code != "INVALID_KEYTERM" || failure.Description != "invalid keyterm" {
		t.Errorf("unexpected ConfigureFailure: %#v", failure)
	}
}

func TestConfigureCallbackRouter(t *testing.T) {
	handler := &configureCallback{}
	router := NewCallbackRouter(handler)

	if err := router.Message([]byte(`{"type":"ConfigureSuccess","request_id":"request-1","keyterms":["Deepgram"]}`)); err != nil {
		t.Fatalf("route ConfigureSuccess: %v", err)
	}
	if err := router.Message([]byte(`{"type":"ConfigureFailure","request_id":"request-1","code":"INVALID_KEYTERM","description":"invalid keyterm"}`)); err != nil {
		t.Fatalf("route ConfigureFailure: %v", err)
	}
	if len(handler.successes) != 1 || handler.successes[0].RequestID != "request-1" {
		t.Errorf("unexpected ConfigureSuccess routing: %#v", handler.successes)
	}
	if len(handler.failures) != 1 || handler.failures[0].Code != "INVALID_KEYTERM" {
		t.Errorf("unexpected ConfigureFailure routing: %#v", handler.failures)
	}

	legacy := &legacyConfigureCallback{}
	legacyRouter := NewCallbackRouter(legacy)
	if err := legacyRouter.Message([]byte(`{"type":"ConfigureSuccess","request_id":"request-2"}`)); err != nil {
		t.Fatalf("route ConfigureSuccess for legacy handler: %v", err)
	}
	if len(legacy.unhandled) != 1 {
		t.Errorf("legacy handler received %d unhandled events, want 1", len(legacy.unhandled))
	}
}

func TestConfigureChanRouter(t *testing.T) {
	handler := newConfigureChannels()
	router := NewChanRouter(handler)

	if err := router.Message([]byte(`{"type":"ConfigureSuccess","request_id":"request-1","keyterms":["Deepgram"]}`)); err != nil {
		t.Fatalf("route ConfigureSuccess: %v", err)
	}
	if err := router.Message([]byte(`{"type":"ConfigureFailure","request_id":"request-1","code":"INVALID_KEYTERM","description":"invalid keyterm"}`)); err != nil {
		t.Fatalf("route ConfigureFailure: %v", err)
	}

	select {
	case success := <-handler.success:
		if success.RequestID != "request-1" || len(success.Keyterms) != 1 {
			t.Errorf("unexpected ConfigureSuccess: %#v", success)
		}
	case <-time.After(time.Second):
		t.Fatal("ConfigureSuccess was not delivered")
	}
	select {
	case failure := <-handler.failure:
		if failure.Code != "INVALID_KEYTERM" || failure.Description != "invalid keyterm" {
			t.Errorf("unexpected ConfigureFailure: %#v", failure)
		}
	case <-time.After(time.Second):
		t.Fatal("ConfigureFailure was not delivered")
	}

	legacy := newLegacyConfigureChannels()
	legacyRouter := NewChanRouter(legacy)
	if err := legacyRouter.Message([]byte(`{"type":"ConfigureSuccess","request_id":"request-2"}`)); err != nil {
		t.Fatalf("route ConfigureSuccess for legacy handler: %v", err)
	}
	select {
	case raw := <-legacy.unhandled:
		if string(*raw) != `{"type":"ConfigureSuccess","request_id":"request-2"}` {
			t.Errorf("legacy handler received %s", *raw)
		}
	case <-time.After(time.Second):
		t.Fatal("legacy channel handler did not receive ConfigureSuccess")
	}
}
