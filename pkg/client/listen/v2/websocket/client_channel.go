// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// Connect performs a WebSocket connection with DefaultConnectRetry retries.
func (c *WSChannel) Connect() bool {
	c.ctx, c.ctxCancel = context.WithCancel(c.ctx)
	return c.ConnectWithCancel(c.ctx, c.ctxCancel, int(DefaultConnectRetry))
}

// ConnectWithCancel performs a WebSocket connection with a caller-supplied context.
func (c *WSChannel) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	return c.WSClient.ConnectWithCancel(ctx, ctxCancel, retryCnt)
}

// AttemptReconnect reconnects after exhausting retries.
func (c *WSChannel) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.AttemptReconnectWithCancel(c.ctx, c.ctxCancel, retries)
}

// AttemptReconnectWithCancel reconnects with a caller-supplied cancel function.
func (c *WSChannel) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	return c.WSClient.AttemptReconnectWithCancel(ctx, ctxCancel, retries)
}

// GetURL builds the Flux WebSocket URL: wss://host/v2/listen?<FluxTranscriptionOptions>
func (c *WSChannel) GetURL(host string) (string, error) {
	url, err := version.GetFluxAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.tOptions)
	if err != nil {
		klog.V(1).Infof("version.GetFluxAPI failed. Err: %v\n", err)
		return "", err
	}
	klog.V(5).Infof("Flux connecting to %s\n", url)
	return url, nil
}

// Start launches the WebSocket-level ping goroutine if EnableKeepAlive is set.
func (c *WSChannel) Start() {
	if c.cOptions.EnableKeepAlive {
		go c.ping()
	}
}

// ping sends WebSocket protocol-level ping frames on a fixed interval to keep
// the connection alive while the context is active.
func (c *WSChannel) ping() {
	klog.V(6).Infof("flux.WSChannel.ping() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("ping panic: %v\n%s\n", r, debug.Stack())
			sendErr := c.ProcessError(common.ErrFatalPanicRecovered)
			if sendErr != nil {
				klog.V(1).Infof("ping: ProcessError failed. Err: %v\n", sendErr)
			}
		}
		klog.V(6).Infof("flux.WSChannel.ping() LEAVE\n")
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("flux.WSChannel.ping() Exiting\n")
			return
		case <-ticker.C:
			klog.V(5).Infof("Sending WebSocket ping...\n")
			if err := c.WritePing(); err != nil {
				klog.V(1).Infof("ping: WritePing failed. Err: %v\n", err)
			}
		}
	}
}

// ProcessMessage routes incoming WebSocket text frames through the Flux router.
func (c *WSChannel) ProcessMessage(wsType int, byMsg []byte) error {
	klog.V(6).Infof("flux.WSChannel.ProcessMessage() ENTER\n")

	if wsType == websocket.TextMessage {
		err := (*c.router).Message(byMsg)
		if err != nil {
			klog.V(1).Infof("ProcessMessage: router.Message failed. Err: %v\n", err)
			klog.V(6).Infof("flux.WSChannel.ProcessMessage() LEAVE\n")
			return err
		}
	} else {
		klog.V(7).Infof("ProcessMessage: binary frame ignored (type %d, len %d)\n", wsType, len(byMsg))
	}

	klog.V(6).Infof("flux.WSChannel.ProcessMessage() LEAVE\n")
	return nil
}

// ProcessError converts a transport error into an ErrorResponse and routes it.
func (c *WSChannel) ProcessError(err error) error {
	response := c.errorToResponse(err)
	sendErr := (*c.router).Error(response)
	if sendErr != nil {
		klog.V(1).Infof("ProcessError: router.Error failed. Err: %v\n", sendErr)
	}
	return err
}

// Write sends raw audio bytes as a binary WebSocket frame.
// Implements io.Writer so the client can be used directly with streaming helpers.
func (c *WSChannel) Write(p []byte) (int, error) {
	klog.V(7).Infof("flux.WSChannel.Write() ENTER\n")

	byteLen := len(p)
	if err := c.WriteBinary(p); err != nil {
		klog.V(1).Infof("Write failed. Err: %v\n", err)
		klog.V(7).Infof("flux.WSChannel.Write() LEAVE\n")
		return 0, err
	}

	klog.V(7).Infof("flux.WSChannel.Write() LEAVE\n")
	return byteLen, nil
}

