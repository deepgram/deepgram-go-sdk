// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"net/http"
	"net/url"
)

// ClientOptions defines any options for the client
type ClientOptions struct {
	APIKey            string
	Host              string                                // override for the host endpoint
	APIVersion        string                                // override for the version used
	Path              string                                // override for the endpoint path usually <version/listen> or <version/projects>
	SelfHosted        bool                                  // set to true if using on-prem
	Proxy             func(*http.Request) (*url.URL, error) // provide function for proxy -- e.g. http.ProxyFromEnvironment
	WSHeaderProcessor func(http.Header)                     // process headers before dialing for websocket connection

	// shared client options
	SkipServerAuth bool // keeps the client from authenticating with the server

	// prerecorded client options

	// speech-to-text client options
	RedirectService     bool  // allows HTTP redirects to be followed
	EnableKeepAlive     bool  // enables the keep alive feature
	AutoFlushReplyDelta int64 // enables the auto flush feature based on the delta in milliseconds

	// text-to-speech client options
	AutoFlushSpeakDelta int64 // enables the auto flush feature based on the delta in milliseconds
}
