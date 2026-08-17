---
title: "Validation Report — Tooling Category"
category: tooling
created: 2026-03-20
validator: validator-tooling
nodes_checked: [merge-hooks-story.md, sync-pipeline.md]
---

## Sources Inspected

| File | Path |
|------|------|
| merge_hooks.py | `/home/prem-modha/projects/dsl/src/code/user-management-templ/merge_hooks.py` |
| sync.sh | `/home/prem-modha/projects/dsl/src/code/user-management-templ/sync.sh` |

---

## Fact Checks

### 1. Python regex for Go function signatures
**Claim (merge-hooks-story.md):** Python regex captures receiver type and function name to form identity keys like `*UserHandler.preCreateD`.

**Verdict: ✅ CONFIRMED**

Actual code (merge_hooks.py:36):
```python
FUNC_SIG_RE = re.compile(r'^func\s+\(\w+\s+(\*?\w+)\)\s+(\w+)\s*\(', re.MULTILINE)
```
- Group 1 `(\*?\w+)`: captures receiver type, including optional `*` pointer prefix
- Group 2 `(\w+)`: captures function name
- `re.MULTILINE`: `^` matches start of each line, not just string start
- Key formed as `f"{receiver_type}.{func_name}"` → e.g. `*UserHandler.preCreateD`
- Regex is correct and sufficient for Go method declarations on a single line

---

### 2. #HOOKS_START / #HOOKS_END marker pattern
**Claim (merge-hooks-story.md):** Hooks are generated between `#HOOKS_START` / `#HOOKS_END` markers inside each entity's generated `.go` file.

**Verdict: ✅ CONFIRMED — with minor documentation imprecision**

Actual markers as searched in code (merge_hooks.py:43-44):
```python
start_marker = "// #HOOKS_START"
end_marker   = "// #HOOKS_END"
```
- The docs describe the pattern as `#HOOKS_START/#HOOKS_END` but the actual Go-comment-prefixed form is `// #HOOKS_START` / `// #HOOKS_END`.
- The docs omit the `//` prefix — minor imprecision but the intent is clear.
- Extraction uses string `find()` in a while loop, not regex — handles multiple marker pairs per file.
- This is a standard code-generation marker technique (used in Android tools, ANTLR, etc.).

---

### 3. rsync for SDK syncing
**Claim (sync-pipeline.md):** SDK is synced from `~/projects/nextgen/gosdk/` via rsync.

**Verdict: ✅ CONFIRMED**

Actual code (sync.sh:42):
```bash
rsync -a --delete "$SDK_SRC/" "$PROJECT_DIR/motadata-go-sdk/"
```
- `-a` (archive): recursive, preserves symlinks, permissions, timestamps, owner, group
- `--delete`: removes files in destination that no longer exist in source — ensures clean sync
- Trailing `/` on source is intentional: syncs directory contents, not the directory itself

---

### 4. sed for Go import path rewriting
**Claim (sync-pipeline.md):** `sed -i` rewrites module name in go.mod and import paths in .go files.

**Verdict: ✅ CONFIRMED**

Actual code (sync.sh:44-48):
```bash
sed -i 's|^module motadatagosdk|module dev.azure.com/Motadata/NextGen/motadata-go-sdk|' \
    "$PROJECT_DIR/motadata-go-sdk/go.mod"

find "$PROJECT_DIR/motadata-go-sdk" -name '*.go' -exec \
    sed -i 's|"motadatagosdk/|"dev.azure.com/Motadata/NextGen/motadata-go-sdk/|g' {} +
```
- `sed -i`: in-place edit (GNU sed, standard on Linux)
- `|` delimiter avoids escaping slashes in paths
- `^module` anchors to line start to avoid accidental replacements in comments
- `find ... -exec ... {} +`: batches multiple files per sed call for efficiency

---

### 5. fuser -k for port killing
**Claim:** `fuser -k` kills processes on a port before Docker starts.

**Verdict: ✅ CONFIRMED IN CODE — but ABSENT from sync-pipeline.md documentation**

Actual code (sync.sh:66):
```bash
fuser -k 4229/tcp 8222/tcp 2>/dev/null || true
```
- `fuser -k <port>/tcp`: sends SIGKILL to all processes listening on that TCP port (Linux `psmisc` package)
- `2>/dev/null`: suppresses "no process" errors when port is already free
- `|| true`: prevents `set -euo pipefail` from aborting the script when no process is found
- Ports killed: 4229 (NATS external), 8222 (NATS monitoring)
- **Gap**: sync-pipeline.md lists 7 pipeline steps but does not include the `fuser -k` step. The actual step happens between step 7 (implied Docker) and the docker compose call.

---

## Summary

| Fact | Status | Notes |
|------|--------|-------|
| Python re module parses Go method signatures | ✅ Confirmed | Regex correct; MULTILINE + `\*?` handles pointer receivers |
| `#HOOKS_START`/`#HOOKS_END` marker pattern | ✅ Confirmed | Actual form is `// #HOOKS_START` (Go comment prefix omitted in docs) |
| `rsync -a --delete` for SDK directory sync | ✅ Confirmed | Exact flags match; clean sync semantics correct |
| `sed -i` for Go import path rewriting | ✅ Confirmed | Two-pass: go.mod module name, then all .go import strings |
| `fuser -k` for port killing (Linux) | ✅ Confirmed in code | **Not documented in sync-pipeline.md** — documentation gap |

## Recommended Fix

`sync-pipeline.md` should add a step 7.5 (or append to step 7):
> **Free ports**: `fuser -k 4229/tcp 8222/tcp` kills any process holding NATS or monitoring ports before Docker compose starts, preventing bind errors.
