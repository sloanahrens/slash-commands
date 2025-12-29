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

- **Avoid colons in list items** - even inside quotes, `- echo "Service URL: foo"` becomes `{'echo "Service URL': 'foo"'}`. Use dashes instead: `- echo "Service URL - foo"`
- **Use `|` for multiline scripts** - avoids escaping issues
- **Validate with Python** - `python3 -c "import yaml; print(yaml.safe_load(open('file.yml')))"` reveals parsing surprises

---

## General

- Read the repo's `CLAUDE.md` before making changes
- Run tests after making changes
- Keep changes focused and minimal

---

Now continue with your previous task.
