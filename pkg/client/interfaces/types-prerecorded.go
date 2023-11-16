// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package contains the interface to manage the prerecorded interfaces for the Deepgram API
*/
package interfaces

/*
PreRecordedTranscriptionOptions contain all of the knobs and dials to control a Prerecorded transcription
from the Deepgram API

Please see the documentation for live/streaming for more details:
https://developers.deepgram.com/reference/pre-recorded
*/
type PreRecordedTranscriptionOptions struct {
	Alternatives       int         `json:"alternatives" url:"alternatives,omitempty" `
	Callback           string      `json:"callback" url:"callback,omitempty" `
	DetectEntities     bool        `json:"detect_entities" url:"detect_entities,omitempty"`
	DetectLanguage     bool        `json:"detect_language" url:"detect_language,omitempty" `
	DetectTopics       bool        `json:"detect_topics" url:"detect_topics,omitempty" `
	Diarize            bool        `json:"diarize" url:"diarize,omitempty" `
	Diarize_version    string      `json:"diarize_version" url:"diarize_version,omitempty" `
	Keywords           []string    `json:"keywords" url:"keywords,omitempty" `
	Language           string      `json:"language" url:"language,omitempty" `
	Model              string      `json:"model" url:"model,omitempty" `
	Multichannel       bool        `json:"multichannel" url:"multichannel,omitempty" `
	Numerals           bool        `json:"numerals" url:"numerals,omitempty" ` // Same as Numbers, old name for same option
	Paragraphs         bool        `json:"paragraphs" url:"paragraphs,omitempty" `
	Profanity_filter   bool        `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate          bool        `json:"punctuate" url:"punctuate,omitempty" `
	Redact             []string    `json:"redact" url:"redact,omitempty" `
	Replace            []string    `json:"replace" url:"replace,omitempty" `
	Search             []string    `json:"search" url:"search,omitempty" `
	Sentiment          bool        `json:"sentiment" url:"sentiment,omitempty" `
	SentimentThreshold float64     `json:"sentiment_threshold" url:"sentiment_threshold,omitempty" `
	SmartFormat        bool        `json:"smart_format" url:"smart_format,omitempty" `
	Summarize          interface{} `json:"summarize" url:"summarize,omitempty" ` // bool | string
	Tag                []string    `json:"tag" url:"tag,omitempty"`
	Tier               string      `json:"tier" url:"tier,omitempty" `
	Utterances         bool        `json:"utterances" url:"utterances,omitempty" `
	Utt_split          float64     `json:"utt_split" url:"utt_split,omitempty" `
	Version            string      `json:"version" url:"version,omitempty" `
	FillerWords        string      `json:"filler_words" url:"filler_words,omitempty" `
}
