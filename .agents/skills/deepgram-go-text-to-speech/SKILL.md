---
name: deepgram-go-text-to-speech
description: Use when writing or reviewing Go code in this repo that synthesizes audio with Speak v1 REST or Speak WebSockets. Route transcription work to deepgram-go-speech-to-text, voice conversation runtime work to deepgram-go-voice-agent, and repository maintenance work to deepgram-go-maintaining-sdk.
---

# Using Deepgram Text-to-Speech from the Go SDK

## When to use this product

Use this skill for `pkg/client/speak` work:

- file or stream synthesis over REST
- low-latency synthesis over WebSockets
- callback-based or channel-based audio playback pipelines

Use a different skill when:

- you need STT (`deepgram-go-speech-to-text`)
- you need live voice-agent orchestration (`deepgram-go-voice-agent`)
- you need repo workflow guidance (`deepgram-go-maintaining-sdk`)

## Authentication

Set `DEEPGRAM_API_KEY` before creating Speak clients.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

Use the repo's env-backed client defaults instead of embedding secrets in code.

## Quick start

REST synthesis to file:

```go
package main

import (
	"context"

	speak "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func main() error {
	ctx := context.Background()

	client, err := speak.NewRESTWithDefaults()
	if err != nil {
		return err
	}

	if err := client.ToFile(
		ctx,
		"hello.wav",
		"Hello from the Deepgram Go SDK.",
		&interfaces.SpeakOptions{Model: "aura-2-thalia-en"},
	); err != nil {
		return err
	}

	return nil
}
```

Streaming synthesis with callbacks or channels:

```go
package main

import (
	"context"
	"fmt"

	speak "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func main() error {
	ctx := context.Background()

	conn, err := speak.NewWSUsingChanWithDefaults(
		ctx,
		&interfaces.WSSpeakOptions{Model: "aura-2-thalia-en"},
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

	if err := conn.SpeakWithText("Streaming TTS from Go."); err != nil {
		return err
	}

	if err := conn.Flush(); err != nil {
		return err
	}

	return nil
}
```

## Key parameters

- `interfaces.SpeakOptions`
  - typical fields: `Model`, `Encoding`, `Container`, `SampleRate`
- `interfaces.WSSpeakOptions`
  - typical fields: `Model`, streaming audio format settings
- REST methods
  - `ToStream`, `ToFile`, `ToSave`
- WS methods
  - `SpeakWithText`, `Speak`, `Flush`, `Reset`, `KeepAlive`
- constructors
  - `speak.NewRESTWithDefaults()` / `speak.NewREST(...)`
  - `speak.NewWSUsingCallback...`
  - `speak.NewWSUsingChan...`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/speak/client.go`
   - `pkg/client/speak/v1/rest/client.go`
   - `pkg/client/speak/v1/websocket/client_callback.go`
   - `pkg/client/speak/v1/websocket/client_channel.go`
   - `pkg/client/interfaces/v1/types-speak.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/reference/text-to-speech/speak-request`
   - `https://developers.deepgram.com/reference/text-to-speech/speak-streaming`
   - `https://developers.deepgram.com/docs/tts-models`

## Gotchas

1. REST and WebSocket Speak clients use different option structs.
2. WebSocket flows usually need `Connect()`, `Start()`, one or more `Speak...` calls, then `Flush()` or `Reset()`.
3. Keep the playback or message-consumer goroutine running while audio frames arrive.
4. Follow the examples for file handling and output format selection instead of guessing encodings.

## Example files in this repo

- `examples/text-to-speech/rest/file/hello-world/main.go`
- `examples/text-to-speech/websocket/simple_callback/main.go`

## Central product skills

For cross-language Deepgram product knowledge — the consolidated API reference, documentation finder, focused runnable recipes, third-party integration examples, and MCP setup — install the central skills:

```bash
npx skills add deepgram/skills
```

This SDK ships language-idiomatic code skills; `deepgram/skills` ships cross-language product knowledge (see `api`, `docs`, `recipes`, `examples`, `starters`, `setup-mcp`).
