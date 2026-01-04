# Setup Workspace Design

**Date:** 2026-01-04
**Status:** Approved

## Summary

Consolidate `/install-devbot`, `/setup-symlinks`, and `/setup-plugins` into a unified `/setup-workspace` command that handles complete workspace initialization in one idempotent flow.

## Motivation

Current setup requires three separate commands:
- `/setup-symlinks` — Links commands to `~/.claude/commands`
- `/install-devbot` — Builds Go CLI
- `/setup-plugins` — Installs 25+ plugins

Additionally, `config.yaml` must be manually created from the example file. This creates friction for new machine setup and onboarding.

## Design

### Command Signature

```
/setup-workspace [--skip-config] [--skip-plugins]
```

### Step-by-Step Flow

Each step prompts only when work is needed; auto-skips if already complete.

#### Step 1: Configure Workspace

**If no `config.yaml` exists:**
```
→ Step 1/4: Configure workspace

  Scanning ~/code/mono-claude...
  Found 12 repositories:

  │ Repo                    │ Stack         │ Status      │
  │ mango                   │ go, nextjs    │ NEW         │
  │ slash-commands          │ go            │ NEW         │
  │ mlx-hub-claude-plugin   │ typescript    │ NEW         │
  │ ...                     │               │             │

  Generate config.yaml? (yes / edit / skip)
```

**If `config.yaml` exists:**
```
  Found existing config.yaml with 10 repos.
  Discovered 2 new repos not in config:

  │ new-project             │ python        │ NEW         │

  Add these to config.yaml? (yes / skip)
```

Uses `devbot detect` for language/framework detection.

#### Step 2: Install devbot

```
→ Step 2/4: Install devbot CLI

  Checking devbot...
  ✓ Already installed at /Users/sloan/go/bin/devbot (v1.0.0)

  (skip)
```

If not installed:
```
  ✗ Not found in PATH

  Install devbot? (yes / skip)

  → Building from slash-commands/devbot...
  ✓ Installed to ~/go/bin/devbot
```

#### Step 3: Create Symlinks

```
→ Step 3/4: Setup command symlinks

  Checking ~/.claude/commands...
  ✓ 28 symlinks already configured

  (skip)
```

If incomplete:
```
  Found 20 symlinks, 8 missing

  Create missing symlinks? (yes / skip)
```

#### Step 4: Install Plugins

Invokes `/setup-plugins` to keep logic DRY.

```
→ Step 4/4: Install plugins

  │ Category              │ Installed │ Available │
  │ Superpowers           │ 7/7       │ ✓         │
  │ Official Anthropic    │ 9/11      │ 2 new     │
  │ Language Servers      │ 2/4       │ 2 new     │

  Install missing plugins? (yes / skip / choose)
```

#### Completion Summary

```
→ Workspace setup complete!

  ┌─────────────────────────────────────────────────┐
  │ Config      │ ✓ 12 repos configured             │
  │ devbot      │ ✓ Installed (v1.0.0)              │
  │ Symlinks    │ ✓ 28 commands available           │
  │ Plugins     │ ✓ 25 installed                    │
  └─────────────────────────────────────────────────┘

  Quick start:
    /status              Show all repo status
    /super <repo>        Start brainstorming

  NOTE: Restart Claude Code to activate new plugins.
```

### Config Schema Change

Simplify from two paths to one:

**Before:**
```yaml
base_path: ~/code/my-workspace
code_path: ~/code
```

**After:**
```yaml
workspace: ~/code/mono-claude
```

Special directories (`.trees/`, `clones/`) are relative to `workspace`.

### File Changes

| File | Action |
|------|--------|
| `setup-workspace.md` | CREATE — New unified command |
| `setup-plugins.md` | KEEP — Called by setup-workspace, also standalone |
| `install-devbot.md` | DELETE — Logic moves to setup-workspace |
| `setup-symlinks.md` | DELETE — Logic moves to setup-workspace |
| `config.yaml.example` | UPDATE — Use single `workspace:` path |
| `README.md` | UPDATE — Reflect new setup flow |
| `_shared-repo-logic.md` | UPDATE — Use `workspace` instead of `base_path`/`code_path` |

### Idempotency

- Config merge: Only adds new repos, never removes existing
- devbot: Checks `which devbot` before building
- Symlinks: Uses `-sf` flag, only creates missing
- Plugins: Checks installed list before installing

## Implementation Notes

- Workspace detection: Uses Claude session's current directory
- Repo discovery: Leverages existing `devbot` exclusion logic (skips `.trees/`, `clones/`, `node_modules/`)
- Language detection: Uses `devbot detect <path>`

## Open Questions

None — design approved.
