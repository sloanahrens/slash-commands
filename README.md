# Claude Code Slash Commands

Portable slash commands for managing multi-repo workspaces with Claude Code.

## Features

- Multi-source repository discovery (builtin packages, worktrees, clones, standalone repos)
- Fuzzy matching for quick repo access
- Git worktree support for feature branch isolation
- Integration with devbot for fast parallel operations
- Optional Linear and GitHub MCP integration
- Optional local model acceleration via MLX

## Installation

### Option 1: Clone directly to commands folder

```bash
git clone https://github.com/sloanahrens/slash-commands.git ~/.claude/commands
cd ~/.claude/commands/devbot && make install
```

### Option 2: Clone anywhere and create symlinks

```bash
# Clone to your preferred location
git clone https://github.com/sloanahrens/slash-commands.git ~/code/slash-commands

# Create symlinks (run this command in Claude Code)
/setup-symlinks

# Install devbot
/install-devbot
```

### Post-Installation

Create your workspace config:
```bash
cp config.yaml.example config.yaml
# Edit config.yaml for your workspace
```

## Commands

| Command | Description |
|---------|-------------|
| `/super <repo>` | Start brainstorming session with context |
| `/find-tasks <repo>` | Find tasks from code, Linear, GitHub |
| `/run-tests <repo>` | Run lint, type-check, build, and tests |
| `/make-test <repo>` | Test Makefile targets interactively |
| `/yes-commit <repo>` | Draft and commit changes |
| `/push <repo>` | Push commits to origin |
| `/update-docs <repo>` | Update documentation |
| `/review-project <repo>` | Technical review with analysis |
| `/resolve-pr <url>` | Resolve GitHub PR review feedback |
| `/add-repo <url>` | Clone repo (reference or working) |
| `/status [repo]` | Show repository status |
| `/sync [repo]` | Pull latest changes |
| `/switch <repo>` | Context switch with suggestions |
| `/quick-explain <code>` | Quick code explanation |
| `/quick-gen <desc>` | Quick code generation |
| `/yes-proceed` | Accept recommendation and proceed |
| `/dev-rules` | Load workspace development rules |
| `/setup-plugins` | Install recommended plugins |
| `/install-devbot` | Build and install devbot CLI |
| `/setup-symlinks` | Create global command symlinks |
| `/list-commands` | List all available commands |
| `/list-skills` | List available skills from plugins |

All repo commands support fuzzy matching (e.g., `/run-tests cli`).

## Configuration

Create a `config.yaml` based on your workspace structure:

```yaml
# Primary workspace (monorepo or main project)
base_path: ~/code/my-workspace

# Fixed packages within base_path
builtin:
  - name: my-cli
    group: packages
    path: packages/my-cli
    language: typescript

  - name: my-server
    group: apps
    path: apps/my-server
    language: python

# Worktree directory name (under base_path)
worktrees_dir: .trees

# Reference repos config (optional)
clones_config: clones/clone-config.json

# Working code directory
code_path: ~/code

# Standalone repos at code_path
repos:
  - name: my-project
    group: projects
    aliases: [proj]
```

## Repository Types

| Type | Location | Description |
|------|----------|-------------|
| **Builtin** | `<base_path>/<path>` | Fixed packages in your monorepo |
| **Worktrees** | `<base_path>/.trees/` | Git worktrees for feature branches |
| **Clones** | `<base_path>/clones/` | Read-only reference repos |
| **Repos** | `<code_path>/<name>` | Standalone working repos |

## devbot CLI

The included `devbot` Go CLI provides fast parallel operations across your workspace.

### Installation

```bash
cd devbot && make install
# Or use: /install-devbot
```

### Commands

#### status - Parallel Git Status (~0.03s for 12 repos)

```bash
devbot status           # Show dirty repos (clean count summarized)
devbot status --all     # Show all repos
devbot status --dirty   # Only dirty repos
devbot status <repo>    # Single repo details
```

#### diff - Git Diff Summary (~0.02s)

```bash
devbot diff <repo>      # Staged/unstaged files with line counts
```

Shows branch, staged files, unstaged files with +/- counts in a single call.

#### check - Auto-Detected Quality Checks

```bash
devbot check <repo>              # Run all checks (lint, typecheck, build, test)
devbot check <repo> --only=lint  # Run specific checks
devbot check <repo> --fix        # Auto-fix where possible
```

Auto-detects stack (go, ts, nextjs, python, rust) and runs appropriate commands:
- Lint and typecheck run in parallel
- Build and test run sequentially
- Exits with code 1 on first failure

#### run - Parallel Command Execution

```bash
devbot run -- git pull              # Pull all repos in parallel
devbot run -- npm install           # Install deps in all repos
devbot run -f myapp -- make build   # Filter to repos matching "myapp"
devbot run -q -- git fetch          # Quiet mode (suppress empty output)
```

