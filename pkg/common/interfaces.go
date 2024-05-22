// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package common

type Config interface {
	// Parse the config and setup
	Parse() error

	// Getters
	GetHost() string
	GetAPIKey() string
	GetAPIVersion() string
	GetPath() string
	GetSelfHosted() bool

	// prerecorded/manage client options
	GetSkipServerAuth() bool

	// live client options
	GetRedirectService() bool
	GetEnableKeepAlive() bool
}
