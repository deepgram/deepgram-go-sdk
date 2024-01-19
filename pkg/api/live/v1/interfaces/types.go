// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

/***********************************/
// shared/common structs
/***********************************/
// Word is a single word in a transcript
type Word struct {
	Confidence     float64 `json:"confidence,omitempty"`
	End            float64 `json:"end,omitempty"`
	PunctuatedWord string  `json:"punctuated_word,omitempty"`
	Start          float64 `json:"start,omitempty"`
	Word           string  `json:"word,omitempty"`
	Speaker        int     `json:"speaker,omitempty"`
}

// Alternative is a single alternative in a transcript
type Alternative struct {
	Confidence float64 `json:"confidence,omitempty"`
	Transcript string  `json:"transcript,omitempty"`
	Words      []Word  `json:"words,omitempty"`
}

// Channel is a single channel in a transcript
type Channel struct {
	Alternatives []Alternative `json:"alternatives,omitempty"`
}

// ModelInfo is the model information for a transcript
type ModelInfo struct {
	Arch    string `json:"arch,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Metadata is the metadata for a transcript
type Metadata struct {
	ModelInfo ModelInfo `json:"model_info,omitempty"`
	ModelUUID string    `json:"model_uuid,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
}

/***********************************/
// Results from Live Transcription
/***********************************/
// MessageResponse is the response from a live transcription
type MessageResponse struct {
	Channel      Channel  `json:"channel,omitempty"`
	ChannelIndex []int    `json:"channel_index,omitempty"`
	Duration     float64  `json:"duration,omitempty"`
	IsFinal      bool     `json:"is_final,omitempty"`
	Metadata     Metadata `json:"metadata,omitempty"`
	SpeechFinal  bool     `json:"speech_final,omitempty"`
	Start        float64  `json:"start,omitempty"`
	Type         string   `json:"type,omitempty"`
}

// MetadataResponse is the response from a live transcription
type MetadataResponse struct {
	Channels       int                  `json:"channels,omitempty"`
	Created        string               `json:"created,omitempty"`
	Duration       float64              `json:"duration,omitempty"`
	ModelInfo      map[string]ModelInfo `json:"model_info,omitempty"`
	Models         []string             `json:"models,omitempty"`
	RequestID      string               `json:"request_id,omitempty"`
	Sha256         string               `json:"sha256,omitempty"`
	TransactionKey string               `json:"transaction_key,omitempty"`
	Type           string               `json:"type,omitempty"`
}

// UtteranceEndResponse is the response from a live transcription
type UtteranceEndResponse struct {
	Type        string  `json:"type,omitempty"`
	Channel     []int   `json:"channel,omitempty"`
	LastWordEnd float64 `json:"last_word_end,omitempty"`
}

// ErrorResponse is the response from a live transcription
type ErrorResponse struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Variant     string `json:"variant"`
}
