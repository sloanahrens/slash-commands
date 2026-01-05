---
name: suggest-prefix
enabled: true
event: bash
pattern: ^cd\s+\S+
action: warn
---

**Use --prefix or -C instead of cd**

| Instead of | Use |
|------------|-----|
| `cd /path; npm run build` | `npm run build --prefix /path` |
| `cd /path; make target` | `make -C /path target` |
| `cd /path; timeout 5 npm run dev` | `timeout 5 npm run dev --prefix /path` |

These patterns:
- Match pre-approved permissions (`Bash(npm:*)`, `Bash(make:*)`)
- Don't require `&&` or `;` which are blocked
- Run without manual approval
