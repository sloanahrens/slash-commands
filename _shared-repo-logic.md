# Shared Repo Logic

This file contains shared patterns used by all repo-targeting slash commands.

---

## Configuration

Commands use `config.yaml` in this directory. The config supports:

- `base_path` - primary workspace for monorepo/packages (e.g., ~/code/my-workspace)
- `builtin` - fixed packages within base_path
- `worktrees_dir` - auto-discovers `.trees/*` worktrees
- `clones_config` - reads reference repos from `clones/clone-config.json`
- `code_path` - location of working code repos (~/code)
- `repos` - working repos at code_path

```yaml
base_path: ~/code/my-workspace

builtin:
  - name: my-cli
    group: packages
    path: packages/my-cli
    language: typescript

worktrees_dir: .trees
clones_config: clones/clone-config.json

code_path: ~/code

repos:
  - name: my-project
    group: projects
    aliases: [proj, mp]
```

---

## Critical Rule

**CRITICAL**: Stay within the configured paths:
- `<base_path>` for monorepo/workspace work
- `<code_path>` for standalone repos

Never navigate above these directories.

---

## Repo Discovery

Parse repos from multiple sources in order:

### 1. Builtin Components

Read `config.yaml` → `builtin[]`:
- These are fixed packages within your monorepo
- Path resolved as `<base_path>/<path>`

### 2. Worktrees (Dynamic)

Use devbot for fast worktree discovery:
```bash
devbot worktrees
```

This scans `.trees/`, `worktrees/`, `.worktrees/` directories across all repos in parallel (~0.01s) and returns:
- Name: directory name (e.g., `feature-new-auth`)
- Path: full path to worktree
- Branch: current branch
- Status: dirty file count

### 3. Clones (Reference Repos)

Read `<base_path>/clones/clone-config.json`:
```json
{
  "repositories": {
    "some-sdk": { "url": "...", "description": "..." },
    "other-lib": { "url": "...", "description": "..." }
  }
}
```

For each clone:
- Name: key name
- Path: `<base_path>/clones/<name>`
- Group: `clones`
- **Note**: These are read-only reference repos

### 4. Working Repos

Read `config.yaml` → `repos[]`:
- These are working repos at `<code_path>/<name>`
- Added via `/add-repo` or manually to config

---

## Groups

| Group | Source | Location | Description |
|-------|--------|----------|-------------|
| `packages` | builtin | `<base_path>/packages/` | Monorepo packages |
| `apps` | builtin | `<base_path>/apps/` | Monorepo applications |
| `worktrees` | dynamic | `<base_path>/.trees/` | Active feature branches |
| `clones` | clone-config.json | `<base_path>/clones/` | Reference repos |
| `projects` | repos | `<code_path>/` | Active project work |
| `devops` | repos | `<code_path>/` | Infrastructure repos |
| `personal` | repos | `<code_path>/` | Personal/exploration repos |

---

## Repo Selection

**If `$ARGUMENTS` is empty:**

Display grouped list and ask user to select:

```
Select a repository:

Packages:
  1. my-cli
  2. my-server

Worktrees:
  3. feature/new-auth

Clones:
  4. some-sdk

Projects:
  5. my-project
  6. another-project

Enter number or name:
```

**If `$ARGUMENTS` is provided:**

Fuzzy match against:
1. Directory names
2. Configured aliases
3. Worktree branch names

| Input | Matches |
|-------|---------|
| `cli` | my-cli |
| `server` | my-server |
| `sdk` | some-sdk |
| `proj` | my-project |

**IMPORTANT**: After fuzzy matching, always use the **full repo name** (from config.yaml `name` field) for all devbot commands, NOT the alias that was matched. Devbot requires exact repo names.

---

## Path Resolution

Once a repo is selected, resolve its full path:

