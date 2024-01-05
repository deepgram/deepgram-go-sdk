// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the Pre-recorded API for Deepgram
package prerecorded

import (
	"context"
	"io"
	"net/http"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/prerecorded/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
)

type PrerecordedClient struct {
	*client.Client
}

func New(client *client.Client) *PrerecordedClient {
	return &PrerecordedClient{client}
}

func (c *PrerecordedClient) FromFile(ctx context.Context, file string, options interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error) {
	klog.V(6).Infof("FromFile ENTER\n")
	klog.V(3).Infof("filePath: %s\n", file)

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("TranscribeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	var resp api.PreRecordedResponse
	err = c.Client.DoFile(ctx, file, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("prerecorded.FromFile ENTER\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("prerecorded.FromFile ENTER\n")
		return nil, err
	}

	klog.V(3).Infof("FromFile Succeeded\n")
	klog.V(6).Infof("prerecorded.FromFile ENTER\n")
	return &resp, nil
}

func (c *PrerecordedClient) FromStream(ctx context.Context, src io.Reader, options interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error) {
	klog.V(6).Infof("prerecorded.FromStream ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("TranscribeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	var resp api.PreRecordedResponse
	err = c.Client.DoStream(ctx, src, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("prerecorded.FromStream LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("prerecorded.FromStream LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("FromStream Succeeded\n")
	klog.V(6).Infof("prerecorded.FromStream LEAVE\n")
	return &resp, nil
}

func (c *PrerecordedClient) FromURL(ctx context.Context, url string, options interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error) {
	klog.V(6).Infof("prerecorded.FromURL ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("TranscribeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	var resp api.PreRecordedResponse
	err = c.Client.DoURL(ctx, url, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("prerecorded.FromURL LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("prerecorded.FromURL LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("FromURL Succeeded\n")
	klog.V(6).Infof("prerecorded.FromURL LEAVE\n")

	return &resp, nil
}
