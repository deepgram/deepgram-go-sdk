// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the interface to manage the live/streaming interfaces for the Deepgram API
*/
package interfaces

/*
LiveTranscriptionOptions contain all of the knobs and dials to control the live transcription
from the Deepgram API

Please see the documentation for live/streaming for more details:
https://developers.deepgram.com/reference/streaming
*/
type LiveTranscriptionOptions struct {
	Alternatives     int      `json:"alternatives" url:"alternatives,omitempty" `
	Callback         string   `json:"callback" url:"callback,omitempty" `
	Channels         int      `json:"channels" url:"channels,omitempty" `
	Dates            bool     `json:"dates" url:"dates,omitempty"` // Indicates whether to convert dates from written format (e.g., january first) to numerical format (e.g., 01-01).
	Diarize          bool     `json:"diarize" url:"diarize,omitempty" `
	Diarize_version  string   `json:"diarize_version" url:"diarize_version,omitempty" `
	Dictation        bool     `json:"dictation" url:"dictation,omitempty"` // Option to format punctuated commands. Eg: "i went to the store period new paragraph then i went home" --> "i went to the store. <\n> then i went home"
	Encoding         string   `json:"encoding" url:"encoding,omitempty" `
	Endpointing      string   `json:"endpointing" url:"endpointing,omitempty" ` // Can be "false" to disable endpointing, or can be the milliseconds of silence to wait before returning a transcript. Default is 10 milliseconds. Is string here so it can accept "false" as a value.
	Interim_results  bool     `json:"interim_results" url:"interim_results,omitempty" `
	Keywords         []string `json:"keywords" url:"keywords,omitempty" `
	KeywordBoost     string   `json:"keyword_boost" url:"keyword_boost,omitempty" `
	Language         string   `json:"language" url:"language,omitempty" `
	Measurements     bool     `json:"measurements" url:"measurements,omitempty" `
	Model            string   `json:"model" url:"model,omitempty" `
	Multichannel     bool     `json:"multichannel" url:"multichannel,omitempty" `
	Ner              bool     `json:"ner" url:"ner,omitempty" `
	Numbers          bool     `json:"numbers" url:"numbers,omitempty" `
	Numerals         bool     `json:"numerals" url:"numerals,omitempty" `
	Profanity_filter bool     `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate        bool     `json:"punctuate" url:"punctuate,omitempty" `
	Redact           []string `json:"redact" url:"redact,omitempty" `
	Replace          string   `json:"replace" url:"replace,omitempty" `
	Sample_rate      int      `json:"sample_rate" url:"sample_rate,omitempty" `
	Search           []string `json:"search" url:"search,omitempty" `
	Smart_format     bool     `json:"smart_format" url:"smart_format,omitempty" `
	Tag              []string `json:"tag" url:"tag,omitempty" `
	Tier             string   `json:"tier" url:"tier,omitempty" `
	Times            bool     `json:"times" url:"times,omitempty" `
	Vad_turnoff      int      `json:"vad_turnoff" url:"vad_turnoff,omitempty" `
	Version          string   `json:"version" url:"version,omitempty" `
	FillerWords      string   `json:"filler_words" url:"filler_words,omitempty" `
}
