# Changelog

## [3.9.0](https://github.com/deepgram/deepgram-go-sdk/compare/v3.8.0...v3.9.0) (2026-10-06)


### Features

* **agent:** support current function call protocol ([#372](https://github.com/deepgram/deepgram-go-sdk/issues/372)) ([c5e3b5c](https://github.com/deepgram/deepgram-go-sdk/commit/c5e3b5c329cec32e5edaa940e5062d44bbdd3f72))
* **manage:** add reusable agent configurations and agent variables REST endpoints ([#362](https://github.com/deepgram/deepgram-go-sdk/issues/362)) ([4d2996a](https://github.com/deepgram/deepgram-go-sdk/commit/4d2996a661869c5657ce3d7219f721137ca6fdc8))

## [3.8.0](https://github.com/deepgram/deepgram-go-sdk/compare/v3.7.1...v3.8.0) (2026-09-21)

### Features

* **Listen v2 (Flux):** `WSCallback.ForceEndTurn()` and `WSChannel.ForceEndTurn()` let applications end an active turn from an external signal. `TurnInfoResponse.Trigger` identifies whether the server ended a turn through `model`, `manual`, or `timeout`; set `EotThreshold` to `1.0` to suppress native detection, although `EotTimeoutMs` can still end an idle turn. Available on hosted Flux, not self-hosted deployments. ([#350](https://github.com/deepgram/deepgram-go-sdk/issues/350)) ([b295dfc](https://github.com/deepgram/deepgram-go-sdk/commit/b295dfc76d6e6f923e5c3547060301ec28f44614))
* **Speak v2 (Flux TTS):** add batch `POST /v2/speak` synthesis via `ToStream`, `ToFile`, and `ToSave`, plus callback and channel WebSocket clients with `Speak`, `Flush`, `Interrupt`, `InterruptWithOffset`, and mid-session speed `Configure`. `Finish(ctx)` drains queued audio and awaits final `SessionMetadata`; streaming audio and server lifecycle events are typed. ([#351](https://github.com/deepgram/deepgram-go-sdk/issues/351)) ([e1a420b](https://github.com/deepgram/deepgram-go-sdk/commit/e1a420bb19a5f1481b3ca6982e4ca8e9acd79721))

## [3.7.1](https://github.com/deepgram/deepgram-go-sdk/compare/v3.7.0...v3.7.1) (2026-09-21)

### Bug Fixes

* **SDK:** Correct the SDK `User-Agent` version. ([#357](https://github.com/deepgram/deepgram-go-sdk/issues/357)) ([279b669](https://github.com/deepgram/deepgram-go-sdk/commit/279b669986243113fe1cf018e1708b89a2cdd59e))
* **Speak v1 WebSocket:** Preserve `model_name`, `model_version`, `model_uuid`, and `additional_model_uuids` in `MetadataResponse`. ([#345](https://github.com/deepgram/deepgram-go-sdk/issues/345)) ([b96e163](https://github.com/deepgram/deepgram-go-sdk/commit/b96e16384dc816dada114f579a6acffb1abec5b3))
* **Speak v1 WebSocket:** `Clear()` now sends the supported `Clear` control message. `Reset()` remains as a deprecated alias. ([94d5e2a](https://github.com/deepgram/deepgram-go-sdk/commit/94d5e2a1ce1cbc52f9a0ead4f3959fa9179885ea))
