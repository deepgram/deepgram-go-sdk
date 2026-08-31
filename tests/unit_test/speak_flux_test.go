// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	speakrestv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/rest"
	speakmsg "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
	speakclientrest "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/rest"
	speakclientws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/websocket"
)

const (
	testFluxModel  = "flux-haley-en"
	testSpeechID   = "dg_sp_a1b2c3d4e5f6"
	speakV2URLPath = "/v2/speak"
	testAPIHost    = "api.deepgram.com"
)

// Test_FluxSpeakClientMessageSerialization verifies every client control message
// serializes to the exact wire format from the /v2/speak AsyncAPI spec.
func Test_FluxSpeakClientMessageSerialization(t *testing.T) {
	cases := []struct {
		name string
		msg  interface{}
		want string
	}{
		{
			name: "Speak",
			msg: speakmsg.SpeakMessage{
				Type: speakclientws.MessageTypeSpeak,
				Text: "Sure, I can help you cancel your subscription.",
			},
			want: `{"type":"Speak","text":"Sure, I can help you cancel your subscription."}`,
		},
		{
			name: "Flush",
			msg:  speakmsg.FlushMessage{Type: speakclientws.MessageTypeFlush},
			want: `{"type":"Flush"}`,
		},
		{
			name: "Interrupt without offset",
			msg:  speakmsg.InterruptMessage{Type: speakclientws.MessageTypeInterrupt},
			want: `{"type":"Interrupt"}`,
		},
		{
			name: "Interrupt with playback offset",
			msg: speakmsg.InterruptMessage{
				Type: speakclientws.MessageTypeInterrupt,
				PlaybackOffset: &speakmsg.PlaybackOffset{
					Type:  speakmsg.PlaybackOffsetTypeTimeMs,
					Value: 2340,
				},
			},
			want: `{"type":"Interrupt","playback_offset":{"type":"time_ms","value":2340}}`,
		},
		{
			name: "Configure with speed",
			msg:  speakmsg.ConfigureMessage{Type: speakclientws.MessageTypeConfigure, Speed: 1.05},
			want: `{"type":"Configure","speed":1.05}`,
		},
		{
			name: "Close",
			msg:  speakmsg.CloseMessage{Type: speakclientws.MessageTypeClose},
			want: `{"type":"Close"}`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.msg)
			if err != nil {
				t.Fatalf("marshal failed: %s", err)
			}
			if string(data) != tc.want {
				t.Errorf("wire format mismatch:\n got %s\nwant %s", data, tc.want)
			}
		})
	}
}

