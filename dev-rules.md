---
description: Remind Claude of workspace development rules
---

# Dev Rules

Continue whatever you were doing, but remember these workspace rules:

---

## Critical Rules

- **NO Claude/Anthropic attribution** in commits (no "Generated with", no co-author lines)
- **Read repo CLAUDE.md first** - each repo has specific guidance
- **Use exact repo names** from config.yaml
- **Use devbot** - prefer over manual git/file operations
- **Simple bash only** - no `&&`, `$()`, `;`, or `git -C` (blocked by hookify)

---

## Tool Selection Guide

**STOP before running bash.** Use the right tool:

| Need | Tool | Example |
|------|------|---------|
| Git status/diff/branch | `devbot` | `devbot diff slash-commands --full` |
| Read file content | `Read` tool | Read tool on any file path |
| Search files | `Grep`/`Glob` | Grep for pattern, Glob for filenames |
| File operations | `Read`/`Edit`/`Write` | Never use cat/sed/awk |

**Never improvise bash commands.** If devbot doesn't have a command for it, use Claude Code's built-in tools (Read, Grep, Glob, Edit, Write).

---

## Git Operations

**NEVER use `git -C`** - this flag cannot be auto-approved, requiring manual user approval every time.

| Use This | NOT This |
|----------|----------|
| `devbot status <repo>` | `git status` or `git -C /path status` |
| `devbot diff <repo>` | `git diff` or `git -C /path diff` |
| `devbot branch <repo>` | `git branch -vv` or `git -C /path branch` |
| `devbot log <repo>` | `git log` or `git -C /path log` |
| `devbot show <repo> [ref]` | `git show` or `git -C /path show` |
| `devbot fetch <repo>` | `git fetch` or `git -C /path fetch` |
| `devbot switch <repo> <branch>` | `git switch` or `git checkout` |
| `devbot check <repo>` | `npm test && npm run lint` |
| `devbot last-commit <repo> [file]` | `git log -1 --format="%ar"` |

**Pattern for other git commands:**
```bash
devbot path my-repo        # → /full/path/to/my-repo
cd /full/path/to/my-repo
git commit                 # Regular git commands work after cd
```

---

## Subdirectory Commands (NO cd &&)

Run commands in subdirectories using flags, not `cd /path && cmd`:

| Tool | Pattern | Example |
|------|---------|---------|
| npm | `npm run <cmd> --prefix <path>` | `npm run build --prefix /path/to/app` |
| make | `make -C <path> <target>` | `make -C /path/to/app build` |
| timeout | `timeout <sec> <cmd>` | `timeout 5 npm run dev --prefix /path` |

**Sequential commands:** Run each command separately, one tool call at a time. Do NOT combine with `&&` or `;`.

---

## Path Safety

- **Run `pwd` before bash commands** - verify current location before file/path operations
- **Use absolute paths** - always use full paths from your configured workspace
- **Stay within workspace** - never navigate above your base_path or code_path
- **Respect worktree isolation** - `.trees/` worktrees are separate git environments

---

## File Creation

- **NO `/tmp` files** - create temporary/working files in `docs/` directories
- **Prefer editing over creating** - modify existing files when possible
- **Plans go in** `docs/plans/YYYY-MM-DD-<topic>-<type>.md`

---

## Commit Messages

| Do | Don't |
|----|-------|
| Use imperative mood ("Add feature") | Include Claude/Anthropic attribution |
| Keep summary under 72 characters | Include co-author lines |
| Focus on WHY not just WHAT | Include "Generated with" tags |
| Consider compliance implications | Commit secrets (.env, credentials) |

---

## devbot Quick Reference

**NAME commands** (take repo name):
- `path`, `status`, `diff`, `branch`, `log`, `show`, `fetch`, `switch`, `check`, `make`, `todos`, `last-commit`, `config`, `deps`, `remote`, `worktrees`, `pulumi`, `deploy`

**PATH commands** (take file path, use `devbot path` first):
- `tree`, `stats`, `detect`

**Other:**
- `run` - parallel command across repos
- `find-repo` - GitHub org/repo lookup

**Critical:**
- `pulumi <repo>` - **MANDATORY before any Pulumi operation** - shows infrastructure state

**Git wrappers** (faster, auto-approved):
- `log <repo>` - git log (default: --oneline -20)
- `show <repo> [ref]` - git show (default: HEAD)
- `fetch <repo>` - git fetch --all --prune
- `switch <repo> <branch>` - git switch <branch>

---

## YAML Gotchas

- **Avoid colons in list items** - even inside quotes, `- echo "Service URL: foo"` becomes `{'echo "Service URL': 'foo"'}`. Use dashes instead.
- **Use `|` for multiline scripts** - avoids escaping issues
- **Validate with Python** - `python3 -c "import yaml; print(yaml.safe_load(open('file.yml')))"` reveals parsing surprises

---

## Pulumi (CRITICAL - READ BEFORE ANY PULUMI OPERATION)

### Backend Rules (NON-NEGOTIABLE)

**Azure Blob is the ONLY allowed backend.** The remote state is the source of truth.

| Allowed | Forbidden (blocked by hookify) |
|---------|-------------------------------|
| `azblob://state?storage_account=...` | `pulumi login file://` |
| Already configured in shell | `pulumi login --local` |
| | `pulumi login -l` |

**Local state files (`~/.pulumi/stacks/`) are caches only.** Never create, rely on, or import from them.

To verify backend: `pulumi whoami -v` should show `azblob://`

### Mandatory Pre-flight Check

**BEFORE ANY PULUMI COMMAND, RUN:**

```bash
devbot pulumi <repo>
```

This shows existing stacks, resources, and warns if infrastructure exists. **NEVER skip this step.**

### Forbidden Commands (blocked by hookify)

| Command | Why It's Dangerous |
|---------|-------------------|
| `pulumi stack init` | Creates new empty stack, orphans existing infrastructure |
| `pulumi destroy` | Deletes all resources |
| `pulumi stack rm` | Removes stack and loses state |
| `pulumi login file://` | Switches to local backend, orphans remote state |
| `pulumi login --local` | Same as above |

### Safe Workflow

```bash
# 1. ALWAYS check state first
devbot pulumi my-repo

# 2. cd to the pulumi directory
cd /path/to/platform

# 3. If "no stack selected" error, select existing stack (DON'T init new one!)
pulumi stack select dev

# 4. Then run preview/up
pulumi preview
pulumi up
```

### Other Rules

- **Never prefix with `PULUMI_CONFIG_PASSPHRASE=""`** - already set in zsh
- **Never run `pulumi stack init`** - blocked by hookify, user must explicitly allow
- **If you see "no stack selected"** → run `pulumi stack ls --all` to find existing stacks
- **If stack doesn't exist in blob** → ask user, don't create automatically

---

## Language-Specific Rules

### TypeScript (packages/)
- ES2022 target, CommonJS output, strict mode
- Node.js >=18.0.0 required
- Run `npm run build` before testing

### Python
- Python >=3.12 required
- Use `uv` for dependency management
- Run `uv sync` before testing
- FastMCP with sub-server composition

---

## General

- Read the repo's `CLAUDE.md` before making changes
- Read global `~/.claude/CLAUDE.md` for user-wide settings
- Run tests after making changes
- Keep changes focused and minimal
- Consider clone repos as read-only references

---

Now continue with your previous task.
