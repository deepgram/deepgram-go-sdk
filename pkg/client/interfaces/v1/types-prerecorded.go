// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

/*
PreRecordedTranscriptionOptions contain all of the knobs and dials to control a Prerecorded transcription
from the Deepgram API

Please see the prerecorded audio documentation for more details:
https://developers.deepgram.com/reference/pre-recorded
*/
type PreRecordedTranscriptionOptions struct {
	Alternatives     int      `json:"alternatives,omitempty" schema:"alternatives,omitempty"`
	Callback         string   `json:"callback,omitempty" schema:"callback,omitempty"`
	CallbackMethod   string   `json:"callback_method,omitempty" schema:"callback_method,omitempty"`
	Channels         int      `json:"channels,omitempty" schema:"channels,omitempty"`
	CustomIntent     []string `json:"custom_intent,omitempty" schema:"custom_intent,omitempty"`
	CustomIntentMode string   `json:"custom_intent_mode,omitempty" schema:"custom_intent_mode,omitempty"`
	CustomTopic      []string `json:"custom_topic,omitempty" schema:"custom_topic,omitempty"`
	CustomTopicMode  string   `json:"custom_topic_mode,omitempty" schema:"custom_topic_mode,omitempty"`
	DetectEntities   bool     `json:"detect_entities,omitempty" schema:"detect_entities,omitempty"`
	DetectLanguage   bool     `json:"detect_language,omitempty" schema:"detect_language,omitempty"`
	DetectTopics     bool     `json:"detect_topics,omitempty" schema:"detect_topics,omitempty"`
	Diarize          bool     `json:"diarize,omitempty" schema:"diarize,omitempty"`
	DiarizeVersion   string   `json:"diarize_version,omitempty" schema:"diarize_version,omitempty"`
	Dictation        bool     `json:"dictation,omitempty" schema:"dictation,omitempty"`
	Encoding         string   `json:"encoding,omitempty" schema:"encoding,omitempty"`
	Extra            []string `json:"extra,omitempty" schema:"extra,omitempty"`
	FillerWords      bool     `json:"filler_words,omitempty" schema:"filler_words,omitempty"`
	Intents          bool     `json:"intents,omitempty" schema:"intents,omitempty"`
	Keywords         []string `json:"keywords,omitempty" schema:"keywords,omitempty"`
	Keyterm          []string `json:"keyterm,omitempty" schema:"keyterm,omitempty"`
	Language         string   `json:"language,omitempty" schema:"language,omitempty"`
	Measurements     bool     `json:"measurements,omitempty" schema:"measurements,omitempty"`
	Model            string   `json:"model,omitempty" schema:"model,omitempty"`
	Multichannel     bool     `json:"multichannel,omitempty" schema:"multichannel,omitempty"`
	Numerals         bool     `json:"numerals,omitempty" schema:"numerals,omitempty"`
	Paragraphs       bool     `json:"paragraphs,omitempty" schema:"paragraphs,omitempty"`
	ProfanityFilter  bool     `json:"profanity_filter,omitempty" schema:"profanity_filter,omitempty"`
	Punctuate        bool     `json:"punctuate,omitempty" schema:"punctuate,omitempty"`
	Redact           []string `json:"redact,omitempty" schema:"redact,omitempty"`
	Replace          []string `json:"replace,omitempty" schema:"replace,omitempty"`
	SampleRate       int      `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	Search           []string `json:"search,omitempty" schema:"search,omitempty"`
	Sentiment        bool     `json:"sentiment,omitempty" schema:"sentiment,omitempty"`
	SmartFormat      bool     `json:"smart_format,omitempty" schema:"smart_format,omitempty"`
	Summarize        string   `json:"summarize,omitempty" schema:"summarize,omitempty"`
	Tag              []string `json:"tag,omitempty" schema:"tag,omitempty"`
	Topics           bool     `json:"topics,omitempty" schema:"topics,omitempty"`
	UttSplit         float64  `json:"utt_split,omitempty" schema:"utt_split,omitempty"`
	Utterances       bool     `json:"utterances,omitempty" schema:"utterances,omitempty"`
	Version          string   `json:"version,omitempty" schema:"version,omitempty"`
}
