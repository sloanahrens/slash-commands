# Shared Repo Logic

This file contains shared patterns used by all repo-targeting slash commands.

---

## Configuration

Commands use `config.yaml` in this directory:

```yaml
base_path: ~/code/mono-claude
code_path: ~/code/mono-claude

repos:
  - name: my-project      # Must match directory name exactly
    group: apps
    language: typescript
    work_dir: nextapp     # Optional: subdirectory for actual code
```

**Fields:**
- `name` - Directory name (exact match required)
- `group` - Organization category (apps, devops, tools, experimental)
- `language` - Primary language (typescript, python, go, etc.)
- `work_dir` - Optional subdirectory for nested projects

---

## Repo Resolution

### When user provides `@directory/`

The `@` prefix means Claude Code passed a directory context:

1. Extract directory name from path (e.g., `@fractals-nextjs/` → `fractals-nextjs`)
2. Use `devbot path <name>` to get full path
3. If not found, show suggestion and ask user

### When user provides plain name

1. Use `devbot path <name>` to get full path
2. If found → use that path
3. If not found → show suggestion, ask user to confirm

### Getting the full path

**ALWAYS use devbot to get paths:**

```bash
devbot path fractals-nextjs
# Output: /Users/sloan/code/mono-claude/fractals-nextjs
```

**NEVER construct paths manually.** Do not assume `~/code/<name>` or any other pattern.

### When name not found

devbot suggests similar names:

```bash
devbot path fractals
# Output: Repository 'fractals' not found. Did you mean:
#   fractals-nextjs
```

Show this to the user and ask them to confirm or provide the correct name.

---

## Standard Process

1. Extract repo name from `$ARGUMENTS`
2. Run `devbot path <name>` to get full path
3. If not found, show suggestion and ask user
4. Confirm: "Working on: <repo-name>"
5. Load context (see below)

---

## Context Loading

After resolving the repo path:

1. Read `~/.claude/CLAUDE.md` (global settings)
2. Read `<repo-path>/CLAUDE.md` (repo-specific guidance)

---

## Commit Rules

When committing changes:

- **NO** Claude/Anthropic attribution
- **NO** co-author lines
- **NO** "generated with" tags
- Use imperative mood ("Add feature" not "Added feature")
- Keep summary under 72 characters

---

## devbot CLI

Fast operations across repos:

| Command | Purpose |
|---------|---------|
| `devbot path <repo>` | Get full filesystem path for repo |
| `devbot status` | Git status across all repos |
| `devbot status <repo>` | Single repo git details |
| `devbot diff <repo>` | Git diff summary |
| `devbot branch <repo>` | Branch tracking info |
| `devbot check <repo>` | Run lint/typecheck/build/test |
| `devbot make <repo>` | Makefile target analysis |
| `devbot tree <path>` | Directory tree |
| `devbot stats <path>` | Code metrics |

All commands require exact repo names from config.yaml.

Install: `/install-devbot`

---

## Local Model Acceleration (Optional)

For simple tasks, use local models via mlx-hub plugin:

| Use Local Model For | Stay on Claude For |
|---------------------|-------------------|
| Commit messages | Security analysis |
| Code explanation | Architecture decisions |
| Simple code gen | Multi-file refactoring |

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="...",
  max_tokens=200
)
```

Always prefix output: `[local] Generated: "..."`
