// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the initialization code for the Deepgram Go SDK
*/
package common

import (
	"flag"
	"strconv"

	klog "k8s.io/klog/v2"
)

// LogLevel expressed as an int64
type LogLevel int64

// The verbosity of the logging to the console or logfile.
// Default is LogLevelStandard
// LogLevelFull contains INFO related messages that could be helpful in debugging (recommended and default)
// LogLevelTrace is very detailed function enter, highly verbose statements, function exit
// LogLevelVerbose contains data movement on top of LogLevelTrace. This is extremely chatty.
const (
	LogLevelDefault   LogLevel = iota
	LogLevelErrorOnly          = 1
	LogLevelStandard           = 2
	LogLevelElevated           = 3
	LogLevelFull               = 4
	LogLevelDebug              = 5
	LogLevelTrace              = 6
	LogLevelVerbose            = 7
)

/*
Initialization options for this SDK.
*/
type InitLib struct {
	LogLevel      LogLevel
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
	if init.LogLevel == LogLevelDefault {
		init.LogLevel = LogLevelStandard
	}

	klog.InitFlags(nil)
	flag.Set("v", strconv.FormatInt(int64(init.LogLevel), 10))
	if init.DebugFilePath != "" {
		flag.Set("logtostderr", "false")
		flag.Set("log_file", init.DebugFilePath)
	}
	flag.Parse()
}
