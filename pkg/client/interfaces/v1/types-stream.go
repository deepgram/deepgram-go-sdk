// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

/*
LiveTranscriptionOptions contain all of the knobs and dials to control the live transcription
from the Deepgram API

Please see the live/streaming documentation for more details:
https://developers.deepgram.com/reference/streaming
*/
type LiveTranscriptionOptions struct {
	Alternatives    int      `json:"alternatives,omitempty" schema:"alternatives,omitempty"`
	Callback        string   `json:"callback,omitempty" schema:"callback,omitempty"`
	CallbackMethod  string   `json:"callback_method,omitempty" schema:"callback_method,omitempty"`
	Channels        int      `json:"channels,omitempty" schema:"channels,omitempty"`
	Diarize         bool     `json:"diarize,omitempty" schema:"diarize,omitempty"`
	DiarizeVersion  string   `json:"diarize_version,omitempty" schema:"diarize_version,omitempty"`
	Dictation       bool     `json:"dictation,omitempty" schema:"dictation,omitempty"` // Option to format spoken punctuated commands, must be enabled with punctuate parameter to true. Eg: "i went to the store comma new paragraph then i went home period" --> "i went to the store, <\n> then i went home."
	Encoding        string   `json:"encoding,omitempty" schema:"encoding,omitempty"`
	Endpointing     string   `json:"endpointing,omitempty" schema:"endpointing,omitempty"`
	Extra           []string `json:"extra,omitempty" schema:"extra,omitempty"`
	FillerWords     bool     `json:"filler_words,omitempty" schema:"filler_words,omitempty"`
	InterimResults  bool     `json:"interim_results,omitempty" schema:"interim_results,omitempty"`
	Keywords        []string `json:"keywords,omitempty" schema:"keywords,omitempty"`
	Keyterm         []string `json:"keyterm,omitempty" schema:"keyterm,omitempty"`
	Language        string   `json:"language,omitempty" schema:"language,omitempty"`
	Model           string   `json:"model,omitempty" schema:"model,omitempty"`
	Multichannel    bool     `json:"multichannel,omitempty" schema:"multichannel,omitempty"`
	NoDelay         bool     `json:"no_delay,omitempty" schema:"no_delay,omitempty"`
	Numerals        bool     `json:"numerals,omitempty" schema:"numerals,omitempty"`
	ProfanityFilter bool     `json:"profanity_filter,omitempty" schema:"profanity_filter,omitempty"`
	Punctuate       bool     `json:"punctuate,omitempty" schema:"punctuate,omitempty"`
	Redact          []string `json:"redact,omitempty" schema:"redact,omitempty"`
	Replace         []string `json:"replace,omitempty" schema:"replace,omitempty"`
	SampleRate      int      `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	Search          []string `json:"search,omitempty" schema:"search,omitempty"`
	SmartFormat     bool     `json:"smart_format,omitempty" schema:"smart_format,omitempty"`
	Tag             []string `json:"tag,omitempty" schema:"tag,omitempty"`
	Tier            string   `json:"tier,omitempty" schema:"tier,omitempty"` // Tier is only used for legacy models. Deprecated: This field is deprecated and will be removed in a future release.
	UtteranceEndMs  string   `json:"utterance_end_ms,omitempty" schema:"utterance_end_ms,omitempty"`
	VadEvents       bool     `json:"vad_events,omitempty" schema:"vad_events,omitempty"`
	Version         string   `json:"version,omitempty" schema:"version,omitempty"`
}
