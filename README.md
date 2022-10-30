# Deepgram Go SDK

Community Go SDK for [Deepgram](https://www.deepgram.com/)'s automated speech recognition APIs.

> This SDK only supports hosted usage of api.deepgram.com.

To access the API you will need a Deepgram account. Sign up for free at [signup][signup].

---

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

go get .
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

- [Open an issue](https://github.com/deepgram-devs/deepgram-go-sdk/issues/new) on this repository
- Tweet at us! We're [@DeepgramAI on Twitter](https://twitter.com/DeepgramAI)

## Further Reading

Check out the Developer Documentation at [https://developers.deepgram.com/](https://developers.deepgram.com/)

[signup]: https://console.deepgram.com/signup?utm_medium=github&utm_source=DEVREL&utm_content=deepgram-go-sdk
[license]: LICENSE.txt