// Test_FluxSpeakServerEventParsing verifies every server event from the /v2/speak
// AsyncAPI spec parses into its typed response struct.
func Test_FluxSpeakServerEventParsing(t *testing.T) {
	t.Run("Connected", func(t *testing.T) {
		raw := `{"type":"Connected","request_id":"550e8400-e29b-41d4-a716-446655440000","model_name":"` + testFluxModel + `","model_version":"2026.06.01","model_uuids":["c0d1e2f3-a4b5-6789-abcd-ef0123456789"]}`
		var msg speakmsg.ConnectedResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.ModelName != testFluxModel || msg.ModelVersion != "2026.06.01" || len(msg.ModelUUIDs) != 1 {
			t.Errorf("unexpected Connected: %+v", msg)
		}
	})

	t.Run("SpeechStarted", func(t *testing.T) {
		raw := `{"type":"SpeechStarted","speech_id":"` + testSpeechID + `"}`
		var msg speakmsg.SpeechStartedResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.SpeechID != testSpeechID {
			t.Errorf("unexpected SpeechStarted: %+v", msg)
		}
	})

	t.Run("SpeechMetadata", func(t *testing.T) {
		raw := `{
			"type":"SpeechMetadata","speech_id":"` + testSpeechID + `",
			"audio_duration_ms":1875,"input_character_count":47,"billable_character_count":47,
			"controls_applied":{"pronunciations_applied":0,"breaks_applied":0,"pronunciation_warnings":0}
		}`
		var msg speakmsg.SpeechMetadataResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.AudioDurationMs != 1875 || msg.BillableCharacterCount != 47 {
			t.Errorf("unexpected SpeechMetadata: %+v", msg)
		}
	})

	t.Run("SpeechInterrupted with offset-derived text split", func(t *testing.T) {
		raw := `{
			"type":"SpeechInterrupted","audio_played_ms":2340,
			"text_spoken":"Sure, I can help you","text_remaining":" cancel your subscription.",
			"metadata":{"speech_id":"` + testSpeechID + `","audio_duration_ms":4200,
				"input_character_count":47,"billable_character_count":47,
				"controls_applied":{"pronunciations_applied":0,"breaks_applied":0,"pronunciation_warnings":0}}
		}`
		var msg speakmsg.SpeechInterruptedResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.AudioPlayedMs != 2340 || msg.TextSpoken == "" || msg.Metadata.SpeechID != testSpeechID {
			t.Errorf("unexpected SpeechInterrupted: %+v", msg)
		}
	})

	t.Run("SpeechInterrupted without offset omits text fields", func(t *testing.T) {
		raw := `{"type":"SpeechInterrupted","audio_played_ms":1200,
			"metadata":{"speech_id":"` + testSpeechID + `","audio_duration_ms":1200,
				"input_character_count":10,"billable_character_count":10,
				"controls_applied":{"pronunciations_applied":0,"breaks_applied":0,"pronunciation_warnings":0}}}`
		var msg speakmsg.SpeechInterruptedResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.TextSpoken != "" || msg.TextRemaining != "" {
			t.Errorf("expected empty text fields, got %+v", msg)
		}
	})

	t.Run("Flushed", func(t *testing.T) {
		raw := `{"type":"Flushed","speech_id":"` + testSpeechID + `"}`
		var msg speakmsg.FlushedResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.SpeechID != testSpeechID {
			t.Errorf("unexpected Flushed: %+v", msg)
		}
	})

	t.Run("SessionMetadata", func(t *testing.T) {
		raw := `{"type":"SessionMetadata","total_audio_duration_ms":15400,"total_input_character_count":312,"total_billable_character_count":310}`
		var msg speakmsg.SessionMetadataResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.TotalAudioDurationMs != 15400 || msg.TotalBillableCharacterCount != 310 {
			t.Errorf("unexpected SessionMetadata: %+v", msg)
		}
	})

	t.Run("ConfigureSuccess", func(t *testing.T) {
		raw := `{"type":"ConfigureSuccess","applied":{"speed":1.05}}`
		var msg speakmsg.ConfigureSuccessResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.Applied.Speed != 1.05 {
			t.Errorf("unexpected ConfigureSuccess: %+v", msg)
		}
	})

	t.Run("ConfigureFailure", func(t *testing.T) {
		raw := `{"type":"ConfigureFailure","code":"SPEED_OUT_OF_RANGE","field":"speed","value":3.5,"description":"speed must be between 0.85 and 1.15 in 0.05 increments"}`
		var msg speakmsg.ConfigureFailureResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.Code != "SPEED_OUT_OF_RANGE" || msg.Field != "speed" || msg.Value != 3.5 {
			t.Errorf("unexpected ConfigureFailure: %+v", msg)
		}
	})

	t.Run("Warning", func(t *testing.T) {
		raw := `{"type":"Warning","code":"NO_ACTIVE_SPEECH","description":"There is no active turn. The request will be ignored."}`
		var msg speakmsg.WarningResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.Code != "NO_ACTIVE_SPEECH" {
			t.Errorf("unexpected Warning: %+v", msg)
		}
	})

	t.Run("Error", func(t *testing.T) {
		raw := `{"type":"Error","code":"MESSAGE-0000","description":"The message could not be parsed."}`
		var msg speakmsg.FatalErrorResponse
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			t.Fatalf("unmarshal failed: %s", err)
		}
		if msg.Code != "MESSAGE-0000" {
			t.Errorf("unexpected Error: %+v", msg)
		}
	})
}

