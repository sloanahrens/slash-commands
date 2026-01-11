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
3. **Prime with notes** (optional but recommended):
   - Run `/prime <repo>` to surface relevant patterns and hindsight
   - Or manually search `docs/patterns/` and `~/.claude/notes/`

### Note Locations

| Location | Contents |
|----------|----------|
| `docs/patterns/` | Versioned patterns (git tracked) |
| `~/.claude/notes/hindsight/` | Local failure captures |
| `~/.claude/notes/sessions/` | Local session summaries |

See `/prime`, `/capture-hindsight`, and `/promote-pattern` for the full knowledge management workflow.

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
| `devbot exec <repo>[/subdir] <cmd>` | repo + command | Run command in repo directory |
| `devbot port <port> [--kill]` | port number | Check/kill process on port |
| `devbot prereq <repo>[/subdir]` | repo name | Validate tools, deps, and env vars |

### CRITICAL: Path vs Name Commands

**Commands that take repo NAME:** `path`, `status`, `diff`, `branch`, `check`, `make`, `config`, `exec`

**Commands that take filesystem PATH:** `tree`, `stats`

**ALWAYS get the path first, then use it:**

```bash
# CORRECT - two separate commands
devbot path fractals-nextjs
# Output: /Users/sloan/code/mono-claude/fractals-nextjs
devbot tree /Users/sloan/code/mono-claude/fractals-nextjs  # Use literal path

# WRONG - compound commands or manual paths
REPO_PATH=$(devbot path repo) && devbot tree "$REPO_PATH"  # ❌ Compound
devbot tree ~/code/fractals-nextjs                          # ❌ Guessed path
```

All repo-name commands require exact names from config.yaml.

Install: `/setup-workspace` (or `make -C ~/code/mono-claude/slash-commands/devbot install`)

---

## Running Commands in Repos

Use `devbot exec` instead of `cd && command` (which is blocked by hookify):

```bash
# Instead of: cd /path/to/repo && npm run build
devbot exec <repo-name> npm run build

# For monorepo subprojects
devbot exec <repo-name>/<subdir> <command>

# Override work_dir and use repo root (trailing slash)
devbot exec <repo-name>/ <command>
```

**Directory resolution:**
1. If `/subdir` specified → `{repo_path}/{subdir}`
2. If trailing slash (`repo/`) → repo root (ignores work_dir)
3. If `work_dir` in config.yaml → `{repo_path}/{work_dir}`
4. Otherwise → `{repo_path}`

**Examples:**

| Command | Runs in |
|---------|---------|
| `devbot exec atap-automation2 npm test` | `.../atap-automation2/nextapp` (uses work_dir) |
| `devbot exec mango/go-api go build` | `.../mango/go-api` (explicit subdir) |
| `devbot exec slash-commands/devbot make` | `.../slash-commands/devbot` |
| `devbot exec atap-automation2/ docker build .` | `.../atap-automation2` (root, ignores work_dir) |

---

## Local Model Acceleration (Optional)

For simple tasks, use local models via mlx-hub plugin:

| Use Local Model For | Stay on Claude For |
|---------------------|-------------------|
| Commit messages | Security analysis |
| Code explanation | Architecture decisions |
| Simple code gen | Multi-file refactoring |

### Availability Check

**Before using local model, verify availability:**

```bash
# Check if mlx-hub plugin is installed
claude plugin list 2>/dev/null | grep -q mlx-hub
```

**If mlx-hub is unavailable:**
- Skip dual-model evaluation entirely
- Use Claude directly for all generation
- Note in output: "(local model unavailable, using Claude)"

**Requirements for local model:**
- Apple Silicon Mac (M1/M2/M3/M4)
- mlx-hub plugin installed (`/setup-plugins`)
- Model downloaded (`mlx-community/Qwen2.5-Coder-14B-Instruct-4bit`)

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

For text generation tasks, use this pattern to build confidence in local model.

**Skip this section if local model is unavailable** - use Claude directly instead.

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

### Inline Markers (Comparison Only)

During evaluation, show both for comparison:
```
[local]  Add retry logic for API failures
[claude] Refactor authentication to use OAuth2 (local missed requirement)
```

This helps evaluate local model quality over time.

### Commit Message Suffix

- If local model message is used → append ` [local]` to commit message
- If Claude message is used → no suffix (Claude is the default, assumed)

This creates visible audit trail: `git log --oneline | grep "\[local\]"`
