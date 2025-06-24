// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"net/http"
	"net/url"
	"sync"
)

// ClientOptions defines any options for the client
type ClientOptions struct {
	APIKey            string
	AccessToken       string                                // JWT access token for Bearer authentication
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

	// Thread safety for credential management
	credentialsMutex sync.RWMutex // protects AccessToken and APIKey fields
}

// SetAccessToken dynamically sets the access token for Bearer authentication (thread-safe)
func (o *ClientOptions) SetAccessToken(accessToken string) {
	o.credentialsMutex.Lock()
	defer o.credentialsMutex.Unlock()
	o.AccessToken = accessToken
}

// SetAPIKey dynamically sets the API key for Token authentication (thread-safe)
func (o *ClientOptions) SetAPIKey(apiKey string) {
	o.credentialsMutex.Lock()
	defer o.credentialsMutex.Unlock()
	o.APIKey = apiKey
}

// GetAuthToken returns the effective authentication token following priority order (thread-safe):
// 1. AccessToken (Bearer) - highest priority
// 2. APIKey (Token) - fallback
func (o *ClientOptions) GetAuthToken() (token string, isBearer bool) {
	o.credentialsMutex.RLock()
	defer o.credentialsMutex.RUnlock()
	if o.AccessToken != "" {
		return o.AccessToken, true
	}
	return o.APIKey, false
}
