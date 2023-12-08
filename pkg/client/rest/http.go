// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

// New allocated a Simple HTTP client
func NewHTTPClient() *HttpClient {
	bDisable := true
	if v := os.Getenv("DEEPGRAM_SSL_HOST_VERIFICATION"); v != "" {
		bDisable = strings.EqualFold(strings.ToLower(v), "false")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: bDisable},
	}

	c := HttpClient{
		Client: http.Client{
			Transport: tr,
		},
		d:         newDebug(),
		UserAgent: interfaces.DgAgent,
	}
	return &c
}

// Do performs a simple HTTP-style call
func (c *HttpClient) Do(ctx context.Context, req *http.Request, f func(*http.Response) error) error {
	// checks
	if ctx == nil {
		ctx = context.Background()
	}

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
