# Flux ForceEndTurn Example (Bring Your Own Turn Detection)

Streams microphone audio to the [Deepgram Flux (v2/listen)](https://developers.deepgram.com/reference/speech-to-text/listen-flux) turn-based transcription API and ends each turn **manually** by sending a `ForceEndTurn` control message — the pattern you would use for push-to-talk, DTMF input, a UI send button, or your own VAD/endpointing stack.

By default the example connects with `eot_threshold=1.0`, which suppresses Flux's native end-of-turn detection entirely, so turns end only when you press **ENTER** (which sends `{"type":"ForceEndTurn"}`) or when `eot_timeout_ms` elapses. Each `EndOfTurn` event prints its `trigger` field so you can see what ended the turn.

> **Note:** `ForceEndTurn` is gated per deployment on the Deepgram side. If your project does not have it enabled yet, contact Deepgram support.

## Prerequisites

- Go 1.19+
- A Deepgram API key with access to the Flux endpoint (and `ForceEndTurn` enabled for your deployment)
- A working microphone (ALSA/PulseAudio on Linux, Core Audio on macOS)

## Usage

```bash
export DEEPGRAM_API_KEY=<your-key>

# Fully manual turn detection (eot_threshold=1.0, the default here)
go run main.go

# Blend manual and native detection: Flux may still end turns on its own,
# but ENTER forces an immediate EndOfTurn at any time
go run main.go -eot-threshold 0.7
```

While running:

- **Speak**, then press **ENTER** to end the turn — the final transcript prints with `trigger=manual`
- Type **q** then **ENTER** to exit

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-model` | `flux-general-en` | Flux model to use. Valid values: `flux-general-en`, `flux-general-multi` |
| `-eot-threshold` | `1.0` | End-of-turn confidence threshold (`0.5`–`1.0`). `1.0` suppresses native detection so turns end only via `ForceEndTurn` or `eot_timeout_ms`. |

## The `trigger` field

Every `EndOfTurn` event carries a `trigger` field reporting what ended the turn:

| Value | Meaning |
|-------|---------|
| `model` | Flux's native end-of-turn detection ended the turn |
| `manual` | The client sent a `ForceEndTurn` message |
| `timeout` | `eot_timeout_ms` elapsed |

The field is an open string — the server may add new values over time.
