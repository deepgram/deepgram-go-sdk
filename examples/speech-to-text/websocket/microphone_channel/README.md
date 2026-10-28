# Live API (Real-Time) Example

This example uses the Microphone as input in order to detect conversation insights in what is being said. This example required additional components (for the microphone) to be installed in order for this example to function correctly.

## Configuration

The SDK (and this example) needs to be initialized with your account's credentials `DEEPGRAM_API_KEY`, which are available in your [Deepgram Console][dg-console]. If you don't have a Deepgram account, you can [sign up here][dg-signup] for free.

You must add your `DEEPGRAM_API_KEY` to your list of environment variables. We use environment variables because they are easy to configure, support PaaS-style deployments, and work well in containerized environments like Docker and Kubernetes.

```sh
export DEEPGRAM_API_KEY=YOUR-APP-KEY-HERE
```

## Installation

The Live API (Real-Time) example makes use of a [microphone package](https://github.com/deepgram/deepgram-go-sdk/tree/main/pkg/audio/microphone) contained within the repository. That package makes use of the [PortAudio library](http://www.portaudio.com/) which is a cross-platform open source audio library. If you are on Linux, you can install this library using whatever package manager is available (yum, apt, etc.) on your operating system. If you are on macOS, you can install this library using [brew](https://brew.sh/).

## Language hints (Nova-3 multilingual)

By default the example transcribes English (`model=nova-3`, `language=en-US`). Pass one or more `-language` flags to switch to Nova-3 multilingual (`language=multi`) and send each value as a `language_hint` query parameter. Hints bias language detection toward the languages you expect, which helps with short utterances and code-switching.

```go
tOptions.Language = "multi"
tOptions.LanguageHint = []string{"en", "es"} // sent as language_hint=en&language_hint=es
```

The message channel handler prints, after every final result, the detected languages (`Channel.Alternatives[0].Languages`, most-spoken first) and the hints the server applied (`Channel.Alternatives[0].LanguagesHinted`).

Things to know:

- Hints only take effect with `model=nova-3` and `language=multi`. With any other model or language they are ignored.
- The server normalizes hints to canonical codes. For example, `ES` becomes `es` and `en-GB` becomes `en`, so `LanguagesHinted` may differ from what you sent.
- A code that is not in the Nova-3 multilingual language list makes the server reject the WebSocket handshake with HTTP 400.

## How to test

Run these commands from the repository root with `DEEPGRAM_API_KEY` exported and a working microphone.

1. Baseline (English only, no language output):

   ```sh
   go run ./examples/speech-to-text/websocket/microphone_channel
   ```

   Speak, then press ENTER to exit. You should see `[MessageResponse] (Interim)` and `[MessageResponse] (Final)` lines but no `[Languages]` lines.

2. Multilingual with hints:

   ```sh
   go run ./examples/speech-to-text/websocket/microphone_channel -language en -language es
   ```

   The program prints `Using language hints: en, es`. Say a sentence in English, then one in Spanish. Each final result is followed by a line like:

   ```text
   [Languages]: languages=[es] hinted=[en es]
   ```

   `hinted` always shows `[en es]`, and `languages` follows what you are speaking.

3. Normalization: run with `-language en-GB -language ES`. `hinted` shows `[en es]`, the canonical codes the server applied.

4. Invalid hint: run with `-language xx`. The server rejects the connection with HTTP 400, and the program prints `Client.Connect failed` and exits.

[dg-console]: https://console.deepgram.com/
[dg-signup]: https://console.deepgram.com/signup
