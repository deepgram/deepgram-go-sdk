---
trigger: always_on
description: "Select the repository skill that matches a Deepgram SDK product or workflow request."
---

When a request concerns Deepgram SDK product usage, API behavior, examples, or SDK-specific workflows:

1. Inspect `.agents/skills/` to identify the skill whose name and description match the request.
2. Read only that relevant skill before acting.
3. Do not load unrelated skills.

For repository-wide conventions, follow `AGENTS.md`.
