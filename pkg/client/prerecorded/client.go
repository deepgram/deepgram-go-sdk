// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the prerecorded client implementation for the Deepgram API
*/
package prerecorded

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	klog "k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	rest "github.com/deepgram/deepgram-go-sdk/pkg/client/rest"
)

type urlSource struct {
	Url string `json:"url"`
}

/*
NewWithDefaults creates a new websocket connection with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWithDefaults() *Client {
	return New("", &interfaces.ClientOptions{})
}

/*
New creates a new websocket connection with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func New(apiKey string, options *interfaces.ClientOptions) *Client {
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
DoFile posts a file capturing a conversation to a given REST endpoint

Input parameters:
- filePath: string containing the path to the file to be posted
- req: PreRecordedTranscriptionOptions which allows overriding things like language, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoFile(ctx context.Context, filePath string, req interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	klog.V(6).Infof("prerecorded.DoFile() ENTER\n")

	// file?
	fileInfo, err := os.Stat(filePath)
	if err != nil || errors.Is(err, os.ErrNotExist) {
		klog.V(1).Infof("File %s does not exist. Err : %v\n", filePath, err)
		klog.V(6).Infof("prerecorded.DoFile() LEAVE\n")
		return err
	}

	if fileInfo.IsDir() || fileInfo.Size() == 0 {
		klog.V(1).Infof("%s is a directory not a file\n", filePath)
		klog.V(6).Infof("prerecorded.DoFile() LEAVE\n")
		return ErrInvalidInput
	}

	file, err := os.Open(filePath)
	if err != nil {
		klog.V(1).Infof("os.Open(%s) failed. Err : %v\n", filePath, err)
		klog.V(6).Infof("prerecorded.DoFile() LEAVE\n")
		return err
	}
	defer file.Close()

	klog.V(6).Infof("prerecorded.DoFile() LEAVE\n")

	return c.DoStream(ctx, file, req, resBody)
}

/*
DoStream posts a stream capturing a conversation to a given REST endpoint

Input parameters:
- src: io.Reader containing the stream to be posted
- req: PreRecordedTranscriptionOptions which allows overriding things like language, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoStream(ctx context.Context, src io.Reader, options interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	klog.V(6).Infof("rest.DoStream() ENTER\n")

	// obtain URL for the REST API call
	URI, err := version.GetPrerecordedAPI(ctx, c.cOptions.Host, c.cOptions.ApiVersion, c.cOptions.Path, options)
	if err != nil {
		klog.V(1).Infof("version.GetPrerecordedAPI failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoStream() LEAVE\n")
		return err
	}
	klog.V(4).Infof("Connecting to %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "POST", URI, src)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoStream() LEAVE\n")
		return err
	}

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
				klog.V(6).Infof("rest.DoStream() LEAVE\n")
				return &interfaces.StatusError{res}
			}
			klog.V(6).Infof("rest.DoStream() LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
		}

		if resBody == nil {
			klog.V(1).Infof("resBody == nil\n")
			klog.V(6).Infof("rest.DoStream() LEAVE\n")
			return nil
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			klog.V(3).Infof("RawResponse\n")
			klog.V(6).Infof("rest.DoStream() LEAVE\n")
			return res.Write(b)
		case io.Writer:
			klog.V(3).Infof("io.Writer\n")
			klog.V(6).Infof("rest.DoStream() LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			d := json.NewDecoder(res.Body)
			klog.V(3).Infof("json.NewDecoder\n")
			klog.V(6).Infof("rest.DoStream() LEAVE\n")
			return d.Decode(resBody)
		}
	})

	if err != nil {
		klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoStream() LEAVE\n")
		return err
	}

	klog.V(3).Infof("rest.DoStream() Succeeded\n")
	klog.V(6).Infof("rest.DoStream() LEAVE\n")
	return nil
}

// IsUrl returns true if a string is of a URL format
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

/*
DoURL posts a URL capturing a conversation to a given REST endpoint

Input parameters:
- url: string containing the URL to be posted
- req: PreRecordedTranscriptionOptions which allows overriding things like language, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoURL(ctx context.Context, url string, options interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	klog.V(6).Infof("rest.DoURL() ENTER\n")
	klog.V(4).Infof("apiURI: %s\n", url)

	// checks
	validURL := IsUrl(url)
	if !validURL {
		klog.V(1).Infof("Invalid URL: %s\n", url)
		klog.V(6).Infof("rest.DoURL() LEAVE\n")
		return ErrInvalidInput
	}

	// obtain URL
	URI, err := version.GetPrerecordedAPI(ctx, c.cOptions.Host, c.cOptions.ApiVersion, c.cOptions.Path, options)
	if err != nil {
		klog.V(1).Infof("version.GetPrerecordedAPI failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoURL() LEAVE\n")
		return err
	}
	klog.V(4).Infof("Connecting to %s\n", URI)

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(urlSource{Url: url})
	if err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoURL() LEAVE\n")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", URI, &buf)
	if err != nil {
		klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoURL() LEAVE\n")
		return err
	}

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
				klog.V(1).Infof("io.ReadAll failed. Err: %e\n", errBody)
				klog.V(6).Infof("rest.DoURL() LEAVE\n")
				return &interfaces.StatusError{res}
			}
			klog.V(6).Infof("rest.DoURL() LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
		}

		if resBody == nil {
			klog.V(1).Infof("resBody == nil\n")
			klog.V(6).Infof("rest.DoURL() LEAVE\n")
			return nil
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			klog.V(3).Infof("RawResponse\n")
			klog.V(6).Infof("rest.DoURL() LEAVE\n")
			return res.Write(b)
		case io.Writer:
			klog.V(3).Infof("io.Writer\n")
			klog.V(6).Infof("rest.DoURL() LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			d := json.NewDecoder(res.Body)
			klog.V(3).Infof("json.NewDecoder\n")
			klog.V(6).Infof("rest.DoURL() LEAVE\n")
			return d.Decode(resBody)
		}
	})

	if err != nil {
		klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		klog.V(6).Infof("rest.DoURL() LEAVE\n")
		return err
	}

	klog.V(3).Infof("rest.DoURL() Succeeded\n")
	klog.V(6).Infof("rest.DoURL() LEAVE\n")
	return nil
}

/*
Do is a generic REST API call to the platform

Input parameters:
- req: http.Request object

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
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

	req.Header.Set("Host", c.cOptions.Host)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.cOptions.ApiKey)
	req.Header.Set("User-Agent", interfaces.DgAgent)

	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		klog.V(3).Infof("Content-Type = application/json\n")
		req.Header.Set("Content-Type", "application/json")
	}

	err := c.HttpClient.Do(ctx, req, func(res *http.Response) error {
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
				return &interfaces.StatusError{res}
			}
			klog.V(6).Infof("rest.Do() LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
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
			klog.V(3).Infof("io.Writer\n")
			klog.V(6).Infof("rest.Do() LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			d := json.NewDecoder(res.Body)
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

	klog.V(3).Infof("rest.Do() Succeeded\n")
	klog.V(6).Infof("rest.Do() LEAVE\n")
	return nil
}
