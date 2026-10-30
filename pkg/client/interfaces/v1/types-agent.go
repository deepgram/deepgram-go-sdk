// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

import "encoding/json"

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
	//nolint:stylecheck // Url is preserved for backwards compatibility.
	Url     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Method  string            `json:"method,omitempty"`
}
type Item struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}
type Properties struct {
	// Deprecated: use Parameters.PropertySchemas for arbitrary JSON Schema properties.
	Item Item `json:"item,omitempty"`
}
type Parameters struct {
	Type string `json:"type,omitempty"`
	// Deprecated: use PropertySchemas for arbitrary JSON Schema properties.
	Properties      Properties             `json:"properties,omitempty"`
	PropertySchemas map[string]interface{} `json:"-"`
	Required        []string               `json:"required,omitempty"`
}

// MarshalJSON sends arbitrary property schemas when provided while preserving
// the legacy Properties representation for existing callers.
//
//nolint:gocritic // A value receiver preserves marshaling for value fields.
func (p Parameters) MarshalJSON() ([]byte, error) {
	type parametersJSON struct {
		Type       string      `json:"type,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Required   []string    `json:"required,omitempty"`
	}

	properties := interface{}(p.Properties)
	if p.PropertySchemas != nil {
		properties = p.PropertySchemas
	}

	return json.Marshal(parametersJSON{
		Type:       p.Type,
		Properties: properties,
		Required:   p.Required,
	})
}

type Headers struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
type Functions struct {
	Name          string     `json:"name,omitempty"`
	Description   string     `json:"description,omitempty"`
	Parameters    Parameters `json:"parameters,omitempty"`
	Endpoint      Endpoint   `json:"endpoint,omitempty"`
	DeferUntilEOT bool       `json:"defer_until_eot,omitempty"`
}

// MarshalJSON omits an unset endpoint so functions without a server endpoint
// are correctly registered as client-side functions.
//
//nolint:gocritic // A value receiver preserves marshaling for value fields.
func (f Functions) MarshalJSON() ([]byte, error) {
	// functionJSON mirrors Functions. Add any new Functions field here too, or it
	// is silently dropped from the wire.
	type functionJSON struct {
		Name          string     `json:"name,omitempty"`
		Description   string     `json:"description,omitempty"`
		Parameters    Parameters `json:"parameters,omitempty"`
		Endpoint      *Endpoint  `json:"endpoint,omitempty"`
		DeferUntilEOT bool       `json:"defer_until_eot,omitempty"`
	}

	var endpoint *Endpoint
	if f.Endpoint.Url != "" || len(f.Endpoint.Headers) != 0 || f.Endpoint.Method != "" {
		endpoint = &f.Endpoint
	}

	return json.Marshal(functionJSON{
		Name:          f.Name,
		Description:   f.Description,
		Parameters:    f.Parameters,
		Endpoint:      endpoint,
		DeferUntilEOT: f.DeferUntilEOT,
	})
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

// ThinkLiveProvider configures Deepgram's Think Live model. Think Live runs
// alongside the reasoning model configured by Think.
type ThinkLiveProvider struct {
	Type                        string   `json:"type"`
	Version                     string   `json:"version"`
	AgentBackchannelSensitivity *float64 `json:"agent_backchannel_sensitivity,omitempty"`
	PresenceCheckMaxCount       *int     `json:"presence_check_max_count,omitempty"`
	PresenceCheckIntervalMS     *int     `json:"presence_check_interval_ms,omitempty"`
	FrustrationMaxCount         *int     `json:"frustration_max_count,omitempty"`
	RestrictedTerms             []string `json:"restricted_terms,omitempty"`
}

// MarshalJSON preserves an explicit empty restricted_terms array, which clears
// the service default, while omitting nil values.
func (p ThinkLiveProvider) MarshalJSON() ([]byte, error) {
	type providerJSON struct {
		Type                        string    `json:"type"`
		Version                     string    `json:"version"`
		AgentBackchannelSensitivity *float64  `json:"agent_backchannel_sensitivity,omitempty"`
		PresenceCheckMaxCount       *int      `json:"presence_check_max_count,omitempty"`
		PresenceCheckIntervalMS     *int      `json:"presence_check_interval_ms,omitempty"`
		FrustrationMaxCount         *int      `json:"frustration_max_count,omitempty"`
		RestrictedTerms             *[]string `json:"restricted_terms,omitempty"`
	}

	var restrictedTerms *[]string
	if p.RestrictedTerms != nil {
		restrictedTerms = &p.RestrictedTerms
	}
	return json.Marshal(providerJSON{
		Type:                        p.Type,
		Version:                     p.Version,
		AgentBackchannelSensitivity: p.AgentBackchannelSensitivity,
		PresenceCheckMaxCount:       p.PresenceCheckMaxCount,
		PresenceCheckIntervalMS:     p.PresenceCheckIntervalMS,
		FrustrationMaxCount:         p.FrustrationMaxCount,
		RestrictedTerms:             restrictedTerms,
	})
}

// ThinkLiveSettings enables and configures Deepgram's Think Live model.
type ThinkLiveSettings struct {
	Provider *ThinkLiveProvider `json:"provider,omitempty"`
}

// NewThinkLiveSettings creates a Think Live configuration with its required
// Deepgram provider and version.
func NewThinkLiveSettings() *ThinkLiveSettings {
	return &ThinkLiveSettings{
		Provider: &ThinkLiveProvider{
			Type:    "deepgram",
			Version: "v1",
		},
	}
}

type Speak struct {
	Provider map[string]interface{} `json:"provider,omitempty"`
	Endpoint *Endpoint              `json:"endpoint,omitempty"`
}
type Agent struct {
	Language      string             `json:"language,omitempty"`
	Listen        Listen             `json:"listen,omitempty"`
	Think         Think              `json:"think,omitempty"`
	ThinkLive     *ThinkLiveSettings `json:"think_live,omitempty"`
	Speak         Speak              `json:"speak,omitempty"`
	SpeakFallback *[]Speak           `json:"speak_fallback,omitempty"`
	Greeting      string             `json:"greeting,omitempty"`
}
