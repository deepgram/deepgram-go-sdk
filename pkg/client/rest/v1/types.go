// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	"net/http"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// HTTPClient which extends HTTP client
type HTTPClient struct {
	http.Client

	d         *debugContainer
	UserAgent string

	options *interfaces.ClientOptions
}

// Client which extends HTTPClient to support REST
type Client struct {
	*HTTPClient

	Options *interfaces.ClientOptions
}
