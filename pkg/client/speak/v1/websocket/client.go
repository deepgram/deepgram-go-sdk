// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the speak/streaming client implementation for the Deepgram API
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

	speak "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

type controlMessage struct {
	Type string `json:"type"`
}
type TextSource struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

/*
NewWebSocketForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWebSocketForDemo(ctx context.Context, options *interfaces.SpeakOptions) (*Client, error) {
	return NewWebSocket(ctx, "", &interfaces.ClientOptions{}, options, nil)
}

/*
NewWebSocketWithDefaults creates a new websocket connection with all default options

Notes:
  - The callback handler is set to the default handler
*/
func NewWebSocketWithDefaults(ctx context.Context, options *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*Client, error) {
	return NewWebSocket(ctx, "", &interfaces.ClientOptions{}, options, callback)
}

/*
NewWebSocket creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewWebSocket(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*Client, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewWebSocketWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}

/*
NewWebSocketWithCancel creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages

Notes:
  - If apiKey is an empty string, the Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewWebSocketWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*Client, error) {
	klog.V(6).Infof("speak.New() ENTER\n")

	if apiKey != "" {
		cOptions.APIKey = apiKey
	}
	err := cOptions.Parse()
	if err != nil {
		klog.V(1).Infof("ClientOptions.Parse() failed. Err: %v\n", err)
		return nil, err
	}
	err = sOptions.Check()
	if err != nil {
		klog.V(1).Infof("SpeakOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	if callback == nil {
		klog.V(2).Infof("Using DefaultCallbackHandler.\n")
		callback = speak.NewDefaultCallbackHandler()
	}

	// init
	conn := Client{
		cOptions:  cOptions,
		sOptions:  sOptions,
		sendBuf:   make(chan []byte, 1),
		callback:  callback,
		router:    speak.NewStream(callback),
		ctx:       ctx,
		ctxCancel: ctxCancel,
		retry:     true,
	}

	klog.V(3).Infof("NewDeepGramWSClient Succeeded\n")
	klog.V(6).Infof("speak.New() LEAVE\n")

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
	c.retry = true
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(retries), true) != nil
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *Client) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.retry = true
	return c.internalConnectWithCancel(ctx, ctxCancel, int(retries), true) != nil
}

func (c *Client) internalConnect() *websocket.Conn {
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt), false)
}

//nolint:funlen // this is a complex function. keep as is
func (c *Client) internalConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int, lock bool) *websocket.Conn {
	klog.V(7).Infof("speak.internalConnectWithCancel() ENTER\n")

	// set the context
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	c.retryCnt = int64(retryCnt)

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(7).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		klog.V(7).Infof("speak.internalConnectWithCancel() LEAVE\n")
		return nil
	}

	// lock conn access
	if lock {
		klog.V(3).Infof("Locking connection mutex\n")
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	// if the connection is good, return it otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Connection is not valid\n")
			klog.V(7).Infof("speak.internalConnectWithCancel() LEAVE\n")
			return nil
		default:
			klog.V(7).Infof("Connection is good. Return object.")
			klog.V(7).Infof("speak.internalConnectWithCancel() LEAVE\n")
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
		url, err := version.GetSpeakStreamAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.sOptions)
		if err != nil {
			klog.V(1).Infof("version.GetSpeakAPI failed. Err: %v\n", err)
			klog.V(7).Infof("speak.internalConnectWithCancel() LEAVE\n")
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

		// kick off threads to listen for messages
		go c.listen()

		// fire off open connection
		err = c.router.OpenHelper(&msginterfaces.OpenResponse{
			Type: msginterfaces.TypeOpenResponse,
		})
		if err != nil {
			klog.V(1).Infof("router.OpenHelper failed. Err: %v\n", err)
		}

		klog.V(3).Infof("WebSocket Connection Successful!")
		klog.V(7).Infof("speak.internalConnectWithCancel() LEAVE\n")

		return c.wsconn
	}

	// if we get here, we failed to connect
	klog.V(1).Infof("Failed to connect to websocket: %s\n", c.cOptions.Host)
	klog.V(7).Infof("speak.ConnectWithRetry() LEAVE\n")

	return nil
}

//nolint:funlen,gocyclo // this is a complex function. keep as is
func (c *Client) listen() {
	klog.V(6).Infof("speak.listen() ENTER\n")

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
			c.closeWs(true)

			klog.V(6).Infof("live.flush() LEAVE\n")
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
				c.closeWs(false)

				klog.V(6).Infof("speak.listen() LEAVE\n")
				return
			case strings.Contains(errStr, UseOfClosedSocket):
				klog.V(3).Infof("Probable graceful websocket close: %v\n", err)

				// fatal close
				c.closeWs(false)

				klog.V(6).Infof("speak.listen() LEAVE\n")
				return
			case strings.Contains(errStr, FatalReadSocketErr):
				klog.V(1).Infof("Fatal socket error: %v\n", err)

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("speak.listen(): Fatal socket error. Err: %v\n", sendErr)
				}

				// fatal close
				c.closeWs(true)

				klog.V(6).Infof("speak.listen() LEAVE\n")
				return
			case strings.Contains(errStr, "Deepgram"):
				klog.V(1).Infof("speak.listen(): Deepgram error. Err: %v\n", err)

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("speak.listen(): Deepgram ErrorMsg. Err: %v\n", sendErr)
				}

				// close the connection
				c.closeWs(false)

				klog.V(6).Infof("speak.listen() LEAVE\n")
				return
			case (err == io.EOF || err == io.ErrUnexpectedEOF) && !c.retry:
				klog.V(3).Infof("Client object EOF\n")

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("speak.listen(): EOF error. Err: %v\n", sendErr)
				}

				// close the connection
				c.closeWs(true)

				klog.V(6).Infof("speak.listen() LEAVE\n")
				return
			default:
				klog.V(1).Infof("speak.listen(): Cannot read websocket message. Err: %v\n", err)

				// send error on callback
				sendErr := c.sendError(err)
				if sendErr != nil {
					klog.V(1).Infof("speak.listen(): EOF error. Err: %v\n", sendErr)
				}

				// close the connection
				c.closeWs(true)

				klog.V(6).Infof("speak.listen() LEAVE\n")
				return
			}
		}

		if len(byMsg) == 0 {
			klog.V(7).Infof("listen(): message empty")
			continue
		}

		// inspect the message
		// if c.cOptions.InspectMessage() {
		// 	err := c.inspect(byMsg)
		// 	if err != nil {
		// 		klog.V(1).Infof("speak.listen(): inspect failed. Err: %v\n", err)
		// 	}
		// }

		switch msgType {
		case websocket.TextMessage:
			err := c.router.Message(byMsg)
			if err != nil {
				klog.V(1).Infof("speak.listen(): router.Message failed. Err: %v\n", err)
			}
		case websocket.BinaryMessage:
			err := c.router.Binary(byMsg)
			if err != nil {
				klog.V(1).Infof("speak.listen(): router.Message failed. Err: %v\n", err)
			}
		default:
			klog.V(7).Infof("speak.listen(): msg recv: type %d, len: %d\n", msgType, len(byMsg))
		}
	}
}

// SpeakWithText writes binary data to the websocket server
func (c *Client) SpeakWithText(text string) error {
	klog.V(6).Infof("speak.SpeakText() ENTER\n")
	klog.V(4).Infof("text: %s\n", text)

	err := c.WriteJSON(TextSource{
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

// SpeakWithStream writes binary data to the websocket server
// NOTE: This is unimplemented on the server side
func (c *Client) SpeakWithStream(byData []byte) error {
	klog.V(6).Infof("speak.SpeakText() ENTER\n")

	err := c.WriteBinary(byData)
	if err == nil {
		klog.V(4).Infof("SpeakText Succeeded\n")
	} else {
		klog.V(1).Infof("SpeakText failed. Err: %v\n", err)
	}

	klog.V(6).Infof("speak.SpeakText() LEAVE\n")

	return err
}

// WriteBinary writes binary data to the websocket server
func (c *Client) WriteBinary(byData []byte) error {
	klog.V(6).Infof("speak.WriteBinary() ENTER\n")

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	// get the connection
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(1).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(6).Infof("speak.WriteBinary() LEAVE\n")

		return err
	}

	if err := ws.WriteMessage(
		websocket.BinaryMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteBinary WriteMessage failed. Err: %v\n", err)
		klog.V(6).Infof("speak.WriteBinary() LEAVE\n")
		return err
	}

	klog.V(6).Infof("WriteBinary Successful\n")
	klog.V(7).Infof("payload: %x\n", byData)
	klog.V(6).Infof("speak.WriteBinary() LEAVE\n")

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the text-to-speech session on the Deepgram server.
*/
func (c *Client) WriteJSON(payload interface{}) error {
	klog.V(6).Infof("speak.WriteJSON() ENTER\n")

	byData, err := json.Marshal(payload)
	if err != nil {
		klog.V(1).Infof("WriteJSON: Error marshaling JSON. Data: %v, Err: %v\n", payload, err)
		klog.V(6).Infof("speak.WriteJSON() LEAVE\n")
		return err
	}

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(1).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(6).Infof("speak.WriteJSON() LEAVE\n")

		return err
	}
	if err := ws.WriteMessage(
		websocket.TextMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteJSON WriteMessage failed. Err: %v\n", err)
		klog.V(6).Infof("speak.WriteJSON() LEAVE\n")
		return err
	}

	klog.V(4).Infof("WriteJSON succeeded.\n")
	klog.V(7).Infof("payload: %s\n", string(byData))
	klog.V(6).Infof("speak.WriteJSON() LEAVE\n")

	return nil
}

