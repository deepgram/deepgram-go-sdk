// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates the Deepgram Flux TTS batch (REST) API: POST /v2/speak.
// It synthesizes a complete block of text in a single request and saves the audio
// to an mp3 file — the pre-rendering path for fixed audio like IVR prompts,
// notifications, or narration.
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

	speakv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
	speakv2client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2"
)

const outputFile = "output.mp3"

var (
	model string
	text  string
)

func init() {
	flag.StringVar(&model, "model", "flux-haley-en", "Flux TTS model (flux-{voice}-{language}); Aura models are rejected on /v2/speak")
	flag.StringVar(&text, "text", "Your appointment is confirmed for 3pm tomorrow.", "Text to synthesize")
}

func main() {
	flag.Parse()

	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelDefault,
	})

	// context
	ctx := context.Background()

	// create the Flux TTS batch (REST) client
	dg, err := speakv2client.NewRESTWithDefaults()
	if err != nil {
		fmt.Printf("ERROR creating Flux TTS batch client: %v\n", err)
		os.Exit(1)
	}
	speakClient := speakv2.New(dg)

	options := &interfaces.SpeakV2Options{
		Model: model,
		// Encoding defaults to mp3 on the batch transport
	}

	resp, err := speakClient.ToSave(ctx, outputFile, text, options)
	if err != nil {
		fmt.Printf("Speech synthesis failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Saved synthesized audio to %s\n", resp.Filename)
	fmt.Printf("  request_id: %s\n", resp.RequestID)
	fmt.Printf("  model:      %s\n", resp.ModelName)
	fmt.Printf("  characters: %d\n", resp.Characters)
	if resp.Warnings != "" {
		fmt.Printf("  warnings:   %s\n", resp.Warnings)
	}
}
