---
name: deepgram-go-conversational-stt
description: Use when planning or reviewing Go SDK work for Deepgram conversational STT / Flux v2. This repo ships first-class callback and channel Listen v2 WebSocket clients with typed options and mid-session configuration. Route v1 transcription to deepgram-go-speech-to-text.
---

# Using Deepgram Conversational STT from the Go SDK

## When to use this product

Use this skill when someone asks for Flux or conversational STT v2 in the Go SDK.

Current repo status: this SDK ships a first-class `pkg/client/listen/v2` WebSocket client for Flux STT. Flux STT does not have a REST client.

Use a different skill when:

- v1 Listen is acceptable (`deepgram-go-speech-to-text`)
- voice-agent runtime is the real target (`deepgram-go-voice-agent`)

## Authentication

Flux constructors with defaults resolve credentials from `DEEPGRAM_API_KEY` or `DEEPGRAM_ACCESS_TOKEN`.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

Use the versioned v2 package and typed options:

```go
package main

import (
  "context"
  "log"

  api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
  interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
  listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2"
)

func main() {
  ctx := context.Background()
  handler := api.NewDefaultCallbackHandler()
  client, err := listen.NewWSUsingCallbackWithDefaults(ctx, &interfaces.FluxTranscriptionOptions{
    Model: "flux-general-en",
  }, handler)
  if err != nil {
    log.Fatal(err)
  }
  defer client.Stop()

  if !client.Connect() {
    log.Fatal("connect failed")
  }
}
```

## Key parameters

- `interfaces.FluxTranscriptionOptions`
  - `Model` is required: `flux-general-en` or `flux-general-multi`
  - `Encoding`, `SampleRate`, `Keyterm`, `Tag`, and `LanguageHint` configure connection behavior
  - `EagerEotThreshold`, `EotThreshold`, and `EotTimeoutMs` configure turn detection
- `interfaces.FluxConfigureOptions`
  - send with `WSCallback.Configure` or `WSChannel.Configure` to update thresholds, keyterms, or language hints mid-session
- constructors
  - `listen.NewWSUsingCallbackWithDefaults` and `listen.NewWSUsingChanWithDefaults` use environment-backed credentials
  - `listen.NewWSUsingCallback` and `listen.NewWSUsingChan` accept explicit v2 client options and handlers

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/listen/v2/client.go`
   - `pkg/client/listen/v2/websocket`
   - `pkg/api/listen/v2/websocket/interfaces`
   - `pkg/client/interfaces/v2/types-flux.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/docs/conversational-speech-recognition`

## Gotchas

1. Flux STT is WebSocket-only; route v1 prerecorded or Nova streaming requests to `deepgram-go-speech-to-text`.
2. Use `pkg/client/listen/v2`, not the v1 `pkg/client/listen` facade, for Flux STT.
3. Flux callback and channel clients use their own v2 event interfaces; do not mix them with Listen v1 handlers.

## Example files in this repo

- `examples/speech-to-text/websocket/flux_callback/main.go`
- `examples/speech-to-text/websocket/flux_channel/main.go`
- `examples/speech-to-text/websocket/flux_force_end_turn/main.go`

## Central product skills

For cross-language Deepgram product knowledge — the consolidated API reference, documentation finder, focused runnable recipes, third-party integration examples, and MCP setup — install the central skills:

```bash
npx skills add deepgram/skills
```

This SDK ships language-idiomatic code skills; `deepgram/skills` ships cross-language product knowledge (see `api`, `docs`, `recipes`, `examples`, `starters`, `setup-mcp`).
