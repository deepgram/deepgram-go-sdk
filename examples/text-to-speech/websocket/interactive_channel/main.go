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
	"sync"
	"time"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	speak "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

const (
	TTS_TEXT   = "Hello, this is a text to speech example using Deepgram."
	AUDIO_FILE = "output.wav"
)

type MyHandler struct {
	binaryChan    chan *[]byte
	openChan      chan *msginterfaces.OpenResponse
	metadataChan  chan *msginterfaces.MetadataResponse
	flushChan     chan *msginterfaces.FlushedResponse
	clearChan     chan *msginterfaces.ClearedResponse
	closeChan     chan *msginterfaces.CloseResponse
	warningChan   chan *msginterfaces.WarningResponse
	errorChan     chan *msginterfaces.ErrorResponse
	unhandledChan chan *[]byte
}

func NewMyHandler() MyHandler {
	handler := MyHandler{
		binaryChan:    make(chan *[]byte),
		openChan:      make(chan *msginterfaces.OpenResponse),
		metadataChan:  make(chan *msginterfaces.MetadataResponse),
		flushChan:     make(chan *msginterfaces.FlushedResponse),
		clearChan:     make(chan *msginterfaces.ClearedResponse),
		closeChan:     make(chan *msginterfaces.CloseResponse),
		warningChan:   make(chan *msginterfaces.WarningResponse),
		errorChan:     make(chan *msginterfaces.ErrorResponse),
		unhandledChan: make(chan *[]byte),
	}

	go func() {
		handler.Run()
	}()

	return handler
}

// GetUnhandled returns the binary event channels
func (dch MyHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&dch.binaryChan}
}

// GetOpen returns the open channels
func (dch MyHandler) GetOpen() []*chan *msginterfaces.OpenResponse {
	return []*chan *msginterfaces.OpenResponse{&dch.openChan}
}

// GetMetadata returns the metadata channels
func (dch MyHandler) GetMetadata() []*chan *msginterfaces.MetadataResponse {
	return []*chan *msginterfaces.MetadataResponse{&dch.metadataChan}
}

// GetFlushed returns the flush channels
func (dch MyHandler) GetFlush() []*chan *msginterfaces.FlushedResponse {
	return []*chan *msginterfaces.FlushedResponse{&dch.flushChan}
}

// Getcleared returns the flush channels
func (dch MyHandler) GetClear() []*chan *msginterfaces.ClearedResponse {
	return []*chan *msginterfaces.ClearedResponse{&dch.clearChan}
}

// GetClose returns the close channels
func (dch MyHandler) GetClose() []*chan *msginterfaces.CloseResponse {
	return []*chan *msginterfaces.CloseResponse{&dch.closeChan}
}

// GetWarning returns the warning channels
func (dch MyHandler) GetWarning() []*chan *msginterfaces.WarningResponse {
	return []*chan *msginterfaces.WarningResponse{&dch.warningChan}
}

// GetError returns the error channels
func (dch MyHandler) GetError() []*chan *msginterfaces.ErrorResponse {
	return []*chan *msginterfaces.ErrorResponse{&dch.errorChan}
}

// GetUnhandled returns the unhandled event channels
func (dch MyHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&dch.unhandledChan}
}

// Open is the callback for when the connection opens
// golintci: funlen
func (dch MyHandler) Run() error {
	wgReceivers := sync.WaitGroup{}

	// open channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.openChan {
			fmt.Printf("\n\n[OpenResponse]\n\n")
		}
	}()

	// binary channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for br := range dch.binaryChan {
			fmt.Printf("\n\n[Binary Data]\n\n")
			fmt.Printf("Size: %d\n\n", len(*br))

			file, err := os.OpenFile(AUDIO_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				fmt.Printf("Failed to open file. Err: %v\n", err)
				continue
			}

			_, err = file.Write(*br)
			file.Close()

			if err != nil {
				fmt.Printf("Failed to write to file. Err: %v\n", err)
				continue
			}
		}
	}()

	// metadata channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for mr := range dch.metadataChan {
			fmt.Printf("\n[FlushedResponse]\n")
			fmt.Printf("RequestID: %s\n", strings.TrimSpace(mr.RequestID))
		}
	}()

	// flushed channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.flushChan {
			fmt.Printf("\n[FlushedResponse]\n")
		}
	}()

	// cleared channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.clearChan {
			fmt.Printf("\n[ClearedResponse]\n")
		}
	}()

	// close channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for _ = range dch.closeChan {
			fmt.Printf("\n\n[CloseResponse]\n\n")
		}
	}()

	// warning channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for er := range dch.warningChan {
			fmt.Printf("\n[WarningResponse]\n")
			fmt.Printf("\nWarning.Type: %s\n", er.WarnCode)
			fmt.Printf("Warning.Message: %s\n", er.WarnMsg)
			fmt.Printf("Warning.Description: %s\n\n", er.Description)
			fmt.Printf("Warning.Variant: %s\n\n", er.Variant)
		}
	}()

	// error channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for er := range dch.errorChan {
			fmt.Printf("\n[ErrorResponse]\n")
			fmt.Printf("\nError.Type: %s\n", er.ErrCode)
			fmt.Printf("Error.Message: %s\n", er.ErrMsg)
			fmt.Printf("Error.Description: %s\n\n", er.Description)
			fmt.Printf("Error.Variant: %s\n\n", er.Variant)
		}
	}()

	// unhandled event channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for byData := range dch.unhandledChan {
			fmt.Printf("\n[UnhandledEvent]")
			fmt.Printf("Dump:\n%s\n\n", string(*byData))
		}
	}()

	// wait for all receivers to finish
	wgReceivers.Wait()

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
	callback := NewMyHandler()

	// create a new stream using the NewStream function
	dgClient, err := speak.NewWSUsingChan(ctx, "", cOptions, ttsOptions, callback)
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
