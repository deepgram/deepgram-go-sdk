package deepgram_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/deepgram-devs/deepgram-go-sdk/deepgram"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func TestPrerecordedFromURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var preRecordedHandler = func(r *http.Request) (*http.Response, error) {
		// Auth Checking
		authToken := fmt.Sprintf("token %s", MockApiKey)
		if r.Header.Get("Authorization") != authToken {
			return httpmock.NewJsonResponse(401, map[string]any{
				"err_code":   "INVALID_AUTH",
				"err_msg":    "Invalid credentials.",
				"request_id": MockRequestId})
		}

		// Content checking
		var body map[string]any
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return httpmock.NewJsonResponse(400, map[string]any{
				"err_code":   "Bad Request",
				"err_msg":    "Content-type was application/json, but we could not process the JSON payload.",
				"request_id": MockRequestId})
		}

		return httpmock.NewStringResponse(200, MockPrerecordedResponseJSON), nil

	}

	// Register Handlers to endpoints
	httpmock.RegisterResponder("POST", MockEndPoint, preRecordedHandler)

	t.Run("Test Basic PreRecordedFromURL", func(t *testing.T) {
		dg := deepgram.NewClient(MockApiKey).WithHost(MockApiDomain)
		resp, err := dg.PreRecordedFromURL(MockUrlSource, MockPrerecordedOptions)

		if err != nil {
			t.Errorf("should succeed, but got %s", err)
		}

		if !cmp.Equal(resp, MockPrerecordedResponse) {
			t.Errorf("results does not match, expected %+v, got %+v", MockPrerecordedResponse, resp)
		}

	})
}
