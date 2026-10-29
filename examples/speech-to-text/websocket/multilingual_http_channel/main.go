// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

const streamURL = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"

func main() {
	ctx := context.Background()
	client, err := listen.NewWSUsingChanForDemo(ctx, &interfaces.LiveTranscriptionOptions{
		Model:        "nova-3",
		Language:     "multi",
		LanguageHint: []string{"en", "es"},
		Punctuate:    true,
	})
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
