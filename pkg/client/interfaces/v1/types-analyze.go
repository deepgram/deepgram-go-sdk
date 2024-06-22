// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

/*
AnalyzeOptions contain all of the knobs and dials to control a read transcription
from the Deepgram API

Please see the text intelligence documentation for more details:
https://developers.deepgram.com/reference/text-intelligence-apis
*/

type AnalyzeOptions struct {
	Callback         string   `json:"callback,omitempty" schema:"callback,omitempty"`
	CallbackMethod   string   `json:"callback_method,omitempty" schema:"callback_method,omitempty"`
	CustomIntent     []string `json:"custom_intent,omitempty" schema:"custom_intent,omitempty"`
	CustomIntentMode string   `json:"custom_intent_mode,omitempty" schema:"custom_intent_mode,omitempty"`
	CustomTopic      []string `json:"custom_topic,omitempty" schema:"custom_topic,omitempty"`
	CustomTopicMode  string   `json:"custom_topic_mode,omitempty" schema:"custom_topic_mode,omitempty"`
	Intents          bool     `json:"intents,omitempty" schema:"intents,omitempty"`
	Language         string   `json:"language,omitempty" schema:"language,omitempty"`
	Summarize        bool     `json:"summarize,omitempty" schema:"summarize,omitempty"`
	Sentiment        bool     `json:"sentiment,omitempty" schema:"sentiment,omitempty"`
	Topics           bool     `json:"topics,omitempty" schema:"topics,omitempty"`
}
