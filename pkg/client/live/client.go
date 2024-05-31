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
func (c *Client) Connect() *websocket.Conn {
	return c.ConnectWithRetry(c.ctx, c.ctxCancel, int(DefaultConnectRetry))
}

// AttemptReconnect performs a reconnect after failing retries
func (c *Client) AttemptReconnect(ctx context.Context, retries int64) *websocket.Conn {
	c.retry = true
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.ConnectWithRetry(c.ctx, c.ctxCancel, int(retries))
}

// AttemptReconnect performs a reconnect after failing retries
func (c *Client) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) *websocket.Conn {
	c.retry = true
	return c.ConnectWithRetry(ctx, ctxCancel, int(retries))
}

// ConnectWithRetry allows for connecting with specified retry attempts
//
//nolint:funlen // this is a complex function. keep as is
func (c *Client) ConnectWithRetry(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) *websocket.Conn {
	klog.V(7).Infof("live.Connect() ENTER\n")

	// set the context
	c.ctx = ctx
	c.ctxCancel = ctxCancel

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(7).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		klog.V(7).Infof("live.Connect() LEAVE\n")
		return nil
	}

	// set the retry count
	if retryCnt <= 0 {
		c.retryCnt = DefaultConnectRetry
	} else {
		c.retryCnt = int64(retryCnt)
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
		c.mu.Lock()

		// perform the websocket connection
		ws, res, err := dialer.DialContext(c.ctx, url, myHeader)
		if res != nil {
			klog.V(3).Infof("HTTP Response: %s\n", res.Status)
			res.Body.Close()
		}
		if err != nil {
			klog.V(1).Infof("Cannot connect to websocket: %s\n", c.cOptions.Host)
			klog.V(1).Infof("Dialer failed. Err: %v\n", err)
			c.mu.Unlock()
			continue
		}

		// set the object to allow threads to function
		c.wsconn = ws
		c.retry = true

		// unlock the connection
		c.mu.Unlock()

		// kick off threads to listen for messages and ping/keepalive
		go c.listen()
		if c.cOptions.EnableKeepAlive {
			go c.ping()
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
			klog.V(6).Infof("live.listen() Done\n")
			klog.V(6).Infof("live.listen() LEAVE\n")
			return
		case <-ticker.C:
			for {
				ws := c.Connect()
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
						klog.V(6).Infof("live.listen() LEAVE\n")
						return
					case strings.Contains(errStr, FatalReadSocketErr):
						klog.V(1).Infof("Fatal socket error: %v\n", err)
						c.closeWs(true)
						return
					case strings.Contains(errStr, "Deepgram"):
						klog.V(1).Infof("listen: Deepgram error. Err: %v\n", err)

						// extract DG error
						response := c.ErrorToResponse(err)
						if c.callback != nil {
							err := c.router.ErrorHelper(response)
							if err != nil {
								klog.V(1).Infof("listen: router.Error failed. Err: %v\n", err)
							}
						} else {
							klog.V(7).Infof("listen: Deepgram Error: %v\n", err)
						}

						// reset connection
						c.closeWs(true)
						return
					case (err == io.EOF || err == io.ErrUnexpectedEOF) && !c.retry:
						klog.V(3).Infof("stream object EOF\n")
						klog.V(6).Infof("live.listen() LEAVE\n")
						return
					default:
						klog.V(1).Infof("listen: Cannot read websocket message. Err: %v\n", err)
					}
					continue
				}

				if len(byMsg) == 0 {
					klog.V(7).Infof("listen: message empty")
					continue
				}

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
					return nil
				case (err == io.EOF || err == io.ErrUnexpectedEOF) && !c.retry:
					klog.V(3).Infof("stream object EOF\n")
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return nil
				case err != nil:
					klog.V(1).Infof("r.Read encountered EOF. Err: %v\n", err)
					klog.V(6).Infof("live.Stream() LEAVE\n")
					return err
				}
			}

			if bytesRead == 0 {
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
	ws := c.Connect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")

		return err
	}

	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

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
	ws := c.Connect()
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
	c.mu.Lock()
	defer c.mu.Unlock()

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
	ws := c.Connect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("live.Finalize() LEAVE\n")

		return err
	}

	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

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
	c.mu.Lock()
	defer c.mu.Unlock()

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
			klog.V(6).Infof("live.ping() LEAVE\n")
			return
		case <-ticker.C:
			klog.V(5).Infof("Starting ping...")
			counter++

			ws := c.Connect()
			if ws == nil {
				klog.V(1).Infof("ping Connection is not valid\n")
				klog.V(6).Infof("live.ping() LEAVE\n")
				return
			}

			// doing a write, need to lock
			c.mu.Lock()

			// deepgram keepalive message
			klog.V(5).Infof("Sending Deepgram KeepAlive message...\n")
			err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"KeepAlive\" }"))
			if err == nil {
				klog.V(5).Infof("Ping sent!")
			} else {
				klog.V(1).Infof("Failed to send Deepgram KeepAlive. Err: %v\n", err)
			}

			// release
			c.mu.Unlock()
		}
	}
}

func (c *Client) ErrorToResponse(err error) *msginterfaces.ErrorResponse {
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
		Message:     fmt.Sprintf("%s %s", errorCode, errorNum),
		Description: errorDesc,
		Variant:     errorNum,
	}
	return response
}
