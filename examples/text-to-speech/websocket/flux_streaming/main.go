// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates the Deepgram Flux TTS streaming WebSocket API
// (wss://api.deepgram.com/v2/speak) with the callback-based client. Text is sent
// with Speak, the turn is ended with Flush, and the synthesized audio arrives as
// binary frames which are appended to output.wav.
//
// Run:
//
//	DEEPGRAM_API_KEY=<your-key> go run main.go
//	DEEPGRAM_API_KEY=<your-key> go run main.go -model flux-alexis-en -text "Hello from Flux."
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
	speakv2client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2"
)

const audioFile = "output.wav"

var (
	model string
	text  string
)

func init() {
	flag.StringVar(&model, "model", "flux-haley-en", "Flux TTS model (flux-{voice}-{language}); Aura models are rejected on /v2/speak")
	flag.StringVar(&text, "text", "Hello! This is the Deepgram Flux text-to-speech API, streaming over a websocket.", "Text to synthesize")
}

// MyCallback implements api.FluxSpeakMessageCallback. Binary audio frames are
// appended to the output file; turn completion is signaled on metadataDone.
type MyCallback struct {
	metadataDone chan struct{}
}

func (c MyCallback) Open(or *api.OpenResponse) error {
	fmt.Printf("\n[Open] WebSocket connection established\n")
	return nil
}

func (c MyCallback) Connected(cr *api.ConnectedResponse) error {
	fmt.Printf("\n[Connected] request_id=%s model=%s version=%s\n", cr.RequestID, cr.ModelName, cr.ModelVersion)
	return nil
}

func (c MyCallback) Binary(byMsg []byte) error {
	fmt.Printf("[Binary] %d bytes of audio received\n", len(byMsg))

	file, err := os.OpenFile(audioFile, os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		fmt.Printf("ERROR opening %s: %v\n", audioFile, err)
		return err
	}
	defer file.Close()

	if _, err := file.Write(byMsg); err != nil {
		fmt.Printf("ERROR writing audio to %s: %v\n", audioFile, err)
		return err
	}
	return nil
}

func (c MyCallback) SpeechStarted(ss *api.SpeechStartedResponse) error {
	fmt.Printf("\n[SpeechStarted] speech_id=%s\n", ss.SpeechID)
	return nil
}

func (c MyCallback) SpeechMetadata(sm *api.SpeechMetadataResponse) error {
	fmt.Printf("\n[SpeechMetadata] speech_id=%s duration_ms=%d billable_chars=%d\n",
		sm.SpeechID, sm.AudioDurationMs, sm.BillableCharacterCount)
	// SpeechMetadata arrives after all audio for the turn has been sent
	close(c.metadataDone)
	return nil
}

func (c MyCallback) SpeechInterrupted(si *api.SpeechInterruptedResponse) error {
	fmt.Printf("\n[SpeechInterrupted] audio_played_ms=%d\n", si.AudioPlayedMs)
	return nil
}

func (c MyCallback) Flushed(fl *api.FlushedResponse) error {
	fmt.Printf("\n[Flushed] speech_id=%s\n", fl.SpeechID)
	return nil
}

func (c MyCallback) SessionMetadata(sm *api.SessionMetadataResponse) error {
	fmt.Printf("\n[SessionMetadata] total_duration_ms=%d total_billable_chars=%d\n",
		sm.TotalAudioDurationMs, sm.TotalBillableCharacterCount)
	return nil
}

func (c MyCallback) ConfigureSuccess(cs *api.ConfigureSuccessResponse) error {
	fmt.Printf("\n[ConfigureSuccess] Received\n")
	return nil
}

func (c MyCallback) ConfigureFailure(cf *api.ConfigureFailureResponse) error {
	fmt.Printf("\n[ConfigureFailure] code=%s description=%s\n", cf.Code, cf.Description)
	return nil
}

func (c MyCallback) Warning(wr *api.WarningResponse) error {
	fmt.Printf("\n[Warning] code=%s description=%s\n", wr.Code, wr.Description)
	return nil
}

func (c MyCallback) FatalError(fe *api.FatalErrorResponse) error {
	fmt.Printf("\n[FatalError] code=%s description=%s\n", fe.Code, fe.Description)
	return nil
}

func (c MyCallback) Close(cr *api.CloseResponse) error {
	fmt.Printf("\n[Close] WebSocket connection closed\n")
	return nil
}

func (c MyCallback) Error(er *api.ErrorResponse) error {
	fmt.Printf("\n[Error] type=%s message=%s description=%s\n", er.ErrCode, er.ErrMsg, er.Description)
	return nil
}

func (c MyCallback) UnhandledEvent(byData []byte) error {
	fmt.Printf("\n[UnhandledEvent] %s\n", string(byData))
	return nil
}

func main() {
	flag.Parse()

	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault,
	})

	ctx := context.Background()

	cOptions := &interfaces.ClientOptionsV2{
		EnableKeepAlive: true,
	}

	sOptions := &interfaces.SpeakV2WSOptions{
		Model:      model,
		Encoding:   "linear16",
		SampleRate: 48000,
	}

	// start the output file with a WAV container header so the raw linear16
	// audio can be played with a standard media player. Sizes are left as
	// placeholders, which most players accept for streamed WAV files.
	header := []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x00, 0x00, 0x00, 0x00, // Placeholder for file size
		0x57, 0x41, 0x56, 0x45, // "WAVE"
		0x66, 0x6d, 0x74, 0x20, // "fmt "
		0x10, 0x00, 0x00, 0x00, // Chunk size (16)
		0x01, 0x00, // Audio format (1 for PCM)
		0x01, 0x00, // Number of channels (1)
		0x80, 0xbb, 0x00, 0x00, // Sample rate (48000)
		0x00, 0x77, 0x01, 0x00, // Byte rate (48000 * 2)
		0x02, 0x00, // Block align (2)
		0x10, 0x00, // Bits per sample (16)
		0x64, 0x61, 0x74, 0x61, // "data"
		0x00, 0x00, 0x00, 0x00, // Placeholder for data size
	}
	file, err := os.OpenFile(audioFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		fmt.Printf("ERROR creating %s: %v\n", audioFile, err)
		os.Exit(1)
	}
	if _, err := file.Write(header); err != nil {
		fmt.Printf("ERROR writing WAV header: %v\n", err)
		os.Exit(1)
	}
	file.Close()

	callback := MyCallback{
		metadataDone: make(chan struct{}),
	}

	dgClient, err := speakv2client.NewWSUsingCallback(ctx, "", cOptions, sOptions, callback)
	if err != nil {
		fmt.Printf("ERROR creating Flux TTS client: %v\n", err)
		os.Exit(1)
	}

	if !dgClient.Connect() {
		fmt.Println("ERROR: failed to connect to Deepgram Flux TTS endpoint")
		os.Exit(1)
	}

	// send the text and end the turn
	if err := dgClient.Speak(text); err != nil {
		fmt.Printf("ERROR sending Speak: %v\n", err)
		os.Exit(1)
	}
	if err := dgClient.Flush(); err != nil {
		fmt.Printf("ERROR sending Flush: %v\n", err)
		os.Exit(1)
	}

	// wait for the turn to complete (SpeechMetadata arrives after all audio)
	select {
	case <-callback.metadataDone:
		fmt.Printf("\nTurn complete — audio saved to %s\n", audioFile)
	case <-time.After(60 * time.Second):
		fmt.Println("\nTimed out waiting for the turn to complete")
	}

	dgClient.Stop()

	fmt.Printf("\nProgram exiting...\n")
}
