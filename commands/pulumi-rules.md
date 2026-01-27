---
description: Critical Pulumi infrastructure rules
---

# Pulumi Rules

**READ BEFORE ANY PULUMI OPERATION**

---

## Backend Rules (NON-NEGOTIABLE)

**Azure Blob is the ONLY allowed backend.** The remote state is the source of truth.

| Allowed | Forbidden (blocked by hookify) |
|---------|-------------------------------|
| `azblob://state?storage_account=...` | `pulumi login file://` |
| Already configured in shell | `pulumi login --local` |
| | `pulumi login -l` |

**Local state files (`~/.pulumi/stacks/`) are caches only.** Never create, rely on, or import from them.

To verify backend: `pulumi whoami -v` should show `azblob://`

---

## Mandatory Pre-flight Check

**BEFORE ANY PULUMI COMMAND, RUN:**

```bash
devbot pulumi <repo>
```

This shows existing stacks, resources, and warns if infrastructure exists. **NEVER skip this step.**

---

## Forbidden Commands (blocked by hookify)

| Command | Why It's Dangerous |
|---------|-------------------|
| `pulumi stack init` | Creates new empty stack, orphans existing infrastructure |
| `pulumi destroy` | Deletes all resources |
| `pulumi stack rm` | Removes stack and loses state |
| `pulumi login file://` | Switches to local backend, orphans remote state |
| `pulumi login --local` | Same as above |

---

## Safe Workflow

```bash
# 1. ALWAYS check state first
devbot pulumi my-repo

# 2. cd to the pulumi directory
cd /path/to/platform

# 3. If "no stack selected" error, select existing stack (DON'T init new one!)
pulumi stack select dev

# 4. Then run preview/up
pulumi preview
pulumi up
```

---

## Other Rules

- **Never prefix with `PULUMI_CONFIG_PASSPHRASE=""`** - already set in zsh
- **Never run `pulumi stack init`** - blocked by hookify, user must explicitly allow
- **If you see "no stack selected"** → run `pulumi stack ls --all` to find existing stacks
- **If stack doesn't exist in blob** → ask user, don't create automatically

---

Now continue with your previous task.
