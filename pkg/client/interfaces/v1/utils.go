// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// constants
const (
	sdkVersion string = "v1.4.0"
)

// DgAgent is the user agent string for the SDK
var DgAgent string = "@deepgram/sdk/" + sdkVersion + " go/" + goVersion()

func goVersion() string {
	version := runtime.Version()
	if strings.HasPrefix(version, "go") {
		return version[2:]
	}
	return version
}

/*
	custom headers and configuration options
*/
// Signer callback for the certificant signer
type Signer interface {
	SignRequest(*http.Request) error
}

// SignerContext blackbox of data
type SignerContext struct{}

// WithSigner appends a signer to the given context
func WithSigner(ctx context.Context, s Signer) context.Context {
	return context.WithValue(ctx, SignerContext{}, s)
}

// HeadersContext blackbox of data
type HeadersContext struct{}

// WithCustomHeaders appends a header to the given context
func WithCustomHeaders(ctx context.Context, headers http.Header) context.Context {
	return context.WithValue(ctx, HeadersContext{}, headers)
}

// ParametersContext blackbox of data
type ParametersContext struct{}

// WithCustomParameters
func WithCustomParameters(ctx context.Context, params map[string][]string) context.Context {
	return context.WithValue(ctx, ParametersContext{}, params)
}

/*
RawResponse may be used with the Do method as the resBody argument in order
to capture the raw response data.
*/
type RawResponse struct {
	bytes.Buffer
}

// DeepgramWarning is the Deepgram specific response warning
type DeepgramWarning struct {
	Type        string
	WarnCode    string `json:"warn_code,omitempty"`
	WarnMsg     string `json:"warn_msg,omitempty"`
	Description string `json:"description,omitempty"`
	Variant     string `json:"variant,omitempty"`
}

// DeepgramError is the Deepgram specific response error
type DeepgramError struct {
	Type        string
	ErrCode     string `json:"err_code,omitempty"`
	ErrMsg      string `json:"err_msg,omitempty"`
	Description string `json:"description,omitempty"`
	Variant     string `json:"variant,omitempty"`
}

// StatusError captures a REST error in the library
type StatusError struct {
	Resp          *http.Response
	DeepgramError *DeepgramError
}

// Error string representation for a given error
func (e *StatusError) Error() string {
	if e.DeepgramError != nil && e.DeepgramError.ErrMsg != "" {
		return fmt.Sprintf("%s %s: %s - %s", e.Resp.Request.Method, e.Resp.Request.URL, e.Resp.Status, e.DeepgramError.ErrMsg)
	}
	return fmt.Sprintf("%s %s: %s", e.Resp.Request.Method, e.Resp.Request.URL, e.Resp.Status)
}
