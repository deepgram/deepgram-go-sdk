package deepgram

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
	Model string `json:"model"`
	Language string `json:"language"`
	Version string `json:"version"`
	Punctuate bool `json:"punctuate"`
	Profanity_filter bool `json:"profanity_filter"`
	Redact bool `json:"redact"`
	Diarize bool `json:"diarize"`
	Multichannel bool `json:"multichannel"`
	Alternatives int `json:"alternatives"`
	Numerals bool `json:"numerals"`
	Search []string `json:"search"`
	Callback string `json:"callback"`
	Keywords []string `json:"keywords"`
	Interim_results bool `json:"interim_results"`
	Endpointing bool `json:"endpointing"`
	Vad_turnoff int `json:"vad_turnoff"`
	Encoding string `json:"encoding"`
	Channels int `json:"channels"`
	Sample_rate int `json:"sample_rate"`
}