// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package restv2 defines the Flux TTS batch (REST) API for Deepgram: POST /v2/speak.
package restv2

import (
	"context"
	"io"
	"os"
	"strconv"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/rest/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	speak "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/rest"
)

// Client wraps the Flux TTS batch (REST) transport client.
type Client struct {
	*speak.Client
}

// New creates a Client from a Flux TTS batch transport client.
func New(client *speak.Client) *Client {
	return &Client{client}
}

// ToStream synthesizes text and streams the audio into buf.
// When options.Callback is set, the request is processed asynchronously and buf
// receives the JSON acknowledgement {"request_id":"..."} instead of audio bytes.
func (c *Client) ToStream(ctx context.Context, text string, options *interfaces.SpeakV2Options, buf *interfaces.RawResponse) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speakv2.ToStream ENTER\n")

	keys := initializeKeys()

	err := options.Check()
	if err != nil {
		klog.V(1).Infof("SpeakV2Options.Check() failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToStream LEAVE\n")
		return nil, err
	}

	action := func() (map[string]string, error) {
		return c.Client.DoText(ctx, text, options, keys, buf)
	}

	result, err := c.performAction(action)
	if err != nil {
		klog.V(1).Infof("performAction failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("Speech synthesis successful\n")
	}
	klog.V(6).Infof("speakv2.ToStream LEAVE\n")

	return result, err
}

// ToFile synthesizes text and writes the audio to w.
// When options.Callback is set, the request is processed asynchronously and w
// receives the JSON acknowledgement {"request_id":"..."} instead of audio bytes.
func (c *Client) ToFile(ctx context.Context, text string, options *interfaces.SpeakV2Options, w io.Writer) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speakv2.ToFile ENTER\n")

	keys := initializeKeys()

	err := options.Check()
	if err != nil {
		klog.V(1).Infof("SpeakV2Options.Check() failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToFile LEAVE\n")
		return nil, err
	}

	action := func() (map[string]string, error) {
		return c.Client.DoText(ctx, text, options, keys, w)
	}

	result, err := c.performAction(action)
	if err != nil {
		klog.V(1).Infof("performAction failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("Speech synthesis successful\n")
	}
	klog.V(6).Infof("speakv2.ToFile LEAVE\n")

	return result, err
}

// ToSave synthesizes text and saves the audio to the named file.
func (c *Client) ToSave(ctx context.Context, filename, text string, options *interfaces.SpeakV2Options) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speakv2.ToSave ENTER\n")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		klog.V(1).Infof("os.OpenFile failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}
	defer file.Close()

	result, err := c.ToFile(ctx, text, options, file)
	if err != nil {
		klog.V(1).Infof("speakv2.ToFile failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}

	result.Filename = filename

	klog.V(3).Infof("Saved to file: %v\n", filename)
	klog.V(6).Infof("speakv2.ToSave LEAVE\n")

	return result, nil
}

// helper function
func initializeKeys() []string {
	return []string{
		"content-type",
		"request-id",
		"model-uuid",
		"model-name",
		"char-count",
		"transfer-encoding",
		"date",
		"warnings",
	}
}

// performAction performs the common actions of sending text to the Deepgram API
// and assembling the header-derived response.
func (c *Client) performAction(action func() (map[string]string, error)) (*api.SpeakResponse, error) {
	var resp api.SpeakResponse
	retVal, err := action()
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
			return nil, err
		}
		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		return nil, err
	}

	// char-count is informational; absent on asynchronous (callback) acknowledgements
	if charCntStr, ok := retVal["char-count"]; ok && charCntStr != "" {
		charCnt, err := strconv.Atoi(charCntStr)
		if err != nil {
			klog.V(1).Infof("strconv.Atoi(char-count) failed. Err: %v\n", err)
		} else {
			resp.Characters = charCnt
		}
	}

	resp.ContextType = retVal["content-type"]
	resp.RequestID = retVal["request-id"]
	resp.ModelUUID = retVal["model-uuid"]
	resp.ModelName = retVal["model-name"]
	resp.TransferEncoding = retVal["transfer-encoding"]
	resp.Date = retVal["date"]
	resp.Warnings = retVal["warnings"]

	return &resp, nil
}
