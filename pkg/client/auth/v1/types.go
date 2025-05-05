// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package contains the code for the Token APIs in the Deepgram Auth API
package authv1

import (
	common "github.com/deepgram/deepgram-go-sdk/v2/pkg/client/common/v1"
)

// Client is the client for the Deepgram Auth API
type Client struct {
	*common.RESTClient
}
