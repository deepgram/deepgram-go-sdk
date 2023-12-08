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

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/live"
)

const (
	STREAM_URL = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"
)

func main() {
	// init library
	client.InitWithDefault()

	// Go context
	ctx := context.Background()

	// set the Transcription options
	transcriptOptions := interfaces.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}

	// create a Deepgram client
	dgClient, err := client.NewForDemo(ctx, transcriptOptions)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// connect the websocket to Deepgram
	wsconn := dgClient.Connect()
	if wsconn == nil {
		log.Println("Client.Connect failed")
		os.Exit(1)
	}

	// get the HTTP stream
	httpClient := new(http.Client)

	res, err := httpClient.Get(STREAM_URL)
	if err != nil {
		log.Printf("httpClient.Get failed. Err: %v\n", err)
		return
	}

	log.Printf("Stream is up and running %s\n", reflect.TypeOf(res))

	go func() {
		// feed the HTTP stream to the Deepgram client (this is a blocking call)
		dgClient.Stream(bufio.NewReader(res.Body))
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close HTTP stream
	res.Body.Close()

	// close client
	dgClient.Stop()

	log.Printf("Program exiting...\n")
}
