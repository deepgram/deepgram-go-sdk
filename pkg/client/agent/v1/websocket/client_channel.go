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
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
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
	// we dont send the SettingsConfigurationOptions because that is sent as a WS message to the server
	url, err := version.GetAgentAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path /*, c.tOptions*/)
	if err != nil {
		klog.V(1).Infof("version.GetAgentAPI failed. Err: %v\n", err)
		return "", err
	}
	klog.V(5).Infof("Connecting to %s\n", url)
	return url, nil
}

// Start the keepalive and flush threads
func (c *WSChannel) Start() {
	// send ConfigurationOptions to server
	if c.tOptions != nil {
		// send the configuration settings to the server
		klog.V(4).Infof("Sending ConfigurationSettings to server\n")
		err := c.WriteJSON(c.tOptions)
		if err != nil {
			klog.V(1).Infof("w.WriteJSON ConfigurationSettings failed. Err: %v\n", err)

			// terminate the connection
			c.WSClient.Stop()

			return
		}
	}

	if c.cOptions.EnableKeepAlive {
		go c.ping()
	}
}

// ProcessMessage processes the message and sends it to the callback
func (c *WSChannel) ProcessMessage(wsType int, byMsg []byte) error {
	klog.V(6).Infof("ProcessMessage() ENTER\n")

	switch wsType {
	case websocket.TextMessage:
		// route the message
		err := (*c.router).Message(byMsg)
		if err != nil {
			klog.V(1).Infof("agent.listen(): router.Message failed. Err: %v\n", err)
		}
	case websocket.BinaryMessage:
		// audio data!
		err := (*c.router).Binary(byMsg)
		if err != nil {
			klog.V(1).Infof("agent.listen(): router.Binary failed. Err: %v\n", err)
		}
	default:
		klog.V(7).Infof("agent.listen(): msg recv: type %d, len: %d\n", wsType, len(byMsg))
	}

	klog.V(6).Infof("ProcessMessage Succeeded\n")
	klog.V(6).Infof("ProcessMessage() LEAVE\n")

	return nil
}

// Stream is a helper function to stream audio data from a io.Reader object to deepgram
func (c *WSChannel) Stream(r io.Reader) error {
	klog.V(6).Infof("agent.Stream() ENTER\n")

	chunk := make([]byte, ChunkSize)

	for {
		select {
		case <-c.ctx.Done():
			klog.V(2).Infof("stream object Done()\n")
			klog.V(6).Infof("agent.Stream() LEAVE\n")
			return nil
		default:
			bytesRead, err := r.Read(chunk)
			if err != nil {
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, common.SuccessfulSocketErr):
					klog.V(3).Infof("Graceful websocket close\n")
					klog.V(6).Infof("agent.Stream() LEAVE\n")
					return nil
				case strings.Contains(errStr, common.UseOfClosedSocket):
					klog.V(3).Infof("Graceful websocket close\n")
					klog.V(6).Infof("agent.Stream() LEAVE\n")
					return nil
				case strings.Contains(errStr, common.FatalReadSocketErr):
					klog.V(1).Infof("Fatal socket error: %v\n", err)
					klog.V(6).Infof("agent.Stream() LEAVE\n")
					return err
				case (err == io.EOF || err == io.ErrUnexpectedEOF):
					klog.V(3).Infof("stream object EOF\n")
					klog.V(6).Infof("agent.Stream() LEAVE\n")
					return err
				default:
					klog.V(1).Infof("r.Read error. Err: %v\n", err)
					klog.V(6).Infof("agent.Stream() LEAVE\n")
					return err
				}
			}

			if bytesRead == 0 {
				klog.V(7).Infof("Skipping. bytesRead == 0\n")
				continue
			}

			err = c.WriteBinary(chunk[:bytesRead])
			if err != nil {
				klog.V(1).Infof("w.Write failed. Err: %v\n", err)
				klog.V(6).Infof("agent.Stream() LEAVE\n")
				return err
			}
			klog.V(7).Infof("io.Writer succeeded\n")
		}
	}
}

/*
Write performs the lower level websocket write operation.
This is needed to implement the io.Writer interface. (aka the streaming interface)
*/
func (c *WSChannel) Write(p []byte) (int, error) {
	klog.V(7).Infof("agent.Write() ENTER\n")

	byteLen := len(p)
	err := c.WriteBinary(p)
	if err != nil {
		klog.V(1).Infof("Write failed. Err: %v\n", err)
		klog.V(7).Infof("agent.Write() LEAVE\n")
		return 0, err
	}

	klog.V(7).Infof("agent.Write Succeeded\n")
	klog.V(7).Infof("agent.Write() LEAVE\n")
	return byteLen, nil
}

/*
Kick off the keepalive message to the server
*/
func (c *WSChannel) KeepAlive() error {
	klog.V(7).Infof("agent.KeepAlive() ENTER\n")

	keepAlive := msginterfaces.KeepAlive{
		Type: msginterfaces.TypeKeepAlive,
	}
	err := c.WriteJSON(keepAlive)
	if err != nil {
		klog.V(1).Infof("KeepAlive failed. Err: %v\n", err)
		klog.V(7).Infof("agent.KeepAlive() LEAVE\n")

		return err
	}

	klog.V(4).Infof("KeepAlive Succeeded\n")
	klog.V(7).Infof("agent.KeepAlive() LEAVE\n")

	return err
}

// GetCloseMsg sends an application level message to Deepgram
func (c *WSChannel) GetCloseMsg() []byte {
	close := msginterfaces.Close{
		Type: msginterfaces.TypeClose,
	}

	byMsg, err := json.Marshal(close)
	if err != nil {
		klog.V(1).Infof("GetCloseMsg failed. Err: %v\n", err)
		return nil
	}

	return byMsg
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
	klog.V(6).Infof("agent.ping() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")

			// send error on callback
			err := common.ErrFatalPanicRecovered
			sendErr := c.ProcessError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			klog.V(6).Infof("agent.ping() LEAVE\n")
			return
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("agent.ping() Exiting\n")
			klog.V(6).Infof("agent.ping() LEAVE\n")
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
