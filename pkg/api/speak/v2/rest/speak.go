// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package restv2 defines the Flux TTS batch (REST) API for Deepgram: POST /v2/speak.
package restv2

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/rest/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
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

// validate rejects the setup mistakes that would otherwise panic or corrupt
// output: a nil transport client, nil options, and options that fail Check().
func (c *Client) validate(options *interfaces.SpeakV2Options) error {
	if c == nil || c.Client == nil {
		return ErrNilClient
	}
	if options == nil {
		return interfacesv2.ErrOptionsRequired
	}
	return options.Check()
}

// ToStream synthesizes text and streams the audio into buf.
// When options.Callback is set, the request is processed asynchronously and buf
// receives the JSON acknowledgement {"request_id":"..."} instead of audio bytes;
// prefer ToAsync for a typed acknowledgement.
func (c *Client) ToStream(ctx context.Context, text string, options *interfaces.SpeakV2Options, buf *interfaces.RawResponse) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speakv2.ToStream ENTER\n")

	keys := initializeKeys()

	if err := c.validate(options); err != nil {
		klog.V(1).Infof("speakv2.ToStream validation failed. Err: %v\n", err)
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
// receives the JSON acknowledgement {"request_id":"..."} instead of audio bytes;
// prefer ToAsync for a typed acknowledgement.
func (c *Client) ToFile(ctx context.Context, text string, options *interfaces.SpeakV2Options, w io.Writer) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speakv2.ToFile ENTER\n")

	keys := initializeKeys()

	if err := c.validate(options); err != nil {
		klog.V(1).Infof("speakv2.ToFile validation failed. Err: %v\n", err)
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

// ToSave synthesizes text and saves the audio to the named file. The audio is
// written to a temporary sibling file and renamed into place only on success, so
// invalid options or a failed request never truncate or corrupt an existing file.
//
// Asynchronous (callback) requests are rejected: their response body is a JSON
// acknowledgement, not audio, and saving it under an audio filename would be
// misleading. Use ToAsync instead.
func (c *Client) ToSave(ctx context.Context, filename, text string, options *interfaces.SpeakV2Options) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speakv2.ToSave ENTER\n")

	// validate everything before touching the destination
	if err := c.validate(options); err != nil {
		klog.V(1).Infof("speakv2.ToSave validation failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}
	if options.Callback != "" {
		klog.V(1).Infof("speakv2.ToSave rejected callback options\n")
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, ErrCallbackNotSupported
	}

	tmp, err := os.CreateTemp(filepath.Dir(filename), filepath.Base(filename)+".*.tmp")
	if err != nil {
		klog.V(1).Infof("os.CreateTemp failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}
	tmpName := tmp.Name()
	cleanup := func() {
		tmp.Close()
		if rmErr := os.Remove(tmpName); rmErr != nil && !os.IsNotExist(rmErr) {
			klog.V(1).Infof("os.Remove(%s) failed. Err: %v\n", tmpName, rmErr)
		}
	}

	result, err := c.ToFile(ctx, text, options, tmp)
	if err != nil {
		klog.V(1).Infof("speakv2.ToFile failed. Err: %v\n", err)
		cleanup()
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}

	if err := tmp.Close(); err != nil {
		klog.V(1).Infof("tmp.Close failed. Err: %v\n", err)
		cleanup()
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}
	if err := os.Rename(tmpName, filename); err != nil {
		klog.V(1).Infof("os.Rename failed. Err: %v\n", err)
		cleanup()
		klog.V(6).Infof("speakv2.ToSave LEAVE\n")
		return nil, err
	}

	result.Filename = filename

	klog.V(3).Infof("Saved to file: %v\n", filename)
	klog.V(6).Infof("speakv2.ToSave LEAVE\n")

	return result, nil
}

// ToAsync submits text for asynchronous synthesis. options.Callback must be set;
// the server replies immediately with a JSON acknowledgement and delivers the
// audio to the callback URL when synthesis completes. No local file is created.
func (c *Client) ToAsync(ctx context.Context, text string, options *interfaces.SpeakV2Options) (*api.AsyncResponse, error) {
	klog.V(6).Infof("speakv2.ToAsync ENTER\n")

	if err := c.validate(options); err != nil {
		klog.V(1).Infof("speakv2.ToAsync validation failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToAsync LEAVE\n")
		return nil, err
	}
	if options.Callback == "" {
		klog.V(1).Infof("speakv2.ToAsync requires a callback URL\n")
		klog.V(6).Infof("speakv2.ToAsync LEAVE\n")
		return nil, ErrCallbackRequired
	}

	var buf interfaces.RawResponse
	if _, err := c.Client.DoText(ctx, text, options, initializeKeys(), &buf); err != nil {
		klog.V(1).Infof("speakv2.DoText failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToAsync LEAVE\n")
		return nil, err
	}

	var resp api.AsyncResponse
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		klog.V(1).Infof("json.Unmarshal(AsyncResponse) failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.ToAsync LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("Asynchronous synthesis accepted: request_id=%s\n", resp.RequestID)
	klog.V(6).Infof("speakv2.ToAsync LEAVE\n")
	return &resp, nil
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
