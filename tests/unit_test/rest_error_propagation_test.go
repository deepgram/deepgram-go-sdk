// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authapi "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1"
	manageapi "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
	manageinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	authclient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/auth"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	manageclient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

const forbiddenErrorJSON = `{"err_code":"FORBIDDEN","err_msg":"Insufficient permissions."}`

// newForbiddenServer starts an httptest server that answers every request with
// a 403 and a Deepgram error body.
func newForbiddenServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		if _, err := w.Write([]byte(forbiddenErrorJSON)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// assertForbidden checks that an API error reaches the caller as a StatusError.
func assertForbidden(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error for a 403 response, got nil")
	}
	var statusErr *interfaces.StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected *interfaces.StatusError, got %T: %v", err, err)
	}
	if statusErr.Resp == nil || statusErr.Resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status %d, got %+v", http.StatusForbidden, statusErr.Resp)
	}
	if statusErr.DeepgramError == nil || statusErr.DeepgramError.ErrCode != "FORBIDDEN" {
		t.Fatalf("expected err_code FORBIDDEN, got %+v", statusErr.DeepgramError)
	}
}

func Test_ManageAPIErrorsArePropagated(t *testing.T) {
	srv := newForbiddenServer(t)
	dg := manageapi.New(manageclient.New("testkey", &interfaces.ClientOptions{Host: srv.URL}))
	ctx := context.Background()

	calls := map[string]func() error{
		"ListBalances": func() error { _, err := dg.ListBalances(ctx, testProjectID); return err },
		"GetBalance":   func() error { _, err := dg.GetBalance(ctx, testProjectID, "balance"); return err },
		"ListKeys":     func() error { _, err := dg.ListKeys(ctx, testProjectID); return err },
		"GetKey":       func() error { _, err := dg.GetKey(ctx, testProjectID, "key"); return err },
		"CreateKey": func() error {
			_, err := dg.CreateKey(ctx, testProjectID, &manageinterfaces.KeyCreateRequest{Comment: "c", Scopes: []string{"member"}})
			return err
		},
		"DeleteKey":    func() error { _, err := dg.DeleteKey(ctx, testProjectID, "key"); return err },
		"ListMembers":  func() error { _, err := dg.ListMembers(ctx, testProjectID); return err },
		"RemoveMember": func() error { _, err := dg.RemoveMember(ctx, testProjectID, testMemberID); return err },
		"ListModels":   func() error { _, err := dg.ListModels(ctx, nil); return err },
		"GetModel":     func() error { _, err := dg.GetModel(ctx, "model"); return err },
		"ListProjectModels": func() error {
			_, err := dg.ListProjectModels(ctx, testProjectID, nil)
			return err
		},
		"GetProjectModel": func() error { _, err := dg.GetProjectModel(ctx, testProjectID, "model"); return err },
		"ListProjects":    func() error { _, err := dg.ListProjects(ctx); return err },
		"GetProject":      func() error { _, err := dg.GetProject(ctx, testProjectID); return err },
		"UpdateProject": func() error {
			_, err := dg.UpdateProject(ctx, testProjectID, &manageinterfaces.ProjectUpdateRequest{Name: "n"})
			return err
		},
		"DeleteProject":   func() error { _, err := dg.DeleteProject(ctx, testProjectID); return err },
		"GetMemberScopes": func() error { _, err := dg.GetMemberScopes(ctx, testProjectID, testMemberID); return err },
		"UpdateMemberScopes": func() error {
			_, err := dg.UpdateMemberScopes(ctx, testProjectID, testMemberID, &manageinterfaces.ScopeUpdateRequest{Scope: "member"})
			return err
		},
		"ListRequests": func() error {
			_, err := dg.ListRequests(ctx, testProjectID, &manageinterfaces.UsageListRequest{})
			return err
		},
		"GetRequest": func() error { _, err := dg.GetRequest(ctx, testProjectID, "request"); return err },
		"GetFields": func() error {
			_, err := dg.GetFields(ctx, testProjectID, &manageinterfaces.UsageListRequest{})
			return err
		},
		"GetUsage": func() error {
			_, err := dg.GetUsage(ctx, testProjectID, &manageinterfaces.UsageRequest{})
			return err
		},
	}

	for name, call := range calls {
		call := call
		t.Run(name, func(t *testing.T) {
			assertForbidden(t, call())
		})
	}
}

func Test_GrantTokenErrorIsPropagated(t *testing.T) {
	srv := newForbiddenServer(t)
	dg := authapi.New(authclient.New("testkey", &interfaces.ClientOptions{Host: srv.URL}))

	resp, err := dg.GrantToken(context.Background(), nil)
	assertForbidden(t, err)
	if resp != nil && resp.AccessToken != "" {
		t.Fatalf("expected no access token, got %q", resp.AccessToken)
	}
}
