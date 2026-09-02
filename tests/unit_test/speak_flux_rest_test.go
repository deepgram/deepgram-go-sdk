// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	speakrestv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
	speakclientrest "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/rest"
)

const testAsyncRequestID = "550e8400-e29b-41d4-a716-446655440000"

// newFluxRESTClient builds a Flux TTS batch API client pointed at srvURL.
func newFluxRESTClient(t *testing.T, srvURL string) *speakrestv2.Client {
	t.Helper()
	dg, err := speakclientrest.New("testkey", &interfaces.ClientOptions{Host: srvURL})
	if err != nil {
		t.Fatalf("REST client constructor failed: %s", err)
	}
	return speakrestv2.New(dg)
}

// Test_FluxSpeakToSaveCallbackRejected verifies callback (asynchronous) requests
// cannot masquerade as saved audio (Greg's B4 on PR #351): ToSave rejects them
// without creating a file, and ToAsync returns the typed acknowledgement.
func Test_FluxSpeakToSaveCallbackRejected(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"request_id":"` + testAsyncRequestID + `"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer srv.Close()

	client := newFluxRESTClient(t, srv.URL)
	dir := t.TempDir()

	t.Run("ToSave rejects callback options and creates no file", func(t *testing.T) {
		filename := filepath.Join(dir, "audio.mp3")
		_, err := client.ToSave(context.Background(), filename, "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel, Callback: "https://example.com/hook"})
		if !errors.Is(err, speakrestv2.ErrCallbackNotSupported) {
			t.Fatalf("expected ErrCallbackNotSupported, got: %v", err)
		}
		if requests != 0 {
			t.Errorf("no HTTP request should be sent, got %d", requests)
		}
		if _, statErr := os.Stat(filename); !os.IsNotExist(statErr) {
			t.Errorf("ToSave must not create a file for callback requests")
		}
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Errorf("no temp files may be left behind, found %v", entries)
		}
	})

	t.Run("ToAsync returns the typed acknowledgement", func(t *testing.T) {
		resp, err := client.ToAsync(context.Background(), "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel, Callback: "https://example.com/hook"})
		if err != nil {
			t.Fatalf("ToAsync failed: %s", err)
		}
		if resp.RequestID != testAsyncRequestID {
			t.Errorf("expected the acknowledged request_id, got %q", resp.RequestID)
		}
	})

	t.Run("ToAsync requires a callback URL", func(t *testing.T) {
		_, err := client.ToAsync(context.Background(), "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel})
		if !errors.Is(err, speakrestv2.ErrCallbackRequired) {
			t.Fatalf("expected ErrCallbackRequired, got: %v", err)
		}
	})
}

// Test_FluxSpeakToSavePreservesExistingFile verifies failed requests leave an
// existing destination untouched (Greg's S1 on PR #351): validation runs before
// the file is opened, and the audio lands via a temp file renamed only on success.
func Test_FluxSpeakToSavePreservesExistingFile(t *testing.T) {
	const precious = "previously synthesized audio"

	newDest := func(t *testing.T) string {
		t.Helper()
		dest := filepath.Join(t.TempDir(), "audio.mp3")
		if err := os.WriteFile(dest, []byte(precious), 0o600); err != nil {
			t.Fatalf("seeding destination failed: %s", err)
		}
		return dest
	}
	assertUntouched := func(t *testing.T, dest string) {
		t.Helper()
		got, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("destination unreadable after failure: %s", err)
		}
		if string(got) != precious {
			t.Errorf("existing file was clobbered: got %q", got)
		}
		if entries, _ := os.ReadDir(filepath.Dir(dest)); len(entries) != 1 {
			t.Errorf("temp files left behind: %v", entries)
		}
	}

	t.Run("validation failure leaves the file untouched", func(t *testing.T) {
		client := newFluxRESTClient(t, "http://127.0.0.1:1")
		dest := newDest(t)
		if _, err := client.ToSave(context.Background(), dest, "Hello.", &interfaces.SpeakV2Options{}); err == nil {
			t.Fatal("expected a validation error for a missing model")
		}
		assertUntouched(t, dest)
	})

	t.Run("HTTP failure leaves the file untouched", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := newFluxRESTClient(t, srv.URL)
		dest := newDest(t)
		if _, err := client.ToSave(context.Background(), dest, "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel}); err == nil {
			t.Fatal("expected an error from a 500 response")
		}
		assertUntouched(t, dest)
	})

	t.Run("success replaces the file atomically", func(t *testing.T) {
		audio := []byte{0x49, 0x44, 0x33, 0x04}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "audio/mpeg")
			if _, err := w.Write(audio); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer srv.Close()

		client := newFluxRESTClient(t, srv.URL)
		dest := newDest(t)
		resp, err := client.ToSave(context.Background(), dest, "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel})
		if err != nil {
			t.Fatalf("ToSave failed: %s", err)
		}
		if resp.Filename != dest {
			t.Errorf("expected Filename %q, got %q", dest, resp.Filename)
		}
		got, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("reading saved audio failed: %s", err)
		}
		if !bytes.Equal(got, audio) {
			t.Errorf("saved audio mismatch: got %v, want %v", got, audio)
		}
		if entries, _ := os.ReadDir(filepath.Dir(dest)); len(entries) != 1 {
			t.Errorf("temp files left behind: %v", entries)
		}
	})
}

