---
title: "Getting Started"
description: "Install the Deepgram Go SDK, understand the problem it solves, and ship a first transcription quickly."
---

The Deepgram Go SDK gives Go applications typed clients for speech-to-text, text intelligence, text-to-speech, management APIs, authentication, and Voice Agent streaming.

## The Problem

- Speech and audio workloads usually split across unrelated REST and WebSocket clients, which makes auth, retries, and request shaping inconsistent.
- Raw HTTP code forces you to hand-build query strings for dozens of transcription and synthesis options, then decode large JSON payloads yourself.
- Realtime speech pipelines need transport details like keepalive, reconnection, finalization, and event routing that are easy to get wrong.
- Teams often need one SDK for transcription, TTS, project management, and token minting, not four separate integrations.

## The Solution

The SDK separates low-level transports in `pkg/client/...` from typed API wrappers in `pkg/api/...`. The client packages own authentication, request setup, and WebSocket lifecycle management, while the API packages expose typed methods like `FromURL`, `ToSave`, and `GrantToken`. That split is visible in files such as `pkg/client/listen/v1/rest/client.go`, `pkg/api/listen/v1/rest/rest.go`, and `pkg/client/common/v1/websocket.go`.

```go
package main

import (
  "context"
  "fmt"

  api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
  client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
  interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
  ctx := context.Background()

  c := client.NewRESTWithDefaults()
  dg := api.New(c)

  res, err := dg.FromURL(ctx, "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav", &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
    SmartFormat: true,
  })
  if err != nil {
    panic(err)
  }

  fmt.Println(res.Results.Channels[0].Alternatives[0].Transcript)
}
```

## Installation

<Tabs items={["go get", "go.mod", "workspace", "local replace"]}>
<Tab value="go get">

```bash
go get github.com/deepgram/deepgram-go-sdk/v3
```

</Tab>
<Tab value="go.mod">

```go
module your-app

go 1.22

require github.com/deepgram/deepgram-go-sdk/v3 latest
```

</Tab>
<Tab value="workspace">

```bash
go work use .
go get github.com/deepgram/deepgram-go-sdk/v3
```

</Tab>
<Tab value="local replace">

```go
require github.com/deepgram/deepgram-go-sdk/v3 v3.0.0

replace github.com/deepgram/deepgram-go-sdk/v3 => ../deepgram-go-sdk
```

</Tab>
</Tabs>

The SDK targets Go 1.19 or newer according to `go.mod`.

## Quick Start

Set `DEEPGRAM_API_KEY` in your environment, then run the smallest useful transcription flow:

```go
package main

import (
  "context"
  "fmt"
  "log"

  api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
  client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
  interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
  ctx := context.Background()

  dg := api.New(client.NewRESTWithDefaults())
  res, err := dg.FromURL(ctx, "https://dpgr.am/spacewalk.wav", &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
    SmartFormat: true,
  })
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Transcript:", res.Results.Channels[0].Alternatives[0].Transcript)
}
```

Expected output:

```text
Transcript: one small step for man one giant leap for mankind
```

## Key Features

- Dual authentication: API key or bearer token, with explicit priority in `pkg/client/interfaces/v1/options.go`.
- Typed REST wrappers for prerecorded audio, text analysis, TTS, auth, and project management.
- Realtime WebSocket clients for speech-to-text, text-to-speech, and Voice Agent sessions.
- Shared request plumbing for headers, custom query parameters, error decoding, and raw-body handling.
- Optional audio helpers for microphone capture and replay in local streaming workflows.

<Cards>
  <Card title="Architecture" href="/docs/architecture">See how the client, API, and transport layers fit together.</Card>
  <Card title="Core Concepts" href="/docs/authentication-and-client-options">Learn the credential model, typed wrappers, and realtime event flow.</Card>
  <Card title="API Reference" href="/docs/api-reference/client-interfaces">Jump straight to constructor signatures, options, and public methods.</Card>
</Cards>
