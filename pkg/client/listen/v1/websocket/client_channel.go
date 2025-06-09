// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the live/streaming client implementation for the Deepgram API
package websocketv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
)

// Connect performs a websocket connection with "DefaultConnectRetry" number of retries.
func (c *WSChannel) Connect() bool {
	c.ctx, c.ctxCancel = context.WithCancel(c.ctx)
	return c.ConnectWithCancel(c.ctx, c.ctxCancel, int(DefaultConnectRetry))
}

// ConnectWithCancel performs a websocket connection with specified number of retries and providing a
// cancel function to stop the connection
func (c *WSChannel) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	return c.WSClient.ConnectWithCancel(ctx, ctxCancel, retryCnt)
}

// AttemptReconnect performs a reconnect after failing retries
func (c *WSChannel) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.AttemptReconnectWithCancel(c.ctx, c.ctxCancel, retries)
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *WSChannel) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	return c.WSClient.AttemptReconnectWithCancel(ctx, ctxCancel, retries)
}

// GetURL returns the websocket URL
func (c *WSChannel) GetURL(host string) (string, error) {
	url, err := version.GetLiveAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.tOptions)
	if err != nil {
		klog.V(1).Infof("version.GetLiveAPI failed. Err: %v\n", err)
		return "", err
	}
	klog.V(5).Infof("Connecting to %s\n", url)
	return url, nil
}

// Start the keepalive and flush threads
func (c *WSChannel) Start() {
	if c.cOptions.EnableKeepAlive {
		go c.ping()
	}
	if c.cOptions.AutoFlushReplyDelta != 0 {
		go c.flush()
	}
}

// ProcessMessage processes the message and sends it to the callback
func (c *WSChannel) ProcessMessage(wsType int, byMsg []byte) error {
	klog.V(6).Infof("ProcessMessage() ENTER\n")

	// inspect the message
	if c.cOptions.InspectListenMessage() {
		err := c.inspect(byMsg)
		if err != nil {
			klog.V(1).Infof("ProcessMessage: inspect failed. Err: %v\n", err)
		}
	}

	// callback
	if wsType == websocket.TextMessage {
		err := (*c.router).Message(byMsg)
		if err != nil {
			klog.V(1).Infof("ProcessMessage: router.Message failed. Err: %v\n", err)
			klog.V(6).Infof("ProcessMessage() LEAVE\n")

			return err
		}
	} else {
		// this shouldn't happen, but let's log it
		klog.V(7).Infof("ProcessMessage: msg recv: type %d, len: %d\n", wsType, len(byMsg))
	}

	klog.V(6).Infof("ProcessMessage Succeeded\n")
	klog.V(6).Infof("ProcessMessage() LEAVE\n")

	return nil
}

// Stream is a helper function to stream audio data from a io.Reader object to deepgram
func (c *WSChannel) Stream(r io.Reader) error {
	klog.V(6).Infof("live.Stream() ENTER\n")

	chunk := make([]byte, ChunkSize)

	for {
		select {
		case <-c.ctx.Done():
			klog.V(2).Infof("stream object Done()\n")
			klog.V(6).Infof("live.Stream() LEAVE\n")
			return nil
		default:
			bytesRead, err := r.Read(chunk)
			if err != nil {
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, common.SuccessfulSocketErr):
					klog.V(3).Infof("Graceful websocket close\n")
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return nil
				case strings.Contains(errStr, common.UseOfClosedSocket):
					klog.V(3).Infof("Graceful websocket close\n")
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return nil
				case strings.Contains(errStr, common.FatalReadSocketErr):
					klog.V(1).Infof("Fatal socket error: %v\n", err)
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return err
				case (err == io.EOF || err == io.ErrUnexpectedEOF):
					klog.V(3).Infof("stream object EOF\n")
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return err
				default:
					klog.V(1).Infof("r.Read error. Err: %v\n", err)
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return err
				}
			}

			if bytesRead == 0 {
				klog.V(7).Infof("Skipping. bytesRead == 0\n")
				continue
			}

			byteCount, err := c.Write(chunk[:bytesRead])
			if err != nil {
				klog.V(1).Infof("w.Write failed. Err: %v\n", err)
				klog.V(6).Infof("live.Stream() LEAVE\n")
				return err
			}
			klog.V(7).Infof("io.Writer succeeded. Bytes written: %d\n", byteCount)
		}
	}
}

