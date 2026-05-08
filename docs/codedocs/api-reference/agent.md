---
title: "Agent"
description: "Reference for the Voice Agent WebSocket client, settings objects, and event interfaces."
---

The Voice Agent package is WebSocket-only. It exposes one public session type, `WSChannel`, plus a nested settings object and a broad event interface for text, audio, and tool-call messages.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/agent/v1/websocket/interfaces`

## Constructors And Settings Helper

Source: `pkg/client/agent/client.go`

```go
func NewSettingsConfigurationOptions() *interfaces.SettingsOptions
func NewWSUsingChanForDemo(ctx context.Context, options *interfaces.SettingsOptions) (*listenv1ws.WSChannel, error)
func NewWSUsingChanWithDefaults(ctx context.Context, options *interfaces.SettingsOptions, chans msginterfaces.AgentMessageChan) (*listenv1ws.WSChannel, error)
func NewWSUsingChan(ctx context.Context, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.SettingsOptions, chans msginterfaces.AgentMessageChan) (*listenv1ws.WSChannel, error)
func NewWSUsingChanWithCancel(ctx context.Context, ctxCancel context.CancelFunc, apiKey string, cOptions *interfaces.ClientOptions, tOptions *interfaces.SettingsOptions, chans msginterfaces.AgentMessageChan) (*listenv1ws.WSChannel, error)
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `ctx` | `context.Context` | — | Connection-scoped context. |
| `ctxCancel` | `context.CancelFunc` | — | Optional caller-owned cancel function. |
| `apiKey` | `string` | `""` | Optional explicit API key. |
| `cOptions` | `*interfaces.ClientOptions` | required | Shared transport options. |
| `tOptions` | `*interfaces.SettingsOptions` | required | Agent session configuration payload. |
| `chans` | `msginterfaces.AgentMessageChan` | required | Channel router implementation. |

## Settings Types

Source: `pkg/client/interfaces/v1/types-agent.go`

```go
type SettingsOptions struct {
  Type string
  Tags []string
  Experimental bool
  MipOptOut bool
  Audio Audio
  Agent Agent
}
```

Important nested types:

```go
type Audio struct { Input *Input; Output *Output }
type Agent struct {
  Language string
  Listen Listen
  Think Think
  Speak Speak
  SpeakFallback *[]Speak
  Greeting string
}
type Think struct {
  Provider map[string]interface{}
  Endpoint *Endpoint
  Functions *[]Functions
  Prompt string
  ContextLength any
}
```

Example:

```go
settings := agent.NewSettingsConfigurationOptions()
settings.Agent.Greeting = "Hello, thanks for calling."
settings.Agent.Think.Prompt = "Be concise and confirm important details."
```

## `WSChannel` Methods

Source: `pkg/client/agent/v1/websocket/client_channel.go`

```go
func (c *WSChannel) Connect() bool
func (c *WSChannel) ConnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retryCnt int) bool
func (c *WSChannel) AttemptReconnect(ctx context.Context, retries int64) bool
func (c *WSChannel) AttemptReconnectWithCancel(ctx context.Context, ctxCancel context.CancelFunc, retries int64) bool
func (c *WSChannel) GetURL(host string) (string, error)
func (c *WSChannel) Start()
func (c *WSChannel) ProcessMessage(wsType int, byMsg []byte) error
func (c *WSChannel) Stream(r io.Reader) error
func (c *WSChannel) Write(p []byte) (int, error)
func (c *WSChannel) KeepAlive() error
func (c *WSChannel) GetCloseMsg() []byte
func (c *WSChannel) Finish()
func (c *WSChannel) ProcessError(err error) error
```

Example:

```go
dg, _ := agent.NewWSUsingChan(ctx, "", cOptions, settings, chans)
if dg.Connect() {
  _ = dg.Stream(reader)
}
```

## Agent Event And Command Types

Source: `pkg/api/agent/v1/websocket/interfaces/types.go`, `interfaces/interfaces.go`

```go
type UpdatePrompt struct {
  Type string `json:"type,omitempty"`
  Prompt string `json:"instructions,omitempty"`
}

type UpdateSpeak struct {
  Type string `json:"type,omitempty"`
  Speak interfaces.SpeakOptions `json:"speak,omitempty"`
}

type InjectAgentMessage struct {
  Type string `json:"type,omitempty"`
  Content string `json:"content,omitempty"`
}

type FunctionCallResponse struct {
  Type string `json:"type,omitempty"`
  FunctionCallID string `json:"function_call_id,omitempty"`
  Output string `json:"output,omitempty"`
}
```

Event routing interface:

```go
type AgentMessageChan interface {
  GetBinary() []*chan *[]byte
  GetOpen() []*chan *OpenResponse
  GetWelcome() []*chan *WelcomeResponse
  GetConversationText() []*chan *ConversationTextResponse
  GetUserStartedSpeaking() []*chan *UserStartedSpeakingResponse
  GetAgentThinking() []*chan *AgentThinkingResponse
  GetFunctionCallRequest() []*chan *FunctionCallRequestResponse
  GetAgentStartedSpeaking() []*chan *AgentStartedSpeakingResponse
  GetAgentAudioDone() []*chan *AgentAudioDoneResponse
  GetClose() []*chan *CloseResponse
  GetError() []*chan *ErrorResponse
  GetUnhandled() []*chan *[]byte
  GetInjectionRefused() []*chan *InjectionRefusedResponse
  GetKeepAlive() []*chan *KeepAlive
  GetSettingsApplied() []*chan *SettingsAppliedResponse
}
```

Related page: [Voice Agent Settings](/workspace/home/codedocs-template/content/docs/voice-agent-settings.mdx)
