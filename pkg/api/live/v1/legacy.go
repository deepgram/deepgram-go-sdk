// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// *********** WARNING ***********
// This package provides the Live API
//
// Deprecated: This package is deprecated. Use the listen package instead. This will be removed in a future release.
//
// This package is frozen and no new functionality will be added.
// *********** WARNING ***********
package legacy

import (
	websocketv1 "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
)

const (
	PackageVersion = websocketv1.PackageVersion
)

// Alias
type LiveMessageCallback = interfacesv1.LiveMessageCallback
type DefaultCallbackHandler = websocketv1.DefaultCallbackHandler
type MessageRouter = websocketv1.MessageRouter

// NewDefaultCallbackHandler
func NewDefaultCallbackHandler() websocketv1.DefaultCallbackHandler {
	return DefaultCallbackHandler{}
}

// MessageRouter
func NewWithDefault() *websocketv1.MessageRouter {
	return websocketv1.NewWithDefault()
}

// New creates a MessageRouter with a user-defined callback
func New(callback LiveMessageCallback) *websocketv1.MessageRouter {
	return websocketv1.New(callback)
}
