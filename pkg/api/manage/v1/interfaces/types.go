// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the types for the Deepgram Manage API.
*/
package interfaces

import "time"

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

// Config provides a config
type Config struct {
	Diarize        bool   `json:"diarize,omitempty"`
	Language       string `json:"language,omitempty"`
	Model          string `json:"model,omitempty"`
	Punctuate      bool   `json:"punctuate,omitempty"`
	Utterances     bool   `json:"utterances,omitempty"`
	InterimResults bool   `json:"interim_results,omitempty"`
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
	Details   Details `json:"details,omitempty"`
	Code      int     `json:"code,omitempty"`
	Completed string  `json:"completed,omitempty"`
}

// Request provides a request
type Request struct {
	RequestID string      `json:"request_id,omitempty"`
	Created   string      `json:"created,omitempty"`
	Path      string      `json:"path,omitempty"`
	Accessor  string      `json:"accessor" url:"accessor,omitempty"`
	APIKeyID  string      `json:"api_key_id,omitempty"`
	Response  Response    `json:"response,omitempty"`
	Callback  interface{} `json:"callback,omitempty"`
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
}

// RequestList provides a list of requests
type RequestList struct {
	Page     int       `json:"page,omitempty"`
	Limit    int       `json:"limit,omitempty"`
	Requests []Request `json:"requests,omitempty"`
}

// UsageField provides a usage field
type UsageField struct {
	Tags              []any    `json:"tags,omitempty"`
	Models            []Model  `json:"models,omitempty"`
	ProcessingMethods []string `json:"processing_methods,omitempty"`
	Features          []string `json:"features,omitempty"`
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
	Name    string `json:"name,omitempty"`
	Company string `json:"company,omitempty"`
}

// InvitationRequest provides an invitation request
type InvitationRequest struct {
	Email string `json:"email"`
	Scope string `json:"scope"`
}

// KeyCreateRequest provides a key create request
type KeyCreateRequest struct {
	Comment        string    `json:"comment"`
	Scopes         []string  `json:"scopes"`
	ExpirationDate time.Time `json:"expiration_date"`
	TimeToLive     int       `json:"time_to_live,omitempty"`
	Tags           []string  `json:"tags"`
}

// ScopeUpdateRequest provides a scope update request
type ScopeUpdateRequest struct {
	Scope string `json:"scope"`
}

// UsageListRequest provides a usage request
type UsageListRequest struct {
	Start  string `json:"start" url:"start,omitempty"`
	End    string `json:"end" url:"end,omitempty"`
	Page   int    `json:"page" url:"page,omitempty"`
	Limit  int    `json:"limit" url:"limit,omitempty"`
	Status string `json:"status" url:"status,omitempty"`
}

// UsageRequest provides a usage request
type UsageRequest struct {
	Accessor           string   `json:"accessor" url:"accessor,omitempty"`
	Alternatives       bool     `json:"alternatives" url:"alternatives,omitempty"`
	AnalyzeSentiment   bool     `json:"analyze_sentiment" url:"analyze_sentiment,omitempty"`
	DetectEntities     bool     `json:"detect_entities" url:"detect_entities,omitempty"`
	DetectLanguage     bool     `json:"detect_language" url:"detect_language,omitempty"`
	DetectTopics       bool     `json:"detect_topics" url:"detect_topics,omitempty"`
	Diarize            bool     `json:"diarize" url:"diarize,omitempty"`
	End                string   `json:"end" url:"end,omitempty"`
	InterimResults     bool     `json:"interim_results" url:"interim_results,omitempty"`
	Keywords           bool     `json:"keywords" url:"keywords,omitempty"`
	Method             string   `json:"method" url:"method,omitempty"` // Must be one of "sync" | "async" | "streaming"
	Model              string   `json:"model" url:"model,omitempty"`
	Multichannel       bool     `json:"multichannel" url:"multichannel,omitempty"`
	Ner                bool     `json:"ner" url:"ner,omitempty"`
	Numbers            bool     `json:"numbers" url:"numbers,omitempty"`
	Numerals           bool     `json:"numerals" url:"numerals,omitempty"`
	Paragraphs         bool     `json:"paragraphs" url:"paragraphs,omitempty"`
	ProfanityFilter    bool     `json:"profanity_filter" url:"profanity_filter,omitempty"`
	Punctuate          bool     `json:"punctuate" url:"punctuate,omitempty"`
	Redact             bool     `json:"redact" url:"redact,omitempty"`
	Replace            bool     `json:"replace" url:"replace,omitempty"`
	Search             bool     `json:"search" url:"search,omitempty"`
	Sentiment          bool     `json:"sentiment" url:"sentiment,omitempty"`
	SentimentThreshold float64  `json:"sentiment_threshold" url:"sentiment_threshold,omitempty"`
	SmartFormat        bool     `json:"smart_format" url:"smart_format,omitempty"`
	Start              string   `json:"start" url:"start,omitempty"`
	Summarize          bool     `json:"summarize" url:"summarize,omitempty"`
	Tag                []string `json:"tag" url:"tag,omitempty"`
	Translate          bool     `json:"translate" url:"translate,omitempty"`
	Utterances         bool     `json:"utterances" url:"utterances,omitempty"`
	UttSplit           bool     `json:"utt_split" url:"utt_split,omitempty"`
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
