---
title: "merge_hooks.py — Preserving User Code Across Regeneration"
category: tooling
project: v3
created: 2026-03-20
---

merge_hooks.py exists because of the **MPS 1:1 concept-to-file constraint** — you can't generate `userDefinedHooks.go` without a root concept for it, which would clutter the sandbox with one-line files.

## The Root Cause (MPS Constraint Again)
- `userDefinedHooks.go` needs to be a separate file from entity handlers
- But in MPS, each generated file needs its own root concept in the sandbox
- Creating a root concept just for hook stubs would clutter the sandbox with trivial one-line definitions
- Same problem affected the SQL file generation
- **Solution**: generate hooks INSIDE entity files, then use an EXTERNAL script to extract them

## The Three-Stage Evolution

### Stage 1: Generate hooks in entity files
Hooks are generated between `#HOOKS_START` / `#HOOKS_END` markers inside each entity's generated `.go` file. This works with MPS's 1:1 constraint.

### Stage 2: Extract and merge
`merge_hooks.py` strips hooks from entity files and appends them to `userDefinedHooks.go`:
- New hooks → appended
- Existing hooks → preserved (user implementations kept)
- Identity key: `*ReceiverType.funcName` (avoids cross-entity collisions)

### Stage 3: Signature change detection (added after a real bug)
When a hook was changed from sync to async (or vice versa), the return type changed but the name stayed the same. This **broke the code** because merge_hooks.py kept the old implementation with the wrong signature.

Fix: scan for signature mismatches and replace the stub when the signature changes, printing the old implementation as a warning so the developer can migrate.

## Key Design: Identity by Receiver+Name
```python
# Key: "*UserHandler.preCreateD"
def make_func_key(receiver_type, func_name):
    return f"{receiver_type}.{func_name}"
```

## Stats
- Manages 30+ hook stubs across 4 handler types
- Handles 4 distinct hook signatures correctly
- Signature change detection triggered by real sync↔async changes

Links: [[hooks/four-hook-signatures]] [[tooling/sync-pipeline]] [[mps-gotchas/one-concept-one-file-constraint]] [[projects/v3-user-managment-mps]]
