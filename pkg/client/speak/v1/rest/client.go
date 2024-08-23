// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package provides the speak client implementation for the Deepgram API
package restv1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	klog "k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
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
	return New("", &interfaces.ClientOptions{})
}

/*
New creates a new speak client with the specified options

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
DoText posts the text to be spoken to a given REST endpoint

Input parameters:
- ctx: context.Context object
- text: contains the text for Text-to-Speech
- req: SpeakOptions which allows configuring things like the speech model, etc.

Output parameters:
- resBody: interface{} which will be populated with the response from the server
*/
func (c *Client) DoText(ctx context.Context, text string, options *interfaces.SpeakOptions, keys []string, resBody interface{}) (map[string]string, error) {
	klog.V(6).Infof("speak.DoText() ENTER\n")

	// obtain URL for the REST API call
	uri, err := version.GetSpeakAPI(ctx, c.Options.Host, c.Options.APIVersion, c.Options.Path, options)
	if err != nil {
		klog.V(1).Infof("version.GetSpeakAPI failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoText() LEAVE\n")
		return nil, err
	}

	// TODO: detect if the source is JSON. If not, then wrap the text in a JSON object
	// and then marshal it to bytes
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(textSource{Text: text})
	if err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoURL() LEAVE\n")
		return nil, err
	}

	req, err := c.SetupRequest(ctx, "POST", uri, strings.NewReader(buf.String()))

	// using the RESTClient SetupRequest (c.SetupRequest vs c.RESTClient.SetupRequest) method which
	// also sets the common headers including the content-type (for example)
	// req, err := c.SetupRequest(ctx, "POST", uri, strings.NewReader(text))
	if err != nil {
		klog.V(1).Infof("SetupRequest failed. Err: %v\n", err)
		klog.V(6).Infof("prerecorded.DoStream() LEAVE\n")
		return nil, err
	}

	// we need to use the HTTPClient + HandleResponse method in order to extract the
	// response headers from the HTTP response. HandleResponse allows us to do that.
	var kv map[string]string
	err = c.HTTPClient.Do(ctx, req, func(res *http.Response) error {
		kv, err = c.HandleResponse(res, keys, resBody)
		return err
	})

	if err != nil {
		klog.V(1).Infof("RESTClient.Do() failed. Err: %v\n", err)
	} else {
		klog.V(4).Infof("DoStream successful\n")
	}
	klog.V(6).Infof("prerecorded.DoStream() LEAVE\n")

	return kv, nil
}
