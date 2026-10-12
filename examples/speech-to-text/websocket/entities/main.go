// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"

	listenws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

const streamURL = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"

type entityHandler struct{}

func (h *entityHandler) Open(*listenws.OpenResponse) error { return nil }

func (h *entityHandler) Message(message *listenws.MessageResponse) error {
	if !message.IsFinal {
		return nil
	}

	for _, alternative := range message.Channel.Alternatives {
		fmt.Printf("Transcript: %s\n", alternative.Transcript)
	}
	for _, entity := range message.Entities {
		fmt.Printf(
			"Entity: label=%s value=%q raw_value=%q confidence=%.2f words=%d-%d\n",
			entity.Label,
			entity.Value,
			entity.RawValue,
			entity.Confidence,
			entity.StartWord,
			entity.EndWord,
		)
	}
	return nil
}

func (h *entityHandler) Metadata(*listenws.MetadataResponse) error { return nil }

func (h *entityHandler) SpeechStarted(*listenws.SpeechStartedResponse) error { return nil }

func (h *entityHandler) UtteranceEnd(*listenws.UtteranceEndResponse) error { return nil }

func (h *entityHandler) Close(*listenws.CloseResponse) error { return nil }

func (h *entityHandler) Error(*listenws.ErrorResponse) error { return nil }

func (h *entityHandler) UnhandledEvent([]byte) error { return nil }

func main() {
	client.InitWithDefault()
	ctx := context.Background()

	options := &interfaces.LiveTranscriptionOptions{
		Model:          "nova-3",
		DetectEntities: true,
		Language:       "en-US",
		Punctuate:      true,
	}
	dgClient, err := client.NewWSUsingCallbackWithDefaults(ctx, options, &entityHandler{})
	if err != nil {
		fmt.Printf("ERROR creating LiveTranscription connection: %v\n", err)
		return
	}
	defer dgClient.Stop()

	response, err := http.Get(streamURL)
	if err != nil {
		fmt.Printf("http.Get failed: %v\n", err)
		return
	}
	defer response.Body.Close()
	fmt.Printf("Stream is up and running %s\n", reflect.TypeOf(response))

	if !dgClient.Connect() {
		fmt.Println("Client.Connect failed")
		return
	}
	go dgClient.Stream(bufio.NewReader(response.Body))

	fmt.Print("Press ENTER to exit!\n\n")
	bufio.NewScanner(os.Stdin).Scan()
}
