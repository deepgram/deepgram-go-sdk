// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package common

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var helperApplicationFlag *string
var helperFlagsParsed bool

func init() {
	if os.Getenv("DEEPGRAM_COMMON_INIT_HELPER") != "1" {
		return
	}

	// This matches a consumer defining flags and initializing the SDK during
	// package initialization, before the Go test runner defines its own flags.
	helperApplicationFlag = flag.String("custom", "", "application flag")
	Init(InitLib{LogLevel: LogLevelStandard})
	helperFlagsParsed = flag.Parsed()
}

func Test_InitDefersGoTestFlagParsing(t *testing.T) {
	if os.Getenv("DEEPGRAM_COMMON_INIT_HELPER") == "1" {
		if helperApplicationFlag == nil {
			t.Fatal("helper application flag was not initialized")
		}
		if helperFlagsParsed {
			t.Fatal("Init must not parse flags before the Go test runner initializes them")
		}
		if *helperApplicationFlag != "expected" {
			t.Fatalf("expected custom flag value %q, got %q", "expected", *helperApplicationFlag)
		}
		verbosity := flag.CommandLine.Lookup("v")
		if verbosity == nil {
			t.Fatal("Init must preserve klog flags on the application flag set")
		}
		if got := verbosity.Value.String(); got != "2" {
			t.Fatalf("expected klog verbosity %q, got %q", "2", got)
		}
		return
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("find test executable: %v", err)
	}

	renamedExecutable := filepath.Join(t.TempDir(), "common-test-helper")
	executableContents, err := os.ReadFile(executable)
	if err != nil {
		t.Fatalf("read test executable: %v", err)
	}
	if err := os.WriteFile(renamedExecutable, executableContents, 0o755); err != nil {
		t.Fatalf("write renamed test executable: %v", err)
	}

	command := exec.Command(renamedExecutable, "--test.run=^Test_InitDefersGoTestFlagParsing$", "-custom=expected", "-v=2")
	command.Env = append(os.Environ(), "DEEPGRAM_COMMON_INIT_HELPER=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("Init modified application flags:\n%s\n%v", output, err)
	}
	if strings.Contains(string(output), "flag provided but not defined") {
		t.Fatalf("Init rejected a test runner flag:\n%s", output)
	}
}
