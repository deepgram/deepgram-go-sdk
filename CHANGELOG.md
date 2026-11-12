# Changelog

## [3.11.0](https://github.com/deepgram/deepgram-go-sdk/compare/v3.10.0...v3.11.0) (2026-11-12)


### Features

* **listen:** add Flux safety options ([#393](https://github.com/deepgram/deepgram-go-sdk/issues/393)) ([c443a2c](https://github.com/deepgram/deepgram-go-sdk/commit/c443a2ce59dd7a8ce404d543449e9534260cfa86))
* **listen:** add support to language_hint in the v1 API ([#380](https://github.com/deepgram/deepgram-go-sdk/issues/380)) ([ab7f944](https://github.com/deepgram/deepgram-go-sdk/commit/ab7f944196d55c46d87fbf777d7be9f9446351f6))
* **listen:** configure Flux numerals ([#394](https://github.com/deepgram/deepgram-go-sdk/issues/394)) ([45cb5fb](https://github.com/deepgram/deepgram-go-sdk/commit/45cb5fb0ca94ffefc29a298eb67cd911c353672c))
* **speak:** expose raw control frames ([#386](https://github.com/deepgram/deepgram-go-sdk/issues/386)) ([02bdba9](https://github.com/deepgram/deepgram-go-sdk/commit/02bdba96f258624c2da82b1989dc3d7228e9b859))


### Bug Fixes

* allow nil client options ([#381](https://github.com/deepgram/deepgram-go-sdk/issues/381)) ([47a0b70](https://github.com/deepgram/deepgram-go-sdk/commit/47a0b708cdf01ce2a7f54012843d713747ee5d1f))
* **listen:** synchronize Flux STT reconnect workers ([#400](https://github.com/deepgram/deepgram-go-sdk/issues/400)) ([a995c77](https://github.com/deepgram/deepgram-go-sdk/commit/a995c77fb0505251ce3b28e911a863a9d38a119e))
* **websocket:** close sessions on context cancellation ([#384](https://github.com/deepgram/deepgram-go-sdk/issues/384)) ([272d09b](https://github.com/deepgram/deepgram-go-sdk/commit/272d09b2770508cc35183f6c4f41331618901d6e))
* **WebSocket:** Expose the server's `code` on WebSocket `Error` and `Warning` events in a new `Code` field, and fill it from the close reason for transport errors the SDK reports itself. Re-marshaled `Error` and `Warning` values, including REST errors, now use the lowercase `type` JSON key instead of `Type`. ([a4cd6ab](https://github.com/deepgram/deepgram-go-sdk/commit/a4cd6ab5461e3e1f0abcea8b5491b4738ecbe2b5))
* **websocket:** surface handshake errors ([#383](https://github.com/deepgram/deepgram-go-sdk/issues/383)) ([4913c47](https://github.com/deepgram/deepgram-go-sdk/commit/4913c47cc76f7121fe7308350a84f2300b9cde2f))

## [3.10.0](https://github.com/deepgram/deepgram-go-sdk/compare/v3.9.0...v3.10.0) (2026-10-13)

### Features

* **Listen v1 WebSocket:** Add `LiveTranscriptionOptions.DetectEntities` to enable named entity recognition. Entity results are available in top-level `MessageResponse.Entities`, including raw values and zero-valued fields. ([#318](https://github.com/deepgram/deepgram-go-sdk/issues/318)) ([d79153d](https://github.com/deepgram/deepgram-go-sdk/commit/d79153da75475244e7663046e570d1cd9e8f9866))
* **Listen v1 WebSocket:** Preserve connection `extra` metadata in `UtteranceEndResponse`. Callback and channel clients retain metadata from prior `Results` and `Metadata` messages when the UtteranceEnd wire message omits it. ([#373](https://github.com/deepgram/deepgram-go-sdk/issues/373)) ([d81545f](https://github.com/deepgram/deepgram-go-sdk/commit/d81545f74a2983918f3bfed00a96d42bbfb812b1))

### Bug Fixes

* **Common:** Defer SDK flag parsing in Go test binaries so the test runner retains ownership of `-test.*` flags, while application and `klog` flags remain available. ([#312](https://github.com/deepgram/deepgram-go-sdk/issues/312)) ([7c5b1d3](https://github.com/deepgram/deepgram-go-sdk/commit/7c5b1d353268cdcced5dd8a21d1d50e01dc96f11))
* **Listen v1 WebSocket:** Ignore result messages with no transcript alternatives where a transcript is required, preventing index-out-of-range panics in callback and channel auto-flush and default handlers. ([#369](https://github.com/deepgram/deepgram-go-sdk/issues/369)) ([d5222cf](https://github.com/deepgram/deepgram-go-sdk/commit/d5222cfea61a187cc62c0c10c64254fe5cd385f2))
* **REST:** Configure the HTTP transport with `ClientOptions.Proxy`, matching WebSocket behavior and preventing direct connections when a proxy is supplied. ([#316](https://github.com/deepgram/deepgram-go-sdk/issues/316)) ([a6a4895](https://github.com/deepgram/deepgram-go-sdk/commit/a6a4895c5d7e8259c04f082d6ae7953d25edca48))
* **REST v1:** Surface JSON and non-JSON response bodies for non-400 errors through `StatusError.DeepgramError.Description`, while preserving typed errors and existing HTTP 400 handling. ([#344](https://github.com/deepgram/deepgram-go-sdk/issues/344)) ([e0c7563](https://github.com/deepgram/deepgram-go-sdk/commit/e0c7563eec5336e001dcac119f10da8bd465d820))
* **WebSocket:** Stop keepalive, auto-flush, and ping workers when Listen v1 and v2, Speak v1, and Agent v1 connections close or reconnect, preventing duplicate background workers and control messages. ([#370](https://github.com/deepgram/deepgram-go-sdk/issues/370)) ([d28007a](https://github.com/deepgram/deepgram-go-sdk/commit/d28007a6f0c6a6b43c45c2ae9a1137b2e7888454))

## [3.9.0](https://github.com/deepgram/deepgram-go-sdk/compare/v3.8.0...v3.9.0) (2026-10-06)

### Features

* **Voice Agent:** Support the current function-call protocol with typed function requests, responses, and cancellation events. Function settings accept arbitrary JSON Schema properties, and client-side functions can omit an endpoint. See the runnable function-call example for the response and cancellation flow. ([#372](https://github.com/deepgram/deepgram-go-sdk/issues/372)) ([c5e3b5c](https://github.com/deepgram/deepgram-go-sdk/commit/c5e3b5c329cec32e5edaa940e5062d44bbdd3f72))
* **Management:** Add CRUD endpoints for reusable Voice Agent configurations and Agent template variables. Agent configuration is stored and returned as a JSON string; configuration updates replace metadata only, while variable values accept any JSON type. ([#362](https://github.com/deepgram/deepgram-go-sdk/issues/362)) ([4d2996a](https://github.com/deepgram/deepgram-go-sdk/commit/4d2996a661869c5657ce3d7219f721137ca6fdc8))

## [3.8.0](https://github.com/deepgram/deepgram-go-sdk/compare/v3.7.1...v3.8.0) (2026-09-21)

### Features

* **Listen v2 (Flux):** `WSCallback.ForceEndTurn()` and `WSChannel.ForceEndTurn()` let applications end an active turn from an external signal. `TurnInfoResponse.Trigger` identifies whether the server ended a turn through `model`, `manual`, or `timeout`; set `EotThreshold` to `1.0` to suppress native detection, although `EotTimeoutMs` can still end an idle turn. Available on hosted Flux, not self-hosted deployments. ([#350](https://github.com/deepgram/deepgram-go-sdk/issues/350)) ([b295dfc](https://github.com/deepgram/deepgram-go-sdk/commit/b295dfc76d6e6f923e5c3547060301ec28f44614))
* **Speak v2 (Flux TTS):** add batch `POST /v2/speak` synthesis via `ToStream`, `ToFile`, and `ToSave`, plus callback and channel WebSocket clients with `Speak`, `Flush`, `Interrupt`, `InterruptWithOffset`, and mid-session speed `Configure`. `Finish(ctx)` drains queued audio and awaits final `SessionMetadata`; streaming audio and server lifecycle events are typed. ([#351](https://github.com/deepgram/deepgram-go-sdk/issues/351)) ([e1a420b](https://github.com/deepgram/deepgram-go-sdk/commit/e1a420bb19a5f1481b3ca6982e4ca8e9acd79721))

## [3.7.1](https://github.com/deepgram/deepgram-go-sdk/compare/v3.7.0...v3.7.1) (2026-09-21)

### Bug Fixes

* **SDK:** Correct the SDK `User-Agent` version. ([#357](https://github.com/deepgram/deepgram-go-sdk/issues/357)) ([279b669](https://github.com/deepgram/deepgram-go-sdk/commit/279b669986243113fe1cf018e1708b89a2cdd59e))
* **Speak v1 WebSocket:** Preserve `model_name`, `model_version`, `model_uuid`, and `additional_model_uuids` in `MetadataResponse`. ([#345](https://github.com/deepgram/deepgram-go-sdk/issues/345)) ([b96e163](https://github.com/deepgram/deepgram-go-sdk/commit/b96e16384dc816dada114f579a6acffb1abec5b3))
* **Speak v1 WebSocket:** `Clear()` now sends the supported `Clear` control message. `Reset()` remains as a deprecated alias. ([94d5e2a](https://github.com/deepgram/deepgram-go-sdk/commit/94d5e2a1ce1cbc52f9a0ead4f3959fa9179885ea))
