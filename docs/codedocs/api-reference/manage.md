---
title: "Manage"
description: "Reference for project, key, member, invitation, usage, model, and billing management APIs."
---

The management API surface is the broadest part of the SDK. All resource-specific methods are thin wrappers around a shared request helper in `pkg/client/manage/v1/client.go`, while typed request and response models live in `pkg/api/manage/v1/interfaces/types.go`.

## Import Paths

- `github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1`
- `github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces`

## Constructors

Source: `pkg/client/manage/client.go`, `pkg/api/manage/v1/manage.go`

```go
func NewWithDefaults() *Client
func New(apiKey string, options *interfaces.ClientOptions) *Client
func New(client interface{}) *api.Client
```

## Shared Transport Method

Source: `pkg/client/manage/v1/client.go`

```go
func (c *Client) APIRequest(ctx context.Context, method, apiPath string, body io.Reader, resBody interface{}, params ...interface{}) error
```

Every method below eventually calls `APIRequest()`.

## Projects

Source: `pkg/api/manage/v1/projects.go`

```go
func (c *Client) ListProjects(ctx context.Context) (*api.ProjectsResult, error)
func (c *Client) GetProject(ctx context.Context, projectID string) (*api.ProjectResult, error)
func (c *Client) UpdateProject(ctx context.Context, projectID string, proj *api.ProjectUpdateRequest) (*api.MessageResult, error)
func (c *Client) DeleteProject(ctx context.Context, projectID string) (*api.MessageResult, error)
```

Example:

```go
projects, _ := dg.ListProjects(ctx)
projectID := projects.Projects[0].ProjectID
project, _ := dg.GetProject(ctx, projectID)
_, _ = dg.UpdateProject(ctx, projectID, &manageinterfaces.ProjectUpdateRequest{Name: project.Name})
```

## Keys

Source: `pkg/api/manage/v1/keys.go`

```go
func (c *Client) ListKeys(ctx context.Context, projectID string) (*api.KeysResult, error)
func (c *Client) GetKey(ctx context.Context, projectID, keyID string) (*api.KeyResult, error)
func (c *Client) CreateKey(ctx context.Context, projectID string, key *api.KeyCreateRequest) (*api.APIKey, error)
func (c *Client) DeleteKey(ctx context.Context, projectID, keyID string) (*api.MessageResult, error)
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `projectID` | `string` | — | Project identifier. |
| `keyID` | `string` | — | Key identifier for lookup or deletion. |
| `key` | `*api.KeyCreateRequest` | required for create | Request body containing `Comment`, `Scopes`, and optional expiration fields. |

Example:

```go
created, _ := dg.CreateKey(ctx, projectID, &manageinterfaces.KeyCreateRequest{
  Comment: "ci worker",
  Scopes:  []string{"usage:read"},
})
_, _ = dg.DeleteKey(ctx, projectID, created.KeyID)
```

## Members And Scopes

Source: `pkg/api/manage/v1/members.go`, `pkg/api/manage/v1/scopes.go`

```go
func (c *Client) ListMembers(ctx context.Context, projectID string) (*api.MembersResult, error)
func (c *Client) RemoveMember(ctx context.Context, projectID, memberID string) (*api.MessageResult, error)
func (c *Client) GetMemberScopes(ctx context.Context, projectID, memberID string) (*api.ScopeResult, error)
func (c *Client) UpdateMemberScopes(ctx context.Context, projectID, memberID string, scope *api.ScopeUpdateRequest) (*api.MessageResult, error)
```

## Invitations

Source: `pkg/api/manage/v1/invitations.go`

```go
func (c *Client) ListInvitations(ctx context.Context, projectID string) (*api.InvitationsResult, error)
func (c *Client) SendInvitation(ctx context.Context, projectID string, invite *api.InvitationRequest) (*api.MessageResult, error)
func (c *Client) DeleteInvitation(ctx context.Context, projectID, email string) (*api.MessageResult, error)
func (c *Client) LeaveProject(ctx context.Context, projectID string) (*api.MessageResult, error)
```

Example:

```go
_, _ = dg.SendInvitation(ctx, projectID, &manageinterfaces.InvitationRequest{
  Email: "teammate@example.com",
  Scope: "member",
})
```

## Usage

Source: `pkg/api/manage/v1/usage.go`

```go
func (c *Client) ListRequests(ctx context.Context, projectID string, use *api.UsageListRequest) (*api.UsageListResult, error)
func (c *Client) GetRequest(ctx context.Context, projectID, requestID string) (*api.UsageRequestResult, error)
func (c *Client) GetFields(ctx context.Context, projectID string, use *api.UsageListRequest) (*api.UsageFieldResult, error)
func (c *Client) GetUsage(ctx context.Context, projectID string, use *api.UsageRequest) (*api.UsageResult, error)
```

Example:

```go
summary, _ := dg.GetUsage(ctx, projectID, &manageinterfaces.UsageRequest{})
fields, _ := dg.GetFields(ctx, projectID, &manageinterfaces.UsageListRequest{})
_, _ = summary, fields
```

## Billing

Source: `pkg/api/manage/v1/balances.go`

```go
func (c *Client) ListBalances(ctx context.Context, projectID string) (*api.BalancesResult, error)
func (c *Client) GetBalance(ctx context.Context, projectID, balanceID string) (*api.BalanceResult, error)
```

## Models

Source: `pkg/api/manage/v1/models.go`

```go
func (c *Client) ListModels(ctx context.Context, model *api.ModelRequest) (*api.ModelsResult, error)
func (c *Client) GetModels(ctx context.Context, model *api.ModelRequest) (*api.ModelsResult, error)
func (c *Client) GetModel(ctx context.Context, modelID string) (*api.ModelResult, error)
func (c *Client) ListProjectModels(ctx context.Context, projectID string, model *api.ModelRequest) (*api.ModelsResult, error)
func (c *Client) GetProjectModels(ctx context.Context, projectID string, model *api.ModelRequest) (*api.ModelsResult, error)
func (c *Client) GetProjectModel(ctx context.Context, projectID, modelID string) (*api.ModelResult, error)
```

The typed request and result structs are large, so the best source for field-level inspection is `pkg/api/manage/v1/interfaces/types.go`.
