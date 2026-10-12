//go:build !windows

package main

import "os"

// replaceOutput atomically replaces target with temporary when both paths are
// in the same directory, which newTemporaryOutput guarantees.
func replaceOutput(temporary, target string) error {
	return os.Rename(temporary, target)
}
