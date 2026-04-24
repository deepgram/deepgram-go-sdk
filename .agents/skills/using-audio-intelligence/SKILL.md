---
name: using-audio-intelligence
description: Use when writing or reviewing Go code in this repo that applies summaries, topics, intents, sentiment, language detection, diarization, redaction, or entity extraction to audio inputs through Listen v1 REST. Route plain transcription to using-speech-to-text and plain-text Read requests to using-text-intelligence.
---

# Using Deepgram Audio Intelligence from the Go SDK

## When to use this product

Use this skill for `/v1/listen` REST requests that combine transcription with analytics overlays.

- summaries
- topics
- intents
- sentiments
- entity extraction
- language detection
- diarization and redaction where supported

Use a different skill when:

- you only need transcription (`using-speech-to-text`)
- your input is already text (`using-text-intelligence`)

## Authentication

Set `DEEPGRAM_API_KEY` before creating the listen REST client.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

```go
package main

import (
	"context"
	"fmt"

	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func main() error {
	ctx := context.Background()

	client, err := listen.NewRESTWithDefaults()
	if err != nil {
		return err
	}

	resp, err := client.FromURL(
		ctx,
		"https://dpgr.am/spacewalk.wav",
		&interfaces.PreRecordedTranscriptionOptions{
			Model:     "nova-3",
			Summarize: true,
			Topics:    true,
			Sentiment: true,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(resp.Results.Summary)
	return nil
}
```

## Key parameters

- analytics live on `interfaces.PreRecordedTranscriptionOptions`
  - `Summarize`
  - `Topics`
  - `Intents`
  - `Sentiment`
  - `DetectLanguage`
  - `DetectEntities`
  - `Diarize`
  - `Redact`
- response payloads are in `pkg/api/listen/v1/rest/interfaces/types.go`
  - `Sentiments`
  - `Topics`
  - `Intents`
  - `Entities`
  - `SummaryV2`
- this SDK is REST-first for audio intelligence; the live WS option struct does not expose the same analytics surface

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/listen/v1/rest/client.go`
   - `pkg/client/interfaces/v1/types-prerecorded.go`
   - `pkg/client/interfaces/v1/types-stream.go`
   - `pkg/api/listen/v1/rest/interfaces/types.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/docs/stt-intelligence-feature-overview`
   - `https://developers.deepgram.com/docs/summarization`
   - `https://developers.deepgram.com/docs/topic-detection`
   - `https://developers.deepgram.com/docs/intent-recognition`
   - `https://developers.deepgram.com/docs/sentiment-analysis`
   - `https://developers.deepgram.com/docs/language-detection`
   - `https://developers.deepgram.com/docs/redaction`
   - `https://developers.deepgram.com/docs/diarization`

## Gotchas

1. In this SDK snapshot, audio intelligence is effectively a prerecorded REST feature set.
2. `LiveTranscriptionOptions` does not expose the same summarize/topic/intent/sentiment/entity surface.
3. Start from the speech-to-text REST examples, then layer analytics flags onto `PreRecordedTranscriptionOptions`.

## Example files in this repo

- `examples/speech-to-text/rest/summary/main.go`
- `examples/speech-to-text/rest/sentiment/main.go`
- `examples/speech-to-text/rest/topic/main.go`
- `examples/speech-to-text/rest/intent/main.go`
