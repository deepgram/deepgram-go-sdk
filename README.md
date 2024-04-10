# Deepgram Go SDK

[![Discord](https://dcbadge.vercel.app/api/server/xWRaCDBtW4?style=flat)](https://discord.gg/xWRaCDBtW4)

Official Go SDK for [Deepgram](https://www.deepgram.com/). Start building with our powerful transcription & speech understanding API.

* [Deepgram Go SDK](#deepgram-go-sdk)
* [SDK Documentation](#sdk-documentation)
* [Getting an API Key](#getting-an-api-key)
* [Installation](#installation)
* [Requirements](#requirements)
* [Quickstarts](#quickstarts)
  * [PreRecorded Audio Transcription Quickstart](#prerecorded-audio-transcription-quickstart)
  * [Live Audio Transcription Quickstart](#live-audio-transcription-quickstart)
* [Examples](#examples)
- [Logging](#logging)
* [Testing](#testing)
* [Development and Contributing](#development-and-contributing)
* [Getting Help](#getting-help)

# SDK Documentation

This SDK implements the Deepgram API found at [https://developers.deepgram.com](https://developers.deepgram.com).

Documentation for specifics about the structs, interfaces, and functions of this SDK can be found here: [Go SDK Documentation](https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main)

For documentation relating to Live Audio Transcription:
- Live Client - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/client/live
- Live API - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/live/v1
    - Live API Interfaces - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/live/v1/interfaces

For documentation relating to PreRecorded Audio Transcription and Intelligence:
- PreRecorded Client - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/client/prerecorded
- PreRecorded API - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/prerecorded/v1
    - PreRecorded API Interfaces - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/prerecorded/v1/interfaces

For documentation relating to Text-to-Speech:
- Speak Client - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/client/speak
- Speak API - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/speak/v1
    - Speak API Interfaces - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/speak/v1/interfaces

For documentation relating to Text Intelligence:
- Analyze Client - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/client/analyze
- Analyze API - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/analyze/v1
    - Analyze API Interfaces - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/analyze/v1/interfaces

For documentation relating to Manage API:
- Management Client - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/manage/live
- Manage API - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/api/manage/v1
    - Manage API Interfaces - https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk@main/pkg/manage/live/v1/interfaces

# Getting an API Key

🔑 To access the Deepgram API you will need a [free Deepgram API Key](https://console.deepgram.com/signup?jump=keys).

# Installation

To incorporate this SDK into your project's `go.mod` file, run the following command from your repo:

```bash
go get github.com/deepgram/deepgram-go-sdk
```

# Requirements

[Go](https://go.dev/doc/install) (version ^1.18)

# Quickstarts

This SDK aims to reduce complexity and abtract/hide some internal Deepgram details that clients shouldn't need to know about.  However you can still tweak options and settings if you need.

## PreRecorded Audio Transcription Quickstart

You can find a [walkthrough](https://developers.deepgram.com/docs/pre-recorded-audio-transcription) on our documentation site. Transcribing Pre-Recorded Audio can be done using the following sample code:

```go
// context
ctx := context.Background()

//client
c := client.NewWithDefaults()
dg := prerecorded.New(c)

// transcription options
options := PreRecordedTranscriptionOptions{
	Punctuate:  true,
	Diarize:    true,
	Language:   "en-US",
}	

// send URL
URL := "https://my-domain.com/files/my-conversation.mp3"
res, err := dg.FromURL(ctx, URL, options)
if err != nil {
	log.Fatalf("FromURL failed. Err: %v\n", err)
}
```

## Live Audio Transcription Quickstart

You can find a [walkthrough](https://developers.deepgram.com/docs/live-streaming-audio-transcription) on our documentation site. Transcribing Live Audio can be done using the following sample code:

```go
// options
transcriptOptions := interfaces.LiveTranscriptionOptions{
	Language:    "en-US",
	Punctuate:   true,
	Encoding:    "linear16",
	Channels:    1,
	Sample_rate: 16000,
}

// create a callback for transcription messages
// for example, you can take a look at this example project:
// https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/streaming/microphone/main.go

// create the client
dgClient, err := client.NewWithDefaults(ctx, transcriptOptions, callback)
if err != nil {
	log.Println("ERROR creating LiveTranscription connection:", err)
	return
}

// call connect!
wsconn := dgClient.Connect()
if wsconn == nil {
	log.Println("Client.Connect failed")
	os.Exit(1)
}
```

# Examples

There are examples for **every** API call in this SDK. You can find all of these examples in the [examples folder](https://github.com/deepgram/deepgram-go-sdk/tree/main/examples) at the root of this repo.

These examples provide:

- Speech-to-Text: PreRecorded Audio:

    - From an Audio File - [examples/prerecorded/file](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/prerecorded/file/main.go)
    - From an URL - [examples/prerecorded/url](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/prerecorded/url/main.go)
    - From an Audio Stream - [examples/prerecorded/stream](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/prerecorded/stream/main.go)

- Speech-to-Text: Live Audio:

    - From a Microphone - [examples/streaming/microphone](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/streaming/microphone/main.go)
    - From an HTTP Endpoint - [examples/streaming/http](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/streaming/http/main.go)

- Text-to-Speech

    - Save audio to a Path - [examples/speak/save](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/speak/save/main.go)
    - Save audio to a user-defined Writer - [examples/speak/file](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/speak/file/main.go)
    - Save audio to a Stream/Buffer - [examples/speak/stream](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/speak/stream/main.go)

- Management API exercise the full [CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete) operations for:

    - Balances - [examples/manage/balances](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/balances/main.go)
    - Invitations - [examples/manage/invitations](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/invitations/main.go)
    - Keys - [examples/manage/keys](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/keys/main.go)
    - Members - [examples/manage/members](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/members/main.go)
    - Projects - [examples/manage/projects](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/projects/main.go)
    - Scopes - [examples/manage/scopes](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/scopes/main.go)
    - Usage - [examples/manage/usage](https://github.com/deepgram/deepgram-go-sdk/blob/main/examples/manage/usage/main.go)

To run each example set the `DEEPGRAM_API_KEY` as an environment variable, then `cd` into each example folder and execute the example: `go run main.go`.

# Logging

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

# Testing

TBD

## Development and Contributing

Interested in contributing? We ❤️ pull requests!

To make sure our community is safe for all, be sure to review and agree to our [Code of Conduct](https://github.com/deepgram/deepgram-go-sdk/blob/main/.github/CODE_OF_CONDUCT.md). Then see the [Contribution](https://github.com/deepgram/deepgram-go-sdk/blob/main/.github/CONTRIBUTING.md) guidelines for more information.

## Getting Help

We love to hear from you so if you have questions, comments or find a bug in the
project, let us know! You can either:

- [Open an issue in this repository](https://github.com/deepgram/deepgram-dotnet-sdk/issues/new)
- [Join the Deepgram Github Discussions Community](https://github.com/orgs/deepgram/discussions)
- [Join the Deepgram Discord Community](https://discord.gg/xWRaCDBtW4)
