// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
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
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// gabs "github.com/Jeffail/gabs/v2"
	"github.com/dvonthenen/websocket"

	live "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1"
	msginterfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	version "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
)

/*
NewWithDefaults() creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler which just prints all messages to the console
*/
func NewWithDefaults(ctx context.Context, apiKey string, options interfaces.LiveTranscriptionOptions) (*Client, error) {
	return New(ctx, apiKey, &ClientOptions{}, options, nil)
}

/*
New() creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- tOptions: LiveTranscriptionOptions which allows overriding things like language, model, etc.
- callback: LiveMessageCallback which is a callback that allows you to perform actions based on the transcription
*/
func New(ctx context.Context, apiKey string, cOptions *ClientOptions, tOptions interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*Client, error) {
	if apiKey == "" {
		if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
			log.Println("DEEPGRAM_API_KEY found")
			apiKey = v
		} else {
			return nil, errors.New("DEEPGRAM_API_KEY not found")
		}
	}
	if callback == nil {
		log.Printf("NewDeepGramWSClient callback is nil. Using DefaultCallbackHandler.\n")
		callback = live.NewDefaultCallbackHandler()
	}

	// init
	conn := Client{
		apiKey:   apiKey,
		cOptions: cOptions,
		tOptions: tOptions,
		sendBuf:  make(chan []byte, 1),
		callback: callback,
		router:   live.New(callback),
		org:      ctx,
		retry:    true,
	}
	conn.ctx, conn.ctxCancel = context.WithCancel(ctx)

	log.Printf("NewDeepGramWSClient Succeeded\n")
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
	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		log.Printf("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		return nil
	}

	// if the connection is good, return it
	// otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			// continue through to reconnect by recreating the wsconn object
			// log.Printf("Connection is broken. Will attempt reconnect.")
			c.ctx, c.ctxCancel = context.WithCancel(c.org)
		default:
			// log.Printf("Connection is good. Return object.")
			return c.wsconn
		}
	}

	// TODO: Disable the Hostname validation for now
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		RedirectService:  c.cOptions.RedirectService,
		SkipServerAuth:   c.cOptions.SkipServerAuth,
	}

	// set websocket headers
	myHeader := http.Header{}

	// restore application options to HTTP header
	if headers, ok := c.ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				log.Printf("Connect() RESTORE Header: %s = %s\n", k, v)
				myHeader.Add(k, v)
			}
		}
	}

	// sets the API key
	myHeader.Set("Host", c.cOptions.Host)
	myHeader.Set("Authorization", "token "+c.apiKey)
	myHeader.Set("User-Agent", interfaces.DgAgent)

	// attempt to establish connection
	i := int64(0)
	for {
		if retries != connectionRetryInfinite && i >= retries {
			log.Printf("Connect timeout... exiting!\n")
			break
		}

		// delay on subsequent calls
		if i > 0 {
			log.Printf("Sleep for retry #%d...\n", i)
			time.Sleep(time.Second * time.Duration(defaultDelayBetweenRetry))
		}

		i++

		// create new connection
		url, err := version.GetLiveAPI(c.org, c.cOptions.Host, c.cOptions.ApiVersion, version.LivePath, c.tOptions)
		if err != nil {
			log.Printf("version.GetLiveAPI failed. Err: %v\n", err)
			return nil // no point in retrying because this is going to fail on every retry
		}
		// TODO: DO NOT PRINT
		log.Printf("Connecting to %s\n", url)

		// perform the websocket connection
		ws, _, err := dialer.DialContext(c.ctx, url, myHeader)
		if err != nil {
			log.Printf("Cannot connect to websocket: %s\n", c.cOptions.Host)
			continue
		}

		// set the object to allow threads to function
		log.Printf("WebSocket Connection Successful!")
		c.wsconn = ws
		c.retry = true

		// kick off threads to listen for messages and ping/keepalive
		go c.listen()
		go c.ping()

		return c.wsconn
	}

	return nil
}

func (c *Client) listen() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			for {
				ws := c.Connect()
				if ws == nil {
					log.Printf("WebSocketClient::listen: Connection is not valid\n")
					break
				}

				msgType, byMsg, err := ws.ReadMessage()
				if err != nil {
					log.Printf("WebSocketClient::listen: Cannot read websocket message. Err: %v\n", err)
					break
				}

				if len(byMsg) == 0 {
					log.Printf("WebSocketClient::listen: message empty")
					continue
				}

				if c.callback != nil {
					c.router.Message(byMsg)
				} else {
					log.Printf("WebSocketClient::listen: msg recv (type %d): %s\n", msgType, string(byMsg))
				}
			}
		}
	}
}

