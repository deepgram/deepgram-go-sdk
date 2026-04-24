// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv2

// FluxTranscriptionOptions are the query parameters for wss://api.deepgram.com/v2/listen
// (the Deepgram Flux turn-based audio API).
type FluxTranscriptionOptions struct {
	// Model is required. Valid values: "flux-general-en", "flux-general-multi"
	Model string `json:"model" schema:"model"`
	// Encoding of the audio stream (e.g. "linear16", "opus"). Omit for containerized formats.
	Encoding string `json:"encoding,omitempty" schema:"encoding,omitempty"`
	// SampleRate of the audio in Hz (8000–48000). Omit for containerized formats.
	SampleRate int `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	// EagerEotThreshold (0.3–0.9): triggers early end-of-turn signal before silence completes.
	EagerEotThreshold float64 `json:"eager_eot_threshold,omitempty" schema:"eager_eot_threshold,omitempty"`
	// EotThreshold (0.5–0.9, default 0.7): confidence threshold for declaring EndOfTurn.
	EotThreshold float64 `json:"eot_threshold,omitempty" schema:"eot_threshold,omitempty"`
	// EotTimeoutMs (500–10000, default 5000): silence duration in ms before forcing EndOfTurn.
	EotTimeoutMs int `json:"eot_timeout_ms,omitempty" schema:"eot_timeout_ms,omitempty"`
	// Keyterm provides domain-specific vocabulary hints for improved recognition.
	Keyterm []string `json:"keyterm,omitempty" schema:"keyterm,omitempty"`
	// Tag attaches tracking metadata to the request.
	Tag []string `json:"tag,omitempty" schema:"tag,omitempty"`
	// MipOptOut opts the request out of Deepgram's model improvement program.
	MipOptOut bool `json:"mip_opt_out,omitempty" schema:"mip_opt_out,omitempty"`
}

// FluxThresholds holds the turn-detection thresholds used in Configure messages
// and ConfigureSuccess responses.
type FluxThresholds struct {
	EagerEotThreshold float64 `json:"eager_eot_threshold,omitempty"`
	EotThreshold      float64 `json:"eot_threshold,omitempty"`
	EotTimeoutMs      int     `json:"eot_timeout_ms,omitempty"`
}

// FluxConfigureOptions are the parameters for a mid-session Configure control message.
// Send via WSCallback.Configure() or WSChannel.Configure() to adjust turn-detection
// behavior without reconnecting.
type FluxConfigureOptions struct {
	Thresholds    *FluxThresholds `json:"thresholds,omitempty"`
	Keyterms      []string        `json:"keyterms,omitempty"`
	LanguageHints []string        `json:"language_hints,omitempty"`
}
