// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the initialization code for the Deepgram Go SDK
*/
package common

import (
	"flag"
	"fmt"
	"strconv"

	klog "k8s.io/klog/v2"
)

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
	err := flag.Set("v", strconv.FormatInt(int64(init.LogLevel), 10))
	if err != nil {
		fmt.Printf("Error setting log level: %v", err)
	}
	if init.DebugFilePath != "" {
		err = flag.Set("logtostderr", "false")
		if err != nil {
			fmt.Printf("Error setting logtostderr: %v", err)
		}
		err = flag.Set("log_file", init.DebugFilePath)
		if err != nil {
			fmt.Printf("Error setting log_file: %v", err)
		}
	}
	flag.Parse()
}
