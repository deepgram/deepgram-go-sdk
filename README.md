# Deepgram Go SDK

[![Discord](https://dcbadge.vercel.app/api/server/xWRaCDBtW4?style=flat)](https://discord.gg/xWRaCDBtW4) [![Go Reference](https://pkg.go.dev/badge/github.com/deepgram/deepgram-go-sdk/v3.svg)](https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk/v3) [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg?style=flat-rounded)](./.github/CODE_OF_CONDUCT.md)

Official Go SDK for [Deepgram](https://www.deepgram.com/). Power your apps with world-class speech and Language AI models.

- [Deepgram Go SDK](#deepgram-go-sdk)
  - [Documentation](#documentation)
  - [Migrating from earlier versions](#migrating-from-earlier-versions)
    - [V1.2 to V1.3](#v12-to-v13)
    - [V1 to V2](#v1-to-v2)
    - [V2.\* to V3](#v2-to-v3)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Getting an API Key](#getting-an-api-key)
  - [Initialization](#initialization)
  - [Pre-Recorded (Synchronous)](#pre-recorded-synchronous)
    - [Remote Files (Synchronous)](#remote-files-synchronous)
    - [Local Files (Synchronous)](#local-files-synchronous)
  - [Pre-Recorded (Asynchronous / Callbacks)](#pre-recorded-asynchronous--callbacks)
    - [Remote Files (Asynchronous)](#remote-files-asynchronous)
    - [Local Files (Asynchronous)](#local-files-asynchronous)
  - [Streaming Audio](#streaming-audio)
  - [Voice Agent](#voice-agent)
  - [Text to Speech REST](#text-to-speech-rest)
  - [Text to Speech Streaming](#text-to-speech-streaming)
  - [Text Intelligence](#text-intelligence)
  - [Authentication](#authentication)
    - [Authentication Methods](#authentication-methods)
    - [Authentication Priority](#authentication-priority)
    - [Environment Variables](#environment-variables)
    - [OAuth-Style Workflow: API Key → Bearer Token](#oauth-style-workflow-api-key--bearer-token)
    - [Grant Token](#grant-token)
    - [Migration Guide](#migration-guide)
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
    - [Get Request](#get-request)
    - [Get Fields](#get-fields)
    - [Summarize Usage](#summarize-usage)
  - [Billing](#billing)
    - [Get All Balances](#get-all-balances)
    - [Get Balance](#get-balance)
  - [Models](#models)
    - [Get All Project Models](#get-all-project-models)
    - [Get Model](#get-model)
  - [On-Prem APIs](#on-prem-apis)
    - [List On-Prem credentials](#list-on-prem-credentials)
    - [Get On-Prem credentials](#get-on-prem-credentials)
    - [Create On-Prem credentials](#create-on-prem-credentials)
    - [Delete On-Prem credentials](#delete-on-prem-credentials)
  - [Logging](#logging)
  - [Testing](#testing)
  - [Backwards Compatibility](#backwards-compatibility)
  - [Development and Contributing](#development-and-contributing)
    - [Getting Help](#getting-help)

## Documentation

You can learn more about the Deepgram API at [developers.deepgram.com](https://developers.deepgram.com/docs).

Documentation for specifics about the structs, interfaces, and functions of this SDK can be found here: [Go SDK Documentation](https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main)

## Migrating from earlier versions

### V1.2 to V1.3

See the [migration guide](https://developers.deepgram.com/sdks/go-sdk/v12-to-v13-migration) for more details.

### V1 to V2

The Voice Agent interfaces have been updated to use the new Voice Agent V1 API. Please refer to our [Documentation](https://developers.deepgram.com/docs/voice-agent-v1-migration) on Migration to new V1 Agent API.

### V2.\* to V3

V3 Introduced a generic object approach for Agent Providers to ease the maintenance overhead of adding new providers see this [PR](https://github.com/deepgram/deepgram-go-sdk/pull/296) for more details.

## Requirements

[Go](https://go.dev/doc/install) (version ^1.19)

## Installation

To incorporate this SDK into your project's `go.mod` file, run the following command from your repo:

```bash
go get github.com/deepgram/deepgram-go-sdk/v3
```

## Getting an API Key

🔑 To access the Deepgram API you will need a [free Deepgram API Key](https://console.deepgram.com/signup?jump=keys).

## Initialization

All of the examples below will require initializing the Deepgram client and inclusion of imports. The SDK supports **dual authentication** with both API Keys and Bearer tokens.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
    interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
    client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

ctx := context.Background()

// Option 1: API Key Authentication (Traditional)
c := client.NewREST("YOUR_API_KEY", &interfaces.ClientOptions{
    Host: "https://api.deepgram.com",
})

// Option 2: Bearer Token Authentication (Modern)
c := client.NewREST("", &interfaces.ClientOptions{
    Host:        "https://api.deepgram.com",
    AccessToken: "YOUR_ACCESS_TOKEN",
})

// Option 3: Environment Variables (Recommended)
// Set DEEPGRAM_API_KEY or DEEPGRAM_ACCESS_TOKEN
c := client.NewWithDefaults()
```

## Pre-Recorded (Synchronous)

### Remote Files (Synchronous)

Transcribe audio from a URL.

```go
// Define Deepgram options
options := &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
}
// Define url to use
URL := "https://dpgr.am/spacewalk.wav"
res, err := dg.FromURL(ctx, URL, options)
if err != nil {
    log.Fatalf("FromURL failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Transcript: %s\n", res.Results.Channels[0].Alternatives[0].Transcript)
```

[See our API reference for more info](https://developers.deepgram.com/reference/speech-to-text-api/listen).

### Local Files (Synchronous)

Transcribe audio from a file.

```go
// Define Deepgram options
options := &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
}
// Define file to use
filePath := "path/to/your/audio.wav"
res, err := dg.FromFile(ctx, filePath, options)
if err != nil {
    log.Fatalf("FromFile failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Transcript: %s\n", res.Results.Channels[0].Alternatives[0].Transcript)
```

[See our API reference for more info](https://developers.deepgram.com/reference/speech-to-text-api/listen).

## Pre-Recorded (Asynchronous / Callbacks)

### Remote Files (Asynchronous)

Transcribe audio from a URL with callback.

```go
// Define Deepgram options
options := &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
}
// Define URL to use
URL := "https://dpgr.am/spacewalk.wav"
callbackURL := "https://your-callback-url.com/webhook"
res, err := dg.FromURLAsync(ctx, URL, callbackURL, options)
if err != nil {
    log.Fatalf("FromURLAsync failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Request ID: %s\n", res.RequestID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/speech-to-text-api/listen).

### Local Files (Asynchronous)

Transcribe audio from a file with callback.

```go
// Define Deepgram options
options := &interfaces.PreRecordedTranscriptionOptions{
    Model:       "nova-3",
}
// Define file to use and Callback URL
filePath := "path/to/your/audio.wav"
callbackURL := "https://your-callback-url.com/webhook"
res, err := dg.FromFileAsync(ctx, filePath, callbackURL, options)
if err != nil {
    log.Fatalf("FromFileAsync failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Request ID: %s\n", res.RequestID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/speech-to-text-api/listen).

## Streaming Audio

Transcribe streaming audio.

```go
import (
    "context"
    "fmt"
    "os"

    interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
    client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

// Initiate Client
client.InitWithDefault()
ctx := context.Background()

// Create WebSocket client with default handler
dgClient, err := client.NewWSUsingChanForDemo(ctx, &interfaces.LiveTranscriptionOptions{})
// Define Deepgram options
options := &interfaces.LiveTranscriptionOptions{
    Model:     "nova-3",
}

// Define Streaming URL
const audioURL = "streaming_audio.url"

// Connect to Deepgram WebSocket
dgClient.Connect()

// Fetch audio from URL
resp, err := http.Get(audioURL)
if err != nil {
    fmt.Printf("Error fetching audio: %v\n", err)
    os.Exit(1)
}
defer resp.Body.Close()

// Stream audio data to Deepgram in background
go dgClient.Stream(bufio.NewReader(resp.Body))

// Wait for user input to exit
fmt.Println("Press ENTER to exit...")
bufio.NewScanner(os.Stdin).Scan()

// Cleanup and close connection
dgClient.Stop()
```

[See our API reference for more info](https://developers.deepgram.com/reference/speech-to-text-api/listen-streaming).

## Voice Agent

Configure a Voice Agent using WebSocket.

```go
import (
    "context"
    "fmt"
    "os"
    "time"

    interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
    client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/agent"
)

// Initialize the SDK
client.InitWithDefault()

// Create context
ctx := context.Background()

// Configure agent settings
options := &interfaces.SettingsOptions{}
options.Language = "en"
options.Agent.Think.Provider.Type = "open_ai"
options.Agent.Think.Provider.Model = "gpt-4o-mini"
options.Agent.Think.Prompt = "You are a helpful AI assistant."
options.Agent.Listen.Provider["type"] = "deepgram"
options.Agent.Listen.Provider["model"] = "nova-3"
// Set speak provider - supports both single and multiple providers for fallback
// Option 1: Single provider (backward compatible)
options.Agent.Speak = interfaces.Speak{
    Provider: map[string]interface{}{
        "type":  "deepgram",
        "model": "aura-2-thalia-en",
    },
}

// Option 2: Add fallback providers (new feature)
// Primary provider is still configured above, fallback providers are additional
options.Agent.SpeakFallback = &[]interfaces.Speak{
    {
        Provider: map[string]interface{}{
            "type":  "open_ai",
            "model": "tts-1",
            "voice": "shimmer",
        },
        Endpoint: &interfaces.Endpoint{
            URL: "https://api.openai.com/v1/audio/speech",
            Headers: map[string]interface{}{
                "authorization": "Bearer {{OPENAI_API_KEY}}",
            },
        },
    },
    {
        Provider: map[string]interface{}{
            "type":  "elevenlabs",
            "model": "eleven_turbo_v2",
            "voice": "alice",
        },
        Endpoint: &interfaces.Endpoint{
            URL: "https://api.elevenlabs.io/v1/text-to-speech",
            Headers: map[string]interface{}{
                "authorization": "Bearer {{ELEVENLABS_API_KEY}}",
            },
        },
    },
}

// Note: The system will use the primary provider first, then fallback to
// the providers in SpeakFallback if the primary fails

options.Agent.Greeting = "Hello, I'm your AI assistant."

// Create Deepgram client (uses default handler that prints to console)
dgClient, err := client.NewWSUsingChanForDemo(ctx, options)
if err != nil {
    fmt.Printf("Error creating client: %v\n", err)
    os.Exit(1)
}

// Connect to Deepgram
dgClient.Connect()

// Keep connection alive
time.Sleep(30 * time.Second)

// Cleanup
dgClient.Stop()
```

This example demonstrates:

- Setting up a WebSocket connection for Voice Agent
- Configuring the agent with speech, language, and audio settings
- Handling various agent events (speech, transcripts, audio)
- Sending audio data and keeping the connection alive

For a complete implementation, you would need to:

1. Add your audio input source (e.g., microphone)
2. Implement audio playback for the agent's responses
3. Handle any function calls if your agent uses them
4. Add proper error handling and connection management

[See our API reference for more info](https://developers.deepgram.com/reference/voice-agent-api/agent).

## Text to Speech REST

Convert text into speech using the REST API.

```go
// Define Deepgram options
options := &interfaces.SpeakOptions{
    Model:      "aura-2-thalia-en",
}

// Convert text to speech and save to file
text := "Hello world!"
filePath := "output.wav"
res, err := dg.ToSave(ctx, filePath, text, options)
if err != nil {
    fmt.Printf("ToSave failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Audio saved to: %s\n", filePath)
```

[See our API reference for more info](https://developers.deepgram.com/reference/text-to-speech-api/speak).

## Text to Speech Streaming

Convert streaming text into speech using a WebSocket.

```go
// Define Deepgram options
options := &interfaces.SpeakWSOptions{
    Model:      "aura-2-thalia-en",
    Encoding:   "linear16",
    SampleRate: 16000,
}

// Create Deepgram client with custom callback
dgClient, err := client.NewWSUsingCallback(ctx, "", &interfaces.ClientOptions{}, options, callback)
    if err != nil {
    fmt.Printf("Error creating client: %v\n", err)
        os.Exit(1)
    }

    // Connect to Deepgram
dgClient.Connect()

// Send text to convert to speech
text := "Hello, this is a text to speech example."
dgClient.SendText(text)
dgClient.Flush()

// Wait for completion and cleanup
dgClient.WaitForComplete()
dgClient.Stop()
```

[See our API reference for more info](https://developers.deepgram.com/reference/text-to-speech-api/speak-streaming).

## Text Intelligence

Analyze text.

```go
// Define Deepgram options
options := &interfaces.AnalyzeOptions{
    Model: "Nova-3"
    // Read options
}

// Define text file to analyze
filePath := "text_to_analyze.txt"

// Analyze text content from file
res, err := dg.FromFile(ctx, filePath, options)
if err != nil {
    fmt.Printf("FromFile failed. Err: %v\n", err)
    os.Exit(1)
}

// Display results
fmt.Printf("Analysis Results: %+v\n", res.Results)
```

[See our API reference for more info](https://developers.deepgram.com/reference/text-intelligence-api/text-read).

## Authentication

The Deepgram Go SDK supports **dual authentication** with both traditional API Keys and modern JWT Bearer tokens. This provides flexibility for different authentication workflows while maintaining 100% backward compatibility.

### Authentication Methods

#### 1. API Key Authentication (Traditional)

Use your Deepgram API Key directly:

```go
// Using API key in client options
c := client.NewREST("YOUR_API_KEY", &interfaces.ClientOptions{
    Host: "https://api.deepgram.com",
})

// Or using environment variable
// Set DEEPGRAM_API_KEY=your_api_key
c := client.NewWithDefaults()
```

#### 2. Bearer Token Authentication (Modern)

Use JWT Bearer tokens for enhanced security:

```go
// Using Bearer token in client options
c := client.NewREST("", &interfaces.ClientOptions{
    Host:        "https://api.deepgram.com",
    AccessToken: "YOUR_ACCESS_TOKEN",
})

// Or using environment variable
// Set DEEPGRAM_ACCESS_TOKEN=your_access_token
c := client.NewWithDefaults()
```

#### 3. Dynamic Authentication Switching

Switch between authentication methods at runtime:

```go
// Initialize client
c := client.NewREST("YOUR_API_KEY", &interfaces.ClientOptions{
    Host: "https://api.deepgram.com",
})

// Switch to Bearer token authentication
c.SetAccessToken("YOUR_ACCESS_TOKEN")

// Switch back to API key authentication
c.SetAPIKey("YOUR_API_KEY")

// Get current effective authentication
token, isBearer := c.GetAuthToken()
if isBearer {
    fmt.Printf("Using Bearer token: %s\n", token)
} else {
    fmt.Printf("Using API key: %s\n", token)
}
```

### Authentication Priority

When multiple authentication methods are provided, the SDK uses the following priority order:

1. **Explicit AccessToken parameter** (highest priority)
2. **Explicit APIKey parameter**
3. **DEEPGRAM_ACCESS_TOKEN environment variable**
4. **DEEPGRAM_API_KEY environment variable** (lowest priority)

```go
// Example: AccessToken takes precedence over APIKey
c := client.NewREST("fallback_api_key", &interfaces.ClientOptions{
    AccessToken: "preferred_access_token", // This will be used
    Host:        "https://api.deepgram.com",
})
```

### Environment Variables

Set your credentials using environment variables:

```bash
# API Key authentication
export DEEPGRAM_API_KEY="your_api_key"

# Bearer token authentication (takes precedence if both are set)
export DEEPGRAM_ACCESS_TOKEN="your_access_token"

# Optional: Custom host
export DEEPGRAM_HOST="https://api.deepgram.com"
```

Then initialize the client without explicit credentials:

```go
// Automatically uses environment variables based on priority
c := client.NewWithDefaults()
```

### OAuth-Style Workflow: API Key → Bearer Token

The SDK supports a complete OAuth-style authentication workflow where you use an API key to obtain a Bearer token, then use that token for subsequent requests:

```go
import (
    "context"

    authAPI "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1"
    authClient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/auth"
    interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
    listenClient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

func main() {
    ctx := context.Background()

    // Step 1: Use API Key to get Bearer token
    authOptions := &interfaces.ClientOptions{
        APIKey: "YOUR_API_KEY",
    }

    authClient := authClient.NewWithOptions(authOptions)
    tokenResponse, err := authClient.GrantToken(ctx, &authAPI.GrantTokenRequest{})
    if err != nil {
        log.Fatal("Failed to get token:", err)
    }

    // Step 2: Use Bearer token for API calls
    clientOptions := &interfaces.ClientOptions{
        AccessToken: tokenResponse.AccessToken, // Bearer token
    }

    listenClient := listenClient.NewRESTWithOptions(clientOptions)

    // Step 3: Make API calls with Bearer authentication
    result, err := listenClient.FromURL(ctx, "https://dpgr.am/spacewalk.wav", &interfaces.PreRecordedTranscriptionOptions{
        Model: "nova-3",
    })
    if err != nil {
        log.Fatal("Transcription failed:", err)
    }

    fmt.Printf("Transcript: %s\n", result.Results.Channels[0].Alternatives[0].Transcript)
}
```

### Grant Token

Creates a temporary token with a 30-second TTL using your API key:

```go
// Grant token using API key
res, err := dg.GrantToken(ctx)
if err != nil {
    fmt.Printf("GrantToken failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Access Token: %s\n", res.AccessToken)
fmt.Printf("Expires In: %d seconds\n", res.ExpiresIn)
```

The generated token can then be used for Bearer authentication:

```go
// Use the granted token for Bearer authentication
clientOptions := &interfaces.ClientOptions{
    AccessToken: res.AccessToken,
}
dgClient := client.NewRESTWithOptions(clientOptions)
```

### Migration Guide

**Existing code continues to work unchanged!** No migration is required for current implementations.

**To adopt Bearer tokens:**

```go
// Old way (still works)
c := client.NewREST("YOUR_API_KEY", &interfaces.ClientOptions{})

// New way with Bearer token
c := client.NewREST("", &interfaces.ClientOptions{
    AccessToken: "YOUR_ACCESS_TOKEN",
})

// Best practice: Use environment variables
// Set DEEPGRAM_ACCESS_TOKEN or DEEPGRAM_API_KEY
c := client.NewWithDefaults()
```

[See our API reference for more info](https://developers.deepgram.com/reference/token-based-auth-api/grant-token).

## Projects

### Get Projects

Returns all projects accessible by the API key.

```go
// Get all projects
res, err := dg.ListProjects(ctx)
if err != nil {
    fmt.Printf("ListProjects failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Projects: %+v\n", res.Projects)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/projects/list).

### Get Project

Retrieves a specific project based on the provided project_id.

```go
// Get specific project
res, err := dg.GetProject(ctx, myProjectId)
if err != nil {
    fmt.Printf("GetProject failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Project: %+v\n", res.Project)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/projects/get).

### Update Project

Update a project.

```go
// Update project
options := &interfaces.ProjectUpdateRequest{
    Name: "Updated Project Name",
}
res, err := dg.UpdateProject(ctx, myProjectId, options)
if err != nil {
    fmt.Printf("UpdateProject failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Update result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/projects/update).

### Delete Project

Delete a project.

```go
// Delete project
res, err := dg.DeleteProject(ctx, myProjectId)
if err != nil {
    fmt.Printf("DeleteProject failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Delete result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/projects/delete).

## Keys

### List Keys

Retrieves all keys associated with the provided project_id.

```go
// List all keys
res, err := dg.ListKeys(ctx, myProjectId)
if err != nil {
    fmt.Printf("ListKeys failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Keys: %+v\n", res.APIKeys)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/keys/list)

### Get Key

Retrieves a specific key associated with the provided project_id.

```go
// Get specific key
res, err := dg.GetKey(ctx, myProjectId, myKeyId)
if err != nil {
    fmt.Printf("GetKey failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Key: %+v\n", res.APIKey)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/keys/get)

### Create Key

Creates an API key with the provided scopes.

```go
// Create new key
options := &interfaces.KeyCreateRequest{
    Comment: "My API Key",
    Scopes:  []string{"admin"},
}
res, err := dg.CreateKey(ctx, myProjectId, options)
if err != nil {
    fmt.Printf("CreateKey failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Created key: %s\n", res.APIKeyID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/keys/create)

### Delete Key

Deletes a specific key associated with the provided project_id.

```go
// Delete key
res, err := dg.DeleteKey(ctx, myProjectId, myKeyId)
if err != nil {
    fmt.Printf("DeleteKey failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Delete result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/keys/delete)

## Members

### Get Members

Retrieves account objects for all of the accounts in the specified project_id.

```go
// List all members
res, err := dg.ListMembers(ctx, myProjectId)
if err != nil {
    fmt.Printf("ListMembers failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Members: %+v\n", res.Members)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/members/list).

### Remove Member

Removes member account for specified member_id.

```go
// Remove member
res, err := dg.RemoveMember(ctx, myProjectId, memberId)
    if err != nil {
    fmt.Printf("RemoveMember failed. Err: %v\n", err)
        os.Exit(1)
    }

fmt.Printf("Remove result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/members/delete).

## Scopes

### Get Member Scopes

Retrieves scopes of the specified member in the specified project.

```go
// Get member scopes
res, err := dg.GetMemberScopes(ctx, myProjectId, memberId)
if err != nil {
    fmt.Printf("GetMemberScopes failed. Err: %v\n", err)
        os.Exit(1)
    }

fmt.Printf("Scopes: %+v\n", res.Scopes)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/scopes/list).

### Update Scope

Updates the scope for the specified member in the specified project.

```go
// Update member scope
options := &interfaces.ScopeUpdateRequest{
    Scope: "admin",
}
res, err := dg.UpdateMemberScopes(ctx, myProjectId, memberId, options)
if err != nil {
    fmt.Printf("UpdateMemberScopes failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Update result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/scopes/update).

## Invitations

### List Invites

Retrieves all invitations associated with the provided project_id.

```go
// List all invitations
res, err := dg.ListInvitations(ctx, myProjectId)
if err != nil {
    fmt.Printf("ListInvitations failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Invitations: %+v\n", res.Invites)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/invitations/list).

### Send Invite

Sends an invitation to the provided email address.

```go
// Send invitation
options := &interfaces.InvitationCreateRequest{
    Email: "user@example.com",
    Scope: "admin",
}
res, err := dg.SendInvitation(ctx, myProjectId, options)
if err != nil {
    fmt.Printf("SendInvitation failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Invitation sent: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/invitations/create).

### Delete Invite

Removes the specified invitation from the project.

```go
// Delete invitation
res, err := dg.DeleteInvitation(ctx, myProjectId, email)
if err != nil {
    fmt.Printf("DeleteInvitation failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Delete result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/invitations/delete).

### Leave Project

```go
// Leave project
res, err := dg.LeaveProject(ctx, myProjectId)
if err != nil {
    fmt.Printf("LeaveProject failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Leave result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/invitations/leave).

## Usage

### Get All Requests

Retrieves all requests associated with the provided project_id based on the provided options.

```go
// Get all requests
res, err := dg.GetUsageRequests(ctx, myProjectId)
if err != nil {
    fmt.Printf("GetUsageRequests failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Requests: %+v\n", res.Requests)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/usage/list-requests).

### Get Request

Retrieves a specific request associated with the provided project_id

```go
// Get specific request
res, err := dg.GetUsageRequest(ctx, myProjectId, requestId)
if err != nil {
    fmt.Printf("GetUsageRequest failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Request: %+v\n", res.Request)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/usage/get-request).

### Get Fields

Lists the features, models, tags, languages, and processing method used for requests in the specified project.

```go
// Get usage fields
res, err := dg.GetUsageFields(ctx, myProjectId)
if err != nil {
    fmt.Printf("GetUsageFields failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Fields: %+v\n", res.Fields)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/usage/list-fields).

### Summarize Usage

`Deprecated` Retrieves the usage for a specific project. Use Get Project Usage Breakdown for a more comprehensive usage summary.

```go
// Get usage summary
res, err := dg.GetUsageSummary(ctx, myProjectId)
if err != nil {
    fmt.Printf("GetUsageSummary failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Usage summary: %+v\n", res.Usage)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/usage/get).

## Billing

### Get All Balances

Retrieves the list of balance info for the specified project.

```go
// Get all balances
res, err := dg.GetBalances(ctx, myProjectId)
if err != nil {
    fmt.Printf("GetBalances failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Balances: %+v\n", res.Balances)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/balances/list).

### Get Balance

Retrieves the balance info for the specified project and balance_id.

```go
// Get specific balance
res, err := dg.GetBalance(ctx, myProjectId, balanceId)
if err != nil {
    fmt.Printf("GetBalance failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Balance: %+v\n", res.Balance)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/balances/get).

## Models

### Get All Project Models

Retrieves all models available for a given project.

```go
// Get all project models
res, err := dg.GetProjectModels(ctx, myProjectId)
if err != nil {
    fmt.Printf("GetProjectModels failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Models: %+v\n", res)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/projects/list-models).

### Get Model

Retrieves details of a specific model.

```go
// Get specific model
res, err := dg.GetProjectModel(ctx, myProjectId, modelId)
    if err != nil {
    fmt.Printf("GetProjectModel failed. Err: %v\n", err)
        os.Exit(1)
    }

fmt.Printf("Model: %+v\n", res.Model)
```

[See our API reference for more info](https://developers.deepgram.com/reference/management-api/projects/get-model).

## On-Prem APIs

### List On-Prem credentials

Lists sets of distribution credentials for the specified project.

```go
// List on-prem credentials
res, err := dg.ListSelfhostedCredentials(ctx, projectId)
if err != nil {
    fmt.Printf("ListSelfhostedCredentials failed. Err: %v\n", err)
        os.Exit(1)
    }

fmt.Printf("Credentials: %+v\n", res.Credentials)
```

[See our API reference for more info](https://developers.deepgram.com/reference/self-hosted-api/list-credentials).

### Get On-Prem credentials

Returns a set of distribution credentials for the specified project.

```go
// Get specific on-prem credentials
res, err := dg.GetSelfhostedCredentials(ctx, projectId, distributionCredentialsId)
if err != nil {
    fmt.Printf("GetSelfhostedCredentials failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Credentials: %+v\n", res.Credentials)
```

[See our API reference for more info](https://developers.deepgram.com/reference/self-hosted-api/get-credentials).

### Create On-Prem credentials

Creates a set of distribution credentials for the specified project.

```go
// Create on-prem credentials
options := &interfaces.SelfhostedCredentialsCreateRequest{
    Comment: "My on-prem credentials",
}
res, err := dg.CreateSelfhostedCredentials(ctx, projectId, options)
if err != nil {
    fmt.Printf("CreateSelfhostedCredentials failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Created credentials: %s\n", res.CredentialsID)
```

[See our API reference for more info](https://developers.deepgram.com/reference/self-hosted-api/create-credentials).

### Delete On-Prem credentials

Deletes a set of distribution credentials for the specified project.

```go
// Delete on-prem credentials
res, err := dg.DeleteSelfhostedCredentials(ctx, projectId, distributionCredentialId)
if err != nil {
    fmt.Printf("DeleteSelfhostedCredentials failed. Err: %v\n", err)
    os.Exit(1)
}

fmt.Printf("Delete result: %s\n", res.Message)
```

[See our API reference for more info](https://developers.deepgram.com/reference/self-hosted-api/delete-credentials).

## Logging

This SDK provides logging as a means to troubleshoot and debug issues encountered. By default, this SDK will enable `Information` level messages and higher (ie `Warning`, `Error`, etc) when you initialize the library as follows:

```go
client.InitWithDefault();
```

To increase the logging output/verbosity for debug or troubleshooting purposes, you can set the `TRACE` level but using this code:

```go
// init library
client.Init(client.InitLib{
    LogLevel: client.LogLevelTrace,
})
```

## Testing

There are several test folders in [/tests](./tests/) you can run:

- unit_test/ - Unit tests
- daily_test/ - Integration/daily tests
- edge_cases/ - Edge case testing
- response_data/ - Test data
- utils/ - Test utilities

To run the tests, you can use the following commands:

Run specific tests in a directory:

```bash
go run filename
```

## Backwards Compatibility

We follow semantic versioning (semver) to ensure a smooth upgrade experience. Within a major version (like `3.*`), we will maintain backward compatibility so your code will continue to work without breaking changes. When we release a new major version (like moving from `2.*` to `3.*`), we may introduce breaking changes to improve the SDK. We'll always document these changes clearly in our release notes to help you upgrade smoothly.

Older SDK versions will receive Priority 1 (P1) bug support only. Security issues, both in our code and dependencies, are promptly addressed. Significant bugs without clear workarounds are also given priority attention.

## Development and Contributing

Interested in contributing? We ❤️ pull requests!

To make sure our community is safe for all, be sure to review and agree to our [Code of Conduct](https://github.com/deepgram/deepgram-go-sdk/blob/main/.github/CODE_OF_CONDUCT.md). Then see the [Contribution](https://github.com/deepgram/deepgram-go-sdk/blob/main/.github/CONTRIBUTING.md) guidelines for more information.

### Getting Help

We love to hear from you so if you have questions, comments or find a bug in the
project, let us know! You can either:

- [Open an issue in this repository](https://github.com/deepgram/deepgram-go-sdk/issues/new)
- [Join the Deepgram Github Discussions Community](https://github.com/orgs/deepgram/discussions)
- [Join the Deepgram Discord Community](https://discord.gg/xWRaCDBtW4)
