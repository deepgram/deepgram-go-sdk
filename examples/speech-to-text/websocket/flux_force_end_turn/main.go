// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates bring-your-own turn detection with the Deepgram Flux
// (v2/listen) WebSocket API. It streams microphone audio with eot_threshold=1.0
// (native end-of-turn detection suppressed) and ends each turn manually by sending
// a ForceEndTurn control message when you press ENTER — the same pattern you would
// use for push-to-talk, DTMF input, or a UI send button.
//
// Each EndOfTurn event prints its Trigger field: "manual" for turns you ended with
// ENTER, "timeout" if eot_timeout_ms elapsed first, and "model" for turns ended by
// native detection (only when -eot-threshold is set below 1.0).
//
// NOTE: ForceEndTurn is gated per deployment on the Deepgram side. If your project
// does not have it enabled yet, contact Deepgram support.
//
// Run:
//
//	DEEPGRAM_API_KEY=<your-key> go run main.go
//	DEEPGRAM_API_KEY=<your-key> go run main.go -eot-threshold 0.7   # blend manual + native detection
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	microphone "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/microphone"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2"
)

var (
	model        string
	eotThreshold float64
)

func init() {
	flag.StringVar(&model, "model", "flux-general-en", "Flux model: flux-general-en or flux-general-multi")
	flag.Float64Var(&eotThreshold, "eot-threshold", 1.0, "End-of-turn confidence threshold (0.5-1.0); 1.0 suppresses native detection so turns end only via ForceEndTurn or timeout")
}

// MyCallback implements api.FluxMessageCallback to handle all Flux server events.
type MyCallback struct{}

func (c MyCallback) Open(or *api.OpenResponse) error {
	fmt.Printf("\n[Open] WebSocket connection established\n")
	return nil
}

func (c MyCallback) Connected(cr *api.ConnectedResponse) error {
	fmt.Printf("\n[Connected] Flux session ready (request_id: %s)\n\n", cr.RequestID)
	fmt.Printf("Speak, then press ENTER to end the turn. Type q + ENTER to exit.\n\n")
	return nil
}

func (c MyCallback) TurnInfo(tr *api.TurnInfoResponse) error {
	switch tr.EventType {
	case api.TurnEventStartOfTurn:
		fmt.Printf("[Turn %d] StartOfTurn\n", tr.TurnIndex)

	case api.TurnEventUpdate:
		sentence := strings.TrimSpace(tr.Transcript)
		if sentence != "" {
			fmt.Printf("[Turn %d] Interim: %s\n", tr.TurnIndex, sentence)
		}

	case api.TurnEventEndOfTurn:
		fmt.Printf("[Turn %d] FINAL (trigger=%s): %s\n\n", tr.TurnIndex, tr.Trigger, strings.TrimSpace(tr.Transcript))
	}
	return nil
}

func (c MyCallback) ConfigureSuccess(cs *api.ConfigureSuccessResponse) error {
	fmt.Printf("\n[ConfigureSuccess] Mid-session config accepted\n")
	return nil
}

func (c MyCallback) ConfigureFailure(cf *api.ConfigureFailureResponse) error {
	fmt.Printf("\n[ConfigureFailure] request_id=%s sequence_id=%d\n", cf.RequestID, cf.SequenceID)
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

	microphone.Initialize()

	fmt.Print("\n\nFlux (v2/listen) ForceEndTurn example\n")

	listen.Init(listen.InitLib{
		LogLevel: listen.LogLevelDefault,
	})

	ctx := context.Background()

	cOptions := &interfaces.ClientOptionsV2{
		EnableKeepAlive: true,
	}

	tOptions := &interfaces.FluxTranscriptionOptions{
		Model:        model,
		Encoding:     "linear16",
		SampleRate:   16000,
		EotThreshold: eotThreshold,
	}

	dgClient, err := client.NewWSUsingCallback(ctx, "", cOptions, tOptions, MyCallback{})
	if err != nil {
		fmt.Printf("ERROR creating Flux client: %v\n", err)
		return
	}

	if !dgClient.Connect() {
		fmt.Println("ERROR: failed to connect to Deepgram Flux endpoint")
		os.Exit(1)
	}

	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  16000,
	})
	if err != nil {
		fmt.Printf("ERROR initializing microphone: %v\n", err)
		os.Exit(1)
	}

	if err = mic.Start(); err != nil {
		fmt.Printf("ERROR starting microphone: %v\n", err)
		os.Exit(1)
	}

	var micWg sync.WaitGroup
	micWg.Add(1)
	go func() {
		defer micWg.Done()
		if streamErr := mic.Stream(dgClient); streamErr != nil {
			fmt.Printf("ERROR streaming microphone audio: %v\n", streamErr)
		}
	}()

	// ENTER ends the current turn manually; q + ENTER exits.
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if strings.EqualFold(strings.TrimSpace(input.Text()), "q") {
			break
		}
		if err := dgClient.ForceEndTurn(); err != nil {
			fmt.Printf("ERROR sending ForceEndTurn: %v\n", err)
			continue
		}
		fmt.Printf("-- ForceEndTurn sent --\n")
	}

	if err = mic.Stop(); err != nil {
		fmt.Printf("ERROR stopping microphone: %v\n", err)
	}
	micWg.Wait()

	microphone.Teardown()
	dgClient.Stop()

	fmt.Printf("\nProgram exiting...\n")
}
