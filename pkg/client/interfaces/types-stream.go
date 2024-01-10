// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

/*
LiveTranscriptionOptions contain all of the knobs and dials to control the live transcription
from the Deepgram API

Please see the documentation for live/streaming for more details:
https://developers.deepgram.com/reference/streaming
*/
type LiveTranscriptionOptions struct {
	Alternatives    int      `json:"alternatives,omitempty" url:"alternatives,omitempty"`
	Callback        string   `json:"callback,omitempty" url:"callback,omitempty"`
	CallbackMethod  string   `json:"callback_method,omitempty" url:"callback_method,omitempty"`
	Channels        int      `json:"channels,omitempty" url:"channels,omitempty"`
	Diarize         bool     `json:"diarize,omitempty" url:"diarize,omitempty"`
	DiarizeVersion  string   `json:"diarize_version,omitempty" url:"diarize_version,omitempty"`
	Encoding        string   `json:"encoding,omitempty" url:"encoding,omitempty"`
	Endpointing     string   `json:"endpointing,omitempty" url:"endpointing,omitempty"`
	Extra           string   `json:"extra,omitempty" url:"extra,omitempty"`
	FillerWords     bool     `json:"filler_words,omitempty" url:"filler_words,omitempty"`
	InterimResults  bool     `json:"interim_results,omitempty" url:"interim_results,omitempty"`
	Keywords        []string `json:"keywords,omitempty" url:"keywords,omitempty"`
	Language        string   `json:"language,omitempty" url:"language,omitempty"`
	Model           string   `json:"model,omitempty" url:"model,omitempty"`
	Multichannel    bool     `json:"multichannel,omitempty" url:"multichannel,omitempty"`
	ProfanityFilter bool     `json:"profanity_filter,omitempty" url:"profanity_filter,omitempty"`
	Punctuate       bool     `json:"punctuate,omitempty" url:"punctuate,omitempty"`
	Redact          []string `json:"redact,omitempty" url:"redact,omitempty"`
	Replace         []string `json:"replace,omitempty" url:"replace,omitempty"`
	SampleRate      int      `json:"sample_rate,omitempty" url:"sample_rate,omitempty"`
	Search          []string `json:"search,omitempty" url:"search,omitempty"`
	SmartFormat     bool     `json:"smart_format,omitempty" url:"smart_format,omitempty"`
	Tag             []string `json:"tag,omitempty" url:"tag,omitempty"`
	UtteranceEndMs  string   `json:"utterance_end_ms,omitempty" url:"utterance_end_ms,omitempty"`
	Version         string   `json:"version,omitempty" url:"version,omitempty"`
}
