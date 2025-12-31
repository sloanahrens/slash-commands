---
description: Add a new repository to the workspace
---

# Add Repository (Trabian Branch)

Clone a repository and add it to the workspace configuration, integrating with trabian's structure.

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
- **Host**: github.com, gitlab.com, code.q2developer.com, etc.

Examples:
| URL | Repo Name |
|-----|-----------|
| `https://github.com/user/my-app.git` | my-app |
| `git@github.com:org/my-app.git` | my-app |
| `git@code.q2developer.com:q2/sdk.git` | sdk |

### Step 3: Determine Repository Type

Ask user to classify the repository:

```
What type of repository is this?

1. Reference clone (read-only, for knowledge base)
   → Will be added to clones/clone-config.json
   → Located in ~/trabian/clones/

2. App/project (active development)
   → Will be added to config.yaml repos[]
   → Located in ~/trabian/<name>/

Choose (1/2):
```

### Step 4: Confirm Details

**For reference clones:**
```
Adding reference clone:
  URL: <url>
  Name: <repo-name>
  Location: ~/trabian/clones/<repo-name>
  Config: clones/clone-config.json
  Description: <ask user>

Proceed? (yes/edit/cancel)
```

**For app repos:**
```
Adding app repository:
  URL: <url>
  Name: <repo-name>
  Location: ~/trabian/<repo-name>
  Group: apps
  Config: .claude/commands/sloan/config.yaml

Proceed? (yes/edit/cancel)
```

### Step 5: Clone Repository

**For reference clones:**
```bash
cd ~/trabian/clones && git clone <url>
```

**For app repos:**
```bash
cd ~/trabian && git clone <url>
```

If clone fails (e.g., SSH access required), report error:
```
Clone failed. This might require SSH access.
For Q2/Tecton repos, ensure you have SSH keys configured for the host.
```

### Step 6: Update Configuration

**For reference clones** - Update `clones/clone-config.json`:

```bash
cat ~/trabian/clones/clone-config.json
```

Add new entry following existing pattern:
```json
{
  "repos": [
    ...existing...,
    {
      "name": "<repo-name>",
      "url": "<url>",
      "description": "<user-provided description>"
    }
  ]
}
```

**For app repos** - Update `config.yaml`:

Add new entry to `repos[]`:
```yaml
repos:
  - name: <repo-name>
    path: <repo-name>
    language: <typescript|python|go|other>
    linear_project: <optional - ask if they want to link>
```

### Step 7: Check for CLAUDE.md

```bash
ls <repo-path>/CLAUDE.md
```

If missing for app repos, ask:
```
No CLAUDE.md found. Would you like to create one? (yes/no)
```

If yes, create a basic template following trabian patterns.

### Step 8: Confirm Success

**For reference clones:**
```
Reference clone added successfully:
  Location: ~/trabian/clones/<repo-name>
  Config: Updated clones/clone-config.json

This repo is for reference only. Use it to:
  - Search for patterns: grep -r "pattern" ~/trabian/clones/<repo-name>/
  - Read documentation: cat ~/trabian/clones/<repo-name>/README.md
  - Load knowledge base: /kb/<tag> (if applicable)
```

**For app repos:**
```
Repository added successfully:
  Location: ~/trabian/<repo-name>
  Config: Updated config.yaml

Quick actions:
  /sloan/switch <repo-name>      Switch to this repo
  /sloan/find-tasks <repo-name>  Find tasks
  /sloan/super <repo-name>       Start brainstorming
```

---

## Group Selection (App Repos)

Ask user or infer from repo contents:

| Indicator | Language |
|-----------|----------|
| `package.json` | typescript |
| `pyproject.toml`, `requirements.txt` | python |
| `go.mod` | go |
| `Cargo.toml` | rust |

---

## Linear Integration (Optional)

For app repos, ask:
```
Link to a Linear project? This enables:
  - /sloan/find-tasks integration
  - Issue tracking in /sloan/status

Enter Linear project name (or skip):
```

If provided, add `linear_project: "<name>"` to config.yaml entry.

---

## Examples

```bash
/sloan/add-repo https://github.com/user/my-new-app.git       # App repo
/sloan/add-repo git@code.q2developer.com:q2/new-sdk.git      # Reference clone
/sloan/add-repo                                               # Prompts for URL
```
