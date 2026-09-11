package main

import (
	"os"
	"path/filepath"
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

func TestTemporaryOutputPreservesOrReplacesDestination(t *testing.T) {
	destination := filepath.Join(t.TempDir(), audioFile)
	const existingAudio = "existing audio"
	if err := os.WriteFile(destination, []byte(existingAudio), 0o600); err != nil {
		t.Fatalf("creating existing output file: %v", err)
	}

	t.Run("discard preserves existing output", func(t *testing.T) {
		output, err := newTemporaryOutput(destination)
		if err != nil {
			t.Fatalf("creating temporary output: %v", err)
		}
		if _, err := output.file.Write([]byte("partial audio")); err != nil {
			t.Fatalf("writing temporary output: %v", err)
		}
		output.discard()

		contents, err := os.ReadFile(destination)
		if err != nil {
			t.Fatalf("reading existing output: %v", err)
		}
		if string(contents) != existingAudio {
			t.Errorf("existing output was modified: got %q, want %q", contents, existingAudio)
		}
	})

	t.Run("commit replaces existing output", func(t *testing.T) {
		output, err := newTemporaryOutput(destination)
		if err != nil {
			t.Fatalf("creating temporary output: %v", err)
		}
		if _, err := output.file.Write([]byte("replacement audio")); err != nil {
			t.Fatalf("writing temporary output: %v", err)
		}
		if err := output.commit(); err != nil {
			t.Fatalf("committing temporary output: %v", err)
		}

		contents, err := os.ReadFile(destination)
		if err != nil {
			t.Fatalf("reading replacement output: %v", err)
		}
		if string(contents) != "replacement audio" {
			t.Errorf("output was not replaced: got %q", contents)
		}
	})
}
