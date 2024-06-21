// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the speak/streaming client implementation for the Deepgram API
*/
package speak

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

	speak "github.com/deepgram/deepgram-go-sdk/pkg/api/speak-stream/v1"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak-stream/v1/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/*
NewStreamForDemo creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewStreamForDemo(ctx context.Context, options *interfaces.SpeakOptions) (*StreamClient, error) {
	return NewStream(ctx, "", &interfaces.ClientOptions{}, options, nil)
}

/*
NewStreamWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
  - The callback handler is set to the default handler
*/
func NewStreamWithDefaults(ctx context.Context, options *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*StreamClient, error) {
	return NewStream(ctx, "", &interfaces.ClientOptions{}, options, callback)
}

/*
NewStream creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages
*/
func NewStream(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*StreamClient, error) {
	ctx, ctxCancel := context.WithCancel(ctx)
	return NewStreamWithCancel(ctx, ctxCancel, apiKey, cOptions, sOptions, callback)
}

/*
NewStream creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- ctxCancel: allow passing in own cancel
- apiKey: string containing the Deepgram API key
- cOptions: ClientOptions which allows overriding things like hostname, version of the API, etc.
- sOptions: SpeakOptions which allows overriding things like model, etc.
- callback: SpeakMessageCallback is a callback which lets you perform actions based on platform messages
*/
func NewStreamWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.SpeakOptions, callback msginterfaces.SpeakMessageCallback) (*StreamClient, error) {
	klog.V(6).Infof("StreamClient.New() ENTER\n")

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
	conn := StreamClient{
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
	klog.V(6).Infof("StreamClient.New() LEAVE\n")

	return &conn, nil
}

// Connect performs a websocket connection with "DefaultConnectRetry" number of retries.
func (c *StreamClient) Connect() bool {
	// set the retry count
	if c.retryCnt == 0 {
		c.retryCnt = DefaultConnectRetry
	}
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt)) != nil
}

// ConnectWithCancel performs a websocket connection with specified number of retries and providing a
// cancel function to stop the connection
func (c *StreamClient) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	return c.internalConnectWithCancel(ctx, ctxCancel, retryCnt) != nil
}

// AttemptReconnect performs a reconnect after failing retries
func (c *StreamClient) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.retry = true
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(retries)) != nil
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *StreamClient) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.retry = true
	return c.internalConnectWithCancel(ctx, ctxCancel, int(retries)) != nil
}

func (c *StreamClient) internalConnect() *websocket.Conn {
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt))
}

//nolint:funlen // this is a complex function. keep as is
func (c *StreamClient) internalConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) *websocket.Conn {
	klog.V(7).Infof("StreamClient.Connect() ENTER\n")

	// set the context
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	c.retryCnt = int64(retryCnt)

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(7).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		klog.V(7).Infof("StreamClient.Connect() LEAVE\n")
		return nil
	}

	// if the connection is good, return it otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Connection is not valid\n")
			klog.V(7).Infof("StreamClient.ConnectWithRetry() LEAVE\n")
			return nil
		default:
			klog.V(7).Infof("Connection is good. Return object.")
			klog.V(7).Infof("StreamClient.ConnectWithRetry() LEAVE\n")
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
		url, err := version.GetSpeakStreamAPI(c.ctx, c.cOptions.Host, c.cOptions.APIVersion, c.cOptions.Path, c.sOptions)
		if err != nil {
			klog.V(1).Infof("version.GetSpeakAPI failed. Err: %v\n", err)
			klog.V(7).Infof("StreamClient.ConnectWithRetry() LEAVE\n")
			return nil // no point in retrying because this is going to fail on every retry
		}
		klog.V(5).Infof("Connecting to %s\n", url)

		// a single connection attempt
		// Note: not using defer here because we aren't leaving the scope of the function
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
		klog.V(7).Infof("StreamClient.ConnectWithRetry() LEAVE\n")

		return c.wsconn
	}

	// if we get here, we failed to connect
	klog.V(1).Infof("Failed to connect to websocket: %s\n", c.cOptions.Host)
	klog.V(7).Infof("StreamClient.ConnectWithRetry() LEAVE\n")

	return nil
}

