// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

// ClientOptions defines any options for the client
type ClientOptions struct {
	APIKey     string
	Host       string // override for the host endpoint
	APIVersion string // override for the version used
	Path       string // override for the endpoint path usually <version/listen> or <version/projects>
	SelfHosted bool   // set to true if using on-prem

	// shared client options
	SkipServerAuth bool // keeps the client from authenticating with the server

	// prerecorded client options

	// live client options
	RedirectService bool // allows HTTP redirects to be followed
	EnableKeepAlive bool // enables the keep alive feature

	// these require inspecting messages, therefore you must update the InspectMessage() method
	AutoFlushReplyDelta int64 // enables the auto flush feature
}
