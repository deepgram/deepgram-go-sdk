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
	Alternatives       int      `json:"alternatives,omitempty" url:"alternatives,omitempty"`
	Callback           string   `json:"callback,omitempty" url:"callback,omitempty"`
	DetectEntities     bool     `json:"detect_entities,omitempty" url:"detect_entities,omitempty"`
	DetectLanguage     bool     `json:"detect_language,omitempty" url:"detect_language,omitempty"`
	DetectTopics       bool     `json:"detect_topics,omitempty" url:"detect_topics,omitempty"`
	Diarize            bool     `json:"diarize,omitempty" url:"diarize,omitempty"`
	DiarizeVersion     string   `json:"diarize_version,omitempty" url:"diarize_version,omitempty"`
	FillerWords        string   `json:"filler_words,omitempty" url:"filler_words,omitempty"`
	Keywords           []string `json:"keywords,omitempty" url:"keywords,omitempty"`
	Language           string   `json:"language,omitempty" url:"language,omitempty"`
	Model              string   `json:"model,omitempty" url:"model,omitempty"`
	Multichannel       bool     `json:"multichannel,omitempty" url:"multichannel,omitempty"`
	Numerals           bool     `json:"numerals,omitempty" url:"numerals,omitempty"`
	Paragraphs         bool     `json:"paragraphs,omitempty" url:"paragraphs,omitempty"`
	ProfanityFilter    bool     `json:"profanity_filter,omitempty" url:"profanity_filter,omitempty"`
	Punctuate          bool     `json:"punctuate,omitempty" url:"punctuate,omitempty"`
	Redact             []string `json:"redact,omitempty" url:"redact,omitempty"`
	Replace            []string `json:"replace,omitempty" url:"replace,omitempty"`
	Search             []string `json:"search,omitempty" url:"search,omitempty"`
	Sentiment          bool     `json:"sentiment,omitempty" url:"sentiment,omitempty"`
	SentimentThreshold float64  `json:"sentiment_threshold,omitempty" url:"sentiment_threshold,omitempty"`
	SmartFormat        bool     `json:"smart_format,omitempty" url:"smart_format,omitempty"`
	Summarize          string   `json:"summarize,omitempty" url:"summarize,omitempty"`
	Tag                []string `json:"tag,omitempty" url:"tag,omitempty"`
	Tier               string   `json:"tier,omitempty" url:"tier,omitempty"`
	UttSplit           float64  `json:"utt_split,omitempty" url:"utt_split,omitempty"`
	Utterances         bool     `json:"utterances,omitempty" url:"utterances,omitempty"`
	Version            string   `json:"version,omitempty" url:"version,omitempty"`
}