func (c *StreamClient) listen() {
	klog.V(6).Infof("StreamClient.listen() ENTER\n")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			c.closeWs(false)
			klog.V(6).Infof("StreamClient.listen() Signal Exit\n")
			klog.V(6).Infof("StreamClient.listen() LEAVE\n")
			return
		case <-ticker.C:
			ws := c.internalConnect()
			if ws == nil {
				klog.V(3).Infof("StreamClient.listen(): Connection is not valid\n")
				klog.V(6).Infof("StreamClient.listen() LEAVE\n")
				return
			}

			// msgType can be binary or text
			msgType, byMsg, err := ws.ReadMessage()
			if err != nil {
				errStr := err.Error()
				switch {
				case strings.Contains(errStr, SuccessfulSocketErr):
					klog.V(3).Infof("Graceful websocket close\n")

					// graceful close
					c.closeWs(false)

					klog.V(6).Infof("StreamClient.listen() LEAVE\n")
					return
				case strings.Contains(errStr, FatalReadSocketErr):
					klog.V(1).Infof("Fatal socket error: %v\n", err)

					// send error on callback
					sendErr := c.sendError(err)
					if sendErr != nil {
						klog.V(1).Infof("StreamClient.listen(): Fatal socket error. Err: %v\n", sendErr)
					}

					// fatal close
					c.closeWs(true)

					klog.V(6).Infof("StreamClient.listen() LEAVE\n")
					return
				case strings.Contains(errStr, "Deepgram"):
					klog.V(1).Infof("StreamClient.listen(): Deepgram error. Err: %v\n", err)

					// send error on callback
					sendErr := c.sendError(err)
					if sendErr != nil {
						klog.V(1).Infof("StreamClient.listen(): Deepgram ErrorMsg. Err: %v\n", sendErr)
					}

					// close the connection
					c.closeWs(false)

					klog.V(6).Infof("StreamClient.listen() LEAVE\n")
					return
				case (err == io.EOF || err == io.ErrUnexpectedEOF) && !c.retry:
					klog.V(3).Infof("StreamClient object EOF\n")

					// send error on callback
					sendErr := c.sendError(err)
					if sendErr != nil {
						klog.V(1).Infof("StreamClient.listen(): EOF error. Err: %v\n", sendErr)
					}

					// close the connection
					c.closeWs(true)

					klog.V(6).Infof("StreamClient.listen() LEAVE\n")
					return
				default:
					klog.V(1).Infof("StreamClient.listen(): Cannot read websocket message. Err: %v\n", err)

					// send error on callback
					sendErr := c.sendError(err)
					if sendErr != nil {
						klog.V(1).Infof("StreamClient.listen(): EOF error. Err: %v\n", sendErr)
					}

					// close the connection
					c.closeWs(true)

					klog.V(6).Infof("StreamClient.listen() LEAVE\n")
					return
				}
			}

			if len(byMsg) == 0 {
				klog.V(7).Infof("StreamClient.listen(): message empty")
				continue
			}

			// inspect the message
			// if c.cOptions.InspectMessage() {
			// 	err := c.inspect(byMsg)
			// 	if err != nil {
			// 		klog.V(1).Infof("StreamClient.listen(): inspect failed. Err: %v\n", err)
			// 	}
			// }

			if c.callback != nil {
				if msgType == websocket.TextMessage {
					err := c.router.Message(byMsg)
					if err != nil {
						klog.V(1).Infof("StreamClient.listen(): router.Message failed. Err: %v\n", err)
					}

				} else if msgType == websocket.BinaryMessage {
					err := c.router.Binary(byMsg)
					if err != nil {
						klog.V(1).Infof("StreamClient.listen(): router.Message failed. Err: %v\n", err)
					}
				} else {
					klog.V(7).Infof("StreamClient.listen(): msg recv (type %d): %s\n", msgType, string(byMsg))
				}
			} else {
				klog.V(7).Infof("StreamClient.listen(): msg recv (type %d): %s\n", msgType, string(byMsg))
			}
		}
	}
}

