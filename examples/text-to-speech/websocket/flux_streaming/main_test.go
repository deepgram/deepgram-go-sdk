package main

import (
	"os"
	"testing"
)

func TestRunWithoutCredentialsLeavesNoOutput(t *testing.T) {
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

	if err := run(); err == nil {
		t.Fatal("expected client construction to fail without credentials")
	}
	if _, err := os.Stat(audioFile); !os.IsNotExist(err) {
		t.Errorf("expected no partial output file, got stat error %v", err)
	}
}
