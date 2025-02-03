// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the speak/streaming client implementation for the Deepgram API
package websocketv1

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
)

// Connect performs a websocket connection with "DefaultConnectRetry" number of retries.
func (c *WSCallback) Connect() bool {
	c.ctx, c.ctxCancel = context.WithCancel(c.ctx)
	return c.ConnectWithCancel(c.ctx, c.ctxCancel, int(DefaultConnectRetry))
}

// ConnectWithCancel performs a websocket connection with specified number of retries and providing a
// cancel function to stop the connection
func (c *WSCallback) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	return c.WSClient.ConnectWithCancel(ctx, ctxCancel, retryCnt)
}

// AttemptReconnect performs a reconnect after failing retries
func (c *WSCallback) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.AttemptReconnectWithCancel(c.ctx, c.ctxCancel, retries)
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *WSCallback) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	return c.WSClient.AttemptReconnectWithCancel(ctx, ctxCancel, retries)
}

// GetURL returns the websocket URL
func (c *WSCallback) GetURL(host string) (string, error) {
	url, err := version.GetSpeakStreamAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.sOptions)
	if err != nil {
		klog.V(1).Infof("version.GetSpeakStreamAPI failed. Err: %v\n", err)
		return "", err
	}
	klog.V(5).Infof("Connecting to %s\n", url)
	return url, nil
}

// Start the callback
func (c *WSCallback) Start() {
	if c.cOptions.AutoFlushSpeakDelta != 0 {
		go c.flush()
	}
}

// ProcessMessage processes the incoming message
func (c *WSCallback) ProcessMessage(wsType int, byMsg []byte) error {
	klog.V(6).Infof("ProcessMessage() ENTER\n")

	switch wsType {
	case websocket.TextMessage:
		// inspect the message
		if c.cOptions.InspectSpeakMessage() {
			err := c.inspect(byMsg)
			if err != nil {
				klog.V(1).Infof("speak: inspect failed. Err: %v\n", err)
			}
		}

		// route the message
		err := (*c.router).Message(byMsg)
		if err != nil {
			klog.V(1).Infof("speak.listen(): router.Message failed. Err: %v\n", err)
		}
	case websocket.BinaryMessage:
		// audio data!
		err := (*c.router).Binary(byMsg)
		if err != nil {
			klog.V(1).Infof("speak.listen(): router.Binary failed. Err: %v\n", err)
		}
	default:
		klog.V(7).Infof("speak.listen(): msg recv: type %d, len: %d\n", wsType, len(byMsg))
	}

	klog.V(6).Infof("ProcessMessage Succeeded\n")
	klog.V(6).Infof("ProcessMessage() LEAVE\n")

	return nil
}

// SpeakWithText writes text to the websocket server to obtain corresponding audio
//
// This function will automatically wrap the text in the appropriate JSON structure
// and send it to the server
//
// Args:
//
//	text: string containing the text to be spoken
//
// Return:
//
//	error: if successful, returns nil otherwise an error object
func (c *WSCallback) SpeakWithText(text string) error {
	klog.V(6).Infof("speak.SpeakText() ENTER\n")
	klog.V(4).Infof("text: %s\n", text)

	err := c.WSClient.WriteJSON(TextSource{
		Type: MessageTypeSpeak,
		Text: text,
	})
	if err == nil {
		klog.V(4).Infof("SpeakText Succeeded\n")
	} else {
		klog.V(1).Infof("SpeakText failed. Err: %v\n", err)
	}

	klog.V(6).Infof("speak.SpeakText() LEAVE\n")

	return err
}

// Speak is an alias function for SpeakWithText
func (c *WSCallback) Speak(text string) error {
	return c.SpeakWithText(text)
}

// WriteJSON writes a JSON message to the websocket
func (c *WSCallback) WriteJSON(playload controlMessage) error {
	if playload.Type == MessageTypeFlush {
		c.muFinal.Lock()
		c.flushCount++
		klog.V(5).Infof("Flush Count: %d\n", c.flushCount)
		c.muFinal.Unlock()
	}

	return c.WSClient.WriteJSON(playload)
}

