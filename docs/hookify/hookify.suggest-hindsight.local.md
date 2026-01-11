---
name: suggest-hindsight
enabled: true
event: bash
pattern: error:|Error:|ERROR:|failed|Failed|FAILED|exit code [1-9]|command not found|No such file|Permission denied
action: warn
---

**Consider capturing this as hindsight**

An error or failure was detected in the output. If this required troubleshooting or multiple attempts to resolve, consider running:

```
/capture-hindsight
```

This saves the issue and solution for future sessions, preventing repeated mistakes.

**When to capture:**
- You had to try multiple approaches
- The error wasn't obvious from the message
- The fix involved a non-obvious pattern
- This could happen again in similar contexts
