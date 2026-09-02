// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv2

import "errors"

const (
	PackageVersion string = "v2.0"
)

var (
	// ErrNilClient the transport client is nil (its constructor returned an error)
	ErrNilClient = errors.New("speak client is nil — check the error returned by the REST client constructor")

	// ErrCallbackNotSupported ToSave was called with a callback URL set
	ErrCallbackNotSupported = errors.New("ToSave synthesizes audio synchronously and does not support callback requests; use ToAsync for asynchronous synthesis")

	// ErrCallbackRequired ToAsync was called without a callback URL
	ErrCallbackRequired = errors.New("ToAsync requires options.Callback to be set")
)
