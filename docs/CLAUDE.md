# CLAUDE.md

Global Claude Code configuration. Each repo also has its own CLAUDE.md - read that first.

## Setup

```bash
git clone git@bitbucket.org:sloanahrens/app-builder-commands.git ~/code/slash-commands
/setup-workspace  # Handles config, devbot, symlinks, plugins
```

## Critical Rules

- **NO Claude/Anthropic attribution** in commits
- **Read repo CLAUDE.md first** - each has specific guidance
- **Use exact repo names** from config.yaml
- **Use devbot** - prefer over manual git/file operations
- **Simple bash only** - no `&&`, `$()`, `;`, or `git -C`

## Bash Patterns

```bash
# Get path, cd, then git
devbot path my-repo        # → /full/path/to/my-repo
cd /full/path/to/my-repo
git log                    # Regular git commands
```

| Use This | Not This |
|----------|----------|
| `devbot status <repo>` | `git status` |
| `devbot diff <repo>` | `git diff` |
| `devbot branch <repo>` | `git branch -vv` |
| `devbot check <repo>` | `npm test && npm run lint` |
| `devbot last-commit <repo> [file]` | `git log -1 --format="%ar"` |

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

## Repository Registry

Defined in `~/code/slash-commands/config.yaml`.

| Repo | Stack | Notes |
|------|-------|-------|
| `atap-automation2` | Next.js + Playwright | **NO SANDBOX** - workdir: `nextapp/` |
| `fractals-nextjs` | Next.js + Canvas | Mandelbrot visualizer |
| `mango` | Go + Next.js + DuckDB | **CGO required** |
| `devops-pulumi-ts` | Pulumi + TypeScript | GCP Cloud Run |
| `slash-commands` | Go + Markdown | This workspace's tools |

## devbot CLI

**NAME commands:** `path`, `status`, `diff`, `branch`, `check`, `make`, `todos`, `last-commit`

**PATH commands:** `tree`, `stats` (use `devbot path` first)

```bash
devbot path fractals-nextjs   # Get path
devbot tree /full/path        # Then use path
```

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
| `<repo>/CLAUDE.md` | Repo-specific guidance |
| `slash-commands/docs/` | Reference copies of config |

## Local Model

Use Qwen for simple tasks (commit messages, explanations). Prefix output with `[qwen]`.

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="...", max_tokens=100
)
```
