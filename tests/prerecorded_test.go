package deepgram_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/deepgram-devs/deepgram-go-sdk/deepgram"
	"github.com/jarcoal/httpmock"
)

func TestPrerecordedFromURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const preRecordedEndPoint = "https://api.deepgram.com/v1/listen"
	const betaEndPoint = "https://beta.api.deepgram.com/v1/listen"
	const betaHost = "beta.api.deepgram.com"

	// Specify query params that are acceptable. A nil means no check
	var acceptParams = map[string][]string{
		"model":   nil,
		"tier":    {"nova", "enhanced", "base"},
		"version": nil,
		"language": {"da", "en", "en-AU", "en-GB", "en-IN", "en-NZ",
			"en-US", "es", "es-419", "fr", "fr-CA", "hi", "hi-Latn", "id",
			"it", "ja", "ko", "nl", "pl", "pt", "pt-PT", "pt-BR", "ru", "sv",
			"tr", "uk", "zh-CN", "zh-TW"},
		"detect_language":  {"true", "false"},
		"punctuate":        {"true", "false"},
		"profanity_filter": {"true", "false"},
		"redact":           {"true", "false", "pci", "ssn", "numbers"},
		"diarize":          {"true", "false"},
		"diarize_version":  nil,
		"smart_format":     {"true", "false"},
		"multichannel":     {"true", "false"},
		"alternatives":     nil,
		"numerals":         {"true", "false"},
		"search":           nil,
		"replace":          nil,
		"callback":         nil,
		"keywords":         nil,
		"paragraphs":       {"true", "false"},
		"summarize":        {"true", "false", "v2"},
		"detect_topics":    {"true", "false"},
		"utterances":       {"true", "false"},
		"utt_split":        nil,
		"tag":              nil,
	}

	var preRecordedFromURLHandler = func(r *http.Request) (*http.Response, error) {
		// Content-type Checking
		if r.Header.Get("Content-type") != "application/json" {
			t.Errorf("expect content-type: application/json, got %s", r.Header.Get("Content-type"))
		}

		// Auth Checking
		authToken := fmt.Sprintf("token %s", MockAPIKey)
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

		// Param checking
		options := r.URL.Query()
		for k, values := range options {
			acceptValues, found := acceptParams[k]
			if !found {
				t.Errorf("unexpected query key in prerecorded options: %v", k)
			}

			if acceptValues == nil {
				continue
			}

		OUTER:
			for _, v := range values {
				for _, accept := range acceptValues {
					if accept == v {
						continue OUTER
					}
				}

				t.Errorf("unexpected value of key %v in prerecored options: %v", k, v)
			}
		}

		// Based on query parameters, send Mock responses
		var resp = &http.Response{}
		if options.Get("summarize") == "v2" {
			resp = httpmock.NewStringResponse(200, MockSummarizeV2Response)
		} else if options.Get("summarize") == "true" {
			resp = httpmock.NewStringResponse(200, MockSummarizeV1Response)
		} else {
			resp = httpmock.NewStringResponse(200, MockBasicPreRecordedResponse)
		}

		return resp, nil
	}

	// Register Handlers to endpoints
	httpmock.RegisterResponder("POST", preRecordedEndPoint, preRecordedFromURLHandler)
	httpmock.RegisterResponder("POST", betaEndPoint, preRecordedFromURLHandler)

	t.Run("Test Basic PreRecordedFromURL", func(t *testing.T) {
		dg := deepgram.NewClient(MockAPIKey)
		_, err := dg.PreRecordedFromURL(
			deepgram.UrlSource{Url: MockAudioURL},
			deepgram.PreRecordedTranscriptionOptions{})

		if err != nil {
			t.Errorf("should succeed, but got %s", err)
		}
	})

	t.Run("Test PreRecordedFromURL with summarize v1", func(t *testing.T) {
		dg := deepgram.NewClient(MockAPIKey)
		_, err := dg.PreRecordedFromURL(
			deepgram.UrlSource{Url: MockAudioURL},
			deepgram.PreRecordedTranscriptionOptions{
				Summarize: true,
			})

		if err != nil {
			t.Errorf("Summarize v1 should succeed, but got %s", err)
		}
	})

	t.Run("Test PreRecordedFromURL with summarize v2", func(t *testing.T) {
		dg := deepgram.NewClient(MockAPIKey).WithHost(betaHost)
		_, err := dg.PreRecordedFromURL(
			deepgram.UrlSource{Url: MockAudioURL},
			deepgram.PreRecordedTranscriptionOptions{
				Summarize: "v2",
			})

		if err != nil {
			t.Errorf("Summarize v2 should succeed, but got %s", err)
		}
	})
}
