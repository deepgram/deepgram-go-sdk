---
title: "Analyze"
description: "Reference for Deepgram's text intelligence analysis client and typed API wrapper."
---

The analyze surface converts text, files, streams, or URLs into typed sentiment, topic, summary, and intent results. The low-level transport lives in `pkg/client/analyze/v1/client.go`; the typed wrapper lives in `pkg/api/analyze/v1/analyze.go`.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/analyze`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/analyze/v1`

## Constructors

Source: `pkg/client/analyze/client.go`

```go
func NewWithDefaults() *Client
func New(apiKey string, options *interfaces.ClientOptions) *Client
```

## Typed API Wrapper

Source: `pkg/api/analyze/v1/analyze.go`

```go
func New(client *analyze.Client) *Client
func (c *Client) FromFile(ctx context.Context, file string, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error)
func (c *Client) FromStream(ctx context.Context, src io.Reader, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error)
func (c *Client) FromURL(ctx context.Context, url string, options *interfaces.AnalyzeOptions) (*api.AnalyzeResponse, error)
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `ctx` | `context.Context` | — | Request-scoped context. |
| `file` | `string` | — | Text file path. |
| `src` | `io.Reader` | — | Reader for text input. |
| `url` | `string` | — | URL pointing to text content. |
| `options` | `*interfaces.AnalyzeOptions` | `nil` or explicit struct | Analysis toggles such as `Sentiment`, `Topics`, `Intents`, and `Summarize`. |

Example:

```go
dg := analyzeapi.New(analyze.NewWithDefaults())
res, err := dg.FromFile(ctx, "conversation.txt", &interfaces.AnalyzeOptions{
  Sentiment: true,
  Topics:    true,
})
```

## Low-Level Transport

Source: `pkg/client/analyze/v1/client.go`

```go
func (c *Client) DoFile(ctx context.Context, filePath string, req *interfaces.AnalyzeOptions, resBody interface{}) error
func (c *Client) DoStream(ctx context.Context, src io.Reader, options *interfaces.AnalyzeOptions, resBody interface{}) error
func (c *Client) DoText(ctx context.Context, text string, options *interfaces.AnalyzeOptions, resBody interface{}) error
func (c *Client) DoURL(ctx context.Context, uri string, options *interfaces.AnalyzeOptions, resBody interface{}) error
func IsURL(str string) bool
```

Use `DoText` when your application already has the text in memory.

```go
var rawResponse map[string]any
err := analyze.NewWithDefaults().DoText(ctx, "Customer sounded upset but issue was resolved.", &interfaces.AnalyzeOptions{
  Sentiment: true,
  Summarize: true,
}, &rawResponse)
```

## Response Shape

The typed return is `*github.com/deepgram/deepgram-go-sdk/v3/pkg/api/analyze/v1/interfaces.AnalyzeResponse`. Important top-level fields are:

- `Metadata`
- `Results`
- `Results.Sentiments`
- `Results.Topics`
- `Results.Intents`
- `Results.Summary`

These definitions are in `pkg/api/analyze/v1/interfaces/types.go`.
