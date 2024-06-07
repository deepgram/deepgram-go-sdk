// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the live/streaming client implementation for the Deepgram API
*/
package live

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

	live "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

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
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt)) != nil
}

// ConnectWithCancel performs a websocket connection with specified number of retries and providing a
// cancel function to stop the connection
func (c *Client) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	return c.internalConnectWithCancel(ctx, ctxCancel, retryCnt) != nil
}

// AttemptReconnect performs a reconnect after failing retries
func (c *Client) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.retry = true
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(retries)) != nil
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *Client) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.retry = true
	return c.internalConnectWithCancel(ctx, ctxCancel, int(retries)) != nil
}

func (c *Client) internalConnect() *websocket.Conn {
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt))
}

//nolint:funlen // this is a complex function. keep as is
func (c *Client) internalConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) *websocket.Conn {
	klog.V(7).Infof("live.Connect() ENTER\n")

	// set the context
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	c.retryCnt = int64(retryCnt)

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(7).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		klog.V(7).Infof("live.Connect() LEAVE\n")
		return nil
	}

	// if the connection is good, return it otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Connection is not valid\n")
			klog.V(7).Infof("live.ConnectWithRetry() LEAVE\n")
			return nil
		default:
			klog.V(7).Infof("Connection is good. Return object.")
			klog.V(7).Infof("live.ConnectWithRetry() LEAVE\n")
			return c.wsconn
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
				klog.V(3).Infof("Connect() RESTORE Header: %s = %s\n", k, v)
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
			klog.V(7).Infof("live.ConnectWithRetry() LEAVE\n")
			return nil // no point in retrying because this is going to fail on every retry
		}
		klog.V(5).Infof("Connecting to %s\n", url)

		// a single connection attempt
		// Note: not using defer here because we arent leaving the scope of the function
		c.muConn.Lock()

		// perform the websocket connection
		ws, res, err := dialer.DialContext(c.ctx, url, myHeader)
		if res != nil {
			klog.V(3).Infof("HTTP Response: %s\n", res.Status)
			res.Body.Close()
		}
		if err != nil {
			klog.V(1).Infof("Cannot connect to websocket: %s\n", c.cOptions.Host)
			klog.V(1).Infof("Dialer failed. Err: %v\n", err)
			c.muConn.Unlock()
			continue
		}

		// set the object to allow threads to function
		c.wsconn = ws
		c.retry = true

		// unlock the connection
		c.muConn.Unlock()

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
		klog.V(7).Infof("live.ConnectWithRetry() LEAVE\n")

		return c.wsconn
	}

	// if we get here, we failed to connect
	klog.V(1).Infof("Failed to connect to websocket: %s\n", c.cOptions.Host)
	klog.V(7).Infof("live.ConnectWithRetry() LEAVE\n")

	return nil
}

func (c *Client) listen() {
	klog.V(6).Infof("live.listen() ENTER\n")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			c.closeWs(false)
			klog.V(6).Infof("live.listen() Signal Exit\n")
			klog.V(6).Infof("live.listen() LEAVE\n")
			return
		case <-ticker.C:
			ws := c.internalConnect()
			if ws == nil {
				klog.V(3).Infof("listen: Connection is not valid\n")
				klog.V(6).Infof("live.listen() LEAVE\n")
				return
			}

			msgType, byMsg, err := ws.ReadMessage()
			if err != nil {
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, SuccessfulSocketErr):
					klog.V(3).Infof("Graceful websocket close\n")

					// graceful close
					c.closeWs(false)

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
					c.closeWs(true)

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
					c.closeWs(false)

					klog.V(6).Infof("live.listen() LEAVE\n")
					return
				case (err == io.EOF || err == io.ErrUnexpectedEOF) && !c.retry:
					klog.V(3).Infof("stream object EOF\n")

					// send error on callback
					sendErr := c.sendError(err)
					if sendErr != nil {
						klog.V(1).Infof("listen: EOF error. Err: %v\n", sendErr)
					}

					// close the connection
					c.closeWs(true)

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
					c.closeWs(true)

					klog.V(6).Infof("live.listen() LEAVE\n")
					return
				}
			}

			if len(byMsg) == 0 {
				klog.V(7).Infof("listen: message empty")
				continue
			}

			// inspect the message
			if c.cOptions.InspectMessage() {
				err := c.inspect(byMsg)
				if err != nil {
					klog.V(1).Infof("listen: inspect failed. Err: %v\n", err)
				}
			}

			// callback!
			if c.callback != nil {
				err := c.router.Message(byMsg)
				if err != nil {
					klog.V(1).Infof("listen: router.Message failed. Err: %v\n", err)
				}
			} else {
				klog.V(7).Infof("listen: msg recv (type %d): %s\n", msgType, string(byMsg))
			}
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
				case (err == io.EOF || err == io.ErrUnexpectedEOF) && !c.retry:
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
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")

		return err
	}

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	if err := ws.WriteMessage(
		websocket.BinaryMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteBinary WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteBinary Successful\n")
	klog.V(7).Infof("WriteBinary payload:\nData: %x\n", byData)
	klog.V(7).Infof("live.WriteBinary() LEAVE\n")

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the live transcription session on the Deepgram server.
*/
func (c *Client) WriteJSON(payload interface{}) error {
	klog.V(7).Infof("live.WriteJSON() ENTER\n")

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.WriteJSON() LEAVE\n")

		return err
	}

	byData, err := json.Marshal(payload)
	if err != nil {
		klog.V(1).Infof("WriteJSON json.Marshal failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteJSON() LEAVE\n")
		return err
	}

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	if err := ws.WriteMessage(
		websocket.TextMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteJSON WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteJSON() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteJSON payload:\nData: %s\n", string(byData))
	klog.V(7).Infof("live.WriteJSON() LEAVE\n")

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

func (c *Client) Finalize() error {
	klog.V(7).Infof("live.Finalize() ENTER\n")

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.Finalize() LEAVE\n")

		return err
	}

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"Finalize\" }"))

	klog.V(4).Infof("Finalize Succeeded\n")
	klog.V(7).Infof("live.Finalize() LEAVE\n")

	return err
}

