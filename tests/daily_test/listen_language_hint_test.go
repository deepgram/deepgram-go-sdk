// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package daily_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	rest "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	replay "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/replay"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

type languageHintDailyHandler struct {
	messages chan *msginterfaces.MessageResponse
	errors   chan *msginterfaces.ErrorResponse
}

func (h *languageHintDailyHandler) GetOpen() []*chan *msginterfaces.OpenResponse { return nil }
func (h *languageHintDailyHandler) GetMessage() []*chan *msginterfaces.MessageResponse {
	return []*chan *msginterfaces.MessageResponse{&h.messages}
}
func (h *languageHintDailyHandler) GetMetadata() []*chan *msginterfaces.MetadataResponse {
	return nil
}
func (h *languageHintDailyHandler) GetSpeechStarted() []*chan *msginterfaces.SpeechStartedResponse {
	return nil
}
func (h *languageHintDailyHandler) GetUtteranceEnd() []*chan *msginterfaces.UtteranceEndResponse {
	return nil
}
func (h *languageHintDailyHandler) GetClose() []*chan *msginterfaces.CloseResponse { return nil }
func (h *languageHintDailyHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&h.errors}
}
func (h *languageHintDailyHandler) GetUnhandled() []*chan *[]byte { return nil }

func languageHintFixturePath(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate language hint test source")
	}
	return filepath.Join(filepath.Dir(filename), "..", "..", "examples", "speech-to-text", "websocket", "replay", "testing.wav")
}

func TestDaily_ListenLanguageHint(t *testing.T) {
	if os.Getenv("DEEPGRAM_API_KEY") == "" && os.Getenv("DEEPGRAM_ACCESS_TOKEN") == "" {
		t.Skip("DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN is required for the on-demand production contract test")
	}

	hints := []string{"en", "es"}

	t.Run("REST returns applied hints", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := rest.New(listen.NewRESTWithDefaults()).FromURL(ctx, "https://dpgr.am/spacewalk.wav", &interfaces.PreRecordedTranscriptionOptions{
			Model:        "nova-3",
			Language:     "multi",
			LanguageHint: hints,
			Utterances:   true,
		})
		if err != nil {
			t.Fatalf("transcribe multilingual audio: %v", err)
		}
		if len(response.Results.Channels) == 0 || len(response.Results.Channels[0].Alternatives) == 0 {
			t.Fatal("REST response contained no alternatives")
		}
		if got := response.Results.Channels[0].Alternatives[0].LanguagesHinted; !reflect.DeepEqual(got, hints) {
			t.Fatalf("LanguagesHinted = %q, want %q", got, hints)
		}
		if len(response.Results.Utterances) == 0 {
			t.Fatal("REST response contained no utterances")
		}
		if len(response.Results.Utterances[0].Languages) == 0 {
			t.Fatal("REST utterance did not include detected languages")
		}
	})

	t.Run("streaming returns applied hints", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		handler := &languageHintDailyHandler{
			messages: make(chan *msginterfaces.MessageResponse, 8),
			errors:   make(chan *msginterfaces.ErrorResponse, 1),
		}
		client, err := listen.NewWSUsingChan(ctx, "", &interfaces.ClientOptions{}, &interfaces.LiveTranscriptionOptions{
			Model:        "nova-3",
			Language:     "multi",
			LanguageHint: hints,
			Encoding:     "mulaw",
			Channels:     1,
			SampleRate:   8000,
		}, handler)
		if err != nil {
			t.Fatalf("create multilingual streaming client: %v", err)
		}
		defer client.Stop()

		if !client.Connect() {
			t.Fatal("connect multilingual streaming client")
		}

		audio, err := replay.New(replay.Options{FullFilename: languageHintFixturePath(t)})
		if err != nil {
			t.Fatalf("open audio fixture: %v", err)
		}
		defer audio.Stop()
		if err := audio.Start(); err != nil {
			t.Fatalf("start audio fixture: %v", err)
		}
		if err := audio.Stream(client); err != nil {
			t.Fatalf("stream audio fixture: %v", err)
		}
		if err := client.Finalize(); err != nil {
			t.Fatalf("finalize streaming transcription: %v", err)
		}

		for {
			select {
			case response := <-handler.messages:
				if response == nil || !response.IsFinal {
					continue
				}
				if len(response.Channel.Alternatives) == 0 {
					t.Fatal("final streaming response contained no alternatives")
				}
				for _, alternative := range response.Channel.Alternatives {
					if got := alternative.LanguagesHinted; !reflect.DeepEqual(got, hints) {
						t.Fatalf("LanguagesHinted = %q, want %q", got, hints)
					}
				}
				return
			case response := <-handler.errors:
				if response == nil {
					t.Fatal("streaming language hint error channel closed")
				}
				t.Fatalf("streaming language hint error: %s", response.Description)
			case <-ctx.Done():
				t.Fatalf("wait for final streaming result: %v", ctx.Err())
			}
		}
	})
}
