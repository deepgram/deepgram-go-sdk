---
name: deepgram-go-conversational-stt
description: "Use when planning, prototyping, or reviewing Go SDK work for Deepgram conversational speech-to-text / Flux v2 real-time streaming transcription. Guides raw WebSocket integration against the shared client plumbing, configures host and version overrides, and routes supported v1 listen transcription to deepgram-go-speech-to-text."
---

# Using Deepgram Conversational STT from the Go SDK

Use when someone asks for Flux, conversational STT v2, or real-time streaming transcription v2 in the Go SDK.

**Current status:** this SDK does **not** ship a first-class `listen/v2` client. The guidance below covers the raw WebSocket fallback.

Use a different skill when:

- v1 Listen REST or WebSocket is acceptable (`deepgram-go-speech-to-text`)
- voice-agent runtime is the real target (`deepgram-go-voice-agent`)

## Authentication

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Raw v2 WebSocket integration

No `pkg/client/listen/v2` constructor exists. To prototype v2, use the shared WebSocket plumbing with host/version overrides:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// 1. Build client options with v2 path override
	opts := interfaces.ClientOptions{
		APIKey: os.Getenv("DEEPGRAM_API_KEY"),
		Host:   "api.deepgram.com",
		Path:   "/v2/listen", // override default v1 path
	}

	// 2. Create a raw WS connection using shared plumbing
	ws := common.New(ctx, "", &opts)

	// 3. Connect and validate
	if ok := ws.Connect(); !ok {
		return fmt.Errorf("v2 WebSocket connect failed — verify API key and endpoint")
	}
	defer ws.Stop()

	// 4. Start read/write pumps
	ws.Start()

	// 5. Stream PCM audio frames
	//    for chunk := range audioSource { ws.WriteBinary(chunk) }

	// 6. Finalize and shut down
	//    ws.Finalize()
	return nil
}
```

**Treat this as a manual integration**, not an SDK-supported surface. Do not invent `listen/v2` package paths or constructors.

## Key files for raw integration

| Purpose | File |
|---------|------|
| Shared WS lifecycle and reconnect | `pkg/client/common/v1/websocket.go` |
| Client option types (host, path overrides) | `pkg/client/interfaces/v1/types-client.go` |
| Version and path constants | `pkg/api/version/live-version.go`, `pkg/api/version/constants.go` |

## API reference (layered)

1. In-repo: `pkg/client/common/v1/websocket.go`, `pkg/client/interfaces/v1/types-client.go`, `pkg/api/version/live-version.go`
2. OpenAPI: `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI: `https://developers.deepgram.com/asyncapi.yaml`
4. Product docs: `https://developers.deepgram.com/docs/conversational-speech-recognition`

## Gotchas

1. Conversational STT v2 is **not** a first-class Go SDK client here. Do not invent `listen/v2` constructors.
2. If production-ready Flux support is required, point to raw API usage or another SDK.
3. For WS reconnect and keepalive patterns, see `tests/edge_cases/keepalive/main.go` and `tests/edge_cases/reconnect_client/main.go`.
