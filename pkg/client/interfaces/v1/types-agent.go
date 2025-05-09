// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

/*
SettingsOptions contain all of the knobs and dials to control the Agent API

Please see the live/streaming documentation for more details:
XXXX
*/
type SettingsOptions struct {
	Type         string `json:"type"`
	Experimental bool   `json:"experimental,omitempty"`
	Audio        Audio  `json:"audio"`
	Agent        Agent  `json:"agent"`
}

/*
Sub-structs in SettingsOptions
*/
type Input struct {
	Encoding   string `json:"encoding,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty"`
}
type Output struct {
	Encoding   string `json:"encoding,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"`
	Container  string `json:"container,omitempty"`
}
type Audio struct {
	Input  *Input  `json:"input,omitempty"`
	Output *Output `json:"output,omitempty"`
}
type Endpoint struct {
	Url     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Method  string            `json:"method,omitempty"`
}
type CartesiaVoice struct {
	Mode string `json:"mode,omitempty"`
	Id   string `json:"id,omitempty"`
}
type ListenProvider struct {
	Type     string   `json:"type"`
	Model    string   `json:"model"`
	Keyterms []string `json:"keyterms,omitempty"`
}
type SpeakProvider struct {
	Type         string         `json:"type"`
	Model        string         `json:"model,omitempty"`
	ModelId      string         `json:"model_id,omitempty"`
	Voice        *CartesiaVoice `json:"voice,omitempty"`
	Language     string         `json:"language,omitempty"`
	LanguageCode string         `json:"language_code,omitempty"`
}
type ThinkProvider struct {
	Type        string  `json:"type"`
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature,omitempty"`
}
type Item struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}
type Properties struct {
	Item Item `json:"item,omitempty"`
}
type Parameters struct {
	Type       string     `json:"type,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Required   []string   `json:"required,omitempty"`
}
type Headers struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
type Functions struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Parameters  Parameters `json:"parameters,omitempty"`
	Endpoint    Endpoint   `json:"endpoint,omitempty"`
}
type Listen struct {
	Provider ListenProvider `json:"provider"`
}
type Think struct {
	Provider  ThinkProvider `json:"provider"`
	Endpoint  *Endpoint     `json:"endpoint,omitempty"`
	Functions *[]Functions  `json:"functions,omitempty"`
	Prompt    string        `json:"prompt,omitempty"`
}
type Speak struct {
	Provider SpeakProvider `json:"provider,omitempty"`
	Endpoint *Endpoint     `json:"endpoint,omitempty"`
}
type Agent struct {
	Language string `json:"language,omitempty"`
	Listen   Listen `json:"listen"`
	Think    Think  `json:"think"`
	Speak    Speak  `json:"speak"`
	Greeting string `json:"greeting,omitempty"`
}