// Stream is a helper function to stream audio data from a io.Reader object to deepgram
func (c *Client) Stream(r io.Reader) error {
	chunk := make([]byte, CHUNK_SIZE)

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			bytesRead, err := r.Read(chunk)
			if err != nil {
				// TODO: must put behind verbosity logging
				// log.Printf("r.Read failed. Err: %v\n", err)
				return err
			}

			if bytesRead == 0 {
				continue
			}

			_, err = c.Write(chunk[:bytesRead])
			if err != nil {
				log.Printf("w.Write failed. Err: %v\n", err)
				return err
			}
			// log.Printf("io.Writer succeeded. Bytes written: %d\n", byteCount) // TODO: debug only... delete or implement log levels
		}
	}
}

// WriteBinary writes binary data to the websocket server
func (c *Client) WriteBinary(byData []byte) error {
	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

	ws := c.Connect()
	if ws == nil {
		log.Printf("WebSocketClient::WriteBinary Connection is not valid\n")
		return ErrInvalidConnection
	}

	if err := ws.WriteMessage(
		websocket.BinaryMessage,
		byData,
	); err != nil {
		log.Printf("WebSocketClient::WriteBinary WriteMessage failed. Err: %v\n", err)
		return err
	}

	// log.Printf("WriteBinary Successful\n") // TODO: debug only... delete or implement log levels
	// log.Printf("WriteBinary payload:\nData: %x\n", byData) // TODO: debug only... delete or implement log levels

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the live transcription session on the Deepgram server.
*/
func (c *Client) WriteJSON(payload interface{}) error {
	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

	ws := c.Connect()
	if ws == nil {
		log.Printf("WebSocketClient::WriteJSON Connection is not valid\n")
		return ErrInvalidConnection
	}

	dataStruct, err := json.Marshal(payload)
	if err != nil {
		log.Printf("WebSocketClient::WriteJSON json.Marshal failed. Err: %v\n", err)
		return err
	}

	if err := ws.WriteMessage(
		websocket.TextMessage,
		dataStruct,
	); err != nil {
		log.Printf("WebSocketClient::WriteJSON WriteMessage failed. Err: %v\n", err)
		return err
	}

	// log.Printf("WriteJSON Successful\n") // TODO: debug only... delete or implement log levels
	// log.Printf("WriteJSON payload:\nData: %s\n", string(dataStruct)) // TODO: debug only... delete or implement log levels

	return nil
}

/*
Write performs the lower level websocket write operation.
This is needed to implement the io.Writer interface. (aka the streaming interface)
*/
func (c *Client) Write(p []byte) (int, error) {
	byteLen := len(p)
	err := c.WriteBinary(p)
	if err != nil {
		log.Printf("WebSocketClient::WriteBinary failed. Err: %v\n", err)
		return 0, err
	}
	return byteLen, nil
}

// Stop will send close message and shutdown websocket connection
func (c *Client) Stop() {
	log.Printf("WebSocketClient::Stop Stopping...\n")
	c.ctxCancel()
	c.closeWs()
}

func (c *Client) closeWs() {
	log.Printf("WebSocketClient::closeWs closing channels...\n")

	// doing a write, need to lock
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.wsconn != nil {
		// deepgram requires a close message to be sent
		errDg := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"CloseStream\" }"))
		if errDg != nil {
			log.Printf("Failed to send CloseNormalClosure. Err: %v\n", errDg)
		}
		time.Sleep(TERMINATION_SLEEP) // allow time for server to register closure

		// websocket protocol message
		errProto := c.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if errProto != nil {
			log.Printf("Failed to send CloseNormalClosure. Err: %v\n", errProto)
		}
		time.Sleep(TERMINATION_SLEEP) // allow time for server to register closure
		c.wsconn.Close()
	}
}

func (c *Client) ping() {
	log.Printf("WebSocketClient::ping ENTER\n")

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			log.Printf("Starting ping...")

			ws := c.Connect()
			if ws == nil {
				log.Printf("WebSocketClient::ping Connect is not valid\n")
				break
			}

			// doing a write, need to lock
			c.mu.Lock()
			log.Printf("Sending ping... need reply in %d\n", (pingPeriod / 2))

			// deepgram keepalive message
			errDg := ws.WriteMessage(websocket.BinaryMessage, []byte("{ \"type\": \"KeepAlive\" }"))
			if errDg != nil {
				log.Printf("Failed to send CloseNormalClosure. Err: %v\n", errDg)
			}

			// websocket protocol ping/pong
			errProto := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2))
			if errProto != nil {
				log.Printf("Failed to send CloseNormalClosure. Err: %v\n", errProto)
			}
			c.mu.Unlock()

			if errDg != nil || errProto != nil {
				log.Printf("WebSocketClient::ping failed\n")
				c.closeWs()
			} else {
				log.Printf("Ping sent!")
			}
		}
	}
}
