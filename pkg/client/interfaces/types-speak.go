// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

/*
SpeakOptions contain all of the knobs and dials to transform text into speech
using the Deepgram API

Please see the Text-to-Speech documentation for more details:
https://developers.deepgram.com/reference/text-to-speech-preview-api
*/
type SpeakOptions struct {
	Model      string `json:"model,omitempty" url:"model,omitempty"`
	Encoding   string `json:"encoding,omitempty" url:"encoding,omitempty"`
	Container  string `json:"container,omitempty" url:"container,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty" url:"sample_rate,omitempty"`
	BitRate    int    `json:"bit_rate,omitempty" url:"bit_rate,omitempty"`
}
