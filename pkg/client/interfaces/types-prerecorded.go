// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

type PreRecordedTranscriptionOptions struct {
	Alternatives       int         `json:"alternatives" url:"alternatives,omitempty" `
	AnalyzeSentiment   bool        `json:"analyze_sentiment" url:"analyze_sentiment,omitempty" `
	Callback           string      `json:"callback" url:"callback,omitempty" `
	Dates              bool        `json:"dates" url:"dates,omitempty"` // Indicates whether to convert dates from written format (e.g., january first) to numerical format (e.g., 01-01).
	DetectEntities     bool        `json:"detect_entities" url:"detect_entities,omitempty"`
	DetectLanguage     bool        `json:"detect_language" url:"detect_language,omitempty" `
	DetectTopics       bool        `json:"detect_topics" url:"detect_topics,omitempty" `
	Diarize            bool        `json:"diarize" url:"diarize,omitempty" `
	Diarize_version    string      `json:"diarize_version" url:"diarize_version,omitempty" `
	Dictation          bool        `json:"dictation" url:"dictation,omitempty"` // Option to format punctuated commands. Eg: "i went to the store period new paragraph then i went home" --> "i went to the store. <\n> then i went home"
	Keywords           []string    `json:"keywords" url:"keywords,omitempty" `
	KeywordBoost       string      `json:"keyword_boost" url:"keyword_boost,omitempty" `
	Language           string      `json:"language" url:"language,omitempty" `
	Measurements       bool        `json:"measurements" url:"measurements,omitempty"`
	Model              string      `json:"model" url:"model,omitempty" `
	Multichannel       bool        `json:"multichannel" url:"multichannel,omitempty" `
	Ner                bool        `json:"ner" url:"ner,omitempty" `
	Numbers            bool        `json:"numbers" url:"numbers,omitempty" `
	Numerals           bool        `json:"numerals" url:"numerals,omitempty" ` // Same as Numbers, old name for same option
	Paragraphs         bool        `json:"paragraphs" url:"paragraphs,omitempty" `
	Profanity_filter   bool        `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate          bool        `json:"punctuate" url:"punctuate,omitempty" `
	Redact             bool        `json:"redact" url:"redact,omitempty" `
	Replace            []string    `json:"replace" url:"replace,omitempty" `
	Search             []string    `json:"search" url:"search,omitempty" `
	Sentiment          bool        `json:"sentiment" url:"sentiment,omitempty" `
	SentimentThreshold float64     `json:"sentiment_threshold" url:"sentiment_threshold,omitempty" `
	SmartFormat        bool        `json:"smart_format" url:"smart_format,omitempty" `
	Summarize          interface{} `json:"summarize" url:"summarize,omitempty" ` // bool | string
	Tag                []string    `json:"tag" url:"tag,omitempty"`
	Tier               string      `json:"tier" url:"tier,omitempty" `
	Times              bool        `json:"times" url:"times,omitempty"` // Indicates whether to convert times from written format (e.g., 3:00 pm) to numerical format (e.g., 15:00).
	Translate          string      `json:"translate" url:"translate,omitempty" `
	Utterances         bool        `json:"utterances" url:"utterances,omitempty" `
	Utt_split          float64     `json:"utt_split" url:"utt_split,omitempty" `
	Version            string      `json:"version" url:"version,omitempty" `
	FillerWords        string      `json:"filler_words" url:"filler_words,omitempty" `
}
