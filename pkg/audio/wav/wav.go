// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package wav provides a minimal streaming WAV (RIFF/PCM) writer for saving raw
// linear PCM audio — such as the linear16 frames produced by Deepgram's
// text-to-speech WebSocket APIs — as a playable file when the total length is not
// known up front. It supports uncompressed PCM only (see NewWriter).
//
// The writer emits a standard 44-byte header with placeholder sizes, streams PCM
// data as it arrives, and Finalize patches the RIFF size (offset 4) and data chunk
// size (offset 40) so the finished file declares its real lengths, as conforming
// WAV readers require.
package wav

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

// header layout constants
const (
	headerSize       = 44
	riffSizeOffset   = 4  // little-endian uint32: fileSize - 8
	dataSizeOffset   = 40 // little-endian uint32: PCM data byte count
	pcmFormat        = 1  // WAVE_FORMAT_PCM
	fmtChunkSize     = 16
	maxUint32        = int64(^uint32(0))
	maxPCMDataLength = maxUint32 - headerSize + 8 // data must fit the RIFF uint32 size fields
)

// ErrFinalized is returned by Write after Finalize has been called.
var ErrFinalized = errors.New("wav: writer already finalized")

// ErrTooLarge is returned when the PCM data no longer fits the WAV format's
// 32-bit size fields.
var ErrTooLarge = errors.New("wav: PCM data exceeds the maximum WAV file size")

// ErrInvalidFormat is returned when PCM format values cannot be represented in
// a standard WAV header.
var ErrInvalidFormat = errors.New("wav: invalid PCM format")

// Writer streams PCM audio into w as a WAV file. Create it with NewWriter, call
// Write for each audio chunk, and call Finalize once at the end to patch the
// header's size fields. Writer is not safe for concurrent use.
type Writer struct {
	w         io.WriteSeeker
	dataBytes int64
	finalized bool
}

// NewWriter writes a WAV header with placeholder sizes to w and returns a Writer
// that appends PCM data. channels is the channel count (e.g. 1), sampleRate the
// rate in Hz (e.g. 48000), and bitsPerSample the sample width in bits (e.g. 16
// for linear16).
//
// The header always declares WAVE_FORMAT_PCM (format code 1), so the writer is
// for uncompressed linear PCM (linear16) only. Do NOT feed it companded audio:
// mulaw and alaw would need format codes 7 and 6, and a PCM-labeled file of
// companded bytes plays as noise.
func NewWriter(w io.WriteSeeker, channels uint16, sampleRate uint32, bitsPerSample uint16) (*Writer, error) {
	if isNilWriter(w) || channels == 0 || sampleRate == 0 || bitsPerSample == 0 || bitsPerSample%8 != 0 {
		return nil, ErrInvalidFormat
	}

	blockAlign64 := uint64(channels) * uint64(bitsPerSample) / 8
	byteRate64 := uint64(sampleRate) * blockAlign64
	if blockAlign64 > 1<<16-1 || byteRate64 > 1<<32-1 {
		return nil, ErrInvalidFormat
	}
	blockAlign := uint16(blockAlign64)
	byteRate := uint32(byteRate64)

	header := make([]byte, 0, headerSize)
	header = append(header, 'R', 'I', 'F', 'F')
	header = binary.LittleEndian.AppendUint32(header, 0) // RIFF size, patched by Finalize
	header = append(header, 'W', 'A', 'V', 'E', 'f', 'm', 't', ' ')
	header = binary.LittleEndian.AppendUint32(header, fmtChunkSize)
	header = binary.LittleEndian.AppendUint16(header, pcmFormat)
	header = binary.LittleEndian.AppendUint16(header, channels)
	header = binary.LittleEndian.AppendUint32(header, sampleRate)
	header = binary.LittleEndian.AppendUint32(header, byteRate)
	header = binary.LittleEndian.AppendUint16(header, blockAlign)
	header = binary.LittleEndian.AppendUint16(header, bitsPerSample)
	header = append(header, 'd', 'a', 't', 'a')
	header = binary.LittleEndian.AppendUint32(header, 0) // data size, patched by Finalize

	if _, err := w.Write(header); err != nil {
		return nil, err
	}
	return &Writer{w: w}, nil
}

func isNilWriter(w io.WriteSeeker) bool {
	if w == nil {
		return true
	}

	value := reflect.ValueOf(w)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

// Write appends PCM audio bytes to the data chunk.
func (wr *Writer) Write(p []byte) (int, error) {
	if wr.finalized {
		return 0, ErrFinalized
	}
	if wr.dataBytes+int64(len(p)) > maxPCMDataLength {
		return 0, ErrTooLarge
	}
	n, err := wr.w.Write(p)
	wr.dataBytes += int64(n)
	return n, err
}

// DataBytes reports how many PCM bytes have been written so far.
func (wr *Writer) DataBytes() int64 {
	return wr.dataBytes
}

// Finalize seeks back and patches the RIFF size (offset 4, fileSize-8) and data
// chunk size (offset 40, exact PCM byte count), then seeks to the end. It must be
// called exactly once, after the last Write; the file is not a conforming WAV
// until it succeeds. Finalize does not close the underlying writer.
func (wr *Writer) Finalize() error {
	if wr.finalized {
		return ErrFinalized
	}

	var patch [4]byte

	// RIFF chunk size: total file size minus the 8-byte RIFF header
	binary.LittleEndian.PutUint32(patch[:], uint32(headerSize-8+wr.dataBytes))
	if _, err := wr.w.Seek(riffSizeOffset, io.SeekStart); err != nil {
		return err
	}
	if _, err := wr.w.Write(patch[:]); err != nil {
		return err
	}

	// data chunk size: exact PCM byte count
	binary.LittleEndian.PutUint32(patch[:], uint32(wr.dataBytes))
	if _, err := wr.w.Seek(dataSizeOffset, io.SeekStart); err != nil {
		return err
	}
	if _, err := wr.w.Write(patch[:]); err != nil {
		return err
	}

	if _, err := wr.w.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	wr.finalized = true
	return nil
}
