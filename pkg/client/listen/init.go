// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package listen

import (
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/common"
)

// please see pkg/common/init.go for more information
const (
	LogLevelDefault   = common.LogLevelDefault
	LogLevelErrorOnly = common.LogLevelErrorOnly
	LogLevelStandard  = common.LogLevelStandard
	LogLevelElevated  = common.LogLevelElevated
	LogLevelFull      = common.LogLevelFull
	LogLevelDebug     = common.LogLevelDebug
	LogLevelTrace     = common.LogLevelTrace
	LogLevelVerbose   = common.LogLevelVerbose
)

// Initialization options for this SDK.
type InitLib struct {
	LogLevel      common.LogLevel
	DebugFilePath string
}

// InitWithDefault is the SDK Init function for this library using default values.
func InitWithDefault() {
	Init(InitLib{
		LogLevel: LogLevelDefault,
	})
}

// The SDK Init function for this library.
// Allows you to set the logging level and use of a log file.
// Default is output to the stdout.
func Init(init InitLib) {
	common.Init(common.InitLib{
		LogLevel:      init.LogLevel,
		DebugFilePath: init.DebugFilePath,
	})
}
