// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Test_ModuleMajorScriptReadsGofmtAlignedSDKVersion(t *testing.T) {
	versionFile := filepath.Join(t.TempDir(), "utils.go")
	if err := os.WriteFile(versionFile, []byte("const (\n\totherConstant string = \"value\"\n\tsdkVersion    string = \"v4.0.0\"\n)\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	command := exec.Command("bash", "-c", "source hack/check/check-module-major.sh; read_sdk_major \"$1\"", "bash", versionFile)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("read_sdk_major failed: %v\n%s", err, output)
	}

	if actual := strings.TrimSpace(string(output)); actual != "4" {
		t.Errorf("read_sdk_major() = %q, want %q", actual, "4")
	}
}
