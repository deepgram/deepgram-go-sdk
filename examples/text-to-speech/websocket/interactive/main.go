// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
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

	file, err := os.OpenFile(AUDIO_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
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

func (c MyCallback) Error(er *msginterfaces.ErrorResponse) error {
	fmt.Printf("\n[Error] Received\n")
	fmt.Printf("Error.Type: %s\n", er.Type)
	fmt.Printf("Error.Description: %s\n\n", er.Description)
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
	speak.InitWithDefault()

	// Go context
	ctx := context.Background()

	// print instructions
	fmt.Print("\n\nPress ENTER to exit!\n\n")

	// set the TTS options
	ttsOptions := &interfaces.SpeakOptions{
		Model: "aura-asteria-en",
	}

	// set the Client options
	cOptions := &interfaces.ClientOptions{}

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
	err = dgClient.WriteJSON(map[string]interface{}{
		"type": "Speak",
		"text": TTS_TEXT,
	})
	if err != nil {
		fmt.Printf("Error sending text input: %v\n", err)
		return
	}

	// Simulate user input to reset the buffer, flush, send new text, or just exit
	fmt.Print("\n\nPress 'r' and ENTER to reset the buffer, 'f' and ENTER to flush, enter new text to send it, or just ENTER to exit...\n\n")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch input.Text() {
		case "r":
			err = dgClient.Reset()
			if err != nil {
				fmt.Printf("Error resetting buffer: %v\n", err)
			} else {
				fmt.Println("Buffer reset successfully.")
			}
		case "f":
			err = dgClient.Flush()
			if err != nil {
				fmt.Printf("Error flushing buffer: %v\n", err)
			} else {
				fmt.Println("Buffer flushed successfully.")
			}
		case "":
			goto EXIT
		default:
			err = dgClient.WriteJSON(map[string]interface{}{
				"type": "Speak",
				"text": input.Text(),
			})
			if err != nil {
				fmt.Printf("Error sending text input: %v\n", err)
			} else {
				fmt.Println("Text sent successfully.")
			}
		}
	}

EXIT:

	// close the connection
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
}
