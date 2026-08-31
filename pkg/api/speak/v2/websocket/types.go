// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
)

// CallbackRouter routes incoming Flux TTS server messages to a
// FluxSpeakMessageCallback implementation.
type CallbackRouter struct {
	debugWebsocket bool
	callback       interfaces.FluxSpeakMessageCallback
}

// ChanRouter routes incoming Flux TTS server messages to registered Go channels.
type ChanRouter struct {
	debugWebsocket        bool
	openChan              []*chan *interfaces.OpenResponse
	connectedChan         []*chan *interfaces.ConnectedResponse
	binaryChan            []*chan *[]byte
	speechStartedChan     []*chan *interfaces.SpeechStartedResponse
	speechMetadataChan    []*chan *interfaces.SpeechMetadataResponse
	speechInterruptedChan []*chan *interfaces.SpeechInterruptedResponse
	flushedChan           []*chan *interfaces.FlushedResponse
	sessionMetadataChan   []*chan *interfaces.SessionMetadataResponse
	configureSuccessChan  []*chan *interfaces.ConfigureSuccessResponse
	configureFailureChan  []*chan *interfaces.ConfigureFailureResponse
	warningChan           []*chan *interfaces.WarningResponse
	fatalErrorChan        []*chan *interfaces.FatalErrorResponse
	closeChan             []*chan *interfaces.CloseResponse
	errorChan             []*chan *interfaces.ErrorResponse
	unhandledChan         []*chan *[]byte
}

// DefaultCallbackHandler is a FluxSpeakMessageCallback that prints all events to
// stdout. Used when no callback is provided to the factory functions.
type DefaultCallbackHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool
}

// DefaultChanHandler is a FluxSpeakMessageChan that exposes all events as Go
// channels and prints them to stdout via its Run() goroutine.
// Used when no channel handler is provided to the factory functions.
type DefaultChanHandler struct {
	debugWebsocket        bool
	debugWebsocketVerbose bool

	openChan              chan *interfaces.OpenResponse
	connectedChan         chan *interfaces.ConnectedResponse
	binaryChan            chan *[]byte
	speechStartedChan     chan *interfaces.SpeechStartedResponse
	speechMetadataChan    chan *interfaces.SpeechMetadataResponse
	speechInterruptedChan chan *interfaces.SpeechInterruptedResponse
	flushedChan           chan *interfaces.FlushedResponse
	sessionMetadataChan   chan *interfaces.SessionMetadataResponse
	configureSuccessChan  chan *interfaces.ConfigureSuccessResponse
	configureFailureChan  chan *interfaces.ConfigureFailureResponse
	warningChan           chan *interfaces.WarningResponse
	fatalErrorChan        chan *interfaces.FatalErrorResponse
	closeChan             chan *interfaces.CloseResponse
	errorChan             chan *interfaces.ErrorResponse
	unhandledChan         chan *[]byte
}
