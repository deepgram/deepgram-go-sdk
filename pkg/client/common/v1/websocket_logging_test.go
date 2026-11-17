// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

var (
	klogFlagsOnce sync.Once
	klogOutputMu  sync.Mutex
)

type lockedBuffer struct {
	mu             sync.Mutex
	buffer         bytes.Buffer
	listenerExited chan struct{}
	listenerOnce   sync.Once
}

func (b *lockedBuffer) Write(data []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	n, err := b.buffer.Write(data)
	if strings.Contains(string(data), "common.listen() LEAVE") {
		b.listenerOnce.Do(func() {
			close(b.listenerExited)
		})
	}
	return n, err
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.String()
}

func (b *lockedBuffer) waitForListener(t *testing.T) {
	t.Helper()

	select {
	case <-b.listenerExited:
	case <-time.After(5 * time.Second):
		t.Fatal("websocket listener did not stop")
	}
}

func Test_WSClientDoesNotLogPayloads(t *testing.T) {
	const secret = "websocket-log-secret"

	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		for i := 0; i < 3; i++ {
			if _, _, err := conn.ReadMessage(); err != nil {
				t.Errorf("read websocket message: %v", err)
				return
			}
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := &contextCloseHandler{
		url:      "ws" + strings.TrimPrefix(server.URL, "http"),
		finished: make(chan struct{}),
	}
	router := &contextCloseRouter{closed: make(chan struct{})}
	var handlerInterface commoninterfaces.WebSocketHandler = handler
	var routerInterface commoninterfaces.Router = router
	client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)
	if !client.ConnectWithCancel(ctx, cancel, 1) {
		t.Fatal("connect websocket")
	}

	output, actionErr := captureKlogOutput(t, func() error {
		if err := client.WriteJSON(map[string]string{"authorization": secret}); err != nil {
			return err
		}
		if err := client.WriteRawJSON(json.RawMessage(`{"authorization":"` + secret + `"}`)); err != nil {
			return err
		}
		if err := client.WriteBinary([]byte(secret)); err != nil {
			return err
		}
		if err := client.WriteJSON(map[string]interface{}{
			"authorization": secret,
			"invalid":       math.Inf(1),
		}); err == nil {
			return errors.New("write JSON with an invalid value succeeded")
		}
		return nil
	})
	client.Stop()
	output.waitForListener(t)
	if actionErr != nil {
		t.Fatalf("write payload: %v", actionErr)
	}
	if strings.Contains(output.String(), secret) || strings.Contains(output.String(), hex.EncodeToString([]byte(secret))) {
		t.Fatal("payload leaked to klog output")
	}
}

func captureKlogOutput(t *testing.T, action func() error) (*lockedBuffer, error) {
	t.Helper()

	klogOutputMu.Lock()
	t.Cleanup(klogOutputMu.Unlock)
	klogFlagsOnce.Do(func() {
		if flag.Lookup("v") == nil {
			klog.InitFlags(nil)
		}
	})

	verbosity := flag.Lookup("v")
	logToStderr := flag.Lookup("logtostderr")
	if verbosity == nil || logToStderr == nil {
		t.Fatal("klog flags are not registered")
	}
	previousVerbosity := verbosity.Value.String()
	previousLogToStderr := logToStderr.Value.String()
	output := &lockedBuffer{listenerExited: make(chan struct{})}

	if err := flag.Set("v", "6"); err != nil {
		t.Fatalf("set klog verbosity: %v", err)
	}
	if err := flag.Set("logtostderr", "false"); err != nil {
		t.Fatalf("disable klog stderr output: %v", err)
	}
	klog.SetOutput(output)
	t.Cleanup(func() {
		klog.Flush()
		klog.SetOutput(os.Stderr)
		_ = flag.Set("v", previousVerbosity)
		_ = flag.Set("logtostderr", previousLogToStderr)
	})

	actionErr := action()
	klog.Flush()
	return output, actionErr
}
