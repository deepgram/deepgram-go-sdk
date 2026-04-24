// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv2

import (
	"context"
	"time"

	"github.com/dvonthenen/websocket"
	klog "k8s.io/klog/v2"

	commonv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
	commonv2interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
)

// gocritic:ignore
// NewWS creates a new v2 WSClient.
func NewWS(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, options *clientinterfaces.ClientOptions, processMessages *commonv2interfaces.WebSocketHandler, router *commonv2interfaces.Router) *WSClient {
	v1client := commonv1.NewWS(ctx, ctxCancel, apiKey, options, processMessages, router)
	if v1client == nil {
		return nil
	}
	return &WSClient{WSClient: v1client}
}

// WritePing sends a WebSocket protocol-level ping control frame to the server.
// The server should respond with a pong frame automatically.
// WritePing is only available in v2.
func (c *WSClient) WritePing() error {
	klog.V(7).Infof("commonv2.WritePing() ENTER\n")

	deadline := time.Now().Add(10 * time.Second)
	if err := c.WSClient.WriteControl(websocket.PingMessage, []byte{}, deadline); err != nil {
		klog.V(1).Infof("WritePing failed. Err: %v\n", err)
		klog.V(7).Infof("commonv2.WritePing() LEAVE\n")
		return err
	}

	klog.V(5).Infof("WritePing succeeded\n")
	klog.V(7).Infof("commonv2.WritePing() LEAVE\n")
	return nil
}
