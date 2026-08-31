// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"context"

	msginterface "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v2/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// WSCallback is a Flux TTS WebSocket client that delivers server events via a
// FluxSpeakMessageCallback.
type WSCallback struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptionsV2
	sOptions *interfaces.SpeakV2WSOptions

	callback msginterface.FluxSpeakMessageCallback
	router   *commoninterfaces.Router
}

// WSChannel is a Flux TTS WebSocket client that delivers server events via Go channels.
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptionsV2
	sOptions *interfaces.SpeakV2WSOptions

	chans  []*msginterface.FluxSpeakMessageChan
	router *commoninterfaces.Router
}
