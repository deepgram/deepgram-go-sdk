// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the Analyze API for Deepgram
package analyze

import (
	"context"
	"io"
	"net/http"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/analyze/v1/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/analyze"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

type AnalyzeClient struct {
	*client.Client
}

func New(client *client.Client) *AnalyzeClient {
	return &AnalyzeClient{client}
}

func (c *AnalyzeClient) FromFile(ctx context.Context, file string, options interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	klog.V(6).Infof("analyze.FromFile ENTER\n")
	klog.V(3).Infof("filePath: %s\n", file)

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("AnalyzeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	var resp api.AnalyzeResponse
	err = c.Client.DoFile(ctx, file, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("analyze.FromFile ENTER\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("analyze.FromFile ENTER\n")
		return nil, err
	}

	klog.V(3).Infof("FromFile Succeeded\n")
	klog.V(6).Infof("analyze.FromFile ENTER\n")
	return &resp, nil
}

func (c *AnalyzeClient) FromStream(ctx context.Context, src io.Reader, options interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	klog.V(6).Infof("analyze.FromStream ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("AnalyzeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	var resp api.AnalyzeResponse
	err = c.Client.DoStream(ctx, src, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("analyze.FromStream LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("analyze.FromStream LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("FromStream Succeeded\n")
	klog.V(6).Infof("analyze.FromStream LEAVE\n")
	return &resp, nil
}

func (c *AnalyzeClient) FromURL(ctx context.Context, url string, options interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	klog.V(6).Infof("analyze.FromURL ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("AnalyzeOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	var resp api.AnalyzeResponse
	err = c.Client.DoURL(ctx, url, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("analyze.FromURL LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("analyze.FromURL LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("FromURL Succeeded\n")
	klog.V(6).Infof("analyze.FromURL LEAVE\n")

	return &resp, nil
}
