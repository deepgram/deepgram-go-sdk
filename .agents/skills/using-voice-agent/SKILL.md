---
name: using-voice-agent
description: Use when writing or reviewing Go code in this repo that runs a Deepgram Voice Agent session over WebSockets, including runtime settings, prompt updates, speak updates, injected messages, and event handling. Route standalone STT/TTS work to using-speech-to-text or using-text-to-speech.
---

# Using Deepgram Voice Agent from the Go SDK

## When to use this product

Use this skill for live Voice Agent runtime flows in `pkg/client/agent`.

- open an agent WebSocket session
- send initial settings
- stream audio
- react to agent events and function-call messages
- update prompt/speak settings during the session

Use a different skill when:

- you only need STT (`using-speech-to-text`)
- you only need TTS (`using-text-to-speech`)
- you need management/admin endpoints (`using-management-api`)

## Authentication

Set `DEEPGRAM_API_KEY` before opening the agent connection.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

```go
package main

import (
	"context"
	"fmt"

	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
)

func main() error {
	ctx := context.Background()
	settings := agent.NewSettingsConfigurationOptions()

	conn, err := agent.NewWSUsingChanWithDefaults(ctx, settings)
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

	// Stream audio frames, watch agent events, and respond to function calls as needed.
	return nil
}
```

## Key parameters

- constructors
  - `agent.NewSettingsConfigurationOptions()`
  - `agent.NewWSUsingChanWithDefaults(...)`
  - `agent.NewWSUsingChan(...)`
- runtime methods
  - `Connect`, `Start`, `ProcessMessage`, `Stream`, `Write`, `KeepAlive`
  - reconnect helpers like `AttemptReconnect` and error handling helpers in the WS client
- message and event payloads in `pkg/api/agent/v1/websocket/interfaces/types.go`
  - `UpdatePrompt`
  - `UpdateSpeak`
  - `InjectAgentMessage`
  - `InjectUserMessage`
  - `FunctionCallResponse`
  - server events such as `WelcomeResponse`, `ConversationTextResponse`, `FunctionCallRequestResponse`, `AgentStartedSpeakingResponse`, `AgentAudioDoneResponse`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `pkg/client/agent/client.go`
   - `pkg/client/agent/v1/websocket/client_channel.go`
   - `pkg/client/agent/v1/websocket/new_using_chan.go`
   - `pkg/client/interfaces/v1/types-agent.go`
   - `pkg/api/agent/v1/websocket/interfaces/types.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/reference/voice-agent/voice-agent`
   - `https://developers.deepgram.com/docs/voice-agent`
   - `https://developers.deepgram.com/docs/configure-voice-agent`
   - `https://developers.deepgram.com/docs/voice-agent-message-flow`

## Gotchas

1. This repo exposes live Voice Agent runtime over WebSockets, not a persisted configuration-management surface.
2. Keep audio streaming, event handling, and any function-call response loop running concurrently.
3. Follow the example session setup instead of inventing your own event names; the repo already defines concrete message structs.
4. Use `defer conn.Close()` and keep keepalive/reconnect behavior consistent with the shared WS client.

## Example files in this repo

- `examples/agent/websocket/simple/main.go`
- `examples/agent/websocket/no_mic/main.go`
- `examples/agent/websocket/arbitrary_keys/main.go`
- `tests/unit_test/agent/agent_speak_test.go`
