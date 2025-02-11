// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"context"

	msginterface "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	commoninterfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// client messages
type SettingsConfigurationOptions interfaces.SettingsConfigurationOptions
type UpdateInstructions msginterface.UpdateInstructions
type UpdateSpeak msginterface.UpdateSpeak
type InjectAgentMessage msginterface.InjectAgentMessage
type FunctionCallResponse msginterface.FunctionCallResponse
type KeepAlive msginterface.KeepAlive

// WSChannel is a struct representing the websocket client connection using channels
type WSChannel struct {
	*common.WSClient
	ctx       context.Context
	ctxCancel context.CancelFunc

	cOptions *interfaces.ClientOptions
	tOptions *interfaces.SettingsConfigurationOptions

	chans  []*msginterface.AgentMessageChan
	router *commoninterfaces.Router
}
