# Use PWD Command

**Purpose**: Ensures Claude always checks the current working directory before executing commands.

## Instructions

Before proceeding with any task that involves:
- Running bash commands
- File operations
- Path-dependent operations
- Testing or build commands

**ALWAYS run `pwd` first** to verify your current working directory.

## Pattern

```bash
# 1. Check where you are
pwd

# 2. Then proceed with the actual task
cd /correct/path && npm test
```

## Why This Matters

- Prevents running commands in the wrong directory
- Avoids file operation errors due to incorrect paths
- Ensures commands execute in the expected context
- Makes debugging easier by showing the starting location

## After Running PWD

Continue with whatever task you were originally planning to do, but now with the confidence that you know your current location in the filesystem.
