// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

type blockingCloseRouter struct {
	closeStarted chan struct{}
	closeRelease chan struct{}
	closeOnce    sync.Once
}

func (r *blockingCloseRouter) Open(*commoninterfaces.OpenResponse) error { return nil }
func (r *blockingCloseRouter) Message([]byte) error                      { return nil }
func (r *blockingCloseRouter) Binary([]byte) error                       { return nil }
func (r *blockingCloseRouter) Error(*commoninterfaces.ErrorResponse) error {
	return nil
}
func (r *blockingCloseRouter) Close(*commoninterfaces.CloseResponse) error {
	r.closeOnce.Do(func() { close(r.closeStarted) })
	<-r.closeRelease
	return nil
}

func newDefaultHandlerLifecycleClient(t *testing.T, onFrame func(*websocket.Conn, string), onConnect func(*websocket.Conn)) *WSChannel {
	return newDefaultHandlerLifecycleClientWithContext(t, context.Background(), onFrame, onConnect)
}

func newDefaultHandlerLifecycleClientWithContext(t *testing.T, ctx context.Context, onFrame func(*websocket.Conn, string), onConnect func(*websocket.Conn)) *WSChannel {
	t.Helper()

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		if onConnect != nil {
			onConnect(conn)
		}
		for {
			msgType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if msgType == websocket.TextMessage && onFrame != nil {
				onFrame(conn, string(message))
			}
		}
	}))
	t.Cleanup(srv.Close)

	client, err := NewUsingChan(ctx, "test", &clientinterfaces.ClientOptionsV2{
		Host:           strings.TrimPrefix(srv.URL, "https://"),
		SkipServerAuth: true,
	}, &clientinterfaces.SpeakV2WSOptions{
		Model:      "test-model",
		Encoding:   "linear16",
		SampleRate: 48000,
	}, nil)
	if err != nil {
		t.Fatalf("NewUsingChan failed: %v", err)
	}
	if !client.Connect() {
		t.Fatal("Connect failed")
	}
	return client
}

func awaitDefaultHandlerStopped(t *testing.T, client *WSChannel) {
	t.Helper()
	select {
	case <-client.defaultHandlerDone:
	case <-time.After(2 * time.Second):
		t.Fatal("default channel handler did not stop")
	}
}

func Test_KeepaliveStateReplacesAndStopsPingers(t *testing.T) {
	var keepalive keepaliveState
	firstStarted := make(chan struct{})
	firstStopped := make(chan struct{})
	secondStarted := make(chan struct{})
	secondStopped := make(chan struct{})

	keepalive.restart(context.Background(), func(ctx context.Context) {
		close(firstStarted)
		<-ctx.Done()
		close(firstStopped)
	})
	<-firstStarted

	keepalive.restart(context.Background(), func(ctx context.Context) {
		close(secondStarted)
		<-ctx.Done()
		close(secondStopped)
	})
	<-firstStopped
	<-secondStarted

	keepalive.stop(true)
	select {
	case <-secondStopped:
	case <-time.After(time.Second):
		t.Fatal("active keepalive did not stop")
	}
}

func Test_WSChannelDefaultHandlerStopsWithLifecycle(t *testing.T) {
	t.Run("Stop", func(t *testing.T) {
		client := newDefaultHandlerLifecycleClient(t, nil, nil)
		client.Stop()
		client.Stop()
		awaitDefaultHandlerStopped(t, client)
	})

	t.Run("Finish", func(t *testing.T) {
		client := newDefaultHandlerLifecycleClient(t, func(conn *websocket.Conn, frame string) {
			if !strings.Contains(frame, `"Close"`) {
				return
			}
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"SessionMetadata","total_audio_duration_ms":0,"total_input_character_count":0,"total_billable_character_count":0}`))
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
		}, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := client.Finish(ctx); err != nil {
			t.Fatalf("Finish failed: %v", err)
		}
		awaitDefaultHandlerStopped(t, client)
	})

	t.Run("peer close", func(t *testing.T) {
		client := newDefaultHandlerLifecycleClient(t, nil, func(conn *websocket.Conn) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			}()
		})
		awaitDefaultHandlerStopped(t, client)

		oldDone := client.defaultHandlerDone
		if !client.AttemptReconnect(context.Background(), 3) {
			t.Fatal("AttemptReconnect failed")
		}
		if client.defaultHandlerDone == oldDone {
			t.Fatal("reconnect reused the closed default handler")
		}
		awaitDefaultHandlerStopped(t, client)
	})
}

func Test_WSChannelDefaultHandlerContextCancellation(t *testing.T) {
	serverClosed := make(chan struct{})
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		_, _, _ = conn.ReadMessage()
		close(serverClosed)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := NewUsingChan(ctx, "test", &clientinterfaces.ClientOptionsV2{
		Host:           strings.TrimPrefix(srv.URL, "https://"),
		SkipServerAuth: true,
	}, &clientinterfaces.SpeakV2WSOptions{
		Model:      "test-model",
		Encoding:   "linear16",
		SampleRate: 48000,
	}, nil)
	if err != nil {
		t.Fatalf("NewUsingChan failed: %v", err)
	}
	if !client.Connect() {
		t.Fatal("Connect failed")
	}

	cancel()
	select {
	case <-serverClosed:
	case <-time.After(2 * time.Second):
		t.Fatal("context cancellation did not close the connection")
	}
	awaitDefaultHandlerStopped(t, client)

	// Context cancellation routes Close before shutting down the owned channels.
	// Stop must remain safe after that terminal cleanup.
	client.Stop()
}

func Test_WSChannelStopWaitsForContextCancellationClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := newDefaultHandlerLifecycleClientWithContext(t, ctx, nil, nil)

	// The common client and the Flux TTS channel client share this router pointer.
	// Block cancellation after it owns terminal teardown but before it can call
	// onFinish, which reproduces the former Stop race deterministically.
	blockingRouter := &blockingCloseRouter{
		closeStarted: make(chan struct{}),
		closeRelease: make(chan struct{}),
	}
	var releaseClose sync.Once
	release := func() { releaseClose.Do(func() { close(blockingRouter.closeRelease) }) }
	defer release()
	var router commoninterfaces.Router = blockingRouter
	*client.router = router

	cancel()
	select {
	case <-blockingRouter.closeStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("context cancellation did not route Close")
	}

	stopDone := make(chan struct{})
	go func() {
		client.Stop()
		close(stopDone)
	}()

	select {
	case <-client.defaultHandlerDone:
		t.Fatal("Stop shut down the default handler before cancellation routed Close")
	case <-time.After(100 * time.Millisecond):
	}
	select {
	case <-stopDone:
		t.Fatal("Stop returned before cancellation terminal teardown completed")
	case <-time.After(100 * time.Millisecond):
	}

	release()
	select {
	case <-stopDone:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop did not return after cancellation terminal teardown completed")
	}
	awaitDefaultHandlerStopped(t, client)
}
