# Flux TTS Batch (REST) Example

Synthesizes a complete block of text with the [Deepgram Flux TTS batch endpoint](https://developers.deepgram.com/docs/flux-tts/batch) — `POST /v2/speak` — and saves the audio to `output.mp3`.

Use the batch transport to pre-generate fixed audio (IVR prompts, notifications, audiobook lines) where the whole text is known up front. For live, interruptible conversations, use the streaming WebSocket example instead ([flux_streaming](../../websocket/flux_streaming/)).

## Prerequisites

- Go 1.19+
- A Deepgram API key with access to the Flux TTS endpoint

## Usage

```bash
export DEEPGRAM_API_KEY=<your-key>

# Defaults: flux-haley-en, mp3 output
go run main.go

# Pick a different Flux voice and text
go run main.go -model flux-alexis-en -text "Hello from Flux."
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-model` | `flux-haley-en` | Flux TTS model, in the form `flux-{voice}-{language}`. Required by the API — there is no default server-side, and Aura model strings are rejected on `/v2/speak` (use `/v1/speak` for Aura voices). |
| `-text` | a sample sentence | The text to synthesize. |

## Notes

- The batch transport defaults to `mp3` encoding; `opus`, `flac`, `aac`, and raw `linear16` / `mulaw` / `alaw` are also available via `SpeakV2Options.Encoding`.
- Supplying `SpeakV2Options.Callback` switches the request to asynchronous processing: the response body becomes a JSON acknowledgement (`{"request_id": "..."}`) and the audio is delivered to the callback URL.
