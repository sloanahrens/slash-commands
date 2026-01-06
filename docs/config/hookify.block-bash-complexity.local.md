---
name: block-bash-complexity
enabled: true
event: bash
pattern: (&&|;\s*\w|(?<!\S)\$\(|\|\s*(head|tail|grep|wc|sort|awk|sed))
action: block
---

**Simplify this command**

**Compound commands - run separately:**
- `cmd1 && cmd2` - run cmd1, then cmd2
- `cmd1; cmd2` - run cmd1, then cmd2
- `$(cmd)` - run cmd first, use output

**Pipes - use Claude tools:**
- `| grep` - Grep tool
- `| head/tail` - Read tool with limit/offset
- `| wc/sort/sed/awk` - process in response

**Directory commands - use flags:**
- `cd /path && npm ...` - `npm run ... --prefix /path`
- `cd /path && make ...` - `make -C /path ...`

**Git commands - use devbot:**
- `devbot status/diff/branch/log/show/fetch <repo>`
