// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

func main() {
	client := listen.NewRESTWithDefaults()
	response, err := api.New(client).FromURL(context.Background(), "https://dpgr.am/spacewalk.wav", &interfaces.PreRecordedTranscriptionOptions{
		Model:        "nova-3",
		Language:     "multi",
		LanguageHint: []string{"en", "es"},
		Utterances:   true,
	})
	if err != nil {
		log.Fatal(err)
	}

	alternative := response.Results.Channels[0].Alternatives[0]
	fmt.Printf("Transcript: %s\nApplied language hints: %v\n", alternative.Transcript, alternative.LanguagesHinted)
	for _, utterance := range response.Results.Utterances {
		fmt.Printf("Utterance languages: %v\n", utterance.Languages)
	}
}
