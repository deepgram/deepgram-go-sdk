# Flux TTS Streaming (WebSocket) Example

Synthesizes text with the [Deepgram Flux TTS streaming WebSocket](https://developers.deepgram.com/docs/flux-tts/quickstart) — `wss://api.deepgram.com/v2/speak` — using the **callback** client, and saves the audio to `output.wav`.

The streaming transport is the conversational path: text streams in (`Speak`), turns end explicitly (`Flush`), audio streams back as binary frames, and turns are interruptible (`Interrupt`). This example sends one turn and waits for its `SpeechMetadata` event, which arrives after all audio for the turn has been delivered.

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
| `Configure(opts)` | `{"type":"Configure","speed":1.05}` | Adjust the speech rate mid-session |
| `Stop()` | `{"type":"Close"}` | Drain remaining audio, receive `SessionMetadata`, close |

## Notes

- The streaming transport emits raw (non-containerized) audio: `linear16` (default), `mulaw`, or `alaw`. This example writes a WAV header before appending the linear16 frames so the file plays in standard media players.
- Server events not shown here (`SpeechInterrupted`, `ConfigureSuccess`/`ConfigureFailure`, `Warning`) are all handled by the callback interface — see [main.go](main.go).
