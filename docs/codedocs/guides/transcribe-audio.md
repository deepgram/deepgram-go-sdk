---
title: "Transcribe Audio"
description: "Transcribe files, streams, or URLs with Deepgram's prerecorded speech-to-text client."
---

This guide shows the practical path for prerecorded transcription: construct the listen REST client, wrap it with the typed API package, and send either a remote URL or local file. The code paths come directly from `pkg/client/listen/v1/rest/client.go` and `pkg/api/listen/v1/rest/rest.go`.

<Steps>
<Step>
### Initialize the client

```go
package main

import (
  "context"
  "fmt"
  "os"

  api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
  client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
  interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
  ctx := context.Background()

  dg := api.New(client.NewRESTWithDefaults())
```

</Step>
<Step>
### Pick transcription options

```go
  options := &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
    SmartFormat: true,
    Punctuate:   true,
    Diarize:     true,
    Utterances:  true,
    Language:    "en-US",
  }
```

</Step>
<Step>
### Send either a URL or a file

```go
  res, err := dg.FromURL(ctx, "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav", options)
  if err != nil {
    fmt.Println("transcription failed:", err)
    os.Exit(1)
  }

  fmt.Println("request:", res.Metadata.RequestID)
  fmt.Println("transcript:", res.Results.Channels[0].Alternatives[0].Transcript)
}
```

</Step>
</Steps>

## Complete Runnable Example

```go
package main

import (
  "context"
  "fmt"
  "os"

  api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
  client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
  interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

func main() {
  ctx := context.Background()

  dg := api.New(client.NewRESTWithDefaults())
  options := &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
    SmartFormat: true,
    Punctuate:   true,
    Diarize:     true,
    Utterances:  true,
    Language:    "en-US",
    Redact:      []string{"pci", "ssn"},
  }

  res, err := dg.FromURL(ctx, "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav", options)
  if err != nil {
    fmt.Println("Deepgram request failed:", err)
    os.Exit(1)
  }

  alt := res.Results.Channels[0].Alternatives[0]
  fmt.Println("Request ID:", res.Metadata.RequestID)
  fmt.Println("Transcript:", alt.Transcript)
  fmt.Println("Confidence:", alt.Confidence)
}
```

If you already have audio in memory or on disk, switch to `FromStream(ctx, src, options)` or `FromFile(ctx, path, options)`. Internally those methods still funnel through the same `sendAudio()` path, so the response shape stays identical.

For production services, the main extra step is structured error handling. The SDK returns `*interfaces.StatusError` when Deepgram sends a typed API error, so you can inspect HTTP status, platform error code, and message without parsing raw JSON yourself. That behavior is implemented in `pkg/client/common/v1/rest.go` and reused across the prerecorded, analyze, speak, manage, and auth clients.
