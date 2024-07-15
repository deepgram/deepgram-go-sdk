// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/gorilla/schema"
	"github.com/jarcoal/httpmock"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"

	utils "github.com/deepgram/deepgram-go-sdk/tests/utils"
)

const (
	model string = "2-general-nova"
)

const (
	FromURLSmartFormat = "Yep. I said it before, and I'll say it again. Life moves pretty fast. You don't stop and look around once in a while, you could miss it."
	FromURLSummarize   = "Yep. I said it before, and I'll say it again. Life moves pretty fast. You don't stop and look around once in a while, you could miss it."
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func Test_PrerecordedFromURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const preRecordedEndPoint = "https://api.deepgram.com/v1/listen"
	// const betaEndPoint = "https://beta.api.deepgram.com/v1/listen"

	// Specify query params that are acceptable. A nil means no check
	var acceptParams = map[string][]string{
		"model":   nil,
		"tier":    {"2-general-nova", "nova-2", "nova", "enhanced", "base"},
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
				"request_id": MockRequestID})
		}

		// Content checking
		var body map[string]any
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return httpmock.NewJsonResponse(400, map[string]any{
				"err_code":   "Bad Request",
				"err_msg":    "Content-type was application/json, but we could not process the JSON payload.",
				"request_id": MockRequestID})
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

				t.Errorf("unexpected value of key %v in prerecorded options: %v", k, v)
			}
		}

		// Based on query parameters, send Mock responses
		var optionStruct interfaces.PreRecordedTranscriptionOptions

		decoder := schema.NewDecoder()
		err = decoder.Decode(&optionStruct, r.URL.Query())
		if err != nil {
			t.Errorf("error decoding options: %s", err)
		}

		// save the options
		data, err := json.Marshal(optionStruct)
		if err != nil {
			t.Errorf("json.Marshal Err: %v", err)
		}

		sha256sum := fmt.Sprintf("%x", sha256.Sum256(data))
		filename := fmt.Sprintf("tests/response_data/%s-response.json", sha256sum)

		// Check if the file exists
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			result, err := utils.ReadMetadataString(filename)
			if err == nil {
				resp := httpmock.NewStringResponse(200, result)
				return resp, nil
			} else {
				t.Errorf("error reading response file: %s", err)
			}
		}

		return httpmock.NewStringResponse(404, ""), nil
	}

	// Register Handlers to endpoints
	httpmock.RegisterResponder("POST", preRecordedEndPoint, preRecordedFromURLHandler)
	// httpmock.RegisterResponder("POST", betaEndPoint, preRecordedFromURLHandler)

	t.Run("Test Basic PreRecordedFromURL", func(t *testing.T) {
		c := client.New(MockAPIKey, &interfaces.ClientOptions{})
		httpmock.ActivateNonDefault(&c.Client.HTTPClient.Client)
		dg := prerecorded.New(c)

		res, err := dg.FromURL(
			context.Background(),
			MockAudioURL,
			&interfaces.PreRecordedTranscriptionOptions{
				Model:       "nova-2",
				SmartFormat: true,
			})
		if err != nil {
			t.Errorf("should succeed, but got %s", err)
		}

		// check the response
		for _, value := range res.Metadata.ModelInfo {
			if strings.Compare(model, value.Name) != 0 {
				t.Errorf("%s: %s != %s", t.Name(), model, value.Name)
			}
		}
		transcript := res.Results.Channels[0].Alternatives[0].Transcript
		if strings.Compare(FromURLSmartFormat, transcript) != 0 {
			t.Errorf("%s: %s != %s", t.Name(), FromURLSmartFormat, transcript)
		}
	})

	t.Run("Test PreRecordedFromURL with summarize v2", func(t *testing.T) {
		c := client.New(MockAPIKey, &interfaces.ClientOptions{})
		httpmock.ActivateNonDefault(&c.Client.HTTPClient.Client)
		dg := prerecorded.New(c)
		res, err := dg.FromURL(
			context.Background(),
			MockAudioURL,
			&interfaces.PreRecordedTranscriptionOptions{
				Model:     "nova-2",
				Summarize: "v2",
			})
		if err != nil {
			t.Errorf("Summarize v2 should succeed, but got %s", err)
		}

		// check the response
		for _, value := range res.Metadata.ModelInfo {
			if strings.Compare(model, value.Name) != 0 {
				t.Errorf("%s: %s != %s", t.Name(), model, value.Name)
			}
		}
		summary := res.Results.Summary.Short
		if strings.Compare(FromURLSummarize, summary) != 0 {
			t.Errorf("%s: %s != %s", t.Name(), FromURLSummarize, summary)
		}
	})
}
