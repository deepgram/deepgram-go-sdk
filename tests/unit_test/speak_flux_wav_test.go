// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"

	gowav "github.com/youpy/go-wav"

	wav "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/wav"
)

// Test_WAVWriterDeclaresRealSizes is the parser-based regression test for the
// streaming-WAV size fields (Greg's B3 on PR #351): after Finalize, the RIFF size
// at byte offset 4 is fileSize-8 and the data length at offset 40 is the exact
// PCM byte count — not the zero placeholders a streamed file starts with — and a
// conforming WAV parser accepts the file.
func Test_WAVWriterDeclaresRealSizes(t *testing.T) {
	const (
		channels      = 1
		sampleRate    = 48000
		bitsPerSample = 16
		headerSize    = 44
	)

	filename := filepath.Join(t.TempDir(), "output.wav")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatalf("creating output file failed: %s", err)
	}
	defer file.Close()

	writer, err := wav.NewWriter(file, channels, sampleRate, bitsPerSample)
	if err != nil {
		t.Fatalf("NewWriter failed: %s", err)
	}

	// stream PCM in several chunks, as the websocket delivers it
	var pcmBytes int
	for _, chunkLen := range []int{3200, 6400, 1600} {
		chunk := bytes.Repeat([]byte{0x12, 0x34}, chunkLen/2)
		n, err := writer.Write(chunk)
		if err != nil {
			t.Fatalf("Write failed: %s", err)
		}
		pcmBytes += n
	}

	if err := writer.Finalize(); err != nil {
		t.Fatalf("Finalize failed: %s", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Close failed: %s", err)
	}

	raw, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading finished file failed: %s", err)
	}
	if len(raw) != headerSize+pcmBytes {
		t.Fatalf("expected %d bytes on disk, got %d", headerSize+pcmBytes, len(raw))
	}

	// the two fields the finding is about, at their exact byte offsets
	if got, want := binary.LittleEndian.Uint32(raw[4:8]), uint32(len(raw)-8); got != want {
		t.Errorf("RIFF size at offset 4: got %d, want fileSize-8 = %d", got, want)
	}
	if got, want := binary.LittleEndian.Uint32(raw[40:44]), uint32(pcmBytes); got != want {
		t.Errorf("data size at offset 40: got %d, want exact PCM length %d", got, want)
	}

	// a conforming WAV parser must accept the finished file and report the format
	reader := gowav.NewReader(bytes.NewReader(raw))
	format, err := reader.Format()
	if err != nil {
		t.Fatalf("WAV parser rejected the file: %s", err)
	}
	if format.AudioFormat != 1 || format.NumChannels != channels ||
		format.SampleRate != sampleRate || format.BitsPerSample != bitsPerSample {
		t.Errorf("unexpected parsed format: %+v", format)
	}

	// the declared data chunk yields exactly pcmBytes of samples
	sampleCount := 0
	for {
		samples, err := reader.ReadSamples(1024)
		sampleCount += len(samples)
		if err != nil {
			break
		}
	}
	if wantSamples := pcmBytes / (channels * bitsPerSample / 8); sampleCount != wantSamples {
		t.Errorf("parser read %d samples, want %d", sampleCount, wantSamples)
	}

	t.Run("Write after Finalize is rejected", func(t *testing.T) {
		if _, err := writer.Write([]byte{0x00, 0x00}); err == nil {
			t.Error("expected an error writing after Finalize")
		}
	})
}

func Test_WAVWriterRejectsInvalidFormat(t *testing.T) {
	for _, test := range []struct {
		name          string
		channels      uint16
		sampleRate    uint32
		bitsPerSample uint16
	}{
		{"zero channels", 0, 48000, 16},
		{"zero sample rate", 1, 0, 16},
		{"zero bit depth", 1, 48000, 0},
		{"non-byte-aligned bit depth", 1, 48000, 12},
		{"block alignment overflow", 65535, 48000, 16},
		{"byte rate overflow", 1, ^uint32(0), 16},
	} {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.CreateTemp(t.TempDir(), "invalid.wav")
			if err != nil {
				t.Fatalf("creating output file failed: %s", err)
			}
			defer file.Close()

			if _, err := wav.NewWriter(file, test.channels, test.sampleRate, test.bitsPerSample); !errors.Is(err, wav.ErrInvalidFormat) {
				t.Errorf("expected ErrInvalidFormat, got: %v", err)
			}
		})
	}
}
