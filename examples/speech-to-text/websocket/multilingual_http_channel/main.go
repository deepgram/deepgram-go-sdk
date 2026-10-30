// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	listenws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

const streamURL = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"

type languageHandler struct {
	listenws.DefaultCallbackHandler
}

func (h languageHandler) Message(response *msginterfaces.MessageResponse) error {
	for _, alternative := range response.Channel.Alternatives {
		fmt.Printf("Transcript: %s\nDetected languages: %v\nApplied language hints: %v\n", alternative.Transcript, alternative.Languages, alternative.LanguagesHinted)
	}
	return nil
}

func main() {
	ctx := context.Background()
	options := &interfaces.LiveTranscriptionOptions{
		Model:        "nova-3",
		Language:     "multi",
		LanguageHint: []string{"en", "es"},
		Punctuate:    true,
	}
	client, err := listen.NewWSUsingCallback(ctx, "", &interfaces.ClientOptions{}, options, languageHandler{DefaultCallbackHandler: *listenws.NewDefaultCallbackHandler()})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Stop()

	response, err := http.Get(streamURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if !client.Connect() {
		log.Fatal("connect Listen WebSocket")
	}
	go client.Stream(bufio.NewReader(response.Body))

	log.Print("Streaming multilingual audio. Press Enter to stop.")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
