---
title: "Manage Projects And Keys"
description: "Use the management API wrapper to list projects, inspect keys, and update account resources from Go."
---

The management surface is broad, but the usage pattern is consistent: create `pkg/client/manage.Client`, wrap it with `pkg/api/manage/v1.Client`, then call typed methods. Internally every method funnels through `Client.APIRequest()` in `pkg/client/manage/v1/client.go`.

<Steps>
<Step>
### Create the management client

```go
ctx := context.Background()
dg := manageapi.New(manage.NewWithDefaults())
```

</Step>
<Step>
### List projects and pick one

```go
projects, err := dg.ListProjects(ctx)
if err != nil {
  panic(err)
}

projectID := projects.Projects[0].ProjectID
```

</Step>
<Step>
### Read or mutate resources

```go
details, _ := dg.GetProject(ctx, projectID)
keys, _ := dg.ListKeys(ctx, projectID)
_, _ = details, keys
```

</Step>
</Steps>

## Complete Runnable Example

```go
package main

import (
  "context"
  "fmt"

  manageapi "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1"
  manageinterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/manage/v1/interfaces"
  "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
)

func main() {
  ctx := context.Background()
  dg := manageapi.New(manage.NewWithDefaults())

  projects, err := dg.ListProjects(ctx)
  if err != nil {
    panic(err)
  }
  if len(projects.Projects) == 0 {
    panic("no projects available")
  }

  projectID := projects.Projects[0].ProjectID
  fmt.Println("project:", projects.Projects[0].Name, projectID)

  project, err := dg.GetProject(ctx, projectID)
  if err != nil {
    panic(err)
  }
  fmt.Println("status:", project.Status)

  _, err = dg.UpdateProject(ctx, projectID, &manageinterfaces.ProjectUpdateRequest{
    Name: project.Name,
  })
  if err != nil {
    panic(err)
  }

  keys, err := dg.ListKeys(ctx, projectID)
  if err != nil {
    panic(err)
  }
  fmt.Println("keys:", len(keys.ApiKeys))
}
```

For write operations, the API wrapper files in `pkg/api/manage/v1` are the best place to inspect request shapes. `CreateKey`, `SendInvitation`, `UpdateProject`, and `UpdateMemberScopes` all marshal their request structs before delegating to the shared request helper.

That shared helper matters when you need consistent behavior across many endpoints. `pkg/client/manage/v1.Client.APIRequest()` builds the final URI through `pkg/api/version`, sets Deepgram headers through the common REST client, and decodes the JSON response into whichever typed result struct the wrapper allocated. Once you understand one management method, the rest of the surface is mechanically similar.
