package main

import (
	"os"
	"testing"
)

func TestRunWithoutCredentialsPreservesExistingOutput(t *testing.T) {
	t.Setenv("DEEPGRAM_API_KEY", "")
	t.Setenv("DEEPGRAM_ACCESS_TOKEN", "")

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Errorf("restoring working directory: %v", err)
		}
	})
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("switching to temporary directory: %v", err)
	}
	const existingAudio = "existing audio"
	if err := os.WriteFile(audioFile, []byte(existingAudio), 0o600); err != nil {
		t.Fatalf("creating existing output file: %v", err)
	}

	if err := run(); err == nil {
		t.Fatal("expected client construction to fail without credentials")
	}
	output, err := os.ReadFile(audioFile)
	if err != nil {
		t.Fatalf("reading existing output file: %v", err)
	}
	if string(output) != existingAudio {
		t.Errorf("existing output was modified: got %q, want %q", output, existingAudio)
	}
}
