// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Implementation for a replay device. In this case, replays an audio file to stream into a listener
package replay

import (
	"io"
	"os"

	wav "github.com/youpy/go-wav"
	klog "k8s.io/klog/v2"
)

// New creates an audio replay device
func New(opts ReplayOptions) (*Client, error) {
	klog.V(6).Infof("Replay.New ENTER\n")

	client := &Client{
		options:  opts,
		stopChan: make(chan struct{}),
		muted:    false,
	}

	// create wav decoder instance
	f, err := os.Open(opts.FullFilename)
	if err != nil {
		klog.V(1).Infof("ReplayClient.New os.Open failed. Err: %v\n", err)
		klog.V(6).Infof("Replay.New LEAVE\n")
		return nil, err
	}

	// housekeeping
	client.file = f

	klog.V(3).Infof("Replay.New Succeeded\n")
	klog.V(6).Infof("Replay.New LEAVE\n")

	return client, nil
}

// Start begins streaming the audio for the device
func (c *Client) Start() error {
	reader := wav.NewReader(c.file)
	if reader == nil {
		klog.V(1).Infof("ReplayClient.New wav.NewDecoder is nil\n")
		klog.V(6).Infof("Replay.New LEAVE\n")
		return ErrInvalidInput
	}

	// housekeeping
	c.decoder = reader

	return nil
}

// Read bits from the replay device
func (c *Client) Read() ([]byte, error) {
	buf := make([]byte, defaultBytesToRead)

	byteCount, err := c.decoder.Read(buf)
	if err == io.EOF {
		klog.V(3).Infof("byteBuf.Read failed. Err: %v\n", err)
		return []byte{}, err
	} else if err != nil {
		klog.V(1).Infof("byteBuf.Read failed. Err: %v\n", err)
		return []byte{}, err
	}
	klog.V(7).Infof("byteBuf.Read bytes copied: %d\n", byteCount)

	return buf, nil
}

// Stream is a helper function to stream the replay device data to a source
func (c *Client) Stream(w io.Writer) error {
	for {
		select {
		case <-c.stopChan:
			klog.V(6).Infof("stopChan signal exit\n")
			return nil
		default:
			byData, err := c.Read()
			if err == io.EOF {
				klog.V(3).Infof("decoder.Read EOF\n")
				return nil
			} else if err != nil {
				klog.V(1).Infof("decoder.Read failed. Err: %v\n", err)
				return err
			}

			c.mute.Lock()
			isMuted := c.muted
			c.mute.Unlock()

			if isMuted {
				klog.V(7).Infof("Mic is MUTED!\n")
				byData = make([]byte, len(byData))
			}

			byteCount, err := w.Write(byData)
			if err != nil {
				klog.V(1).Infof("w.Write failed. Err: %v\n", err)
				return err
			}
			klog.V(7).Infof("io.Writer succeeded. Bytes written: %d\n", byteCount)
		}
	}
}

// Mute silences the replay device
func (c *Client) Mute() {
	c.mute.Lock()
	c.muted = true
	c.mute.Unlock()
}

// Unmute restores playback on the replay device
func (c *Client) Unmute() {
	c.mute.Lock()
	c.muted = false
	c.mute.Unlock()
}

// Stop terminates the playback on the replay device
func (c *Client) Stop() error {
	c.decoder = nil

	if c.file != nil {
		c.file.Close()
	}
	c.file = nil

	close(c.stopChan)
	<-c.stopChan

	return nil
}
