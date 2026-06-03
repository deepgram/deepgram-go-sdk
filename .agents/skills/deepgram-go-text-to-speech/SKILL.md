---
name: deepgram-go-text-to-speech
description: "Use when writing or reviewing Go code in this repo that synthesizes audio with Speak v1 REST or Speak WebSockets -- configure output formats, save to file, stream low-latency audio, and manage WebSocket playback pipelines. Route transcription to deepgram-go-speech-to-text and voice conversation runtime to deepgram-go-voice-agent."
---

# Using Deepgram Text-to-Speech from the Go SDK

Use this skill for `pkg/client/speak` work: file/stream synthesis over REST, low-latency synthesis over WebSockets, callback or channel audio playback.

Use a different skill when:

- you need STT (`deepgram-go-speech-to-text`)
- you need live voice-agent orchestration (`deepgram-go-voice-agent`)

## Authentication

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start -- REST synthesis to file

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v1/rest"
	speak "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	client := speak.NewRESTWithDefaults()
	dg := api.New(client)

	if _, err := dg.ToSave(
		ctx,
		"hello.wav",
		"Hello from the Deepgram Go SDK.",
		&interfaces.SpeakOptions{Model: "aura-2-thalia-en"},
	); err != nil {
		return err
	}

	// Validate output was written
	info, err := os.Stat("hello.wav")
	if err != nil || info.Size() == 0 {
		return fmt.Errorf("output file missing or empty")
	}
	return nil
}
```

## Streaming synthesis over WebSocket

```go
// After imports: speakws, speak, interfaces (same pattern as REST)
func runWS() error {
	ctx := context.Background()
	handler := speakws.NewDefaultChanHandler()

	conn, err := speak.NewWSUsingChanWithDefaults(
		ctx,
		&interfaces.WSSpeakOptions{Model: "aura-2-thalia-en"},
		handler,
	)
	if err != nil {
		return err
	}
	defer conn.Stop()

	if ok := conn.Connect(); !ok {
		return fmt.Errorf("TTS WebSocket connect failed")
	}
	conn.Start()

	if err := conn.SpeakWithText("Streaming TTS from Go."); err != nil {
		return err
	}
	// Flush waits for the server to confirm all audio has been sent
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}
```

## Key parameters

| Layer | Types / Methods |
|-------|----------------|
| REST options | `interfaces.SpeakOptions` -- `Model`, `Encoding`, `Container`, `SampleRate` |
| WS options | `interfaces.WSSpeakOptions` -- `Model`, streaming format settings |
| REST methods | `api.New(client).ToStream`, `ToFile`, `ToSave` |
| WS methods | `SpeakWithText`, `Speak`, `Flush`, `Reset` |
| REST constructors | `speak.NewRESTWithDefaults()`, `speak.NewREST(...)` |
| WS constructors | `speak.NewWSUsingCallback...`, `speak.NewWSUsingChan...` |

## API reference (layered)

1. In-repo: `pkg/client/speak/client.go`, `pkg/client/speak/v1/rest/client.go`, `pkg/client/speak/v1/websocket/client_channel.go`, `pkg/client/interfaces/v1/types-speak.go`
2. OpenAPI: `https://developers.deepgram.com/openapi.yaml`
3. Product docs: `https://developers.deepgram.com/reference/text-to-speech/speak-request`, `https://developers.deepgram.com/docs/tts-models`

## Gotchas

1. REST and WebSocket Speak clients use **different option structs** (`SpeakOptions` vs `WSSpeakOptions`).
2. `Connect()` returns `bool`, not `error` -- always check it.
3. Keep the playback or message-consumer goroutine running while audio frames arrive on the WS handler.
4. Follow the in-repo examples for output format selection instead of guessing encodings.

## Example files

- `examples/text-to-speech/rest/file/hello-world/main.go`
- `examples/text-to-speech/websocket/simple_channel/main.go`
- `examples/text-to-speech/websocket/simple_callback/main.go`
