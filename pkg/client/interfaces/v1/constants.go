// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"errors"
)

const (
	PackageVersion string = "v1.0"
)

const (
	TypeSettings = "Settings"
)

// errors
var (
	// ErrNoAPIKey no api key found
	ErrNoAPIKey = errors.New("no api key found")
)
