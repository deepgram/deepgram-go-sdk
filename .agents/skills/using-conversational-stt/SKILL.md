---
name: using-conversational-stt
description: Use when planning or reviewing Go SDK work for Deepgram conversational STT / Flux v2. This repo does not currently ship a first-class v2 listen client, so route supported v1 transcription to using-speech-to-text and document raw WebSocket fallback honestly when v2 is requested.
---

# Using Deepgram Conversational STT from the Go SDK

## When to use this product

Use this skill when someone asks for Flux or conversational STT v2 in the Go SDK.

Current repo status: this SDK snapshot does **not** implement a first-class `listen/v2` client.

Use a different skill when:

- v1 Listen is acceptable (`using-speech-to-text`)
- voice-agent runtime is the real target (`using-voice-agent`)

## Authentication

If you hand-roll a raw WebSocket client, it still uses the same Deepgram credentials.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

There is no supported `pkg/client/listen/v2` constructor in this repo today.

What you can do instead:

- use `using-speech-to-text` for supported `listen` v1 flows
- if you must prototype v2, start from the low-level shared WebSocket plumbing and path/version hooks already in the repo:
  - `pkg/client/common/v1/websocket.go`
  - `pkg/client/interfaces/v1/types-client.go`
  - `pkg/api/version/live-version.go`
  - `pkg/api/version/constants.go`

Treat that path as a manual integration, not an SDK-supported product surface.

## Key parameters

- unsupported in current public package layout
  - no `pkg/client/listen/v2`
  - no `pkg/api/listen/v2`
- possible raw integration hooks
  - client host/version/path overrides from the shared client options
  - shared WS lifecycle and reconnect logic in `pkg/client/common/v1/websocket.go`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/common/v1/websocket.go`
   - `pkg/client/interfaces/v1/types-client.go`
   - `pkg/api/version/live-version.go`
   - `pkg/api/version/constants.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/docs/conversational-speech-recognition`

## Gotchas

1. Be explicit: conversational STT v2 is not yet implemented as a first-class Go SDK client here.
2. Do not invent `listen/v2` package paths, constructors, or examples.
3. If a task requires production-ready Flux support, point to raw API usage or another SDK until this repo adds first-class coverage.

## Example files in this repo

- No dedicated conversational STT v2 examples exist in this repo.
- For shared WS patterns, inspect `tests/edge_cases/keepalive/main.go` and `tests/edge_cases/reconnect_client/main.go`.
