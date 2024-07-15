// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

/*
SpeakOptions contain all of the knobs and dials to transform text into speech
using the Deepgram API

Please see the text-to-speech documentation for more details:
https://developers.deepgram.com/reference/text-to-speech-api
*/
type SpeakOptions struct {
	Model          string `json:"model,omitempty" schema:"model,omitempty"`
	Encoding       string `json:"encoding,omitempty" schema:"encoding,omitempty"`
	Container      string `json:"container,omitempty" schema:"container,omitempty"`
	SampleRate     int    `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	BitRate        int    `json:"bit_rate,omitempty" schema:"bit_rate,omitempty"`
	Callback       string `json:"callback,omitempty" schema:"callback,omitempty"`
	CallbackMethod string `json:"callback_method,omitempty" schema:"callback_method,omitempty"`
}

/*
WSSpeakOptions contain all of the knobs and dials to transform text into speech
using the Deepgram API

Please see the text-to-speech documentation for more details:
https://developers.deepgram.com/reference/transform-text-to-speech-websocket
*/
type WSSpeakOptions struct {
	Model      string `json:"model,omitempty" schema:"model,omitempty"`
	Encoding   string `json:"encoding,omitempty" schema:"encoding,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	BitRate    int    `json:"bit_rate,omitempty" schema:"bit_rate,omitempty"`
}
