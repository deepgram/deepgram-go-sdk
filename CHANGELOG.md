# Changelog

## [3.7.1](https://github.com/deepgram/deepgram-go-sdk/compare/v3.7.0...v3.7.1) (2026-09-21)

### Bug Fixes

* **SDK:** Correct the SDK `User-Agent` version. ([#357](https://github.com/deepgram/deepgram-go-sdk/issues/357)) ([279b669](https://github.com/deepgram/deepgram-go-sdk/commit/279b669986243113fe1cf018e1708b89a2cdd59e))
* **Speak v1 WebSocket:** Preserve `model_name`, `model_version`, `model_uuid`, and `additional_model_uuids` in `MetadataResponse`. ([#345](https://github.com/deepgram/deepgram-go-sdk/issues/345)) ([b96e163](https://github.com/deepgram/deepgram-go-sdk/commit/b96e16384dc816dada114f579a6acffb1abec5b3))
* **Speak v1 WebSocket:** `Clear()` now sends the supported `Clear` control message. `Reset()` remains as a deprecated alias. ([94d5e2a](https://github.com/deepgram/deepgram-go-sdk/commit/94d5e2a1ce1cbc52f9a0ead4f3959fa9179885ea))
