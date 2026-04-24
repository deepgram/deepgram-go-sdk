# Flux Channel Example

Streams microphone audio to the [Deepgram Flux (v2/listen)](https://developers.deepgram.com/reference/speech-to-text/listen-flux) turn-based transcription API using the **channel** pattern.

Each Flux server event is delivered to a typed Go channel. `MyHandler.Run()` starts a goroutine for every channel and prints events to stdout. The connection stays open until you press **ENTER**.

## Prerequisites

- Go 1.19+
- A Deepgram API key with access to the Flux endpoint
- A working microphone (ALSA/PulseAudio on Linux, Core Audio on macOS)

## Usage

```bash
export DEEPGRAM_API_KEY=<your-key>

# English-only model (default)
go run main.go

# Explicitly select the English model
go run main.go -model flux-general-en

# Multilingual model — no language hints
go run main.go -model flux-general-multi

# Multilingual model with language hints (triggers a mid-session Configure)
go run main.go -model flux-general-multi -language en -language es -language fr
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-model` | `flux-general-en` | Flux model to use. Valid values: `flux-general-en`, `flux-general-multi` |
| `-language` | _(none)_ | Language hint BCP-47 tag (e.g. `en`, `es`, `fr`, `pt`). Repeat the flag for multiple languages. Only sent when `-model flux-general-multi` is used. |

## Behaviour

- **Single model** (`flux-general-en`): Transcribes English audio with server-side turn detection. Language hints are ignored.
- **Multi model** (`flux-general-multi`): Transcribes multilingual audio. If `-language` flags are provided, a `Configure` message is sent immediately after connecting, passing the hints as `language_hints`. The server will include the detected and hinted languages in each `EndOfTurn` event.

## Turn Events

| Event | Description |
|-------|-------------|
| `StartOfTurn` | Server detected the start of a new speech turn |
| `Update` | Interim transcript update (speech still in progress) |
| `EagerEndOfTurn` | Early end-of-turn signal; speech may resume |
| `TurnResumed` | Speech resumed after an `EagerEndOfTurn` |
| `EndOfTurn` | Final transcript for the turn, including detected languages |
