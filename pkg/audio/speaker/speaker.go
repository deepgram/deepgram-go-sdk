// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Implementation of a speaker using portaudio
package speaker

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/gordonklaus/portaudio"
	klog "k8s.io/klog/v2"

	common "github.com/deepgram/deepgram-go-sdk/pkg/audio/common"
)

// Initialize inits the library. This handles OS level init of the library.
func Initialize() {
	common.Initialize()
}

// Teardown cleans up the library. This handles OS level cleanup.
func Teardown() {
	common.Teardown()
}

// New creates a new speaker using portaudio
func New(cfg AudioConfig) (*Speaker, error) {
	klog.V(6).Infof("Speaker.New ENTER\n")

	s := &Speaker{
		stopChan: make(chan struct{}),
		intBuf:   make([]int16, defaultBytesToRead),
		muted:    false,
	}

	// Open a stream for audio output (0 input channels, cfg.OutputChannels output channels)
	stream, err := portaudio.OpenDefaultStream(0, cfg.OutputChannels, float64(cfg.SamplingRate), len(s.intBuf), s.intBuf)
	if err != nil {
		klog.V(1).Infof("OpenDefaultStream failed. Err: %v\n", err)
		klog.V(6).Infof("Speaker.New LEAVE\n")
		return nil, err
	}

	// housekeeping
	s.stream = stream

	klog.V(3).Infof("OpenDefaultStream succeeded\n")
	klog.V(6).Infof("Speaker.New LEAVE\n")

	return s, nil
}

// Start begins the listening on the microphone
func (s *Speaker) Start() error {
	err := s.stream.Start()
	if err != nil {
		klog.V(1).Infof("Speaker failed to start. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("Speaker.Start() succeeded\n")
	return nil
}

// Write sends audio data to the speaker
func (s *Speaker) Write(p []byte) (n int, err error) {
	// if err := binary.Read(bytes.NewReader(p), binary.LittleEndian, s.intBuf); err != nil {
	// 	klog.V(1).Infof("binary.Read failed. Err: %v\n", err)
	// 	return 0, err
	// }

	buf := s.littleEndianByteToInt16(p)
	byteCopied := copy(buf, s.intBuf)
	klog.V(7).Infof("stream.Read bytes copied: %d\n", byteCopied)

	if err := s.stream.Write(); err != nil {
		klog.V(1).Infof("Speaker.Write failed. Err: %v\n", err)
		return 0, err
	}

	klog.V(7).Infof("Speaker.Write succeeded. Bytes written: %d\n", len(p))
	return len(p), nil
}

// Stream is a helper function to stream audio data to the speaker
func (s *Speaker) Stream(r io.Reader) error {
	for {
		select {
		case <-s.stopChan:
			return nil
		default:
			buf := make([]byte, len(s.intBuf)*2) // 2 bytes per int16
			n, err := r.Read(buf)
			if err != nil || n == 0 {
				klog.V(1).Infof("r.Read failed. Err: %v\n", err)
				return err
			}

			if _, err := s.Write(buf); err != nil {
				klog.V(1).Infof("s.Write failed. Err: %v\n", err)
				return err
			}
		}
	}
}

// Mute silences the mic
func (s *Speaker) Mute() {
	s.mute.Lock()
	s.muted = true
	s.mute.Unlock()
}

// Unmute restores recording on the mic
func (s *Speaker) Unmute() {
	s.mute.Lock()
	s.muted = false
	s.mute.Unlock()
}

// Stop terminates listening on the mic
func (s *Speaker) Stop() error {
	err := s.stream.Stop()
	if err != nil {
		klog.V(1).Infof("stream.Stop failed. Err: %v\n", err)
		return err
	}

	close(s.stopChan)
	<-s.stopChan

	return nil
}

func (s *Speaker) littleEndianByteToInt16(b []byte) []int16 {
	s.mute.Lock()
	isMuted := s.muted
	s.mute.Unlock()

	if isMuted {
		klog.V(7).Infof("Mic is MUTED!\n")
		return make([]int16, len(b)/2)
	}

	var result []int16
	buf := bytes.NewReader(b)
	for buf.Len() > 0 {
		var value int16
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil {
			klog.V(1).Infof("binary.Read failed. Err: %v\n", err)
			break
		}
		result = append(result, value)
	}

	return result
}
