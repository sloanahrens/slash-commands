# Simplify Repo Resolution

*Design document - 2026-01-03*

## Problem

Aliases and fuzzy matching cause Claude to guess paths incorrectly. When running `/make-test @fractals-nextjs/`, Claude used `devbot make fractals` (worked), then manually constructed `~/code/fractals-nextjs` (wrong - should be `~/code/mono-claude/fractals-nextjs`).

**Root causes:**
1. Aliases add indirection that Claude misremembers
2. Fuzzy matching in devbot means Claude doesn't ask for clarification
3. Path construction logic scattered across 14+ files
4. Examples in docs don't match actual config values

## Solution

Remove aliases entirely. Require exact repo names. Add `devbot path` command. Consolidate all resolution logic into `_shared-repo-logic.md`.

## Design

### 1. Config Simplification

Remove `aliases` field from all repo entries in `config.yaml`:

```yaml
# BEFORE
repos:
  - name: fractals-nextjs
    group: apps
    aliases: [fractals]
    language: typescript

# AFTER
repos:
  - name: fractals-nextjs
    group: apps
    language: typescript
```

**Constraint:** The `name` field MUST exactly match the directory name under `code_path`.

### 2. devbot Changes

**Remove:**
- Alias matching in `FindRepoByName()`
- Fuzzy (contains) matching in `FindRepoByName()`

**Add:**
- `devbot path <repo>` command

**New `FindRepoByName()` logic:**

```go
func FindRepoByName(name string) *RepoConfig {
    cfg, _ := LoadConfig()
    if cfg == nil {
        return nil
    }

    for i := range cfg.Repos {
        if cfg.Repos[i].Name == name {  // Exact match only
            return &cfg.Repos[i]
        }
    }
    return nil
}
```

**New `devbot path` command:**

```bash
$ devbot path fractals-nextjs
/Users/sloan/code/mono-claude/fractals-nextjs

$ devbot path fractals
Repository 'fractals' not found. Did you mean:
  fractals-nextjs

$ devbot path nonexistent
Repository 'nonexistent' not found.
```

The "did you mean" is output only - no auto-resolution.

### 3. Slash Command Consolidation

**`_shared-repo-logic.md` becomes single source of truth.**

Individual commands change from:
```markdown
**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.
```

To:
```markdown
**Arguments**: `$ARGUMENTS` - Repo name. See `_shared-repo-logic.md`.
```

**New resolution logic in _shared-repo-logic.md:**

```markdown
## Repo Resolution

### When user provides `@directory/`
1. Extract directory name from path
2. Run `devbot path <name>` to get full path
3. If not found, show devbot's suggestion and ask user

### When user provides plain name
1. Run `devbot path <name>`
2. If found, use that path
3. If not found, show suggestion and ask user to confirm

### Getting the full path
ALWAYS use devbot to get paths:
```bash
devbot path fractals-nextjs
```

NEVER construct paths manually.
```

### 4. CLAUDE.md Cleanup

**Remove from `~/.claude/CLAUDE.md`:**
- `Alias` column from Repository Registry table
- Line: "All repo commands support fuzzy matching"
- Any mention of aliases

**Keep:**
- Repo table (without aliases) for quick reference
- Gotchas section

## Files to Modify

| File | Changes |
|------|---------|
| `config.yaml` | Remove `aliases` from all repos |
| `devbot/internal/workspace/wscfg.go` | Exact match only |
| `devbot/cmd/devbot/main.go` | Add `path` command |
| `_shared-repo-logic.md` | Rewrite repo resolution |
| `~/.claude/CLAUDE.md` | Remove alias references |
| 10 command `.md` files | Reference shared logic |

## Testing

```bash
devbot path fractals-nextjs    # /Users/sloan/.../fractals-nextjs
devbot path fractals           # "not found, did you mean: fractals-nextjs"
devbot path nonexistent        # "not found"
```

## Migration

No user action needed. Old aliases stop working. devbot suggests full names when partial names used.
