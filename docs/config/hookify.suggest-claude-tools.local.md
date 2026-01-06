---
name: suggest-claude-tools
enabled: true
event: bash
pattern: ^(cat|head|tail)\s+[^|]|^grep\s+|^find\s+.*-name|^(sed|awk)\s+
action: warn
---

**Use Claude Code tools instead**

| Shell command | Claude tool |
|---------------|-------------|
| `cat file` | Read tool |
| `head -n file` | Read tool with `limit` |
| `tail -n file` | Read tool with `offset` |
| `grep pattern dir` | Grep tool |
| `find -name "*.ts"` | Glob tool |
| `sed` / `awk` | Edit tool |

Benefits: auto-approved, better context, faster.
