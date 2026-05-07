---
name: deepgram-go-speech-to-text
description: "Use when writing or reviewing Go code in this repo that transcribes prerecorded audio with Listen v1 REST or streams live audio with Listen v1 WebSockets. Route text generation to deepgram-go-text-to-speech, text analysis to deepgram-go-text-intelligence, audio analytics overlays to deepgram-go-audio-intelligence, and Flux v2 conversational work to deepgram-go-conversational-stt."
---

# Using Deepgram Speech-to-Text from the Go SDK

Use this skill for `pkg/client/listen` work: prerecorded transcription (`FromURL`, `FromFile`, `FromStream`), live transcription via WebSocket, and channel or callback streaming flows.

Use a different skill when:

- you need TTS output (`deepgram-go-text-to-speech`)
- you need text analysis on plain text (`deepgram-go-text-intelligence`)
- you need analytics overlays like summaries, topics, or sentiments (`deepgram-go-audio-intelligence`)
- you need Flux / conversational STT v2 (`deepgram-go-conversational-stt`)

## Authentication

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start -- prerecorded REST

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

## Live WebSocket -- lifecycle

Follow these steps in order:

1. **Construct** -- `listen.NewWSUsingChanWithDefaults(ctx, opts, handler)` (returns `conn, err`)
2. **Connect** -- `conn.Connect()` returns `bool`; fail if `false`
3. **Start** -- `conn.Start()` launches read/write pumps
4. **Stream** -- send PCM chunks: `conn.WriteBinary(chunk)` (check returned `error`)
5. **KeepAlive** -- call `conn.KeepAlive()` during idle periods to prevent timeout
6. **Finalize** -- `conn.Finalize()` flushes trailing audio; check `error`
7. **Stop** -- `defer conn.Stop()` near construction

```go
// After imports: listenws, listen, interfaces (same pattern as REST)
func runWS() error {
	ctx := context.Background()
	handler := listenws.NewDefaultChanHandler()

	conn, err := listen.NewWSUsingChanWithDefaults(
		ctx,
		&interfaces.LiveTranscriptionOptions{Model: "nova-3", InterimResults: true},
		handler,
	)
	if err != nil {
		return err
	}
	defer conn.Stop()

	if ok := conn.Connect(); !ok {
		return fmt.Errorf("STT WebSocket connect failed")
	}
	conn.Start()

	// Stream audio and finalize when done:
	//   for chunk := range audioSource { conn.WriteBinary(chunk) }
	//   conn.Finalize()
	return nil
}
```

## Key parameters

| Layer | Types / Fields |
|-------|---------------|
| Prerecorded options | `interfaces.PreRecordedTranscriptionOptions` -- `Model`, `Language`, `Punctuate`, `SmartFormat`, `Diarize`, `Redact`, `Utterances` |
| Live options | `interfaces.LiveTranscriptionOptions` -- `Model`, `Language`, `Encoding`, `SampleRate`, `Channels`, `InterimResults`, `Endpointing` |
| REST constructors | `listen.NewRESTWithDefaults()`, `listen.NewREST(apiKey, options)` |
| WS constructors | `listen.NewWSUsingCallback...`, `listen.NewWSUsingChan...` |
| REST methods | `api.New(client).FromURL`, `FromFile`, `FromStream` |

## API reference (layered)

1. In-repo: `pkg/client/listen/client.go`, `pkg/client/listen/v1/rest/client.go`, `pkg/client/listen/v1/websocket/client_channel.go`, `pkg/client/interfaces/v1/types-prerecorded.go`, `pkg/client/interfaces/v1/types-stream.go`
2. OpenAPI: `https://developers.deepgram.com/openapi.yaml`
3. Product docs: `https://developers.deepgram.com/reference/speech-to-text/listen-pre-recorded`, `https://developers.deepgram.com/reference/speech-to-text/listen-streaming`

## Gotchas

1. This repo uses `listen` package names for STT v1, not `transcription`.
2. Streaming code is split into callback and channel variants; match the surrounding package style.
3. Live and prerecorded option structs are different; analytics-only fields (e.g. `Summarize`) exist only on prerecorded.
4. `Connect()` returns `bool`, not `error` -- always check it.
5. Always call `Finalize()` before `Stop()` to flush trailing audio.

## Example files

- `examples/speech-to-text/rest/url/main.go`
- `examples/speech-to-text/rest/file/main.go`
- `examples/speech-to-text/websocket/microphone_channel/main.go`
- `examples/speech-to-text/websocket/microphone_callback/main.go`
- `tests/edge_cases/keepalive/main.go`