// Test_FluxSpeakURLBuilding verifies the /v2/speak URLs for both transports.
func Test_FluxSpeakURLBuilding(t *testing.T) {
	t.Run("websocket URL", func(t *testing.T) {
		uri, err := version.GetSpeakV2StreamAPI(context.Background(), "", "", "", &interfacesv2.SpeakV2WSOptions{
			Model:      testFluxModel,
			Encoding:   "linear16",
			SampleRate: 24000,
		})
		if err != nil {
			t.Fatalf("GetSpeakV2StreamAPI failed: %s", err)
		}
		parsed, err := url.Parse(uri)
		if err != nil {
			t.Fatalf("failed to parse URI %q: %s", uri, err)
		}
		if parsed.Scheme != "wss" || parsed.Host != testAPIHost || parsed.Path != speakV2URLPath {
			t.Errorf("unexpected URL: %s", uri)
		}
		if got := parsed.Query().Get("model"); got != testFluxModel {
			t.Errorf("expected model=flux-haley-en, got %q (full URI: %s)", got, uri)
		}
		if got := parsed.Query().Get("sample_rate"); got != "24000" {
			t.Errorf("expected sample_rate=24000, got %q (full URI: %s)", got, uri)
		}
	})

	t.Run("REST URL", func(t *testing.T) {
		uri, err := version.GetSpeakV2API(context.Background(), "", "", "", &interfacesv2.SpeakV2Options{
			Model:    testFluxModel,
			Encoding: "mp3",
		})
		if err != nil {
			t.Fatalf("GetSpeakV2API failed: %s", err)
		}
		parsed, err := url.Parse(uri)
		if err != nil {
			t.Fatalf("failed to parse URI %q: %s", uri, err)
		}
		if parsed.Scheme != "https" || parsed.Host != testAPIHost || parsed.Path != speakV2URLPath {
			t.Errorf("unexpected URL: %s", uri)
		}
		if got := parsed.Query().Get("encoding"); got != "mp3" {
			t.Errorf("expected encoding=mp3, got %q (full URI: %s)", got, uri)
		}
	})

	t.Run("model is required", func(t *testing.T) {
		opts := &interfacesv2.SpeakV2WSOptions{}
		if err := opts.Check(); err == nil {
			t.Error("expected Check() to fail without a model")
		}
		restOpts := &interfacesv2.SpeakV2Options{}
		if err := restOpts.Check(); err == nil {
			t.Error("expected Check() to fail without a model")
		}
	})
}

