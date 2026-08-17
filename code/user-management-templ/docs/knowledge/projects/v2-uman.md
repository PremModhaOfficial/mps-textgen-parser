---
title: "UMAN — The First Serious Attempt (Prisma-Style, Single-File Problem)"
category: architecture
project: v2-uman
created: 2026-03-20
---

UMAN was the **first** serious attempt at the User Management DSL. It came BEFORE UserManagement/UserManagmentMps.

**Path**: `/home/prem-modha/MPSProjects/UMAN`
**Timeline**: Feb 2 - Mar 3, 2025 (earliest of the main projects)
**Commits**: 5
**Status**: Working generation — abandoned due to architectural mistake

## The Core Mistake: Prisma-Style Single-File Relations
All entities and relations were defined in a single file, similar to how Prisma schema files work. This seemed natural at first but created a problem: the project needed **separate files per entity** for clean design and management of those entities.

The single-file approach made it hard to:
- Generate separate handler files per entity
- Manage entity-specific concerns independently
- Scale the DSL as more entities were added

## Codebase Facts
- Language namespace: `uman`
- TextGen file: 560 KB — largest of all projects
- Dual generators: `generator/` + `generator1/` — pivot happened mid-project
- Has README.md with full DSL specification
- Generated output: WORKING — `.go` and `.sql` in `source_gen/`

## Key Concepts
- `Models` → `Schema` → `Field` (Schema-centric)
- `Field` has `dataType` enum (int64, string, float, boolean, timestamp)
- `Sensitive` boolean on fields (auto-redact in JSON)
- `FieldRefrence` for M:N joins
- NO hook system at all

## The Lesson
The shift from UMAN to UserManagement was driven by **entity separation**: each entity needs its own definition scope, not a shared Prisma-style model file. This led to the Entity-centric architecture where each Entity is a standalone concept with its own fields, operations, hooks, and relations.

Links: [[projects/v1-user-management]] [[projects/v3-user-managment-mps]] [[architecture/schema-vs-entity-centric]]
