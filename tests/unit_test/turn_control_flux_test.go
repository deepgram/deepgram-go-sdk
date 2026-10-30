// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	version "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/version"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
	listenv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2"
	listenv2ws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2/websocket"
)

// forceEndTurnWireFormat is the exact JSON text frame the SDK must emit for
// a ForceEndTurn control message.
const forceEndTurnWireFormat = `{"type":"ForceEndTurn"}`

// Test_FluxForceEndTurnSerialization verifies the ForceEndTurn control message
// serializes to the exact wire format {"type":"ForceEndTurn"}.
func Test_FluxForceEndTurnSerialization(t *testing.T) {
	msg := msginterfaces.ForceEndTurnMessage{
		Type: listenv2ws.MessageTypeForceEndTurn,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal ForceEndTurnMessage failed: %s", err)
	}

	want := forceEndTurnWireFormat
	if string(data) != want {
		t.Errorf("ForceEndTurnMessage wire format mismatch: got %s, want %s", data, want)
	}
}

// Test_FluxTurnInfoTriggerParsing verifies the trigger field on EndOfTurn
// TurnInfo messages parses for each documented value and stays empty when absent.
func Test_FluxTurnInfoTriggerParsing(t *testing.T) {
	for _, trigger := range []string{
		msginterfaces.TurnTriggerModel,
		msginterfaces.TurnTriggerManual,
		msginterfaces.TurnTriggerTimeout,
	} {
		trigger := trigger
		t.Run("trigger="+trigger, func(t *testing.T) {
			raw := `{
				"type": "TurnInfo",
				"request_id": "550e8400-e29b-41d4-a716-446655440000",
				"sequence_id": 42,
				"event": "EndOfTurn",
				"turn_index": 3,
				"audio_window_start": 1.5,
				"audio_window_end": 4.2,
				"transcript": "hello world",
				"words": [{"word": "hello", "confidence": 0.98}],
				"end_of_turn_confidence": 0.91,
				"trigger": "` + trigger + `"
			}`

			var resp msginterfaces.TurnInfoResponse
			if err := json.Unmarshal([]byte(raw), &resp); err != nil {
				t.Fatalf("unmarshal TurnInfoResponse failed: %s", err)
			}

			if resp.EventType != msginterfaces.TurnEventEndOfTurn {
				t.Errorf("expected event EndOfTurn, got %q", resp.EventType)
			}
			if resp.Trigger != trigger {
				t.Errorf("expected trigger %q, got %q", trigger, resp.Trigger)
			}
		})
	}

	t.Run("trigger absent leaves the field empty", func(t *testing.T) {
		raw := `{"type":"TurnInfo","event":"Update","turn_index":1,"transcript":"partial"}`

		var resp msginterfaces.TurnInfoResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			t.Fatalf("unmarshal TurnInfoResponse failed: %s", err)
		}
		if resp.Trigger != "" {
			t.Errorf("expected empty trigger, got %q", resp.Trigger)
		}
	})

	t.Run("unknown trigger values pass through as-is", func(t *testing.T) {
		raw := `{"type":"TurnInfo","event":"EndOfTurn","trigger":"some_future_value"}`

		var resp msginterfaces.TurnInfoResponse
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			t.Fatalf("unmarshal TurnInfoResponse failed: %s", err)
		}
		if resp.Trigger != "some_future_value" {
			t.Errorf("expected trigger to pass through open string, got %q", resp.Trigger)
		}
	})
}

// Test_FluxEotThresholdOne verifies eot_threshold=1.0 (suppress native end-of-turn
// detection) serializes onto the wss://.../v2/listen query string.
func Test_FluxEotThresholdOne(t *testing.T) {
	uri, err := version.GetFluxAPI(
		context.Background(),
		"",
		"",
		"",
		&interfaces.FluxTranscriptionOptions{
			Model:        "flux-general-en",
			Encoding:     "linear16",
			SampleRate:   16000,
			EotThreshold: 1.0,
		})
	if err != nil {
		t.Fatalf("GetFluxAPI should succeed, but got %s", err)
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		t.Fatalf("failed to parse generated URI %q: %s", uri, err)
	}

	raw := parsed.Query().Get("eot_threshold")
	got, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		t.Fatalf("eot_threshold %q is not a number (full URI: %s): %s", raw, uri, err)
	}
	if got != 1.0 {
		t.Errorf("expected eot_threshold=1.0 in query string, got %q (full URI: %s)", raw, uri)
	}
}

// Test_FluxConfigureThresholdOne verifies a mid-session Configure message carrying
// eot_threshold 1.0 serializes the full-suppression value on the wire.
func Test_FluxConfigureThresholdOne(t *testing.T) {
	msg := msginterfaces.ConfigureMessage{
		Type: listenv2ws.MessageTypeConfigure,
		Thresholds: &interfaces.FluxThresholds{
			EotThreshold: 1.0,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal ConfigureMessage failed: %s", err)
	}

	want := `{"type":"Configure","thresholds":{"eot_threshold":1}}`
	if string(data) != want {
		t.Errorf("ConfigureMessage wire format mismatch: got %s, want %s", data, want)
	}
}

// Test_FluxConfigureNumeralsSerialization verifies Configure uses a pointer so
// both enabled numerals and an explicit false value reach the wire unchanged.
func Test_FluxConfigureNumeralsSerialization(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value bool
		want  string
	}{
		{name: "true", value: true, want: `{"type":"Configure","numerals":true}`},
		{name: "false", value: false, want: `{"type":"Configure","numerals":false}`},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			msg := msginterfaces.ConfigureMessage{
				Type:     listenv2ws.MessageTypeConfigure,
				Numerals: &tc.value,
			}

			data, err := json.Marshal(msg)
			if err != nil {
				t.Fatalf("marshal ConfigureMessage failed: %s", err)
			}
			if string(data) != tc.want {
				t.Errorf("ConfigureMessage wire format mismatch: got %s, want %s", data, tc.want)
			}
		})
	}
}

