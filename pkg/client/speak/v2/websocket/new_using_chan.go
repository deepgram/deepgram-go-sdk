// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"

	klog "k8s.io/klog/v2"

	websocketv2api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
)

// NewUsingChanForDemo creates a Flux TTS WebSocket client with all default options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewUsingChanForDemo(ctx context.Context, options *clientinterfaces.SpeakV2WSOptions) (*WSChannel, error) {
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptionsV2{}, options, nil)
}

// NewUsingChanWithDefaults creates a Flux TTS WebSocket client with default connection options.
// The Deepgram API key is read from the DEEPGRAM_API_KEY environment variable.
func NewUsingChanWithDefaults(ctx context.Context, options *clientinterfaces.SpeakV2WSOptions, chans msginterfaces.FluxSpeakMessageChan) (*WSChannel, error) {
	return NewUsingChan(ctx, "", &clientinterfaces.ClientOptionsV2{}, options, chans)
}

// NewUsingChan creates a Flux TTS WebSocket client with the specified options.
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
// If chans is nil, the default stdout-printing channel handler is used.
func NewUsingChan(ctx context.Context, apiKey string, cOptions *clientinterfaces.ClientOptionsV2, sOptions *clientinterfaces.SpeakV2WSOptions, chans msginterfaces.FluxSpeakMessageChan) (*WSChannel, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewUsingChanWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, chans)
}

// NewUsingChanWithCancel creates a Flux TTS WebSocket client and lets the caller
// supply their own context cancel function (BYOC — Bring Your Own Cancel).
//
// If apiKey is empty, it is read from the DEEPGRAM_API_KEY environment variable.
// If chans is nil, the default stdout-printing channel handler is used.
func NewUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *clientinterfaces.ClientOptionsV2, sOptions *clientinterfaces.SpeakV2WSOptions, chans msginterfaces.FluxSpeakMessageChan) (*WSChannel, error) {
	klog.V(6).Infof("fluxspeak.NewUsingChanWithCancel() ENTER\n")

	if cOptions == nil {
		cOptions = &clientinterfaces.ClientOptionsV2{}
	}
	if sOptions == nil {
		klog.V(1).Infof("SpeakV2WSOptions is nil\n")
		return nil, interfacesv2.ErrOptionsRequired
	}
	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	if err := cOptions.Parse(); err != nil {
		klog.V(1).Infof("ClientOptions.Parse() failed. Err: %v\n", err)
		return nil, err
	}
	if err := sOptions.Check(); err != nil {
		klog.V(1).Infof("SpeakV2WSOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	usesDefaultHandler := chans == nil
	var defaultHandlerShutdown func()
	var defaultHandlerDone chan struct{}
	if usesDefaultHandler {
		klog.V(2).Infof("Using DefaultChanHandler.\n")
		handler := websocketv2api.NewDefaultChanHandler()
		defaultHandlerDone = make(chan struct{})
		go func() {
			defer close(defaultHandlerDone)
			if err := handler.Run(); err != nil {
				klog.V(1).Infof("DefaultChanHandler.Run failed. Err: %v\n", err)
			}
		}()
		chans = handler
		defaultHandlerShutdown = handler.Shutdown
	}

	conn := WSChannel{
		cOptions:  cOptions,
		sOptions:  sOptions,
		chans:     []*msginterfaces.FluxSpeakMessageChan{&chans},
		ctx:       ctx,
		ctxCancel: ctxCancel,
		finish:    newFinishState(),

		usesDefaultHandler:     usesDefaultHandler,
		defaultHandlerShutdown: defaultHandlerShutdown,
		defaultHandlerDone:     defaultHandlerDone,
	}

	// the observer exposes the user channels unchanged and adds internal ones that
	// mark the graceful-close milestones Finish(ctx) waits on
	observer := newChanFinishObserver(chans)
	observer.run(ctx, conn.currentFinishState(), false)
	conn.observer = observer
	var router commoninterfaces.Router
	router = websocketv2api.NewChanRouter(observer)
	conn.router = &router

	var handler commoninterfaces.WebSocketHandler
	handler = &wsHandlerShim{core: &conn}
	conn.WSClient = common.NewWS(ctx, ctxCancel, apiKey, cOptions, &handler, &router)

	klog.V(3).Infof("fluxspeak WSChannel created\n")
	klog.V(6).Infof("fluxspeak.NewUsingChanWithCancel() LEAVE\n")
	return &conn, nil
}

// restartDefaultHandler replaces a factory-owned handler that was closed for a
// previous session, so reconnects never route messages to closed channels.
func (c *WSChannel) restartDefaultHandler() {
	c.defaultHandlerMu.Lock()
	defer c.defaultHandlerMu.Unlock()

	if !c.usesDefaultHandler || c.defaultHandlerShutdown != nil {
		return
	}

	handler := websocketv2api.NewDefaultChanHandler()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := handler.Run(); err != nil {
			klog.V(1).Infof("DefaultChanHandler.Run failed. Err: %v\n", err)
		}
	}()

	chans := msginterfaces.FluxSpeakMessageChan(handler)
	*c.chans[0] = chans
	c.observer = newChanFinishObserver(chans)
	var router commoninterfaces.Router
	router = websocketv2api.NewChanRouter(c.observer)
	*c.router = router
	c.defaultHandlerShutdown = handler.Shutdown
	c.defaultHandlerDone = done
}
