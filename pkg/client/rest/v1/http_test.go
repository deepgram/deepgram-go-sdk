// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	"net/http"
	"net/url"
	"testing"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func Test_NewHTTPClientProxy(t *testing.T) {
	proxyURL, err := url.Parse("http://proxy.example:8080")
	if err != nil {
		t.Fatalf("parse proxy URL: %v", err)
	}

	t.Run("configured", func(t *testing.T) {
		calls := 0
		client := NewHTTPClient(&interfaces.ClientOptions{
			Proxy: func(*http.Request) (*url.URL, error) {
				calls++
				return proxyURL, nil
			},
		})
		transport := client.Client.Transport.(*http.Transport)
		request, err := http.NewRequest(http.MethodGet, "https://api.deepgram.com/v1/listen", nil)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}

		got, err := transport.Proxy(request)
		if err != nil {
			t.Fatalf("call proxy: %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected proxy function to be called once, got %d", calls)
		}
		if got == nil {
			t.Fatal("expected a proxy URL")
		}
		if got.String() != proxyURL.String() {
			t.Fatalf("expected proxy URL %q, got %q", proxyURL, got)
		}
	})

	t.Run("unset", func(t *testing.T) {
		client := NewHTTPClient(&interfaces.ClientOptions{})
		transport := client.Client.Transport.(*http.Transport)
		if transport.Proxy != nil {
			t.Fatal("expected direct connections when no proxy is configured")
		}
	})
}
