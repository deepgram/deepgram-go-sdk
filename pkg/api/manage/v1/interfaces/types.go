// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the types for the Deepgram Manage API.
*/
package interfaces

import (
	"encoding/json"
	"time"

	"github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

/***********************************/
// shared/common structs
/***********************************/
// Balance provides a balance
type Balance struct {
	BalanceID       string  `json:"balance_id,omitempty"`
	Amount          float64 `json:"amount,omitempty"`
	Units           string  `json:"units,omitempty"`
	PurchaseOrderID string  `json:"purchase_order_id,omitempty"`
}

// BalanceList provides a list of balances
type BalanceList struct {
	Balances []Balance `json:"balances,omitempty"`
}

// Invitation provides an invitation
type Invite struct {
	Email string `json:"email,omitempty"`
	Scope string `json:"scope,omitempty"`
}

// InvitationList provides a list of invitations
type InvitesList struct {
	Invites []Invite `json:"invites,omitempty"`
}

// APIKeyPermission Provides a user and key pairing
type APIKeyPermission struct {
	Member Member `json:"member,omitempty"`
	APIKey APIKey `json:"api_key,omitempty"`
}

// Key provides a key
type APIKey struct {
	APIKeyID string   `json:"api_key_id,omitempty"`
	Key      string   `json:"key,omitempty"`
	Comment  string   `json:"comment,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
	Created  string   `json:"created,omitempty"`
}

// KeyList provides a list of keys
type APIKeysList struct {
	APIKeys []APIKeyPermission `json:"api_keys,omitempty"`
}

// ScopeList provides a list of scopes
type ScopeList struct {
	Scopes []string `json:"scopes,omitempty"`
}

// Member provides a member
type Member struct {
	MemberID  string   `json:"member_id,omitempty"`
	Email     string   `json:"email,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
}

// MemberList provides a list of members
type MemberList struct {
	Members []Member `json:"members,omitempty"`
}

// Project provides a project
type Project struct {
	ProjectID string `json:"project_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Company   string `json:"company,omitempty"`
}

// ProjectList provides a list of projects
type ProjectList struct {
	Projects []Project `json:"projects,omitempty"`
}

// Stt provides a STT info
type Stt struct {
	Name            string   `json:"name,omitempty"`
	CanonicalName   string   `json:"canonical_name,omitempty"`
	Architecture    string   `json:"architecture,omitempty"`
	Languages       []string `json:"languages,omitempty"`
	Version         string   `json:"version,omitempty"`
	UUID            string   `json:"uuid,omitempty"`
	Batch           bool     `json:"batch,omitempty"`
	Streaming       bool     `json:"streaming,omitempty"`
	FormattedOutput bool     `json:"formatted_output,omitempty"`
}

// Metadata provides a metadata about a given model
type Metadata struct {
	Accent string   `json:"accent,omitempty"`
	Color  string   `json:"color,omitempty"`
	Image  string   `json:"image,omitempty"`
	Sample string   `json:"sample,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

// Tts provides a TTS info
type Tts struct {
	Name          string   `json:"name,omitempty"`
	CanonicalName string   `json:"canonical_name,omitempty"`
	Architecture  string   `json:"architecture,omitempty"`
	Languages     []string `json:"languages,omitempty"`
	Version       string   `json:"version,omitempty"`
	UUID          string   `json:"uuid,omitempty"`
	Metadata      Metadata `json:"metadata,omitempty"`
}

// Token provides a token
type Token struct {
	In  int `json:"in,omitempty"`
	Out int `json:"out,omitempty"`
}

// TTS provides a TTS
type TTS struct {
	Characters int `json:"characters,omitempty"`
	Requests   int `json:"requests,omitempty"`
}

// TokenDetails provides token details
type TokenDetails struct {
	Feature string `json:"feature,omitempty"`
	Input   int    `json:"input,omitempty"`
	Output  int    `json:"output,omitempty"`
	Model   string `json:"model,omitempty"`
}

// SpeechSegment provides a speech segment
type SpeechSegment struct {
	Characters int    `json:"characters,omitempty"`
	Model      string `json:"model,omitempty"`
	Tier       string `json:"tier,omitempty"`
}

// TTSDetails provides token details
type TTSDetails struct {
	Duration       float64         `json:"duration,omitempty"`
	SpeechSegments []SpeechSegment `json:"speech_segments,omitempty"`
	// TODO: audio_metadata
}

// Config provides a config
type Config struct {
	Language       string  `json:"language,omitempty"`
	Model          string  `json:"model,omitempty"`
	Punctuate      *bool   `json:"punctuate,omitempty"`
	Utterances     *bool   `json:"utterances,omitempty"`
	Diarize        *bool   `json:"diarize,omitempty"`
	SmartFormat    *bool   `json:"smart_format,omitempty"`
	InterimResults *bool   `json:"interim_results,omitempty"`
	Topics         *string `json:"topics,omitempty"`
	Intents        *string `json:"intents,omitempty"`
	Sentiment      *bool   `json:"sentiment,omitempty"`
	Summarize      *string `json:"summarize,omitempty"`
}

// Details provides details
type Details struct {
	Usd        float64  `json:"usd,omitempty"`
	Duration   float64  `json:"duration,omitempty"`
	TotalAudio float64  `json:"total_audio,omitempty"`
	Channels   int      `json:"channels,omitempty"`
	Streams    int      `json:"streams,omitempty"`
	Models     []string `json:"models,omitempty"`
	Method     string   `json:"method,omitempty"`
	Tier       string   `json:"tier,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Features   []string `json:"features,omitempty"`
	Config     Config   `json:"config,omitempty"`
}

// Response provides a response
type Response struct {
	Details      Details        `json:"details,omitempty"`
	Code         int            `json:"code,omitempty"`
	Completed    string         `json:"completed,omitempty"`
	TTSDetails   *TTSDetails    `json:"tts_details,omitempty"`
	TokenDetails []TokenDetails `json:"token_details,omitempty"`
}

type Callback struct {
	Attempts  int    `json:"attempts,omitempty"`
	Code      int    `json:"code,omitempty"`
	Completed string `json:"completed,omitempty"`
}

// Request provides a request
type Request struct {
	RequestID   string   `json:"request_id,omitempty"`
	ProjectUUID string   `json:"project_uuid,omitempty"`
	Created     string   `json:"created,omitempty"`
	Path        string   `json:"path,omitempty"`
	Accessor    string   `json:"accessor,omitempty"`
	APIKeyID    string   `json:"api_key_id,omitempty"`
	Response    Response `json:"response,omitempty"`
	Callback    Callback `json:"callback,omitempty"`
}

// Model provides a list of models
type Model struct {
	Name     string `json:"name,omitempty"`
	Language string `json:"language,omitempty"`
	Version  string `json:"version,omitempty"`
	ModelID  string `json:"model_id,omitempty"`
}

// Resolution provides a resolution
type Resolution struct {
	Units  string `json:"units,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

// Result provides a list of results
type Result struct {
	Start      string  `json:"start,omitempty"`
	End        string  `json:"end,omitempty"`
	Hours      float64 `json:"hours,omitempty"`
	TotalHours float64 `json:"total_hours,omitempty"`
	Requests   int     `json:"requests,omitempty"`
	Tokens     Token   `json:"tokens,omitempty"`
	TTS        TTS     `json:"tts,omitempty"`
}

// RequestList provides a list of requests
type RequestList struct {
	Page     int       `json:"page,omitempty"`
	Limit    int       `json:"limit,omitempty"`
	Requests []Request `json:"requests,omitempty"`
}

// UsageField provides a usage field
type UsageField struct {
	Tags              []string `json:"tags,omitempty"`
	Models            []Model  `json:"models,omitempty"`
	ProcessingMethods []string `json:"processing_methods,omitempty"`
	Features          []string `json:"features,omitempty"`
	Languages         []string `json:"languages,omitempty"`
}

// Usage provides a usage
type Usage struct {
	Start      string     `json:"start,omitempty"`
	End        string     `json:"end,omitempty"`
	Resolution Resolution `json:"resolution,omitempty"`
	Results    []Result   `json:"results,omitempty"`
}

// UsageBreakdownGrouping identifies the dimension used to group usage results.
type UsageBreakdownGrouping string

const (
	UsageBreakdownGroupingAccessor   UsageBreakdownGrouping = "accessor"
	UsageBreakdownGroupingEndpoint   UsageBreakdownGrouping = "endpoint"
	UsageBreakdownGroupingFeatureSet UsageBreakdownGrouping = "feature_set"
	UsageBreakdownGroupingModels     UsageBreakdownGrouping = "models"
	UsageBreakdownGroupingMethod     UsageBreakdownGrouping = "method"
	UsageBreakdownGroupingTags       UsageBreakdownGrouping = "tags"
	UsageBreakdownGroupingDeployment UsageBreakdownGrouping = "deployment"
)

// UsageBreakdownGroupingResult identifies the values for one usage group.
type UsageBreakdownGroupingResult struct {
	Start      string    `json:"start,omitempty"`
	End        string    `json:"end,omitempty"`
	Accessor   *string   `json:"accessor,omitempty"`
	Endpoint   *string   `json:"endpoint,omitempty"`
	FeatureSet *string   `json:"feature_set,omitempty"`
	Models     []*string `json:"models,omitempty"`
	Method     *string   `json:"method,omitempty"`
	Tags       []string  `json:"tags,omitempty"`
	Deployment *string   `json:"deployment,omitempty"`
}

// UsageBreakdownEntry contains the usage totals for one grouping.
type UsageBreakdownEntry struct {
	Hours         float64                      `json:"hours,omitempty"`
	TotalHours    float64                      `json:"total_hours,omitempty"`
	AgentHours    float64                      `json:"agent_hours,omitempty"`
	TokensIn      float64                      `json:"tokens_in,omitempty"`
	TokensOut     float64                      `json:"tokens_out,omitempty"`
	TTSCharacters float64                      `json:"tts_characters,omitempty"`
	Requests      float64                      `json:"requests,omitempty"`
	Grouping      UsageBreakdownGroupingResult `json:"grouping,omitempty"`
}

// UsageBreakdown provides project usage grouped by a reporting dimension.
type UsageBreakdown struct {
	Start      string                `json:"start,omitempty"`
	End        string                `json:"end,omitempty"`
	Resolution ReportingResolution   `json:"resolution,omitempty"`
	Results    []UsageBreakdownEntry `json:"results,omitempty"`
}

// BillingBreakdownGrouping identifies a dimension used to group billing results.
type BillingBreakdownGrouping string

const (
	BillingBreakdownGroupingAccessor   BillingBreakdownGrouping = "accessor"
	BillingBreakdownGroupingDeployment BillingBreakdownGrouping = "deployment"
	BillingBreakdownGroupingLineItem   BillingBreakdownGrouping = "line_item"
	BillingBreakdownGroupingTags       BillingBreakdownGrouping = "tags"
)

// BillingBreakdownGroupingResult identifies the values for one billing group.
type BillingBreakdownGroupingResult struct {
	Start      string   `json:"start,omitempty"`
	End        string   `json:"end,omitempty"`
	Accessor   *string  `json:"accessor,omitempty"`
	Deployment *string  `json:"deployment,omitempty"`
	LineItem   *string  `json:"line_item,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

// BillingBreakdownEntry contains the USD cost for one grouping.
type BillingBreakdownEntry struct {
	Dollars  float64                        `json:"dollars,omitempty"`
	Grouping BillingBreakdownGroupingResult `json:"grouping,omitempty"`
}

// BillingBreakdown provides project billing grouped by reporting dimensions.
type BillingBreakdown struct {
	Start      string                  `json:"start,omitempty"`
	End        string                  `json:"end,omitempty"`
	Resolution ReportingResolution     `json:"resolution,omitempty"`
	Results    []BillingBreakdownEntry `json:"results,omitempty"`
}

// ReportingResolution specifies the reporting period's time unit and amount.
type ReportingResolution struct {
	Units  string  `json:"units,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

// BillingFields lists values that can filter a billing breakdown.
type BillingFields struct {
	Accessors   []string          `json:"accessors,omitempty"`
	Deployments []string          `json:"deployments,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	LineItems   map[string]string `json:"line_items,omitempty"`
}

// Purchase contains the original amount on an order transaction.
type Purchase struct {
	OrderID    string     `json:"order_id,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
	Created    *time.Time `json:"created,omitempty"`
	Amount     float64    `json:"amount,omitempty"`
	Units      string     `json:"units,omitempty"`
	OrderType  string     `json:"order_type,omitempty"`
}

// Purchases lists original amounts on a project's order transactions.
type Purchases struct {
	Orders []Purchase `json:"orders,omitempty"`
}

// AgentConfiguration provides a reusable voice-agent configuration.
//
// Config is the "agent" block of a Settings message as a JSON *string*, not a
// nested JSON object: the server stores the exact string the create request
// sent and echoes it back verbatim, uninterpolated, so DG_<VARIABLE_NAME>
// placeholders appear as-is. The round trip is therefore symmetric with
// AgentCreateRequest.Config — what you send is what you read back.
//
// Metadata values are arbitrary JSON, not just strings: the API stores and
// returns whatever a client wrote, so an agent whose metadata was set by the
// console or another SDK can carry numbers, booleans, and nested objects.
// Typing this as map[string]string fails the decode of the *entire* response,
// not just the field.
//
// CreateAgent answers with the new UUID only, so AgentID is the single
// populated field on a create result; GetAgent and ListAgents fill in the rest.
type AgentConfiguration struct {
	AgentID    string                 `json:"agent_uuid,omitempty"`
	MemberID   string                 `json:"member_id,omitempty"`
	APIVersion int                    `json:"api_version,omitempty"`
	Config     string                 `json:"config,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// AgentsList provides a list of reusable voice-agent configurations.
//
// The list endpoint answers with a bare JSON array at the top level rather than
// an object wrapping one, so this type encodes and decodes as that array while
// still exposing the .Agents field every other list type in this package uses.
type AgentsList struct {
	Agents []AgentConfiguration
}

// UnmarshalJSON decodes the bare JSON array returned by the list agents endpoint.
func (l *AgentsList) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.Agents)
}

// MarshalJSON re-encodes the list as the bare JSON array the API uses, so a
// decode/encode round trip reproduces the wire format.
func (l AgentsList) MarshalJSON() ([]byte, error) {
	if l.Agents == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(l.Agents)
}

// AgentVariable provides a template variable for agent configurations.
//
// Value can be any valid JSON type (string, number, boolean, object, or array),
// and IsSensitive is reported on every read, so neither carries omitempty: a
// variable holding false, 0, or "" must survive a decode/encode round trip.
//
// CreateAgentVariable answers with the new UUID only, so VariableID is the
// single populated field on a create result.
type AgentVariable struct {
	VariableID  string      `json:"agent_variable_uuid,omitempty"`
	MemberID    string      `json:"member_id,omitempty"`
	APIVersion  int         `json:"api_version,omitempty"`
	Key         string      `json:"key,omitempty"`
	Value       interface{} `json:"value"`
	IsSensitive bool        `json:"is_sensitive"`
}

// AgentVariablesList provides a list of agent template variables.
//
// Like the agents list, the endpoint answers with a bare JSON array at the top
// level, so this type encodes and decodes as that array.
type AgentVariablesList struct {
	Variables []AgentVariable
}

// UnmarshalJSON decodes the bare JSON array returned by the list agent
// variables endpoint.
func (l *AgentVariablesList) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.Variables)
}

// MarshalJSON re-encodes the list as the bare JSON array the API uses.
func (l AgentVariablesList) MarshalJSON() ([]byte, error) {
	if l.Variables == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(l.Variables)
}

/***********************************/
// Request/Input structs
/***********************************/
// ProjectUpdateRequest provides a project update
type ProjectUpdateRequest struct {
	Name    string `json:"name,omitempty" url:"name,omitempty"`
	Company string `json:"company,omitempty" url:"company,omitempty"`
}

// ModelRequest provides a model request
type ModelRequest struct {
	IncludeOutdated bool `json:"include_outdated,omitempty" url:"include_outdated,omitempty"`
}

// InvitationRequest provides an invitation request
type InvitationRequest struct {
	Email string `json:"email,omitempty" url:"email,omitempty"`
	Scope string `json:"scope,omitempty" url:"scope,omitempty"`
}

// KeyCreateRequest provides a key create request
type KeyCreateRequest struct {
	Comment        string    `json:"comment,omitempty" url:"comment,omitempty"`
	Scopes         []string  `json:"scopes,omitempty" url:"scopes,omitempty"`
	ExpirationDate time.Time `json:"expiration_date,omitempty" url:"expiration_date,omitempty"`
	TimeToLive     int       `json:"time_to_live_in_seconds,omitempty" url:"time_to_live_in_seconds,omitempty"`
	Tags           []string  `json:"tags,omitempty" url:"tags,omitempty"`
}

// ScopeUpdateRequest provides a scope update request
type ScopeUpdateRequest struct {
	Scope string `json:"scope,omitempty" url:"scope,omitempty"`
}

// AgentCreateRequest provides an agent configuration create request.
// Config is a JSON string holding the "agent" block of a Settings message, not
// a nested JSON object. Template variables are referenced bare (unquoted)
// inside it — "prompt": DG_SYSTEM_PROMPT — so a config that uses them is not
// itself parseable JSON until the server substitutes their values.
// Metadata values are arbitrary JSON, matching the read side, so a caller can
// write back metadata it just read without losing non-string values.
type AgentCreateRequest struct {
	Config     string                 `json:"config"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	APIVersion int                    `json:"api_version,omitempty"`
}

// AgentMetadataUpdateRequest provides an agent metadata update request.
// Only metadata can be updated; the config itself is immutable — to change the
// configuration, delete the existing agent and create a new one.
// Metadata values are arbitrary JSON, matching what GetAgent returns, so a
// caller can round-trip metadata it did not author. The update REPLACES the
// entire metadata object: any key absent from this map is deleted, silently and
// with a 200. Read the current metadata with GetAgent first and resend every key
// you mean to keep.
type AgentMetadataUpdateRequest struct {
	Metadata map[string]interface{} `json:"metadata"`
}

// AgentVariableCreateRequest provides an agent variable create request.
// Key follows the DG_<VARIABLE_NAME> naming format; Value can be any valid
// JSON type (string, number, boolean, object, or array).
//
// IsSensitive is optional on the wire — the API defaults it to false when the
// field is absent — but false is the only value it accepts today. The field is
// serialized unconditionally (no omitempty) so the request states that intent
// explicitly rather than relying on a server-side default, and the Go zero
// value is the value callers want.
type AgentVariableCreateRequest struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	IsSensitive bool        `json:"is_sensitive"`
	APIVersion  int         `json:"api_version,omitempty"`
}

// AgentVariableUpdateRequest provides an agent variable update request
type AgentVariableUpdateRequest struct {
	Value interface{} `json:"value"`
}

// UsageListRequest provides a usage request
type UsageListRequest struct {
	Start  string `json:"start,omitempty" url:"start,omitempty"`
	End    string `json:"end,omitempty" url:"end,omitempty"`
	Page   int    `json:"page,omitempty" url:"page,omitempty"`
	Limit  int    `json:"limit,omitempty" url:"limit,omitempty"`
	Status string `json:"status,omitempty" url:"status,omitempty"`
}

// UsageRequest provides a usage request
type UsageRequest struct {
	Accessor           string   `json:"accessor,omitempty" url:"accessor,omitempty"`
	Alternatives       bool     `json:"alternatives,omitempty" url:"alternatives,omitempty"`
	AnalyzeSentiment   bool     `json:"analyze_sentiment,omitempty" url:"analyze_sentiment,omitempty"`
	DetectEntities     bool     `json:"detect_entities,omitempty" url:"detect_entities,omitempty"`
	DetectLanguage     bool     `json:"detect_language,omitempty" url:"detect_language,omitempty"`
	DetectTopics       bool     `json:"detect_topics,omitempty" url:"detect_topics,omitempty"`
	Diarize            bool     `json:"diarize,omitempty" url:"diarize,omitempty"`
	End                string   `json:"end,omitempty" url:"end,omitempty"`
	InterimResults     bool     `json:"interim_results,omitempty" url:"interim_results,omitempty"`
	Keywords           bool     `json:"keywords,omitempty" url:"keywords,omitempty"`
	Method             string   `json:"method,omitempty" url:"method,omitempty"` // Must be one of "sync" | "async" | "streaming"
	Model              string   `json:"model,omitempty" url:"model,omitempty"`
	Multichannel       bool     `json:"multichannel,omitempty" url:"multichannel,omitempty"`
	Ner                bool     `json:"ner,omitempty" url:"ner,omitempty"`
	Numbers            bool     `json:"numbers,omitempty" url:"numbers,omitempty"`
	Numerals           bool     `json:"numerals,omitempty" url:"numerals,omitempty"`
	Paragraphs         bool     `json:"paragraphs,omitempty" url:"paragraphs,omitempty"`
	ProfanityFilter    bool     `json:"profanity_filter,omitempty" url:"profanity_filter,omitempty"`
	Punctuate          bool     `json:"punctuate,omitempty" url:"punctuate,omitempty"`
	Redact             bool     `json:"redact,omitempty" url:"redact,omitempty"`
	Replace            bool     `json:"replace,omitempty" url:"replace,omitempty"`
	Search             bool     `json:"search,omitempty" url:"search,omitempty"`
	Sentiment          bool     `json:"sentiment,omitempty" url:"sentiment,omitempty"`
	SentimentThreshold float64  `json:"sentiment_threshold,omitempty" url:"sentiment_threshold,omitempty"`
	SmartFormat        bool     `json:"smart_format,omitempty" url:"smart_format,omitempty"`
	Start              string   `json:"start,omitempty" url:"start,omitempty"`
	Summarize          bool     `json:"summarize,omitempty" url:"summarize,omitempty"`
	Tag                []string `json:"tag,omitempty" url:"tag,omitempty"`
	Translate          bool     `json:"translate,omitempty" url:"translate,omitempty"`
	Utterances         bool     `json:"utterances,omitempty" url:"utterances,omitempty"`
	UttSplit           bool     `json:"utt_split,omitempty" url:"utt_split,omitempty"`
}

// UsageBreakdownRequest filters and groups usage breakdowns. Pointer boolean
// fields preserve an explicit false filter, which excludes requests using that
// feature according to the API contract.
type UsageBreakdownRequest struct {
	Start            string                 `schema:"start,omitempty"`
	End              string                 `schema:"end,omitempty"`
	Grouping         UsageBreakdownGrouping `schema:"grouping,omitempty"`
	Accessor         string                 `schema:"accessor,omitempty"`
	Alternatives     *bool                  `schema:"alternatives,omitempty"`
	CallbackMethod   *bool                  `schema:"callback_method,omitempty"`
	Callback         *bool                  `schema:"callback,omitempty"`
	Channels         *bool                  `schema:"channels,omitempty"`
	CustomIntentMode *bool                  `schema:"custom_intent_mode,omitempty"`
	CustomIntent     *bool                  `schema:"custom_intent,omitempty"`
	CustomTopicMode  *bool                  `schema:"custom_topic_mode,omitempty"`
	CustomTopic      *bool                  `schema:"custom_topic,omitempty"`
	Deployment       string                 `schema:"deployment,omitempty"`
	DetectEntities   *bool                  `schema:"detect_entities,omitempty"`
	DetectLanguage   *bool                  `schema:"detect_language,omitempty"`
	Diarize          *bool                  `schema:"diarize,omitempty"`
	Dictation        *bool                  `schema:"dictation,omitempty"`
	Encoding         *bool                  `schema:"encoding,omitempty"`
	Endpoint         string                 `schema:"endpoint,omitempty"`
	Extra            *bool                  `schema:"extra,omitempty"`
	FillerWords      *bool                  `schema:"filler_words,omitempty"`
	Intents          *bool                  `schema:"intents,omitempty"`
	Keyterm          *bool                  `schema:"keyterm,omitempty"`
	Keywords         *bool                  `schema:"keywords,omitempty"`
	Language         *bool                  `schema:"language,omitempty"`
	Measurements     *bool                  `schema:"measurements,omitempty"`
	Method           string                 `schema:"method,omitempty"`
	Model            string                 `schema:"model,omitempty"`
	Multichannel     *bool                  `schema:"multichannel,omitempty"`
	Numerals         *bool                  `schema:"numerals,omitempty"`
	Paragraphs       *bool                  `schema:"paragraphs,omitempty"`
	ProfanityFilter  *bool                  `schema:"profanity_filter,omitempty"`
	Punctuate        *bool                  `schema:"punctuate,omitempty"`
	Redact           *bool                  `schema:"redact,omitempty"`
	Replace          *bool                  `schema:"replace,omitempty"`
	SampleRate       *bool                  `schema:"sample_rate,omitempty"`
	Search           *bool                  `schema:"search,omitempty"`
	Sentiment        *bool                  `schema:"sentiment,omitempty"`
	SmartFormat      *bool                  `schema:"smart_format,omitempty"`
	Summarize        *bool                  `schema:"summarize,omitempty"`
	Tag              string                 `schema:"tag,omitempty"`
	Topics           *bool                  `schema:"topics,omitempty"`
	UttSplit         *bool                  `schema:"utt_split,omitempty"`
	Utterances       *bool                  `schema:"utterances,omitempty"`
	Version          *bool                  `schema:"version,omitempty"`
}

// BillingBreakdownRequest filters and groups billing breakdowns.
type BillingBreakdownRequest struct {
	Start      string                     `schema:"start,omitempty"`
	End        string                     `schema:"end,omitempty"`
	Accessor   string                     `schema:"accessor,omitempty"`
	Deployment string                     `schema:"deployment,omitempty"`
	Tag        string                     `schema:"tag,omitempty"`
	LineItem   string                     `schema:"line_item,omitempty"`
	Grouping   []BillingBreakdownGrouping `schema:"grouping,omitempty"`
}

// BillingFieldsRequest filters the reporting period for billing fields.
type BillingFieldsRequest struct {
	Start string `schema:"start,omitempty"`
	End   string `schema:"end,omitempty"`
}

// PurchasesRequest controls purchase-history pagination.
type PurchasesRequest struct {
	Limit int `schema:"limit,omitempty"`
}

/***********************************/
// Result/Output structs
/***********************************/
// BookmarksResult provides a generic message results
type MessageResult struct {
	Message string `json:"message"`
}

// BalanceResult provides a result with a list of balances
type BalancesResult struct {
	BalanceList
}

// BalanceResult provides a result with a single balance
type BalanceResult struct {
	Balance
}

// InvitationResult provides a result with a single invitation
type InvitationsResult struct {
	InvitesList
}

// InvitationResult provides a result with a single invitation
type KeysResult struct {
	APIKeysList
}

// KeyResult provides a result with a single key
type KeyResult struct {
	APIKeyPermission
}

// MembersResult provides a result with a list of members
type MembersResult struct {
	MemberList
}

// ProjectsResult provides a result with a list of projects
type ProjectsResult struct {
	ProjectList
}

// ProjectResult provides a result with a single project
type ProjectResult struct {
	Project
}

// ModelsResult provides a result with a list of models
type ModelsResult struct {
	Stt []Stt `json:"stt,omitempty"`
	Tts []Tts `json:"tts,omitempty"`
}

// ModelResult provides a result with a single model
type ModelResult struct {
	Name          string   `json:"name,omitempty"`
	CanonicalName string   `json:"canonical_name,omitempty"`
	Architecture  string   `json:"architecture,omitempty"`
	Language      string   `json:"language,omitempty"`
	Version       string   `json:"version,omitempty"`
	UUID          string   `json:"uuid,omitempty"`
	Metadata      Metadata `json:"metadata,omitempty"`
}

// ScopeResult provides a result with a list of scopes
type ScopeResult struct {
	ScopeList
}

// UsageResult provides a result with a list of usage
type UsageListResult struct {
	RequestList
}

// UsageRequestResult provides a result with a single usage request
type UsageRequestResult struct {
	Request
}

// UsageFieldResult provides a result with a list of fields
type UsageFieldResult struct {
	UsageField
}

// UsageSummary provides a result with a list of usage
type UsageResult struct {
	Usage
}

// UsageBreakdownResult provides a project usage breakdown.
type UsageBreakdownResult struct {
	UsageBreakdown
}

// BillingBreakdownResult provides a project billing breakdown.
type BillingBreakdownResult struct {
	BillingBreakdown
}

// BillingFieldsResult provides values that can filter billing breakdowns.
type BillingFieldsResult struct {
	BillingFields
}

// PurchasesResult provides project purchase history.
type PurchasesResult struct {
	Purchases
}

// AgentsResult provides a result with a list of agent configurations
type AgentsResult struct {
	AgentsList
}

// AgentResult provides a result with a single agent configuration
type AgentResult struct {
	AgentConfiguration
}

// AgentVariablesResult provides a result with a list of agent variables
type AgentVariablesResult struct {
	AgentVariablesList
}

// AgentVariableResult provides a result with a single agent variable
type AgentVariableResult struct {
	AgentVariable
}

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
