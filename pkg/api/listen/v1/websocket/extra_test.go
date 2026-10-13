// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"reflect"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
)

const utteranceEndWithoutExtra = `{"type":"UtteranceEnd","channel":[0],"last_word_end":1}`

type extraCallback struct {
	utteranceEnds []*interfaces.UtteranceEndResponse
}

func (c *extraCallback) Open(or *interfaces.OpenResponse) error { return nil }
func (c *extraCallback) Message(mr *interfaces.MessageResponse) error {
	return nil
}
func (c *extraCallback) Metadata(md *interfaces.MetadataResponse) error {
	return nil
}
func (c *extraCallback) SpeechStarted(ssr *interfaces.SpeechStartedResponse) error {
	return nil
}
func (c *extraCallback) UtteranceEnd(ur *interfaces.UtteranceEndResponse) error {
	c.utteranceEnds = append(c.utteranceEnds, ur)
	return nil
}
func (c *extraCallback) Close(cr *interfaces.CloseResponse) error { return nil }
func (c *extraCallback) Error(er *interfaces.ErrorResponse) error { return nil }
func (c *extraCallback) UnhandledEvent(byData []byte) error       { return nil }

type extraChannels struct {
	utteranceEnd chan *interfaces.UtteranceEndResponse
}

func newExtraChannels() *extraChannels {
	return &extraChannels{utteranceEnd: make(chan *interfaces.UtteranceEndResponse, 2)}
}

func (c *extraChannels) GetOpen() []*chan *interfaces.OpenResponse         { return nil }
func (c *extraChannels) GetMessage() []*chan *interfaces.MessageResponse   { return nil }
func (c *extraChannels) GetMetadata() []*chan *interfaces.MetadataResponse { return nil }
func (c *extraChannels) GetSpeechStarted() []*chan *interfaces.SpeechStartedResponse {
	return nil
}
func (c *extraChannels) GetUtteranceEnd() []*chan *interfaces.UtteranceEndResponse {
	return []*chan *interfaces.UtteranceEndResponse{&c.utteranceEnd}
}
func (c *extraChannels) GetClose() []*chan *interfaces.CloseResponse { return nil }
func (c *extraChannels) GetError() []*chan *interfaces.ErrorResponse { return nil }
func (c *extraChannels) GetUnhandled() []*chan *[]byte               { return nil }

func Test_UtteranceEndExtraMetadataCallback(t *testing.T) {
	tests := []struct {
		name            string
		source          string
		utteranceEnd    string
		want            map[string]string
		verifyIsolation bool
	}{
		{
			name:            "metadata response",
			source:          `{"type":"Metadata","extra":{"session":"cached"}}`,
			utteranceEnd:    utteranceEndWithoutExtra,
			want:            map[string]string{"session": "cached"},
			verifyIsolation: true,
		},
		{
			name:         "results response",
			source:       `{"type":"Results","metadata":{"extra":{"session":"cached"}}}`,
			utteranceEnd: utteranceEndWithoutExtra,
			want:         map[string]string{"session": "cached"},
		},
		{
			name:         "wire extra takes precedence",
			source:       `{"type":"Metadata","extra":{"session":"cached"}}`,
			utteranceEnd: `{"type":"UtteranceEnd","extra":{"session":"wire"}}`,
			want:         map[string]string{"session": "wire"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &extraCallback{}
			router := NewCallbackRouter(handler)

			if err := router.Message([]byte(tt.source)); err != nil {
				t.Fatalf("route source message: %v", err)
			}
			if err := router.Message([]byte(tt.utteranceEnd)); err != nil {
				t.Fatalf("route utterance end: %v", err)
			}

			if len(handler.utteranceEnds) != 1 {
				t.Fatalf("expected one utterance end, got %d", len(handler.utteranceEnds))
			}
			if got := handler.utteranceEnds[0].Extra; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extra = %#v, want %#v", got, tt.want)
			}

			if tt.verifyIsolation {
				handler.utteranceEnds[0].Extra["session"] = "mutated"
				if err := router.Message([]byte(utteranceEndWithoutExtra)); err != nil {
					t.Fatalf("route second utterance end: %v", err)
				}
				if got := handler.utteranceEnds[1].Extra; !reflect.DeepEqual(got, tt.want) {
					t.Errorf("cached extra after caller mutation = %#v, want %#v", got, tt.want)
				}
			}
		})
	}
}

func Test_UtteranceEndExtraMetadataChannel(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		utteranceEnd string
		want         map[string]string
	}{
		{
			name:         "metadata response",
			source:       `{"type":"Metadata","extra":{"session":"cached"}}`,
			utteranceEnd: utteranceEndWithoutExtra,
			want:         map[string]string{"session": "cached"},
		},
		{
			name:         "results response",
			source:       `{"type":"Results","metadata":{"extra":{"session":"cached"}}}`,
			utteranceEnd: utteranceEndWithoutExtra,
			want:         map[string]string{"session": "cached"},
		},
		{
			name:         "wire extra takes precedence",
			source:       `{"type":"Metadata","extra":{"session":"cached"}}`,
			utteranceEnd: `{"type":"UtteranceEnd","extra":{"session":"wire"}}`,
			want:         map[string]string{"session": "wire"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newExtraChannels()
			router := NewChanRouter(handler)

			if err := router.Message([]byte(tt.source)); err != nil {
				t.Fatalf("route source message: %v", err)
			}
			if err := router.Message([]byte(tt.utteranceEnd)); err != nil {
				t.Fatalf("route utterance end: %v", err)
			}

			got := <-handler.utteranceEnd
			if !reflect.DeepEqual(got.Extra, tt.want) {
				t.Errorf("extra = %#v, want %#v", got.Extra, tt.want)
			}
		})
	}
}
