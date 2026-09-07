# Deepgram Go SDK

## Repository Purpose

This is the official Go SDK for Deepgram APIs. The module path is
`github.com/deepgram/deepgram-go-sdk/v3`, and the supported Go version is 1.19
or newer.

Read `README.md` and `.github/CONTRIBUTING.md` before changing public behavior.
Do not hardcode API keys or access tokens; examples and tests should use the
SDK's environment-backed client defaults.

## Layout

- `pkg/api/<product>/<version>`: public API facades, request/response models,
  routers, and callbacks.
- `pkg/client/<product>/<version>`: REST and WebSocket transport clients.
- `pkg/client/interfaces`: shared public option types. Versioned types live in
  `v1` and `v2` and are re-exported here when appropriate.
- `pkg/client/common`: shared REST and WebSocket behavior.
- `pkg/api/version`: endpoint URL and query-parameter construction.
- `examples`: runnable, product-specific usage.
- `tests/unit_test`: deterministic unit tests. `tests/daily_test` uses live
  service responses and may refresh response fixtures.

Use `listen` for Speech-to-Text, `speak` for Text-to-Speech, `analyze` for
text intelligence, `agent` for Voice Agent, and `manage` or `auth` for
management and token operations. Legacy `live`, `prerecorded`, and `rest`
packages are deprecated; do not use them for new code.

## Implementation Conventions

- Keep public changes additive whenever possible. Preserve deprecated aliases
  when replacing a public method or constant.
- Put version-specific protocol types in the matching `v1` or `v2` package.
  Do not reuse a model across REST and WebSocket protocols unless their wire
  formats are identical.
- Use `context.Context` for request and connection lifetime. WebSocket clients
  expose separate callback and channel variants; preserve both when changing
  their behavior.
- Match JSON tags and optionality to the wire contract. Test fields whose zero
  value is meaningful when the model is re-marshaled.
- Return errors to callers. Do not log credentials, and do not change global
  logging or process flag behavior from ordinary client code.
- Keep examples idiomatic and runnable. Update affected examples and public
  documentation with every API change.

## Validation

Run the narrowest relevant check first:

```bash
go test ./tests/unit_test
go test -v -run Test_ ./...
go vet ./...
go mod verify
```

CI uses Go 1.19 and runs `go test -v -run Test_ ./...`. Name unit tests with
the `Test_` prefix so CI executes them. Run `make lint` with Go 1.19 for Go
changes. Run `make mdlint` for Markdown changes.

Use `go test ./...` deliberately: it includes `tests/daily_test`, which makes
live service calls and can update tracked response fixtures.

## Pull Requests

- Keep diffs focused on one issue or behavior change.
- Add regression coverage for bug fixes, including exact REST or WebSocket
  wire behavior when the bug is a serialization or protocol mismatch.
- State the commands run and link the related issue. Use `Fixes #<issue>` only
  when the PR fully resolves it.
