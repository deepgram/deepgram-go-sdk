// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func TestDoReturnsDeepgramErrorForNonBadRequestStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"err_code":"ENGINE_ERROR","err_msg":"engine failed","description":"model failed"}`))
	}))
	defer server.Close()

	client := New(&interfaces.ClientOptions{APIKey: "test-api-key", Host: server.URL})
	if client == nil {
		t.Fatal("New returned nil")
	}

	req, err := client.SetupRequest(context.Background(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("SetupRequest returned error: %v", err)
	}

	err = client.Do(context.Background(), req, nil)
	if err == nil {
		t.Fatal("Do returned nil error")
	}

	var statusErr *interfaces.StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("Do returned %T, want *interfaces.StatusError", err)
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

func TestDoReturnsPlainTextBodyForNonBadRequestStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("upstream unavailable"))
	}))
	defer server.Close()

	client := New(&interfaces.ClientOptions{APIKey: "test-api-key", Host: server.URL})
	if client == nil {
		t.Fatal("New returned nil")
	}

	req, err := client.SetupRequest(context.Background(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("SetupRequest returned error: %v", err)
	}

	err = client.Do(context.Background(), req, nil)
	if err == nil {
		t.Fatal("Do returned nil error")
	}

	got := err.Error()
	if !strings.Contains(got, "503 Service Unavailable") {
		t.Fatalf("error = %q, want status text", got)
	}
	if !strings.Contains(got, "upstream unavailable") {
		t.Fatalf("error = %q, want response body", got)
	}
}
