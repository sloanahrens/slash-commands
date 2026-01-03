---
description: Build and install devbot CLI
---

# Install Devbot

Build and install the devbot CLI from slash-commands/devbot/.

---

## Process

1. Navigate to devbot directory and run make install:
   ```bash
   cd ~/code/slash-commands/devbot && make install
   ```

2. Verify installation:
   ```bash
   which devbot && devbot --help | head -5
   ```

---

## When to Use

- After cloning slash-commands for the first time
- After pulling updates that changed devbot source
- If `devbot` command is not found

---

## Output

```
Installing devbot...
go install ./cmd/devbot

Installed: /Users/<user>/go/bin/devbot

devbot - Fast parallel workspace operations
Usage: devbot [command]
```
