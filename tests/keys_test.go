package deepgram_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/deepgram-devs/deepgram-go-sdk/deepgram"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func TestKeys(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var keyHandler = func(r *http.Request) (*http.Response, error) {
		// Auth Checking
		authToken := fmt.Sprintf("token %s", MockApiKey)
		if r.Header.Get("Authorization") != authToken {
			return httpmock.NewJsonResponse(401, map[string]any{
				"err_code":   "INVALID_AUTH",
				"err_msg":    "Invalid credentials.",
				"request_id": MockRequestId})
		}

		return httpmock.NewJsonResponse(200, MockListKeys)

	}

	// Register Handlers to endpoints
	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://%s/v1/projects/%s/keys", MockApiDomain, MockProjectId),
		keyHandler)

	t.Run("Test List Keys should success", func(t *testing.T) {
		dg := deepgram.NewClient(MockApiKey).WithHost(MockApiDomain)
		resp, err := dg.ListKeys(MockProjectId)

		if err != nil {
			t.Errorf("should succeed, but got %s", err)
		}

		if !cmp.Equal(MockListKeys, resp) {
			t.Errorf("results does not match, expected %+v, got %+v", MockListKeys, resp)
		}
	})

}
