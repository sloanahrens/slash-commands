# Shared Repo Logic

This file contains shared patterns used by all repo-targeting slash commands.

---

## Configuration

Commands use `config.yaml` in this directory. The config supports:

- `base_path` - trabian product workspace (~/code/trabian-ai)
- `builtin` - trabian's own packages (fixed paths under base_path)
- `worktrees_dir` - auto-discovers `.trees/*` worktrees
- `clones_config` - reads trabian's `clones/clone-config.json`
- `code_path` - location of working code repos (~/code)
- `repos` - working repos at code_path

```yaml
base_path: ~/code/trabian-ai

builtin:
  - name: trabian-cli
    group: packages
    path: packages/trabian-cli
    language: typescript

worktrees_dir: .trees
clones_config: clones/clone-config.json

code_path: ~/code

repos:
  - name: hanscom-fcu-poc-plaid-token-manager
    group: projects
    aliases: [hanscom, plaid-poc]
```

---

## Critical Rule

**CRITICAL**: Stay within the configured paths:
- `~/code/trabian-ai/` for trabian product work
- `~/code/` for working repos

Never navigate above these directories.

---

## Repo Discovery

Parse repos from multiple sources in order:

### 1. Builtin Components

Read `config.yaml` → `builtin[]`:
- These are fixed trabian packages
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
    "q2-sdk": { "url": "...", "description": "..." },
    "tecton": { "url": "...", "description": "..." }
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
| `packages` | builtin | ~/code/trabian-ai/packages/ | trabian TypeScript packages |
| `mcp` | builtin | ~/code/trabian-ai/mcp/ | trabian MCP server |
| `worktrees` | dynamic | ~/code/trabian-ai/.trees/ | Active feature branches |
| `clones` | clone-config.json | ~/code/trabian-ai/clones/ | Reference repos (Q2 SDK, Tecton) |
| `projects` | repos | ~/code/ | Active project work |
| `devops` | repos | ~/code/ | Infrastructure repos |
| `personal` | repos | ~/code/ | Personal/exploration repos |

---

## Repo Selection

**If `$ARGUMENTS` is empty:**

Display grouped list and ask user to select:

```
Select a repository:

Packages:
  1. trabian-cli
  2. trabian-server

Worktrees:
  3. feature/new-auth

Clones:
  4. q2-sdk
  5. tecton

Projects:
  6. hanscom-fcu-poc-plaid-token-manager
  7. service-cu-cloud-services-platform

DevOps:
  8. devops-pulumi-ts

Personal:
  9. fractals-nextjs
  10. mango

Enter number or name:
```

**If `$ARGUMENTS` is provided:**

Fuzzy match against:
1. Directory names
2. Configured aliases
3. Worktree branch names

| Input | Matches |
|-------|---------|
| `cli` | trabian-cli |
| `server` | trabian-server |
| `q2` | q2-sdk |
| `hanscom` | hanscom-fcu-poc-plaid-token-manager |
| `pulumi` | devops-pulumi-ts |
| `fractals` | fractals-nextjs |

---

## Path Resolution

Once a repo is selected, resolve its full path:

| Source | Path Pattern |
|--------|--------------|
| builtin | `<base_path>/<path>` (e.g., `~/code/trabian-ai/packages/trabian-cli`) |
| worktree | `<base_path>/.trees/<name>` (e.g., `~/code/trabian-ai/.trees/feature-x`) |
| clone | `<base_path>/clones/<name>` (e.g., `~/code/trabian-ai/clones/q2-sdk`) |
| repo | `<code_path>/<name>` (e.g., `~/code/hanscom-fcu-poc-plaid-token-manager`) |

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

1. **For trabian repos**: Read `~/code/trabian-ai/CLAUDE.md`
2. **For code repos**: Read `<repo-path>/CLAUDE.md` if it exists
3. **For MCP server**: Note Python/uv patterns
4. **For packages**: Note TypeScript/npm patterns

---

## Standard Process Start

1. Parse `config.yaml` for base_path, code_path, and repo definitions
2. Discover worktrees from `<base_path>/.trees/`
3. Discover clones from `clones/clone-config.json`
4. If `$ARGUMENTS` empty → show selection prompt
5. If `$ARGUMENTS` provided → fuzzy match to repo
6. Confirm selection: "Working on: <repo-name>"
7. Read repo's CLAUDE.md (if exists) for repo-specific guidance

---

## Local Model Acceleration

Commands can use local Qwen model for 5-18x speed gains. Requires `mlx-hub` plugin (installed via `/setup-plugins`).

**See workspace `CLAUDE.md` → "Automatic Local Acceleration" for full routing rules.**

### Quick Reference

| Use Qwen For | Stay on Claude For |
|--------------|-------------------|
| Commit messages | Security analysis |
| Code explanation | Architecture decisions |
| Simple code gen | Multi-file refactoring |
| Type fixes | Complex debugging |

### Output Format

Always prefix local model output:
```
[qwen] Drafting commit message...
[qwen] Generated: "feat(utils): add validation helper"
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

## Linear Integration

For repos with `linear_project` in config:
- Use trabian's Linear MCP plugin tools
- `mcp__plugin_linear_linear__list_issues`
- `mcp__plugin_linear_linear__get_issue`
- `mcp__plugin_linear_linear__create_issue`

---

## GitHub Integration

Use trabian's GitHub MCP for project data:
- `mcp__trabian__github_get_assigned_issues_with_project_status`
- `mcp__trabian__github_get_project_items`
- `mcp__trabian__github_find_project_by_name`

---

## devbot CLI

Fast parallel operations across ~/code/ repos. Use devbot for speed-critical operations:

| Command | Purpose | Speed |
|---------|---------|-------|
| `devbot status` | Git status across all repos | ~0.03s |
| `devbot status <repo>` | Single repo details | ~0.01s |
| `devbot run -- <cmd>` | Parallel command execution | ~0.5s |
| `devbot todos` | TODO/FIXME scanning | ~0.1s |
| `devbot make` | Makefile target analysis | ~0.01s |
| `devbot worktrees` | Worktree discovery | ~0.01s |
| `devbot detect <path>` | Stack detection | instant |
| `devbot config` | Config file discovery | ~0.01s |
| `devbot stats <path>` | File/directory code metrics | ~0.01s |

Install: Run `/install-devbot` or `cd ~/code/slash-commands/devbot && make install`

