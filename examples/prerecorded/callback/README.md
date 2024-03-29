# Callback Example

This example shows how to use the [Callback](https://developers.deepgram.com/docs/callback) functionality on the Prerecorded API.

> **_NOTE:_** To use this example, the `endpoint` component must run somewhere with a public-facing IP address. You cannot run this example locally behind your typical firewall.

## Configuration

This example consists of two components:
- `endpoint`: which is an example of what a callback endpoint would look like. Reminder: this requires running with a public-facing IP address
- `callback`: which is just a Deepgram client posts a PreRecorded transcription request using an audio file located at the following URL: [https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav](https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav)

The `callback` component requires the Deepgram API Key environment variable to be configured to run.

```sh
export DEEPGRAM_API_KEY=YOUR-APP-KEY-HERE
```

## Installation

Compile the `endpoint` application and upload the binary and sample SSL certificates to a system on the public-facing internet. This could be an EC2 instance on AWS, an instance on GCP, or etc. Then, run the `endpoint` binary.

On the `callback` project, modify the IP address constant in the code with the public-facing IP address of your EC2, GCP, etc instance.

```Go
hostport string = "<REPLACE WITH YOUR HOSTPORT - FORMAT: 127.0.0.1:3000>"
```

Then run the `callback` application. This can be done from your local laptop.
