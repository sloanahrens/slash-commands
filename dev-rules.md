---
description: Remind Claude of workspace development rules
---

# Dev Rules

Continue whatever you were doing, but remember these rules:

---

## Path Safety

- **Run `pwd` before bash commands** - verify current location before file/path operations
- **Use absolute paths** - always use full paths from your configured `base_path`
- **Stay within workspace** - never navigate above your configured `base_path`

---

## File Creation

- **NO `/tmp` files** - create temporary/working files in a `docs/` directory within the workspace
- **Prefer editing over creating** - modify existing files when possible

---

## Commit Messages

- **NO** Claude/Anthropic attribution
- **NO** co-author lines
- **NO** "generated with" tags
- Use imperative mood ("Add feature" not "Added feature")
- Keep summary under 72 characters

---

## YAML Gotchas

- **Quote strings with colons** - `echo "Service URL:"` not `echo Service URL:` (YAML interprets unquoted colons as key-value separators)
- **Use `|` for multiline scripts** - avoids escaping issues

---

## General

- Read the repo's `CLAUDE.md` before making changes
- Run tests after making changes
- Keep changes focused and minimal

---

Now continue with your previous task.
