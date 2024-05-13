// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetQueryParams(url string) string {
	pos := strings.Index(url, "?")
	if pos == -1 {
		return ""
	}
	return url[pos+1:]
}

func CreateDirs(fullpath string) error {
	basedir := filepath.Dir(fullpath)
	return os.MkdirAll(basedir, 0o700)
}

func SaveMetadataBytes(filename string, data []byte) error {
	return SaveMetadataString(filename, string(data))
}

func SaveMetadataString(filename, data string) error {
	// create directory
	err := CreateDirs(filename)
	if err != nil {
		return err
	}

	// save metadata
	dataFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o700)
	if err != nil {
		return err
	}
	_, err = dataFile.WriteString(data)
	if err != nil {
		return err
	}
	dataFile.Close()

	return err
}

func ReadMetadataString(filename string) (string, error) {
	byData, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(byData), err
}

func ReadMetadataBytes(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func StringMatchFailure(expected, actual string) error {
	if len(expected) != len(actual) {
		return errors.New("string lengths don't match")
	}

	found := -1
	for i := 0; i < len(expected); i++ {
		expectedPos := expected[i:1]
		actualPos := actual[i:1]
		if expectedPos != actualPos {
			found = i
			break
		}
	}

	// expected
	for i := 0; i < len(expected); i++ {
		expectedPos := expected[i:1]
		if i == found {
			fmt.Fprintf(os.Stdout, "\033[0;31m %s", expectedPos)
		} else {
			fmt.Fprintf(os.Stdout, "\033[0m %s", expectedPos)
		}
	}
	fmt.Printf("\n")

	// actual
	for i := 0; i < len(expected); i++ {
		actualPos := actual[i:1]
		if i == found {
			fmt.Fprintf(os.Stdout, "\033[0;31m %s", actualPos)
		} else {
			fmt.Fprintf(os.Stdout, "\033[0m %s", actualPos)
		}
	}
	fmt.Printf("\n")

	return nil
}
