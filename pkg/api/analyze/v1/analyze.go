// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the Analyze API for Deepgram
package analyze

import (
	"context"
	"io"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/analyze/v1/interfaces"
	analyze "github.com/deepgram/deepgram-go-sdk/pkg/client/analyze"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

type Client struct {
	*analyze.Client
}

// New creates a new Client
func New(client *analyze.Client) *Client {
	return &Client{client}
}

// FromFile transcribes a prerecorded audio file from a file
func (c *Client) FromFile(ctx context.Context, file string, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	return c.sendAnalysis(ctx, func(ctx context.Context, opts *interfaces.AnalyzeOptions, resp interface{}) error {
		return c.sendFile(ctx, file, opts, resp)
	}, options)
}

// FromStream transcribes a prerecorded audio file from a stream
func (c *Client) FromStream(ctx context.Context, src io.Reader, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	return c.sendAnalysis(ctx, func(ctx context.Context, opts *interfaces.AnalyzeOptions, resp interface{}) error {
		return c.sendStream(ctx, src, opts, resp)
	}, options)
}

// FromURL transcribes a prerecorded audio file from a URL
func (c *Client) FromURL(ctx context.Context, url string, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	return c.sendAnalysis(ctx, func(ctx context.Context, opts *interfaces.AnalyzeOptions, resp interface{}) error {
		return c.sendURL(ctx, url, opts, resp)
	}, options)
}

// private functions
type sendFunc func(ctx context.Context, options *interfaces.AnalyzeOptions, resp interface{}) error

func (c *Client) sendFile(ctx context.Context, filePath string, options *interfaces.AnalyzeOptions, resp interface{}) error {
	return c.Client.DoFile(ctx, filePath, options, resp)
}

func (c *Client) sendStream(ctx context.Context, src io.Reader, options *interfaces.AnalyzeOptions, resp interface{}) error {
	return c.Client.DoStream(ctx, src, options, resp)
}

func (c *Client) sendURL(ctx context.Context, url string, options *interfaces.AnalyzeOptions, resp interface{}) error {
	return c.Client.DoURL(ctx, url, options, resp)
}

func (c *Client) sendAnalysis(ctx context.Context, sender sendFunc, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error) {
	klog.V(6).Infof("analyze.sendAnalysis ENTER\n")

	err := options.Check()
	if err != nil {
		klog.V(1).Infof("AnalyzeOptions.Check() failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.sendAnalysis LEAVE\n")
		return nil, err
	}

	var resp api.AnalyzeResponse

	err = sender(ctx, options, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
		}
		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("analyze.sendAnalysis LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("sendAnalysis Succeeded\n")
	klog.V(6).Infof("analyze.sendAnalysis LEAVE\n")
	return &resp, nil
}
