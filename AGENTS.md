# Agents

Instructions for AI coding agents (Claude Code, Cursor, Codex, Copilot) and for humans working with them in this repository. `CLAUDE.md` includes this file.

## Repository purpose

This is the official Go SDK for the Deepgram API. The module path is `github.com/deepgram/deepgram-go-sdk/v3`, `go.mod` declares `go 1.19`, and the latest release tag is `v3.7.0`. The SDK is hand-written: there is no code generator, no `fern/` folder, and no `.fernignore`. Edit the source directly.

Read `README.md` and `.github/CONTRIBUTING.md` before changing public behavior. Never hardcode API keys or access tokens; examples and tests use the SDK's environment-backed client defaults (`DEEPGRAM_API_KEY`, or `DEEPGRAM_ACCESS_TOKEN` for a bearer token).

## Repository map

| Path | What lives there |
| --- | --- |
| `pkg/api/<product>/<version>` | Public API facades, request and response models, message routers, and callback interfaces |
| `pkg/client/<product>/<version>` | REST and WebSocket transport clients |
| `pkg/client/interfaces` | Shared public option types; versioned types live in `v1` (Nova REST and WebSocket, Speak, Agent, Analyze) and `v2` (`types-flux.go`, the Flux STT options) |
| `pkg/client/common` | Shared REST and WebSocket behavior (`v1` for the Nova, Speak, and Agent sockets, `v2` for the Flux STT socket) |
| `pkg/api/version` | Endpoint URL and query-parameter construction; `constants.go` holds shared default paths and product-specific helpers define their own paths |
| `pkg/audio` | Microphone capture (`microphone`, CGO via PortAudio) and WAV replay helpers |
| `examples/` | Runnable, product-specific `main.go` programs, one directory per scenario |
| `tests/unit_test` | Deterministic unit tests (`package deepgram_test`); no network access |
| `tests/daily_test` | `TestDaily_*` tests that call the live API and refresh fixtures under `tests/response_data` |
| `tests/edge_cases` | Standalone `main.go` programs that exercise reconnect, cancel, keepalive, and timeout paths against the live API |
| `hack/` | Dependency installers (`ensure-deps/`) and lint wrappers (`check/`) used by the `Makefile` and CI |
| `.github/workflows` | CI: `tests-unit.yaml`, `check-lint.yaml`, `check-mdlint.yaml`, `check-shell.yaml`, `check-yaml.yaml`, `check-actionlint.yaml`, `check-all.yaml`, `tests-daily.yaml`, `context7.yml` |
| `.agents/skills` | Agent-agnostic skills for using this SDK (speech-to-text, conversational STT, text-to-speech, voice agent, audio intelligence, text intelligence, management API) |

Legacy `pkg/api/live`, `pkg/api/prerecorded`, `pkg/client/live`, `pkg/client/prerecorded`, and `pkg/client/rest` are deprecated aliases. Use `listen`, `speak`, `analyze`, `agent`, `manage`, and `auth` for new code.

## Client surfaces

Every row below was checked against the source tree on 2026-09-13.

| Product | Endpoint | Package | Status |
| --- | --- | --- | --- |
| Speech-to-text, pre-recorded | `POST /v1/listen` | `pkg/client/listen` (`NewREST`, `NewRESTWithDefaults`) with `pkg/api/listen/v1/rest` | Shipped |
| Speech-to-text, streaming (Nova) | `wss /v1/listen` | `pkg/client/listen` (`NewWSUsingCallback*`, `NewWSUsingChan*`) with `pkg/api/listen/v1/websocket` | Shipped; callback and channel variants |
| Flux STT (conversational speech-to-text) | `wss /v2/listen` | `pkg/client/listen/v2/websocket` with `pkg/api/listen/v2/websocket`; options in `pkg/client/interfaces/v2/types-flux.go` | Shipped; callback and channel variants, `Configure` for mid-session changes |
| Text-to-speech, batch (Aura) | `POST /v1/speak` | `pkg/client/speak` (`NewREST`) with `pkg/api/speak/v1/rest` | Shipped |
| Text-to-speech, streaming (Aura) | `wss /v1/speak` | `pkg/client/speak` (`NewWSUsingCallback*`, `NewWSUsingChan*`) with `pkg/api/speak/v1/websocket` | Shipped |
| Flux TTS | `/v2/speak` REST and WebSocket | none | Not shipped. No first-class Flux TTS client or typed v2 Speak options are available; existing client URL overrides are not Flux TTS support |
| Voice Agent | `wss agent.deepgram.com/v1/agent/converse` | `pkg/client/agent` (`NewWSUsingChan*`) with `pkg/api/agent/v1/websocket` | Shipped; channel-based only |
| Text intelligence | `POST /v1/read` | `pkg/client/analyze` with `pkg/api/analyze/v1` | Shipped |
| Management API | `/v1/projects/...` and `/v1/models...` | `pkg/client/manage` with `pkg/api/manage/v1` | Shipped |
| Auth (grant token) | `POST /v1/auth/grant` | `pkg/client/auth` with `pkg/api/auth/v1` | Shipped |

