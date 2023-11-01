// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the initialization code for the Deepgram Go SDK
*/
package live

import (
	common "github.com/deepgram-devs/deepgram-go-sdk/pkg/common"
)

// InitWithDefault is the SDK Init function for this library using default values.
func InitWithDefault() {
	Init(common.InitLib{
		LogLevel: common.LogLevelDefault,
	})
}

// The SDK Init function for this library.
// Allows you to set the logging level and use of a log file.
// Default is output to the stdout.
func Init(init common.InitLib) {
	common.Init(init)
}