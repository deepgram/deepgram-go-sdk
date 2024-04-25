// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package implements a reusable REST client
*/
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
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

// Do is a generic REST API call to the platform
func (c *Client) Do(ctx context.Context, req *http.Request, resBody interface{}) error {
	klog.V(6).Infof("rest.Do() ENTER\n")

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
	req.Header.Set("Authorization", "token "+c.Options.APIKey)
	req.Header.Set("User-Agent", interfaces.DgAgent)

	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		klog.V(3).Infof("Content-Type = application/json\n")
		req.Header.Set("Content-Type", "application/json")
	}

	klog.V(4).Infof("%s %s\n", req.Method, req.URL.String())

	err := c.HTTPClient.Do(ctx, req, func(res *http.Response) error {
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusCreated:
		case http.StatusNoContent:
		case http.StatusBadRequest:
			klog.V(1).Infof("HTTP Error Code: %d\n", res.StatusCode)
			detail, errBody := io.ReadAll(res.Body)
			if errBody != nil {
				klog.V(1).Infof("io.ReadAll failed. Err: %e\n", errBody)
				klog.V(6).Infof("rest.Do() LEAVE\n")
				return &interfaces.StatusError{Resp: res}
			}
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
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
			klog.V(3).Infof("RawResponse\n")
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return res.Write(b)
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
			d := json.NewDecoder(strings.NewReader(string(resultStr)))
			klog.V(3).Infof("json.NewDecoder\n")
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return d.Decode(resBody)
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
