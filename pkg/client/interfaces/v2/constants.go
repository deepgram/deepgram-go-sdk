// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv2

import (
	"errors"

	v1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const (
	PackageVersion string = "v2.0"
)

const (
	TypeSettings = v1.TypeSettings
)

// errors
var (
	// ErrNoAPIKey no api key found
	ErrNoAPIKey = v1.ErrNoAPIKey

	// ErrModelRequired no model specified where one is required
	ErrModelRequired = errors.New("model is required")

	// ErrOptionsRequired a nil options struct was passed where one is required
	ErrOptionsRequired = errors.New("options cannot be nil")

	// ErrConfigureNoFields a Configure carried no fields to change
	ErrConfigureNoFields = errors.New("configure requires at least one field to change (set Speed)")

	// ErrSpeedOutOfRange the speed value is outside the supported set
	ErrSpeedOutOfRange = errors.New("speed must be between 0.5 and 1.5 in 0.05 increments")
)
