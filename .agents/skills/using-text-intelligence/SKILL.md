---
name: using-text-intelligence
description: Use when writing or reviewing Go code in this repo that sends text to Deepgram Read via the analyze client. Route speech/audio inputs to using-speech-to-text or using-audio-intelligence, and management/admin work to using-management-api.
---

# Using Deepgram Text Intelligence from the Go SDK

## When to use this product

Use this skill for plain-text analysis handled by `pkg/client/analyze`.

- summarization
- sentiment on text input
- URL, file, or stream based Read requests

Use a different skill when:

- your source material is audio and should go through `/v1/listen` (`using-audio-intelligence`)
- you need transcription first (`using-speech-to-text`)
- you need admin endpoints (`using-management-api`)

## Authentication

Set `DEEPGRAM_API_KEY` before creating the analyze client.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"strings"

	analyze "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/analyze"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func main() error {
	ctx := context.Background()

	client, err := analyze.NewWithDefaults()
	if err != nil {
		return err
	}

	resp, err := client.FromStream(
		ctx,
		strings.NewReader("Deepgram provides APIs for speech, voice, and media intelligence."),
		&interfaces.AnalyzeOptions{Summarize: true, Sentiment: true},
	)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}
```

## Key parameters

- `interfaces.AnalyzeOptions`
  - common fields include `Summarize`, `Sentiment`, and other Read features exposed by the analyze client
- request helpers
  - `FromFile`
  - `FromStream`
  - `FromURL`
- constructors
  - `analyze.NewWithDefaults()`
  - `analyze.New(apiKey, options)`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/analyze/client.go`
   - `pkg/client/analyze/v1/client.go`
   - `pkg/client/interfaces/v1/types-analyze.go`
   - `pkg/api/version/analyze-version.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/reference/text-intelligence/analyze-text`
   - `https://developers.deepgram.com/docs/text-intelligence`
   - `https://developers.deepgram.com/docs/text-sentiment-analysis`

## Gotchas

1. The Go package is named `analyze`, but the product route maps to Read / text intelligence.
2. Do not send raw audio here; audio intelligence features live on the listen REST path.
3. Reuse `context.Context` and stream readers for large inputs instead of materializing everything into one string.

## Example files in this repo

- `examples/analyze/summary/main.go`
- `examples/analyze/sentiment/main.go`
