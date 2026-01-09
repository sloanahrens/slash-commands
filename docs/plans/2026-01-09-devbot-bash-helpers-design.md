# devbot Bash Helper Commands Design

**Date:** 2026-01-09
**Status:** Draft

## Problem

Claude wastes time and tokens iterating on bash commands that fail. Memory analysis revealed three main patterns:

1. **CLI tool orchestration** - Check tool exists → run → fail → discover missing config → retry
2. **Port conflicts** - Finding/killing processes requires `lsof -i` which is error-prone
3. **Directory-scoped commands** - Hookify blocks `cd && cmd`, but `--prefix` only works for npm/make

## Proposed Commands

### 1. `devbot exec <repo>[/subdir] <command...>`

Run any command in a repo's directory with smart directory resolution.

**Directory resolution order:**
1. If `/subdir` specified: `{repo_path}/{subdir}`
2. If `work_dir` in config.yaml: `{repo_path}/{work_dir}`
3. Otherwise: `{repo_path}`

```bash
# Uses work_dir from config (nextapp for atap-automation2)
devbot exec atap-automation2 npm run build
# → Runs in /Users/sloan/code/mono-claude/atap-automation2/nextapp

# Explicit subdir for monorepos
devbot exec mango/go-api go test ./...
# → Runs in /Users/sloan/code/mono-claude/mango/go-api

devbot exec mango/nextapp npm run build
# → Runs in /Users/sloan/code/mono-claude/mango/nextapp

# Subdir within slash-commands
devbot exec slash-commands/devbot make test
# → Runs in /Users/sloan/code/mono-claude/slash-commands/devbot

# Root directory (override work_dir)
devbot exec atap-automation2/ docker build .
# → Trailing slash means repo root, ignores work_dir
```

**Implementation:** ~100 lines. Parse repo/subdir → resolve path → chdir → exec.Command → stream output.

### 2. `devbot port <port> [--kill]`

Check/kill processes on ports.

```bash
devbot port 3000           # Show what's running
devbot port 3000 --kill    # Kill it
```

**Implementation:** ~80 lines. Parse `lsof -i :<port>`, format output, optionally kill.

### 3. `devbot prereq <repo>[/subdir]`

Validate prerequisites before work. Uses same directory resolution as `exec`.

```bash
devbot prereq atap-automation2
# → Checks nextapp/ (from work_dir)
# Shows: tools installed, deps installed, env vars set/missing
```

**Implementation:** ~200 lines. Use detect package, compare .env to .env.example.

## Config.yaml Integration

The existing `work_dir` field already supports this pattern:

```yaml
repos:
  - name: atap-automation2
    work_dir: nextapp  # Default execution directory

  - name: mango
    # No work_dir - must specify subdir: mango/go-api or mango/nextapp

  - name: slash-commands
    # No work_dir - root is default, use slash-commands/devbot for CLI
```

**Optional enhancement:** Add `subprojects` field for monorepos:

```yaml
  - name: mango
    subprojects:
      - name: go-api
        language: go
      - name: nextapp
        language: typescript
```

This would enable `devbot prereq mango` to check ALL subprojects.

## Documentation Updates Required

### 1. Update `~/.claude/CLAUDE.md`

Add to devbot CLI section:

```markdown
## devbot CLI

**NAME commands:** `path`, `status`, `diff`, `branch`, `log`, `show`, `fetch`, `switch`, `check`, `make`, `todos`, `last-commit`, `config`, `deps`, `remote`, `worktrees`, `pulumi`, `deploy`

**NEW: Execution helpers:**
- `exec <repo>[/subdir] <cmd...>` - Run command in repo directory
- `port <port> [--kill]` - Check/kill process on port
- `prereq <repo>` - Validate prerequisites

**Exec examples:**
| Command | Runs in |
|---------|---------|
| `devbot exec atap-automation2 npm test` | `.../atap-automation2/nextapp` |
| `devbot exec mango/go-api go build` | `.../mango/go-api` |
| `devbot exec slash-commands/devbot make` | `.../slash-commands/devbot` |
```

### 2. Update Slash Commands

**`_shared-repo-logic.md`** - Add exec pattern:

```markdown
### Running Commands in Repos

Use `devbot exec` instead of `cd && command`:

```bash
# Instead of: cd /path/to/repo && npm run build
devbot exec <repo-name> npm run build

# For monorepo subprojects
devbot exec <repo-name>/<subdir> <command>
```
```

**`run-tests.md`** - Replace complex bash with devbot exec:

```bash
# Instead of compound commands
devbot exec mango/go-api go test ./...
devbot exec mango/nextapp npm test
```

### 3. Update Hookify Rules

**`hookify.suggest-prefix.local.md`** - Add devbot exec suggestion:

```markdown
**Use devbot exec instead of cd**

| Instead of | Use |
|------------|-----|
| `cd /path; npm run build` | `devbot exec <repo> npm run build` |
| `cd /path/sub; make` | `devbot exec <repo>/sub make` |
```

## What NOT to Add

- More git commands (already well covered)
- Deployment wrappers (too repo-specific)
- Complex environment management (prereq covers validation)

## Implementation Priority

1. **`exec`** - Highest value, enables all other patterns (~100 lines)
2. **`port`** - Simple, frequent need (~80 lines)
3. **`prereq`** - Most complex, catches issues early (~200 lines)
4. **Documentation** - Update CLAUDE.md, slash commands, hookify rules

## Architecture Notes

Each command follows existing devbot patterns:
- Package in `internal/<command>/`
- Structured result types
- Uses `workspace` package for repo resolution
- Command wiring in `cmd/devbot/main.go`

The `exec` command can reuse:
- `workspace.LoadConfig()` for repo resolution
- `workspace.RepoInfo.Path` for base path
- New `workspace.RepoInfo.WorkDir` field (already exists in config)
