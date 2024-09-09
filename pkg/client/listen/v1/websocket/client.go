// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the live/streaming client implementation for the Deepgram API
package websocketv1

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	live "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

type controlMessage struct {
	Type string `json:"type"`
}

/*
NewForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*Client, error) {
	return New(ctx, "", &interfaces.ClientOptions{}, options, nil)
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewWithDefaults(ctx context.Context, options *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	return New(ctx, "", &interfaces.ClientOptions{}, options, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription
*/
func New(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewWithCancel(ctx, ctxCancel, apiKey, cOptions, tOptions, callback)
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription
*/
func NewWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	klog.V(6).Infof("live.New() ENTER\n")

	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	err := cOptions.Parse()
	if err != nil {
		klog.V(1).Infof("ClientOptions.Parse() failed. Err: %v\n", err)
		return nil, err
	}
	err = tOptions.Check()
	if err != nil {
		klog.V(1).Infof("TranscribeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	if callback == nil {
		klog.V(2).Infof("Using DefaultCallbackHandler.\n")
		callback = live.NewDefaultCallbackHandler()
	}

	// init
	conn := Client{
		cOptions:  cOptions,
		tOptions:  tOptions,
		sendBuf:   make(chan []byte, 1),
		callback:  callback,
		router:    live.New(callback),
		ctx:       ctx,
		ctxCancel: ctxCancel,
		retry:     true,
	}

	klog.V(3).Infof("NewDeepGramWSClient Succeeded\n")
	klog.V(6).Infof("live.New() LEAVE\n")

	return &conn, nil
}

// Connect performs a websocket connection with "DefaultConnectRetry" number of retries.
func (c *Client) Connect() bool {
	// set the retry count
	if c.retryCnt == 0 {
		c.retryCnt = DefaultConnectRetry
	}
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt), true) != nil
}

// ConnectWithCancel performs a websocket connection with specified number of retries and providing a
// cancel function to stop the connection
func (c *Client) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	return c.internalConnectWithCancel(ctx, ctxCancel, retryCnt, true) != nil
}

// AttemptReconnect performs a reconnect after failing retries
func (c *Client) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.muConn.Lock()
	c.retry = true
	c.muConn.Unlock()
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(retries), true) != nil
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *Client) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.muConn.Lock()
	c.retry = true
	c.muConn.Unlock()
	return c.internalConnectWithCancel(ctx, ctxCancel, int(retries), true) != nil
}

func (c *Client) internalConnect() *websocket.Conn {
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt), false)
}