// Flush will instruct the server to flush the current text buffer
func (c *WSCallback) Flush() error {
	klog.V(6).Infof("speak.Flush() ENTER\n")

	err := c.WriteJSON(controlMessage{Type: MessageTypeFlush})
	if err != nil {
		klog.V(1).Infof("Flush failed. Err: %v\n", err)
		klog.V(6).Infof("speak.Flush() LEAVE\n")

		return err
	}
	c.flushCount++

	klog.V(4).Infof("Flush Succeeded\n")
	klog.V(6).Infof("speak.Flush() LEAVE\n")

	return err
}

// Reset will instruct the server to reset the current buffer
func (c *WSCallback) Reset() error {
	klog.V(6).Infof("speak.Reset() ENTER\n")

	err := c.WriteJSON(controlMessage{Type: MessageTypeReset})
	if err != nil {
		klog.V(1).Infof("Reset failed. Err: %v\n", err)
		klog.V(6).Infof("speak.Reset() LEAVE\n")

		return err
	}

	klog.V(4).Infof("Reset Succeeded\n")
	klog.V(6).Infof("speak.Reset() LEAVE\n")
	return nil
}

// GetCloseMsg sends an application level message to Deepgram
func (c *WSCallback) GetCloseMsg() []byte {
	return []byte("{ \"type\": \"Close\" }")
}

// Finish the callback
func (c *WSCallback) Finish() {
	// NA
}

// ProcessError sends an error message to the callback handler
func (c *WSCallback) ProcessError(err error) error {
	response := c.errorToResponse(err)
	sendErr := (*c.router).Error(response)
	if err != nil {
		klog.V(1).Infof("speak.listen(): router.Error failed. Err: %v\n", sendErr)
	}

	return err
}

// flush thread
func (c *WSCallback) flush() {
	klog.V(6).Infof("speak.flush() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")
			klog.V(1).Infof("Panic: %v\n", r)
			klog.V(1).Infof("Stack trace: %s\n", string(debug.Stack()))

			// send error on callback
			err := common.ErrFatalPanicRecovered
			sendErr := c.ProcessError(err)
			if sendErr != nil {
				klog.V(1).Infof("speak: Fatal socket error. Err: %v\n", sendErr)
			}

			klog.V(6).Infof("speak.flush() LEAVE\n")
			return
		}
	}()

	ticker := time.NewTicker(flushPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("speak.flush() Exiting\n")
			klog.V(6).Infof("speak.flush() LEAVE\n")
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
			trigger := c.lastDatagram.Add(time.Millisecond * time.Duration(c.cOptions.AutoFlushSpeakDelta))
			now := time.Now()
			klog.V(6).Infof("Time (Last): %s\n", trigger.String())
			klog.V(6).Infof("Time (Now ): %s\n", now.String())
			bNeedFlush := trigger.Before(now)
			if bNeedFlush {
				c.lastDatagram = nil
			}

			// release
			c.muFinal.Unlock()

			if bNeedFlush {
				klog.V(5).Infof("Sending Flush message...\n")
				err := c.Flush()
				if err == nil {
					klog.V(5).Infof("Flush sent!")
				} else {
					klog.V(1).Infof("Failed to send Flush. Err: %v\n", err)
				}
			}
		}
	}
}

// errorToResponse converts an error into a Deepgram error response
func (c *WSCallback) errorToResponse(err error) *msginterfaces.ErrorResponse {
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

// inspect will check the message and determine the type to
// see if we should do  actionable based on those types of messages
func (c *WSCallback) inspect(byMsg []byte) error {
	klog.V(7).Infof("speak.inspect() ENTER\n")

	var mt msginterfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(7).Infof("speak.inspect() LEAVE\n")
		return err
	}

	switch msginterfaces.TypeResponse(mt.Type) {
	case msginterfaces.TypeFlushedResponse:
		klog.V(7).Infof("TypeFlushedResponse\n")

		// decrement the flush count
		c.muFinal.Lock()
		c.flushCount--
		klog.V(5).Infof("Flush Count: %d\n", c.flushCount)
		c.muFinal.Unlock()
	default:
		klog.V(5).Infof("MessageType: %s\n", mt.Type)
	}

	klog.V(7).Info("inspect() succeeded\n")
	klog.V(7).Infof("speak.inspect() LEAVE\n")
	return nil
}
