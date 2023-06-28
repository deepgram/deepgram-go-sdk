# Deepgram Go SDK

Community Go SDK for [Deepgram](https://www.deepgram.com/).  Start building with our powerful transcription & speech understanding API.

> This SDK only supports hosted usage of api.deepgram.com.
## Getting an API Key

🔑 To access the Deepgram API you will need a [free Deepgram API Key](https://console.deepgram.com/signup?jump=keys).
## Documentation

You can learn more about the Deepgram API at [developers.deepgram.com](https://developers.deepgram.com/docs).

## Current Status

There is minimal functionality on the SDK but we hope to add more features soon.

While we don't have a stable release, that is because we don't have feature parity with the other SDKs. It is safe to use, but not all features are available.

To process live transcriptions, see the [example](/examples/liveTranscription_example.go).


## Development and Contributing

Interested in contributing? We ❤️ pull requests!

To make sure our community is safe for all, be sure to review and agree to our
[Code of Conduct](./CODE_OF_CONDUCT.md). Then see the
[Contribution](./CONTRIBUTING.md) guidelines for more information.

## Local Installation and Example Set up

Requirements: [Go](https://go.dev/doc/install) (version ^1.18)

- Clone the repository:
```
git clone https://github.com/deepgram-devs/deepgram-go-sdk/
```

- Move into the directory and install the dependencies:

```
cd deepgram-go-sdk

cd deepgram

go get
```

- Add the API key in the `examples/liveTranscription_example.go` file:

```go
dg := *deepgram.NewClient("YOUR_API_KEY")
```

- Run the example:

```
go run examples/liveTranscription_example.go
```

## Getting Help

We love to hear from you so if you have questions, comments or find a bug in the
project, let us know! You can either:

- [Open an issue in this repository](https://github.com/deepgram-devs/deepgram-go-sdk/issues/new)
- [Join the Deepgram Github Discussions Community](https://github.com/orgs/deepgram/discussions)
- [Join the Deepgram Discord Community](https://discord.gg/xWRaCDBtW4)

[license]: LICENSE.txt
