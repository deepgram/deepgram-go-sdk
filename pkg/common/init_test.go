// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package common

import (
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test_InitLeavesApplicationFlagsAlone(t *testing.T) {
	if os.Getenv("DEEPGRAM_COMMON_INIT_HELPER") == "1" {
		applicationFlags := flag.NewFlagSet("application", flag.ContinueOnError)
		flag.CommandLine = applicationFlags

		Init(InitLib{LogLevel: LogLevelStandard})
		if applicationFlags.Parsed() {
			t.Fatal("Init must not parse application flags")
		}
		if applicationFlags.Lookup("v") != nil {
			t.Fatal("Init must not register klog flags on the application flag set")
		}

		custom := applicationFlags.String("custom", "", "application flag")
		if err := applicationFlags.Parse([]string{"-custom=expected"}); err != nil {
			t.Fatalf("parse application flags: %v", err)
		}
		if *custom != "expected" {
			t.Fatalf("expected custom flag value %q, got %q", "expected", *custom)
		}

		// klog supports distinct flag sets, so repeated SDK initialization must also
		// leave the consumer's global flag set untouched.
		Init(InitLib{LogLevel: LogLevelFull})
		if applicationFlags.Lookup("v") != nil {
			t.Fatal("repeated Init must not register klog flags on the application flag set")
		}
		return
	}

	command := exec.Command(os.Args[0], "-test.run=^Test_InitLeavesApplicationFlagsAlone$")
	command.Env = append(os.Environ(), "DEEPGRAM_COMMON_INIT_HELPER=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("Init modified application flags:\n%s\n%v", output, err)
	}
	if strings.Contains(string(output), "flag provided but not defined") {
		t.Fatalf("Init rejected a test runner flag:\n%s", output)
	}
}