#### todos - Parallel TODO/FIXME Scanning

```bash
devbot todos                    # All TODOs across workspace
devbot todos --type FIXME       # Filter by marker type
devbot todos <repo>             # Single repo
```

Scans for: TODO, FIXME, HACK, XXX, BUG in .go, .ts, .tsx, .js, .jsx, .py, .md, .yaml, .yml files.

#### stats - Code Metrics and Complexity

```bash
devbot stats <path>             # Analyze file or directory
```

Reports: files, lines (code/comments/blank), functions, complexity flags.

#### detect - Project Stack Detection

```bash
devbot detect               # Current directory
devbot detect <path>        # Specific path
# Output: Detected: go, ts, nextjs
```

#### deps - Dependency Analysis

```bash
devbot deps             # Show shared dependencies (2+ repos)
devbot deps --all       # Show all dependencies by usage
devbot deps <repo>      # Analyze single repo
```

#### tree - Gitignore-Aware Tree

```bash
devbot tree                 # Current directory
devbot tree <path>          # Specific path
devbot tree -d 5            # Depth limit (default: 3)
```

#### config - Config File Discovery

```bash
devbot config               # All config files by type
devbot config --type go     # Filter by config type
devbot config <repo>        # Single repo
```

#### make - Makefile Target Analysis

```bash
devbot make                     # All targets grouped by category
devbot make --category test     # Filter by category
devbot make <repo>              # Single repo
```

#### worktrees - Git Worktree Discovery

```bash
devbot worktrees                # All worktrees across repos
devbot worktrees <repo>         # Single repo
```

### Architecture

```
devbot/
├── cmd/devbot/main.go       # CLI entry point (cobra)
└── internal/
    ├── workspace/           # Repo discovery and parallel git status
    ├── runner/              # Parallel command execution
    ├── deps/                # Dependency analysis
    ├── tree/                # Gitignore-aware directory tree
    ├── detect/              # Project stack detection
    ├── todos/               # Parallel TODO/FIXME scanning
    ├── stats/               # Code metrics and complexity
    ├── config/              # Config file discovery
    ├── makefile/            # Makefile target parsing
    ├── worktrees/           # Git worktree discovery
    └── output/              # Terminal rendering
```

## Worktree Workflow

For feature development, use git worktrees:

```bash
# Create worktree for feature branch
git worktree add .trees/feature-name -b feature/feature-name

# Switch to worktree
/switch feature-name

# Commands work in worktrees
/run-tests feature-name
/yes-commit feature-name
/push feature-name
```

## Local Model Acceleration (Optional)

Commands can use a local MLX model for faster processing of simple tasks. Requires the mlx-hub plugin.

**Model**: `mlx-community/Qwen2.5-Coder-14B-Instruct-4bit`

| Command | Local Model Use |
|---------|-----------------|
| `/yes-commit` | Draft commit messages |
| `/quick-explain` | Code explanations |
| `/quick-gen` | Simple code generation |

Output is labeled `[local]` vs `[claude]` to indicate the source.

## MCP Integration (Optional)

### Linear

For repos with `linear_project` in config:
```
mcp__plugin_linear_linear__list_issues
mcp__plugin_linear_linear__get_issue
mcp__plugin_linear_linear__create_issue
```

### GitHub

If you have GitHub MCP tools configured:
- Search issues by project
- Get assigned issues with status
- Update issue status

## Commit Rules

| Do | Don't |
|----|-------|
| Use imperative mood | Include Claude/Anthropic attribution |
| Keep summary under 72 chars | Include co-author lines |
| Focus on WHY not WHAT | Include "Generated with" tags |

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Go 1.23+ (for devbot)
- Node.js >=18.0.0 (for TypeScript projects)
- Python >=3.12 (for Python projects)

## Recommended Plugins

Run `/setup-plugins` to install all recommended plugins, or install manually:

### Core Plugins

```bash
claude plugin marketplace add obra/superpowers-marketplace
claude plugin install superpowers@superpowers-marketplace
claude plugin install episodic-memory@superpowers-marketplace
```

### Official Plugins

```bash
claude plugin marketplace add anthropics/claude-plugins-official
claude plugin install code-review@claude-plugins-official
claude plugin install commit-commands@claude-plugins-official
claude plugin install pr-review-toolkit@claude-plugins-official
```

## Files

| File | Purpose |
|------|---------|
| `config.yaml` | Your workspace config (gitignored) |
| `config.yaml.example` | Template for config |
| `_shared-repo-logic.md` | Multi-source repo discovery logic |
| `devbot/` | Go CLI for parallel operations |

## Contributing

This repo uses conventional commits:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation updates
- `refactor:` Code refactoring
