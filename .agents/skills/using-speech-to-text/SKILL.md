---
name: using-speech-to-text
description: Use when writing or reviewing Go code in this repo that transcribes prerecorded audio with Listen v1 REST or streams live audio with Listen v1 WebSockets. Route text generation to using-text-to-speech, text analysis to using-text-intelligence, audio analytics overlays to using-audio-intelligence, and Flux or other v2 conversational work to using-conversational-stt.
---

# Using Deepgram Speech-to-Text from the Go SDK

## When to use this product

Use this skill for `pkg/client/listen` work:

- prerecorded transcription with `FromURL`, `FromFile`, or `FromStream`
- live transcription with `pkg/client/listen/v1/websocket`
- channel-based or callback-based streaming flows

Use a different skill when:

- you need TTS output (`using-text-to-speech`)
- you need text analysis on plain text (`using-text-intelligence`)
- you need analytics overlays like summaries, topics, or sentiments (`using-audio-intelligence`)
- you need Flux / conversational STT v2 (`using-conversational-stt`)

## Authentication

Set `DEEPGRAM_API_KEY` before constructing clients.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

This SDK reads env-backed defaults via the client option layer. Prefer API key or token auth supported by the repo's client options; do not hardcode credentials.

## Quick start

Prerecorded REST:

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
			Model:       "nova-3",
			SmartFormat: true,
			Punctuate:   true,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(resp.Results.Channels[0].Alternatives[0].Transcript)
	return nil
}
```

Live WebSocket with channel fan-out:

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

	conn, err := listen.NewWSUsingChanWithDefaults(
		ctx,
		&interfaces.LiveTranscriptionOptions{Model: "nova-3", InterimResults: true},
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Connect(); err != nil {
		return err
	}

	go func() {
		for msg := range conn.MessageChan() {
			fmt.Println(msg)
		}
	}()

	if err := conn.Start(); err != nil {
		return err
	}

	// Write PCM/audio chunks from your mic or file reader, then finalize.
	if err := conn.Finalize(); err != nil {
		return err
	}

	return nil
}
```

## Key parameters

- `interfaces.PreRecordedTranscriptionOptions`
  - common fields: `Model`, `Language`, `Punctuate`, `SmartFormat`, `Diarize`, `Redact`, `Utterances`
  - transport methods: `FromURL`, `FromFile`, `FromStream`
- `interfaces.LiveTranscriptionOptions`
  - common fields: `Model`, `Language`, `Encoding`, `SampleRate`, `Channels`, `InterimResults`, `Endpointing`
- constructor families
  - REST: `listen.NewRESTWithDefaults()`, `listen.NewREST(apiKey, options)`
  - WS callbacks: `listen.NewWSUsingCallback...`
  - WS channels: `listen.NewWSUsingChan...`
- lifecycle
  - connect, start, stream/write audio, keepalive as needed, `Finalize()`, then `defer conn.Close()`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/listen/client.go`
   - `pkg/client/listen/v1/rest/client.go`
   - `pkg/client/listen/v1/websocket/client_callback.go`
   - `pkg/client/listen/v1/websocket/client_channel.go`
   - `pkg/client/interfaces/v1/types-prerecorded.go`
   - `pkg/client/interfaces/v1/types-stream.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/reference/speech-to-text/listen-pre-recorded`
   - `https://developers.deepgram.com/reference/speech-to-text/listen-streaming`
   - `https://developers.deepgram.com/docs/speech-to-text`

## Gotchas

1. This repo uses `listen` package names for STT v1, not `transcription`.
2. Streaming code is split into callback and channel variants; copy the style that matches the surrounding package.
3. For WebSockets, keep `defer conn.Close()` near construction and finalize before shutdown.
4. Live and prerecorded option structs are different; do not assume analytics-only prerecorded fields exist in live mode.
5. Use `context.Context` and return `error`; do not translate examples into exception-style control flow.

## Example files in this repo

- `examples/speech-to-text/rest/url/main.go`
- `examples/speech-to-text/rest/file/main.go`
- `examples/speech-to-text/websocket/microphone_channel/main.go`
- `examples/speech-to-text/websocket/microphone_callback/main.go`
- `tests/edge_cases/keepalive/main.go`
- `tests/edge_cases/reconnect_client/main.go`
