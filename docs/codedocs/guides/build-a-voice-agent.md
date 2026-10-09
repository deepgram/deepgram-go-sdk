---
title: "Build A Voice Agent"
description: "Create a Voice Agent session that streams microphone audio, receives text events, and returns synthesized speech."
---

This guide follows the architecture in `examples/agent/websocket/simple/main.go`: an agent WebSocket session consumes microphone audio, emits conversation text events, and can return binary audio packets for playback or storage.

<Steps>
<Step>
### Configure the agent session

```go
settings := agent.NewSettingsConfigurationOptions()
settings.Experimental = true
settings.Agent.Greeting = "Hello, thanks for calling."
settings.Agent.Listen.Provider["model"] = "nova-3"
settings.Agent.Think.Provider["model"] = "gpt-4o-mini"
settings.Agent.Speak.Provider["model"] = "aura-2-thalia-en"
```

</Step>
<Step>
### Create a channel handler and connect

```go
ctx := context.Background()

dg, err := agent.NewWSUsingChan(ctx, "", &interfaces.ClientOptions{
  EnableKeepAlive: true,
}, settings, chans)
if err != nil {
  panic(err)
}

if !dg.Connect() {
  panic("connect failed")
}
```

</Step>
<Step>
### Stream microphone audio into the session

```go
microphone.Initialize()
defer microphone.Teardown()

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

  agentapi "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/audio/microphone"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

type Handler struct {
  open      chan *agentapi.OpenResponse
  text      chan *agentapi.ConversationTextResponse
  audio     chan *[]byte
  close     chan *agentapi.CloseResponse
  err       chan *agentapi.ErrorResponse
  unhandled chan *[]byte
}

func NewHandler() *Handler {
  return &Handler{
    open:      make(chan *agentapi.OpenResponse, 1),
    text:      make(chan *agentapi.ConversationTextResponse, 8),
    audio:     make(chan *[]byte, 8),
    close:     make(chan *agentapi.CloseResponse, 1),
    err:       make(chan *agentapi.ErrorResponse, 1),
    unhandled: make(chan *[]byte, 1),
  }
}

func (h Handler) GetBinary() []*chan *[]byte { return []*chan *[]byte{&h.audio} }
func (h Handler) GetOpen() []*chan *agentapi.OpenResponse { return []*chan *agentapi.OpenResponse{&h.open} }
func (h Handler) GetWelcome() []*chan *agentapi.WelcomeResponse { return nil }
func (h Handler) GetConversationText() []*chan *agentapi.ConversationTextResponse { return []*chan *agentapi.ConversationTextResponse{&h.text} }
func (h Handler) GetUserStartedSpeaking() []*chan *agentapi.UserStartedSpeakingResponse { return nil }
func (h Handler) GetAgentThinking() []*chan *agentapi.AgentThinkingResponse { return nil }
func (h Handler) GetFunctionCallRequest() []*chan *agentapi.FunctionCallRequestResponse { return nil }
func (h Handler) GetAgentStartedSpeaking() []*chan *agentapi.AgentStartedSpeakingResponse { return nil }
func (h Handler) GetAgentAudioDone() []*chan *agentapi.AgentAudioDoneResponse { return nil }
func (h Handler) GetClose() []*chan *agentapi.CloseResponse { return []*chan *agentapi.CloseResponse{&h.close} }
func (h Handler) GetError() []*chan *agentapi.ErrorResponse { return []*chan *agentapi.ErrorResponse{&h.err} }
func (h Handler) GetUnhandled() []*chan *[]byte { return []*chan *[]byte{&h.unhandled} }
func (h Handler) GetInjectionRefused() []*chan *agentapi.InjectionRefusedResponse { return nil }
func (h Handler) GetKeepAlive() []*chan *agentapi.KeepAlive { return nil }
func (h Handler) GetSettingsApplied() []*chan *agentapi.SettingsAppliedResponse { return nil }

func main() {
  ctx := context.Background()
  chans := NewHandler()

  settings := agent.NewSettingsConfigurationOptions()
  settings.Agent.Greeting = "Hello, how can I help?"

  dg, err := agent.NewWSUsingChan(ctx, "", &interfaces.ClientOptions{
    EnableKeepAlive: true,
  }, settings, chans)
  if err != nil {
    panic(err)
  }
  if !dg.Connect() {
    panic("connect failed")
  }
  defer dg.Stop()

  microphone.Initialize()
  defer microphone.Teardown()

  mic, err := microphone.New(microphone.AudioConfig{InputChannels: 1, SamplingRate: 16000})
  if err != nil {
    panic(err)
  }
  if err := mic.Start(); err != nil {
    panic(err)
  }
  defer mic.Stop()

  go mic.Stream(dg)
  go func() {
    for msg := range chans.text {
      fmt.Printf("%s: %s
", msg.Role, msg.Content)
    }
  }()

  fmt.Println("Press ENTER to stop.")
  bufio.NewScanner(os.Stdin).Scan()
}
```

The agent client only exposes channel-oriented constructors publicly, so most production wrappers create an internal adapter like the `Handler` above and then fan out the events they actually care about.