When you add a surface, add its path to the appropriate `pkg/api/version` helper, update `constants.go` when the helper uses `APIPathMap`, add its options to `pkg/client/interfaces/<version>`, and include a runnable example under `examples/`.

## Prerequisites

- Go 1.19. The Go build, test, and lint workflows pin `go-version: "1.19"`. Newer toolchains compile and test the module (verified with Go 1.27.1), but `golangci-lint` v1.48.0 (the version `hack/check/tools` builds) panics while loading packages under Go 1.27, so run `make lint` with Go 1.19.
- PortAudio development headers. `pkg/audio/microphone` imports `github.com/gordonklaus/portaudio`, which needs `pkg-config --cflags portaudio-2.0` to succeed. Without it `go build ./...`, `go vet ./...`, and `go test ./...` fail on `pkg/audio/*` and on every microphone example. Install with `apt-get install -y portaudio19-dev pkg-config` (Debian and Ubuntu) or `brew install portaudio` (macOS), or run `make ensure-deps`, which also installs actionlint, shellcheck, jq, and the GitHub CLI.
- Docker, for `make mdlint` (it runs `ghcr.io/igorshubovych/markdownlint-cli` against `*.md`).

## Build, test, lint, format

Run the narrowest check first, then the CI command. Every command in this table was run on 2026-09-13 inside `golang:1.19-bullseye` and `golang:latest` (Go 1.27.1) containers with the repository mounted at `/work` and PortAudio installed; the exit codes are recorded in the pull request that added this file.

| Task | Command | Notes |
| --- | --- | --- |
| Download and verify modules | `go mod download && go mod verify` | CI runs `go mod tidy`, `go mod download`, `go mod verify` before tests |
| Build everything | `go build ./...` | Needs PortAudio headers (see above) |
| Vet | `go vet ./...` | Not run by CI; run it anyway |
| Unit tests, CI command | `go test -v -run Test_ ./...` | What `.github/workflows/tests-unit.yaml` runs on every pull request, with a 5 minute job timeout. Only functions named `Test_*` execute, so name unit tests with the `Test_` prefix |
| Unit tests, fastest | `go test ./tests/unit_test/...` | Deterministic, no network, about 1 second |
| One test | `go test ./tests/unit_test -run Test_PrerecordedDiarizeModel` | |
| Lint | `make lint` | Builds `golangci-lint` v1.48.0 from `hack/check/tools` and runs it with `.golangci.yaml`. It exits 0 even when the linter reports issues or panics, so read the output. At `v3.7.0` it reports 18 pre-existing findings (`deadcode`, `goconst`, `gocritic` hugeParam); fix only the ones in files you touch. In a container, pass `ROOT_DIR=/work` because the `Makefile` derives the root from `git rev-parse` |
| Format | `gofmt -l ./pkg ./tests ./examples` | Prints unformatted files; `gofmt` is also a `golangci-lint` linter here. As of 2026-09-13, `tests/unit_test/speak_websocket_metadata_test.go` is listed; format it only in a change that touches that file |
| Markdown lint | `make mdlint` | Rules in `.markdownlintrc` (line length off, fenced code blocks, `MD024` allows duplicate headings at different nesting) |
| Shell, YAML, Actions lint | `make shellcheck`, `make yamllint`, `make actionlint` | `make check` runs all five linters; `check-all.yaml` runs it on every push to `main` and `release-*` |

Use `go test ./...` deliberately. It includes `tests/daily_test`, which calls the live API when `DEEPGRAM_API_KEY` or `DEEPGRAM_ACCESS_TOKEN` is set and rewrites the fixtures under `tests/response_data`. `tests-daily.yaml` runs `go test -v -run TestDaily_ ./...` on a schedule at 09:00 UTC and opens a pull request with the refreshed fixtures.

## Run an example against the live API

Examples that call the Deepgram API resolve credentials from the environment when they construct a client with a `*WithDefaults` constructor, so nothing needs editing.

```bash
# The SDK resolves DEEPGRAM_API_KEY from the environment.
export DEEPGRAM_API_KEY="<your key>"

# Transcribe a hosted file with Nova-3 and print the JSON response.
go run ./examples/speech-to-text/rest/url

# Stream a live radio feed over the Nova WebSocket and print transcripts.
go run ./examples/speech-to-text/websocket/http_channel

# Synthesize speech with Aura and write a WAV file.
go run ./examples/text-to-speech/rest/file/hello-world
```

Microphone examples (`microphone_callback`, `microphone_channel`, `flux_callback`, `flux_channel`, `test`, `agent/websocket/simple`) capture audio through PortAudio, so they need an audio device and do not run in a container.

## Implementation conventions

