---
description: Add a new repository to the workspace
---

# Add Repository

Clone a repository and add it to the workspace configuration.

**Arguments**: `$ARGUMENTS` - Git URL (HTTPS or SSH). If empty, prompts user.

---

## Process

### Step 1: Get Repository URL

**If `$ARGUMENTS` is empty:**

Ask the user:
```
Enter the git repository URL (HTTPS or SSH):
```

**If `$ARGUMENTS` is provided:**

Use the URL from arguments.

### Step 2: Parse Repository Info

Extract from URL:
- **Repo name**: Last path segment without `.git`
- **Host**: github.com, gitlab.com, etc.

Examples:
| URL | Repo Name |
|-----|-----------|
| `https://github.com/user/my-app.git` | my-app |
| `git@github.com:org/my-app.git` | my-app |

### Step 3: Determine Repository Type

Ask user to classify the repository:

```
What type of repository is this?

1. Reference clone (read-only, for knowledge base)
   → Will be added to clones/clone-config.json
   → Located in <base_path>/clones/

2. Working repo (active development)
   → Will be added to config.yaml repos[]
   → Located in <code_path>/<name>/

Choose (1/2):
```

### Step 4: Confirm Details

**For reference clones:**
```
Adding reference clone:
  URL: <url>
  Name: <repo-name>
  Location: <base_path>/clones/<repo-name>
  Config: clones/clone-config.json
  Description: <ask user>

Proceed? (yes/edit/cancel)
```

**For working repos:**
```
Adding working repository:
  URL: <url>
  Name: <repo-name>
  Location: <code_path>/<repo-name>
  Group: projects
  Config: config.yaml

Proceed? (yes/edit/cancel)
```

### Step 5: Clone Repository

**For reference clones:**
```bash
git clone <url> <base_path>/clones/<repo-name>
```

**For working repos:**
```bash
git clone <url> <code_path>/<repo-name>
```

If clone fails (e.g., SSH access required), report error:
```
Clone failed. This might require SSH access.
Ensure you have SSH keys configured for the host.
```

### Step 6: Update Configuration

**For reference clones** - Update `clones/clone-config.json`:

Add new entry following existing pattern:
```json
{
  "repositories": {
    ...existing...,
    "<repo-name>": {
      "url": "<url>",
      "description": "<user-provided description>"
    }
  }
}
```

**For working repos** - Update `config.yaml`:

Add new entry to `repos[]`:
```yaml
repos:
  - name: <repo-name>
    group: projects
    language: <typescript|python|go|other>
```

### Step 7: Check for CLAUDE.md

```bash
ls <repo-path>/CLAUDE.md
```

If missing for working repos, ask:
```
No CLAUDE.md found. Would you like to create one? (yes/no)
```

If yes, create a basic template.

### Step 8: Confirm Success

**For reference clones:**
```
Reference clone added successfully:
  Location: <base_path>/clones/<repo-name>
  Config: Updated clones/clone-config.json

This repo is for reference only. Use it to:
  - Search for patterns: grep -r "pattern" <base_path>/clones/<repo-name>/
  - Read documentation: cat <base_path>/clones/<repo-name>/README.md
```

**For working repos:**
```
Repository added successfully:
  Location: <code_path>/<repo-name>
  Config: Updated config.yaml

Quick actions:
  /prime-context <repo-name>       Load context
  /find-tasks <repo-name>  Find tasks
  /super-plan <repo-name>       Start brainstorming
```

---

## Group Selection (Working Repos)

Use devbot for fast language detection after cloning:

```bash
devbot detect <repo-path>
# Output: Detected: typescript, nextjs
```

This checks for package.json, go.mod, pyproject.toml, Cargo.toml, etc. and identifies frameworks (Next.js, React, etc.) in ~0.01s.

---

## Linear Integration (Optional)

For working repos, ask:
```
Link to a Linear project? This enables:
  - /find-tasks integration
  - Issue tracking via Linear MCP

Enter Linear project name (or skip):
```

If provided, add `linear_project: "<name>"` to config.yaml entry.

---

## Examples

```bash
/add-repo https://github.com/user/my-new-app.git   # Working repo
/add-repo git@github.com:org/some-sdk.git          # Reference clone
/add-repo                                           # Prompts for URL
```