/*
Write performs the lower level websocket write operation.
This is needed to implement the io.Writer interface. (aka the streaming interface)
*/
func (c *WSChannel) Write(p []byte) (int, error) {
	klog.V(7).Infof("live.Write() ENTER\n")

	byteLen := len(p)
	err := c.WriteBinary(p)
	if err != nil {
		klog.V(1).Infof("Write failed. Err: %v\n", err)
		klog.V(7).Infof("live.Write() LEAVE\n")
		return 0, err
	}

	klog.V(7).Infof("live.Write Succeeded\n")
	klog.V(7).Infof("live.Write() LEAVE\n")
	return byteLen, nil
}

/*
Kick off the keepalive message to the server
*/
func (c *WSChannel) KeepAlive() error {
	klog.V(7).Infof("live.KeepAlive() ENTER\n")

	err := c.WriteJSON(controlMessage{Type: MessageTypeKeepAlive})
	if err != nil {
		klog.V(1).Infof("KeepAlive failed. Err: %v\n", err)
		klog.V(7).Infof("live.KeepAlive() LEAVE\n")

		return err
	}

	klog.V(4).Infof("KeepAlive Succeeded\n")
	klog.V(7).Infof("live.KeepAlive() LEAVE\n")

	return err
}

/*
Finalize the live transcription utterance/sentence/fragment
*/
func (c *WSChannel) Finalize() error {
	klog.V(7).Infof("live.KeepAlive() ENTER\n")

	err := c.WriteJSON(controlMessage{Type: MessageTypeFinalize})
	if err != nil {
		klog.V(1).Infof("Finalize failed. Err: %v\n", err)
		klog.V(7).Infof("live.Finalize() LEAVE\n")

		return err
	}

	klog.V(4).Infof("Finalize Succeeded\n")
	klog.V(7).Infof("live.Finalize() LEAVE\n")

	return err
}

// GetCloseMsg sends an application level message to Deepgram
func (c *WSChannel) GetCloseMsg() []byte {
	return []byte("{ \"type\": \"CloseStream\" }")
}

// Finish the websocket connection
func (c *WSChannel) Finish() {
	// NA
}

// ProcessError processes the error and sends it to the callback
func (c *WSChannel) ProcessError(err error) error {
	response := c.errorToResponse(err)
	sendErr := (*c.router).Error(response)
	if err != nil {
		klog.V(1).Infof("ProcessError failed. Err: %v\n", sendErr)
	}

	return err
}

// ping thread
func (c *WSChannel) ping() {
	klog.V(6).Infof("live.ping() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")
			klog.V(1).Infof("Panic: %v\n", r)
			klog.V(1).Infof("Stack trace: %s\n", string(debug.Stack()))

			// send error on callback
			err := common.ErrFatalPanicRecovered
			sendErr := c.ProcessError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			klog.V(6).Infof("live.ping() LEAVE\n")
			return
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("live.ping() Exiting\n")
			klog.V(6).Infof("live.ping() LEAVE\n")
			return
		case <-ticker.C:
			klog.V(5).Infof("Starting ping...")

			// deepgram keepalive message
			klog.V(5).Infof("Sending Deepgram KeepAlive message...\n")
			err := c.KeepAlive()
			if err == nil {
				klog.V(5).Infof("Ping sent!")
			} else {
				klog.V(1).Infof("Failed to send Deepgram KeepAlive. Err: %v\n", err)
			}
		}
	}
}

