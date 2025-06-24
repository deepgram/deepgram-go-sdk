// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"context"
	"sync"

	"github.com/dvonthenen/websocket"

	commonv1interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	clientinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
	restv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/rest/v1"
)

const (
	PackageVersion string = "v1.0"
)

// ***************************
// Common WS Client
// ***************************
// WSClient is a struct representing the websocket client connection
type WSClient struct {
	cOptions *clientinterfaces.ClientOptions

	sendBuf   chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc

	muConn   sync.RWMutex
	wsconn   *websocket.Conn
	retry    bool
	retryCnt int64

	processMessages *commonv1interfaces.WebSocketHandler
	router          *commonv1interfaces.Router
}

// ***************************
// Common REST Client
// ***************************
// RESTClient implements an extensible REST client
type RESTClient struct {
	*restv1.Client
}
