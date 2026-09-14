// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func Test_HandleResponseReturnsDeepgramErrorForNonBadRequestStatus(t *testing.T) {
	client := NewREST("test-api-key", &interfaces.ClientOptions{})
	if client == nil {
		t.Fatal("NewREST returned nil")
	}

	res := newTestResponse(
		http.StatusInternalServerError,
		`{"err_code":"ENGINE_ERROR","err_msg":"engine failed","description":"model failed"}`,
	)
	defer res.Body.Close()

	values, err := client.HandleResponse(res, nil, nil)
	if values != nil {
		t.Fatalf("values = %#v, want nil", values)
	}
	if err == nil {
		t.Fatal("HandleResponse returned nil error")
	}

	var statusErr *interfaces.StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("HandleResponse returned %T, want *interfaces.StatusError", err)
	}
	if statusErr.Resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status code = %d, want %d", statusErr.Resp.StatusCode, http.StatusInternalServerError)
	}
	if statusErr.DeepgramError == nil {
		t.Fatal("DeepgramError = nil")
	}
	if statusErr.DeepgramError.ErrCode != "ENGINE_ERROR" {
		t.Fatalf("ErrCode = %q, want %q", statusErr.DeepgramError.ErrCode, "ENGINE_ERROR")
	}
	if statusErr.DeepgramError.ErrMsg != "engine failed" {
		t.Fatalf("ErrMsg = %q, want %q", statusErr.DeepgramError.ErrMsg, "engine failed")
	}
}

func Test_HandleResponseReturnsPlainTextBodyForNonBadRequestStatus(t *testing.T) {
	client := NewREST("test-api-key", &interfaces.ClientOptions{})
	if client == nil {
		t.Fatal("NewREST returned nil")
	}

	res := newTestResponse(
		http.StatusServiceUnavailable,
		"upstream unavailable",
	)
	defer res.Body.Close()

	values, err := client.HandleResponse(res, nil, nil)
	if values != nil {
		t.Fatalf("values = %#v, want nil", values)
	}
	if err == nil {
		t.Fatal("HandleResponse returned nil error")
	}

	got := err.Error()
	if !strings.Contains(got, "503 Service Unavailable") {
		t.Fatalf("error = %q, want status text", got)
	}
	if !strings.Contains(got, "upstream unavailable") {
		t.Fatalf("error = %q, want response body", got)
	}
}

func newTestResponse(statusCode int, body string) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, "https://api.deepgram.com/v1/listen", nil)
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}
