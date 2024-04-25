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
	"net/http"
	"strings"

	klog "k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/pkg/client/common"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
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
		common.New(apiKey, options),
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

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(textSource{Text: text})
	if err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("speak.DoURL() LEAVE\n")
		return nil, err
	}

	req, err := c.SetupRequest(ctx, "POST", uri, strings.NewReader(buf.String()))
	if err != nil {
		klog.V(1).Infof("SetupRequest failed. Err: %v\n", err)
		klog.V(6).Infof("prerecorded.DoStream() LEAVE\n")
		return nil, err
	}

	var kv map[string]string
	err = c.HTTPClient.Do(ctx, req, func(res *http.Response) error {
		kv, err = c.HandleResponse(res, keys, resBody)
		return err
	})

	if err != nil {
		klog.V(1).Infof("HTTPClient.Do() failed. Err: %v\n", err)
	} else {
		klog.V(4).Infof("DoStream successful\n")
	}
	klog.V(6).Infof("prerecorded.DoStream() LEAVE\n")

	return kv, nil
}
