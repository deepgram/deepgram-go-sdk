---
title: "Auth"
description: "Reference for the Deepgram authentication client and bearer-token grant workflow."
---

The auth surface is intentionally small. It exists to exchange an API key for a temporary bearer token, which you can then reuse through `ClientOptions.AccessToken`.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/auth`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1/interfaces`

## Constructors

Source: `pkg/client/auth/client.go`, `pkg/api/auth/v1/auth.go`

```go
func NewWithDefaults() *Client
func New(apiKey string, options *interfaces.ClientOptions) *Client
func New(client interface{}) *api.Client
```

## Shared Transport Method

Source: `pkg/client/auth/v1/client.go`

```go
func (c *Client) APIRequest(ctx context.Context, method, apiPath string, body io.Reader, resBody interface{}, params ...interface{}) error
```

## Token Grant Method

Source: `pkg/api/auth/v1/grant-token.go`

```go
func (c *Client) GrantToken(ctx context.Context, req *api.GrantTokenRequest) (*api.GrantToken, error)
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `ctx` | `context.Context` | — | Request-scoped context. |
| `req` | `*api.GrantTokenRequest` | `nil` | Optional request body with TTL override. |

### Request And Return Types

Source: `pkg/api/auth/v1/interfaces/types.go`

```go
type GrantTokenRequest struct {
  TTLSeconds *int `json:"ttl_seconds,omitempty"`
}

type GrantToken struct {
  AccessToken string `json:"access_token,omitempty"`
  TokenType string `json:"token_type,omitempty"`
  ExpiresIn float64 `json:"expires_in,omitempty"`
}
```

### Usage Examples

Default TTL:

```go
authAPIClient := authapi.New(auth.NewWithDefaults())
token, err := authAPIClient.GrantToken(ctx, nil)
```

Custom TTL:

```go
ttl := 60
token, err := authAPIClient.GrantToken(ctx, &authinterfaces.GrantTokenRequest{
  TTLSeconds: &ttl,
})
```

Using the returned bearer token:

```go
opts := &interfaces.ClientOptions{AccessToken: token.AccessToken}
client := listen.NewREST("", opts)
_ = client
```

The end-to-end bearer-token flow is demonstrated in `examples/auth/grant-token/main.go`. That example also shows that `ClientOptions.GetAuthToken()` prefers `AccessToken` over `APIKey`, and that you can switch credentials dynamically with `SetAccessToken()` and `SetAPIKey()`. In practice, the common pattern is to keep the management or auth client on API-key auth, mint short-lived bearer tokens, and then pass those bearer tokens to listen, speak, or agent clients that should not hold the full API key.

Because the auth package is intentionally narrow, most application design decisions live one layer up in `ClientOptions`. The auth client gives you the bearer credential; the rest of the SDK decides how that credential is attached to REST requests and WebSocket upgrade headers. That separation is why the same granted token works unchanged with prerecorded transcription, realtime listen sessions, TTS, and Voice Agent connections.
