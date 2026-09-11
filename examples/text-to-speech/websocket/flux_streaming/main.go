// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates the Deepgram Flux TTS streaming WebSocket API
// (wss://api.deepgram.com/v2/speak) with the callback-based client. Text is sent
// with Speak, the turn is ended with Flush, and the synthesized audio arrives as
// binary frames which are streamed into output.wav.
//
// When Finish returns nil before its context deadline, the server drains every
// queued turn, sends all remaining audio, and reports the final SessionMetadata before
// the socket closes — nothing is truncated. The WAV header's RIFF and data sizes
// are then patched, so the finished file is a conforming WAV that any player or
// parser accepts. If Finish times out or fails, no output file is written.
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
	"path/filepath"
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

func (c *MyCallback) Open(or *api.OpenResponse) error {
	fmt.Printf("\n[Open] WebSocket connection established\n")
	return nil
}

func (c *MyCallback) Connected(cr *api.ConnectedResponse) error {
	fmt.Printf("\n[Connected] request_id=%s model=%s version=%s\n", cr.RequestID, cr.ModelName, cr.ModelVersion)
	return nil
}

func (c *MyCallback) Binary(byMsg []byte) error {
	fmt.Printf("[Binary] %d bytes of audio received\n", len(byMsg))

	if _, err := c.wavWriter.Write(byMsg); err != nil {
		fmt.Printf("ERROR writing audio to %s: %v\n", audioFile, err)
		return err
	}
	return nil
}

func (c *MyCallback) SpeechStarted(ss *api.SpeechStartedResponse) error {
	fmt.Printf("\n[SpeechStarted] speech_id=%s\n", ss.SpeechID)
	return nil
}

func (c *MyCallback) SpeechMetadata(sm *api.SpeechMetadataResponse) error {
	fmt.Printf("\n[SpeechMetadata] speech_id=%s duration_ms=%d billable_chars=%d\n",
		sm.SpeechID, sm.AudioDurationMs, sm.BillableCharacterCount)
	return nil
}

func (c *MyCallback) SpeechInterrupted(si *api.SpeechInterruptedResponse) error {
	fmt.Printf("\n[SpeechInterrupted] audio_played_ms=%d\n", si.AudioPlayedMs)
	return nil
}

func (c *MyCallback) Flushed(fl *api.FlushedResponse) error {
	fmt.Printf("\n[Flushed] speech_id=%s\n", fl.SpeechID)
	return nil
}

func (c *MyCallback) SessionMetadata(sm *api.SessionMetadataResponse) error {
	fmt.Printf("\n[SessionMetadata] total_duration_ms=%d total_billable_chars=%d\n",
		sm.TotalAudioDurationMs, sm.TotalBillableCharacterCount)
	return nil
}

func (c *MyCallback) ConfigureSuccess(cs *api.ConfigureSuccessResponse) error {
	fmt.Printf("\n[ConfigureSuccess] Received\n")
	return nil
}

func (c *MyCallback) ConfigureFailure(cf *api.ConfigureFailureResponse) error {
	fmt.Printf("\n[ConfigureFailure] code=%s description=%s\n", cf.Code, cf.Description)
	return nil
}

func (c *MyCallback) Warning(wr *api.WarningResponse) error {
	fmt.Printf("\n[Warning] code=%s description=%s\n", wr.Code, wr.Description)
	return nil
}

func (c *MyCallback) FatalError(fe *api.FatalErrorResponse) error {
	fmt.Printf("\n[FatalError] code=%s description=%s\n", fe.Code, fe.Description)
	return nil
}

func (c *MyCallback) Close(cr *api.CloseResponse) error {
	fmt.Printf("\n[Close] WebSocket connection closed\n")
	return nil
}

func (c *MyCallback) Error(er *api.ErrorResponse) error {
	fmt.Printf("\n[Error] type=%s message=%s description=%s\n", er.ErrCode, er.ErrMsg, er.Description)
	return nil
}

func (c *MyCallback) UnhandledEvent(byData []byte) error {
	fmt.Printf("\n[UnhandledEvent] %s\n", string(byData))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

type temporaryOutput struct {
	file      *os.File
	temporary string
	target    string
	closed    bool
}

func newTemporaryOutput(target string) (*temporaryOutput, error) {
	file, err := os.CreateTemp(filepath.Dir(target), filepath.Base(target)+".*.tmp")
	if err != nil {
		return nil, err
	}
	return &temporaryOutput{file: file, temporary: file.Name(), target: target}, nil
}

func (o *temporaryOutput) discard() {
	if !o.closed {
		_ = o.file.Close()
		o.closed = true
	}
	_ = os.Remove(o.temporary)
}

func (o *temporaryOutput) commit() error {
	if !o.closed {
		if err := o.file.Close(); err != nil {
			return err
		}
		o.closed = true
	}
	return replaceOutput(o.temporary, o.target)
}

func run() error {
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

	callback := &MyCallback{}

	dgClient, err := speakv2client.NewWSUsingCallback(ctx, "", cOptions, sOptions, callback)
	if err != nil {
		return fmt.Errorf("creating Flux TTS client: %w", err)
	}

	if !dgClient.Connect() {
		return fmt.Errorf("connecting to Deepgram Flux TTS endpoint")
	}

	// Write to a sibling file until the full WAV is finalized. A failed run cannot
	// corrupt an existing output.wav.
	output, err := newTemporaryOutput(audioFile)
	if err != nil {
		return fmt.Errorf("creating %s: %w", audioFile, err)
	}
	completed := false
	defer func() {
		if !completed {
			output.discard()
		}
	}()

	wavWriter, err := wav.NewWriter(output.file, 1, sampleRate, 16)
	if err != nil {
		return fmt.Errorf("writing WAV header: %w", err)
	}
	callback.wavWriter = wavWriter

	// send the text and end the turn
	if err := dgClient.Speak(text); err != nil {
		return fmt.Errorf("sending Speak: %w", err)
	}
	if err := dgClient.Flush(); err != nil {
		return fmt.Errorf("sending Flush: %w", err)
	}

	// close gracefully: Finish waits while the server drains all queued audio,
	// sends the final SessionMetadata, and closes the socket — so every audio
	// frame has reached the callback (and the WAV file) when it returns
	finishCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		return fmt.Errorf("closing session gracefully: %w", err)
	}

	// a session that produced no audio is a failure, not an empty-but-valid WAV
	if wavWriter.DataBytes() == 0 {
		return fmt.Errorf("the session produced no audio")
	}

	// patch the WAV header's RIFF and data sizes now that the length is known
	if err := wavWriter.Finalize(); err != nil {
		return fmt.Errorf("finalizing %s: %w", audioFile, err)
	}
	if err := output.commit(); err != nil {
		return fmt.Errorf("saving %s: %w", audioFile, err)
	}
	completed = true

	fmt.Printf("\nSession complete — %d bytes of audio saved to %s\n", wavWriter.DataBytes(), audioFile)
	return nil
}
