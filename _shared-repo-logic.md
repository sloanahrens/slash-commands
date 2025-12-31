# Shared Repo Logic (Trabian Branch)

This file contains shared patterns used by all repo-targeting slash commands in the trabian workspace.

---

## Configuration

Commands use `config.yaml` in this directory. The trabian branch supports:

- `builtin` - trabian's own packages (fixed paths)
- `worktrees_dir` - auto-discovers `.trees/*` worktrees
- `clones_config` - reads trabian's `clones/clone-config.json`
- `repos` - additional repos cloned into base_path

```yaml
base_path: ~/trabian

builtin:
  - name: trabian-cli
    group: packages
    path: packages/trabian-cli
    language: typescript

worktrees_dir: .trees
clones_config: clones/clone-config.json

repos:
  - name: client-project
    group: apps
    aliases: [client]
```

---

## Critical Rule

**CRITICAL**: Always stay within `~/trabian/` - never navigate above this directory.

---

## Repo Discovery

Parse repos from multiple sources in order:

### 1. Builtin Components

Read `config.yaml` → `builtin[]`:
- These are fixed trabian packages
- Always available regardless of clones

### 2. Worktrees (Dynamic)

Scan `<base_path>/.trees/` directory:
```bash
ls -d ~/trabian/.trees/*/ 2>/dev/null
```

For each worktree, extract:
- Name: directory name (e.g., `feature-new-auth`)
- Path: `.trees/<name>`
- Branch: `git -C .trees/<name> branch --show-current`

### 3. Clones (From trabian config)

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
- Path: `clones/<name>`
- Group: `clones`

### 4. Additional Repos

Read `config.yaml` → `repos[]`:
- These are working repos cloned into base_path
- Added via `/sloan/add-repo`

---

## Groups

| Group | Source | Description |
|-------|--------|-------------|
| `packages` | builtin | trabian TypeScript packages |
| `mcp` | builtin | trabian MCP server |
| `worktrees` | dynamic | Active feature branches in .trees/ |
| `clones` | clone-config.json | Reference repositories (Q2 SDK, Tecton) |
| `apps` | repos | Additional working repos |

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
  4. fix/mcp-bug

Clones:
  5. q2-sdk
  6. tecton

Apps:
  7. client-project

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
| `auth` | feature/new-auth (worktree) |

---

## Path Resolution

Once a repo is selected, resolve its full path:

| Source | Path Pattern |
|--------|--------------|
| builtin | `<base_path>/<path>` (e.g., `~/trabian/packages/trabian-cli`) |
| worktree | `<base_path>/.trees/<name>` |
| clone | `<base_path>/clones/<name>` |
| repo | `<base_path>/<name>` |

---

## Commit Rules

When committing changes in any repo:

- **NO** Claude/Anthropic attribution
- **NO** co-author lines
- **NO** "generated with" tags
- Use imperative mood ("Add feature" not "Added feature")
- Keep summary under 72 characters

---

## Trabian Context

After selecting a repo, load relevant context:

1. **Always read**: `~/trabian/CLAUDE.md` (workspace rules)
2. **If repo has CLAUDE.md**: Read `<repo-path>/CLAUDE.md`
3. **For MCP server**: Note Python/uv patterns
4. **For packages**: Note TypeScript/npm patterns

---

## Standard Process Start

1. **Apply dev rules** → Reference trabian's CLAUDE.md
2. Parse `config.yaml` for base path and repo definitions
3. Discover worktrees from `.trees/`
4. Discover clones from `clones/clone-config.json`
5. If `$ARGUMENTS` empty → show selection prompt
6. If `$ARGUMENTS` provided → fuzzy match to repo
7. Confirm selection: "Working on: <repo-name>"
8. Read repo's CLAUDE.md (if exists) for repo-specific guidance

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
