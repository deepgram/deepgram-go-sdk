---
name: deepgram-go-maintaining-sdk
description: Use when changing, reviewing, or extending the Deepgram Go SDK itself, including package layout, examples, docs, tests, linting, dependency bootstrapping, and release-adjacent maintenance. Route product usage questions to the product skills instead of this maintainer workflow.
---

# Maintaining the Deepgram Go SDK

## When to use this product

Use this skill when you are changing the SDK itself rather than consuming it.

- add or update endpoint wrappers
- adjust package layout under `pkg/client` or `pkg/api`
- update examples, tests, docs, or CI-facing scripts
- align the SDK with new Deepgram API behavior

Use a different skill when the task is just product consumption:

- STT: `deepgram-go-speech-to-text`
- TTS: `deepgram-go-text-to-speech`
- text intelligence: `deepgram-go-text-intelligence`
- audio intelligence: `deepgram-go-audio-intelligence`
- voice agent: `deepgram-go-voice-agent`
- management: `deepgram-go-management-api`

## Authentication

Most maintenance work does not require credentials, but examples and integration-style tests do.

```bash
export DEEPGRAM_API_KEY="your_api_key"
```

## Quick start

Typical maintainer loop:

```bash
make ensure-deps
make lint
make check
go test ./...
```

When running examples, follow `examples/README.md` and use `go run main.go` from the specific example directory after `go mod tidy` if needed.

## Key parameters

- repo shape
  - `pkg/client/...` = user-facing constructors and higher-level clients
  - `pkg/api/...` = endpoint-level operations and payload structs
  - `pkg/api/version/...` = path/version builders
  - `pkg/client/common/v1/...` = shared REST and WebSocket plumbing
  - `examples/` = copyable usage code
  - `tests/` = unit, daily, edge-case, and response-data coverage
- maintainer commands from `Makefile`
  - `make ensure-deps`
  - `make lint`
  - `make check`
- helper scripts in `hack/`
  - `hack/ensure-deps/ensure-dependencies.sh`
  - `hack/check/check-shell.sh`
  - `hack/check/check-yaml.sh`
  - `hack/check/check-mdlint.sh`
- contribution flow
  - root `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md` forward to `.github/CONTRIBUTING.md` and `.github/CODE_OF_CONDUCT.md`

## API reference (layered)

1. In-repo reference
   - `README.md`
   - `docs.go`
   - `CONTRIBUTING.md`
   - `.github/CONTRIBUTING.md`
   - `Makefile`
   - `hack/`
   - `go.mod`
   - `tests/`
2. OpenAPI
   - `https://developers.deepgram.com/openapi.yaml`
3. AsyncAPI
   - `https://developers.deepgram.com/asyncapi.yaml`
4. Context7
   - `/llmstxt/developers_deepgram_llms_txt`
5. Product docs
   - `https://developers.deepgram.com/docs`

## Gotchas

1. This repo is hand-maintained; do not apply Fern-specific regeneration guidance unless the repo adds that tooling later.
2. Preserve the existing layering: package constructors in `pkg/client`, endpoint logic in `pkg/api`, shared plumbing in `pkg/client/common/v1`.
3. Match Go idioms already used here: `context.Context`, explicit `error` returns, `defer conn.Close()`, goroutines/channels for streaming paths.
4. When adding a new endpoint, update examples and tests alongside the client wrapper instead of landing transport-only code.
5. Check the nearest package docs and existing examples before renaming fields or reshaping public structs.

## Example files in this repo

- `README.md`
- `docs.go`
- `examples/README.md`
- `Makefile`
- `hack/ensure-deps/ensure-dependencies.sh`
- `tests/unit_test/prerecorded_test.go`
- `tests/unit_test/agent/agent_speak_test.go`