// flush thread
func (c *WSChannel) flush() {
	klog.V(6).Infof("live.flush() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")
			klog.V(1).Infof("Panic: %v\n", r)
			klog.V(1).Infof("Stack trace: %s\n", string(debug.Stack()))

			// send error on callback
			err := common.ErrFatalPanicRecovered
			sendErr := c.ProcessError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			klog.V(6).Infof("live.flush() LEAVE\n")
			return
		}
	}()

	ticker := time.NewTicker(flushPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("live.flush() Exiting\n")
			klog.V(6).Infof("live.flush() LEAVE\n")
			return
		case <-ticker.C:
			// doing a read, need to lock.
			c.muFinal.Lock()

			// have we received anything? no, then skip
			if c.lastDatagram == nil {
				klog.V(7).Infof("No datagram received. Skipping...\n")
				c.muFinal.Unlock()
				continue
			}

			// we have received something, but is it recent?
			trigger := c.lastDatagram.Add(time.Millisecond * time.Duration(c.cOptions.AutoFlushReplyDelta))
			now := time.Now()
			klog.V(7).Infof("Time (Last): %s\n", trigger.String())
			klog.V(7).Infof("Time (Now ): %s\n", now.String())
			bNeedFlush := trigger.Before(now)
			if bNeedFlush {
				c.lastDatagram = nil
			}

			// release
			c.muFinal.Unlock()

			if bNeedFlush {
				klog.V(5).Infof("Sending Finalize message...\n")
				err := c.Finalize()
				if err == nil {
					klog.V(5).Infof("Finalize sent!")
				} else {
					klog.V(1).Infof("Failed to send Finalize. Err: %v\n", err)
				}
			}
		}
	}
}

// errorToResponse converts an error into a Deepgram error response
func (c *WSChannel) errorToResponse(err error) *msginterfaces.ErrorResponse {
	r := regexp.MustCompile(`websocket: ([a-z]+) (\d+) .+: (.+)`)

	var errorCode string
	var errorNum string
	var errorDesc string

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

	response := &msginterfaces.ErrorResponse{
		Type:        string(msginterfaces.TypeErrorResponse),
		ErrMsg:      strings.TrimSpace(fmt.Sprintf("%s %s", errorCode, errorNum)),
		Description: strings.TrimSpace(errorDesc),
		Variant:     errorNum,
	}
	return response
}

// inspectMessage inspects the message and determines the type to
// see if we should do anything with those types of messages
func (c *WSChannel) inspect(byMsg []byte) error {
	klog.V(7).Infof("live.inspect() ENTER\n")

	var mt msginterfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(7).Infof("live.inspect() LEAVE\n")
		return err
	}

	switch msginterfaces.TypeResponse(mt.Type) {
	case msginterfaces.TypeMessageResponse:
		klog.V(7).Infof("TypeMessageResponse\n")

		// convert to MessageResponse
		var mr msginterfaces.MessageResponse
		if err := json.Unmarshal(byMsg, &mr); err != nil {
			klog.V(1).Infof("json.Unmarshal(MessageResponse) failed. Err: %v\n", err)
			klog.V(7).Infof("live.inspect() LEAVE\n")
			return err
		}

		// inspect the message
		err := c.inspectMessage(&mr)
		if err != nil {
			klog.V(1).Infof("inspectMessage() failed. Err: %v\n", err)
			klog.V(7).Infof("live.inspect() LEAVE\n")
			return err
		}
	default:
		klog.V(7).Infof("MessageType: %s\n", mt.Type)
	}

	klog.V(7).Info("inspect() succeeded\n")
	klog.V(7).Infof("live.inspect() LEAVE\n")
	return nil
}

func (c *WSChannel) inspectMessage(mr *msginterfaces.MessageResponse) error {
	klog.V(7).Infof("live.inspectMessage() ENTER\n")

	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)
	if len(mr.Channel.Alternatives) == 0 || sentence == "" {
		klog.V(7).Info("inspectMessage is empty\n")
		klog.V(7).Infof("live.inspectMessage() LEAVE\n")
		return nil
	}

	if mr.IsFinal {
		klog.V(7).Infof("IsFinal received: %s\n", time.Now().String())

		// doing a write, need to lock
		c.muFinal.Lock()
		c.lastDatagram = nil
		c.muFinal.Unlock()
	} else {
		klog.V(7).Infof("Interim received: %s\n", time.Now().String())

		// last datagram received
		c.muFinal.Lock()
		now := time.Now()
		c.lastDatagram = &now
		c.muFinal.Unlock()
	}

	klog.V(7).Info("inspectMessage() succeeded\n")
	klog.V(7).Infof("live.inspectMessage() LEAVE\n")
	return nil
}
