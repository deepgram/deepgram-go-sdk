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
	Type         string   `json:"type"`
	Tags         []string `json:"tags,omitempty"`
	Experimental bool     `json:"experimental,omitempty"`
	MipOptOut    bool     `json:"mip_opt_out,omitempty"`
	Audio        Audio    `json:"audio"`
	Agent        Agent    `json:"agent"`
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
	Provider map[string]interface{} `json:"provider,omitempty"`
}
type Think struct {
	Provider      map[string]interface{} `json:"provider,omitempty"`
	Endpoint      *Endpoint              `json:"endpoint,omitempty"`
	Functions     *[]Functions           `json:"functions,omitempty"`
	Prompt        string                 `json:"prompt,omitempty"`
	ContextLength any                    `json:"context_length,omitempty"` // int or "max"
}
type Speak struct {
	Provider map[string]interface{} `json:"provider,omitempty"`
	Endpoint *Endpoint              `json:"endpoint,omitempty"`
}
type Agent struct {
	Language      string   `json:"language,omitempty"`
	Listen        Listen   `json:"listen,omitempty"`
	Think         Think    `json:"think,omitempty"`
	Speak         Speak    `json:"speak,omitempty"`
	SpeakFallback *[]Speak `json:"speak_fallback,omitempty"`
	Greeting      string   `json:"greeting,omitempty"`
}
