package deepgram

import (
	"bytes"
	"io"
)

type InvitationOptions struct {
	Email string `json:"email"`
	Scope string `json:"scope"`
}

type InvitationList struct {
	Invites []InvitationOptions `json:"invites"`
}

type Message struct {
	Message string `json:"message"`
}

type TranscriptionSource interface {
	ReadStreamSource | UrlSource | BufferSource
}

type ReadStreamSource struct {
	Stream   io.Reader `json:"stream"`
	Mimetype string    `json:"mimetype"`
}

type UrlSource struct {
	Url string `json:"url"`
}

type BufferSource struct {
	Buffer   bytes.Buffer `json:"buffer"`
	Mimetype string       `json:"mimetype"`
}
