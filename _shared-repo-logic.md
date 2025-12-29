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
