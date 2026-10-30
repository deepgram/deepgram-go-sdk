# Flux Callback Example

Streams microphone audio to the [Deepgram Flux (v2/listen)](https://developers.deepgram.com/reference/speech-to-text/listen-flux) turn-based transcription API using the **callback** pattern.

Each Flux server event triggers a method on `MyCallback`. The connection stays open until you press **ENTER**.

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

# Multilingual model with language hints (sent as connect-time query params)
go run main.go -model flux-general-multi -language en -language es -language fr

# Enable numeric formatting after the WebSocket connects
go run main.go -configure-numerals true

# Explicitly disable numeric formatting after the WebSocket connects
go run main.go -configure-numerals false

# Filter profanity, format numerals, and redact sensitive numbers at connection time
go run main.go -profanity-filter -numerals -redact aggressive_numbers
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-model` | `flux-general-en` | Flux model to use. Valid values: `flux-general-en`, `flux-general-multi` |
| `-language` | _(none)_ | Language hint BCP-47 tag (e.g. `en`, `es`, `fr`, `pt`). Repeat the flag for multiple languages. Only sent when `-model flux-general-multi` is used. |
| `-configure-numerals` | _(unset)_ | Sends a post-connect `Configure` message with `numerals` set to `true` or `false`. |
| `-profanity-filter` | `false` | Replace recognized profanity with a non-profane word or remove it. |
| `-numerals` | `false` | Convert written numbers to numerical format. |
| `-redact` | _(none)_ | Sensitive-number redaction mode: `numbers` or `aggressive_numbers`. |

## Behaviour

- **Single model** (`flux-general-en`): Transcribes English audio with server-side turn detection. Language hints are ignored.
- **Multi model** (`flux-general-multi`): Transcribes multilingual audio. If `-language` flags are provided, the hints are passed as `language_hint` query parameters at connect time. The server will include the detected and hinted languages in each `EndOfTurn` event.
- **Numerals**: `-numerals` sets the initial connection-time value. `-configure-numerals` sends `FluxConfigureOptions.Numerals` after the WebSocket connects to update numeric formatting for subsequent turns.
- **Safety options**: `-profanity-filter` and `-redact` are connection-time options. Redaction supports `numbers` and `aggressive_numbers`.

## Turn Events

| Event | Description |
|-------|-------------|
| `StartOfTurn` | Server detected the start of a new speech turn |
| `Update` | Interim transcript update (speech still in progress) |
| `EagerEndOfTurn` | Early end-of-turn signal; speech may resume |
| `TurnResumed` | Speech resumed after an `EagerEndOfTurn` |
| `EndOfTurn` | Final transcript for the turn, including detected languages |
