---
name: deepgram-go-voice-agent
description: "Use when writing or reviewing Go code in this repo that runs a Deepgram Voice Agent session over WebSockets, including runtime settings, prompt updates, speak updates, injected messages, and event handling. Route standalone STT/TTS work to deepgram-go-speech-to-text or deepgram-go-text-to-speech."
---

# Using Deepgram Voice Agent from the Go SDK

Use this skill for live Voice Agent runtime flows in `pkg/client/agent`.

Use a different skill when:

- you only need STT (`deepgram-go-speech-to-text`)
- you only need TTS (`deepgram-go-text-to-speech`)
- you need management/admin endpoints (`deepgram-go-management-api`)

## Authentication

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Session lifecycle

Follow these steps in order:

1. **Configure** -- build settings with `agent.NewSettingsConfigurationOptions()`
2. **Connect** -- call `conn.Connect()` (returns `bool`; fail if `false`)
3. **Start** -- call `conn.Start()` to launch read/write pumps
4. **Stream audio** -- send PCM chunks via `conn.WriteBinary(chunk)`
5. **Handle events** -- read from the channel handler (see example below)
6. **Runtime updates** -- call `conn.ProcessMessage(msg)` with `UpdatePrompt`, `UpdateSpeak`, `InjectAgentMessage`, or `InjectUserMessage`
7. **Respond to function calls** -- send `FunctionCallResponse` when you receive `FunctionCallRequestResponse`
8. **Shut down** -- `defer conn.Stop()` near construction

## Quick start -- connection and event loop

```go
package main

import (
	"context"
	"fmt"
	"log"

	agentws "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket"
	agent "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	settings := agent.NewSettingsConfigurationOptions()
	handler := agentws.NewDefaultChanHandler()

	conn, err := agent.NewWSUsingChanWithDefaults(ctx, settings, handler)
	if err != nil {
		return err
	}
	defer conn.Stop()

	if ok := conn.Connect(); !ok {
		return fmt.Errorf("agent WebSocket connect failed")
	}
	conn.Start()

	// Event handling loop -- run concurrently alongside audio streaming.
	go func() {
		for {
			select {
			case welcome := <-handler.GetWelcome():
				fmt.Printf("Session started: %s\n", welcome.SessionID)
			case text := <-handler.GetConversationText():
				fmt.Printf("[%s] %s\n", text.Role, text.Content)
			case fcReq := <-handler.GetFunctionCallRequest():
				// Respond to function calls:
				resp := &agentws.FunctionCallResponse{
					FunctionCallID: fcReq.FunctionCallID,
					Output:         `{"result": "ok"}`,
				}
				if err := conn.ProcessMessage(ctx, resp); err != nil {
					log.Printf("function call response error: %v", err)
				}
			case <-handler.GetUnhandled():
				// Log or ignore unrecognized events
			case <-ctx.Done():
				return
			}
		}
	}()

	// Stream audio frames from your mic or file reader:
	//   for chunk := range audioSource { conn.WriteBinary(chunk) }

	// Runtime updates mid-session:
	//   conn.ProcessMessage(ctx, &agentws.UpdatePrompt{Prompt: "new system prompt"})
	//   conn.ProcessMessage(ctx, &agentws.InjectAgentMessage{Content: "Hello!"})

	return nil
}
```

## Key parameters

- constructors: `agent.NewSettingsConfigurationOptions()`, `agent.NewWSUsingChanWithDefaults(...)`, `agent.NewWSUsingChan(...)`
- runtime methods: `Connect`, `Start`, `ProcessMessage`, `Stream`, `WriteBinary`, `KeepAlive`, `AttemptReconnect`
- event/message types in `pkg/api/agent/v1/websocket/interfaces/types.go`:
  `WelcomeResponse`, `ConversationTextResponse`, `FunctionCallRequestResponse`, `AgentStartedSpeakingResponse`, `AgentAudioDoneResponse`, `UpdatePrompt`, `UpdateSpeak`, `InjectAgentMessage`, `InjectUserMessage`, `FunctionCallResponse`

## API reference (layered)

1. In-repo: `pkg/client/agent/client.go`, `pkg/client/agent/v1/websocket/client_channel.go`, `pkg/client/interfaces/v1/types-agent.go`, `pkg/api/agent/v1/websocket/interfaces/types.go`
2. OpenAPI: `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI: `https://developers.deepgram.com/asyncapi.yaml`
4. Product docs: `https://developers.deepgram.com/docs/voice-agent`, `https://developers.deepgram.com/docs/voice-agent-message-flow`

## Gotchas

1. `Connect()` returns `bool`, not `error` -- always check it.
2. Keep audio streaming and the event-handling goroutine running concurrently; blocking either stalls the session.
3. Channel-based constructors require an `AgentMessageChan` handler; events are routed there, not pulled from the client.
4. Use the concrete message structs defined in the repo -- do not invent custom event names.

## Example files

- `examples/agent/websocket/simple/main.go`
- `examples/agent/websocket/no_mic/main.go`
- `examples/agent/websocket/arbitrary_keys/main.go`
- `tests/unit_test/agent/agent_speak_test.go`
