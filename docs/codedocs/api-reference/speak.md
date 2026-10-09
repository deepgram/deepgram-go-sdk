---
title: "Speak"
description: "Reference for Deepgram text-to-speech REST and WebSocket clients."
---

The speak surface supports two different workflows: REST for full audio responses and WebSockets for incremental, interactive synthesis. The public constructors live in `pkg/client/speak`, while the typed REST wrapper lives in `pkg/api/speak/v1/rest`.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v1/rest`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v1/websocket/interfaces`

## REST Constructors

Source: `pkg/client/speak/client.go`

```go
func NewRESTWithDefaults() *speakv1rest.RESTClient
func NewREST(apiKey string, options *interfaces.ClientOptions) *speakv1rest.RESTClient
```

## Typed REST API

Source: `pkg/api/speak/v1/rest/speak.go`

```go
func New(client *speak.Client) *Client
func (c *Client) ToStream(ctx context.Context, text string, options *interfaces.SpeakOptions, buf *interfaces.RawResponse) (*api.SpeakResponse, error)
func (c *Client) ToFile(ctx context.Context, text string, options *interfaces.SpeakOptions, w io.Writer) (*api.SpeakResponse, error)
func (c *Client) ToSave(ctx context.Context, filename, text string, options *interfaces.SpeakOptions) (*api.SpeakResponse, error)
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `ctx` | `context.Context` | — | Request context. |
| `text` | `string` | — | Text to synthesize. |
| `options` | `*interfaces.SpeakOptions` | required | TTS model and audio settings. |
| `buf` | `*interfaces.RawResponse` | — | In-memory destination for audio bytes. |
| `w` | `io.Writer` | — | Arbitrary writer destination. |
| `filename` | `string` | — | Output path created by `ToSave`. |

Example:

```go
dg := speakapi.New(speak.NewRESTWithDefaults())
resp, err := dg.ToSave(ctx, "hello.wav", "Hello from Deepgram", &interfaces.SpeakOptions{
  Model:      "aura-2-thalia-en",
  Encoding:   "linear16",
  SampleRate: 48000,
})
```

## Low-Level REST Transport

Source: `pkg/client/speak/v1/rest/client.go`

```go
func (c *Client) DoText(ctx context.Context, text string, options *interfaces.SpeakOptions, keys []string, resBody interface{}) (map[string]string, error)
```

`DoText` is the underlying method used by `ToStream`, `ToFile`, and `ToSave`. It captures response headers like `request-id` and `char-count`, which the typed wrapper converts into `SpeakResponse`.

## WebSocket Constructors

Source: `pkg/client/speak/client.go`

```go
func NewWSUsingCallbackForDemo(ctx context.Context, options *interfaces.WSSpeakOptions) (*speakv1ws.WSCallback, error)
func NewWSUsingCallbackWithDefaults(ctx context.Context, options *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.WSCallback, error)
func NewWSUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.WSCallback, error)
func NewWSUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageCallback) (*speakv1ws.WSCallback, error)

func NewWSUsingChanForDemo(ctx context.Context, options *interfaces.WSSpeakOptions) (*speakv1ws.WSChannel, error)
func NewWSUsingChanWithDefaults(ctx context.Context, options *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageChan) (*speakv1ws.WSChannel, error)
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageChan) (*speakv1ws.WSChannel, error)
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, sOptions *interfaces.WSSpeakOptions, callback msginterfaces.SpeakMessageChan) (*speakv1ws.WSChannel, error)
```

## WebSocket Methods

Source: `pkg/client/speak/v1/websocket/client_callback.go`, `pkg/client/speak/v1/websocket/client_channel.go`

```go
func (c *WSCallback) Connect() bool
func (c *WSCallback) SpeakWithText(text string) error
func (c *WSCallback) Speak(text string) error
func (c *WSCallback) Flush() error
func (c *WSCallback) Reset() error

func (c *WSChannel) Connect() bool
func (c *WSChannel) SpeakWithText(text string) error
func (c *WSChannel) Speak(text string) error
func (c *WSChannel) Flush() error
func (c *WSChannel) Reset() error
```

Example:

```go
dg, _ := speak.NewWSUsingCallback(ctx, "", cOptions, wsOptions, callback)
if dg.Connect() {
  _ = dg.SpeakWithText("First sentence.")
  _ = dg.SpeakWithText("Second sentence.")
  _ = dg.Flush()
}
```

## WebSocket Event Interfaces

Source: `pkg/api/speak/v1/websocket/interfaces/interfaces.go`

```go
type SpeakMessageCallback interface {
  Open(or *OpenResponse) error
  Metadata(md *MetadataResponse) error
  Flush(fl *FlushedResponse) error
  Clear(cl *ClearedResponse) error
  Close(cr *CloseResponse) error
  Warning(er *WarningResponse) error
  Error(er *ErrorResponse) error
  UnhandledEvent(byMsg []byte) error
  Binary(byMsg []byte) error
}
```

The binary callback is where synthesized audio bytes arrive.
