// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
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

// Metadata describes the audio, model, and request that produced a transcription.
type Metadata struct {
	// TransactionKey is a deprecated transaction identifier returned by older API versions.
	TransactionKey string `json:"transaction_key,omitempty"`
	// RequestID identifies this Deepgram request for troubleshooting and support.
	RequestID string `json:"request_id,omitempty"`
	// Sha256 is the SHA-256 hash of the submitted audio.
	Sha256 string `json:"sha256,omitempty"`
	// Created is the time at which Deepgram created the transcription request.
	Created string `json:"created,omitempty"`
	// Duration is the submitted audio duration in seconds.
	Duration float64 `json:"duration,omitempty"`
	// Channels is the number of audio channels in the submitted media.
	Channels int `json:"channels,omitempty"`
	// Models lists the model identifiers used for the request.
	Models []string `json:"models,omitempty"`
	// ModelInfo maps each model identifier in Models to its name, version, and architecture.
	ModelInfo map[string]ModelInfo `json:"model_info,omitempty"`
	// Warnings contains non-fatal warnings returned while processing the request.
	Warnings *[]Warning `json:"warnings,omitempty"`
	// SummaryInfo describes the model usage for a requested summary.
	SummaryInfo *SummaryInfo `json:"summary_info,omitempty"`
	// IntentsInfo describes the model usage for requested intent detection.
	IntentsInfo *IntentsInfo `json:"intents_info,omitempty"`
	// SentimentInfo describes the model usage for requested sentiment analysis.
	SentimentInfo *SentimentInfo `json:"sentiment_info,omitempty"`
	// TopicsInfo describes the model usage for requested topic detection.
	TopicsInfo *TopicsInfo `json:"topics_info,omitempty"`
	// Extra contains caller-supplied metadata returned with the response.
	Extra map[string]string `json:"extra,omitempty"`
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

// Word is one recognized word and its timing and confidence information.
type Word struct {
	// Word is the recognized word without punctuation formatting.
	Word string `json:"word,omitempty"`
	// Start is the word start time in seconds from the beginning of the audio.
	Start float64 `json:"start,omitempty"`
	// End is the word end time in seconds from the beginning of the audio.
	End float64 `json:"end,omitempty"`
	// Confidence is Deepgram's confidence in the recognized word.
	Confidence float64 `json:"confidence,omitempty"`
	// Speaker is the diarized speaker number when diarization is enabled.
	Speaker *int `json:"speaker,omitempty"`
	// SpeakerConfidence is Deepgram's confidence in the diarized speaker assignment.
	SpeakerConfidence *float64 `json:"speaker_confidence,omitempty"`
	// PunctuatedWord is the word after punctuation and formatting are applied.
	PunctuatedWord string `json:"punctuated_word,omitempty"`
	// Sentiment is the word-level sentiment when sentiment analysis is requested.
	Sentiment *string `json:"sentiment,omitempty"`
	// SentimentScore is the word-level sentiment score when sentiment analysis is requested.
	SentimentScore *float64 `json:"sentiment_score,omitempty"`
	// Language is the detected BCP-47 language tag when multilingual detection is used.
	Language string `json:"language,omitempty"`
}

type Translation struct {
	Language    string `json:"language,omitempty"`
	Translation string `json:"translation,omitempty"`
}

// Alternative is one possible transcription for an audio channel.
type Alternative struct {
	// Transcript is the complete recognized text for this alternative.
	Transcript string `json:"transcript,omitempty"`
	// Confidence is Deepgram's confidence in this alternative.
	Confidence float64 `json:"confidence,omitempty"`
	// Words contains the timed words that make up the transcript.
	Words []Word `json:"words,omitempty"`
	// Paragraphs contains paragraph segmentation when paragraphs are requested.
	Paragraphs *Paragraphs `json:"paragraphs,omitempty"`
	// Entities contains detected entities when entity detection is requested.
	Entities *[]Entity `json:"entities,omitempty"`
	// Summaries contains summary segments returned for this alternative.
	Summaries *[]SummaryV1 `json:"summaries,omitempty"`
	// Translation contains the translated text when translation is requested.
	Translation *Translation `json:"translation,omitempty"`
	// Languages lists detected BCP-47 language tags when multilingual detection is used.
	Languages []string `json:"languages,omitempty"`
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

// Channel contains transcription results for one submitted audio channel.
type Channel struct {
	// Search contains search hits when search terms are requested.
	Search *[]Search `json:"search,omitempty"`
	// Alternatives contains the candidate transcriptions for the channel.
	Alternatives []Alternative `json:"alternatives,omitempty"`
	// DetectedLanguage is the dominant detected language for the channel.
	DetectedLanguage string `json:"detected_language,omitempty"`
	// LanguageConfidence is Deepgram's confidence in DetectedLanguage.
	LanguageConfidence float64 `json:"language_confidence,omitempty"`
}

// Utterance is a semantically segmented portion of the transcript.
// It is returned when the utterances option is enabled.
type Utterance struct {
	// Start is the utterance start time in seconds from the beginning of the audio.
	Start float64 `json:"start,omitempty"`
	// End is the utterance end time in seconds from the beginning of the audio.
	End float64 `json:"end,omitempty"`
	// Confidence is Deepgram's confidence in the utterance transcript.
	Confidence float64 `json:"confidence,omitempty"`
	// Channel identifies the submitted audio channel that produced this utterance.
	Channel int `json:"channel,omitempty"`
	// Transcript is the recognized text for the utterance.
	Transcript string `json:"transcript,omitempty"`
	// Words contains the timed words in the utterance.
	Words []Word `json:"words,omitempty"`
	// Speaker is the diarized speaker number when diarization is enabled.
	Speaker *int `json:"speaker,omitempty"`
	// Sentiment is the utterance-level sentiment when sentiment analysis is requested.
	Sentiment *string `json:"sentiment,omitempty"`
	// SentimentScore is the utterance-level sentiment score when sentiment analysis is requested.
	SentimentScore *float64 `json:"sentiment_score,omitempty"`
	// ID is the API identifier for the utterance.
	ID string `json:"id,omitempty"`
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

// Result contains the transcription and optional analysis output.
type Result struct {
	// Channels contains one result for each transcribed audio channel.
	Channels []Channel `json:"channels,omitempty"`
	// Utterances contains semantic transcript segments when utterances are requested.
	Utterances []Utterance `json:"utterances,omitempty"`
	// Summary contains the requested summary output.
	Summary *SummaryV2 `json:"summary,omitempty"`
	// Sentiments contains the requested sentiment-analysis output.
	Sentiments *Sentiments `json:"sentiments,omitempty"`
	// Topics contains the requested topic-detection output.
	Topics *Topics `json:"topics,omitempty"`
	// Intents contains the requested intent-detection output.
	Intents *Intents `json:"intents,omitempty"`
}

/***********************************/
// Request/Input structs
/***********************************/
type PreRecordedTranscriptionOptions interfaces.PreRecordedTranscriptionOptions

/***********************************/
// response/result structs
/***********************************/
// PreRecordedResponse is the successful Listen v1 batch transcription response.
type PreRecordedResponse struct {
	// RequestID identifies an asynchronous callback request when no full response is returned.
	RequestID string `json:"request_id,omitempty"`
	// Metadata describes the request, submitted audio, and models used for the transcription.
	Metadata *Metadata `json:"metadata,omitempty"`
	// Results contains the transcription and optional analysis output.
	Results *Result `json:"results,omitempty"`
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
