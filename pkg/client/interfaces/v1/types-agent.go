// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

/*
SettingsConfigurationOptions contain all of the knobs and dials to control the Agent API

Please see the live/streaming documentation for more details:
XXXX
*/
type SettingsConfigurationOptions struct {
	Type    string   `json:"type"`
	Audio   Audio    `json:"audio"`
	Agent   Agent    `json:"agent"`
	Context *Context `json:"context,omitempty"`
}

/*
Sub-structs in SettingsConfigurationOptions
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
type Listen struct {
	Model    string   `json:"model,omitempty"`
	Keyterms []string `json:"keyterms,omitempty"`
}
type Provider struct {
	Type string `json:"type,omitempty"`
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
	URL         *string    `json:"url,omitempty"`
	Headers     *[]Headers `json:"headers,omitempty"`
	Method      *string    `json:"method,omitempty"`
}
type Think struct {
	Provider     Provider    `json:"provider"`
	Model        string      `json:"model,omitempty"`
	Instructions string      `json:"instructions,omitempty"`
	Functions    []Functions `json:"functions,omitempty"`
}
type Speak struct {
	Model    string `json:"model,omitempty"`
	Provider string `json:"provider,omitempty"`
	VoiceID  string `json:"voice_id,omitempty"`
}
type Agent struct {
	Listen Listen `json:"listen"`
	Think  Think  `json:"think"`
	Speak  Speak  `json:"speak"`
}
type Context struct {
	Messages map[string]string `json:"messages,omitempty"`
	Replay   bool              `json:"replay,omitempty"`
}
