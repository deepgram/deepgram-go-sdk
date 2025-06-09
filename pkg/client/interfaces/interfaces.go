// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const (
	PackageVersion = interfacesv1.PackageVersion
)

// NewSettingsConfigurationOptions creates a new SettingsConfigurationOptions object
func NewSettingsConfigurationOptions() *interfacesv1.SettingsOptions {
	return interfacesv1.NewSettingsOptions()
}

// options
type ClientOptions = interfacesv1.ClientOptions
type SettingsOptions = interfacesv1.SettingsOptions
type PreRecordedTranscriptionOptions = interfacesv1.PreRecordedTranscriptionOptions
type LiveTranscriptionOptions = interfacesv1.LiveTranscriptionOptions
type AnalyzeOptions = interfacesv1.AnalyzeOptions
type SpeakOptions = interfacesv1.SpeakOptions
type WSSpeakOptions = interfacesv1.WSSpeakOptions
