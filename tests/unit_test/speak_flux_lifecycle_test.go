// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"testing"
	"time"

	speakapiws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket"
	speakmsg "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
)

// Test_FluxSpeakDefaultChanHandlerLifecycle covers Greg's S7 on PR #351: the
// default channel handler's constructor no longer starts Run itself (so the
// factory paths that start it once cannot create duplicate consumer sets), events
// flow through a single Run, and Shutdown terminates Run's goroutines.
func Test_FluxSpeakDefaultChanHandlerLifecycle(t *testing.T) {
	handler := speakapiws.NewDefaultChanHandler()

	// the constructor must NOT have started consumers: with unbuffered channels, a
	// send only proceeds once our own Run starts receiving
	select {
	case *handler.GetFlushed()[0] <- &speakmsg.FlushedResponse{Type: "Flushed", SpeechID: testSpeechID}:
		t.Fatal("constructor started its own consumer: an unbuffered send succeeded before Run")
	case <-time.After(100 * time.Millisecond):
	}

	runDone := make(chan error, 1)
	go func() {
		runDone <- handler.Run()
	}()

	// with Run started (once), events are consumed
	select {
	case *handler.GetFlushed()[0] <- &speakmsg.FlushedResponse{Type: "Flushed", SpeechID: testSpeechID}:
	case <-time.After(2 * time.Second):
		t.Fatal("Run is not consuming events")
	}

	// Shutdown closes the owned channels and Run returns
	handler.Shutdown()
	select {
	case err := <-runDone:
		if err != nil {
			t.Errorf("Run returned an error after Shutdown: %s", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after Shutdown")
	}

	// Shutdown is idempotent
	handler.Shutdown()
}
