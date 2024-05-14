// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
)

/***********************************/
// share/common structs
/***********************************/
type IntentsInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type SentimentInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type SummaryInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type TopicsInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type Metadata struct {
	RequestID     string         `json:"request_id,omitempty"`
	Created       string         `json:"created,omitempty"`
	Language      string         `json:"language,omitempty"`
	IntentsInfo   *IntentsInfo   `json:"intents_info,omitempty"`
	SentimentInfo *SentimentInfo `json:"sentiment_info,omitempty"`
	SummaryInfo   *SummaryInfo   `json:"summary_info,omitempty"`
	TopicsInfo    *TopicsInfo    `json:"topics_info,omitempty"`
}

type Average struct {
	Sentiment      string  `json:"sentiment,omitempty"`
	SentimentScore float64 `json:"sentiment_score,omitempty"`
}

type Summary struct {
	Text string `json:"text,omitempty"`
}

type Topic struct {
	Topic           string  `json:"topic,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
}

type Intent struct {
	Intent          string  `json:"intent,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
}

type Segment struct {
	Text           string    `json:"text,omitempty"`
	StartWord      int       `json:"start_word,omitempty"`
	EndWord        int       `json:"end_word,omitempty"`
	Sentiment      *string   `json:"sentiment,omitempty"`
	SentimentScore *float64  `json:"sentiment_score,omitempty"`
	Topics         *[]Topic  `json:"topics,omitempty"`
	Intents        *[]Intent `json:"intents,omitempty"`
}

type Sentiments struct {
	Segments []Segment `json:"segments,omitempty"`
	Average  Average   `json:"average,omitempty"`
}

type Topics struct {
	Segments []Segment `json:"segments,omitempty"`
}

type Intents struct {
	Segments []Segment `json:"segments,omitempty"`
}

type Results struct {
	Sentiments *Sentiments `json:"sentiments,omitempty"`
	Summary    *Summary    `json:"summary,omitempty"`
	Topics     *Topics     `json:"topics,omitempty"`
	Intents    *Intents    `json:"intents,omitempty"`
}

/***********************************/
// Request/Input structs
/***********************************/
type AnalyzeOptions interfaces.AnalyzeOptions

/***********************************/
// response/result structs
/***********************************/
type AnalyzeResponse struct {
	RequestID string   `json:"request_id,omitempty"` // for ?callback=...
	Metadata  Metadata `json:"metadata,omitempty"`
	Results   Results  `json:"results,omitempty"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
