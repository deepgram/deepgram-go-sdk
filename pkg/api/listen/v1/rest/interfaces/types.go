// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
)

/***********************************/
// share/common structs
/***********************************/
type SummaryInfo struct {
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
	ModelUUID    string `json:"model_uuid,omitempty"`
}

type ModelInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Arch    string `json:"arch,omitempty"`
}

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

type TopicsInfo struct {
	ModelUUID    string `json:"model_uuid,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type Metadata struct {
	TransactionKey string               `json:"transaction_key,omitempty"`
	RequestID      string               `json:"request_id,omitempty"`
	Sha256         string               `json:"sha256,omitempty"`
	Created        string               `json:"created,omitempty"`
	Duration       float64              `json:"duration,omitempty"`
	Channels       int                  `json:"channels,omitempty"`
	Models         []string             `json:"models,omitempty"`
	ModelInfo      map[string]ModelInfo `json:"model_info,omitempty"`
	Warnings       *[]Warning           `json:"warnings,omitempty"`
	SummaryInfo    *SummaryInfo         `json:"summary_info,omitempty"`
	IntentsInfo    *IntentsInfo         `json:"intents_info,omitempty"`
	SentimentInfo  *SentimentInfo       `json:"sentiment_info,omitempty"`
	TopicsInfo     *TopicsInfo          `json:"topics_info,omitempty"`
	Extra          map[string]string    `json:"extra,omitempty"`
}

type Warning struct {
	Parameter string `json:"parameter,omitempty"`
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
}

type Hit struct {
	Confidence float64 `json:"confidence,omitempty"`
	Start      float64 `json:"start,omitempty"`
	End        float64 `json:"end,omitempty"`
	Snippet    string  `json:"snippet,omitempty"`
}

type Search struct {
	Query string `json:"query,omitempty"`
	Hits  []Hit  `json:"hits,omitempty"`
}

type Word struct {
	Word              string   `json:"word,omitempty"`
	Start             float64  `json:"start,omitempty"`
	End               float64  `json:"end,omitempty"`
	Confidence        float64  `json:"confidence,omitempty"`
	Speaker           *int     `json:"speaker,omitempty"`
	SpeakerConfidence *float64 `json:"speaker_confidence,omitempty"`
	PunctuatedWord    string   `json:"punctuated_word,omitempty"`
	Sentiment         *string  `json:"sentiment,omitempty"`
	SentimentScore    *float64 `json:"sentiment_score,omitempty"`
	Language          string   `json:"language,omitempty"`
}

type Translation struct {
	Language    string `json:"language,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type Alternative struct {
	Transcript  string       `json:"transcript,omitempty"`
	Confidence  float64      `json:"confidence,omitempty"`
	Words       []Word       `json:"words,omitempty"`
	Paragraphs  *Paragraphs  `json:"paragraphs,omitempty"`
	Entities    *[]Entity    `json:"entities,omitempty"`
	Summaries   *[]SummaryV1 `json:"summaries,omitempty"`
	Translation *Translation `json:"translation,omitempty"`
	Languages   []string     `json:"languages,omitempty"`
}

type Paragraphs struct {
	Transcript string      `json:"transcript,omitempty"`
	Paragraphs []Paragraph `json:"paragraphs,omitempty"`
}

type Paragraph struct {
	Sentences      []Sentence `json:"sentences,omitempty"`
	NumWords       int        `json:"num_words,omitempty"`
	Start          float64    `json:"start,omitempty"`
	End            float64    `json:"end,omitempty"`
	Speaker        *int       `json:"speaker,omitempty"`
	Sentiment      *string    `json:"sentiment,omitempty"`
	SentimentScore *float64   `json:"sentiment_score,omitempty"`
}

type Sentence struct {
	Text           string   `json:"text,omitempty"`
	Start          float64  `json:"start,omitempty"`
	End            float64  `json:"end,omitempty"`
	Sentiment      *string  `json:"sentiment,omitempty"`
	SentimentScore *float64 `json:"sentiment_score,omitempty"`
}

type Entity struct {
	Label      string  `json:"label,omitempty"`
	Value      string  `json:"value,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	StartWord  float64 `json:"start_word,omitempty"`
	EndWord    float64 `json:"end_word,omitempty"`
}

type Channel struct {
	Search             *[]Search     `json:"search,omitempty"`
	Alternatives       []Alternative `json:"alternatives,omitempty"`
	DetectedLanguage   string        `json:"detected_language,omitempty"`
	LanguageConfidence float64       `json:"language_confidence,omitempty"`
}

type Utterance struct {
	Start          float64  `json:"start,omitempty"`
	End            float64  `json:"end,omitempty"`
	Confidence     float64  `json:"confidence,omitempty"`
	Channel        int      `json:"channel,omitempty"`
	Transcript     string   `json:"transcript,omitempty"`
	Words          []Word   `json:"words,omitempty"`
	Speaker        *int     `json:"speaker,omitempty"`
	Sentiment      *string  `json:"sentiment,omitempty"`
	SentimentScore *float64 `json:"sentiment_score,omitempty"`
	ID             string   `json:"id,omitempty"`
}

type Intent struct {
	Intent          string  `json:"intent,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
}

type Average struct {
	Sentiment      string  `json:"sentiment,omitempty"`
	SentimentScore float64 `json:"sentiment_score,omitempty"`
}

type Topic struct {
	Topic           string  `json:"topic,omitempty"`
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

type SummaryV1 struct {
	Summary   string `json:"summary,omitempty"`
	StartWord int    `json:"start_word,omitempty"`
	EndWord   int    `json:"end_word,omitempty"`
}
type Summaries SummaryV1 // internal reference to old name

type SummaryV2 struct {
	Short  string `json:"short,omitempty"`
	Result string `json:"result,omitempty"`
}
type Summary SummaryV2 // internal reference to old name

type Result struct {
	Channels   []Channel   `json:"channels,omitempty"`
	Utterances []Utterance `json:"utterances,omitempty"`
	Summary    *SummaryV2  `json:"summary,omitempty"`
	Sentiments *Sentiments `json:"sentiments,omitempty"`
	Topics     *Topics     `json:"topics,omitempty"`
	Intents    *Intents    `json:"intents,omitempty"`
}

/***********************************/
// Request/Input structs
/***********************************/
type PreRecordedTranscriptionOptions interfaces.PreRecordedTranscriptionOptions

/***********************************/
// response/result structs
/***********************************/
// PreRecordedResponse is the PreRecorded Transcription
type PreRecordedResponse struct {
	RequestID string    `json:"request_id,omitempty"` // for ?callback=...
	Metadata  *Metadata `json:"metadata,omitempty"`
	Results   *Result   `json:"results,omitempty"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
