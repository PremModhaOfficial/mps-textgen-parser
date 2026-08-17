---
title: "v3: UserManagmentMps — Production NATS Microservice"
category: architecture
project: v3
created: 2026-03-20
---

The final, production-grade DSL. Generates a complete NATS microservice with SDK integration, OTEL tracing, and a sophisticated hook system.

**Path**: `/home/prem-modha/UserManagmentMps`
**Timeline**: Mar 12 - Mar 19, 2025 (active, latest)
**Commits**: 11+ (on `prem-dev-pocs` branch)
**Status**: Production — fully functional, 21 e2e tests passing

## Generated Output (1,790 LOC)
- `main.go`: 155 lines — NATS micro.Service setup, SDK connection, endpoint wiring
- `user.go`: 595 lines — User entity + UserRoles relation handlers
- `roles.go`: 481 lines — Roles entity + RolesPermissions relation handlers
- `permissions.go`: 307 lines — Permissions entity handlers
- `userDefinedHooks.go`: 221 lines — 30+ hook stubs (preserved across regeneration)
- `sqlPrem_init_sql.sql`: 31 lines — IAM schema with 3 tables

## Key Concepts (Evolved)
- `Entity` with named, prioritized, async hooks (`Hook` concept with `priority` + `isAsync`)
- `Relation` with `extraFields` support
- `NatsServer` concept (new — microservice-specific)
- 4 hook signatures: pre-sync, pre-async, post-sync, post-async
- Motadata Go SDK integration (events, core, tracer, logger, workerpool)
- OTEL tracing built into every handler

## What's Notable
- Returns to Entity-centric architecture (like v1, not v2-UMAN)
- Hook system is the biggest evolution: boolean → named → prioritized + async
- Python TextGen bridge (`parse_textgen.py`) — new tooling not in v1/v2
- `merge_hooks.py` preserves user implementations across regeneration
- Has 4 TextGen guide files documenting the generation templates

screenshot: images/entity-user-definition-1.png

Links: [[projects/v1-user-management]] [[projects/v2-uman]] [[hooks/four-hook-signatures]] [[tooling/merge-hooks-story]] [[tooling/sync-pipeline]] [[textgen/entity-textgen-template-full]]
