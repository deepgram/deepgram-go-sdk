// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"
	"os"
)

type LogWriterCloser struct {
}

func NewLogWriterCloser() *LogWriterCloser {
	return &LogWriterCloser{}
}

func (lwc *LogWriterCloser) Write(p []byte) (n int, err error) {
	fmt.Fprint(os.Stderr, string(Scrub(p)))
	return len(p), nil
}

func (lwc *LogWriterCloser) Close() error {
	return nil
}

type LogProvider struct {
}

func (s LogProvider) NewFile(p string) io.WriteCloser {
	return NewLogWriterCloser()
}

func (s LogProvider) Flush() {
}
