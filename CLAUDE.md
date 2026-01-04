# CLAUDE.md

Slash commands and devbot CLI for multi-repo workspace management.

## Stack

- **devbot/**: Go 1.23+ CLI for parallel workspace operations
- **\*.md**: Markdown slash commands for Claude Code

## Build & Test

```bash
make -C devbot ci         # Full CI: fmt, vet, test, lint, build
make -C devbot test       # Just tests
make -C devbot install    # Install devbot to PATH
```

## Key Patterns

### Command Structure

All repo-targeting commands follow `_shared-repo-logic.md`:

1. Resolve repo with `devbot path <name>` (NEVER construct paths manually)
2. Load context: global CLAUDE.md, repo CLAUDE.md
3. Execute command logic
4. Follow commit rules (no AI attribution)

### devbot: Name vs Path

**Critical distinction:**

| Takes NAME | Takes PATH |
|------------|------------|
| `devbot path`, `status`, `diff`, `branch`, `check`, `make`, `config`, `todos` | `devbot tree`, `devbot stats` |

```bash
# Correct - two separate commands
devbot path my-repo                              # Returns: /path/to/my-repo
devbot tree /path/to/my-repo                     # Use literal path from above

# Wrong - compound commands or manual paths
REPO_PATH=$(devbot path my-repo) && devbot tree "$REPO_PATH"  # ❌
devbot tree ~/code/my-repo                                     # ❌
```

### Dual-Model Evaluation

Commands like `yes-commit`, `find-tasks`, `update-docs` can use local Qwen model:

1. Local model generates first (fast)
2. Claude generates independently
3. Evaluate local against criteria
4. Use local if it passes; append `[local]` to commit messages only

Requires: mlx-hub plugin + Apple Silicon. Falls back to Claude-only if unavailable.

## Files

| File | Purpose |
|------|---------|
| `config.yaml` | User's workspace config (gitignored) |
| `config.yaml.example` | Template for new users |
| `_shared-repo-logic.md` | Shared patterns imported by all commands |
| `devbot/` | Go CLI source |
| `*.md` | Individual slash commands |

## Gotchas

- **config.yaml is gitignored** - users create from example
- **Symlinks go to ~/.claude/commands/** - created by `/setup-workspace`
- **`_` prefix files** are not symlinked (shared logic, not commands)
- **YAML colons** - quote strings containing colons

## Development

To modify commands:
1. Edit the `.md` file directly
2. Symlinks automatically reflect changes
3. Test with `/command-name` in Claude Code

To modify devbot:
1. Edit Go code in `devbot/`
2. Run `make ci` to verify
3. Run `make install` to update binary
