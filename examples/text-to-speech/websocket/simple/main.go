// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	speak "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

const (
	API_KEY    = ""
	TTS_TEXT   = "Hello, this is a text to speech example using Deepgram."
	AUDIO_FILE = "output.mp3"
)

// Implement your own callback
type MyCallback struct{}

func (c MyCallback) Metadata(md *msginterfaces.MetadataResponse) error {
	fmt.Printf("\n[Metadata] Received\n")
	fmt.Printf("Metadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	return nil
}

func (c MyCallback) Binary(byMsg []byte) error {
	fmt.Printf("\n[Binary] Received\n")

	file, err := os.OpenFile(AUDIO_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", AUDIO_FILE, err)
		return err
	}
	defer file.Close()

	_, err = file.Write(byMsg)
	if err != nil {
		fmt.Printf("Error writing audio data to file: %v\n", err)
		return err
	}

	fmt.Printf("Audio data saved to %s\n", AUDIO_FILE)
	return nil
}

func (c MyCallback) Flush(fl *msginterfaces.FlushedResponse) error {
	fmt.Printf("\n[Flushed] Received\n")
	return nil
}

func (c MyCallback) Warning(wr *msginterfaces.WarningResponse) error {
	fmt.Printf("\n[Warning] Received\n")
	fmt.Printf("Warning.Code: %s\n", wr.WarnCode)
	fmt.Printf("Warning.Description: %s\n\n", wr.WarnMsg)
	return nil
}

func (c MyCallback) Error(er *msginterfaces.ErrorResponse) error {
	fmt.Printf("\n[Error] Received\n")
	fmt.Printf("Error.Code: %s\n", er.ErrCode)
	fmt.Printf("Error.Description: %s\n\n", er.ErrMsg)
	return nil
}

func (c MyCallback) Close(cr *msginterfaces.CloseResponse) error {
	fmt.Printf("\n[Close] Received\n")
	return nil
}

func (c MyCallback) Open(or *msginterfaces.OpenResponse) error {
	fmt.Printf("\n[Open] Received\n")
	return nil
}

func main() {
	// init library
	speak.Init(speak.InitLib{
		LogLevel: speak.LogLevelDefault, // LogLevelDefault, LogLevelFull, LogLevelDebug, LogLevelTrace
	})

	// Go context
	ctx := context.Background()

	// set the Client options
	cOptions := &interfaces.ClientOptions{}

	// set the TTS options
	ttsOptions := &interfaces.SpeakOptions{
		Model: "aura-asteria-en",
	}

	// create the callback
	callback := MyCallback{}

	// create a new stream using the NewStream function
	dgClient, err := speak.NewWebSocket(ctx, "", cOptions, ttsOptions, callback)
	if err != nil {
		fmt.Println("ERROR creating TTS connection:", err)
		return
	}

	// connect the websocket to Deepgram
	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Println("Client.Connect failed")
		os.Exit(1)
	}

	// Send the text input
	err = dgClient.SpeakWithText(TTS_TEXT)
	if err != nil {
		fmt.Printf("Error sending text input: %v\n", err)
		return
	}

	// Flush the text input
	err = dgClient.Flush()
	if err != nil {
		fmt.Printf("Error sending text input: %v\n", err)
		return
	}

	// wait for user input to exit
	time.Sleep(5 * time.Second)

	// close the connection
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
}
