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
	TTS_TEXT   = "Hello, this is a text to speech example using Deepgram."
	AUDIO_FILE = "output.wav"
)

// Implement your own callback
type MyCallback struct{}

func (c MyCallback) Open(or *msginterfaces.OpenResponse) error {
	fmt.Printf("\n[Open] Received\n")
	return nil
}

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

func (c MyCallback) Clear(fl *msginterfaces.ClearedResponse) error {
	fmt.Printf("\n[Cleared] Received\n")
	return nil
}

func (c MyCallback) Close(cr *msginterfaces.CloseResponse) error {
	fmt.Printf("\n[Close] Received\n")
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

func (c MyCallback) UnhandledEvent(byData []byte) error {
	// handle the unhandled event
	fmt.Printf("\n[UnhandledEvent] Received\n")
	fmt.Printf("UnhandledEvent: %s\n\n", string(byData))
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
	cOptions := &interfaces.ClientOptions{
		// AutoFlushSpeakDelta: 1000,
	}

	// set the TTS options
	ttsOptions := &interfaces.WSSpeakOptions{
		Model:      "aura-asteria-en",
		Encoding:   "linear16",
		SampleRate: 48000,
	}

	// create the callback
	callback := MyCallback{}

	// create a new stream using the NewStream function
	dgClient, err := speak.NewWSUsingCallback(ctx, "", cOptions, ttsOptions, callback)
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

	file, err := os.OpenFile(AUDIO_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		fmt.Printf("Failed to open file. Err: %v\n", err)
		return
	}

	// Add a wav audio container header to the file if you want to play the audio
	// using a media player like VLC, Media Player, or Apple Music
	header := []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x00, 0x00, 0x00, 0x00, // Placeholder for file size
		0x57, 0x41, 0x56, 0x45, // "WAVE"
		0x66, 0x6d, 0x74, 0x20, // "fmt "
		0x10, 0x00, 0x00, 0x00, // Chunk size (16)
		0x01, 0x00, // Audio format (1 for PCM)
		0x01, 0x00, // Number of channels (1)
		0x80, 0xbb, 0x00, 0x00, // Sample rate (48000)
		0x00, 0xee, 0x02, 0x00, // Byte rate (48000 * 2)
		0x02, 0x00, // Block align (2)
		0x10, 0x00, // Bits per sample (16)
		0x64, 0x61, 0x74, 0x61, // "data"
		0x00, 0x00, 0x00, 0x00, // Placeholder for data size
	}

	_, err = file.Write(header)
	if err != nil {
		fmt.Printf("Failed to write header to file. Err: %v\n", err)
		return
	}
	file.Close()

	// Send the text input
	err = dgClient.SpeakWithText(TTS_TEXT)
	if err != nil {
		fmt.Printf("Error sending text input: %v\n", err)
		return
	}

	// If AutoFlushSpeakDelta is not set, you Flush the text input manually
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
