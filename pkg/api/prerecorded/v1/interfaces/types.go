// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

/***********************************/
// share/common structs
/***********************************/
type Metadata struct {
	TransactionKey string   `json:"transaction_key,omitempty"`
	RequestID      string   `json:"request_id,omitempty"`
	Sha256         string   `json:"sha256,omitempty"`
	Created        string   `json:"created,omitempty"`
	Duration       float64  `json:"duration,omitempty"`
	Channels       int      `json:"channels,omitempty"`
	Models         []string `json:"models,omitempty"`
	SummaryInfo    struct {
		InputTokens  int    `json:"input_tokens,omitempty"`
		OutputTokens int    `json:"output_tokens,omitempty"`
		ModelUUID    string `json:"model_uuid,omitempty"`
	} `json:"summary_info,omitempty"`
	ModelInfo map[string]struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
		Arch    string `json:"arch,omitempty"`
	} `json:"model_info,omitempty"`
	Warnings []Warning `json:"warnings,omitempty"`
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
	Word              string  `json:"word,omitempty"`
	Start             float64 `json:"start,omitempty"`
	End               float64 `json:"end,omitempty"`
	Confidence        float64 `json:"confidence,omitempty"`
	Speaker           int     `json:"speaker,omitempty"`
	SpeakerConfidence float64 `json:"speaker_confidence,omitempty"`
	PunctuatedWord    string  `json:"punctuated_word,omitempty"`
}

type Translation struct {
	Language    string `json:"language,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type Alternative struct {
	Transcript  string      `json:"transcript,omitempty"`
	Confidence  float64     `json:"confidence,omitempty"`
	Words       []Word      `json:"words,omitempty"`
	Paragraphs  Paragraph   `json:"paragraphs,omitempty"`
	Entities    []Entity    `json:"entities,omitempty"`
	Summaries   []SummaryV1 `json:"summaries,omitempty"`
	Translation Translation `json:"translation,omitempty"`
	Topics      []Topics    `json:"topics,omitempty"`
}

type Paragraphs struct {
	Transcript string      `json:"transcript,omitempty"`
	Paragraphs []Paragraph `json:"paragraphs,omitempty"`
}

type Paragraph struct {
	Sentences []Sentence `json:"sentences,omitempty"`
	NumWords  int        `json:"num_words,omitempty"`
	Start     float64    `json:"start,omitempty"`
	End       float64    `json:"end,omitempty"`
}

type Sentence struct {
	Text  string  `json:"text,omitempty"`
	Start float64 `json:"start,omitempty"`
	End   float64 `json:"end,omitempty"`
}

type Entity struct {
	Label      string  `json:"label,omitempty"`
	Value      string  `json:"value,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	StartWord  float64 `json:"start_word,omitempty"`
	EndWord    float64 `json:"end_word,omitempty"`
}

type Topics struct {
	Text      string  `json:"text,omitempty"`
	StartWord int     `json:"start_word,omitempty"`
	EndWord   int     `json:"end_word,omitempty"`
	Topics    []Topic `json:"topics,omitempty"`
}

type Topic struct {
	Topic      string  `json:"topic,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type Channel struct {
	Search             []Search      `json:"search,omitempty"`
	Alternatives       []Alternative `json:"alternatives,omitempty"`
	DetectedLanguage   string        `json:"detected_language,omitempty"`
	LanguageConfidence float64       `json:"language_confidence,omitempty"`
}

type Utterance struct {
	Start      float64 `json:"start,omitempty"`
	End        float64 `json:"end,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Channel    int     `json:"channel,omitempty"`
	Transcript string  `json:"transcript,omitempty"`
	Words      []Word  `json:"words,omitempty"`
	Speaker    int     `json:"speaker,omitempty"`
	ID         string  `json:"id,omitempty"`
}

type Result struct {
	Channels   []Channel   `json:"channels,omitempty"`
	Utterances []Utterance `json:"utterances,omitempty"`
	Summary    SummaryV2   `json:"summary,omitempty"`
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

/***********************************/
// response/result structs
/***********************************/
type PreRecordedResponse struct {
	RequestID string   `json:"request_id,omitempty"` // for ?callback=...
	Metadata  Metadata `json:"metadata,omitempty"`
	Results   Result   `json:"results,omitempty"`
}
