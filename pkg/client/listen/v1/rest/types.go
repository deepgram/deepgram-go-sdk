// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
)

// RESTClient implements helper functionality for Prerecorded API
type RESTClient struct {
	*common.RESTClient
}

// Client is an alias for WSCallback
type Client = RESTClient
