// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package daily_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	rest "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

func TestDaily_ListenLanguageHint(t *testing.T) {
	if os.Getenv("DEEPGRAM_API_KEY") == "" && os.Getenv("DEEPGRAM_ACCESS_TOKEN") == "" {
		t.Skip("DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	hints := []string{"en", "es"}

	t.Run("REST returns applied hints", func(t *testing.T) {
		response, err := rest.New(listen.NewRESTWithDefaults()).FromURL(ctx, "https://dpgr.am/spacewalk.wav", &interfaces.PreRecordedTranscriptionOptions{
			Model:        "nova-3",
			Language:     "multi",
			LanguageHint: hints,
		})
		if err != nil {
			t.Fatalf("transcribe multilingual audio: %v", err)
		}
		if got := response.Results.Channels[0].Alternatives[0].LanguagesHinted; !reflect.DeepEqual(got, hints) {
			t.Fatalf("LanguagesHinted = %q, want %q", got, hints)
		}
	})

	t.Run("streaming accepts hints", func(t *testing.T) {
		client, err := listen.NewWSUsingChanForDemo(ctx, &interfaces.LiveTranscriptionOptions{
			Model:        "nova-3",
			Language:     "multi",
			LanguageHint: hints,
		})
		if err != nil {
			t.Fatalf("create multilingual streaming client: %v", err)
		}
		defer client.Stop()

		if !client.Connect() {
			t.Fatal("connect multilingual streaming client")
		}
	})
}
