// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
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
	"reflect"
	"time"

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/live"
)

const (
	STREAM_URL             = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"
	CHUNK_SIZE             = 1024 * 2
	TEN_MILLISECONDS_SLEEP = 10 * time.Millisecond
)

func main() {
	// context
	ctx := context.Background()

	// options
	transcriptOptions := interfaces.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}

	dgClient, err := client.NewWithDefaults(ctx, "", transcriptOptions)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// call connect!
	wsconn := dgClient.Connect()
	if wsconn == nil {
		log.Println("Client.Connect failed")
		os.Exit(1)
	}

	// feed the stream to the websocket
	httpClient := new(http.Client)

	res, err := httpClient.Get(STREAM_URL)
	if err != nil {
		log.Printf("httpClient.Get failed. Err: %v\n", err)
		return
	}

	fmt.Printf("Stream is up and running %s\n", reflect.TypeOf(res))

	// this is a blocking call...
	go func() {
		dgClient.Stream(bufio.NewReader(res.Body))
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close HTTP stream
	res.Body.Close()

	// close client
	dgClient.Stop()

	fmt.Printf("Succeeded!\n\n")
}
