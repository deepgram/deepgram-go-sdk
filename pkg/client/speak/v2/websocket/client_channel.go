// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
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
	c.resetFinishStateIfClosed()
	return c.WSClient.ConnectWithCancel(ctx, ctxCancel, retryCnt)
}

// AttemptReconnect reconnects after exhausting retries.
func (c *WSChannel) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.AttemptReconnectWithCancel(c.ctx, c.ctxCancel, retries)
}

// AttemptReconnectWithCancel reconnects with a caller-supplied cancel function.
// The graceful-close milestones are reset so Finish waits on the new session.
func (c *WSChannel) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	c.resetFinishStateIfClosed()
	return c.WSClient.AttemptReconnectWithCancel(ctx, ctxCancel, retries)
}

// GetURL builds the Flux TTS WebSocket URL: wss://host/v2/speak?<SpeakV2WSOptions>
func (c *WSChannel) GetURL(host string) (string, error) {
	url, err := version.GetSpeakV2StreamAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.sOptions)
	if err != nil {
		klog.V(1).Infof("version.GetSpeakV2StreamAPI failed. Err: %v\n", err)
		return "", err
	}
	klog.V(5).Infof("Flux TTS connecting to %s\n", url)
	return url, nil
}

// Start launches the WebSocket-level ping goroutine if EnableKeepAlive is set.
// Start runs on every (re)connect; the guard ensures at most one ping goroutine
// is alive at a time.
func (c *WSChannel) Start() {
	if c.cOptions.EnableKeepAlive && c.pingActive.CompareAndSwap(false, true) {
		go c.ping(c.ctx)
	}
}

// ping sends WebSocket protocol-level ping frames on a fixed interval to keep
// the connection alive while ctx is active.
func (c *WSChannel) ping(ctx context.Context) {
	klog.V(6).Infof("fluxspeak.WSChannel.ping() ENTER\n")

	defer func() {
		c.pingActive.Store(false)
		if r := recover(); r != nil {
			klog.V(1).Infof("ping panic: %v\n%s\n", r, debug.Stack())
			sendErr := c.ProcessError(common.ErrFatalPanicRecovered)
			if sendErr != nil {
				klog.V(1).Infof("ping: ProcessError failed. Err: %v\n", sendErr)
			}
		}
		klog.V(6).Infof("fluxspeak.WSChannel.ping() LEAVE\n")
	}()

	period := c.cOptions.KeepAlivePeriod
	if period <= 0 {
		period = DefaultKeepAlivePeriod
	}
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			klog.V(3).Infof("fluxspeak.WSChannel.ping() Exiting\n")
			return
		case <-ticker.C:
			klog.V(5).Infof("Sending WebSocket ping...\n")
			if err := c.WritePing(); err != nil {
				klog.V(1).Infof("ping: WritePing failed. Err: %v\n", err)
			}
		}
	}
}

