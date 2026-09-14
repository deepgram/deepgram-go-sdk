---
trigger: always_on
---

MANDATORY: Read `.agents/skills/**` before responding to ANY user message.

Do not act on the user's request until you have read `.agents/skills/**`.
It contains the rules that determine your first action based on what the user asked.
Skipping it WILL cause you to take the wrong action.

There are no additional instructions in this file. All detailed instructions are in `.agents/skills/**`.
