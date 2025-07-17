// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"encoding/json"
	"testing"

	authInterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1/interfaces"
)

// TestGrantTokenRequest tests the GrantTokenRequest struct
func TestGrantTokenRequest(t *testing.T) {
	t.Run("Test_GrantTokenRequest_struct_creation", func(t *testing.T) {
		ttlSeconds := 60
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: &ttlSeconds,
		}

		if req.TTLSeconds == nil {
			t.Error("TTLSeconds should not be nil")
		}

		if *req.TTLSeconds != 60 {
			t.Errorf("Expected TTLSeconds to be 60, got %d", *req.TTLSeconds)
		}
	})

	t.Run("Test_GrantTokenRequest_nil_ttl", func(t *testing.T) {
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: nil,
		}

		if req.TTLSeconds != nil {
			t.Error("TTLSeconds should be nil")
		}
	})

	t.Run("Test_GrantTokenRequest_zero_ttl", func(t *testing.T) {
		ttlSeconds := 0
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: &ttlSeconds,
		}

		if req.TTLSeconds == nil {
			t.Error("TTLSeconds should not be nil")
		}

		if *req.TTLSeconds != 0 {
			t.Errorf("Expected TTLSeconds to be 0, got %d", *req.TTLSeconds)
		}
	})

	t.Run("Test_GrantTokenRequest_max_ttl", func(t *testing.T) {
		ttlSeconds := 3600
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: &ttlSeconds,
		}

		if req.TTLSeconds == nil {
			t.Error("TTLSeconds should not be nil")
		}

		if *req.TTLSeconds != 3600 {
			t.Errorf("Expected TTLSeconds to be 3600, got %d", *req.TTLSeconds)
		}
	})
}

// TestGrantTokenRequestJSON tests JSON marshaling and unmarshaling
func TestGrantTokenRequestJSON(t *testing.T) {
	t.Run("Test_GrantTokenRequest_JSON_marshaling", func(t *testing.T) {
		ttlSeconds := 60
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: &ttlSeconds,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			t.Errorf("Failed to marshal GrantTokenRequest: %v", err)
		}

		expectedJSON := `{"ttl_seconds":60}`
		if string(jsonData) != expectedJSON {
			t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(jsonData))
		}
	})

	t.Run("Test_GrantTokenRequest_JSON_marshaling_nil", func(t *testing.T) {
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: nil,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			t.Errorf("Failed to marshal GrantTokenRequest: %v", err)
		}

		expectedJSON := `{}`
		if string(jsonData) != expectedJSON {
			t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(jsonData))
		}
	})

	t.Run("Test_GrantTokenRequest_JSON_unmarshaling", func(t *testing.T) {
		jsonData := `{"ttl_seconds":120}`

		var req authInterfaces.GrantTokenRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		if err != nil {
			t.Errorf("Failed to unmarshal GrantTokenRequest: %v", err)
		}

		if req.TTLSeconds == nil {
			t.Error("TTLSeconds should not be nil after unmarshaling")
		}

		if *req.TTLSeconds != 120 {
			t.Errorf("Expected TTLSeconds to be 120, got %d", *req.TTLSeconds)
		}
	})

	t.Run("Test_GrantTokenRequest_JSON_unmarshaling_empty", func(t *testing.T) {
		jsonData := `{}`

		var req authInterfaces.GrantTokenRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		if err != nil {
			t.Errorf("Failed to unmarshal empty GrantTokenRequest: %v", err)
		}

		if req.TTLSeconds != nil {
			t.Error("TTLSeconds should be nil after unmarshaling empty JSON")
		}
	})
}

// TestGrantTokenRequestValidation tests validation scenarios
func TestGrantTokenRequestValidation(t *testing.T) {
	t.Run("Test_GrantTokenRequest_minimum_value", func(t *testing.T) {
		ttlSeconds := 1
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: &ttlSeconds,
		}

		if *req.TTLSeconds != 1 {
			t.Errorf("Expected TTLSeconds to be 1, got %d", *req.TTLSeconds)
		}
	})

	t.Run("Test_GrantTokenRequest_maximum_value", func(t *testing.T) {
		ttlSeconds := 3600
		req := authInterfaces.GrantTokenRequest{
			TTLSeconds: &ttlSeconds,
		}

		if *req.TTLSeconds != 3600 {
			t.Errorf("Expected TTLSeconds to be 3600, got %d", *req.TTLSeconds)
		}
	})

	t.Run("Test_GrantTokenRequest_typical_values", func(t *testing.T) {
		testCases := []int{30, 60, 300, 600, 1800, 3600}

		for _, ttl := range testCases {
			req := authInterfaces.GrantTokenRequest{
				TTLSeconds: &ttl,
			}

			if *req.TTLSeconds != ttl {
				t.Errorf("Expected TTLSeconds to be %d, got %d", ttl, *req.TTLSeconds)
			}
		}
	})
}
