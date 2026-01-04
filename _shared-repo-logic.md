# Shared Repo Logic

This file contains shared patterns used by all repo-targeting slash commands.

---

## Configuration

Commands use `config.yaml` in this directory:

```yaml
workspace: ~/code/mono-claude

repos:
  - name: my-project      # Must match directory name exactly
    group: apps
    language: typescript
    work_dir: nextapp     # Optional: subdirectory for actual code
```

**Fields:**
- `workspace` - Root directory containing all repos
- `name` - Directory name (exact match required)
- `group` - Organization category (apps, devops, tools, experimental)
- `language` - Primary language (typescript, python, go, etc.)
- `work_dir` - Optional subdirectory for nested projects

**Setup:** Run `/setup-workspace` to auto-generate config from your directory.

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

| Command | Input | Purpose |
|---------|-------|---------|
| `devbot path <repo>` | repo name | Get full filesystem path (USE THIS FIRST) |
| `devbot status` | none | Git status across all repos |
| `devbot status <repo>` | repo name | Single repo git details |
| `devbot diff <repo>` | repo name | Git diff summary |
| `devbot branch <repo>` | repo name | Branch tracking info |
| `devbot check <repo>` | repo name | Run lint/typecheck/build/test |
| `devbot make <repo>` | repo name | Makefile target analysis |
| `devbot config <repo>` | repo name | Show config files |
| `devbot tree <path>` | **filesystem path** | Directory tree |
| `devbot stats <path>` | **filesystem path** | Code metrics |

### CRITICAL: Path vs Name Commands

**Commands that take repo NAME:** `path`, `status`, `diff`, `branch`, `check`, `make`, `config`

**Commands that take filesystem PATH:** `tree`, `stats`

**ALWAYS get the path first, then use it:**

```bash
# CORRECT - two-step process
REPO_PATH=$(devbot path fractals-nextjs)
devbot tree "$REPO_PATH"
devbot stats "$REPO_PATH"

# WRONG - DO NOT construct paths manually
devbot tree ~/code/fractals-nextjs        # ❌ Path may be wrong!
devbot stats ~/code/my-repo               # ❌ Never guess paths!
```

All repo-name commands require exact names from config.yaml.

Install: `/setup-workspace` (or rebuild with `cd devbot && make install`)

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

---

## Dual-Model Evaluation

For text generation tasks, use this pattern to build confidence in local model:

### Pattern

1. **Local model generates first** (fast path, ~0.5s)
2. **Claude generates independently** (for comparison)
3. **Evaluate local against criteria** (task-specific)
4. **Select winner:**
   - If local passes all criteria → use with `[local]` marker
   - If local fails any → use Claude's version

### Criteria Templates

**For prose/docs:**
- Factually accurate (matches actual code/commands)
- Concise (no bloat or unnecessary words)
- Active voice throughout
- No hallucinated features

**For task summaries:**
- Starts with actionable verb (Add, Fix, Refactor)
- Under 100 characters
- Accurate file/line references
- Reasonable priority inference

### Inline Markers

Show provenance in output:
```
[local] Add retry logic for API failures
[claude] Refactor authentication to use OAuth2 (local missed requirement)
```

### Commit Suffix

If >50% of generated content used local model, append ` [local]` to commit message.
This creates visible audit trail: `git log --oneline | grep "\[local\]"`