- Keep public changes additive. Preserve deprecated aliases when you replace a public method or constant.
- Put version-specific protocol types in the matching `v1` or `v2` package. Do not reuse a model across REST and WebSocket protocols unless their wire formats are identical.
- Use `context.Context` for request and connection lifetime. Listen and Speak WebSocket clients ship both callback and channel variants; keep both in step when changing behavior. Voice Agent is channel-based only.
- Match JSON tags and optionality to the wire contract. Add a test for any field whose zero value is meaningful when the model is re-marshaled.
- Return errors to callers. Existing `pkg/api/manage/v1` methods (all but `invitations.go`) return `&resp, nil` after a failed request; that is a known defect. Do not copy it into new code and do not change it without a tracked issue.
- Do not log credentials, and do not change global logging (`klog`) or process flag behavior from ordinary client code.
- Every `.go` file starts with the MIT license header that the `goheader` linter checks. Use the current year or a year range accepted by `.golangci.yaml`:

  ```go
  // Copyright <year or year range> Deepgram SDK contributors. All Rights Reserved.
  // Use of this source code is governed by a MIT license that can be found in the LICENSE file.
  // SPDX-License-Identifier: MIT
  ```

- Imports are grouped with `goimports` and the local prefix `github.com/deepgram/deepgram-go-sdk`. `ioutil` is banned (`make lint` greps for it).
- Files in `tests/unit_test` use `package deepgram_test`; keep new tests in that package.
- Write "Flux STT" or "Flux TTS" in prose and comments; never bare "Flux". Identifiers such as `FluxTranscriptionOptions` and the `flux-general-en` model name stay as they are.
- Keep examples idiomatic and runnable, and update the affected examples and `README.md` with every API change.

## Example: add a Listen v1 option

For a new pre-recorded transcription query parameter:

1. Add the typed field and its `schema` tag to `pkg/client/interfaces/v1/types-prerecorded.go`.
2. Add or update the typed response model under `pkg/api/listen/v1/rest/interfaces` when the API returns new data.
3. Add a focused unit test proving the query parameter or response field reaches the wire contract.
4. Update the closest runnable example under `examples/speech-to-text/rest`.
5. Run the focused test, then `go test ./tests/unit_test/...`, then the CI command `go test -v -run Test_ ./...`.

## Release process

Releases are typically tags on `main`; patch releases for older majors are tagged from their matching `release-v[0-9]+` branch. There is no release workflow or release-please. The full process is in `.github/BRANCH_AND_RELEASE_PROCESS.md`.

1. `main` must stay releasable. Consumers pin a tag (`go get github.com/deepgram/deepgram-go-sdk/v3@v3.7.0`), never `main`.
2. A maintainer tags with semver and a `v` prefix (`git tag -m v3.8.0 v3.8.0 && git push upstream v3.8.0`), then creates a GitHub release from the tag. A breaking interface change bumps the major version and the module path (`/v4`) and gets a `release-v3` branch for patches.
3. `context7.yml` refreshes the Context7 index when the release is published.
4. Patch releases for an older major happen on the matching `release-v[0-9]+` branch with the same tag commands.

## Pull requests

- Keep diffs focused on one issue or behavior change.
- Add regression coverage for bug fixes, including exact REST or WebSocket wire behavior when the bug is a serialization or protocol mismatch.
- State the commands you ran and link the related issue. Use `Fixes #<issue>` only when the pull request fully resolves it.
- Follow `.github/PULL_REQUEST_TEMPLATE.md`. Commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`).

## Documentation

- API reference and product guides: <https://developers.deepgram.com/docs>
- Go build pages: <https://developers.deepgram.com/docs/speech-to-text/build/go>, <https://developers.deepgram.com/docs/speech-to-text/streaming/build/go>, <https://developers.deepgram.com/docs/speech-to-text/flux/build/go>, <https://developers.deepgram.com/docs/text-to-speech/build/go>, <https://developers.deepgram.com/docs/text-to-speech/streaming/build/go>, <https://developers.deepgram.com/docs/voice-agent/build/go>
- SDK feature matrix: <https://developers.deepgram.com/docs/sdks/sdk-features>
- Package docs: <https://pkg.go.dev/github.com/deepgram/deepgram-go-sdk/v3>
- Agent skills that teach this SDK: `.agents/skills/` (install with `npx skills add deepgram/deepgram-go-sdk`)

## Do not

- Do not run `go test ./...` with `DEEPGRAM_API_KEY` or `DEEPGRAM_ACCESS_TOKEN` set unless you intend to call the live API and refresh `tests/response_data`.
- Do not commit fixtures, keys, or `.env` files.
- Do not add new code to the deprecated `live`, `prerecorded`, or `rest` packages.
- Do not reformat files you are not otherwise changing.
- Do not add a `/v2/speak` (Flux TTS) client without also adding `pkg/api/version` constants, `v2` option types, unit tests, and an example.
