---
name: suggest-claude-tools
enabled: true
event: bash
pattern: ^(cat|head|tail)\s+[^|]|^(grep|rg)\s+|^find\s+.*-name|^(sed|awk)\s+|\|\s*(head|tail|grep|rg|wc|sort|awk|sed)
action: warn
---

**Use Claude Code tools instead**

## File Operations

| Shell command | Claude tool |
|---------------|-------------|
| `cat file` | Read tool |
| `head -n file` | Read tool with `limit` |
| `tail -n file` | Read tool with `offset` |
| `grep pattern dir` | Grep tool |
| `rg pattern` | Grep tool |
| `find -name "*.ts"` | Glob tool |
| `sed` / `awk` | Edit tool |

## Piped Commands

Avoid piping to filtering commands:
- `cmd | head` → run cmd, use Read tool
- `cmd | grep` → run cmd, use Grep tool
- `cmd | wc/sort` → run cmd, process output in response

Benefits: auto-approved, better context, faster.

The full output is often useful for diagnosis anyway - Claude Code auto-truncates if needed.
