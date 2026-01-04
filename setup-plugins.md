---
description: Install all recommended plugins for the workspace
---

# Setup Plugins

Idempotently install all recommended Claude Code plugins for this workspace.

**Arguments**: `$ARGUMENTS` - Optional flag: `--check` (only show status, don't install/update)

---

## Plugin Registry

### Superpowers Marketplace

| Plugin | Description |
|--------|-------------|
| `superpowers` | Core skills: TDD, debugging, brainstorming, collaboration patterns |
| `elements-of-style` | Writing guidance (Strunk & White) |
| `episodic-memory` | Persistent memory across sessions via semantic search |
| `double-shot-latte` | Auto-continue without "Would you like me to continue?" |
| `superpowers-developing-for-claude-code` | Plugin/skill development resources + 42 docs |
| `superpowers-lab` | Experimental: vim, menuconfig, REPLs via tmux |
| `superpowers-chrome` | Chrome DevTools Protocol access (BETA) |

### Local Model Acceleration

| Plugin | Description |
|--------|-------------|
| `mlx-hub` | Local MLX model inference for speed (Apple Silicon only) |

> ⚠️ **Dev Setup Note**: `mlx-hub` is symlinked from `~/code/mono-claude/mlx-hub-claude-plugin`
> to `~/.claude/plugins/mlx-hub`. Since it points to our local dev repo, it should **never be
> reinstalled** - just verify the symlink exists. If missing, the repo needs to be cloned and
> symlinked manually, not installed via `claude plugin add`.

### Official Anthropic

| Plugin | Description |
|--------|-------------|
| `frontend-design` | Avoid generic React/Tailwind aesthetics |
| `feature-dev` | Code architect, explorer, reviewer agents |
| `code-review` | Code review workflow |
| `commit-commands` | Git commit helpers (/commit, /commit-push-pr, /clean_gone) |
| `pr-review-toolkit` | PR review: comment-analyzer, test-analyzer, silent-failure-hunter, type-design-analyzer, code-reviewer, code-simplifier |
| `hookify` | Create custom hooks for Claude Code |
| `plugin-dev` | Plugin development tools |
| `agent-sdk-dev` | Agent SDK development helpers |
| `security-guidance` | Security best practices |
| `explanatory-output-style` | Educational explanations during work |
| `learning-output-style` | Interactive learning with code contributions |

### Language Servers (Official Anthropic)

| Plugin | Description |
|--------|-------------|
| `typescript-lsp` | TypeScript language server |
| `gopls-lsp` | Go language server |
| `pyright-lsp` | Python language server |
| `rust-analyzer-lsp` | Rust language server |
| `swift-lsp` | Swift language server |
| `clangd-lsp` | C/C++ language server |
| `jdtls-lsp` | Java language server |
| `php-lsp` | PHP language server |
| `lua-lsp` | Lua language server |
| `csharp-lsp` | C# language server |

---

## Process

### Step 1: Check Marketplaces

Verify both marketplaces are registered:

```bash
cat ~/.claude/plugins/known_marketplaces.json
```

**If `superpowers-marketplace` missing:**
```bash
claude plugin marketplace add obra/superpowers-marketplace
```

**If `claude-plugins-official` missing:**
```bash
claude plugin marketplace add anthropics/claude-plugins-official
```

### Step 2: Check Installed Plugins

Read current installations:

```bash
cat ~/.claude/plugins/installed_plugins.json
```

Build list of what's already installed vs what needs installing.

### Step 3: Install Missing & Update Existing

**If `--check` flag passed:** Only report status, don't install or update.

For each plugin in the registry:

1. **If not installed** → Install it
2. **If already installed** → Run `claude plugin update` to check for updates

```bash
# Superpowers Marketplace
claude plugin install superpowers@superpowers-marketplace
claude plugin install elements-of-style@superpowers-marketplace
claude plugin install episodic-memory@superpowers-marketplace
claude plugin install double-shot-latte@superpowers-marketplace
claude plugin install superpowers-developing-for-claude-code@superpowers-marketplace
claude plugin install superpowers-lab@superpowers-marketplace
claude plugin install superpowers-chrome@superpowers-marketplace

# Official Anthropic
claude plugin install frontend-design@claude-plugins-official
claude plugin install feature-dev@claude-plugins-official
claude plugin install code-review@claude-plugins-official
claude plugin install commit-commands@claude-plugins-official
claude plugin install pr-review-toolkit@claude-plugins-official
claude plugin install hookify@claude-plugins-official
claude plugin install plugin-dev@claude-plugins-official
claude plugin install agent-sdk-dev@claude-plugins-official
claude plugin install security-guidance@claude-plugins-official
claude plugin install explanatory-output-style@claude-plugins-official
claude plugin install learning-output-style@claude-plugins-official

# Language Servers (install as needed)
claude plugin install typescript-lsp@claude-plugins-official
claude plugin install gopls-lsp@claude-plugins-official
claude plugin install pyright-lsp@claude-plugins-official      # Python
claude plugin install rust-analyzer-lsp@claude-plugins-official # Rust
# Additional LSPs: swift-lsp, clangd-lsp, jdtls-lsp, php-lsp, lua-lsp, csharp-lsp
```

**For already-installed plugins**, check for updates:

```bash
claude plugin update <plugin-name>@<marketplace>
```

### Step 4: Handle mlx-hub (Local Dev Plugin)

**Do NOT reinstall mlx-hub.** It's symlinked to our local dev repo.

1. Check if symlink exists: `ls -la ~/.claude/plugins/ | grep mlx-hub`
2. If symlink exists → Report as "OK (local dev)"
3. If symlink missing → Warn user to set up manually:
   ```bash
   ln -s ~/code/mono-claude/mlx-hub-claude-plugin ~/.claude/plugins/mlx-hub
   ```

### Step 5: Report Results

```
Plugin Setup Complete
=====================

Installed (new):
  - superpowers@superpowers-marketplace (v4.0.2)
  - episodic-memory@superpowers-marketplace (v1.0.15)
  ...

Already up to date:
  - elements-of-style@superpowers-marketplace (v1.0.0)
  ...

Updated:
  - feature-dev (v1.0.0 → v1.0.1)
  ...

Local dev plugins:
  - mlx-hub (symlinked to ~/code/mono-claude/mlx-hub-claude-plugin) ✓

Total: X plugins installed, Y updated

NOTE: Restart Claude Code to activate new/updated plugins.
```

---

## Key Plugins for This Workspace

| Repo Type | Recommended Plugins |
|-----------|---------------------|
| Next.js apps | `frontend-design`, `typescript-lsp` |
| Go projects | `gopls-lsp` |
| Python projects | `pyright-lsp` |
| Rust projects | `rust-analyzer-lsp` |
| All repos | `superpowers`, `episodic-memory`, `mlx-hub`, `pr-review-toolkit` |
| Feature dev | `feature-dev` (code-explorer, code-architect agents) |
| Plugin dev | `superpowers-developing-for-claude-code`, `plugin-dev` |

## Plugin Integration with Slash Commands

| Slash Command | Integrates With |
|---------------|-----------------|
| `/run-tests` | `pr-review-toolkit:code-reviewer` (after tests pass) |
| `/review-project` | `feature-dev:code-explorer` (codebase analysis) |
| `/resolve-pr` | `pr-review-toolkit` agents (understand/fix issues) |
| `/update-docs` | `pr-review-toolkit:comment-analyzer` (verify accuracy) |
| `/push` | `pr-review-toolkit:code-reviewer` (suggested before PRs) |

---

## Examples

```bash
/setup-plugins              # Install missing + update existing plugins
/setup-plugins --check      # Only show status (dry run)
```
