// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"

	msginterface "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// client messages
type SettingsOptions interfaces.SettingsOptions
type UpdatePrompt msginterface.UpdatePrompt
type UpdateSpeak msginterface.UpdateSpeak
type InjectAgentMessage msginterface.InjectAgentMessage
type InjectUserMessage msginterface.InjectUserMessage
type FunctionCallResponse msginterface.FunctionCallResponse
type KeepAlive msginterface.KeepAlive

// WSChannel is a struct representing the websocket client connection using channels
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptions
	tOptions *interfaces.SettingsOptions

	chans  []*msginterface.AgentMessageChan
	router *commoninterfaces.Router
}
