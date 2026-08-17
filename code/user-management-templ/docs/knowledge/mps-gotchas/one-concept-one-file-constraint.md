---
title: "MPS Core Constraint: One Concept = One Generated File"
category: mps-gotchas
project: general
created: 2026-03-20
---

The single biggest constraint in MPS that shaped the entire DSL architecture.

## The Constraint
In MPS, there is a **1:1 mapping between a root concept instance and a TextGen output file**. If you define a single "Models" root with all entities inside it, TextGen generates ONE file containing ALL entities.

To get **separate files per entity**, each entity must be a **separate root definition** in the sandbox. This means each Entity instance in the MPS sandbox is its own top-level node.

## Why This Matters
- Defining all entities in a single model file (like Prisma does) produces a single generated file
- That single file violates the **S in SOLID** — one file has multiple responsibilities (multiple entity handlers)
- Handlers are generated alongside their entity definition, so they MUST be in the same file
- To get `user.go`, `roles.go`, `permissions.go` as separate files, User, Roles, and Permissions must be separate root definitions

## How This Shaped the Architecture
- UMAN's single-file approach → forced everything into one generated file → abandoned
- UserManagement/UserManagmentMps → each Entity is a separate root concept → generates separate files
- The `Main` concept references entities via `EntityRef` nodes, not by containing them
- Relations are children of Entity (not standalone) because they generate into the same file as their parent entity

## The Workaround Pattern
```
Sandbox:
  Main (root) → main.go
  User (root) → user.go (includes UserRoles relation handlers)
  Roles (root) → roles.go (includes RolesPermissions relation handlers)
  Permissions (root) → permissions.go
  SqlSchem (root) → sqlPrem_init_sql.sql
```

Each root concept = one TextGen = one output file.

## pcrd Also Hit This
`pcrd` (Prem's CRUD DSL) was a working end-to-end solution: entity/relation declaration → Go Gin server + SQL schema + Docker setup. But it had the same single-file problem — all entities in one model produced one giant file.

Links: [[projects/v2-uman]] [[projects/v1-user-management]] [[architecture/schema-vs-entity-centric]]
