// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
)

// Client implements helper functionality for Prerecorded API
type RESTClient struct {
	*common.RESTClient
}

// Client is an alias for RESTClient
type Client = RESTClient
