# Deepgram Go SDK
[![Discord](https://dcbadge.vercel.app/api/server/xWRaCDBtW4?style=flat)](https://discord.gg/xWRaCDBtW4)

Community Go SDK for [Deepgram](https://www.deepgram.com/). Start building with our powerful transcription & speech understanding API.

> This SDK only supports hosted usage of api.deepgram.com.

* [Deepgram Go SDK](#deepgram-go-sdk)
* [Getting an API Key](#getting-an-api-key)
* [Documentation](#documentation)
* [Installation](#installation)
* [Requirements](#requirements)
* [Configuration](#configuration)
* [Testing](#testing)
  * [Using Example Projects to test new features](#using-example-projects-to-test-new-features)
* [Transcription](#transcription)
  * [Remote Files](#remote-files)
      * [UrlSource](#urlsource)
  * [Local files](#local-files)
      * [ReadStreamSource](#readstreamsource)
      * [PrerecordedTranscriptionOptions](#prerecordedtranscriptionoptions)
* [Generating Captions](#generating-captions)
  * [Live Audio](#live-audio)
      * [LiveTranscriptionOptions](#livetranscriptionoptions)
* [Projects](#projects)
  * [Get Projects](#get-projects)
  * [Get Project](#get-project)
  * [Update Project](#update-project)
  * [Delete Project](#delete-project)
* [Keys](#keys)
  * [List Keys](#list-keys)
  * [Get Key](#get-key)
  * [Create Key](#create-key)
  * [Delete Key](#delete-key)
* [Members](#members)
  * [Get Members](#get-members)
  * [Remove Member](#remove-member)
* [Scopes](#scopes)
  * [Get Member Scopes](#get-member-scopes)
  * [Update Scope](#update-scope)
* [Invitations](#invitations)
  * [List Invites](#list-invites)
  * [Send Invite](#send-invite)
  * [Delete Invite](#delete-invite)
  * [Leave Project](#leave-project)
* [Usage](#usage)
  * [Get All Requests](#get-all-requests)
      * [UsageRequestListOptions](#usagerequestlistoptions)
  * [Get Request](#get-request)
  * [Summarize Usage](#summarize-usage)
      * [UsageOptions](#usageoptions)
  * [Get Fields](#get-fields)
      * [GetUsageFieldsOptions](#getusagefieldsoptions)
  * [Development and Contributing](#development-and-contributing)
  * [Getting Help](#getting-help)

# Getting an API Key

🔑 To access the Deepgram API you will need a [free Deepgram API Key](https://console.deepgram.com/signup?jump=keys).

# Documentation

You can learn more about the full Deepgram API at [https://developers.deepgram.com](https://developers.deepgram.com).

# Installation

```bash
go get github.com/deepgram-devs/deepgram-go-sdk
```

# Requirements

[Go](https://go.dev/doc/install) (version ^1.18)

# Configuration

```go
dg := deepgram.NewClient(DEEPGRAM_API_KEY)
```

# Testing

## Using Example Projects to test new features

Contributors to the SDK can test their changes locally by running the projects in the `examples` folder. This can be done when making changes without adding a unit test, but of course it is recommended that you add unit tests for any feature additions made to the SDK.

Go to the folder `examples` and look for these two projects, which can be used to test out features in the Deepgram Go SDK:

- prerecorded
- streaming

These are standalone projects, so you will need to follow the instructions in the README.md files for each project to get it running.

# Transcription

## Remote Files

```go
dg := prerecorded.NewClient("DEEPGRAM_API_KEY")
prClient := api.New(dg)

res, err := prClient.PreRecordedFromURL(deepgram.UrlSource{Url: "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"},
deepgram.PreRecordedTranscriptionOptions{Punctuate: true, Utterances: true})
```

[See our Pre-Recorded Quickstart for more info](https://developers.deepgram.com/docs/getting-started-with-pre-recorded-audio).

#### UrlSource

| Property | Value  |          Description          |
| -------- | :----- | :---------------------------: |
| Url      | string | Url of the file to transcribe |

## Local files

```go
dg := prerecorded.NewClient("DEEPGRAM_API_KEY")
prClient := api.New(dg)

file, err := os.Open("PATH_TO_LOCAL_FILE")

if err != nil {
  fmt.Printf("error opening file %s:", file.Name())
  }

source := api.ReadStreamSource{Stream: file, Mimetype: "MIMETYPE_OF_YOUR_FILE"}
res, err := prClient.PreRecordedFromStream(source, deepgram.PreRecordedTranscriptionOptions{Punctuate: true})
if err != nil {
  fmt.Println("ERROR", err)
  return
  }
```

[See our Pre-Recorded Quickstart for more info](https://developers.deepgram.com/docs/getting-started-with-pre-recorded-audio).

#### ReadStreamSource

| Property | Value Type |      reason for      |
| -------- | :--------- | :------------------: |
| Stream   | io.Reader  | stream to transcribe |
| MimeType | string     |  MIMETYPE of stream  |

#### PrerecordedTranscriptionOptions

| Property         | Value Type | Example                            |
| ---------------- | ---------- | ---------------------------------- |
| Model            | string     | `Model: "phonecall"`               |
| Tier             | string     | `Tier: "nova"`                     |
| Version          | string     | `Version: "latest"`                |
| Language         | string     | `Language: "es"`                   |
| DetectLanguage   | bool       | `DetectLanguage: true`             |
| Punctuate        | bool       | `Punctuate: true`                  |
| Profanity_filter | bool       | `Profanity_filter: true`           |
| Redact           | bool       | `Redact: true`                     |
| Diarize          | bool       | `Diarize: true`                    |
| SmartFormat      | bool       | `SmartFormat: true`                |
| Multichannel     | bool       | `Multichannel: true`               |
| Alternatives     | int        | `Alternatives: 2`                  |
| Numerals         | bool       | `Numerals: true`                   |
| Search           | []string   | `Search: []string{"apple"}`        |
| Replace          | []string   | `Replace:[]string{"apple:orange"}` |
| Callback         | string     | `Callback: "https://example.com"`  |
| Keywords         | []string   | `Keywords: []string{"Hannah"}`     |
| Paragraphs       | bool       | `Paragraphs: true`                 |
| Summarize        | bool       | `Summarize: true`                  |
| DetectTopics     | bool       | `DetectTopics: true`               |
| Utterances       | bool       | `Utterances: true`                 |
| Utt_split        | int        | `Utt_split: 9`                     |

# Generating Captions

```go
// The request can be from a local file stream or a URL
// Turn on utterances with {Utterances: true} for captions to work
res, err := prClient.PreRecordedFromStream(source, deepgram.PreRecordedTranscriptionOptions{Punctuate: true, Utterances:true})

// Convert the results to WebVTT format
vtt, err := res.ToWebVTT()

// Convert the results to SRT format
stt, err := res.ToSRT()
```

## Live Audio

```go
options := deepgram.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}
dg, _, err := live.New(options, "DEEPGRAM_API_KEY")
```

#### LiveTranscriptionOptions

See [API Reference](https://developers.deepgram.com/reference/streaming)

# Projects

> projectId and memberId are of type`string`

## Get Projects

Returns all projects accessible by the API key.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)
res, err := mgClient.ListProjects()
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-projects).

## Get Project

Retrieves a specific project based on the provided projectId.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)
res, err := mgClient.GetProject(projectId)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-project).

## Update Project

Update a project.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
options := deepgram.ProjectUpdateOptions{
	Name: "NAME_OF_PROJECT",
	Company:"COMPANY",
}

res, err := mgClient.UpdateProject(projectID, options)
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
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"

res, err := mgClient.DeleteProject(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/delete-project).

# Keys

> projectId,keyId and comment are of type`string`

## List Keys

Retrieves all keys associated with the provided project_id.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"

res, err := mgClient.ListKeys(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/list-keys).

## Get Key

Retrieves a specific key associated with the provided project_id.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
keyID := "YOUR_KEY_ID"

res, err := mgClient.GetKey(projectID, keyID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-key).

## Create Key

Creates an API key with the provided scopes.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
comment := "A comment"
scopes := []string{"admin", "member"}
options := deepgram.CreateKeyOptions{
    ExpirationDate: time.Now().AddDate(1, 0, 0),
		TimeToLive:     3600,
		Tags:           []string{"tag1", "tag2"},
}

res, err := mgClient.CreateKey(projectID, comment, scopes, options)

```

[See our API reference for more info](https://developers.deepgram.com/reference/create-key).

## Delete Key

Deletes a specific key associated with the provided project_id.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
keyID := "YOUR_KEY_ID"

res, err := mgClient.DeleteKey(projectID, keyID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/delete-key).

# Members

> projectId and memberId are of type`string`

## Get Members

Retrieves account objects for all of the accounts in the specified project_id.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"

res, err := mgClient.ListMembers(projectID)

```

[See our API reference for more info](https://developers.deepgram.com/reference/get-members).

## Remove Member

Removes member account for specified member_id.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
memberID := "YOUR_MEMBER_ID"

res, err := mgClient.RemoveMember(projectID, memberID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/remove-member).

# Scopes

> projectId and memberId are of type`string`

## Get Member Scopes

Retrieves scopes of the specified member in the specified project.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
memberID := "YOUR_MEMBER_ID"

res, err := mgClient.GetMemberScopes(projectID, memberID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-member-scopes).

## Update Scope

Updates the scope for the specified member in the specified project.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
memberID := "THEIR_MEMBER_ID"
scope := "SCOPE_TO_ASSIGN"

res, err := mgClient.UpdateMemberScopes(projectID, memberID, scope)
```

[See our API reference for more info](https://developers.deepgram.com/reference/update-scope).

# Invitations

## List Invites

Retrieves all invitations associated with the provided project_id.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"

res, err := mgClient.ListInvitations(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/list-invites).

## Send Invite

Sends an invitation to the provided email address.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
invitationOptions := deepgram.InvitationOptions{
    Email: "",
		Scope: "",
}

res, err := mgClient.SendInvitation(projectID, invitationOptions)
```

[See our API reference for more info](https://developers.deepgram.com/reference/send-invites).

## Delete Invite

Removes the specified invitation from the project.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
email := ""

res, err := mgClient.DeleteInvitation(projectID, email)
```

[See our API reference for more info](https://developers.deepgram.com/reference/delete-invite).

## Leave Project

Removes the authenticated user from the project.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"

res, err := mgClient.LeaveProject(projectID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/leave-project).

# Usage

> projectId and requestId type`string`

## Get All Requests

Retrieves all requests associated with the provided projectId based on the provided options.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
options := deepgram.UsageRequestListOptions{
    Start: "2009-11-10",
		End: "2029-11-10",
		Page: 0,
		Limit:0,
		Status: "failed",
}

res, err := mgClient.ListRequests(projectID, options)
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
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
requestID := "REQUEST_ID"
res, err := mgClient.GetRequest(projectID, requestID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/get-request).

## Summarize Usage

Retrieves usage associated with the provided project_id based on the provided options.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
options := deepgram.UsageOptions{
	  Start: "2009-11-10",
		End: "2029-11-10",
}

res, err := mgClient.GetUsage(projectID, options)
```

#### UsageOptions

| Property | Value  | Description                                                                                           |
| -------- | :----- | :---------------------------------------------------------------------------------------------------- |
| Start    | string | Start date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM |
| End      | string | End date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM   |

[See our API reference for more info](https://developers.deepgram.com/reference/summarize-usage).

## Get Fields

Lists the features, models, tags, languages, and processing method used for requests in the specified project.

```go
dg := manage.New("YOUR_API_KEY")
mgClient := api.New(dg)

projectID := "YOUR_PROJECT_ID"
options := deepgram.UsageRequestListOptions{
    Start: "2009-11-10",
		End: "2029-11-10",
}
res, err := mgClient.GetFields(projectID, options)
```

#### GetUsageFieldsOptions

| Property      | Value    | Description                                                                                           |
| ------------- | :------- | :---------------------------------------------------------------------------------------------------- |
| StartDateTime | DateTime | Start date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM |
| EndDateTime   | DateTime | End date of the requested date range, YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or YYYY-MM-DDTHH:MM:SS+HH:MM   |

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
