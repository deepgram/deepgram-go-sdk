// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

/*
AnalyzeOptions contain all of the knobs and dials to control a read transcription
from the Deepgram API

Please see the documentation for live/streaming for more details:
https://developers.deepgram.com/reference/text-intelligence-apis
*/

type AnalyzeOptions struct {
	Callback         string   `json:"callback,omitempty" url:"callback,omitempty"`
	CallbackMethod   string   `json:"callback_method,omitempty" url:"callback_method,omitempty"`
	CustomIntent     []string `json:"custom_intent,omitempty" url:"custom_intent,omitempty"`
	CustomIntentMode string   `json:"custom_intent_mode,omitempty" url:"custom_intent_mode,omitempty"`
	CustomTopic      []string `json:"custom_topic,omitempty" url:"custom_topic,omitempty"`
	CustomTopicMode  string   `json:"custom_topic_mode,omitempty" url:"custom_topic_mode,omitempty"`
	Intents          bool     `json:"intents,omitempty" url:"intents,omitempty"`
	Language         string   `json:"language,omitempty" url:"language,omitempty"`
	Summarize        bool     `json:"summarize,omitempty" url:"summarize,omitempty"`
	Sentiment        bool     `json:"sentiment,omitempty" url:"sentiment,omitempty"`
	Topics           bool     `json:"topics,omitempty" url:"topics,omitempty"`
}
