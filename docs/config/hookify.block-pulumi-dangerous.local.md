---
name: block-pulumi-dangerous
enabled: true
event: bash
pattern: pulumi\s+(stack\s+init|login\s+(file://|--local|-l\b)|destroy|stack\s+rm)
action: block
---

**Dangerous Pulumi operation blocked**

**Before ANY Pulumi command, run:**
```
devbot pulumi <repo>
```

**Blocked operations:**
| Command | Risk |
|---------|------|
| `pulumi stack init` | Orphans existing infrastructure |
| `pulumi destroy` | Deletes all resources |
| `pulumi stack rm` | Loses state permanently |
| `pulumi login file://` | Abandons remote state |
| `pulumi login --local` | Abandons remote state |

**Safe workflow:**
1. `devbot pulumi <repo>` - check state first
2. `cd /path/to/platform`
3. `pulumi stack select <env>` - select existing stack
4. `pulumi preview` - verify changes
5. `pulumi up` - apply

**If "no stack selected":** run `pulumi stack ls`, then `pulumi stack select`.
**Never** `stack init` if stacks already exist.

Backend must be: `azblob://` - verify with `pulumi whoami -v`
