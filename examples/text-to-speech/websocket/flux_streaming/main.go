// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates the Deepgram Flux TTS streaming WebSocket API
// (wss://api.deepgram.com/v2/speak) with the callback-based client. Text is sent
// with Speak, the turn is ended with Flush, and the synthesized audio arrives as
// binary frames which are streamed into output.wav.
//
// The session is closed gracefully with Finish: the server drains every queued
// turn, sends all remaining audio, and reports the final SessionMetadata before
// the socket closes — nothing is truncated. The WAV header's RIFF and data sizes
// are patched once the byte count is known, so the finished file is a conforming
// WAV that any player or parser accepts.
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
	wav "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/wav"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
	speakv2client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2"
)

const (
	audioFile  = "output.wav"
	sampleRate = 48000
)

var (
	model string
	text  string
)

func init() {
	flag.StringVar(&model, "model", "flux-haley-en", "Flux TTS model (flux-{voice}-{language}); Aura models are rejected on /v2/speak")
	flag.StringVar(&text, "text", "Hello! This is the Deepgram Flux text-to-speech API, streaming over a websocket.", "Text to synthesize")
}

// MyCallback implements api.FluxSpeakMessageCallback. Binary audio frames are
// streamed into the WAV writer; all events are printed.
type MyCallback struct {
	wavWriter *wav.Writer
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

	if _, err := c.wavWriter.Write(byMsg); err != nil {
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
		SampleRate: sampleRate,
	}

	// one file handle for the whole session: the WAV writer streams the raw
	// linear16 audio behind a standard header, and Finalize patches the RIFF and
	// data sizes once the total byte count is known
	file, err := os.OpenFile(audioFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		fmt.Printf("ERROR creating %s: %v\n", audioFile, err)
		os.Exit(1)
	}
	defer file.Close()

	wavWriter, err := wav.NewWriter(file, 1, sampleRate, 16)
	if err != nil {
		fmt.Printf("ERROR writing WAV header: %v\n", err)
		os.Exit(1)
	}

	callback := MyCallback{wavWriter: wavWriter}

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

	// close gracefully: Finish waits while the server drains all queued audio,
	// sends the final SessionMetadata, and closes the socket — so every audio
	// frame has reached the callback (and the WAV file) when it returns
	finishCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		fmt.Printf("ERROR closing session gracefully: %v\n", err)
		os.Exit(1)
	}

	// a session that produced no audio is a failure, not an empty-but-valid WAV
	if wavWriter.DataBytes() == 0 {
		fmt.Println("ERROR: the session produced no audio")
		file.Close()
		os.Remove(audioFile)
		os.Exit(1)
	}

	// patch the WAV header's RIFF and data sizes now that the length is known
	if err := wavWriter.Finalize(); err != nil {
		fmt.Printf("ERROR finalizing %s: %v\n", audioFile, err)
		os.Exit(1)
	}
	if err := file.Close(); err != nil {
		fmt.Printf("ERROR closing %s: %v\n", audioFile, err)
		os.Exit(1)
	}

	fmt.Printf("\nSession complete — %d bytes of audio saved to %s\n", wavWriter.DataBytes(), audioFile)
}
