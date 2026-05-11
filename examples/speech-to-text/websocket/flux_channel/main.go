// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package main demonstrates using the Deepgram Flux (v2/listen) WebSocket API
// with the channel-based client. Events are received on Go channels and processed
// in dedicated goroutines inside MyHandler.Run().
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

// MyHandler implements api.FluxMessageChan. Each getter returns the channels that
// the router will send events to. Return nil for events you don't care about.
type MyHandler struct {
	openChan             chan *api.OpenResponse
	connectedChan        chan *api.ConnectedResponse
	turnInfoChan         chan *api.TurnInfoResponse
	configureSuccessChan chan *api.ConfigureSuccessResponse
	configureFailureChan chan *api.ConfigureFailureResponse
	fatalErrorChan       chan *api.FatalErrorResponse
	closeChan            chan *api.CloseResponse
	errorChan            chan *api.ErrorResponse
	unhandledChan        chan *[]byte
}

func NewMyHandler() *MyHandler {
	handler := &MyHandler{
		openChan:             make(chan *api.OpenResponse),
		connectedChan:        make(chan *api.ConnectedResponse),
		turnInfoChan:         make(chan *api.TurnInfoResponse),
		configureSuccessChan: make(chan *api.ConfigureSuccessResponse),
		configureFailureChan: make(chan *api.ConfigureFailureResponse),
		fatalErrorChan:       make(chan *api.FatalErrorResponse),
		closeChan:            make(chan *api.CloseResponse),
		errorChan:            make(chan *api.ErrorResponse),
		unhandledChan:        make(chan *[]byte),
	}
	go func() {
		handler.Run()
	}()
	return handler
}

func (h *MyHandler) GetOpen() []*chan *api.OpenResponse {
	return []*chan *api.OpenResponse{&h.openChan}
}
func (h *MyHandler) GetConnected() []*chan *api.ConnectedResponse {
	return []*chan *api.ConnectedResponse{&h.connectedChan}
}
func (h *MyHandler) GetTurnInfo() []*chan *api.TurnInfoResponse {
	return []*chan *api.TurnInfoResponse{&h.turnInfoChan}
}
func (h *MyHandler) GetConfigureSuccess() []*chan *api.ConfigureSuccessResponse {
	return []*chan *api.ConfigureSuccessResponse{&h.configureSuccessChan}
}
func (h *MyHandler) GetConfigureFailure() []*chan *api.ConfigureFailureResponse {
	return []*chan *api.ConfigureFailureResponse{&h.configureFailureChan}
}
func (h *MyHandler) GetFatalError() []*chan *api.FatalErrorResponse {
	return []*chan *api.FatalErrorResponse{&h.fatalErrorChan}
}
func (h *MyHandler) GetClose() []*chan *api.CloseResponse {
	return []*chan *api.CloseResponse{&h.closeChan}
}
func (h *MyHandler) GetError() []*chan *api.ErrorResponse {
	return []*chan *api.ErrorResponse{&h.errorChan}
}
func (h *MyHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandledChan}
}

// Run starts a goroutine for each event channel and blocks until all channels are closed.
// It is launched automatically by NewMyHandler.
func (h *MyHandler) Run() error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range h.openChan {
			fmt.Printf("\n[Open] WebSocket connection established\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cr := range h.connectedChan {
			fmt.Printf("\n[Connected] Flux session ready\n")
			fmt.Printf("  request_id:  %s\n", cr.RequestID)
			fmt.Printf("  sequence_id: %d\n\n", cr.SequenceID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for tr := range h.turnInfoChan {
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
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range h.configureSuccessChan {
			fmt.Printf("\n[ConfigureSuccess] Mid-session config accepted\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cf := range h.configureFailureChan {
			fmt.Printf("\n[ConfigureFailure] request_id=%s sequence_id=%d\n", cf.RequestID, cf.SequenceID)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for fe := range h.fatalErrorChan {
			fmt.Printf("\n[FatalError] code=%s description=%s\n", fe.Code, fe.Description)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range h.closeChan {
			fmt.Printf("\n[Close] WebSocket connection closed\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for er := range h.errorChan {
			fmt.Printf("\n[Error]\n")
			fmt.Printf("Error.Type: %s\n", er.ErrCode)
			fmt.Printf("Error.Message: %s\n", er.ErrMsg)
			fmt.Printf("Error.Description: %s\n\n", er.Description)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for byData := range h.unhandledChan {
			fmt.Printf("\n[UnhandledEvent] %s\n", string(*byData))
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	flag.Parse()

	microphone.Initialize()

	fmt.Print("\n\nFlux (v2/listen) channel example — press ENTER to exit\n\n")

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

	handler := NewMyHandler()

	dgClient, err := client.NewWSUsingChan(ctx, "", cOptions, tOptions, handler)
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

	var micWg sync.WaitGroup
	micWg.Add(1)
	go func() {
		defer micWg.Done()
		mic.Stream(dgClient)
	}()

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	if err = mic.Stop(); err != nil {
		fmt.Printf("ERROR stopping microphone: %v\n", err)
	}
	micWg.Wait()

	microphone.Teardown()
	dgClient.Stop()

	fmt.Printf("\nProgram exiting...\n")
}
