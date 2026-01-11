---
tags: [bash, devbot, cd, exec]
repos: [all]
created: 2026-01-11
updated: 2026-01-11
---

# Running commands in repository directories

## Problem

Need to run a command (npm, make, go, etc.) in a specific repo's directory, but:
- Hookify blocks `cd /path && command` compound patterns
- Manual path construction is error-prone
- Some repos have `work_dir` config that changes the target directory

## Solution

Use `devbot exec` to run commands in repo directories:

```bash
devbot exec <repo-name> <command>
```

For monorepo subprojects or explicit subdirectories:

```bash
devbot exec <repo-name>/<subdir> <command>
```

To run in repo root (ignoring `work_dir` config):

```bash
devbot exec <repo-name>/ <command>    # Trailing slash = repo root
```

## Why

1. **Hookify safety**: Compound commands (`cd && cmd`) are blocked to prevent unintended side effects
2. **work_dir awareness**: `devbot exec` respects `work_dir` from config.yaml automatically
3. **Path resolution**: No need to remember or construct full paths

## Examples

| Instead of | Use |
|------------|-----|
| `cd /path/to/atap-automation2/nextapp && npm test` | `devbot exec atap-automation2 npm test` |
| `cd /path/to/mango/go-api && go build` | `devbot exec mango/go-api go build` |
| `cd /path/to/slash-commands/devbot && make test` | `devbot exec slash-commands/devbot make test` |

## Directory resolution order

1. If `/subdir` specified → `{repo_path}/{subdir}`
2. If trailing slash (`repo/`) → repo root (ignores work_dir)
3. If `work_dir` in config.yaml → `{repo_path}/{work_dir}`
4. Otherwise → `{repo_path}`

## Related

- `devbot path <repo>` — Get the full path if you need it for other tools
- `devbot prereq <repo>` — Validate tools and deps before running commands