// Test_FluxSpeakRESTBatch exercises the batch transport end to end against a mock
// server: request method, path, query, body, and binary/JSON response handling.
func Test_FluxSpeakRESTBatch(t *testing.T) {
	audioBytes := []byte{0x49, 0x44, 0x33, 0x04, 0x00, 0x01, 0x02, 0x03} // fake mp3 bytes

	t.Run("synchronous request streams audio bytes and parses headers", func(t *testing.T) {
		var recordedMethod, recordedPath, recordedBody, recordedModel string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recordedMethod = r.Method
			recordedPath = r.URL.Path
			recordedModel = r.URL.Query().Get("model")
			body, _ := io.ReadAll(r.Body)
			recordedBody = string(body)
			w.Header().Set("Content-Type", "audio/mpeg")
			w.Header().Set("dg-request-id", "550e8400-e29b-41d4-a716-446655440000")
			w.Header().Set("dg-model-name", testFluxModel)
			w.Header().Set("dg-char-count", "47")
			if _, err := w.Write(audioBytes); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer srv.Close()

		dg := speakclientrest.New("testkey", &interfaces.ClientOptions{Host: srv.URL})
		client := speakrestv2.New(dg)

		var buf bytes.Buffer
		resp, err := client.ToFile(context.Background(), "Your appointment is confirmed for 3pm tomorrow.",
			&interfaces.SpeakV2Options{Model: testFluxModel}, &buf)
		if err != nil {
			t.Fatalf("ToFile failed: %s", err)
		}

		if recordedMethod != http.MethodPost {
			t.Errorf("expected POST, got %s", recordedMethod)
		}
		if recordedPath != speakV2URLPath {
			t.Errorf("expected path /v2/speak, got %s", recordedPath)
		}
		if recordedModel != testFluxModel {
			t.Errorf("expected model query param, got %q", recordedModel)
		}
		var reqBody map[string]string
		if err := json.Unmarshal([]byte(recordedBody), &reqBody); err != nil {
			t.Fatalf("request body is not valid JSON: %s", err)
		}
		if reqBody["text"] != "Your appointment is confirmed for 3pm tomorrow." {
			t.Errorf("unexpected request body: %s", recordedBody)
		}
		if !bytes.Equal(buf.Bytes(), audioBytes) {
			t.Errorf("audio bytes mismatch: got %v, want %v", buf.Bytes(), audioBytes)
		}
		if resp.RequestID != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("expected request-id header parsed, got %q", resp.RequestID)
		}
		if resp.ModelName != testFluxModel {
			t.Errorf("expected model-name header parsed, got %q", resp.ModelName)
		}
		if resp.Characters != 47 {
			t.Errorf("expected char-count 47, got %d", resp.Characters)
		}
	})

	t.Run("callback request returns the JSON acknowledgement bytes", func(t *testing.T) {
		var recordedCallback string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recordedCallback = r.URL.Query().Get("callback")
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"request_id":"550e8400-e29b-41d4-a716-446655440000"}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer srv.Close()

		dg := speakclientrest.New("testkey", &interfaces.ClientOptions{Host: srv.URL})
		client := speakrestv2.New(dg)

		var buf bytes.Buffer
		if _, err := client.ToFile(context.Background(), "Hello.",
			&interfaces.SpeakV2Options{Model: testFluxModel, Callback: "https://example.com/hook"}, &buf); err != nil {
			t.Fatalf("ToFile failed: %s", err)
		}

		if recordedCallback != "https://example.com/hook" {
			t.Errorf("expected callback query param, got %q", recordedCallback)
		}
		var ack map[string]string
		if err := json.Unmarshal(buf.Bytes(), &ack); err != nil {
			t.Fatalf("acknowledgement is not valid JSON: %s", err)
		}
		if ack["request_id"] == "" {
			t.Errorf("expected request_id in acknowledgement, got %s", buf.String())
		}
	})

	t.Run("400 response propagates as an error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{"err_code":"INVALID_MODEL","err_msg":"aura models are not accepted on /v2/speak"}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer srv.Close()

		dg := speakclientrest.New("testkey", &interfaces.ClientOptions{Host: srv.URL})
		client := speakrestv2.New(dg)

		var buf bytes.Buffer
		if _, err := client.ToFile(context.Background(), "Hello.",
			&interfaces.SpeakV2Options{Model: "aura-2-thalia-en"}, &buf); err == nil {
			t.Fatal("expected an error from a 400 response, got nil")
		}
	})

	t.Run("missing model is rejected client-side", func(t *testing.T) {
		dg := speakclientrest.New("testkey", &interfaces.ClientOptions{Host: "http://127.0.0.1:1"})
		client := speakrestv2.New(dg)

		var buf bytes.Buffer
		if _, err := client.ToFile(context.Background(), "Hello.", &interfaces.SpeakV2Options{}, &buf); err == nil {
			t.Fatal("expected an error for a missing model, got nil")
		}
	})
}
