// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	listenws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

const streamURL = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"

type entityHandler struct {
	errorChan chan<- error
}

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

func (h *entityHandler) Error(response *listenws.ErrorResponse) error {
	message := response.ErrMsg
	if message == "" {
		message = response.Description
	}
	if message == "" {
		message = "received a Deepgram streaming error"
	}

	select {
	case h.errorChan <- errors.New(message):
	default:
	}
	return nil
}

func (h *entityHandler) UnhandledEvent([]byte) error { return nil }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client.InitWithDefault()
	ctx := context.Background()

	options := &interfaces.LiveTranscriptionOptions{
		Model:          "nova-3",
		DetectEntities: true,
		Language:       "en-US",
		Punctuate:      true,
	}
	errorChan := make(chan error, 1)
	dgClient, err := client.NewWSUsingCallbackWithDefaults(ctx, options, &entityHandler{errorChan: errorChan})
	if err != nil {
		return fmt.Errorf("create LiveTranscription connection: %w", err)
	}
	defer dgClient.Stop()

	response, err := http.Get(streamURL)
	if err != nil {
		return fmt.Errorf("get audio stream: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("get audio stream: %s", response.Status)
	}
	fmt.Printf("Stream is up and running: %s\n", response.Status)

	if !dgClient.Connect() {
		return errors.New("connect LiveTranscription client")
	}
	streamErr := make(chan error, 1)
	go func() {
		streamErr <- dgClient.Stream(bufio.NewReader(response.Body))
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	inputDone := make(chan struct{})
	go func() {
		bufio.NewScanner(os.Stdin).Scan()
		close(inputDone)
	}()

	select {
	case err := <-errorChan:
		return err
	case err := <-streamErr:
		if err != nil {
			return fmt.Errorf("stream audio: %w", err)
		}
		return errors.New("audio stream ended")
	case <-inputDone:
		return nil
	}
}
