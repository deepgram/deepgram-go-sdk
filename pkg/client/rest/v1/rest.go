// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package implements a reusable REST client
*/
package restv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

const (
	PackageVersion string = "v1.0"
)

// NewWithDefaults creates a REST client with default options
func NewWithDefaults() *Client {
	return New(&interfaces.ClientOptions{})
}

// New REST client
func New(options *interfaces.ClientOptions) *Client {
	err := options.Parse()
	if err != nil {
		klog.V(1).Infof("options.Parse failed. Err: %v\n", err)
		return nil
	}

	c := Client{
		HTTPClient: NewHTTPClient(options),
		Options:    options,
	}
	return &c
}

// SetupRequest prepares and returns a new REST request with common headers set.
func (c *Client) SetupRequest(ctx context.Context, method, uri string, body io.Reader) (*http.Request, error) {
	klog.V(3).Infof("Using SetupRequest from REST Package\n")

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		return nil, err
	}
	klog.V(4).Infof("%s %s\n", req.Method, uri)

	if headers, ok := ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				klog.V(3).Infof("Custom Header: %s = %s\n", k, v)
				req.Header.Add(k, v)
			}
		}
	}

	req.Header.Set("Host", c.Options.Host)
	req.Header.Set("Accept", "application/json")

	// Set Authorization header based on priority: AccessToken (Bearer) > APIKey (Token)
	token, isBearer := c.Options.GetAuthToken()
	if isBearer {
		req.Header.Set("Authorization", "Bearer "+token)
		klog.V(4).Infof("Using Bearer authentication")
	} else {
		req.Header.Set("Authorization", "token "+token)
		klog.V(4).Infof("Using Token authentication")
	}

	req.Header.Set("User-Agent", interfaces.DgAgent)

	return req, nil
}

// Do is a generic REST API call to the platform
func (c *Client) Do(ctx context.Context, req *http.Request, resBody interface{}) error {
	klog.V(6).Infof("rest.Do() ENTER\n")
	klog.V(4).Infof("%s %s\n", req.Method, req.URL.String())

	err := c.HTTPClient.Do(ctx, req, func(res *http.Response) error {
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusCreated:
		case http.StatusNoContent:
		case http.StatusBadRequest:
			klog.V(4).Infof("HTTP Error Code: %d\n", res.StatusCode)
			detail, err := io.ReadAll(res.Body)
			if err != nil {
				klog.V(4).Infof("io.ReadAll failed. Err: %v\n", err)
				return &interfaces.StatusError{Resp: res}
			}

			// attempt to parse out Deepgram error
			var e interfaces.DeepgramError
			if err := json.Unmarshal(detail, &e); err == nil {
				klog.V(6).Infof("Parsed Deepgram Specific Error\n")
				return &interfaces.StatusError{
					Resp:          res,
					DeepgramError: &e,
				}
			}

			// give standard generic error
			byDetails := bytes.TrimSpace(detail)
			klog.V(1).Infof("Unable to parse Deepgram Error. Err: %s: %s\n", res.Status, byDetails)
			return fmt.Errorf("%s: %s", res.Status, byDetails)
		default:
			return &interfaces.StatusError{Resp: res}
		}

		if resBody == nil {
			klog.V(1).Infof("resBody == nil\n")
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return nil
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			_, err := io.Copy(b, res.Body)
			klog.V(3).Infof("RawResponse\n")
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return err
		case io.Writer:
			_, err := io.Copy(b, res.Body)
			klog.V(3).Infof("io.Writer\n")
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return err
		default:
			resultStr, errRead := io.ReadAll(res.Body)
			if errRead != nil {
				klog.V(1).Infof("io.ReadAll failed. Err: %v\n", errRead)
				klog.V(6).Infof("rest.Do() LEAVE\n")
				return errRead
			}
			klog.V(5).Infof("json.NewDecoder Raw:\n\n%s\n\n", resultStr)
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return json.NewDecoder(strings.NewReader(string(resultStr))).Decode(resBody)
		}
	})

	if err != nil {
		klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		klog.V(6).Infof("rest.Do() LEAVE\n")
		return err
	}

	klog.V(3).Infof("rest.Do Succeeded\n")
	klog.V(6).Infof("rest.Do() LEAVE\n")
	return nil
}
