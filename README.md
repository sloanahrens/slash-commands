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
| `/status [repo]` | Show status overview of all or one repo |
| `/sync [repo]` | Pull latest changes for all or one repo |
| `/switch <repo>` | Quick context switch to a repo |
| `/dev-rules` | Remind Claude of workspace rules |
| `/setup-plugins` | Install all recommended plugins |

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
    language: typescript       # Optional: typescript | go | python | rust | shell
    commands:                  # Optional: override default commands
      test: pnpm test
      lint: pnpm lint
      build: pnpm build
```

### Language Detection

If `language` is not specified, it's auto-detected from files:
- `package.json` → typescript
- `go.mod` → go
- `pyproject.toml` → python
- `Cargo.toml` → rust

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

## Recommended Plugins

Run `/setup-plugins` to install all recommended plugins, or install manually:

### Add Marketplaces

```bash
claude plugin marketplace add obra/superpowers-marketplace
claude plugin marketplace add anthropics/claude-plugins-official
```

### Core Plugins (Superpowers Marketplace)

```bash
claude plugin install superpowers@superpowers-marketplace          # TDD, debugging, brainstorming
claude plugin install elements-of-style@superpowers-marketplace    # Writing guidance
claude plugin install episodic-memory@superpowers-marketplace      # Memory across sessions
claude plugin install double-shot-latte@superpowers-marketplace    # Auto-continue
claude plugin install superpowers-developing-for-claude-code@superpowers-marketplace  # Plugin dev
claude plugin install superpowers-lab@superpowers-marketplace      # Experimental (vim, tmux)
claude plugin install superpowers-chrome@superpowers-marketplace   # Chrome DevTools (BETA)
```

### Official Plugins

```bash
claude plugin install frontend-design@claude-plugins-official      # React/Tailwind guidance
claude plugin install feature-dev@claude-plugins-official          # Code architect agents
claude plugin install code-review@claude-plugins-official          # Code review workflow
claude plugin install commit-commands@claude-plugins-official      # Git helpers
claude plugin install pr-review-toolkit@claude-plugins-official    # PR review
claude plugin install hookify@claude-plugins-official              # Custom hooks
claude plugin install typescript-lsp@claude-plugins-official       # TypeScript LSP
claude plugin install gopls-lsp@claude-plugins-official            # Go LSP
```

### Key Skills

The `/super` command uses the `superpowers:brainstorming` skill. Other useful skills:

| Skill | When to Use |
|-------|-------------|
| `superpowers:brainstorming` | Before creative work, designing features |
| `superpowers:writing-plans` | Creating implementation plans |
| `superpowers:systematic-debugging` | Bug investigation (find root cause first) |
| `superpowers:test-driven-development` | Writing new code (test first) |
| `superpowers:verification-before-completion` | Before claiming work is done |

## Commit Rules

These commands enforce:
- No Claude/Anthropic attribution in commits
- Imperative mood ("Add feature" not "Added feature")
- Summary under 72 characters

## Portability

These commands are designed to be portable across different workspaces.

### Using in a New Workspace

1. Copy the entire `.claude/commands/` folder to your workspace
2. Create `config.yaml` from the example template
3. Update `base_path` and add your repos
4. Commands work immediately - no code changes needed

### Customizing for Your Stack

**Different package managers:**
```yaml
commands:
  test: pnpm test
  lint: pnpm lint
  build: pnpm build
```

**Monorepos with multiple languages:**
```yaml
- name: my-fullstack
  commands:
    test: "cd backend && go test ./... && cd ../frontend && npm test"
    lint: "cd backend && golangci-lint run && cd ../frontend && npm run lint"
```

**Custom test patterns:**
```yaml
- name: my-django-app
  language: python
  commands:
    test: python manage.py test
    lint: ruff check . && black --check .
```

### What's Portable vs Local

| Portable (commit these) | Local (gitignored) |
|------------------------|-------------------|
| `*.md` command files | `config.yaml` |
| `config.yaml.example` | |
| `_shared-repo-logic.md` | |

### Forking for Your Organization

1. Fork this commands folder
2. Modify `_shared-repo-logic.md` for org-specific rules
3. Update `config.yaml.example` with your standard repos
4. Add org-specific commands as needed
