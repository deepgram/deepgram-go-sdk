---
title: "Client Interfaces"
description: "Reference for the shared option, error, and helper types used across the Deepgram Go SDK."
---

The shared types in `pkg/client/interfaces` and `pkg/client/interfaces/v1` define the configuration contract for the rest of the SDK. If you understand these types, every constructor page becomes easier to read.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1`

The unversioned package mostly aliases the v1 package and adds `NewSettingsConfigurationOptions()`.

## `ClientOptions`

Source: `pkg/client/interfaces/v1/types-client.go`, `pkg/client/interfaces/v1/options.go`

```go
type ClientOptions struct {
  APIKey string
  AccessToken string
  Host string
  APIVersion string
  Path string
  SelfHosted bool
  Proxy func(*http.Request) (*url.URL, error)
  WSHeaderProcessor func(http.Header)
  SkipServerAuth bool
  RedirectService bool
  EnableKeepAlive bool
  AutoFlushReplyDelta int64
  AutoFlushSpeakDelta int64
}
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `APIKey` | `string` | `DEEPGRAM_API_KEY` if unset explicitly | Token auth credential. |
| `AccessToken` | `string` | `DEEPGRAM_ACCESS_TOKEN` if unset explicitly | Bearer auth credential; takes priority over `APIKey`. |
| `Host` | `string` | Deepgram default host | Override the API host or self-hosted endpoint. |
| `APIVersion` | `string` | product default | Override version selection in URI builders. |
| `Path` | `string` | product default | Override the product path segment. |
| `SelfHosted` | `bool` | `false` | Allows missing auth when talking to self-hosted deployments. |
| `Proxy` | `func(*http.Request) (*url.URL, error)` | `nil` | Custom HTTP or WS proxy function. |
| `WSHeaderProcessor` | `func(http.Header)` | `nil` | Final mutation hook for WebSocket headers. |
| `SkipServerAuth` | `bool` | `false` | Disables server-side TLS verification for WebSockets. |
| `RedirectService` | `bool` | `false` | Allows redirects during WebSocket dial. |
| `EnableKeepAlive` | `bool` | `false` | Starts product keepalive loops where supported. |
| `AutoFlushReplyDelta` | `int64` | `0` | Listen auto-finalization threshold in milliseconds. |
| `AutoFlushSpeakDelta` | `int64` | `0` | Speak auto-flush threshold in milliseconds. |

### Methods

```go
func (o *ClientOptions) Parse() error
func (o *ClientOptions) SetAccessToken(accessToken string)
func (o *ClientOptions) SetAPIKey(apiKey string)
func (o *ClientOptions) GetAuthToken() (token string, isBearer bool)
func (o *ClientOptions) InspectListenMessage() bool
func (o *ClientOptions) InspectSpeakMessage() bool
```

Example:

```go
opts := &interfaces.ClientOptions{
  AccessToken: "temporary-jwt",
  Host:        "https://api.deepgram.com",
}
token, isBearer := opts.GetAuthToken()
_, _ = token, isBearer
```

## Core Request Option Types

Source: `pkg/client/interfaces/v1/types-prerecorded.go`, `types-stream.go`, `types-analyze.go`, `types-speak.go`, `types-agent.go`

```go
type PreRecordedTranscriptionOptions struct { ... }
type LiveTranscriptionOptions struct { ... }
type AnalyzeOptions struct { ... }
type SpeakOptions struct { ... }
type WSSpeakOptions struct { ... }
type SettingsOptions struct { ... }
```

### `PreRecordedTranscriptionOptions`

Key fields include `Model`, `Language`, `SmartFormat`, `Punctuate`, `Diarize`, `Utterances`, `Sentiment`, `Summarize`, `Topics`, `Intents`, `Redact`, and `Keyterm`. These are serialized into REST query parameters by the URI helpers in `pkg/api/version`.

Example:

```go
opts := &interfaces.PreRecordedTranscriptionOptions{
  Model:       "nova-3",
  SmartFormat: true,
  Diarize:     true,
  Utterances:  true,
}
```

### `LiveTranscriptionOptions`

Key realtime fields include `Encoding`, `SampleRate`, `InterimResults`, `UtteranceEndMs`, `VadEvents`, `NoDelay`, and `Endpointing`.

Example:

```go
opts := &interfaces.LiveTranscriptionOptions{
  Model:          "nova-3",
  Encoding:       "linear16",
  SampleRate:     16000,
  InterimResults: true,
}
```

### `AnalyzeOptions`

The analysis surface is intentionally compact: `Language`, `Summarize`, `Sentiment`, `Topics`, `Intents`, and custom topic or intent lists.

### `SpeakOptions` and `WSSpeakOptions`

REST uses `SpeakOptions`, which includes callback fields and container selection. WebSocket TTS uses `WSSpeakOptions`, which omits callback and container fields because the protocol is interactive.

### `SettingsOptions`

Voice Agent settings contain nested `Audio`, `Agent`, `Listen`, `Think`, `Speak`, `Endpoint`, and `Functions` definitions. Use `interfaces.NewSettingsConfigurationOptions()` to start from a valid default object.

## Helper Types

Source: `pkg/client/interfaces/v1/utils.go`

```go
func WithCustomHeaders(ctx context.Context, headers http.Header) context.Context
func WithCustomParameters(ctx context.Context, params map[string][]string) context.Context

type RawResponse struct { bytes.Buffer }
type DeepgramWarning struct { ... }
type DeepgramError struct { ... }
type StatusError struct {
  Resp *http.Response
  DeepgramError *DeepgramError
}
```

`WithCustomHeaders` and `WithCustomParameters` let you attach extra transport metadata without extending every method signature. `RawResponse` is most useful with Speak REST when you want the raw audio bytes in memory.

Example:

```go
headers := http.Header{}
headers.Set("X-Correlation-ID", "req-123")
ctx := interfaces.WithCustomHeaders(context.Background(), headers)
```