// Stop will send close message and shutdown websocket connection
func (c *Client) Stop() {
	klog.V(3).Infof("Stopping...\n")
	c.retry = false

	// exit gracefully
	c.ctxCancel()
	c.closeWs(false)
}

func (c *Client) closeWs(fatal bool) {
	klog.V(6).Infof("live.closeWs() closing channels...\n")

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	if c.wsconn != nil && !fatal {
		// deepgram requires a close message to be sent
		errDg := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"CloseStream\" }"))
		if errDg == websocket.ErrCloseSent {
			klog.V(3).Infof("Failed to send CloseNormalClosure. Err: %v\n", errDg)
		} else if errDg != nil {
			klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", errDg)
		}
		time.Sleep(TerminationSleep) // allow time for server to register closure

		// websocket protocol message
		errProto := c.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if errProto == websocket.ErrCloseSent {
			klog.V(3).Infof("Failed to send CloseNormalClosure. Err: %v\n", errProto)
		} else if errProto != nil {
			klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", errProto)
		}
		time.Sleep(TerminationSleep) // allow time for server to register closure
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

func (c *Client) ping() {
	klog.V(6).Infof("live.ping() ENTER\n")

	counter := 0
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("live.ping() Exiting\n")

			// exit gracefully
			c.closeWs(false)

			klog.V(6).Infof("live.ping() LEAVE\n")
			return
		case <-ticker.C:
			klog.V(5).Infof("Starting ping...")
			counter++

			ws := c.internalConnect()
			if ws == nil {
				klog.V(1).Infof("ping Connection is not valid\n")
				klog.V(6).Infof("live.ping() LEAVE\n")
				return
			}

			// doing a write, need to lock.
			// Note: not using defer here because we arent leaving the scope of the function
			c.muConn.Lock()

			// deepgram keepalive message
			klog.V(5).Infof("Sending Deepgram KeepAlive message...\n")
			err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"KeepAlive\" }"))
			if err == nil {
				klog.V(5).Infof("Ping sent!")
			} else {
				klog.V(1).Infof("Failed to send Deepgram KeepAlive. Err: %v\n", err)
			}

			// release
			c.muConn.Unlock()
		}
	}
}

func (c *Client) flush() {
	klog.V(6).Infof("live.flush() ENTER\n")

	ticker := time.NewTicker(flushPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(3).Infof("live.flush() Exiting\n")

			// exit gracefully
			c.closeWs(false)

			klog.V(6).Infof("live.flush() LEAVE\n")
			return
		case <-ticker.C:
			ws := c.internalConnect()
			if ws == nil {
				klog.V(1).Infof("flush Connection is not valid\n")
				klog.V(6).Infof("live.flush() LEAVE\n")
				return
			}

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
		klog.V(1).Infof("listen: router.Error failed. Err: %v\n", sendErr)
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
		errorDesc = UnknownDeepgramErr
	}

	response := &msginterfaces.ErrorResponse{
		Type:        msginterfaces.TypeErrorResponse,
		Message:     strings.TrimSpace(fmt.Sprintf("%s %s", errorCode, errorNum)),
		Description: strings.TrimSpace(errorDesc),
		Variant:     errorNum,
	}
	return response
}

// inspectMessage inspects the message and determines the type to
// see if we should do anything with those types of messages
func (c *Client) inspect(byMsg []byte) error {
	klog.V(7).Infof("client.inspect() ENTER\n")

	var mt msginterfaces.MessageType
	if err := json.Unmarshal(byMsg, &mt); err != nil {
		klog.V(1).Infof("json.Unmarshal(MessageType) failed. Err: %v\n", err)
		klog.V(7).Infof("client.inspect() LEAVE\n")
		return err
	}

	switch mt.Type {
	case msginterfaces.TypeMessageResponse:
		klog.V(7).Infof("TypeMessageResponse\n")

		// convert to MessageResponse
		var mr msginterfaces.MessageResponse
		if err := json.Unmarshal(byMsg, &mr); err != nil {
			klog.V(1).Infof("json.Unmarshal(MessageResponse) failed. Err: %v\n", err)
			klog.V(7).Infof("client.inspect() LEAVE\n")
			return err
		}

		// inspect the message
		err := c.inspectMessage(&mr)
		if err != nil {
			klog.V(1).Infof("inspectMessage() failed. Err: %v\n", err)
			klog.V(7).Infof("client.inspect() LEAVE\n")
			return err
		}
	default:
		klog.V(7).Infof("MessageType: %s\n", mt.Type)
	}

	klog.V(7).Info("inspect() succeeded\n")
	klog.V(7).Infof("client.inspect() LEAVE\n")
	return nil
}

func (c *Client) inspectMessage(mr *msginterfaces.MessageResponse) error {
	klog.V(7).Infof("client.inspectMessage() ENTER\n")

	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)
	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		klog.V(7).Info("inspectMessage is empty\n")
		klog.V(7).Infof("client.inspectMessage() LEAVE\n")
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
	klog.V(7).Infof("client.inspectMessage() LEAVE\n")
	return nil
}