// ProcessMessage routes incoming WebSocket frames through the Flux TTS router.
// Text frames carry JSON server events; binary frames carry synthesized audio.
func (c *WSChannel) ProcessMessage(wsType int, byMsg []byte) error {
	klog.V(6).Infof("fluxspeak.WSChannel.ProcessMessage() ENTER\n")

	var err error
	switch wsType {
	case websocket.TextMessage:
		err = (*c.router).Message(byMsg)
	case websocket.BinaryMessage:
		err = (*c.router).Binary(byMsg)
	default:
		klog.V(7).Infof("ProcessMessage: frame ignored (type %d, len %d)\n", wsType, len(byMsg))
	}

	if err != nil {
		klog.V(1).Infof("ProcessMessage: router failed. Err: %v\n", err)
	}
	klog.V(6).Infof("fluxspeak.WSChannel.ProcessMessage() LEAVE\n")
	return err
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

// Speak sends text to be synthesized into the active turn.
// The server tracks the active turn internally and assigns a speech_id.
func (c *WSChannel) Speak(text string) error {
	klog.V(7).Infof("fluxspeak.WSChannel.Speak() ENTER\n")

	msg := msginterfaces.SpeakMessage{
		Type: MessageTypeSpeak,
		Text: text,
	}
	if err := c.WriteJSON(msg); err != nil {
		klog.V(1).Infof("Speak failed. Err: %v\n", err)
		klog.V(7).Infof("fluxspeak.WSChannel.Speak() LEAVE\n")
		return err
	}

	klog.V(4).Infof("Speak sent\n")
	klog.V(7).Infof("fluxspeak.WSChannel.Speak() LEAVE\n")
	return nil
}

// Flush ends the active turn. The server drains the buffer, generates the
// remaining audio, and reports the turn with a SpeechMetadata event.
func (c *WSChannel) Flush() error {
	klog.V(7).Infof("fluxspeak.WSChannel.Flush() ENTER\n")

	msg := msginterfaces.FlushMessage{
		Type: MessageTypeFlush,
	}
	if err := c.WriteJSON(msg); err != nil {
		klog.V(1).Infof("Flush failed. Err: %v\n", err)
		klog.V(7).Infof("fluxspeak.WSChannel.Flush() LEAVE\n")
		return err
	}

	klog.V(4).Infof("Flush sent\n")
	klog.V(7).Infof("fluxspeak.WSChannel.Flush() LEAVE\n")
	return nil
}

// Interrupt reports that the user barged in, without a playback offset.
// The server stops active audio generation and replies with SpeechInterrupted;
// without an offset the reply omits text_spoken and text_remaining.
func (c *WSChannel) Interrupt() error {
	return c.sendInterrupt(&msginterfaces.InterruptMessage{
		Type: MessageTypeInterrupt,
	})
}

// InterruptWithOffset reports that the user barged in after playing
// playedMs milliseconds of session audio (cumulative from the start of the
// session, not the current turn). With the offset the server can split the
// turn's text into text_spoken and text_remaining on the SpeechInterrupted reply.
func (c *WSChannel) InterruptWithOffset(playedMs int) error {
	return c.sendInterrupt(&msginterfaces.InterruptMessage{
		Type: MessageTypeInterrupt,
		PlaybackOffset: &msginterfaces.PlaybackOffset{
			Type:  msginterfaces.PlaybackOffsetTypeTimeMs,
			Value: playedMs,
		},
	})
}

func (c *WSChannel) sendInterrupt(msg *msginterfaces.InterruptMessage) error {
	klog.V(7).Infof("fluxspeak.WSChannel.Interrupt() ENTER\n")

	if err := c.WriteJSON(msg); err != nil {
		klog.V(1).Infof("Interrupt failed. Err: %v\n", err)
		klog.V(7).Infof("fluxspeak.WSChannel.Interrupt() LEAVE\n")
		return err
	}

	klog.V(4).Infof("Interrupt sent\n")
	klog.V(7).Infof("fluxspeak.WSChannel.Interrupt() LEAVE\n")
	return nil
}

// Configure sends a mid-session configuration update to the server.
// The server responds with ConfigureSuccess or ConfigureFailure.
//
// opts must be non-nil and name at least one field to change. Speed is validated
// client-side against the documented contract (0.5 to 1.5 in 0.05 increments), so
// an accidental zero value is rejected instead of being dropped from the wire
// message and accepted as an empty update.
func (c *WSChannel) Configure(opts *clientinterfaces.SpeakV2ConfigureOptions) error {
	klog.V(7).Infof("fluxspeak.WSChannel.Configure() ENTER\n")

	if opts == nil {
		klog.V(1).Infof("Configure: nil options\n")
		klog.V(7).Infof("fluxspeak.WSChannel.Configure() LEAVE\n")
		return interfacesv2.ErrOptionsRequired
	}
	if err := opts.Check(); err != nil {
		klog.V(1).Infof("Configure: invalid options. Err: %v\n", err)
		klog.V(7).Infof("fluxspeak.WSChannel.Configure() LEAVE\n")
		return err
	}

	msg := msginterfaces.ConfigureMessage{
		Type:  MessageTypeConfigure,
		Speed: opts.Speed,
	}
	if err := c.WriteJSON(msg); err != nil {
		klog.V(1).Infof("Configure failed. Err: %v\n", err)
		klog.V(7).Infof("fluxspeak.WSChannel.Configure() LEAVE\n")
		return err
	}

	klog.V(4).Infof("Configure sent\n")
	klog.V(7).Infof("fluxspeak.WSChannel.Configure() LEAVE\n")
	return nil
}

// GetCloseMsg returns the JSON bytes for the Close control message. The server
// drains all remaining audio, emits a final SessionMetadata, then closes the socket.
func (c *WSChannel) GetCloseMsg() []byte {
	return []byte(`{"type":"Close"}`)
}

// Finish gracefully closes the session. It sends the Close control message, then
// waits — bounded by ctx — while the server drains every queued turn: all remaining
// audio frames, the final SessionMetadata, and the server's own socket closure.
// When Finish returns nil, the handler has received every event the session
// produced. If the connection closes before SessionMetadata arrives,
// ErrGracefulCloseIncomplete is returned. On ctx expiry the connection is aborted
// and ctx.Err() is returned; audio still queued server-side is discarded.
//
// Use Stop for an immediate abort that does not wait for queued synthesis.
func (c *WSChannel) Finish(ctx context.Context) error {
	klog.V(7).Infof("fluxspeak.WSChannel.Finish() ENTER\n")

	fs := c.currentFinishState()

	if err := c.WriteJSON(msginterfaces.CloseMessage{Type: MessageTypeClose}); err != nil {
		klog.V(1).Infof("Finish: Close write failed. Err: %v\n", err)
		klog.V(7).Infof("fluxspeak.WSChannel.Finish() LEAVE\n")
		return err
	}

	select {
	case <-fs.sessionMetaDone:
	case <-fs.peerCloseDone:
		klog.V(1).Infof("Finish: connection closed before SessionMetadata\n")
		c.Stop()
		klog.V(7).Infof("fluxspeak.WSChannel.Finish() LEAVE\n")
		return ErrGracefulCloseIncomplete
	case <-ctx.Done():
		klog.V(1).Infof("Finish: context expired waiting for SessionMetadata\n")
		c.Stop()
		klog.V(7).Infof("fluxspeak.WSChannel.Finish() LEAVE\n")
		return ctx.Err()
	}

	select {
	case <-fs.peerCloseDone:
	case <-ctx.Done():
		klog.V(1).Infof("Finish: context expired waiting for server closure\n")
		c.Stop()
		klog.V(7).Infof("fluxspeak.WSChannel.Finish() LEAVE\n")
		return ctx.Err()
	}

	// the peer already closed the socket; Stop only releases local resources
	c.Stop()
	klog.V(4).Infof("Finish: graceful close complete\n")
	klog.V(7).Infof("fluxspeak.WSChannel.Finish() LEAVE\n")
	return nil
}

// errorToResponse converts a Go error into a typed ErrorResponse.
func (c *WSChannel) errorToResponse(err error) *msginterfaces.ErrorResponse {
	var errorCode, errorNum, errorDesc string
	matches := wsErrorRegexp.FindStringSubmatch(err.Error())
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
