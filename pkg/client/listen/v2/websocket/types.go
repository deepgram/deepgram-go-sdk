// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"

	msginterface "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// WSCallback is a Flux WebSocket client that delivers server events via a FluxMessageCallback.
type WSCallback struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptionsV2
	tOptions *interfaces.FluxTranscriptionOptions

	callback msginterface.FluxMessageCallback
	router   *commoninterfaces.Router
}

// WSChannel is a Flux WebSocket client that delivers server events via Go channels.
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptionsV2
	tOptions *interfaces.FluxTranscriptionOptions

	chans  []*msginterface.FluxMessageChan
	router *commoninterfaces.Router
}