// Test_FluxSpeakNilGuards verifies plausible first-run mistakes return errors
// instead of panicking (Greg's S2 on PR #351).
func Test_FluxSpeakNilGuards(t *testing.T) {
	t.Run("REST methods reject nil options", func(t *testing.T) {
		client := newFluxRESTClient(t, "http://127.0.0.1:1")
		var buf interfaces.RawResponse
		if _, err := client.ToStream(context.Background(), "Hello.", nil, &buf); !errors.Is(err, interfacesv2.ErrOptionsRequired) {
			t.Errorf("ToStream: expected ErrOptionsRequired, got: %v", err)
		}
		if _, err := client.ToFile(context.Background(), "Hello.", nil, &buf); !errors.Is(err, interfacesv2.ErrOptionsRequired) {
			t.Errorf("ToFile: expected ErrOptionsRequired, got: %v", err)
		}
		if _, err := client.ToSave(context.Background(), filepath.Join(t.TempDir(), "a.mp3"), "Hello.", nil); !errors.Is(err, interfacesv2.ErrOptionsRequired) {
			t.Errorf("ToSave: expected ErrOptionsRequired, got: %v", err)
		}
	})

	t.Run("REST methods reject a nil transport client", func(t *testing.T) {
		client := speakrestv2.New(nil)
		var buf interfaces.RawResponse
		if _, err := client.ToStream(context.Background(), "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel}, &buf); !errors.Is(err, speakrestv2.ErrNilClient) {
			t.Errorf("expected ErrNilClient, got: %v", err)
		}
	})

	t.Run("REST constructor returns an error without credentials", func(t *testing.T) {
		t.Setenv("DEEPGRAM_API_KEY", "")
		t.Setenv("DEEPGRAM_ACCESS_TOKEN", "")
		dg, err := speakclientrest.New("", &interfaces.ClientOptions{})
		if err == nil {
			t.Fatal("expected an error when no credentials are available")
		}
		if dg != nil {
			t.Error("expected a nil client alongside the error")
		}
	})
}

// Test_FluxSpeakRESTErrorDetails verifies every non-success HTTP status preserves
// the Deepgram error payload (Greg's S8 on PR #351), not just HTTP 400.
func Test_FluxSpeakRESTErrorDetails(t *testing.T) {
	for _, status := range []int{
		http.StatusUnauthorized,          // 401
		http.StatusForbidden,             // 403
		http.StatusRequestEntityTooLarge, // 413
		http.StatusTooManyRequests,       // 429
		http.StatusInternalServerError,   // 500
	} {
		status := status
		t.Run(http.StatusText(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				body := `{"err_code":"TEST_CODE","err_msg":"details for status ` + http.StatusText(status) + `"}`
				if _, err := w.Write([]byte(body)); err != nil {
					t.Errorf("failed to write response: %v", err)
				}
			}))
			defer srv.Close()

			client := newFluxRESTClient(t, srv.URL)
			var buf interfaces.RawResponse
			_, err := client.ToStream(context.Background(), "Hello.",
				&interfaces.SpeakV2Options{Model: testFluxModel}, &buf)
			if err == nil {
				t.Fatalf("expected an error for HTTP %d", status)
			}

			var statusErr *interfacesv1.StatusError
			if !errors.As(err, &statusErr) {
				t.Fatalf("expected a StatusError, got %T: %v", err, err)
			}
			if statusErr.Resp.StatusCode != status {
				t.Errorf("expected status %d, got %d", status, statusErr.Resp.StatusCode)
			}
			if statusErr.DeepgramError == nil {
				t.Fatal("expected the Deepgram error payload to be preserved")
			}
			if statusErr.DeepgramError.ErrCode != "TEST_CODE" {
				t.Errorf("expected err_code TEST_CODE, got %q", statusErr.DeepgramError.ErrCode)
			}
		})
	}

	t.Run("non-JSON error bodies surface raw details", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			if _, err := w.Write([]byte("upstream unavailable")); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer srv.Close()

		client := newFluxRESTClient(t, srv.URL)
		var buf interfaces.RawResponse
		_, err := client.ToStream(context.Background(), "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel}, &buf)
		if err == nil {
			t.Fatal("expected an error for HTTP 502")
		}
		if got := err.Error(); !bytes.Contains([]byte(got), []byte("upstream unavailable")) {
			t.Errorf("expected the raw body in the error, got: %s", got)
		}
	})
}
