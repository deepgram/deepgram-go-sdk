// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

/*
Shared defintions for the Deepgram API
*/
type Words struct {
	Confidence     float64 `json:"confidence,omitempty"`
	End            float64 `json:"end,omitempty"`
	PunctuatedWord string  `json:"punctuated_word,omitempty"`
	Start          float64 `json:"start,omitempty"`
	Word           string  `json:"word,omitempty"`
}
type Alternatives struct {
	Confidence float64 `json:"confidence,omitempty"`
	Transcript string  `json:"transcript,omitempty"`
	Words      []Words `json:"words,omitempty"`
}
type Channel struct {
	Alternatives []Alternatives `json:"alternatives,omitempty"`
}

type ModelInfo struct {
	Arch    string `json:"arch,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}
type Metadata struct {
	ModelInfo ModelInfo `json:"model_info,omitempty"`
	ModelUUID string    `json:"model_uuid,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
}

/*
Results from Live Transcription
*/
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

type ErrorResponse struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Variant     string `json:"variant"`
}
