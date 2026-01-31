# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Global Claude Code configuration. Each repo also has its own CLAUDE.md - read that first.

## Architecture

This repo IS `~/.claude/` - the Claude Code configuration directory. Key components:

| Component | Purpose |
|-----------|---------|
| `devbot/` | Go CLI for parallel workspace operations |
| `commands/` | 30+ slash commands for workflow automation |
| `hookify.*.md` | Block dangerous bash patterns |
| `hooks/` | Session start/end automation |

## Setup

```bash
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

**Workflow commands:**
| Command | Description |
|---------|-------------|
| `/super-plan <repo>` | Brainstorm → design → implementation plan |
| `/plan-to-beads <repo>` | Convert plan to exhaustive Beads issues |
| `/execute-plan <repo>` | Resume plan execution with Beads |
| `/run-tests <repo>` | Lint, type-check, build, test |
| `/yes-commit <repo>` | Draft and commit (no AI attribution) |
| `/create-pr <repo>` | Create pull request |

**Knowledge capture:**
| Command | Description |
|---------|-------------|
| `/prime-context <repo>` | Load context + Linear issues + plan/beads status |
| `/capture-session` | Sync Beads + post Linear updates + log decisions |

## devbot CLI

**NAME commands:** `path`, `status`, `diff`, `branch`, `log`, `show`, `fetch`, `switch`, `check`, `make`, `todos`, `last-commit`, `config`, `deps`, `remote`, `worktrees`, `pulumi`, `deploy`, `find-repo`

**PATH commands:** `tree`, `stats`, `detect` (use `devbot path` first)

**GitHub lookup:** `devbot find-repo owner/repo` or `devbot find-repo https://github.com/owner/repo/pull/123` → returns local repo name

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

## Beads Issue Tracking

Beads (`bd`) is a git-backed issue tracker for AI agents. Use for structured work tracking.

### Installation

```bash
brew tap steveyegge/beads
brew install bd
```

Source: https://github.com/steveyegge/beads

### Running Beads Commands

**ALWAYS use `devbot exec` to run beads commands.** This enables pre-approval and avoids `cd &&` patterns.

```bash
# Good - pre-approvable
devbot exec my-repo bd ready
devbot exec my-repo bd show abc123
devbot exec my-repo bd update abc123 --status in_progress
devbot exec my-repo bd close abc123
devbot exec my-repo bd sync

# Bad - requires manual approval each time
cd ~/code/my-repo && bd ready
bd ready  # (unless already in repo directory)
```

### Quick Reference

```bash
devbot exec <repo> bd ready              # See unblocked work
devbot exec <repo> bd show <id>          # View issue details
devbot exec <repo> bd create "Title"     # Create issue (--type task|feature|bug)
devbot exec <repo> bd update <id> --status in_progress  # Start work
devbot exec <repo> bd close <id>         # Complete work
devbot exec <repo> bd sync               # Sync with git
```

### During Work (if repo has `.beads/`)

1. **Start session** → Run `/prime-context <repo>` to pull changes and see ready work
2. **Pick a task** → `devbot exec <repo> bd show <id>` for details, then `devbot exec <repo> bd update <id> --status in_progress`
3. **Create new work** → `devbot exec <repo> bd create "Description" --type task` (or `feature`, `bug`)
4. **Complete work** → `devbot exec <repo> bd close <id>` when done
5. **End session** → Run `/capture-session` to sync and push

**Proactive tracking:** When working on multi-step tasks, create Beads issues to track progress. This helps continuity across sessions.

### Two-Layer System

| Layer | Purpose | Location |
|-------|---------|----------|
| Beads | Structured work tracking | `<repo>/.beads/` |
| Decisions log | Narrative context | `<repo>/.claude/decisions.md` |

**Beads** tracks: tasks, features, bugs, dependencies, status
**Decisions log** tracks: why decisions were made, constraints, learnings

### Initialize Beads in a Repo (Protected Branch Workflow)

```bash
cd /path/to/repo
bd init --branch beads-sync

# Add JSONL files to local exclude (not committed)
cat >> .git/info/exclude << 'EOF'
.beads/issues.jsonl
.beads/interactions.jsonl
.beads/metadata.json
EOF

# Commit config files to main
git add .beads/.gitignore .beads/config.yaml .beads/README.md
git commit -m "Add Beads config for protected branch workflow"

# Push both branches
git push origin main beads-sync
```

**Branch layout:**
- `main`: Config files only (`.beads/.gitignore`, `.beads/config.yaml`)
- `beads-sync`: Issue data (`.beads/issues.jsonl`) - synced via `bd sync`

## Linear Integration

Connect repos to Linear projects for tracking external issues alongside local Beads.

### Config Fields (in config.yaml)

```yaml
repos:
  - name: my-repo
    linear_projects: ["Project Name", "Another Project"]  # Linear project names
    plan_paths: ["docs/plans", "~/code/docs/plans"]       # Where to find plan docs
```