//nolint:funlen // this is a complex function. keep as is
func (c *Client) internalConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int, lock bool) *websocket.Conn {
	klog.V(7).Infof("live.internalConnectWithCancel() ENTER\n")

	// set the context
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	c.retryCnt = int64(retryCnt)

	// lock conn access
	if lock {
		klog.V(3).Infof("Locking connection mutex\n")
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(7).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")
		return nil
	}

	// if the connection is good, return it otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Connection is not valid\n")
			klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")
			return nil
		default:
			klog.V(7).Infof("Connection is good. Return object.")
			klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")
			return c.wsconn
		}
	} else {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Context is not valid. Has been canceled.\n")
			klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")
			return nil
		default:
			klog.V(3).Infof("Context is still valid. Retry...\n")
		}
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		/* #nosec G402 */
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.cOptions.SkipServerAuth},
		RedirectService: c.cOptions.RedirectService,
		SkipServerAuth:  c.cOptions.SkipServerAuth,
	}

	// set websocket headers
	myHeader := http.Header{}

	// restore application options to HTTP header
	if headers, ok := c.ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				klog.V(3).Infof("internalConnectWithCancel RESTORE Header: %s = %s\n", k, v)
				myHeader.Add(k, v)
			}
		}
	}

	// sets the API key
	myHeader.Set("Host", c.cOptions.Host)
	myHeader.Set("Authorization", "token "+c.cOptions.APIKey)
	myHeader.Set("User-Agent", interfaces.DgAgent)

	// attempt to establish connection
	i := int64(0)
	for {
		if i >= c.retryCnt {
			klog.V(3).Infof("Connect timeout... exiting!\n")
			c.retry = false
			break
		}

		// delay on subsequent calls
		if i > 0 {
			klog.V(2).Infof("Sleep for retry #%d...\n", i)
			time.Sleep(time.Second * time.Duration(defaultDelayBetweenRetry))
		}

		i++

		// create new connection
		url, err := version.GetLiveAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.tOptions)
		if err != nil {
			klog.V(1).Infof("version.GetLiveAPI failed. Err: %v\n", err)
			klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")
			return nil // no point in retrying because this is going to fail on every retry
		}
		klog.V(5).Infof("Connecting to %s\n", url)

		// perform the websocket connection
		ws, res, err := dialer.DialContext(c.ctx, url, myHeader)
		if res != nil {
			klog.V(3).Infof("HTTP Response: %s\n", res.Status)
			res.Body.Close()
		}
		if err != nil {
			klog.V(1).Infof("Cannot connect to websocket: %s\n", c.cOptions.Host)
			klog.V(1).Infof("Dialer failed. Err: %v\n", err)
			continue
		}

		// set the object to allow threads to function
		c.wsconn = ws
		c.retry = true

		// kick off threads to listen for messages and ping/keepalive
		go c.listen()
		if c.cOptions.EnableKeepAlive {
			go c.ping()
		}
		if c.cOptions.AutoFlushReplyDelta != 0 {
			go c.flush()
		}

		// fire off open connection
		err = c.router.OpenHelper(&msginterfaces.OpenResponse{
			Type: msginterfaces.TypeOpenResponse,
		})
		if err != nil {
			klog.V(1).Infof("router.OpenHelper failed. Err: %v\n", err)
		}

		klog.V(3).Infof("WebSocket Connection Successful!")
		klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")

		return c.wsconn
	}

	// if we get here, we failed to connect
	klog.V(1).Infof("Failed to connect to websocket: %s\n", c.cOptions.Host)
	klog.V(7).Infof("live.internalConnectWithCancel() LEAVE\n")

	return nil
}

//nolint:funlen,gocyclo // this is a complex function. keep as is
func (c *Client) listen() {
	klog.V(6).Infof("live.listen() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")

			// send error on callback
			err := ErrFatalPanicRecovered
			sendErr := c.sendError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			// fatal close
			c.closeWs(true, false)

			klog.V(6).Infof("live.listen() LEAVE\n")
			return
		}
	}()

	for {
		// doing a read, need to lock
		c.muConn.Lock()

		// get the connection
		ws := c.internalConnect()
		if ws == nil {
			// release
			c.muConn.Unlock()

			klog.V(3).Infof("listen: Connection is not valid\n")
			klog.V(6).Infof("live.listen() LEAVE\n")
			return
		}

		// release the lock
		c.muConn.Unlock()

		// msgType can be binary or text
		msgType, byMsg, err := ws.ReadMessage()

		if err != nil {
			errStr := err.Error()
			switch {
			case strings.Contains(errStr, SuccessfulSocketErr):
				klog.V(3).Infof("Graceful websocket close\n")

				// graceful close
				c.closeWs(false, false)

				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			case strings.Contains(errStr, UseOfClosedSocket):
				klog.V(3).Infof("Probable graceful websocket close: %v\n", err)

				// fatal close
				c.closeWs(false, false)

				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			case strings.Contains(errStr, FatalReadSocketErr):
				klog.V(1).Infof("Fatal socket error: %v\n", err)

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
				}

				// fatal close
				c.closeWs(true, false)

				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			case strings.Contains(errStr, "Deepgram"):
				klog.V(1).Infof("listen: Deepgram error. Err: %v\n", err)

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("listen: Deepgram ErrorMsg. Err: %v\n", sendErr)
				}

				// close the connection
				c.closeWs(false, false)

				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			case (err == io.EOF || err == io.ErrUnexpectedEOF):
				klog.V(3).Infof("stream object EOF\n")

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("listen: EOF error. Err: %v\n", sendErr)
				}

				// close the connection
				c.closeWs(true, false)

				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			default:
				klog.V(1).Infof("listen: Cannot read websocket message. Err: %v\n", err)

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("listen: EOF error. Err: %v\n", sendErr)
				}

				// close the connection
				c.closeWs(true, false)

				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			}
		}

		if len(byMsg) == 0 {
			klog.V(7).Infof("listen(): message empty")
			continue
		}

		// inspect the message
		if c.cOptions.InspectMessage() {
			err := c.inspect(byMsg)
			if err != nil {
				klog.V(1).Infof("listen: inspect failed. Err: %v\n", err)
			}
		}

		// callback
		if msgType == websocket.TextMessage {
			err := c.router.Message(byMsg)
			if err != nil {
				klog.V(1).Infof("live.listen(): router.Message failed. Err: %v\n", err)
			}
		} else {
			// this shouldn't happen, but let's log it
			klog.V(7).Infof("live.listen(): msg recv: type %d, len: %d\n", msgType, len(byMsg))
		}
	}
}

