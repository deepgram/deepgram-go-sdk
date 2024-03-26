// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the speak client implementation for the Deepgram API
*/
package speak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	klog "k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	rest "github.com/deepgram/deepgram-go-sdk/pkg/client/rest"
)

type textSource struct {
	Text string `json:"text"`
}

/*
NewWithDefaults creates a new speak client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWithDefaults() *Client {
	return New("", interfaces.ClientOptions{})
}

/*
New creates a new speak client with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func New(apiKey string, options interfaces.ClientOptions) *Client {
	if apiKey != "" {
		options.ApiKey = apiKey
	}
	err := options.Parse()
	if err != nil {
		klog.V(1).Infof("options.Parse() failed. Err: %v\n", err)
		return nil
	}

	c := Client{
		Client:   rest.New(options),
		cOptions: options,
	}
	return &c
}

/*
DoText posts the text to be spoken to a given REST endpoint

Input parameters:
- ctx: context.Context object
- text: contains the text for Text-to-Speech
- req: SpeakOptions which allows configuring things like the speech model, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoText(ctx context.Context, text string, options interfaces.SpeakOptions, retValues *map[string]string, resBody interface{}) error {
	klog.V(6).Infof("speak.DoText() ENTER\n")

	// obtain URL for the REST API call
	URI, err := version.GetSpeakAPI(ctx, c.cOptions.Host, c.cOptions.ApiVersion, c.cOptions.Path, options)
	if err != nil {
		klog.V(1).Infof("version.GetSpeakAPI failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoText() LEAVE\n")
		return err
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(textSource{Text: text})
	if err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoURL() LEAVE\n")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", URI, &buf)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoText() LEAVE\n")
		return err
	}
	klog.V(4).Infof("%s %s\n", req.Method, URI)

	if headers, ok := ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				klog.V(3).Infof("Custom Header: %s = %s\n", k, v)
				req.Header.Add(k, v)
			}
		}
	}

	req.Header.Set("Host", c.cOptions.Host)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.cOptions.ApiKey)
	req.Header.Set("User-Agent", interfaces.DgAgent)

	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		klog.V(3).Infof("Content-Type = application/json\n")
		req.Header.Set("Content-Type", "application/json")
	}

	err = c.HttpClient.Do(ctx, req, func(res *http.Response) error {
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusCreated:
		case http.StatusNoContent:
		case http.StatusBadRequest:
			klog.V(4).Infof("HTTP Error Code: %d\n", res.StatusCode)
			detail, errBody := io.ReadAll(res.Body)
			if err != nil {
				klog.V(4).Infof("io.ReadAll failed. Err: %e\n", errBody)
				klog.V(6).Infof("speak.DoText() LEAVE\n")
				return &interfaces.StatusError{res}
			}
			klog.V(6).Infof("speak.DoText() LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
		}

		if resBody == nil {
			klog.V(1).Infof("resBody == nil\n")
			klog.V(6).Infof("speak.DoText() LEAVE\n")
			return nil
		}

		// return values in header
		if retValues == nil {
			*retValues = make(map[string]string)
		}
		for k := range *retValues {
			value := res.Header.Get(k)
			if len(value) > 0 {
				klog.V(4).Infof("RetValue Header: %s = %s\n", k, value)
				(*retValues)[k] = value
				continue
			}
			value = res.Header.Get("dg-" + k)
			if len(value) > 0 {
				klog.V(4).Infof("RetValue Header: %s = %s\n", k, value)
				(*retValues)[k] = value
				continue
			}
			value = res.Header.Get("x-dg-" + k)
			if len(value) > 0 {
				klog.V(4).Infof("RetValue Header: %s = %s\n", k, value)
				(*retValues)[k] = value
				continue
			}
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			klog.V(3).Infof("RawResponse\n")
			klog.V(6).Infof("speak.DoText() LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		case io.Writer:
			klog.V(3).Infof("io.Writer\n")
			klog.V(6).Infof("speak.DoText() LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			klog.V(3).Infof("*io.ReadCloser\n")
			klog.V(6).Infof("speak.DoText() LEAVE\n")
			_, err := io.Copy(b.(*bytes.Buffer), res.Body)
			return err
		}
	})

	if err != nil {
		klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoText() LEAVE\n")
		return err
	}

	klog.V(3).Infof("speak.DoText() Succeeded\n")
	klog.V(6).Infof("speak.DoText() LEAVE\n")
	return nil
}