// Flush will instruct the server to flush the current text buffer
func (c *Client) Flush() error {
	klog.V(6).Infof("speak.Flush() ENTER\n")

	err := c.WriteJSON(controlMessage{Type: MessageTypeFlush})
	if err != nil {
		klog.V(1).Infof("Flush failed. Err: %v\n", err)
		klog.V(6).Infof("speak.Flush() LEAVE\n")

		return err
	}

	klog.V(4).Infof("Flush Succeeded\n")
	klog.V(6).Infof("speak.Flush() LEAVE\n")

	return err
}

// Reset will instruct the server to reset the current buffer
func (c *Client) Reset() error {
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

func (c *Client) closeStream(lock bool) error {
	klog.V(6).Infof("speak.closeStream() ENTER\n")

	// doing a write, need to lock
	if lock {
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"Close\" }"))
	if err != nil {
		klog.V(1).Infof("WriteMessage failed. Err: %v\n", err)
		klog.V(6).Infof("speak.closeStream() LEAVE\n")

		return err
	}

	klog.V(4).Infof("closeStream Succeeded\n")
	klog.V(6).Infof("speak.closeStream() LEAVE\n")

	return err
}

func (c *Client) normalClosure(lock bool) error {
	klog.V(6).Infof("speak.normalClosure() ENTER\n")

	// doing a write, need to lock
	if lock {
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(1).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(6).Infof("speak.normalClosure() LEAVE\n")

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

	klog.V(6).Infof("speak.normalClosure() LEAVE\n")

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
	klog.V(6).Infof("speak.closeWs() closing channels...\n")

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

	klog.V(4).Infof("speak.closeWs() Succeeded\n")
	klog.V(6).Infof("speak.closeWs() LEAVE\n")
}

// sendError sends an error message to the callback handler
func (c *Client) sendError(err error) error {
	response := c.errorToResponse(err)
	sendErr := c.router.ErrorHelper(response)
	if err != nil {
		klog.V(1).Infof("speak.listen(): router.Error failed. Err: %v\n", sendErr)
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
