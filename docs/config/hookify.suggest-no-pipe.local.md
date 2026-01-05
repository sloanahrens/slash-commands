---
name: suggest-no-pipe
enabled: true
event: bash
pattern: \|\s*(head|tail|grep|wc|sort)
action: warn
---

**Avoid piping to head/tail/grep/wc/sort**

Piped commands may require manual approval. Instead:

- Run the command without the pipe
- Use Claude Code's built-in tools (Grep, Read) for filtering
- If output is too long, it will be automatically truncated

The full output is often useful for diagnosis anyway.
