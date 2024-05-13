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
	"io"
	"net/http"
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
		cOptions: cOptions,
		tOptions: tOptions,
		sendBuf:  make(chan []byte, 1),
		callback: callback,
		router:   live.New(callback),
		org:      ctx,
		retry:    true,
	}
	conn.ctx, conn.ctxCancel = context.WithCancel(ctx)

	klog.V(3).Infof("NewDeepGramWSClient Succeeded\n")
	klog.V(6).Infof("live.New() LEAVE\n")

	return &conn, nil
}

// Connect performs a websocket connection with "defaultConnectRetry" number of retries.
func (c *Client) Connect() *websocket.Conn {
	return c.ConnectWithRetry(defaultConnectRetry)
}

// AttemptReconnect does exactly that with "retries" number of retries.
// If "retries" is set to -1, then it will retry forever.
func (c *Client) AttemptReconnect(retries int64) *websocket.Conn {
	c.retry = true
	return c.ConnectWithRetry(retries)
}

// ConnectWithRetry is a function to explicitly do a connection with "retries" number of retries.
func (c *Client) ConnectWithRetry(retries int64) *websocket.Conn {
	klog.V(7).Infof("live.ConnectWithRetry() ENTER\n")

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(3).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		return nil
	}

	// if the connection is good, return it otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			// continue through to reconnect by recreating the wsconn object
			klog.V(3).Infof("Connection is broken. Will attempt reconnect.")
			c.ctx, c.ctxCancel = context.WithCancel(c.org)
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
		if retries != connectionRetryInfinite && i >= retries {
			klog.V(3).Infof("Connect timeout... exiting!\n")
			break
		}

		// delay on subsequent calls
		if i > 0 {
			klog.V(2).Infof("Sleep for retry #%d...\n", i)
			time.Sleep(time.Second * time.Duration(defaultDelayBetweenRetry))
		}

		i++

		// create new connection
		url, err := version.GetLiveAPI(c.org, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.tOptions)
		if err != nil {
			klog.V(1).Infof("version.GetLiveAPI failed. Err: %v\n", err)
			klog.V(7).Infof("live.ConnectWithRetry() LEAVE\n")
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
			continue
		}

		// set the object to allow threads to function
		c.wsconn = ws
		c.retry = true

		// kick off threads to listen for messages and ping/keepalive
		go c.listen()
		go c.ping()

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

	return nil
}