// Stream reads audio from r in ChunkSize chunks and writes each chunk to the WebSocket.
// Blocks until r is exhausted or the context is cancelled.
func (c *WSChannel) Stream(r io.Reader) error {
	klog.V(6).Infof("flux.WSChannel.Stream() ENTER\n")

	chunk := make([]byte, ChunkSize)
	for {
		select {
		case <-c.ctx.Done():
			klog.V(2).Infof("Stream context Done()\n")
			klog.V(6).Infof("flux.WSChannel.Stream() LEAVE\n")
			return nil
		default:
			bytesRead, err := r.Read(chunk)
			if err != nil {
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, common.SuccessfulSocketErr),
					strings.Contains(errStr, common.UseOfClosedSocket):
					klog.V(3).Infof("Graceful stream close\n")
					klog.V(6).Infof("flux.WSChannel.Stream() LEAVE\n")
					return nil
				case strings.Contains(errStr, common.FatalReadSocketErr):
					klog.V(1).Infof("Fatal socket error: %v\n", err)
					klog.V(6).Infof("flux.WSChannel.Stream() LEAVE\n")
					return err
				case err == io.EOF || err == io.ErrUnexpectedEOF:
					klog.V(3).Infof("Stream EOF\n")
					klog.V(6).Infof("flux.WSChannel.Stream() LEAVE\n")
					return err
				default:
					klog.V(1).Infof("r.Read error. Err: %v\n", err)
					klog.V(6).Infof("flux.WSChannel.Stream() LEAVE\n")
					return err
				}
			}

			if bytesRead == 0 {
				klog.V(7).Infof("Skipping empty read\n")
				continue
			}

			if _, err := c.Write(chunk[:bytesRead]); err != nil {
				klog.V(1).Infof("Write failed. Err: %v\n", err)
				klog.V(6).Infof("flux.WSChannel.Stream() LEAVE\n")
				return err
			}
		}
	}
}

// Configure sends a mid-session configuration update to the server.
// The server responds with ListenV2ConfigureSuccess or ListenV2ConfigureFailure.
func (c *WSChannel) Configure(opts *clientinterfaces.FluxConfigureOptions) error {
	klog.V(7).Infof("flux.WSChannel.Configure() ENTER\n")

	msg := msginterfaces.ConfigureMessage{
		Type:          MessageTypeConfigure,
		Thresholds:    opts.Thresholds,
		Keyterms:      opts.Keyterms,
		LanguageHints: opts.LanguageHints,
	}
	if err := c.WriteJSON(msg); err != nil {
		klog.V(1).Infof("Configure failed. Err: %v\n", err)
		klog.V(7).Infof("flux.WSChannel.Configure() LEAVE\n")
		return err
	}

	klog.V(4).Infof("Configure sent\n")
	klog.V(7).Infof("flux.WSChannel.Configure() LEAVE\n")
	return nil
}

// GetCloseMsg returns the JSON bytes for the CloseStream control message.
func (c *WSChannel) GetCloseMsg() []byte {
	return []byte(`{"type":"CloseStream"}`)
}

// Finish is a no-op. Flux uses server-side turn detection; no client flush is needed.
func (c *WSChannel) Finish() {}

// errorToResponse converts a Go error into a typed ErrorResponse.
func (c *WSChannel) errorToResponse(err error) *msginterfaces.ErrorResponse {
	r := regexp.MustCompile(`websocket: ([a-z]+) (\d+) .+: (.+)`)

	var errorCode, errorNum, errorDesc string
	matches := r.FindStringSubmatch(err.Error())
	if len(matches) > 3 {
		errorCode = matches[1]
		errorNum = matches[2]
		errorDesc = matches[3]
	} else {
		errorCode = common.UnknownDeepgramErr
		errorNum = common.UnknownDeepgramErr
		errorDesc = err.Error()
	}

	return &msginterfaces.ErrorResponse{
		Type:        string(msginterfaces.TypeErrorResponse),
		ErrMsg:      strings.TrimSpace(fmt.Sprintf("%s %s", errorCode, errorNum)),
		Description: strings.TrimSpace(errorDesc),
		Variant:     errorNum,
	}
}
