# Flux TTS Streaming (WebSocket) Example

Synthesizes text with the [Deepgram Flux TTS streaming WebSocket](https://developers.deepgram.com/docs/flux-tts/quickstart) — `wss://api.deepgram.com/v2/speak` — using the **callback** client, and saves the audio to `output.wav`.

The streaming transport is the conversational path: text streams in (`Speak`), turns end explicitly (`Flush`), audio streams back as binary frames, and turns are interruptible (`Interrupt`). This example sends one turn and closes the session gracefully with `Finish`. When `Finish` completes before its context deadline, the server drains every remaining audio frame and reports the final `SessionMetadata`.

## Prerequisites

- Go 1.19+
- A Deepgram API key with access to the Flux TTS endpoint

## Usage

```bash
export DEEPGRAM_API_KEY=<your-key>

# Defaults: flux-haley-en, linear16 @ 48000 Hz saved as output.wav
go run main.go

# Pick a different Flux voice and text
go run main.go -model flux-alexis-en -text "Hello from Flux."
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-model` | `flux-haley-en` | Flux TTS model, in the form `flux-{voice}-{language}`. Required — Aura model strings are rejected on `/v2/speak`. |
| `-text` | a sample sentence | The text to synthesize. |

## Client messages

| Method | Wire message | Purpose |
|--------|--------------|---------|
| `Speak(text)` | `{"type":"Speak","text":"..."}` | Send text into the active turn |
| `Flush()` | `{"type":"Flush"}` | End the active turn; the server drains the buffer and reports `SpeechMetadata` |
| `Interrupt()` / `InterruptWithOffset(ms)` | `{"type":"Interrupt", ...}` | Report a user barge-in; with an offset the reply splits text into spoken/remaining |
| `Configure(opts)` | `{"type":"Configure","speed":1.05}` | Adjust the speech rate mid-session (`0.5`–`1.5` in `0.05` steps) |
| `Finish(ctx)` | `{"type":"Close"}` | Graceful close: wait (bounded by `ctx`) while the server drains all queued audio, sends `SessionMetadata`, and closes |
| `Stop()` | — | Immediate abort: tears down the connection without waiting; queued synthesis is discarded |

## Notes

- The streaming transport emits raw (non-containerized) audio: `linear16` (default), `mulaw`, or `alaw`. This example streams the linear16 frames into a WAV file via `pkg/audio/wav`, which patches the header's RIFF and data sizes once the total length is known, so the finished file is a conforming WAV that players and parsers accept. The writer is linear16 (PCM) only — for `mulaw`/`alaw`, save the raw frames instead.
- Server events not shown here (`SpeechInterrupted`, `ConfigureSuccess`/`ConfigureFailure`, `Warning`) are all handled by the callback interface — see [main.go](main.go).
