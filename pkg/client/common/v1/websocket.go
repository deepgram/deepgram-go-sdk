// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the live/streaming client implementation for the Deepgram API
package commonv1

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	commonv1interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

// gocritic:ignore
func NewWS(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, options *clientinterfaces.ClientOptions, processMessages *commonv1interfaces.WebSocketHandler, router *commonv1interfaces.Router) *WSClient {
	if apiKey != "" {
		options.APIKey = apiKey
	}
	err := options.Parse()
	if err != nil {
		klog.V(1).Infof("options.Parse() failed. Err: %v\n", err)
		return nil
	}

	c := WSClient{
		cOptions:        options,
		sendBuf:         make(chan []byte, 1),
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		retry:           true,
		processMessages: processMessages,
		router:          router,
	}

	return &c
}

// Connect performs a websocket connection with "DefaultConnectRetry" number of retries.
func (c *WSClient) Connect() bool {
	// set the retry count
	if c.retryCnt == 0 {
		c.retryCnt = DefaultConnectRetry
	}
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt), true) != nil
}

// ConnectWithCancel performs a websocket connection with specified number of retries and providing a
// cancel function to stop the connection
func (c *WSClient) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool {
	return c.internalConnectWithCancel(ctx, ctxCancel, retryCnt, true) != nil
}

// AttemptReconnect performs a reconnect after failing retries
func (c *WSClient) AttemptReconnect(ctx context.Context, retries int64) bool {
	c.muConn.Lock()
	c.retry = true
	c.muConn.Unlock()
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(retries), true) != nil
}

// AttemptReconnect performs a reconnect after failing retries and providing a cancel function
func (c *WSClient) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool {
	c.muConn.Lock()
	c.retry = true
	c.muConn.Unlock()
	return c.internalConnectWithCancel(ctx, ctxCancel, int(retries), true) != nil
}

func (c *WSClient) internalConnect() *websocket.Conn {
	return c.internalConnectWithCancel(c.ctx, c.ctxCancel, int(c.retryCnt), false)
}

//nolint:funlen,gocyclo // this is a complex function. keep as is
func (c *WSClient) internalConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int, lock bool) *websocket.Conn {
	klog.V(7).Infof("common.internalConnectWithCancel() ENTER\n")

	// set the context
	c.ctx = ctx
	c.ctxCancel = ctxCancel
	c.retryCnt = int64(retryCnt)

	// lock conn access
	if lock {
		klog.V(3).Infof("Locking connection mutex\n")
		c.muConn.Lock()
	}

	// we explicitly stopped and should not attempt to reconnect
	if !c.retry {
		klog.V(7).Infof("This connection has been terminated. Please either call with AttemptReconnect or create a new Client object using NewWebSocketClient.")
		klog.V(7).Infof("common.internalConnectWithCancel() LEAVE\n")
		if lock {
			klog.V(3).Infof("Unlocking connection mutex\n")
			c.muConn.Unlock()
		}
		return nil
	}

	// if the connection is good, return it otherwise, attempt reconnect
	if c.wsconn != nil {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Connection is not valid\n")
			klog.V(7).Infof("common.internalConnectWithCancel() LEAVE\n")
			if lock {
				klog.V(3).Infof("Unlocking connection mutex\n")
				c.muConn.Unlock()
			}
			return nil
		default:
			klog.V(7).Infof("Connection is good. Return object.")
			klog.V(7).Infof("common.internalConnectWithCancel() LEAVE\n")
			if lock {
				klog.V(3).Infof("Unlocking connection mutex\n")
				c.muConn.Unlock()
			}
			return c.wsconn
		}
	} else {
		select {
		case <-c.ctx.Done():
			klog.V(1).Infof("Context is not valid. Has been canceled.\n")
			klog.V(7).Infof("common.internalConnectWithCancel() LEAVE\n")
			if lock {
				klog.V(3).Infof("Unlocking connection mutex\n")
				c.muConn.Unlock()
			}
			return nil
		default:
			klog.V(3).Infof("Context is still valid. Retry...\n")
		}
	}

	// set websocket headers
	myHeader := http.Header{}

	// restore application options to HTTP header
	if headers, ok := c.ctx.Value(clientinterfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				klog.V(3).Infof("internalConnectWithCancel RESTORE Header: %s = %s\n", k, v)
				myHeader.Add(k, v)
			}
		}
	}

	// Set Authorization header based on priority: AccessToken (Bearer) > APIKey (Token)
	myHeader.Set("Host", c.cOptions.Host)
	token, isBearer := c.cOptions.GetAuthToken()
	if isBearer {
		myHeader.Set("Authorization", "Bearer "+token)
		klog.V(4).Infof("WebSocket using Bearer authentication")
	} else {
		myHeader.Set("Authorization", "token "+token)
		klog.V(4).Infof("WebSocket using Token authentication")
	}
	myHeader.Set("User-Agent", clientinterfaces.DgAgent)
	if c.cOptions.WSHeaderProcessor != nil {
		c.cOptions.WSHeaderProcessor(myHeader)
	}

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
		url, err := (*c.processMessages).GetURL(c.cOptions.Host)
		if err != nil {
			klog.V(1).Infof("GetURL failed. Err: %v\n", err)
			klog.V(7).Infof("internalConnectWithCancel() LEAVE\n")
			if lock {
				klog.V(3).Infof("Unlocking connection mutex\n")
				c.muConn.Unlock()
			}
			return nil // no point in retrying because this is going to fail on every retry
		}
		klog.V(5).Infof("Connecting to %s\n", url)

		// if host starts with "ws://", then disable TLS
		var dialer websocket.Dialer
		if url[:5] == "ws://" {
			dialer = websocket.Dialer{
				HandshakeTimeout: 15 * time.Second,
				RedirectService:  c.cOptions.RedirectService,
				Proxy:            c.cOptions.Proxy,
			}
		} else {
			dialer = websocket.Dialer{
				HandshakeTimeout: 15 * time.Second,
				/* #nosec G402 */
				TLSClientConfig: &tls.Config{InsecureSkipVerify: c.cOptions.SkipServerAuth},
				RedirectService: c.cOptions.RedirectService,
				SkipServerAuth:  c.cOptions.SkipServerAuth,
				Proxy:           c.cOptions.Proxy,
			}
		}

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
		if lock {
			klog.V(3).Infof("Unlocking connection mutex\n")
			c.muConn.Unlock()
		}

		// start WS specific items
		(*c.processMessages).Start()

		// fire off close connection
		err = (*c.router).Open(&commonv1interfaces.OpenResponse{
			Type: string(commonv1interfaces.TypeOpenResponse),
		})
		if err != nil {
			klog.V(1).Infof("router.Open failed. Err: %v\n", err)
		}

		klog.V(3).Infof("WebSocket Connection Successful!")
		klog.V(7).Infof("common.internalConnectWithCancel() LEAVE\n")

		return c.wsconn
	}

	// if we get here, we failed to connect
	klog.V(1).Infof("Failed to connect to websocket: %s\n", c.cOptions.Host)
	klog.V(7).Infof("common.internalConnectWithCancel() LEAVE\n")

	if lock {
		klog.V(3).Infof("Unlocking connection mutex\n")
		c.muConn.Unlock()
	}

	return nil
}