- `linear_projects`: Array of Linear project names to fetch issues from
- `plan_paths`: Directories to search for plan documents (default: `docs/plans`)

### How It Works

**Hierarchy:** `Linear Issue → Plan File → Beads`

1. `/prime-context <repo>` fetches open Linear issues from configured projects
2. Matches issues to plan files via keyword/URL matching
3. Traces plans to beads for status tracking
4. Shows summary table with actionable gaps

**Session end:** `/capture-session` posts progress updates to matched Linear issues
- Comments only (never changes issue status)
- Idempotent (checks existing comments, skips duplicates)
- Bullet list format, short and human-readable

### Example Output (from /prime-context)

```
| Issue   | Title              | Status      | Plan                  | Beads          |
|---------|--------------------|-------------|-----------------------|----------------|
| XYZ-15  | Auth middleware    | In Progress | ✓ auth-plan.md        | ✓ 5 tasks (2/5)|
| XYZ-18  | API routes         | Backlog     | — No plan             | —              |

Actionable Gaps:
  • XYZ-18 needs planning → /super-plan my-repo api-routes
```

## Key Skills

| Skill | When |
|-------|------|
| `superpowers:brainstorming` | Before creative work |
| `superpowers:writing-plans` | After design, before implementation |
| `superpowers:executing-plans` | Implementing from a plan |
| `superpowers:systematic-debugging` | Bug investigation |
| `superpowers:verification-before-completion` | Before claiming done |

## Feature Development Workflow

Three commands for feature development with multi-session tracking:

```
/super-plan <repo> <topic>           — Design + implementation plan
    │
    └─▶ /plan-to-beads <repo>   — Convert plan to Beads issues
            │
            └─▶ /execute-plan <repo>  — Beads-aware execution
```

| Step | Command | What it does |
|------|---------|--------------|
| 1 | `/super-plan` | Brainstorm → design doc → implementation plan |
| 2 | `/plan-to-beads` | Create one Bead per task with dependencies |
| 3 | `/execute-plan` | Execute tasks, tracking via Beads |

**Beads-Aware Execution:** When executing plans, ALWAYS track via Beads:
- Before task: `devbot exec <repo> bd update <task-id> --status in_progress`
- After task: `devbot exec <repo> bd close <task-id>`
- Sync periodically: `devbot exec <repo> bd sync`
- On completion: `devbot exec <repo> bd close <feature-id>` then `devbot exec <repo> bd sync`

**Resuming:** Use `/execute-plan <repo>` after context compaction or new session to re-prime the task-to-Bead mapping.

## Repo Context Workflow

Repos can use Beads (preferred) or legacy session notes. Check for `.beads/` to determine which.

### With Beads + Linear (preferred)

1. **Start work** → `/prime-context <repo>` shows Linear issues + plans + beads status
2. **Pick work** → Choose from ready beads or plan untracked Linear issues
3. **Do work** → `devbot exec <repo> bd update <id> --status in_progress`, then `devbot exec <repo> bd close <id>`
4. **End session** → `/capture-session` syncs Beads + posts Linear updates

### With Beads only

1. **Start work** → `/prime-context <repo>` runs `bd ready` + shows recent decisions
2. **Pick work** → `devbot exec <repo> bd show <id>` for details, `devbot exec <repo> bd update <id> --status in_progress`
3. **Complete work** → `devbot exec <repo> bd close <id>` when done
4. **End session** → `/capture-session` syncs Beads + logs decisions

### Without Beads (legacy)

1. **Start work** → `/prime-context <repo>` loads project context + most recent session note
2. **Do work** → Session notes link to previous sessions if more context needed
3. **End session** → `/capture-session` saves progress to repo's `.claude/sessions/`

## Files

| Location | Purpose |
|----------|---------|
| `~/.claude/CLAUDE.md` | This file - global instructions |
| `~/.claude/settings.json` | Permissions + plugins |
| `~/.claude/config.yaml` | Workspace config (gitignored) |
| `~/.claude/hookify.*.md` | Hookify rules |
| `~/.claude/commands/` | Slash commands |
| `~/.claude/devbot/` | CLI tool source (Go) |
| `~/.claude/hooks/` | Session start/end hooks |
| `<repo>/.beads/` | Beads config (main) + data (beads-sync branch) |
| `<repo>/.claude/` | Repo-local context (gitignored) |
| `<repo>/.claude/decisions.md` | Key decisions and rationale |
| `<repo>/.claude/sessions/` | Legacy session notes (if no Beads) |
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

## Developing devbot

```bash
make -C ~/.claude/devbot build    # Build binary
make -C ~/.claude/devbot test     # Run tests
make -C ~/.claude/devbot ci       # Full CI: fmt, vet, test, lint
make -C ~/.claude/devbot install  # Install to PATH
```

Add new commands in `devbot/internal/` following existing patterns (cobra CLI in `cmd/devbot/main.go`).
