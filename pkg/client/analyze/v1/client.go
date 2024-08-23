// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the analyze/read client implementation for the Deepgram API
package analyzev1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	klog "k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

type textSource struct {
	Text string `json:"text"`
}
type urlSource struct {
	URL string `json:"url"`
}

/*
NewWithDefaults creates a new analyze/read client with all default options

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWithDefaults() *Client {
	return New("", &interfaces.ClientOptions{})
}

/*
New creates a new analyze/read client with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
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

/*
DoFile posts a file capturing a conversation to a given REST endpoint

Input parameters:
- ctx: context.Context object
- filePath: string containing the path to the file to be posted
- req: PreRecordedTranscriptionOptions which allows overriding things like language, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoFile(ctx context.Context, filePath string, req *interfaces.AnalyzeOptions, resBody interface{}) error {
	klog.V(6).Infof("analyze.DoFile() ENTER\n")

	// file?
	fileInfo, err := os.Stat(filePath)
	if err != nil || errors.Is(err, os.ErrNotExist) {
		klog.V(1).Infof("File %s does not exist. Err : %v\n", filePath, err)
		klog.V(6).Infof("analyze.DoFile() LEAVE\n")
		return err
	}

	if fileInfo.IsDir() || fileInfo.Size() == 0 {
		klog.V(1).Infof("%s is a directory not a file\n", filePath)
		klog.V(6).Infof("analyze.DoFile() LEAVE\n")
		return ErrInvalidInput
	}

	file, err := os.Open(filePath)
	if err != nil {
		klog.V(1).Infof("os.Open(%s) failed. Err : %v\n", filePath, err)
		klog.V(6).Infof("analyze.DoFile() LEAVE\n")
		return err
	}
	defer file.Close()

	klog.V(6).Infof("analyze.DoFile() LEAVE\n")

	return c.DoStream(ctx, file, req, resBody)
}

/*
DoStream posts a stream capturing a conversation to a given REST endpoint

Input parameters:
- ctx: context.Context object
- src: io.Reader containing the stream to be posted
- req: PreRecordedTranscriptionOptions which allows overriding things like language, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoStream(ctx context.Context, src io.Reader, options *interfaces.AnalyzeOptions, resBody interface{}) error {
	klog.V(6).Infof("analyze.DoStream() ENTER\n")

	byFile := new(strings.Builder)
	_, err := io.Copy(byFile, src)
	if err != nil {
		klog.V(1).Infof("io.Copy() failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.DoStream() LEAVE\n")
		return err
	}

	klog.V(6).Infof("analyze.DoStream() LEAVE\n")

	return c.DoText(ctx, byFile.String(), options, resBody)
}

func (c *Client) DoText(ctx context.Context, text string, options *interfaces.AnalyzeOptions, resBody interface{}) error {
	klog.V(6).Infof("analyze.DoText() LEAVE\n")

	uri, err := version.GetAnalyzeAPI(ctx, c.Options.Host, c.Options.APIVersion, c.Options.Path, options)
	if err != nil {
		klog.V(1).Infof("GetAnalyzeAPI failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.DoText() LEAVE\n")
		return err
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(textSource{Text: text})
	if err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoText() LEAVE\n")
		return err
	}

	// using the Common SetupRequest (c.SetupRequest vs c.RESTClient.SetupRequest) method which
	// also sets the common headers including the content-type (for example)
	req, err := c.SetupRequest(ctx, "POST", uri, strings.NewReader(buf.String()))
	if err != nil {
		klog.V(1).Infof("SetupRequest failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.DoText() LEAVE\n")
		return err
	}

	// alternatively, we could have used the Common Client Do method, like this
	// but the default one also sets additional "typical" headers like
	// content-type, etc.
	// This is the Common Client way...
	// err = c.Do(ctx, req, func(res *http.Response) error {
	// 	_, err := c.HandleResponse(res, nil, resBody)
	// 	return err
	// })
	// This uses the RESTClient Do method
	err = c.Do(ctx, req, resBody)
	if err != nil {
		klog.V(1).Infof("RESTClient.Do() failed. Err: %v\n", err)
	} else {
		klog.V(4).Infof("DoText successful\n")
	}
	klog.V(6).Infof("analyze.DoText() LEAVE\n")

	return err
}

// IsURL returns true if a string is of a URL format
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

/*
DoURL posts a URL capturing a conversation to a given REST endpoint

Input parameters:
- ctx: context.Context object
- url: string containing the URL to be posted
- req: PreRecordedTranscriptionOptions which allows overriding things like language, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoURL(ctx context.Context, uri string, options *interfaces.AnalyzeOptions, resBody interface{}) error {
	klog.V(6).Infof("analyze.DoURL() ENTER\n")

	if !IsURL(uri) {
		klog.V(1).Infof("Invalid URL: %s\n", uri)
		klog.V(6).Infof("analyze.DoURL() LEAVE\n")
		return ErrInvalidInput
	}

	uri, err := version.GetAnalyzeAPI(ctx, c.Options.Host, c.Options.APIVersion, c.Options.Path, options)
	if err != nil {
		klog.V(1).Infof("GetAnalyzeAPI failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.DoURL() LEAVE\n")
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(urlSource{URL: uri}); err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.DoURL() LEAVE\n")
		return err
	}

	// using the Common SetupRequest (c.SetupRequest vs c.RESTClient.SetupRequest) method which
	// also sets the common headers including the content-type (for example)
	req, err := c.SetupRequest(ctx, "POST", uri, &buf)
	if err != nil {
		klog.V(1).Infof("SetupRequest failed. Err: %v\n", err)
		klog.V(6).Infof("analyze.DoURL() LEAVE\n")
		return err
	}

	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		klog.V(3).Infof("Content-Type = application/json\n")
		req.Header.Set("Content-Type", "application/json")
	}

	err = c.Do(ctx, req, resBody)
	if err != nil {
		klog.V(1).Infof("RESTClient.Do() failed. Err: %v\n", err)
	} else {
		klog.V(4).Infof("DoURL successful\n")
	}
	klog.V(6).Infof("analyze.DoURL() LEAVE\n")

	return err
}
