// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package contains the code for the Keys APIs in the Deepgram Manage API
package managev1

import (
	"context"
	"io"

	"k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

const (
	PackageVersion string = "v1.0"
)

// New creates a new Client
func NewWithDefaults() *Client {
	return New("", &interfaces.ClientOptions{})
}

func New(apiKey string, options *interfaces.ClientOptions) *Client {
	if apiKey != "" {
		options.APIKey = apiKey
	}
	err := options.Parse()
	if err != nil {
		klog.V(1).Infof("options.Parse() failed. Err: %v\n", err)
		return nil
	}

	c := Client{
		common.NewREST(apiKey, options),
	}
	return &c
}

func (c *Client) APIRequest(ctx context.Context, method, apiPath string, body io.Reader, resBody interface{}, params ...interface{}) error {
	klog.V(6).Infof("manage.%s() ENTER\n", method+apiPath) // Dynamic entry log based on method and path

	// Construct the uri with parameters
	uri, err := version.GetManageAPI(ctx, c.Options.Host, c.Options.APIVersion, apiPath, nil, params...)
	if err != nil {
		klog.V(1).Infof("GetManageAPI failed. Err: %v\n", err)
		klog.V(6).Infof("manage.%s() LEAVE\n", method+apiPath)
		return err
	}

	// using the Common.SetupRequest (c.SetupRequest vs c.RESTClient.SetupRequest) method which
	// also sets the common headers including the content-type (for example)
	req, err := c.SetupRequest(ctx, method, uri, body)
	if err != nil {
		klog.V(6).Infof("manage.%s() LEAVE\n", method+apiPath)
		return err
	}

	// altertatively, we could have used the Common Client Do method, like this
	// but the default one also sets additional "typical" headers like
	// content-type, etc.
	// This is the Common Client way...
	// err = c.Do(ctx, req, func(res *http.Response) error {
	// 	_, err := c.HandleResponse(res, nil, resBody)
	// 	return err
	// })
	// This uses the RESTClient Do method
	err = c.Do(ctx, req, &resBody)
	if err != nil {
		klog.V(6).Infof("manage.%s() LEAVE\n", method+apiPath)
		return err
	}

	klog.V(3).Infof("%s succeeded\n", method+apiPath)
	klog.V(6).Infof("manage.%s() LEAVE\n", method+apiPath)
	return nil
}
