// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type configureTestChanHandler struct{}

func (configureTestChanHandler) GetOpen() []*chan *msginterfaces.OpenResponse       { return nil }
func (configureTestChanHandler) GetMessage() []*chan *msginterfaces.MessageResponse { return nil }
func (configureTestChanHandler) GetMetadata() []*chan *msginterfaces.MetadataResponse {
	return nil
}
func (configureTestChanHandler) GetSpeechStarted() []*chan *msginterfaces.SpeechStartedResponse {
	return nil
}
func (configureTestChanHandler) GetUtteranceEnd() []*chan *msginterfaces.UtteranceEndResponse {
	return nil
}
func (configureTestChanHandler) GetClose() []*chan *msginterfaces.CloseResponse { return nil }
func (configureTestChanHandler) GetError() []*chan *msginterfaces.ErrorResponse { return nil }
func (configureTestChanHandler) GetUnhandled() []*chan *[]byte                  { return nil }

func newConfigureTestServer(t *testing.T) (string, <-chan string) {
	t.Helper()

	frames := make(chan string, 2)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.TextMessage {
				frames <- string(message)
			}
		}
	}))
	t.Cleanup(server.Close)
	return strings.TrimPrefix(server.URL, "https://"), frames
}

func configureTestClientOptions(host string) *clientinterfaces.ClientOptions {
	return &clientinterfaces.ClientOptions{Host: host, SkipServerAuth: true}
}

func configureTestTranscriptionOptions() *clientinterfaces.LiveTranscriptionOptions {
	return &clientinterfaces.LiveTranscriptionOptions{Model: "nova-3"}
}

func TestConfigureRejectsNilAndEmptyOptions(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*clientinterfaces.LiveConfigureOptions) error
	}{
		{name: "callback", configure: (&WSCallback{}).Configure},
		{name: "channel", configure: (&WSChannel{}).Configure},
	}

	for _, test := range tests {
		t.Run(test.name+" nil", func(t *testing.T) {
			if err := test.configure(nil); !errors.Is(err, interfacesv1.ErrOptionsRequired) {
				t.Errorf("Configure(nil) error = %v, want ErrOptionsRequired", err)
			}
		})
		t.Run(test.name+" empty", func(t *testing.T) {
			if err := test.configure(&clientinterfaces.LiveConfigureOptions{}); !errors.Is(err, interfacesv1.ErrConfigureNoFields) {
				t.Errorf("Configure(empty) error = %v, want ErrConfigureNoFields", err)
			}
		})
	}
}

func TestConfigureWritesControlFrame(t *testing.T) {
	keyterms := []string{}
	opts := &clientinterfaces.LiveConfigureOptions{
		Keyterms: &keyterms,
		Features: map[string]bool{"punctuate": false},
	}
	const want = `{"type":"Configure","keyterms":[],"features":{"punctuate":false}}`

	t.Run("callback", func(t *testing.T) {
		host, frames := newConfigureTestServer(t)
		client, err := NewUsingCallback(context.Background(), "test-key", configureTestClientOptions(host), configureTestTranscriptionOptions(), nil)
		if err != nil {
			t.Fatalf("create callback client: %v", err)
		}
		if !client.Connect() {
			t.Fatal("connect callback client")
		}
		defer client.Stop()

		if err := client.Configure(opts); err != nil {
			t.Fatalf("Configure: %v", err)
		}
		assertConfigureFrame(t, frames, want)
	})

	t.Run("channel", func(t *testing.T) {
		host, frames := newConfigureTestServer(t)
		client, err := NewUsingChan(context.Background(), "test-key", configureTestClientOptions(host), configureTestTranscriptionOptions(), configureTestChanHandler{})
		if err != nil {
			t.Fatalf("create channel client: %v", err)
		}
		if !client.Connect() {
			t.Fatal("connect channel client")
		}
		defer client.Stop()

		if err := client.Configure(opts); err != nil {
			t.Fatalf("Configure: %v", err)
		}
		assertConfigureFrame(t, frames, want)
	})
}

func assertConfigureFrame(t *testing.T, frames <-chan string, want string) {
	t.Helper()
	select {
	case got := <-frames:
		if got != want {
			t.Errorf("Configure frame = %s, want %s", got, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not receive Configure frame")
	}
}
