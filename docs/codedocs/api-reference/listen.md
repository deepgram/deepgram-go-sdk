---
title: "Listen"
description: "Reference for prerecorded speech-to-text REST clients and realtime Listen WebSocket clients."
---

The listen surface has two main entry points. `pkg/client/listen` exposes constructors for REST and WebSocket transports. `pkg/api/listen/v1/rest` exposes typed prerecorded transcription methods that return `*PreRecordedResponse`.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces`

## REST Constructors

Source: `pkg/client/listen/client.go`

```go
func NewRESTWithDefaults() *listenv1rest.RESTClient
func NewREST(apiKey string, options *interfaces.ClientOptions) *listenv1rest.RESTClient
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `apiKey` | `string` | `""` | Optional explicit API key. |
| `options` | `*interfaces.ClientOptions` | required | Host, auth, proxy, and transport settings. |

Example:

```go
c := listen.NewREST("", &interfaces.ClientOptions{
  Host: "https://api.deepgram.com",
})
```

## Typed Prerecorded API

Source: `pkg/api/listen/v1/rest/rest.go`

```go
func New(client *rest.Client) *Client
func (c *Client) FromFile(ctx context.Context, file string, options *interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error)
func (c *Client) FromStream(ctx context.Context, src io.Reader, options *interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error)
func (c *Client) FromURL(ctx context.Context, url string, options *interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error)
```

### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `ctx` | `context.Context` | — | Request context and optional custom headers or params. |
| `file` | `string` | — | Path to an audio file on disk. |
| `src` | `io.Reader` | — | Stream of audio bytes. |
| `url` | `string` | — | Remote audio URL. |
| `options` | `*interfaces.PreRecordedTranscriptionOptions` | `nil` or explicit struct | Query-string transcription options. |

### Return Type

Each method returns `(*github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest/interfaces.PreRecordedResponse, error)`.

### Usage Examples

```go
dg := listenapi.New(listen.NewRESTWithDefaults())
res, err := dg.FromFile(ctx, "call.wav", &interfaces.PreRecordedTranscriptionOptions{
  Model:       "nova-3",
  SmartFormat: true,
})
```

```go
f, _ := os.Open("call.wav")
res, err := dg.FromStream(ctx, f, &interfaces.PreRecordedTranscriptionOptions{
  Model:     "nova-3",
  Summarize: "v2",
})
```

```go
res, err := dg.FromURL(ctx, "https://dpgr.am/spacewalk.wav", &interfaces.PreRecordedTranscriptionOptions{
  Model:      "nova-3",
  Punctuate:  true,
  Utterances: true,
})
```

## Low-Level REST Transport

Source: `pkg/client/listen/v1/rest/client.go`

```go
func (c *Client) DoFile(ctx context.Context, filePath string, req *interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error
func (c *Client) DoStream(ctx context.Context, src io.Reader, options *interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error
func (c *Client) DoURL(ctx context.Context, audioURL string, options *interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error
func IsURL(str string) bool
```

Use these when you want to supply your own destination struct or custom raw-body handling.

## Realtime Constructors

Source: `pkg/client/listen/client.go`

```go
func NewWSUsingCallbackForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSCallback, error)
func NewWSUsingCallbackWithDefaults(ctx context.Context, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error)
func NewWSUsingCallback(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error)
func NewWSUsingCallbackWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, callback msginterfaces.LiveMessageCallback) (*listenv1ws.WSCallback, error)

func NewWSUsingChanForDemo(ctx context.Context, options *interfaces.LiveTranscriptionOptions) (*listenv1ws.WSChannel, error)
func NewWSUsingChanWithDefaults(ctx context.Context, options *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error)
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error)
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.LiveTranscriptionOptions, chans msginterfaces.LiveMessageChan) (*listenv1ws.WSChannel, error)
```

## WebSocket Methods

Source: `pkg/client/listen/v1/websocket/client_callback.go`, `pkg/client/listen/v1/websocket/client_channel.go`

```go
func (c *WSCallback) Connect() bool
func (c *WSCallback) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool
func (c *WSCallback) AttemptReconnect(ctx context.Context, retries int64) bool
func (c *WSCallback) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool
func (c *WSCallback) Stream(r io.Reader) error
func (c *WSCallback) Write(p []byte) (int, error)
func (c *WSCallback) KeepAlive() error
func (c *WSCallback) Finalize() error

func (c *WSChannel) Connect() bool
func (c *WSChannel) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool
func (c *WSChannel) AttemptReconnect(ctx context.Context, retries int64) bool
func (c *WSChannel) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool
func (c *WSChannel) Stream(r io.Reader) error
func (c *WSChannel) Write(p []byte) (int, error)
func (c *WSChannel) KeepAlive() error
func (c *WSChannel) Finalize() error
```

Example:

```go
dg, _ := listen.NewWSUsingCallback(ctx, "", cOptions, tOptions, callback)
if dg.Connect() {
  _ = dg.Stream(reader)
  _ = dg.Finalize()
}
```

## Event Interfaces

Source: `pkg/api/listen/v1/websocket/interfaces/interfaces.go`

```go
type LiveMessageChan interface {
  GetOpen() []*chan *OpenResponse
  GetMessage() []*chan *MessageResponse
  GetMetadata() []*chan *MetadataResponse
  GetSpeechStarted() []*chan *SpeechStartedResponse
  GetUtteranceEnd() []*chan *UtteranceEndResponse
  GetClose() []*chan *CloseResponse
  GetError() []*chan *ErrorResponse
  GetUnhandled() []*chan *[]byte
}

type LiveMessageCallback interface {
  Open(or *OpenResponse) error
  Message(mr *MessageResponse) error
  Metadata(md *MetadataResponse) error
  SpeechStarted(ssr *SpeechStartedResponse) error
  UtteranceEnd(ur *UtteranceEndResponse) error
  Close(cr *CloseResponse) error
  Error(er *ErrorResponse) error
  UnhandledEvent(byData []byte) error
}
```

Related pages: [Client Interfaces](/workspace/home/codedocs-template/content/docs/api-reference/client-interfaces.mdx), [Realtime Streaming](/workspace/home/codedocs-template/content/docs/realtime-streaming.mdx)