func Test_FluxConfigureSerializationPreservesEmptySlices(t *testing.T) {
	for _, tc := range []struct {
		name string
		msg  msginterfaces.ConfigureMessage
		want string
	}{
		{
			name: "no-op",
			msg:  msginterfaces.ConfigureMessage{Type: listenv2ws.MessageTypeConfigure},
			want: `{"type":"Configure"}`,
		},
		{
			name: "empty keyterms",
			msg: msginterfaces.ConfigureMessage{
				Type:     listenv2ws.MessageTypeConfigure,
				Keyterms: []string{},
			},
			want: `{"type":"Configure","keyterms":[]}`,
		},
		{
			name: "empty language hints",
			msg: msginterfaces.ConfigureMessage{
				Type:          listenv2ws.MessageTypeConfigure,
				LanguageHints: []string{},
			},
			want: `{"type":"Configure","language_hints":[]}`,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.msg)
			if err != nil {
				t.Fatalf("marshal ConfigureMessage failed: %s", err)
			}
			if string(data) != tc.want {
				t.Errorf("ConfigureMessage wire format mismatch: got %s, want %s", data, tc.want)
			}
		})
	}
}

func Test_FluxConfigureOptionsCheck(t *testing.T) {
	if err := (*interfaces.FluxConfigureOptions)(nil).Check(); !errors.Is(err, interfacesv2.ErrOptionsRequired) {
		t.Errorf("nil Configure options: expected ErrOptionsRequired, got %v", err)
	}

	for _, tc := range []struct {
		name string
		opts interfaces.FluxConfigureOptions
	}{
		{name: "no-op"},
		{name: "empty keyterms", opts: interfaces.FluxConfigureOptions{Keyterms: []string{}}},
		{name: "empty language hints", opts: interfaces.FluxConfigureOptions{LanguageHints: []string{}}},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.opts.Check(); err != nil {
				t.Errorf("Configure options validation failed: %v", err)
			}
		})
	}
}

const fluxSafetyOptionsQuery = "encoding=linear16&model=flux-general-en&numerals=true&profanity_filter=true&redact=aggressive_numbers&sample_rate=16000"

func fluxSafetyOptions() *interfaces.FluxTranscriptionOptions {
	return &interfaces.FluxTranscriptionOptions{
		Model:           "flux-general-en",
		Encoding:        "linear16",
		SampleRate:      16000,
		ProfanityFilter: true,
		Numerals:        true,
		Redact:          interfaces.FluxRedactAggressiveNumbers,
	}
}

func assertFluxSafetyOptionsURL(t *testing.T, uri string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("GetURL failed: %s", err)
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		t.Fatalf("parse URL %q: %s", uri, err)
	}
	if parsed.Scheme != "wss" || parsed.Host != "api.deepgram.com" || parsed.Path != "/v2/listen" {
		t.Fatalf("unexpected Flux WebSocket endpoint: %s", uri)
	}
	if parsed.RawQuery != fluxSafetyOptionsQuery {
		t.Errorf("unexpected Flux safety options query:\n got: %s\nwant: %s", parsed.RawQuery, fluxSafetyOptionsQuery)
	}
}

func Test_FluxSafetyOptionsQueryCallback(t *testing.T) {
	client, err := listenv2.NewWSUsingCallback(
		context.Background(), MockAPIKey, &interfaces.ClientOptionsV2{}, fluxSafetyOptions(), nil)
	if err != nil {
		t.Fatalf("NewWSUsingCallback failed: %s", err)
	}

	uri, err := client.GetURL("")
	assertFluxSafetyOptionsURL(t, uri, err)
}

func Test_FluxSafetyOptionsQueryChannel(t *testing.T) {
	client, err := listenv2.NewWSUsingChan(
		context.Background(), MockAPIKey, &interfaces.ClientOptionsV2{}, fluxSafetyOptions(), nil)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %s", err)
	}

	uri, err := client.GetURL("")
	assertFluxSafetyOptionsURL(t, uri, err)
}

func Test_FluxSafetyOptionsValidateRedactModes(t *testing.T) {
	tests := []struct {
		name string
		mode interfaces.FluxRedact
		want error
	}{
		{name: "unset"},
		{name: "numbers", mode: interfaces.FluxRedactNumbers},
		{name: "aggressive_numbers", mode: interfaces.FluxRedactAggressiveNumbers},
		{name: "invalid", mode: interfaces.FluxRedact("pii"), want: interfacesv2.ErrFluxRedactInvalid},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (&interfaces.FluxTranscriptionOptions{Model: "flux-general-en", Redact: test.mode}).Check()
			if err != test.want {
				t.Errorf("Check() error = %v, want %v", err, test.want)
			}
		})
	}
}
