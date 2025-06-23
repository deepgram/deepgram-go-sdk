// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

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
	restv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/rest/v1"
)

func NewREST(apiKey string, options *interfaces.ClientOptions) *RESTClient {
	if apiKey != "" {
		options.APIKey = apiKey
	}
	err := options.Parse()
	if err != nil {
		klog.V(1).Infof("options.Parse() failed. Err: %v\n", err)
		return nil
	}

	c := RESTClient{
		restv1.New(options),
	}

	return &c
}

// SetupRequest prepares and returns a new HTTP request with common headers set.
func (c *RESTClient) SetupRequest(ctx context.Context, method, uri string, body io.Reader) (*http.Request, error) {
	klog.V(3).Infof("Using SetupRequest from Common Package\n")

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
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// HandleResponse processes the HTTP response for both streaming and URL-based API requests.
func (c *RESTClient) HandleResponse(res *http.Response, keys []string, resBody interface{}) (map[string]string, error) {
	klog.V(6).Infof("Handle HTTP response\n")
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return decodeResponseBody(res, keys, resBody)
	case http.StatusBadRequest:
		klog.V(4).Infof("HTTP Error Code: %d\n", res.StatusCode)
		detail, err := io.ReadAll(res.Body)
		if err != nil {
			klog.V(4).Infof("io.ReadAll failed. Err: %v\n", err)
			return nil, &interfaces.StatusError{Resp: res}
		}

		// attempt to parse out Deepgram error
		var e interfaces.DeepgramError
		if err := json.Unmarshal(detail, &e); err == nil {
			klog.V(6).Infof("Parsed Deepgram Specific Error\n")
			return nil, &interfaces.StatusError{
				Resp:          res,
				DeepgramError: &e,
			}
		}

		// give standard generic error
		byDetails := bytes.TrimSpace(detail)
		klog.V(1).Infof("Unable to parse Deepgram Error. Err: %s: %s\n", res.Status, byDetails)
		return nil, fmt.Errorf("%s: %s", res.Status, byDetails)
	default:
		return nil, &interfaces.StatusError{Resp: res}
	}
}

// decodeResponseBody decodes the HTTP response body into the provided resBody based on its type.
func decodeResponseBody(res *http.Response, keys []string, resBody interface{}) (map[string]string, error) {
	retValues := make(map[string]string)

	// return values in header
	if len(keys) > 0 {
		for _, k := range keys {
			value := res.Header.Get(k)
			if len(value) > 0 {
				klog.V(4).Infof("RetValue Header: %s = %s\n", k, value)
				retValues[k] = value
				continue
			}
			value = res.Header.Get("dg-" + k)
			if len(value) > 0 {
				klog.V(4).Infof("RetValue Header: %s = %s\n", k, value)
				retValues[k] = value
				continue
			}
			value = res.Header.Get("x-dg-" + k)
			if len(value) > 0 {
				klog.V(4).Infof("RetValue Header: %s = %s\n", k, value)
				retValues[k] = value
				continue
			}
		}
	}

	switch b := resBody.(type) {
	case *interfaces.RawResponse:
		klog.V(3).Infof("RawResponse\n")
		_, err := io.Copy(b, res.Body)
		return retValues, err
	case io.Writer:
		klog.V(3).Infof("io.Writer\n")
		_, err := io.Copy(b, res.Body)
		return retValues, err
	default:
		klog.V(3).Infof("*io.ReadCloser\n")
		resultStr, err := io.ReadAll(res.Body)
		if err != nil {
			klog.V(1).Infof("io.ReadAll failed. Err: %v\n", err)
			return nil, err
		}
		klog.V(5).Infof("json.NewDecoder Raw:\n\n%s\n\n", resultStr)
		return retValues, json.NewDecoder(strings.NewReader(string(resultStr))).Decode(resBody)
	}
}
