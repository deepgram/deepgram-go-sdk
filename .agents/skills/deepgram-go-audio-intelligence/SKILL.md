---
name: deepgram-go-audio-intelligence
description: "Use when writing or reviewing Go code in this repo that applies summaries, topics, intents, sentiment, language detection, diarization, redaction, or entity extraction to audio inputs through Listen v1 REST. Route plain transcription to deepgram-go-speech-to-text and plain-text Read requests to deepgram-go-text-intelligence."
---

# Using Deepgram Audio Intelligence from the Go SDK

Use this skill for `/v1/listen` REST requests that combine transcription with analytics overlays: summaries, topics, intents, sentiments, entity extraction, language detection, diarization, and redaction.

Use a different skill when:

- you only need transcription (`deepgram-go-speech-to-text`)
- your input is already text (`deepgram-go-text-intelligence`)

## Authentication

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	client := listen.NewRESTWithDefaults()
	dg := api.New(client)

	resp, err := dg.FromURL(
		ctx,
		"https://dpgr.am/spacewalk.wav",
		&interfaces.PreRecordedTranscriptionOptions{
			Model:     "nova-3",
			Summarize: "v2",
			Topics:    true,
			Sentiment: true,
		},
	)
	if err != nil {
		return err
	}

	if resp.Results.Summary != nil {
		fmt.Println(resp.Results.Summary.Result)
	}
	return nil
}
```

## Key parameters

Analytics flags on `interfaces.PreRecordedTranscriptionOptions`:

| Flag | Type | Notes |
|------|------|-------|
| `Summarize` | `string` | e.g. `"v2"` |
| `Topics` | `bool` | |
| `Intents` | `bool` | |
| `Sentiment` | `bool` | |
| `DetectLanguage` | `bool` | |
| `DetectEntities` | `bool` | |
| `Diarize` | `bool` | |
| `Redact` | `[]string` | |

Response payloads in `pkg/api/listen/v1/rest/interfaces/types.go`: `Sentiments`, `Topics`, `Intents`, `Entities`, `SummaryV2`.

## API reference (layered)

1. In-repo: `pkg/client/listen/v1/rest/client.go`, `pkg/client/interfaces/v1/types-prerecorded.go`, `pkg/api/listen/v1/rest/interfaces/types.go`
2. OpenAPI: `https://developers.deepgram.com/openapi.yaml`
3. Product docs: `https://developers.deepgram.com/docs/stt-intelligence-feature-overview`, `https://developers.deepgram.com/docs/summarization`, `https://developers.deepgram.com/docs/sentiment-analysis`

## Gotchas

1. Audio intelligence is a **prerecorded REST** feature set only; `LiveTranscriptionOptions` does not expose the same analytics surface.
2. Build via `api.New(listen.NewRESTWithDefaults())`, then layer analytics flags onto `PreRecordedTranscriptionOptions`.
3. Do not send raw text here; text analysis belongs in `deepgram-go-text-intelligence`.

## Example files

- `examples/speech-to-text/rest/summary/main.go`
- `examples/speech-to-text/rest/sentiment/main.go`
- `examples/speech-to-text/rest/topic/main.go`
- `examples/speech-to-text/rest/intent/main.go`
