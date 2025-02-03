// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package common

const (
	// default host
	DefaultHost string = "api.deepgram.com"

	// default agent host
	DefaultAgentHost string = "agent.deepgram.com"
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
