---
name: deepgram-go-management-api
description: Use when writing or reviewing Go code in this repo that works with Deepgram management endpoints for projects, keys, members, scopes, invitations, usage, balances, or models. Route live voice runtime to deepgram-go-voice-agent and repo workflow questions to deepgram-go-maintaining-sdk.
---

# Using Deepgram Management API from the Go SDK

## When to use this product

Use this skill for admin and account operations in `pkg/client/manage` and `pkg/api/manage/v1`.

- projects
- API keys
- members and scopes
- invitations
- usage and request history
- balances
- model discovery

Use a different skill when:

- you need live voice sessions (`deepgram-go-voice-agent`)
- you need SDK maintenance workflow (`deepgram-go-maintaining-sdk`)

## Authentication

Set `DEEPGRAM_API_KEY` before constructing the management client.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
	manage "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	client := manage.NewWithDefaults()
	mg := api.New(client)

	projects, err := mg.ListProjects(ctx)
	if err != nil {
		return err
	}

	fmt.Println(projects)
	return nil
}
```

## Key parameters

- constructors
	- `manage.NewWithDefaults()`
	- `manage.New(apiKey, options)`
- common API groups in `pkg/api/manage/v1`, usually via `api.New(client)`
	- projects: `ListProjects`, `GetProject`, `UpdateProject`, `DeleteProject`
	- keys: `ListKeys`, `GetKey`, `CreateKey`, `DeleteKey`
  - members: `ListMembers`, `RemoveMember`
  - scopes: `GetMemberScopes`, `UpdateMemberScopes`
  - invitations: `ListInvitations`, `SendInvitation`, `DeleteInvitation`, `LeaveProject`
  - usage: `ListRequests`, `GetRequest`, `GetFields`, `GetUsage`
  - balances: `ListBalances`, `GetBalance`
  - models: `ListModels`, `GetModels`, `GetModel`, `ListProjectModels`, `GetProjectModels`, `GetProjectModel`
- low-level escape hatch
  - `managev1.Client.APIRequest(...)`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `pkg/client/manage/client.go`
   - `pkg/client/manage/v1/client.go`
   - `pkg/api/manage/v1/projects.go`
   - `pkg/api/manage/v1/keys.go`
   - `pkg/api/manage/v1/members.go`
   - `pkg/api/manage/v1/scopes.go`
   - `pkg/api/manage/v1/invitations.go`
   - `pkg/api/manage/v1/usage.go`
   - `pkg/api/manage/v1/balances.go`
   - `pkg/api/manage/v1/models.go`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/reference/manage/projects/list`
   - `https://developers.deepgram.com/reference/manage/models/list`

## Gotchas

1. The management client is for account/admin APIs, not live Voice Agent runtime.
2. The repo exposes management operations from `pkg/api/manage/v1`; use `api.New(client)` before calling helpers like `ListProjects`.
3. For unsupported convenience helpers, fall back to `APIRequest(...)` rather than inventing a new transport layer.

## Example files in this repo

- `examples/manage/projects/main.go`
- `examples/manage/keys/main.go`
- `examples/manage/models/main.go`
- `examples/manage/usage/main.go`

## Central product skills

For cross-language Deepgram product knowledge — the consolidated API reference, documentation finder, focused runnable recipes, third-party integration examples, and MCP setup — install the central skills:

```bash
npx skills add deepgram/skills
```

This SDK ships language-idiomatic code skills; `deepgram/skills` ships cross-language product knowledge (see `api`, `docs`, `recipes`, `examples`, `starters`, `setup-mcp`).
