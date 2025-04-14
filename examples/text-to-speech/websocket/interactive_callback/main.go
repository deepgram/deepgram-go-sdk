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
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
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

	file, err := os.OpenFile(AUDIO_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", AUDIO_FILE, err)
		return err
	}

	_, err = file.Write(byMsg)
	file.Close()

	if err != nil {
		fmt.Printf("Error writing audio data to file: %v\n", err)
		return err
	}

	return nil
}

func (c MyCallback) Flush(fl *msginterfaces.FlushedResponse) error {
	fmt.Printf("\n[Flushed] Received\n")
	fmt.Printf("\n\nPress 'r' and ENTER to reset the buffer, 'f' and ENTER to flush, enter new text to send it, or just ENTER to exit...\n\n> ")
	return nil
}

func (c MyCallback) Clear(fl *msginterfaces.ClearedResponse) error {
	fmt.Printf("\n[Cleared] Received\n")
	fmt.Printf("\n\nPress 'r' and ENTER to reset the buffer, 'f' and ENTER to flush, enter new text to send it, or just ENTER to exit...\n\n> ")
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
	speak.InitWithDefault()

	// Go context
	ctx := context.Background()

	// print instructions
	fmt.Print("\n\nPress ENTER to exit!\n\n")

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

	// Simulate user input to reset the buffer, flush, send new text, or just exit
	time.Sleep(2 * time.Second)
	fmt.Printf("\n\nPress 'r' and ENTER to reset the buffer, 'f' and ENTER to flush, enter new text to send it, or just ENTER to exit...\n\n> ")
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
			// delete file if exists
			_ = os.Remove(AUDIO_FILE)

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

			err = dgClient.Flush()
			if err != nil {
				fmt.Printf("Error flushing buffer: %v\n", err)
			} else {
				fmt.Println("Buffer flushed successfully.")
			}
		case "":
			goto EXIT
		default:
			err = dgClient.SpeakWithText(input.Text())
			if err != nil {
				fmt.Printf("Error sending text input: %v\n", err)
			} else {
				fmt.Println("Text sent successfully.")
			}
			fmt.Printf("\n\nPress 'r' and ENTER to reset the buffer, 'f' and ENTER to flush, enter new text to send it, or just ENTER to exit...\n\n> ")
		}
	}

EXIT:

	// close the connection
	dgClient.Stop()

	fmt.Printf("Program exiting...\n")
}
