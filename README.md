# Claude Code Commands

Portable slash commands for managing multi-repo workspaces with Claude Code.

## Setup

1. Copy this folder to your workspace at `.claude/commands/`
2. Copy `config.yaml.example` to `config.yaml`
3. Edit `config.yaml` with your base path and repositories

```yaml
base_path: ~/code/my-workspace

repos:
  - name: my-app
    group: apps
    aliases: [app]
```

## Commands

| Command | Description |
|---------|-------------|
| `/super <repo>` | Start brainstorming session with full context |
| `/find-tasks <repo>` | Suggest 3-5 high-priority tasks |
| `/run-tests <repo>` | Run lint, type-check, build, and tests |
| `/commit-progress <repo>` | Draft and commit changes |
| `/update-docs <repo>` | Update CLAUDE.md, README, docs |
| `/review-project <repo>` | Technical review to docs/tech-review.md |
| `/add-repo <url>` | Clone repo and add to config |
| `/dev-rules` | Remind Claude of workspace rules (includes pwd check) |

All repo commands support fuzzy matching via aliases (e.g., `/run-tests app`).

## Configuration

### config.yaml

```yaml
base_path: ~/code/workspace    # Root directory for all repos

repos:
  - name: my-nextjs-app        # Directory name
    group: apps                # 'apps' or 'devops'
    aliases: [app, next]       # Fuzzy match shortcuts
    work_dir: src              # Optional: subdirectory for commands
    test_cmd: pnpm test        # Optional: custom test command
```

### Groups

- **devops**: Infrastructure repos (Pulumi, Terraform, etc.)
- **apps**: Application repos (Next.js, Go, Python, etc.)

## Files

| File | Purpose |
|------|---------|
| `config.yaml.example` | Template (checked in) |
| `config.yaml` | Your config (gitignored) |
| `_shared-repo-logic.md` | Common patterns for repo commands |

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Repos should have a `CLAUDE.md` for best results

## Optional: Superpowers Skills

The `/super` command uses skills from the Superpowers marketplace plugin. To install:

```bash
claude mcp add superpowers-marketplace -- npx -y @anthropic-ai/superpowers-marketplace@latest
```

This adds skills like:
- `superpowers:brainstorming` - Structured creative exploration
- `superpowers:writing-plans` - Implementation planning
- `superpowers:systematic-debugging` - Bug investigation
- `superpowers:test-driven-development` - TDD workflow
- `superpowers:verification-before-completion` - Pre-commit checks

These are optional but recommended. Without them, `/super` will still work but won't invoke the brainstorming skill.

## Commit Rules

These commands enforce:
- No Claude/Anthropic attribution in commits
- Imperative mood ("Add feature" not "Added feature")
- Summary under 72 characters
