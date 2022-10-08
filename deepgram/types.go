package deepgram

import (
	"bytes"
	"io"
)

type InvitationOptions struct {
  Email string `json:"email"`
  Scope string `json:"scope"`
};

type InvitationList struct{
 	Invites []InvitationOptions `json:"invites"`
};

type Message struct {
	Message string `json:"message"`
}

type LiveTranscriptionOptions struct {
	Model string `json:"model" url:"model,omitempty" `
	Language string `json:"language" url:"language,omitempty" `
	Version string `json:"version" url:"version,omitempty" `
	Punctuate bool `json:"punctuate" url:"punctuate,omitempty" `
	Profanity_filter bool `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Redact bool `json:"redact" url:"redact,omitempty" `
	Diarize bool `json:"diarize" url:"diarize,omitempty" `
	Diarize_version string `json:"diarize_version" url:"diarize_version,omitempty" `
	Multichannel bool `json:"multichannel" url:"multichannel,omitempty" `
	Alternatives int `json:"alternatives" url:"alternatives,omitempty" `
	Numerals bool `json:"numerals" url:"numerals,omitempty" `
	Search []string `json:"search" url:"search,omitempty" `
	Callback string `json:"callback" url:"callback,omitempty" `
	Keywords []string `json:"keywords" url:"keywords,omitempty" `
	Interim_results bool `json:"interim_results" url:"interim_results,omitempty" `
	Endpointing bool `json:"endpointing" url:"endpointing,omitempty" `
	Vad_turnoff int `json:"vad_turnoff" url:"vad_turnoff,omitempty" `
	Encoding string `json:"encoding" url:"encoding,omitempty" `
	Channels int `json:"channels" url:"channels,omitempty" `
	Sample_rate int `json:"sample_rate" url:"sample_rate,omitempty" `
	Tier string `json:"tier" url:"tier,omitempty" `
	Replace string `json:"replace" url:"replace,omitempty" `
}

type PreRecordedTranscriptionOptions struct {
	Tier string `json:"tier" url:"tier,omitempty" `
	Model string `json:"model" url:"model,omitempty" `
	Version string `json:"version" url:"version,omitempty" `
	Language string `json:"language" url:"language,omitempty" `
	Punctuate bool `json:"punctuate" url:"punctuate,omitempty" `
	Profanity_filter bool `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Redact bool `json:"redact" url:"redact,omitempty" `
	Diarize bool `json:"diarize" url:"diarize,omitempty" `
	Diarize_version string `json:"diarize_version" url:"diarize_version,omitempty" `
	Ner bool `json:"ner" url:"ner,omitempty" `
	Multichannel bool `json:"multichannel" url:"multichannel,omitempty" `
	Alternatives int `json:"alternatives" url:"alternatives,omitempty" `
	Numerals bool `json:"numerals" url:"numerals,omitempty" `
	Search []string `json:"search" url:"search,omitempty" `
	Replace string `json:"replace" url:"replace,omitempty" `
	Callback string `json:"callback" url:"callback,omitempty" `
	Keywords []string `json:"keywords" url:"keywords,omitempty" `
	Utterances bool `json:"utterances" url:"utterances,omitempty" `
	Utt_split int `json:"utt_split" url:"utt_split,omitempty" `
	Tag string `json:"tag" url:"tag,omitempty"`
}

type TranscriptionSource interface {
	 ReadStreamSource | UrlSource | BufferSource
}

type ReadStreamSource struct {
	Stream io.Reader `json:"stream"`
	Mimetype string `json:"mimetype"`
}

type UrlSource struct {
	Url string `json:"url"`
}

type BufferSource struct {
	Buffer bytes.Buffer `json:"buffer"`
	Mimetype string `json:"mimetype"`
}