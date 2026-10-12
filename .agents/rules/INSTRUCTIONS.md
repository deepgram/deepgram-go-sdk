---
trigger: always_on
description: "Select applicable repository skills for Deepgram SDK product requests."
---

When a request clearly matches one or more skills in `.agents/skills/`:

1. Inspect `.agents/skills/` to identify every skill whose name and description match the request.
2. Read each applicable skill before acting.
3. Do not load unrelated skills.

If no skill matches, follow `AGENTS.md` and the repository guidance.