func (c *Client) listen() {
	klog.V(6).Infof("live.listen() ENTER\n")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			klog.V(6).Infof("live.listen() Done\n")
			klog.V(6).Infof("live.listen() LEAVE\n")
			return
		case <-ticker.C:
			for {
				ws := c.Connect()
				if ws == nil {
					klog.V(3).Infof("WebSocketClient::listen: Connection is not valid\n")
					break
				}

				msgType, byMsg, err := ws.ReadMessage()
				if err != nil {
					klog.V(3).Infof("WebSocketClient::listen: Cannot read websocket message. Err: %v\n", err)
					break
				}

				if len(byMsg) == 0 {
					klog.V(7).Infof("WebSocketClient::listen: message empty")
					continue
				}

				if c.callback != nil {
					err := c.router.Message(byMsg)
					if err != nil {
						klog.V(1).Infof("WebSocketClient::listen: router.Message failed. Err: %v\n", err)
					}
				} else {
					klog.V(7).Infof("WebSocketClient::listen: msg recv (type %d): %s\n", msgType, string(byMsg))
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
			if err == io.EOF && !c.retry {
				klog.V(3).Infof("stream object EOF\n")
				klog.V(6).Infof("live.Stream() LEAVE\n")
				return nil
			} else if err != nil {
				klog.V(1).Infof("r.Read encountered EOF. Err: %v\n", err)
				klog.V(6).Infof("live.Stream() LEAVE\n")
				return err
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
	c.mu.Lock()
	defer c.mu.Unlock()

	ws := c.Connect()
	if ws == nil {
		klog.V(1).Infof("WebSocketClient::WriteBinary Connection is not valid\n")
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")
		return ErrInvalidConnection
	}

	if err := ws.WriteMessage(
		websocket.BinaryMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WebSocketClient::WriteBinary WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteBinary() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteBinary Successful\n")
	klog.V(7).Infof("WriteBinary payload:\nData: %x\n", byData)

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the live transcription session on the Deepgram server.
*/
func (c *Client) WriteJSON(payload interface{}) error {
	klog.V(7).Infof("live.WriteJSON() ENTER\n")

	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

	ws := c.Connect()
	if ws == nil {
		klog.V(1).Infof("WebSocketClient::WriteJSON Connection is not valid\n")
		klog.V(7).Infof("live.WriteJSON() LEAVE\n")
		return ErrInvalidConnection
	}

	byData, err := json.Marshal(payload)
	if err != nil {
		klog.V(1).Infof("WebSocketClient::WriteJSON json.Marshal failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteJSON() LEAVE\n")
		return err
	}

	if err := ws.WriteMessage(
		websocket.TextMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WebSocketClient::WriteJSON WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("live.WriteJSON() LEAVE\n")
		return err
	}

	klog.V(4).Infof("WriteJSON payload:\nData: %s\n", string(byData))
	klog.V(7).Infof("live.Write() LEAVE\n")

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
		klog.V(1).Infof("WebSocketClient::WriteBinary failed. Err: %v\n", err)
		klog.V(7).Infof("live.Write() LEAVE\n")
		return 0, err
	}

	klog.V(7).Infof("live.Write() Succeeded\n")
	klog.V(7).Infof("live.Write() LEAVE\n")
	return byteLen, nil
}

func (c *Client) Finalize() error {
	klog.V(7).Infof("live.Finalize() ENTER\n")

	if c.wsconn == nil {
		err := ErrInvalidConnection

		klog.V(4).Infof("Finalize Failed. Err: %v\n", err)
		klog.V(7).Infof("live.Finalize() LEAVE\n")

		return err
	}

	err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"Finalize\" }"))

	klog.V(4).Infof("Finalize Succeeded\n")
	klog.V(7).Infof("live.Finalize() LEAVE\n")

	return err
}

// Stop will send close message and shutdown websocket connection
func (c *Client) Stop() {
	klog.V(3).Infof("WebSocketClient::Stop Stopping...\n")
	c.retry = false
	c.ctxCancel()
	c.closeWs()
}

func (c *Client) closeWs() {
	klog.V(6).Infof("live.closeWs() closing channels...\n")

	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.wsconn != nil {
		// deepgram requires a close message to be sent
		errDg := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"CloseStream\" }"))
		if errDg == websocket.ErrCloseSent {
			klog.V(3).Infof("Failed to send CloseNormalClosure. Err: %v\n", errDg)
		} else if errDg != nil {
			klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", errDg)
		}
		time.Sleep(TerminationSleep) // allow time for server to register closure

		// fire off close connection
		err := c.router.CloseHelper(&msginterfaces.CloseResponse{
			Type: msginterfaces.TypeCloseResponse,
		})
		if err != nil {
			klog.V(1).Infof("router.CloseHelper failed. Err: %v\n", err)
		}

		// websocket protocol message
		errProto := c.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if errProto == websocket.ErrCloseSent {
			klog.V(3).Infof("Failed to send CloseNormalClosure. Err: %v\n", errDg)
		} else if errProto != nil {
			klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", errProto)
		}
		time.Sleep(TerminationSleep) // allow time for server to register closure
		c.wsconn.Close()
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
				klog.V(1).Infof("WebSocketClient::ping Connect is not valid\n")
				break
			}

			// doing a write, need to lock
			c.mu.Lock()
			klog.V(5).Infof("Sending ping... need reply in %d\n", (pingPeriod / 2))

			var errDg error
			if c.cOptions.EnableKeepAlive {
				klog.V(5).Infof("Sending Deepgram KeepAlive message...\n")
				// deepgram keepalive message
				errDg = ws.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"KeepAlive\" }"))
				if errDg != nil {
					klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", errDg)
				}
			}

			// websocket protocol ping/pong... this loop is every 5 seconds, so ping every 20 seconds
			var errProto error
			errProto = nil
			if counter%4 == 0 {
				klog.V(5).Infof("Sending Protocol KeepAlive message...\n")
				errProto = ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2))
				if errProto != nil {
					klog.V(1).Infof("Failed to send CloseNormalClosure. Err: %v\n", errProto)
				}
			}
			c.mu.Unlock()

			if errDg != nil || errProto != nil {
				klog.V(1).Infof("WebSocketClient::ping failed\n")
				c.closeWs()
			} else {
				klog.V(5).Infof("Ping sent!")
			}
		}
	}
}
