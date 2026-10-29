// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	listenws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

type configureHandler struct {
	listenws.DefaultCallbackHandler
	result chan error
}

func (h *configureHandler) ConfigureSuccess(response *msginterfaces.ConfigureSuccessResponse) error {
	fmt.Printf("Configure succeeded for request %s: %v\n", response.RequestID, response.Keyterms)
	h.result <- nil
	return nil
}

func (h *configureHandler) ConfigureFailure(response *msginterfaces.ConfigureFailureResponse) error {
	h.result <- fmt.Errorf("Configure failed for request %s: %s: %s", response.RequestID, response.Code, response.Description)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	handler := &configureHandler{
		DefaultCallbackHandler: *listenws.NewDefaultCallbackHandler(),
		result:                 make(chan error, 1),
	}

	conn, err := listen.NewWSUsingCallbackWithDefaults(ctx, &interfaces.LiveTranscriptionOptions{
		Model: "nova-3",
	}, handler)
	if err != nil {
		return err
	}
	defer conn.Stop()

	if !conn.Connect() {
		return errors.New("connect failed")
	}

	keyterms := []string{"Deepgram"}
	if err := conn.Configure(&interfaces.LiveConfigureOptions{Keyterms: &keyterms}); err != nil {
		return err
	}

	select {
	case err := <-handler.result:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("timed out waiting for ConfigureSuccess or ConfigureFailure")
	}
}
