# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Overview

Meta-tooling for the mono-claude workspace:
- **Slash commands** (`.md` files) → Claude Code skills
- **devbot CLI** (`devbot/`) → Fast parallel repo operations

**Setup:** Run `/setup-workspace` (or see README.md for manual steps)

## Development

```bash
cd devbot
make test           # Run tests
make ci             # Full CI: fmt, vet, test, lint
make install        # Install to $GOPATH/bin

# Single test
go test ./internal/workspace/... -v -run TestStatusParsing
```

See [devbot/README.md](devbot/README.md) for full command reference.

## Architecture

```
slash-commands/
├── *.md                    # Slash commands
├── _shared-repo-logic.md   # Shared patterns
├── config.yaml             # Workspace config (gitignored)
├── docs/
│   ├── config/             # Hookify rules → ~/.claude/
│   ├── root-claude/        # Global CLAUDE.md template
│   ├── patterns/           # Versioned knowledge
│   └── templates/          # Subagent prompts
└── devbot/                 # Go CLI (see devbot/README.md)
```

## Key Patterns

**Slash Commands:**
- Each `.md` file is a skill Claude Code can invoke
- `_shared-repo-logic.md` contains common patterns
- Use `devbot` for git operations, not raw git

**devbot Organization:**
- One package per command under `internal/`
- `internal/workspace/` handles config + repo discovery
- Tests use `testdata/` fixtures

## Testing Changes

1. `make test` in `devbot/`
2. Verify: `./devbot status`, `./devbot check <repo>`
3. Test slash commands by invoking in Claude Code
