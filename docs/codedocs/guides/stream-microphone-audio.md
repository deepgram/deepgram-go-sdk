---
title: "Stream Microphone Audio"
description: "Capture microphone audio and stream it to Deepgram's realtime speech-to-text WebSocket API."
---

This flow combines two public packages: `pkg/audio/microphone` for local capture and `pkg/client/listen` for the streaming WebSocket. The reference implementation is close to `examples/speech-to-text/websocket/microphone_callback/main.go`.

<Steps>
<Step>
### Initialize logging and PortAudio

```go
listen.Init(listen.InitLib{LogLevel: listen.LogLevelStandard})
microphone.Initialize()
defer microphone.Teardown()
```

</Step>
<Step>
### Create the WebSocket client

```go
ctx := context.Background()

dg, err := listen.NewWSUsingCallback(ctx, "", &interfaces.ClientOptions{
  EnableKeepAlive: true,
}, &interfaces.LiveTranscriptionOptions{
  Model:          "nova-3",
  Encoding:       "linear16",
  SampleRate:     16000,
  Language:       "en-US",
  InterimResults: true,
  UtteranceEndMs: "1000",
  VadEvents:      true,
}, callback)
if err != nil {
  panic(err)
}
```

</Step>
<Step>
### Connect and stream microphone frames

```go
if !dg.Connect() {
  panic("connect failed")
}
defer dg.Stop()

mic, err := microphone.New(microphone.AudioConfig{
  InputChannels: 1,
  SamplingRate:  16000,
})
if err != nil {
  panic(err)
}

if err := mic.Start(); err != nil {
  panic(err)
}
defer mic.Stop()

go mic.Stream(dg)
```

</Step>
</Steps>

## Complete Runnable Example

```go
package main

import (
  "bufio"
  "context"
  "fmt"
  "os"
  "strings"

  liveapi "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/microphone"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

type Callback struct{ buffer strings.Builder }

func (c *Callback) Open(*liveapi.OpenResponse) error { return nil }
func (c *Callback) Metadata(*liveapi.MetadataResponse) error { return nil }
func (c *Callback) SpeechStarted(*liveapi.SpeechStartedResponse) error { return nil }
func (c *Callback) Close(*liveapi.CloseResponse) error { return nil }
func (c *Callback) Error(err *liveapi.ErrorResponse) error {
  fmt.Println("error:", err.ErrMsg)
  return nil
}
func (c *Callback) UnhandledEvent(by []byte) error {
  fmt.Println("unhandled:", string(by))
  return nil
}
func (c *Callback) UtteranceEnd(*liveapi.UtteranceEndResponse) error {
  if text := strings.TrimSpace(c.buffer.String()); text != "" {
    fmt.Println("final utterance:", text)
    c.buffer.Reset()
  }
  return nil
}
func (c *Callback) Message(m *liveapi.MessageResponse) error {
  if len(m.Channel.Alternatives) == 0 {
    return nil
  }
  text := strings.TrimSpace(m.Channel.Alternatives[0].Transcript)
  if text == "" {
    return nil
  }
  if m.IsFinal {
    c.buffer.WriteString(text + " ")
  } else {
    fmt.Println("partial:", text)
  }
  return nil
}

func main() {
  listen.InitWithDefault()
  microphone.Initialize()
  defer microphone.Teardown()

  ctx := context.Background()
  callback := &Callback{}

  dg, err := listen.NewWSUsingCallback(ctx, "", &interfaces.ClientOptions{
    EnableKeepAlive: true,
  }, &interfaces.LiveTranscriptionOptions{
    Model:          "nova-3",
    Encoding:       "linear16",
    SampleRate:     16000,
    Language:       "en-US",
    InterimResults: true,
    UtteranceEndMs: "1000",
    SmartFormat:    true,
  }, callback)
  if err != nil {
    panic(err)
  }

  if !dg.Connect() {
    panic("connect failed")
  }
  defer dg.Stop()

  mic, err := microphone.New(microphone.AudioConfig{InputChannels: 1, SamplingRate: 16000})
  if err != nil {
    panic(err)
  }
  if err := mic.Start(); err != nil {
    panic(err)
  }
  defer mic.Stop()

  go mic.Stream(dg)

  fmt.Println("Press ENTER to stop.")
  bufio.NewScanner(os.Stdin).Scan()
}
```

If you are not using a microphone, you can replace `mic.Stream(dg)` with `dg.Stream(reader)` and send PCM bytes from a file, network source, or in-memory buffer.
