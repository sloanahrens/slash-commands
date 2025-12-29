# Shared Repo Logic

This file contains shared patterns used by all repo-targeting slash commands.

---

## Configuration

Commands use `config.yaml` for repository definitions. Copy `config.yaml.example` to `config.yaml` and customize for your environment.

```yaml
base_path: ~/code/my-workspace
repos:
  - name: my-app
    group: apps
    aliases: [app]
    language: typescript        # optional: typescript | go | python | rust | shell
    work_dir: src               # optional: subdirectory for commands
    commands:                   # optional: override default commands
      test: npm test
      lint: npm run lint
```

---

## Critical Rule

**CRITICAL**: Always stay within your configured `base_path` - never navigate above this directory.

---

## Repo Discovery

Parse `config.yaml` in this commands directory for repository definitions:

| Group | Description |
|-------|-------------|
| `devops` | DevOps/Infrastructure repos |
| `apps` | Application repos |

---

## Language Detection

If `language` is not specified in config, detect from files:

| File Found | Language | Default Commands |
|------------|----------|------------------|
| `package.json` | typescript | `npm run lint`, `npx tsc --noEmit`, `npm run build`, `npm test` |
| `go.mod` | go | `golangci-lint run`, `go build ./...`, `go test ./...` |
| `pyproject.toml` or `requirements.txt` | python | `ruff check .`, `mypy .`, `pytest` |
| `Cargo.toml` | rust | `cargo clippy`, `cargo build`, `cargo test` |
| `Makefile` only | shell | `make lint`, `make build`, `make test` |

Commands can be overridden per-repo in `config.yaml`.

---

## Repo Selection

**If `$ARGUMENTS` is empty:**

Display grouped list and ask user to select:

```
Select a repository:

DevOps/Infrastructure:
  1. my-infra-pulumi
  2. my-terraform

Apps:
  3. my-nextjs-app
  4. my-api

Enter number or name:
```

**If `$ARGUMENTS` is provided:**

Fuzzy match against directory names and configured aliases:

| Input | Matches (example) |
|-------|-------------------|
| `pulumi` | my-infra-pulumi |
| `app` | my-nextjs-app |
| `api` | my-api |

---

## Commit Rules

When committing changes in any repo:

- **NO** Claude/Anthropic attribution
- **NO** co-author lines
- **NO** "generated with" tags
- Use imperative mood ("Add feature" not "Added feature")
- Keep summary under 72 characters

---

## Standard Process Start

1. **Apply dev rules** → `/dev-rules` (path safety, file creation, commit rules)
2. Parse `config.yaml` for base path and repo definitions
3. If `$ARGUMENTS` empty → show selection prompt
4. If `$ARGUMENTS` provided → fuzzy match to repo
5. Confirm selection: "Working on: <repo-name>"
6. Read `<repo>/CLAUDE.md` for repo-specific guidance
