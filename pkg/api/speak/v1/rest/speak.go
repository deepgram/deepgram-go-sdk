// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the Speak REST API for Deepgram
package restv1

import (
	"context"
	"io"
	"os"
	"strconv"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/rest/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	speak "github.com/deepgram/deepgram-go-sdk/pkg/client/speak/v1/rest"
)

type Client struct {
	*speak.Client
}

func New(client *speak.Client) *Client {
	return &Client{client}
}

// ToStream TTS streamed to a buffer
func (c *Client) ToStream(ctx context.Context, text string, options *interfaces.SpeakOptions, buf *interfaces.RawResponse) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speak.ToStream ENTER\n")

	keys := initializeKeys()

	err := options.Check()
	if err != nil {
		klog.V(1).Infof("SpeakOptions.Check() failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToStream LEAVE\n")
		return nil, err
	}

	action := func() (map[string]string, error) {
		return c.Client.DoText(ctx, text, options, keys, buf)
	}

	result, err := c.performAction(action)
	if err != nil {
		klog.V(1).Infof("performAction failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("Transcription successful\n")
	}
	klog.V(6).Infof("speak.ToStream LEAVE\n")

	return result, err
}

// ToFile TTS saved to a file
func (c *Client) ToFile(ctx context.Context, text string, options *interfaces.SpeakOptions, w io.Writer) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speak.ToFile ENTER\n")

	keys := initializeKeys()

	err := options.Check()
	if err != nil {
		klog.V(1).Infof("SpeakOptions.Check() failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToFile LEAVE\n")
		return nil, err
	}

	action := func() (map[string]string, error) {
		return c.Client.DoText(ctx, text, options, keys, w)
	}

	result, err := c.performAction(action)
	if err != nil {
		klog.V(1).Infof("performAction failed. Err: %v\n", err)
	} else {
		klog.V(3).Infof("Transcription successful\n")
	}
	klog.V(6).Infof("speak.ToFile LEAVE\n")

	return result, err
}

// ToSave TTS saved to a file
func (c *Client) ToSave(ctx context.Context, filename, text string, options *interfaces.SpeakOptions) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speak.ToSave ENTER\n")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		klog.V(1).Infof("os.OpenFile failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToSave LEAVE\n")
		return nil, err
	}
	defer file.Close()

	result, err := c.ToFile(ctx, text, options, file)
	if err != nil {
		klog.V(1).Infof("speak.ToFile failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToSave LEAVE\n")
		return nil, err
	}

	result.Filename = filename

	klog.V(3).Infof("Saved to file: %v\n", filename)
	klog.V(6).Infof("speak.ToSave LEAVE\n")

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
	}
}

// performAction performs the common actions of sending text to the Deepgram API and handling the response.
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

	charCnt, err := strconv.Atoi(retVal["char-count"])
	if err != nil {
		klog.V(1).Infof("strconv.Atoi failed. Err: %v\n", err)
		return nil, err
	}

	resp.ContextType = retVal["content-type"]
	resp.RequestID = retVal["request-id"]
	resp.ModelUUID = retVal["model-uuid"]
	resp.ModelName = retVal["model-name"]
	resp.Characters = charCnt
	resp.TransferEncoding = retVal["transfer-encoding"]
	resp.Date = retVal["date"]

	return &resp, nil
}
