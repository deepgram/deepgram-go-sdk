# Reusable Agent Configurations Example

Exercises the [reusable agent configuration](https://developers.deepgram.com/docs/reusable-agent-configurations) management endpoints:

- `POST/GET /v1/projects/{project_id}/agents` — create and list agent configurations
- `GET/PUT/DELETE /v1/projects/{project_id}/agents/{agent_id}` — get, update metadata, delete
- `POST/GET /v1/projects/{project_id}/agent-variables` — create and list template variables
- `GET/PATCH/DELETE /v1/projects/{project_id}/agent-variables/{variable_id}` — get, update value, delete

A reusable agent configuration saves the `agent` block of a Voice Agent Settings message once and returns an `agent_uuid` you can pass in place of the full agent object on every connection. Template variables (`DG_<VARIABLE_NAME>`) let one configuration serve many sessions — the server substitutes their values at connection time, and management endpoints always return configs uninterpolated.

The example creates a variable and an agent configuration, exercises list/get/update on both, then deletes everything it created.

## Prerequisites

- Go 1.19+
- A Deepgram API key with manage scope for the project

## Usage

```bash
DEEPGRAM_API_KEY=<your-key> go run main.go
```

## Notes

- `Config` is a JSON **string** (the `agent` block of a Settings message), not a nested object, in **both** directions: `AgentCreateRequest.Config` sends a string and `AgentConfiguration.Config` reads the same string back verbatim. The round trip is symmetric.
- Template variables are referenced **bare** (unquoted) inside the config string — `"prompt": DG_SYSTEM_PROMPT` — and each one substitutes a whole JSON value (string, number, boolean, object, or array) at connection time. `<VARIABLE_NAME>` uses uppercase alphanumerics, underscores, or hyphens.
- Both create endpoints answer with only the new UUID (`agent_uuid` / `agent_variable_uuid`), so `AgentID` and `VariableID` are the only populated fields on a create result. Call the matching get to read the rest.
- Both list endpoints answer with a **bare JSON array** at the top level rather than an object wrapping one. The SDK decodes that array into `AgentsList.Agents` and `AgentVariablesList.Variables`, so the accessor matches every other list in the manage API.
- The config of an existing agent is immutable: `UpdateAgentMetadata` (PUT) changes only metadata. To change the config, create a new agent and migrate traffic to its UUID.
- Metadata values are arbitrary JSON, not just strings, so `Metadata` is a `map[string]interface{}` on both the request and the response.
- `UpdateAgentMetadata` **replaces** the entire metadata object: any key omitted from the request is deleted, silently and with a `200`. Call `GetAgent` first and resend every key you mean to keep.
- `AgentVariableCreateRequest.IsSensitive` is optional on the wire — the API defaults it to `false` — but `false` is the only value it accepts today, so the SDK always sends it and the Go zero value is the value you want.
- `UpdateAgentMetadata`, `UpdateAgentVariable`, and both deletes answer `200` with an **empty body**, so those methods return only an `error`. Read the stored state back with the matching get.
- Deleting an agent configuration that live sessions still reference can cause a production outage — migrate first, then delete.

## Wire contract

Verified against the live API on 2026-09-23. `openapi.yaml` and the published API reference disagree with production on the id field names, on the list wrapping, and on the type of `config`; the shapes below are what the service actually returns.

| Call | Response |
| --- | --- |
| `POST /agents` | `{"agent_uuid": "..."}` |
| `GET /agents/{id}` | `{"agent_uuid", "member_id", "api_version", "config" (JSON string), "metadata"}` |
| `GET /agents` | bare JSON array of the get shape |
| `PUT /agents/{id}` | `200`, empty body |
| `DELETE /agents/{id}` | `200`, empty body |
| `POST /agent-variables` | `{"agent_variable_uuid": "..."}` |
| `GET /agent-variables/{id}` | `{"agent_variable_uuid", "member_id", "api_version", "key", "value", "is_sensitive"}` |
| `GET /agent-variables` | bare JSON array of the get shape |
| `PATCH /agent-variables/{id}` | `200`, empty body |
| `DELETE /agent-variables/{id}` | `200`, empty body |

Neither read shape carries `created_at` or `updated_at`.