//nolint:funlen // this is a complex function. keep as is
func (c *WSClient) listen() {
	klog.V(6).Infof("common.listen() ENTER\n")

	defer func() {
		if r := recover(); r != nil {
			klog.V(1).Infof("Panic triggered\n")
			klog.V(1).Infof("Panic: %v\n", r)
			klog.V(1).Infof("Stack trace: %s\n", string(debug.Stack()))

			// send error on callback
			err := ErrFatalPanicRecovered
			sendErr := c.sendError(err)
			if sendErr != nil {
				klog.V(1).Infof("listen: Fatal socket error. Err: %v\n", sendErr)
			}

			// fatal close
			c.closeWs(true, false)

			klog.V(6).Infof("common.listen() LEAVE\n")
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
			klog.V(6).Infof("common.listen() LEAVE\n")
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

				klog.V(6).Infof("common.listen() LEAVE\n")
				return
			case strings.Contains(errStr, UseOfClosedSocket):
				klog.V(3).Infof("Probable graceful websocket close: %v\n", err)

				// fatal close
				c.closeWs(false, false)

				klog.V(6).Infof("common.listen() LEAVE\n")
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

				klog.V(6).Infof("common.listen() LEAVE\n")
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

				klog.V(6).Infof("common.listen() LEAVE\n")
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

				klog.V(6).Infof("common.listen() LEAVE\n")
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

				klog.V(6).Infof("common.listen() LEAVE\n")
				return
			}
		}

		if len(byMsg) == 0 {
			klog.V(7).Infof("listen(): message empty")
			continue
		}

		// process WS specific message
		err = (*c.processMessages).ProcessMessage(msgType, byMsg)
		if err != nil {
			klog.V(1).Infof("ProcessMessage failed. Err: %v\n", err)
		}
	}
}

