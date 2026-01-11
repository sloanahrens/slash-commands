# CLAUDE.md

Global Claude Code configuration. Each repo has its own CLAUDE.md - read that first.

## Setup

```bash
git clone https://github.com/sloanahrens/slash-commands.git ~/code/slash-commands
/setup-workspace  # Handles config, devbot, symlinks, plugins
```

## Critical Rules

- **NO Claude/Anthropic attribution** in commits
- **Read repo CLAUDE.md first** - each has specific guidance
- **Use exact repo names** from config.yaml
- **Use devbot** - prefer over manual git/file operations
- **Simple bash only** - no `&&`, `$()`, `;`, or `git -C` (blocked by hookify)

## Tool Selection

| Need | Tool |
|------|------|
| Git operations | `devbot` commands |
| Read files | `Read` tool |
| Search files | `Grep`/`Glob` |
| File operations | `Read`/`Edit`/`Write` (never cat/sed/awk) |

## Bash Patterns

```bash
devbot path my-repo        # Get path first
cd /full/path/to/my-repo   # Then cd
git commit                 # Then git command
```

**Use devbot exec for running commands:**
```bash
devbot exec my-app npm run build    # Respects work_dir
devbot exec my-app/subdir go test   # Explicit subdir
```

## Slash Commands

Run `/list-commands` for full list. Key commands:

| Command | Description |
|---------|-------------|
| `/super <repo>` | Brainstorming with context |
| `/run-tests <repo>` | Lint, type-check, build, test |
| `/yes-commit <repo>` | Draft and commit |
| `/status [repo]` | Repository status |
| `/prime <repo>` | Load patterns and notes |
| `/capture-hindsight` | Save failure as note |

## devbot CLI

See `slash-commands/devbot/README.md` for full reference.

**Quick reference:**
```bash
devbot status              # All repos
devbot check <repo>        # Quality checks
devbot exec <repo> <cmd>   # Run command in repo
devbot pulumi <repo>       # MUST run before any pulumi command
```

## Pulumi (CRITICAL)

**MANDATORY:** Run `devbot pulumi <repo>` BEFORE any Pulumi command.

**Forbidden** (unless devbot pulumi shows NO infrastructure):
- `pulumi stack init` - Orphans existing infrastructure
- `pulumi destroy` - Deletes all resources
- `pulumi stack rm` - Loses state permanently

## Key Skills

| Skill | When |
|-------|------|
| `superpowers:brainstorming` | Before creative work |
| `superpowers:systematic-debugging` | Bug investigation |
| `superpowers:verification-before-completion` | Before claiming done |

## Files

| Location | Purpose |
|----------|---------|
| `~/.claude/settings.json` | Global permissions + plugins |
| `~/.claude/hookify.*.md` | Global hookify rules |
| `~/.claude/notes/` | Local notes (hindsight, sessions) |
| `<repo>/CLAUDE.md` | Repo-specific guidance |
| `slash-commands/docs/patterns/` | Versioned knowledge patterns |

## Local Model

Use Qwen for simple tasks. Prefix output with `[local]`.

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
