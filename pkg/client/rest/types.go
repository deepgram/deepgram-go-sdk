// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package implements a reusable REST client
*/
package rest

import (
	"net/http"
)

// ClientOptions defines any options for the client
type ClientOptions struct {
	Host    string
	Version string
}

// HttpClient which extends HTTP client
type HttpClient struct {
	http.Client

	d         *debugContainer
	UserAgent string
}

// Client which extends HttpClient to support REST
type Client struct {
	*HttpClient

	Options *ClientOptions
	apiKey  string
}
