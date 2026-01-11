# CLAUDE.md

Global Claude Code configuration. Each repo also has its own CLAUDE.md - read that first.

## Setup

```bash
# Clone this repo as ~/.claude
git clone https://github.com/sloanahrens/slash-commands.git ~/.claude
cd ~/.claude
cp config.yaml.example config.yaml  # Edit with your workspace path
make -C devbot install               # Install devbot CLI
```

Or run `/setup-workspace` in Claude Code after cloning.

## Critical Rules

- **NO Claude/Anthropic attribution** in commits
- **Read repo CLAUDE.md first** - each has specific guidance
- **Use exact repo names** from config.yaml
- **Use devbot** - prefer over manual git/file operations
- **Simple bash only** - no `&&`, `$()`, `;`, or `git -C` (blocked by hookify)

## Tool Selection Guide

**STOP before running bash.** Use the right tool:

| Need | Tool | Example |
|------|------|---------|
| Git status/diff/branch | `devbot` | `devbot diff my-repo --full` |
| Read file content | `Read` tool | Read tool on any file path |
| Search files | `Grep`/`Glob` | Grep for pattern, Glob for filenames |
| File operations | `Read`/`Edit`/`Write` | Never use cat/sed/awk |

**Never improvise bash commands.** If devbot doesn't have a command for it, use Claude Code's built-in tools (Read, Grep, Glob, Edit, Write).

## Bash Patterns

```bash
# Get path, cd, then git (for commands without devbot wrappers)
devbot path my-repo        # → /full/path/to/my-repo
cd /full/path/to/my-repo
git commit                 # Commands like commit/push need cd first
```

| Use This | Not This |
|----------|----------|
| `devbot status <repo>` | `git status` |
| `devbot diff <repo>` | `git diff` |
| `devbot branch <repo>` | `git branch -vv` |
| `devbot log <repo>` | `git log` |
| `devbot show <repo> [ref]` | `git show` |
| `devbot fetch <repo>` | `git fetch` |
| `devbot switch <repo> <branch>` | `git switch/checkout` |
| `devbot check <repo>` | `npm test && npm run lint` |
| `devbot last-commit <repo> [file]` | `git log -1 --format="%ar"` |

## Running Commands in Repos (NO cd &&)

**Preferred:** Use `devbot exec` for any command in a repo directory:

```bash
devbot exec my-app npm run build    # Uses work_dir
devbot exec my-app/subdir go test   # Explicit subdir
```

**Fallback patterns** (when devbot exec isn't suitable):

| Tool | Pattern | Example |
|------|---------|---------|
| npm | `npm run <cmd> --prefix <path>` | `npm run build --prefix /path/to/app` |
| make | `make -C <path> <target>` | `make -C /path/to/app build` |

**Sequential commands:** Run each command separately, one tool call at a time. Do NOT combine with `&&` or `;`.

## Slash Commands

All require exact repo names. Run `/list-commands` for full list.

| Command | Description |
|---------|-------------|
| `/super <repo>` | Brainstorming with context |
| `/run-tests <repo>` | Lint, type-check, build, test |
| `/yes-commit <repo>` | Draft and commit (no AI attribution) |
| `/push <repo>` | Push to origin |
| `/status [repo]` | Repository status |
| `/update-docs <repo>` | Update documentation |

## devbot CLI

**NAME commands:** `path`, `status`, `diff`, `branch`, `log`, `show`, `fetch`, `switch`, `check`, `make`, `todos`, `last-commit`, `config`, `deps`, `remote`, `worktrees`, `pulumi`, `deploy`

**PATH commands:** `tree`, `stats`, `detect` (use `devbot path` first)

**Execution helpers:**
- `exec <repo>[/subdir] <cmd...>` - Run command in repo directory (respects work_dir)
- `port <port> [--kill]` - Check/kill process on port
- `prereq <repo>[/subdir]` - Validate tools, deps, and env vars before work

**Git wrappers** (faster, auto-approved):
- `log <repo>` - git log (default: --oneline -20)
- `show <repo> [ref]` - git show (default: HEAD)
- `fetch <repo>` - git fetch --all --prune
- `switch <repo> <branch>` - git switch

**CRITICAL:** `devbot pulumi <repo>` - **MANDATORY before any Pulumi operation**

## Pulumi (CRITICAL)

### MANDATORY: Run `devbot pulumi <repo>` BEFORE any Pulumi command

This prevents destructive operations by showing existing infrastructure state.

### Forbidden Commands (unless devbot pulumi shows NO infrastructure)

| Command | Why Dangerous |
|---------|--------------|
| `pulumi stack init` | Orphans existing infrastructure |
| `pulumi destroy` | Deletes all resources |
| `pulumi stack rm` | Loses state permanently |

## Key Skills

| Skill | When |
|-------|------|
| `superpowers:brainstorming` | Before creative work |
| `superpowers:systematic-debugging` | Bug investigation |
| `superpowers:verification-before-completion` | Before claiming done |

## Files

| Location | Purpose |
|----------|---------|
| `~/.claude/CLAUDE.md` | This file - global instructions |
| `~/.claude/settings.json` | Permissions + plugins |
| `~/.claude/config.yaml` | Workspace config (gitignored) |
| `~/.claude/hookify.*.md` | Hookify rules |
| `~/.claude/commands/` | Slash commands |
| `~/.claude/devbot/` | CLI tool source |
| `<repo>/CLAUDE.md` | Repo-specific guidance |

## Local Model

Use Qwen for simple tasks (commit messages, explanations). Prefix output with `[qwen]`.

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="...", max_tokens=100
)
```

## Usage Monitoring

```bash
npx ccusage@latest          # Daily token usage
npx ccusage@latest blocks   # 5-hour rate limit windows
```
