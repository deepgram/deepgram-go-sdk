---
name: deepgram-go-conversational-stt
description: "Use when writing or reviewing Go code in this repo that streams conversational, turn-based audio through the Deepgram Flux (v2/listen) WebSocket API. Covers the first-class `pkg/client/listen/v2` client (callback and channel variants), `FluxTranscriptionOptions`, turn lifecycle events, and mid-session reconfiguration. Route classic v1 listen transcription to deepgram-go-speech-to-text, text generation to deepgram-go-text-to-speech, and voice-agent runtime work to deepgram-go-voice-agent."
---

# Using Deepgram Conversational STT (Flux) from the Go SDK

Use this skill for `pkg/client/listen/v2` work: real-time conversational transcription against `wss://api.deepgram.com/v2/listen` using Deepgram's Flux turn-based audio API. Available since SDK v3.6.0.

Use a different skill when:

- v1 Listen REST or WebSocket is the target (`deepgram-go-speech-to-text`)
- voice-agent runtime is the real target (`deepgram-go-voice-agent`)
- TTS output is needed (`deepgram-go-text-to-speech`)

## Authentication

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start -- callback client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen/v2"
)

type MyHandler struct{}

func (MyHandler) Open(*api.OpenResponse) error              { return nil }
func (MyHandler) Connected(cr *api.ConnectedResponse) error { fmt.Println("ready:", cr.RequestID); return nil }
func (MyHandler) TurnInfo(tr *api.TurnInfoResponse) error {
	if tr.EventType == api.TurnEventEndOfTurn {
		fmt.Printf("[Turn %d] FINAL: %s\n", tr.TurnIndex, tr.Transcript)
	}
	return nil
}
func (MyHandler) ConfigureSuccess(*api.ConfigureSuccessResponse) error { return nil }
func (MyHandler) ConfigureFailure(*api.ConfigureFailureResponse) error { return nil }
func (MyHandler) FatalError(fe *api.FatalErrorResponse) error          { return fmt.Errorf("fatal: %s", fe.Description) }
func (MyHandler) Close(*api.CloseResponse) error                       { return nil }
func (MyHandler) Error(er *api.ErrorResponse) error                    { return fmt.Errorf("error: %s", er.ErrMsg) }
func (MyHandler) UnhandledEvent([]byte) error                          { return nil }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	cOptions := &interfaces.ClientOptionsV2{EnableKeepAlive: true}
	tOptions := &interfaces.FluxTranscriptionOptions{
		Model:             "flux-general-en",
		Encoding:          "linear16",
		SampleRate:        16000,
		EagerEotThreshold: 0.3,
	}

	dg, err := client.NewWSUsingCallback(ctx, os.Getenv("DEEPGRAM_API_KEY"), cOptions, tOptions, MyHandler{})
	if err != nil {
		return err
	}
	defer dg.Stop()

	if ok := dg.Connect(); !ok {
		return fmt.Errorf("Flux WebSocket connect failed")
	}

	// Stream PCM audio frames:
	//   for chunk := range audioSource { dg.WriteBinary(chunk) }
	return nil
}
```

## WebSocket lifecycle

Follow these steps in order:

1. **Construct** -- `client.NewWSUsingCallback(ctx, apiKey, cOptions, tOptions, callback)` or `client.NewWSUsingChan(ctx, apiKey, cOptions, tOptions, chans)` (returns `conn, err`)
2. **Connect** -- `conn.Connect()` returns `bool`; fail if `false`
3. **Stream** -- send PCM chunks: `conn.WriteBinary(chunk)` (check returned `error`)
4. **KeepAlive** -- set `ClientOptionsV2{EnableKeepAlive: true}` to auto-send keepalives during idle periods
5. **Reconfigure** -- call `conn.Configure(&interfaces.FluxConfigureOptions{...})` to adjust thresholds mid-session without reconnecting
6. **Stop** -- `defer conn.Stop()` near construction

## Turn lifecycle events

`TurnInfo` fires for every turn-detection event. Inspect `TurnInfoResponse.EventType`:

| EventType | Meaning |
|-----------|---------|
| `TurnEventStartOfTurn` | New speech turn detected |
| `TurnEventUpdate` | Interim transcript update |
| `TurnEventEagerEndOfTurn` | Early end-of-turn signal (fires when `EagerEotThreshold` is set) |
| `TurnEventTurnResumed` | Speaker continued after an eager EOT -- treat the prior eager final as not final |
| `TurnEventEndOfTurn` | Final transcript for the turn |

## Channel client (alternative)

When goroutine-per-event ergonomics are preferred over a callback struct, use `client.NewWSUsingChan(...)` with an implementation of `api.FluxMessageChan`. Each `GetX()` getter returns `[]*chan *XResponse`; return `nil` for events you don't care about. See `examples/speech-to-text/websocket/flux_channel/main.go` for the full pattern.

## Key parameters

| Layer | Types / Fields |
|-------|---------------|
| Transcription options | `interfaces.FluxTranscriptionOptions` -- `Model` (`flux-general-en`, `flux-general-multi`), `Encoding`, `SampleRate`, `EagerEotThreshold` (0.3-0.9), `EotThreshold` (0.5-0.9, default 0.7), `EotTimeoutMs` (500-10000, default 5000), `Keyterm`, `LanguageHint`, `Tag`, `MipOptOut` |
| Client options | `interfaces.ClientOptionsV2` -- `APIKey`, `EnableKeepAlive`, host/path overrides |
| Mid-session config | `interfaces.FluxConfigureOptions` -- `Thresholds` (`*FluxThresholds`), `Keyterms`, `LanguageHints` |
| Callback constructors | `client.NewWSUsingCallback`, `NewWSUsingCallbackWithDefaults`, `NewWSUsingCallbackForDemo`, `NewWSUsingCallbackWithCancel` |
| Channel constructors | `client.NewWSUsingChan`, `NewWSUsingChanWithDefaults`, `NewWSUsingChanForDemo`, `NewWSUsingChanWithCancel` |
| Callback contract | `api.FluxMessageCallback` -- `Open`, `Connected`, `TurnInfo`, `ConfigureSuccess`, `ConfigureFailure`, `FatalError`, `Close`, `Error`, `UnhandledEvent` |
| Channel contract | `api.FluxMessageChan` -- `GetOpen`, `GetConnected`, `GetTurnInfo`, `GetConfigureSuccess`, `GetConfigureFailure`, `GetFatalError`, `GetClose`, `GetError`, `GetUnhandled` |

## API reference (layered)

1. In-repo: `pkg/client/listen/v2/client.go`, `pkg/client/listen/v2/websocket/{new_using_callbacks,new_using_chan,client_callback,client_channel,types}.go`, `pkg/client/interfaces/v2/types-flux.go`, `pkg/api/listen/v2/websocket/interfaces/interfaces.go`
2. OpenAPI: `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI: `https://developers.deepgram.com/asyncapi.yaml`
4. Product docs: `https://developers.deepgram.com/docs/flux`

## Gotchas

1. `Connect()` returns `bool`, not `error` -- always check it.
2. `Model` is required on `FluxTranscriptionOptions`; `LanguageHint` only applies when `Model` is `flux-general-multi`.
3. Set `EagerEotThreshold` only if your UX consumes `TurnEventEagerEndOfTurn` -- otherwise stick to `EotThreshold` and the final `TurnEventEndOfTurn`.
4. After `TurnEventEagerEndOfTurn`, a `TurnEventTurnResumed` may arrive; don't commit the eager transcript as final until `TurnEventEndOfTurn` fires.
5. Callback and channel variants are mutually exclusive per connection; pick one per session.

## Example files

- `examples/speech-to-text/websocket/flux_callback/main.go`
- `examples/speech-to-text/websocket/flux_channel/main.go`