// WriteBinary writes binary data to the websocket server
func (c *StreamClient) WriteBinary(byData []byte) error {
	klog.V(7).Infof("StreamClient.WriteBinary() ENTER\n")

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("StreamClient.WriteBinary() LEAVE\n")

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
		klog.V(7).Infof("StreamClient.WriteBinary() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteBinary Successful\n")
	klog.V(7).Infof("WriteBinary payload:\nData: %x\n", byData)
	klog.V(7).Infof("StreamClient.WriteBinary() LEAVE\n")

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the text-to-speech session on the Deepgram server.
*/
func (c *StreamClient) WriteJSON(payload interface{}) error {
	klog.V(7).Infof("StreamClient.WriteJSON() ENTER\n")

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("StreamClient.WriteJSON() LEAVE\n")

		return err
	}

	byData, err := json.Marshal(payload)
	if err != nil {
		klog.V(1).Infof("WriteJSON: Error marshaling JSON. Data: %v, Err: %v\n", payload, err)
		klog.V(7).Infof("StreamClient.WriteJSON() LEAVE\n")
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
		klog.V(7).Infof("StreamClient.WriteJSON() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteJSON payload:\nData: %s\n", string(byData))
	klog.V(7).Infof("StreamClient.WriteJSON() LEAVE\n")

	return nil
}

/*
Write performs the lower level websocket write operation.
This is needed to implement the io.Writer interface. (aka the streaming interface)
*/
func (c *StreamClient) Write(p []byte) (int, error) {
	klog.V(7).Infof("StreamClient.Write() ENTER\n")

	byteLen := len(p)
	err := c.WriteBinary(p)
	if err != nil {
		klog.V(1).Infof("Write failed. Err: %v\n", err)
		klog.V(7).Infof("StreamClient.Write() LEAVE\n")
		return 0, err
	}

	klog.V(7).Infof("StreamClient.Write Succeeded\n")
	klog.V(7).Infof("StreamClient.Write() LEAVE\n")
	return byteLen, nil
}

func (c *StreamClient) Flush() error {
	klog.V(7).Infof("StreamClient.Flush() ENTER\n")

	// doing a write, need to lock
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.Connect() is nil. Err: %v\n", err)
		klog.V(7).Infof("StreamClient.Flush() LEAVE\n")

		return err
	}

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	err := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"Flush\" }"))

	klog.V(4).Infof("Flush Succeeded\n")
	klog.V(7).Infof("StreamClient.Flush() LEAVE\n")

	return err
}

// Reset will instruct the server to reset the current buffer
func (c *StreamClient) Reset() error {
	klog.V(6).Infof("StreamClient.Reset() ENTER\n")

	resetMessage := map[string]string{
		"type": "Reset",
	}

	msg, err := json.Marshal(resetMessage)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("StreamClient.Reset() LEAVE\n")
		return err
	}

	err = c.WriteJSON(msg)
	if err != nil {
		klog.V(1).Infof("WriteJSON failed. Err: %v\n", err)
		klog.V(6).Infof("StreamClient.Reset() LEAVE\n")
		return err
	}

	klog.V(4).Infof("Reset Succeeded\n")
	klog.V(6).Infof("StreamClient.Reset() LEAVE\n")
	return nil
}

// Stop will send close message and shutdown websocket connection
func (c *StreamClient) Stop() {
	klog.V(3).Infof("Stopping...\n")
	c.retry = false

	// exit gracefully
	c.ctxCancel()
	c.closeWs(false)
}

func (c *StreamClient) closeWs(fatal bool) {
	klog.V(6).Infof("StreamClient.closeWs() closing channels...\n")

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	if c.wsconn != nil && !fatal {
		// deepgram requires a close message to be sent
		errDg := c.wsconn.WriteMessage(websocket.TextMessage, []byte("{ \"type\": \"Close\" }"))
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

	klog.V(4).Infof("StreamClient.closeWs() Succeeded\n")
	klog.V(6).Infof("StreamClient.closeWs() LEAVE\n")
}

// sendError sends an error message to the callback handler
func (c *StreamClient) sendError(err error) error {
	response := c.errorToResponse(err)
	sendErr := c.router.ErrorHelper(response)
	if err != nil {
		klog.V(1).Infof("StreamClient.listen(): router.Error failed. Err: %v\n", sendErr)
	}

	return err
}

// errorToResponse converts an error into a Deepgram error response
func (c *StreamClient) errorToResponse(err error) *msginterfaces.ErrorResponse {
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
		ErrMsg:      strings.TrimSpace(fmt.Sprintf("%s %s", errorCode, errorNum)),
		Description: strings.TrimSpace(errorDesc),
		Variant:     errorNum,
	}
	return response
}
