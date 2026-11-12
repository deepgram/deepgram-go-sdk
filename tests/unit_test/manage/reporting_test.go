// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
	apiinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	manageclient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

const reportingProjectID = "9db43a4f-4f24-4c0c-a2fc-1c0f19f77869"

func reportingFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "fixtures", "manage", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}

func newReportingClient(t *testing.T, body string, check func(*http.Request)) *api.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		check(r)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	return api.New(manageclient.New("testkey", &interfaces.ClientOptions{Host: srv.URL}))
}

func Test_ManageReporting(t *testing.T) {
	falseValue := false

	t.Run("GetUsageBreakdown sends false filters in the query and decodes its fixture", func(t *testing.T) {
		client := newReportingClient(t, reportingFixture(t, "usage_breakdown.json"), func(r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if want := "/v1/projects/" + reportingProjectID + "/usage/breakdown"; r.URL.Path != want {
				t.Errorf("expected path %s, got %s", want, r.URL.Path)
			}
			if got := r.URL.Query().Get("diarize"); got != "false" {
				t.Errorf("expected diarize=false, got %q", got)
			}
			if got := r.URL.Query().Get("grouping"); got != "endpoint" {
				t.Errorf("expected grouping=endpoint, got %q", got)
			}
		})

		response, err := client.GetUsageBreakdown(context.Background(), reportingProjectID, &apiinterfaces.UsageBreakdownRequest{
			Grouping: apiinterfaces.UsageBreakdownGroupingEndpoint,
			Diarize:  &falseValue,
		})
		if err != nil {
			t.Fatalf("GetUsageBreakdown failed: %v", err)
		}
		if len(response.Results) != 1 || response.Results[0].Requests != 42 {
			t.Fatalf("unexpected usage results: %+v", response.Results)
		}
		if response.Results[0].Grouping.Accessor != nil {
			t.Errorf("expected null accessor, got %q", *response.Results[0].Grouping.Accessor)
		}
		if models := response.Results[0].Grouping.Models; len(models) != 2 || models[1] != nil {
			t.Errorf("expected nullable model entries, got %+v", models)
		}
	})

	t.Run("GetBillingBreakdown sends repeated grouping values and decodes its fixture", func(t *testing.T) {
		client := newReportingClient(t, reportingFixture(t, "billing_breakdown.json"), func(r *http.Request) {
			if want := "/v1/projects/" + reportingProjectID + "/billing/breakdown"; r.URL.Path != want {
				t.Errorf("expected path %s, got %s", want, r.URL.Path)
			}
			if got, want := r.URL.Query()["grouping"], []string{"deployment", "line_item"}; !reflect.DeepEqual(got, want) {
				t.Errorf("expected grouping %v, got %v", want, got)
			}
		})

		response, err := client.GetBillingBreakdown(context.Background(), reportingProjectID, &apiinterfaces.BillingBreakdownRequest{
			Grouping: []apiinterfaces.BillingBreakdownGrouping{
				apiinterfaces.BillingBreakdownGroupingDeployment,
				apiinterfaces.BillingBreakdownGroupingLineItem,
			},
		})
		if err != nil {
			t.Fatalf("GetBillingBreakdown failed: %v", err)
		}
		if len(response.Results) != 1 || response.Results[0].Dollars != 12.34 {
			t.Fatalf("unexpected billing results: %+v", response.Results)
		}
	})

	t.Run("GetBillingFields sends the reporting period and decodes its fixture", func(t *testing.T) {
		client := newReportingClient(t, reportingFixture(t, "billing_fields.json"), func(r *http.Request) {
			if want := "/v1/projects/" + reportingProjectID + "/billing/fields"; r.URL.Path != want {
				t.Errorf("expected path %s, got %s", want, r.URL.Path)
			}
			if got := r.URL.Query().Get("start"); got != "2026-01-01" {
				t.Errorf("expected start=2026-01-01, got %q", got)
			}
		})

		response, err := client.GetBillingFields(context.Background(), reportingProjectID, &apiinterfaces.BillingFieldsRequest{Start: "2026-01-01"})
		if err != nil {
			t.Fatalf("GetBillingFields failed: %v", err)
		}
		if response.LineItems["streaming::nova-3"] != "Nova-3 streaming" {
			t.Fatalf("unexpected line items: %+v", response.LineItems)
		}
	})

	t.Run("ListPurchases sends limit and decodes its fixture", func(t *testing.T) {
		client := newReportingClient(t, reportingFixture(t, "purchases.json"), func(r *http.Request) {
			if want := "/v1/projects/" + reportingProjectID + "/purchases"; r.URL.Path != want {
				t.Errorf("expected path %s, got %s", want, r.URL.Path)
			}
			if got := r.URL.Query().Get("limit"); got != "25" {
				t.Errorf("expected limit=25, got %q", got)
			}
		})

		response, err := client.ListPurchases(context.Background(), reportingProjectID, &apiinterfaces.PurchasesRequest{Limit: 25})
		if err != nil {
			t.Fatalf("ListPurchases failed: %v", err)
		}
		if len(response.Orders) != 1 || response.Orders[0].Created == nil || response.Orders[0].Amount != 100.5 {
			t.Fatalf("unexpected purchases: %+v", response.Orders)
		}
	})

	t.Run("reporting API errors are returned without a response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"err_code":"BAD_REQUEST","err_msg":"invalid reporting period"}`))
		}))
		defer srv.Close()
		client := api.New(manageclient.New("testkey", &interfaces.ClientOptions{Host: srv.URL}))
		response, err := client.GetUsageBreakdown(context.Background(), reportingProjectID, nil)
		if err == nil {
			t.Fatal("expected a server error, got nil")
		}
		if response != nil {
			t.Fatalf("expected no response on error, got %+v", response)
		}
	})
}

func Test_ManageReportingNilRequests(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		path    string
		call    func(*api.Client) error
	}{
		{
			name:    "GetUsageBreakdown",
			fixture: "usage_breakdown.json",
			path:    "/v1/projects/" + reportingProjectID + "/usage/breakdown",
			call: func(client *api.Client) error {
				_, err := client.GetUsageBreakdown(context.Background(), reportingProjectID, nil)
				return err
			},
		},
		{
			name:    "GetBillingBreakdown",
			fixture: "billing_breakdown.json",
			path:    "/v1/projects/" + reportingProjectID + "/billing/breakdown",
			call: func(client *api.Client) error {
				_, err := client.GetBillingBreakdown(context.Background(), reportingProjectID, nil)
				return err
			},
		},
		{
			name:    "GetBillingFields",
			fixture: "billing_fields.json",
			path:    "/v1/projects/" + reportingProjectID + "/billing/fields",
			call: func(client *api.Client) error {
				_, err := client.GetBillingFields(context.Background(), reportingProjectID, nil)
				return err
			},
		},
		{
			name:    "ListPurchases",
			fixture: "purchases.json",
			path:    "/v1/projects/" + reportingProjectID + "/purchases",
			call: func(client *api.Client) error {
				_, err := client.ListPurchases(context.Background(), reportingProjectID, nil)
				return err
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			client := newReportingClient(t, reportingFixture(t, test.fixture), func(r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != test.path {
					t.Errorf("expected path %s, got %s", test.path, r.URL.Path)
				}
				if r.URL.RawQuery != "" {
					t.Errorf("expected no query parameters, got %q", r.URL.RawQuery)
				}
			})

			if err := test.call(client); err != nil {
				t.Fatalf("%s with nil request failed: %v", test.name, err)
			}
		})
	}
}
