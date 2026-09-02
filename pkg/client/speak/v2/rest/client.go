// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// Package restv2 provides the Flux TTS batch (REST) client for POST /v2/speak.
package restv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	klog "k8s.io/klog/v2"

	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
)

type textSource struct {
	Text string `json:"text"`
}

/*
NewWithDefaults creates a new Flux TTS batch client with all default options.
It returns an error instead of a nil client when configuration fails (for example
when no credentials are available), so callers cannot accidentally dereference a
nil client.

Notes:
  - The Deepgram API KEY is read from the environment variable DEEPGRAM_API_KEY
*/
func NewWithDefaults() (*Client, error) {
	return New("", &interfaces.ClientOptions{})
}

/*
New creates a new Flux TTS batch client with the specified options. It returns an
error instead of a nil client when configuration fails (for example when no
credentials are available).

Input parameters:
- apiKey: string containing the Deepgram API key. If empty, the DEEPGRAM_API_KEY environment variable is used.
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.
*/
func New(apiKey string, options *interfaces.ClientOptions) (*Client, error) {
	if options == nil {
		options = &interfaces.ClientOptions{}
	}
	if apiKey != "" {
		options.APIKey = apiKey
	}
	if err := options.Parse(); err != nil {
		klog.V(1).Infof("options.Parse() failed. Err: %v\n", err)
		return nil, err
	}

	c := Client{
		common.NewREST(apiKey, options),
	}
	return &c, nil
}

/*
DoText posts the text to be synthesized to POST /v2/speak

Input parameters:
  - ctx: context.Context object
  - text: contains the text for Text-to-Speech
  - options: SpeakV2Options which allows configuring things like the model, encoding, etc.
  - keys: response headers to extract into the returned map
  - resBody: the response audio is written here (io.Writer or RawResponse); when a
    callback URL is set, the JSON acknowledgement bytes are written instead

Output parameters:
- map of the extracted response headers
*/
func (c *Client) DoText(ctx context.Context, text string, options *interfacesv2.SpeakV2Options, keys []string, resBody interface{}) (map[string]string, error) {
	klog.V(6).Infof("speakv2.DoText() ENTER\n")

	// obtain URL for the REST API call
	uri, err := version.GetSpeakV2API(ctx, c.Options.Host, c.Options.APIVersion, c.Options.Path, options)
	if err != nil {
		klog.V(1).Infof("version.GetSpeakV2API failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.DoText() LEAVE\n")
		return nil, err
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(textSource{Text: text})
	if err != nil {
		klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.DoText() LEAVE\n")
		return nil, err
	}

	req, err := c.SetupRequest(ctx, "POST", uri, strings.NewReader(buf.String()))
	if err != nil {
		klog.V(1).Infof("SetupRequest failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.DoText() LEAVE\n")
		return nil, err
	}

	// use HTTPClient + handleResponse to extract the response headers alongside the body
	var kv map[string]string
	err = c.HTTPClient.Do(ctx, req, func(res *http.Response) error {
		kv, err = c.handleResponse(res, keys, resBody)
		return err
	})

	if err != nil {
		klog.V(1).Infof("RESTClient.Do() failed. Err: %v\n", err)
		klog.V(6).Infof("speakv2.DoText() LEAVE\n")
		return nil, err
	}

	klog.V(4).Infof("DoText successful\n")
	klog.V(6).Infof("speakv2.DoText() LEAVE\n")
	return kv, nil
}

// handleResponse delegates success responses to the shared handler and preserves
// the Deepgram error payload for EVERY non-success status — the shared v1 handler
// reads the body only for HTTP 400, which loses the details of authentication
// failures (401/403), payload and rate limits (413/429), and 5xx responses.
func (c *Client) handleResponse(res *http.Response, keys []string, resBody interface{}) (map[string]string, error) {
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return c.HandleResponse(res, keys, resBody)
	}

	klog.V(4).Infof("HTTP Error Code: %d\n", res.StatusCode)
	detail, err := io.ReadAll(res.Body)
	if err != nil {
		klog.V(4).Infof("io.ReadAll failed. Err: %v\n", err)
		return nil, &interfaces.StatusError{Resp: res}
	}

	// attempt to parse out a typed Deepgram error
	var e interfaces.DeepgramError
	if jsonErr := json.Unmarshal(detail, &e); jsonErr == nil && (e.ErrCode != "" || e.ErrMsg != "" || e.Description != "") {
		klog.V(6).Infof("Parsed Deepgram Specific Error\n")
		return nil, &interfaces.StatusError{
			Resp:          res,
			DeepgramError: &e,
		}
	}

	// fall back to the raw body so no detail is lost
	byDetails := bytes.TrimSpace(detail)
	if len(byDetails) > 0 {
		klog.V(1).Infof("Unable to parse Deepgram Error. Err: %s: %s\n", res.Status, byDetails)
		return nil, fmt.Errorf("%s: %s", res.Status, byDetails)
	}
	return nil, &interfaces.StatusError{Resp: res}
}
