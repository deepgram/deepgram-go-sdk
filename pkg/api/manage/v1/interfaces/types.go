// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the types for the Deepgram Manage API.
*/
package interfaces

import (
	"time"

	"github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
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

// ErrorResponse is the Deepgram specific response error
type ErrorResponse interfaces.DeepgramError
