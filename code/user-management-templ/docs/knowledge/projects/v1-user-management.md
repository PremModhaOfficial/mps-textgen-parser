---
title: "UserManagement — The Main Project (Same as UserManagmentMps)"
category: architecture
project: v1
created: 2026-03-20
---

UserManagement (in MPSProjects/) and UserManagmentMps (in ~/) are the **same project**. The repo was recloned to fix MPS errors that couldn't be found in the original.

**Path**: `/home/prem-modha/MPSProjects/UserManagement` (original) + `/home/prem-modha/UserManagmentMps` (reclone)
**Timeline**: Feb 9 - Mar 19, 2025 (active across both copies)
**Commits**: 8 (original) + 11 (reclone)
**Status**: Production — the reclone (`UserManagmentMps`) is the active version

## Why the Reclone?
MPS had errors that couldn't be tracked down in the original project. Rather than debug MPS's internal state, it was faster to reclone and start fresh. This is a common MPS pattern — sometimes the IDE state gets corrupted.

## Key Evolution (across both copies)
- Entity-centric architecture (learned from UMAN's single-file mistake)
- Hook system evolved: boolean flags → named hooks → prioritized + async hooks
- TextGen heavily iterated (4+ attempts in git: `textgen 1` → `textgen 4` → `textgen micros`)
- Python TextGen bridge (`parse_textgen.py`) added in reclone
- SDK + OTEL integration added in reclone

## Codebase Facts (reclone = production version)
- 1,790 lines of generated Go code
- 30+ hook stubs managed by `merge_hooks.py`
- 21 e2e tests via `demo.sh`
- Docker Compose deployment

Links: [[projects/v2-uman]] [[projects/v3-user-managment-mps]] [[textgen/textgen-walled-garden]]
