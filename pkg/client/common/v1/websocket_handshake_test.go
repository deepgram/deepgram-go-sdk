// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	commoninterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

type handshakeErrorRouter struct {
	errors chan *commoninterfaces.ErrorResponse
}

func (r *handshakeErrorRouter) Open(*commoninterfaces.OpenResponse) error { return nil }
func (r *handshakeErrorRouter) Message([]byte) error                      { return nil }
func (r *handshakeErrorRouter) Binary([]byte) error                       { return nil }
func (r *handshakeErrorRouter) Close(*commoninterfaces.CloseResponse) error {
	return nil
}
func (r *handshakeErrorRouter) Error(response *commoninterfaces.ErrorResponse) error {
	r.errors <- response
	return nil
}

type handshakeErrorHandler struct {
	url string
}

func (h *handshakeErrorHandler) GetURL(string) (string, error) { return h.url, nil }
func (h *handshakeErrorHandler) ProcessMessage(int, []byte) error {
	return nil
}
func (h *handshakeErrorHandler) ProcessError(error) error { return nil }
func (h *handshakeErrorHandler) Start()                   {}
func (h *handshakeErrorHandler) Finish()                  {}
func (h *handshakeErrorHandler) GetCloseMsg() []byte      { return nil }

func Test_WSClientHandshakeErrors(t *testing.T) {
	tests := []struct {
		name            string
		body            string
		contentType     string
		wantCode        string
		wantMessage     string
		wantDescription string
	}{
		{
			name:            "Deepgram JSON",
			body:            `{"err_code":"INVALID_QUERY","err_msg":"invalid endpointing","description":"endpointing must be an integer"}`,
			contentType:     "application/json",
			wantCode:        "INVALID_QUERY",
			wantMessage:     "invalid endpointing",
			wantDescription: "endpointing must be an integer",
		},
		{
			name:            "plain text",
			body:            "endpointing must be an integer",
			contentType:     "text/plain",
			wantCode:        "400 Bad Request",
			wantMessage:     "400 Bad Request",
			wantDescription: "endpointing must be an integer",
		},
		{
			name:            "malformed JSON",
			body:            `{"err_code":`,
			contentType:     "application/json",
			wantCode:        "400 Bad Request",
			wantMessage:     "400 Bad Request",
			wantDescription: `{"err_code":`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts.Add(1)
				w.Header().Set("Content-Type", test.contentType)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			handler := &handshakeErrorHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
			router := &handshakeErrorRouter{errors: make(chan *commoninterfaces.ErrorResponse, 1)}
			var handlerInterface commoninterfaces.WebSocketHandler = handler
			var routerInterface commoninterfaces.Router = router
			client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

			if client.ConnectWithCancel(ctx, cancel, 3) {
				t.Fatal("ConnectWithCancel succeeded after a 400 response")
			}
			if got := attempts.Load(); got != 1 {
				t.Errorf("handshake attempts = %d, want 1", got)
			}

			select {
			case response := <-router.errors:
				if response.ErrCode != test.wantCode {
					t.Errorf("ErrCode = %q, want %q", response.ErrCode, test.wantCode)
				}
				if test.wantCode == "INVALID_QUERY" && response.Variant != "400 Bad Request" {
					t.Errorf("Variant = %q, want %q", response.Variant, "400 Bad Request")
				}
				if response.ErrMsg != test.wantMessage {
					t.Errorf("ErrMsg = %q, want %q", response.ErrMsg, test.wantMessage)
				}
				if response.Description != test.wantDescription {
					t.Errorf("Description = %q, want %q", response.Description, test.wantDescription)
				}
			default:
				t.Fatal("expected a handshake error routed to the client handler")
			}
		})
	}
}

func Test_WSClientRetriesTransientHandshakeErrors(t *testing.T) {
	for _, statusCode := range []int{http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusInternalServerError} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := attempts.Add(1)
				w.Header().Set("Content-Type", "application/json")
				if attempt == 1 {
					w.WriteHeader(statusCode)
					_, _ = w.Write([]byte(`{"err_code":"TRANSIENT","err_msg":"try again"}`))
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"err_code":"INVALID_QUERY","err_msg":"invalid endpointing"}`))
			}))
			defer server.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			handler := &handshakeErrorHandler{url: "ws" + strings.TrimPrefix(server.URL, "http")}
			router := &handshakeErrorRouter{errors: make(chan *commoninterfaces.ErrorResponse, 1)}
			var handlerInterface commoninterfaces.WebSocketHandler = handler
			var routerInterface commoninterfaces.Router = router
			client := NewWS(ctx, cancel, "test-api-key", &interfaces.ClientOptions{}, &handlerInterface, &routerInterface)

			if client.ConnectWithCancel(ctx, cancel, 3) {
				t.Fatal("ConnectWithCancel succeeded after a final 400 response")
			}
			if got := attempts.Load(); got != 2 {
				t.Errorf("handshake attempts = %d, want 2", got)
			}

			select {
			case response := <-router.errors:
				if response.ErrCode != "INVALID_QUERY" {
					t.Errorf("ErrCode = %q, want %q", response.ErrCode, "INVALID_QUERY")
				}
				if response.Variant != "400 Bad Request" {
					t.Errorf("Variant = %q, want %q", response.Variant, "400 Bad Request")
				}
			default:
				t.Fatal("expected the final handshake error routed to the client handler")
			}
		})
	}
}
