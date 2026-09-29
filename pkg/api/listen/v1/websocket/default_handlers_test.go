// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"testing"
	"time"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
)

func Test_DefaultCallbackHandlerMessageEmptyAlternatives(t *testing.T) {
	handler := NewDefaultCallbackHandler()
	mr := &interfaces.MessageResponse{Type: string(interfaces.TypeMessageResponse)}

	if err := handler.Message(mr); err != nil {
		t.Fatalf("Message returned an error: %v", err)
	}
}

func Test_DefaultChanHandlerMessageEmptyAlternatives(t *testing.T) {
	handler := NewDefaultChanHandler()
	messages := handler.GetMessage()[0]

	// The handler goroutine must survive a message without alternatives and
	// keep receiving; an unbuffered send that times out means it stopped.
	for i := 0; i < 2; i++ {
		select {
		case *messages <- &interfaces.MessageResponse{Type: string(interfaces.TypeMessageResponse)}:
		case <-time.After(2 * time.Second):
			t.Fatalf("message %d was not received by the default handler", i)
		}
	}
}