// Stream is a helper function to stream audio data from a io.Reader object to deepgram
func (c *Client) Stream(r io.Reader) error {
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
				case strings.Contains(errStr, SuccessfulSocketErr):
					klog.V(3).Infof("Graceful websocket close\n")
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return nil
				case strings.Contains(errStr, FatalReadSocketErr):
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

// WriteBinary writes binary data to the websocket server
func (c *Client) WriteBinary(byData []byte) error {
	klog.V(7).Infof("live.WriteBinary() ENTER\n")

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	// get the connection
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")

		return err
	}

	if err := ws.WriteMessage(
		websocket.BinaryMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteBinary WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteBinary Successful\n")
	klog.V(7).Infof("payload: %x\n", byData)
	klog.V(7).Infof("live.WriteBinary() LEAVE\n")

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the live transcription session on the Deepgram server.
*/
func (c *Client) WriteJSON(payload interface{}) error {
	klog.V(6).Infof("live.WriteJSON() ENTER\n")

	byData, err := json.Marshal(payload)
	if err != nil {
		klog.V(1).Infof("WriteJSON: Error marshaling JSON. Data: %v, Err: %v\n", payload, err)
		klog.V(6).Infof("live.WriteJSON() LEAVE\n")
		return err
	}

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(6).Infof("live.WriteJSON() LEAVE\n")

		return err
	}

	if err := ws.WriteMessage(
		websocket.TextMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteJSON WriteMessage failed. Err: %v\n", err)
		klog.V(6).Infof("live.WriteJSON() LEAVE\n")
		return err
	}

	klog.V(4).Infof("live.WriteJSON() Succeeded\n")
	klog.V(6).Infof("payload: %s\n", string(byData))
	klog.V(6).Infof("live.WriteJSON() LEAVE\n")

	return nil
}

/*
Write performs the lower level websocket write operation.
This is needed to implement the io.Writer interface. (aka the streaming interface)
*/
func (c *Client) Write(p []byte) (int, error) {
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
func (c *Client) KeepAlive() error {
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
func (c *Client) Finalize() error {
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

// closeStream sends an application level message to Deepgram
func (c *Client) closeStream(lock bool) error {
	klog.V(7).Infof("live.closeStream() ENTER\n")

	// doing a write, need to lock
	if lock {
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"CloseStream\" }"))
	if err != nil {
		klog.V(1).Infof("WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("live.closeStream() LEAVE\n")

		return err
	}

	klog.V(4).Infof("closeStream Succeeded\n")
	klog.V(7).Infof("live.closeStream() LEAVE\n")

	return err
}

// normalClosure sends a normal closure message to the server
func (c *Client) normalClosure(lock bool) error {
	klog.V(7).Infof("live.normalClosure() ENTER\n")

	// doing a write, need to lock
	if lock {
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.normalClosure() LEAVE\n")

		return err
	}

	err := c.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	switch err {
	case websocket.ErrCloseSent:
		klog.V(3).Infof("ErrCloseSent was sent. Err: %v\n", err)
	case nil:
		klog.V(4).Infof("normalClosure Succeeded\n")
	default:
		klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", err)
	}

	klog.V(7).Infof("live.normalClosure() LEAVE\n")

	return err
}

// Stop will send close message and shutdown websocket connection
func (c *Client) Stop() {
	klog.V(3).Infof("Stopping...\n")
	c.muConn.Lock()
	c.retry = false
	c.muConn.Unlock()

	// exit gracefully
	c.closeWs(false, true)
}

// closeWs closes the websocket connection
func (c *Client) closeWs(fatal bool, perm bool) {
	klog.V(6).Infof("live.closeWs() closing channels...\n")

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	if c.wsconn != nil && !fatal {
		// deepgram requires a close message to be sent
		_ = c.closeStream(false)
		time.Sleep(TerminationSleep) // allow time for server to register closure

		// websocket protocol message
		_ = c.normalClosure(false)
		time.Sleep(TerminationSleep) // allow time for server to register closure
	}

	// cancel the context because we are exiting exiting...
	if perm {
		c.ctxCancel()
	}

	if fatal || c.wsconn != nil {
		// fire off close connection
		err := c.router.CloseHelper(&msginterfaces.CloseResponse{
			Type: msginterfaces.TypeCloseResponse,
		})
		if err != nil {
			klog.V(1).Infof("router.CloseHelper failed. Err: %v\n", err)
		}
	}

	// close the connection
	if c.wsconn != nil {
		c.wsconn.Close()
		c.wsconn = nil
	}

	klog.V(4).Infof("live.closeWs() Succeeded\n")
	klog.V(6).Infof("live.closeWs() LEAVE\n")
}

// ping thread
func (c *Client) ping() {
	klog.V(6).Infof("live.ping() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")

			// send error on callback
			err := ErrFatalPanicRecovered
			sendErr := c.sendError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			// fatal close
			c.closeWs(true, false)

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

			// exit gracefully
			c.closeWs(false, false)

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
func (c *Client) flush() {
	klog.V(6).Infof("live.flush() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")

			// send error on callback
			err := ErrFatalPanicRecovered
			sendErr := c.sendError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			// fatal close
			c.closeWs(true, false)

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

			// exit gracefully
			c.closeWs(false, false)

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

// sendError sends an error message to the callback handler
func (c *Client) sendError(err error) error {
	response := c.errorToResponse(err)
	sendErr := c.router.ErrorHelper(response)
	if err != nil {
		klog.V(1).Infof("live.listen(): router.Error failed. Err: %v\n", sendErr)
	}

	return err
}

// errorToResponse converts an error into a Deepgram error response
func (c *Client) errorToResponse(err error) *msginterfaces.ErrorResponse {
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
		errorCode = UnknownDeepgramErr
		errorNum = UnknownDeepgramErr
		errorDesc = err.Error()
	}

	response := &msginterfaces.ErrorResponse{
		Type:        msginterfaces.TypeErrorResponse,
		ErrMsg:      strings.TrimSpace(fmt.Sprintf("%s %s", errorCode, errorNum)),
		Description: strings.TrimSpace(errorDesc),
		Variant:     errorNum,
	}
	return response
}

// inspectMessage inspects the message and determines the type to
// see if we should do anything with those types of messages
func (c *Client) inspect(byMsg []byte) error {
	klog.V(7).Infof("live.inspect() ENTER\n")

	var mt msginterfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(7).Infof("live.inspect() LEAVE\n")
		return err
	}

	switch mt.Type {
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

func (c *Client) inspectMessage(mr *msginterfaces.MessageResponse) error {
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
