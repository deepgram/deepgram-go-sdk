// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

// share/common structs
type Metadata struct {
	TransactionKey string   `json:"transaction_key"`
	RequestId      string   `json:"request_id"`
	Sha256         string   `json:"sha256"`
	Created        string   `json:"created"`
	Duration       float64  `json:"duration"`
	Channels       int      `json:"channels"`
	Models         []string `json:"models"`
	ModelInfo      map[string]struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Arch    string `json:"arch"`
	} `json:"model_info"`
	Warnings []*Warning `json:"warnings,omitempty"`
}

type Warning struct {
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

type Hit struct {
	Confidence float64 `json:"confidence"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Snippet    string  `json:"snippet"`
}

type Search struct {
	Query string `json:"query"`
	Hits  []Hit  `json:"hits"`
}

type WordBase struct {
	Word              string  `json:"word"`
	Start             float64 `json:"start"`
	End               float64 `json:"end"`
	Confidence        float64 `json:"confidence"`
	Speaker           *int    `json:"speaker,omitempty"`
	SpeakerConfidence float64 `json:"speaker_confidence,omitempty"`
	Punctuated_Word   string  `json:"punctuated_word,omitempty"`
	Sentiment         string  `json:"sentiment,omitempty"`
}

type Alternative struct {
	Transcript string          `json:"transcript"`
	Confidence float64         `json:"confidence"`
	Words      []WordBase      `json:"words"`
	Summaries  []*SummaryV1    `json:"summaries,omitempty"`
	Paragraphs *ParagraphGroup `json:"paragraphs,omitempty"`
	Topics     []*TopicBase    `json:"topics,omitempty"`
	Entities   []*EntityBase   `json:"entities,omitempty"`
}

type ParagraphGroup struct {
	Transcript string          `json:"transcript"`
	Paragraphs []ParagraphBase `json:"paragraphs"`
}

type ParagraphBase struct {
	Sentences []SentenceBase `json:"sentences"`
	NumWords  int            `json:"num_words"`
	Start     float64        `json:"start"`
	End       float64        `json:"end"`
}

type SentenceBase struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type EntityBase struct {
	Label      string  `json:"label"`
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	StartWord  int     `json:"start_word"`
	EndWord    int     `json:"end_word"`
}

type TopicBase struct {
	Text      string  `json:"text"`
	StartWord int     `json:"start_word"`
	EndWord   int     `json:"end_word"`
	Topics    []Topic `json:"topics"`
}

type Topic struct {
	Topic      string  `json:"topic"`
	Confidence float64 `json:"confidence"`
}

type Channel struct {
	Search           []*Search     `json:"search,omitempty"`
	Alternatives     []Alternative `json:"alternatives"`
	DetectedLanguage string        `json:"detected_language,omitempty"`
}

type Utterance struct {
	Start      float64    `json:"start"`
	End        float64    `json:"end"`
	Confidence float64    `json:"confidence"`
	Channel    int        `json:"channel"`
	Transcript string     `json:"transcript"`
	Words      []WordBase `json:"words"`
	Speaker    *int       `json:"speaker,omitempty"`
	Id         string     `json:"id"`
}

type Results struct {
	Utterances []*Utterance `json:"utterances,omitempty"`
	Channels   []Channel    `json:"channels"`
	Summary    *SummaryV2   `json:"summary,omitempty"`
}

type SummaryV1 struct {
	Summary   string `json:"summary"`
	StartWord int    `json:"start_word"`
	EndWord   int    `json:"end_word"`
}

type SummaryV2 struct {
	Short  string `json:"short"`
	Result string `json:"result"`
}

// Response
type PreRecordedResponse struct {
	Request_id string   `json:"request_id,omitempty"`
	Metadata   Metadata `json:"metadata"`
	Results    Results  `json:"results"`
}
