// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates using the Deepgram Flux (v2/listen) WebSocket API
// with the callback-based client. Flux is a turn-based audio transcription API
// that uses server-side turn detection.
//
// Run:
//
//	DEEPGRAM_API_KEY=<your-key> go run main.go -model flux-general-en
//	DEEPGRAM_API_KEY=<your-key> go run main.go -model flux-general-multi -language en -language es
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

// languageFlags is a custom flag type that accumulates multiple -language values.
type languageFlags []string

func (l *languageFlags) String() string { return strings.Join(*l, ",") }
func (l *languageFlags) Set(v string) error {
	for _, existing := range *l {
		if existing == v {
			return nil
		}
	}
	*l = append(*l, v)
	return nil
}

var (
	model         string
	languageHints languageFlags
)

func init() {
	flag.StringVar(&model, "model", "flux-general-en", "Flux model: flux-general-en or flux-general-multi")
	flag.Var(&languageHints, "language", "Language hint (repeat for multiple, e.g. -language en -language es); only applied with flux-general-multi")
}

// MyCallback implements api.FluxMessageCallback to handle all Flux server events.
type MyCallback struct{}

func (c MyCallback) Open(or *api.OpenResponse) error {
	fmt.Printf("\n[Open] WebSocket connection established\n")
	return nil
}

func (c MyCallback) Connected(cr *api.ConnectedResponse) error {
	fmt.Printf("\n[Connected] Flux session ready\n")
	fmt.Printf("  request_id:  %s\n", cr.RequestID)
	fmt.Printf("  sequence_id: %d\n\n", cr.SequenceID)
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

	case api.TurnEventEagerEndOfTurn:
		fmt.Printf("[Turn %d] EagerEndOfTurn: %s\n", tr.TurnIndex, strings.TrimSpace(tr.Transcript))

	case api.TurnEventTurnResumed:
		fmt.Printf("[Turn %d] TurnResumed (speech continued)\n", tr.TurnIndex)

	case api.TurnEventEndOfTurn:
		fmt.Printf("[Turn %d] FINAL: %s\n", tr.TurnIndex, strings.TrimSpace(tr.Transcript))
		if len(tr.Languages) > 0 {
			fmt.Printf("  detected languages:  %s\n", strings.Join(tr.Languages, ", "))
		}
		if len(tr.LanguagesHinted) > 0 {
			fmt.Printf("  hinted languages:    %s\n", strings.Join(tr.LanguagesHinted, ", "))
		}
		fmt.Println()
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

	fmt.Print("\n\nFlux (v2/listen) callback example — press ENTER to exit\n\n")

	listen.Init(listen.InitLib{
		LogLevel: listen.LogLevelDefault,
	})

	ctx := context.Background()

	cOptions := &interfaces.ClientOptionsV2{
		EnableKeepAlive: true,
	}

	tOptions := &interfaces.FluxTranscriptionOptions{
		Model:             model,
		Encoding:          "linear16",
		SampleRate:        16000,
		EagerEotThreshold: 0.3,
	}

	callback := MyCallback{}

	dgClient, err := client.NewWSUsingCallback(ctx, "", cOptions, tOptions, callback)
	if err != nil {
		fmt.Printf("ERROR creating Flux client: %v\n", err)
		return
	}

	if !dgClient.Connect() {
		fmt.Println("ERROR: failed to connect to Deepgram Flux endpoint")
		os.Exit(1)
	}

	if model == "flux-general-multi" && len(languageHints) > 0 {
		fmt.Printf("Configuring language hints: %s\n", strings.Join(languageHints, ", "))
		if err := dgClient.Configure(&interfaces.FluxConfigureOptions{
			LanguageHints: languageHints,
		}); err != nil {
			fmt.Printf("Configure error: %v\n", err)
		}
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		mic.Stream(dgClient)
	}()

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	if err = mic.Stop(); err != nil {
		fmt.Printf("ERROR stopping microphone: %v\n", err)
	}
	wg.Wait()

	microphone.Teardown()
	dgClient.Stop()

	fmt.Printf("\nProgram exiting...\n")
}