// WriteBinary writes binary data to the websocket server
func (c *WSClient) WriteBinary(byData []byte) error {
	klog.V(7).Infof("common.WriteBinary() ENTER\n")

	// doing a write, need to lock
	c.muConn.Lock()
	defer c.muConn.Unlock()

	// get the connection
	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(7).Infof("common.WriteBinary() LEAVE\n")

		return err
	}

	if err := ws.WriteMessage(
		websocket.BinaryMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteBinary WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("common.WriteBinary() LEAVE\n")
		return err
	}

	klog.V(7).Infof("WriteBinary Successful\n")
	klog.V(7).Infof("payload: %x\n", byData)
	klog.V(7).Infof("common.WriteBinary() LEAVE\n")

	return nil
}

/*
WriteJSON writes a JSON control payload to the websocket server. These are control messages for
managing the websocket connection.
*/
func (c *WSClient) WriteJSON(payload interface{}) error {
	klog.V(6).Infof("common.WriteJSON() ENTER\n")

	byData, err := json.Marshal(payload)
	if err != nil {
		klog.V(1).Infof("WriteJSON: Error marshaling JSON. Data: %v, Err: %v\n", payload, err)
		klog.V(6).Infof("common.WriteJSON() LEAVE\n")
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
		klog.V(6).Infof("common.WriteJSON() LEAVE\n")

		return err
	}

	if err := ws.WriteMessage(
		websocket.TextMessage,
		byData,
	); err != nil {
		klog.V(1).Infof("WriteJSON WriteMessage failed. Err: %v\n", err)
		klog.V(6).Infof("common.WriteJSON() LEAVE\n")
		return err
	}

	klog.V(4).Infof("common.WriteJSON() Succeeded\n")
	klog.V(6).Infof("payload: %s\n", string(byData))
	klog.V(6).Infof("common.WriteJSON() LEAVE\n")

	return nil
}

// closeStream sends an application level message to Deepgram
func (c *WSClient) closeStream(lock bool) error {
	klog.V(7).Infof("common.closeStream() ENTER\n")

	// doing a write, need to lock
	if lock {
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	var err error
	byClose := (*c.processMessages).GetCloseMsg()
	if len(byClose) > 0 {
		klog.V(3).Infof("closeStream: Sending close message\n")
		err = c.wsconn.WriteMessage(websocket.TextMessage, byClose)
	} else {
		klog.V(3).Infof("closeStream: No protocol specific close message\n")
	}

	if err != nil {
		klog.V(1).Infof("WriteMessage failed. Err: %v\n", err)
		klog.V(7).Infof("common.closeStream() LEAVE\n")

		return err
	}

	klog.V(4).Infof("closeStream Succeeded\n")
	klog.V(7).Infof("common.closeStream() LEAVE\n")

	return err
}

// normalClosure sends a normal closure message to the server
func (c *WSClient) normalClosure(lock bool) error {
	klog.V(7).Infof("common.normalClosure() ENTER\n")

	// doing a write, need to lock
	if lock {
		c.muConn.Lock()
		defer c.muConn.Unlock()
	}

	ws := c.internalConnect()
	if ws == nil {
		err := ErrInvalidConnection
		klog.V(4).Infof("c.internalConnect() is nil. Err: %v\n", err)
		klog.V(7).Infof("common.normalClosure() LEAVE\n")

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

	klog.V(7).Infof("common.normalClosure() LEAVE\n")

	return err
}

// Stop will send close message and shutdown websocket connection
func (c *WSClient) Stop() {
	klog.V(3).Infof("Stopping...\n")
	c.muConn.Lock()
	c.retry = false
	c.muConn.Unlock()

	// exit gracefully
	c.closeWs(false, true)
}

// closeWs closes the websocket connection
func (c *WSClient) closeWs(fatal bool, perm bool) {
	klog.V(6).Infof("common.closeWs() closing channels...\n")

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
		// process WS specific items
		(*c.processMessages).Finish()

		// fire off close connection
		err := (*c.router).Close(&commonv1interfaces.CloseResponse{
			Type: string(commonv1interfaces.TypeCloseResponse),
		})
		if err != nil {
			klog.V(1).Infof("router.CloseHelper failed. Err: %v\n", err)
		}
	}

	// cancel the context because we are permanently closing the connection
	if perm {
		klog.V(3).Infof("Permanently closing connection\n")
		c.ctxCancel()
	}

	// close the connection
	if c.wsconn != nil {
		c.wsconn.Close()
		c.wsconn = nil
	}

	klog.V(4).Infof("common.closeWs() Succeeded\n")
	klog.V(6).Infof("common.closeWs() LEAVE\n")
}

// sendError sends an error message to the callback handler
func (c *WSClient) sendError(err error) error {
	sendErr := (*c.processMessages).ProcessError(err)
	if err != nil {
		klog.V(1).Infof("ProcessError(%v) failed. Err: %v\n", err, sendErr)
	}

	return err
}
