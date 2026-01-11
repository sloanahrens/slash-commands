---
tags: [hookify, bash, safety, blocked]
repos: [all]
created: 2026-01-11
updated: 2026-01-11
---

# Hookify blocked commands and workarounds

## Problem

Hookify rules block certain bash patterns to prevent mistakes. When a command is blocked, you need to know the correct alternative.

## Blocked patterns

### Compound commands (`&&`, `;`, `$()`)

```bash
# Blocked
cd /path/to/repo && npm test
git add . && git commit -m "msg"
RESULT=$(some-command)
```

**Workaround:** Run commands sequentially in separate Bash calls, or use `devbot exec`:

```bash
devbot exec my-repo npm test
```

For git operations, run each command separately.

### git -C flag

```bash
# Blocked
git -C /path/to/repo status
```

**Workaround:** Use devbot wrappers:

```bash
devbot status my-repo
devbot diff my-repo
devbot branch my-repo
```

### Direct file operations via bash

```bash
# Discouraged (use Claude Code tools instead)
cat file.txt           # Use Read tool
grep pattern files     # Use Grep tool
sed -i 's/x/y/' file   # Use Edit tool
echo "content" > file  # Use Write tool
```

**Workaround:** Use the native Claude Code tools which are faster and safer.

## Why these rules exist

1. **Compound commands**: Partial execution on failure leaves ambiguous state
2. **git -C**: Easy to operate on wrong repo; devbot enforces config.yaml names
3. **File operations**: Native tools have better error handling and don't require shell escaping

## Checking what's blocked

Hookify rules are defined in:
- `~/.claude/hookify.block-dangerous.local.md`
- `~/.claude/hookify.suggest-devbot.local.md`
- `~/.claude/hookify.suggest-claude-tools.local.md`

## When you hit a block

1. Read the hookify error message — it usually suggests the alternative
2. Check this pattern file for workarounds
3. If the block seems wrong, discuss with user before requesting override

## Related

- [bash-execution.md](./bash-execution.md) — devbot exec patterns
- devbot README — Full list of devbot commands
