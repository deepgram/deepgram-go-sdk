// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv2

// SpeakV2Options are the query parameters for POST https://api.deepgram.com/v2/speak
// (the Deepgram Flux text-to-speech batch/REST API).
type SpeakV2Options struct {
	// Model is required. Flux TTS model strings follow the format flux-{voice}-{language}
	// (e.g. "flux-haley-en"). Unlike /v1/speak there is no default, and only flux models
	// are accepted — an Aura model string is rejected.
	Model string `json:"model" schema:"model"`
	// Encoding of the output audio (default "mp3"). Valid values: "mp3", "opus", "flac",
	// "aac" (containerized/compressed), or raw "linear16", "mulaw", "alaw".
	Encoding string `json:"encoding,omitempty" schema:"encoding,omitempty"`
	// Container is the file format wrapper where applicable (e.g. "wav" for linear16,
	// "ogg" for opus, "none" for no container). Defaults per encoding.
	Container string `json:"container,omitempty" schema:"container,omitempty"`
	// SampleRate in Hz; supported values depend on Encoding. linear16: 8000, 16000,
	// 24000, 32000, 44100, 48000. mulaw/alaw: 8000, 16000. flac: 8000, 16000, 22050,
	// 32000, 48000. Not applicable to mp3 or opus. Defaults to the model's native rate.
	SampleRate int `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	// BitRate in bits per second for compressed encodings. mp3: 8000, 16000, 24000,
	// 32000, 40000, 48000 (default). opus: 4000-650000. aac: 4000-192000.
	BitRate int `json:"bit_rate,omitempty" schema:"bit_rate,omitempty"`
	// Speed is the speech-rate multiplier: 0.5 to 1.5 in 0.05 increments,
	// default 1.0. Not supported by every model or language.
	Speed float64 `json:"speed,omitempty" schema:"speed,omitempty"`
	// Expressivity of the generated speech on a calm-to-animated axis: -2, -1,
	// 0 (default), 1, 2. Beta: non-default values increase the risk of hallucinations
	// and pronunciation errors.
	Expressivity int `json:"expressivity,omitempty" schema:"expressivity,omitempty"`
	// Callback URL to receive the result asynchronously instead of on the response body.
	// When set, the response body is a JSON acknowledgement {"request_id": "..."}.
	Callback string `json:"callback,omitempty" schema:"callback,omitempty"`
	// CallbackMethod is the HTTP method for the callback request: "POST" (default) or "PUT".
	CallbackMethod string `json:"callback_method,omitempty" schema:"callback_method,omitempty"`
	// MipOptOut opts the request out of Deepgram's model improvement program.
	MipOptOut bool `json:"mip_opt_out,omitempty" schema:"mip_opt_out,omitempty"`
	// Tag attaches tracking metadata to the request. Repeatable.
	Tag []string `json:"tag,omitempty" schema:"tag,omitempty"`
	// Priority is the processing priority for asynchronous (callback) requests.
	// The only supported value is "low".
	Priority string `json:"priority,omitempty" schema:"priority,omitempty"`
}

// SpeakV2WSOptions are the query parameters for wss://api.deepgram.com/v2/speak
// (the Deepgram Flux text-to-speech streaming WebSocket API).
//
// The streaming transport emits raw (non-containerized) audio, so only
// streaming-compatible encodings are supported; compressed and containerized
// encodings (mp3, opus, flac, aac) are batch/REST only.
type SpeakV2WSOptions struct {
	// Model is required. Flux TTS model strings follow the format flux-{voice}-{language}
	// (e.g. "flux-haley-en"). An Aura model string is rejected on /v2/speak.
	Model string `json:"model" schema:"model"`
	// Encoding of the raw output audio (default "linear16"). Valid values: "linear16",
	// "mulaw", "alaw".
	Encoding string `json:"encoding,omitempty" schema:"encoding,omitempty"`
	// SampleRate in Hz. linear16: 8000, 16000, 24000, 32000, 44100, 48000.
	// mulaw/alaw: 8000, 16000. Defaults to the model's native sample rate.
	SampleRate int `json:"sample_rate,omitempty" schema:"sample_rate,omitempty"`
	// Speed is the initial speech-rate multiplier: 0.5 to 1.5 in 0.05 increments,
	// default 1.0. Can be changed mid-stream with Configure.
	Speed float64 `json:"speed,omitempty" schema:"speed,omitempty"`
	// Expressivity of the generated speech on a calm-to-animated axis: -2, -1,
	// 0 (default), 1, 2. Fixed for the connection — not settable via Configure. Beta.
	Expressivity int `json:"expressivity,omitempty" schema:"expressivity,omitempty"`
	// MipOptOut opts the request out of Deepgram's model improvement program.
	MipOptOut bool `json:"mip_opt_out,omitempty" schema:"mip_opt_out,omitempty"`
	// Tag attaches tracking metadata to the request. Repeatable.
	Tag []string `json:"tag,omitempty" schema:"tag,omitempty"`
}

// SpeakV2ConfigureOptions are the parameters for a mid-session Configure control
// message on the /v2/speak WebSocket. Send via WSCallback.Configure() or
// WSChannel.Configure() to adjust synthesis without reconnecting. Omitted fields
// keep their current value; an accepted change takes effect at the next segment
// boundary. The server answers with ConfigureSuccess or ConfigureFailure.
type SpeakV2ConfigureOptions struct {
	// Speed is the speech-rate multiplier: 0.5 to 1.5 in 0.05 increments. It is a
	// pointer so that "leave speed unchanged" (nil) is distinct from an explicit
	// value; an explicit 0 is rejected as out of range rather than silently
	// dropped from the wire message.
	Speed *float64 `json:"speed,omitempty"`
}
