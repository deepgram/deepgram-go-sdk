# Deepgram Go SDK

Community Go SDK for [Deepgram](https://www.deepgram.com/). Start building with our powerful transcription & speech understanding API.

> This SDK only supports hosted usage of api.deepgram.com.

- [Getting an API Key](#getting-an-api-key)
- [Documentation](#documentation)
- [Installation](#installation)
- [Requirements](#requirements)
- [Configuration](#configuration)
- [Transcription](#transcription)
  - [Remote Files](#remote-files)
    - [UrlSource](####urlsource)
  - [Local Files](#local-files)
    - [StreamSource](####streamsource)
    - [PrerecordedTranscriptionOptions](####prerecordedtranscriptionoptions)
  - [Generating Captions](#generating-captions)
  - [Live Audio](#live-audio)
    - [LiveTranscriptionOptions](####livetranscriptionoptions)
- [Projects](#projects)
  - [Get Projects](#get-projects)
  - [Get Project](#get-project)
  - [Update Project](#update-project)
  - [Delete Project](#delete-project)
- [Keys](#keys)
  - [List Keys](#list-keys)
  - [Get Key](#get-key)
  - [Create Key](#create-key)
  - [Delete Key](#delete-key)
- [Members](#members)
  - [Get Members](#get-members)
  - [Remove Member](#remove-member)
- [Scopes](#scopes)
  - [Get Member Scopes](#get-member-scopes)
  - [Update Scope](#update-scope)
- [Invitations](#invitations)
  - [List Invites](#list-invites)
  - [Send Invite](#send-invite)
  - [Delete Invite](#delete-invite)
  - [Leave Project](#leave-project)
- [Usage](#usage)
  - [Get All Requests](#get-all-requests)
    - [ListAllRequestOptions](####listallrequestoptions)
  - [Get Request](#get-request)
  - [Summarize Usage](#summarize-usage)
    - [GetUsageSummaryOptions](####getusagesummaryoptions)
  - [Get Fields](#get-fields)
    - [GetUsageFieldsOptions](####getusagefieldsoptions)
- [Billing](#billing)
  - [Get All Balances](#get-all-balances)
  - [Get Balance](#get-balance)
- [Logging](#Logging)
- [Development and Contributing](#development-and-contributing)
- [Getting Help](#getting-help)

# Getting an API Key

🔑 To access the Deepgram API you will need a [free Deepgram API Key](https://console.deepgram.com/signup?jump=keys).

# Documentation

You can learn more about the full Deepgram API at [https://developers.deepgram.com](https://developers.deepgram.com).

# Installation

```bash
go get https://github.com/deepgram-devs/deepgram-go-sdk
```

# Requirements

[Go](https://go.dev/doc/install) (version ^1.18)

# Configuration

```go
dg := deepgram.NewClient(DEEPGRAM_API_KEY)
```

# Transcription

## Remote Files

```go
credentials := "DEEPGRAM_API_KEY"
dg := deepgram.NewClient(credentials)

res, err := dg.PreRecordedFromURL(deepgram.UrlSource{Url: "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"},
deepgram.PreRecordedTranscriptionOptions{Punctuate: true, Utterances: true})
```

[See our Pre-Recorded Quickstart for more info](https://developers.deepgram.com/docs/getting-started-with-pre-recorded-audio).

## Local files

```go
credentials := "DEEPGRAM_API_KEY"
dg := deepgram.NewClient(credentials)
file, err := os.Open("PATH_TO_LOCAL_FILE")

if err != nil {
  fmt.Printf("error opening file %s:", file.Name())
  }

source := deepgram.ReadStreamSource{Stream: file, Mimetype: "MIMETYPE_OF_YOUR_FILE"}
res, err := dg.PreRecordedFromStream(source, deepgram.PreRecordedTranscriptionOptions{Punctuate: true})
if err != nil {
  fmt.Println("ERROR", err)
  return
  }
```

[See our Pre-Recorded Quickstart for more info](https://developers.deepgram.com/docs/getting-started-with-pre-recorded-audio).

#### PrerecordedTranscriptionOptions

To be updated

# Generating Captions

```go
// The request can be from a local file stream or a URL
// Turn on utterances with {Utterances: true} for captions to work
res, err := dg.PreRecordedFromStream(source, deepgram.PreRecordedTranscriptionOptions{Punctuate: true, Utterances:true})

// Convert the results to WebVTT format
vtt, err := res.ToWebVTT()

// Convert the results to SRT format
stt, err := res.ToSRT()
```

## Live Audio

To be updated

#### LiveTranscriptionOptions

To be updated

# Projects

> projectId and memberId are of type`string`

## Get Projects

Returns all projects accessible by the API key.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
res, err := dg.ListProjects()
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-projects).

## Get Project

Retrieves a specific project based on the provided projectId.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
res, err := dg.GetProject(projectId)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-project).

## Update Project

Update a project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
options := deepgram.ProjectUpdateOptions{
	Name: "NAME_OF_PROJECT",
	Company:"COMPANY",
}

res, err := dg.UpdateProject(projectID, options)
```

**Project Type**

| Property Name | Type   |                       Description                        |
| ------------- | :----- | :------------------------------------------------------: |
| ProjectId     | string |        Unique identifier of the Deepgram project         |
| Name          | string |                   Name of the project                    |
| Company       | string | Name of the company associated with the Deepgram project |

[See our API reference for more info](https://developers.deepgram.com/reference/update-project).

## Delete Project

Delete a project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"

res, err := dg.DeleteProject(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/delete-project).

# Keys

> projectId,keyId and comment are of type`string`

## List Keys

Retrieves all keys associated with the provided project_id.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"

res, err := dg.ListKeys(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/list-keys).

## Get Key

Retrieves a specific key associated with the provided project_id.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
keyID := "YOUR_KEY_ID"

res, err := dg.GetKey(projectID, keyID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-key).

## Create Key

Creates an API key with the provided scopes.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
comment := "A comment"
scopes := []string{"admin", "member"}
options := deepgram.CreateKeyOptions{
    ExpirationDate: time.Now().AddDate(1, 0, 0),
		TimeToLive:     3600,
		Tags:           []string{"tag1", "tag2"},
}

res, err := dg.CreateKey(projectID, comment, scopes, options)

```

[See our API reference for more info](https://developers.deepgram.com/reference/create-key).

## Delete Key

Deletes a specific key associated with the provided project_id.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
keyID := "YOUR_KEY_ID"

res, err := dg.DeleteKey(projectID, keyID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/delete-key).

# Members

> projectId and memberId are of type`string`

## Get Members

Retrieves account objects for all of the accounts in the specified project_id.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"

res, err := dg.ListMembers(projectID)

```

[See our API reference for more info](https://developers.deepgram.com/reference/get-members).

## Remove Member

Removes member account for specified member_id.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
memberID := "YOUR_MEMBER_ID"

res, err := dg.RemoveMember(projectID, memberID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/remove-member).

# Scopes

> projectId and memberId are of type`string`

## Get Member Scopes

Retrieves scopes of the specified member in the specified project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
memberID := "YOUR_MEMBER_ID"

res, err := dg.GetMemberScopes(projectID, memberID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-member-scopes).

## Update Scope

Updates the scope for the specified member in the specified project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
memberID := "THEIR_MEMBER_ID"
scope := "SCOPE_TO_ASSIGN"

res, err := dg.UpdateMemberScopes(projectID, memberID, scope)
```

[See our API reference for more info](https://developers.deepgram.com/reference/update-scope).

# Invitations

## List Invites

Retrieves all invitations associated with the provided project_id.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"

res, err := dg.ListInvitations(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/list-invites).

## Send Invite

Sends an invitation to the provided email address.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
invitationOptions := deepgram.InvitationOptions{
    Email: "",
		Scope: "",
}

res, err := dg.SendInvitation(projectID, invitationOptions)
```

[See our API reference for more info](https://developers.deepgram.com/reference/send-invites).

## Delete Invite

Removes the specified invitation from the project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
email := ""

res, err := dg.DeleteInvitation(projectID, email)
```

[See our API reference for more info](https://developers.deepgram.com/reference/delete-invite).

## Leave Project

Removes the authenticated user from the project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"

res, err := dg.LeaveProject(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/leave-project).

# Usage

> projectId and requestId type`string`

## Get All Requests

Retrieves all requests associated with the provided projectId based on the provided options.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
options := deepgram.UsageRequestListOptions{
    Start: "2009-11-10",
		End: "2029-11-10",
		Page: 0,
		Limit:0,
		Status: "failed",
}

res, err := dg.ListRequests(projectID, options)
```

#### UsageRequestListOptions

| Property | Type   |                                              Description                                              |
| -------- | :----- | :---------------------------------------------------------------------------------------------------: |
| Start    | string | Start date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM |
| End      | string |  End date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM  |
| Page     | int    |                                           Pages to include                                            |
| Limit    | int    |                                      number of results per page                                       |
| Status   | string |                                     Status of requests to return                                      |

[See our API reference for more info](https://developers.deepgram.com/reference/get-all-requests).

## Get Request

Retrieves a specific request associated with the provided projectId.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
requestID := "REQUEST_ID"
res, err := dg.GetRequest(projectID, requestID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-request).

## Summarize Usage

Retrieves usage associated with the provided project_id based on the provided options.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
options := deepgram.UsageOptions{
	  Start: "2009-11-10",
		End: "2029-11-10",
}

res, err := dg.GetUsage(projectID, options)
```

#### UsageOptions

| Property | Value  |                                              Description                                              |
| -------- | :----- | :---------------------------------------------------------------------------------------------------: |
| Start    | string | Start date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM |
| End      | string |  End date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM  |

[See our API reference for more info](https://developers.deepgram.com/reference/summarize-usage).

## Get Fields

Lists the features, models, tags, languages, and processing method used for requests in the specified project.

```go
dg := deepgram.NewClient("YOUR_API_KEY")
projectID := "YOUR_PROJECT_ID"
options := deepgram.UsageRequestListOptions{
    Start: "2009-11-10",
		End: "2029-11-10",
}
res, err := dg.GetFields(projectID, options)
```

#### GetUsageFieldsOptions

| Property      | Value    |                                              Description                                              |
| ------------- | :------- | :---------------------------------------------------------------------------------------------------: |
| StartDateTime | DateTime | Start date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM |
| EndDateTime   | DateTime |  End date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM  |

[See our API reference for more info](https://developers.deepgram.com/reference/get-fields).

## Development and Contributing

Interested in contributing? We ❤️ pull requests!

To make sure our community is safe for all, be sure to review and agree to our
[Code of Conduct](./.github/CODE_OF_CONDUCT.md). Then see the
[Contribution](./.github/CONTRIBUTING.md) guidelines for more information.

## Getting Help

We love to hear from you so if you have questions, comments or find a bug in the
project, let us know! You can either:

- [Open an issue in this repository](https://github.com/deepgram-devs/deepgram-dotnet-sdk/issues/new)
- [Join the Deepgram Github Discussions Community](https://github.com/orgs/deepgram/discussions)
- [Join the Deepgram Discord Community](https://discord.gg/xWRaCDBtW4)

[license]: LICENSE.txt
