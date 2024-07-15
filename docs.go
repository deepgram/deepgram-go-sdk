// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
	This SDK provides Go library for performing Prerecorded and Live/Streaming operations
	on the Deepgram.com Platform.

	GitHub repo: https://github.com/deepgram/deepgram-go-sdk
	Go SDK Examples: https://github.com/deepgram/deepgram-go-sdk/tree/main/examples

	Deepgram Platform API reference: https://developers.deepgram.com/reference
	Documentation: https://developers.deepgram.com/docs

	The main entry point for this SDK is:
	1. pkg/client/live - contains the SDK objects, functions, etc for performing Live/Stream operations
	2. pkg/client/prerecorded - contains the SDK objects, functions, etc for performing operations on Prerecorded media
*/

package sdk

import (
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/analyze"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/manage"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"

	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/analyze/v1"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/rest"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket"
)
