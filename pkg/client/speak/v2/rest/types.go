// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv2

import (
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
)

// Client implements the Flux TTS batch (REST) transport for POST /v2/speak.
type Client struct {
	*common.RESTClient
}