| Source | Path Pattern |
|--------|--------------|
| builtin | `<base_path>/<path>` (e.g., `~/code/my-workspace/packages/my-cli`) |
| worktree | `<base_path>/.trees/<name>` (e.g., `~/code/my-workspace/.trees/feature-x`) |
| clone | `<base_path>/clones/<name>` (e.g., `~/code/my-workspace/clones/some-sdk`) |
| repo | `<code_path>/<name>` (e.g., `~/code/my-project`) |

---

## Commit Rules

When committing changes in any repo:

- **NO** Claude/Anthropic attribution
- **NO** co-author lines
- **NO** "generated with" tags
- Use imperative mood ("Add feature" not "Added feature")
- Keep summary under 72 characters

---

## Context Loading

After selecting a repo, load relevant context:

1. **Global config**: Read `~/.claude/CLAUDE.md` (if exists) for user-wide settings
2. **Repo-specific**: Read `<repo-path>/CLAUDE.md` if it exists
3. **For Python projects**: Note Python/uv patterns
4. **For TypeScript projects**: Note TypeScript/npm patterns

---

## Standard Process Start

1. Parse `config.yaml` for base_path, code_path, and repo definitions
2. Discover worktrees from `<base_path>/.trees/`
3. Discover clones from `clones/clone-config.json`
4. If `$ARGUMENTS` empty → show selection prompt
5. If `$ARGUMENTS` provided → fuzzy match to repo
6. Confirm selection: "Working on: <repo-name>"
7. Read global `~/.claude/CLAUDE.md` (if exists) for user-wide settings
8. Read repo's CLAUDE.md (if exists) for repo-specific guidance

---

## Local Model Acceleration (Optional)

Commands can use local models for speed gains. Requires `mlx-hub` plugin.

### Quick Reference

| Use Local Model For | Stay on Claude For |
|---------------------|-------------------|
| Commit messages | Security analysis |
| Code explanation | Architecture decisions |
| Simple code gen | Multi-file refactoring |
| Type fixes | Complex debugging |

### Output Format

Always prefix local model output:
```
[local] Drafting commit message...
[local] Generated: "feat(utils): add validation helper"
```

### Usage

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="...",
  max_tokens=200
)
```

---

## Linear Integration (Optional)

For repos with `linear_project` in config:
- `mcp__plugin_linear_linear__list_issues`
- `mcp__plugin_linear_linear__get_issue`
- `mcp__plugin_linear_linear__create_issue`

---

## GitHub Integration (Optional)

If you have GitHub MCP tools configured:
- Search issues by project
- Get assigned issues with project status
- Find project items

---

## devbot CLI

Fast parallel operations across repos. Use devbot for speed-critical operations:

| Command | Purpose | Speed |
|---------|---------|-------|
| `devbot status` | Git status across all repos | ~0.03s |
| `devbot status <repo>` | Single repo details (use full name, not alias) | ~0.01s |
| `devbot diff <repo>` | Git diff summary (staged/unstaged with stats) | ~0.02s |
| `devbot branch <repo>` | Branch tracking, ahead/behind, commits to push | ~0.02s |
| `devbot remote <repo>` | Remote URLs and GitHub identifiers | ~0.01s |
| `devbot find-repo <gh-id>` | Find local repo by GitHub org/repo | ~0.03s |
| `devbot check <repo>` | Auto-detected lint/typecheck/build/test | varies |
| `devbot run -- <cmd>` | Parallel command execution | ~0.5s |
| `devbot todos` | TODO/FIXME scanning | ~0.1s |
| `devbot make` | Makefile target analysis | ~0.01s |
| `devbot worktrees` | Worktree discovery | ~0.01s |
| `devbot detect <path>` | Stack detection | instant |
| `devbot config` | Config file discovery | ~0.01s |
| `devbot stats <path>` | File/directory code metrics | ~0.01s |
| `devbot deps [repo]` | Dependency analysis (shared deps) | ~0.01s |
| `devbot tree <path>` | Directory tree (respects .gitignore) | ~0.01s |

Install: Run `/install-devbot` or `cd ~/code/slash-commands/devbot && make install`

