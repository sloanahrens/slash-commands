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
- **Host**: github.com, gitlab.com, bitbucket.org, etc.

Examples:
| URL | Repo Name |
|-----|-----------|
| `https://github.com/user/my-app.git` | my-app |
| `git@github.com:user/my-app.git` | my-app |
| `https://github.com/user/my-app` | my-app |

### Step 3: Confirm Details

Ask the user to confirm or customize:

```
Adding repository:
  URL: <url>
  Name: <repo-name>
  Group: apps (or devops)
  Aliases: <suggested-alias>

Proceed? (yes/edit/cancel)
```

If user chooses "edit", ask what to change.

### Step 4: Clone Repository

```bash
cd <base_path> && git clone <url>
```

If clone fails, report error and stop.

### Step 5: Update config.yaml

Add new entry to `config.yaml`:

```yaml
  - name: <repo-name>
    group: <apps|devops>
    aliases: [<alias>]
```

### Step 6: Check for CLAUDE.md

```bash
ls <base_path>/<repo-name>/CLAUDE.md
```

If missing, ask:
```
No CLAUDE.md found. Would you like to create one? (yes/no)
```

If yes, create a basic template.

### Step 7: Confirm Success

```
Repository added successfully:
  Location: <base_path>/<repo-name>
  Config: Updated config.yaml

Run `/super <repo-name>` to start working with it.
```

---

## Group Selection

Ask user or infer from repo contents:

| Indicator | Group |
|-----------|-------|
| Contains `pulumi/`, `terraform/`, `infra` | devops |
| Contains `package.json`, `go.mod`, `src/` | apps |
| Unclear | Ask user |

---

## Examples

```bash
/add-repo https://github.com/user/my-new-app.git
/add-repo git@github.com:org/infrastructure.git
/add-repo                                          # Prompts for URL
```
