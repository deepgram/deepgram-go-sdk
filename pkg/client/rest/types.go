// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// HttpClient which extends HTTP client
type HttpClient struct {
	http.Client

	d         *debugContainer
	UserAgent string

	options interfaces.ClientOptions
}

// Client which extends HttpClient to support REST
type Client struct {
	*HttpClient

	Options interfaces.ClientOptions
}
