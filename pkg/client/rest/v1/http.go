// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// New allocated a Simple HTTP client
func NewHTTPClient(options *interfaces.ClientOptions) *HTTPClient {
	/* #nosec G402 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: options.SkipServerAuth},
		Proxy:           options.Proxy,
	}

	c := HTTPClient{
		Client: http.Client{
			Transport: tr,
		},
		d:         newDebug(),
		UserAgent: interfaces.DgAgent,
		options:   options,
	}
	return &c
}

// Do performs a simple HTTP-style call
func (c *HTTPClient) Do(ctx context.Context, req *http.Request, f func(*http.Response) error) error {
	// Create debugging context for this round trip
	d := c.d.newRoundTrip()
	if d.enabled() {
		defer d.done()
	}

	req.Header.Set("User-Agent", c.UserAgent)

	ext := ""
	if d.enabled() {
		ext = d.debugRequest(req)
	}

	tstart := time.Now()
	res, err := c.Client.Do(req.WithContext(ctx))
	tstop := time.Now()

	if d.enabled() {
		name := fmt.Sprintf("%s %s", req.Method, req.URL)
		d.logf("%6dms (%s)", tstop.Sub(tstart)/time.Millisecond, name)
	}

	if err != nil {
		return err
	}

	if d.enabled() {
		d.debugResponse(res, ext)
	}

	defer res.Body.Close()
	return f(res)
}
