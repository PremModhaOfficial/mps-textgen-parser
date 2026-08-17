---
title: "v1 Boolean Hooks — The Starting Point"
category: hooks
project: v1
created: 2026-03-20
---

v1's hook system was a simple boolean: "has pre-hook: yes/no" per operation. No naming, no priority, no async flag.

## What v1 Had
- `preHooks` and `postHooks` as children of Entity
- Binary: either the operation has a hook or it doesn't
- Single hook per operation — no way to compose multiple hooks

## What This Couldn't Do
- Multiple hooks on the same operation (e.g., both validate AND audit on create)
- Async hooks (everything blocked)
- Priority ordering
- Named hooks (couldn't tell them apart in generated code)

## Answers (from interview)
- **When boolean wasn't enough**: The company requirement for customer extensibility demanded multiple named hooks per operation with controllable ordering — boolean "has hook: yes/no" couldn't express this
- **UMAN had no hooks**: Confirmed. UMAN was a Prisma-like learning POC focused on schema/CRUD generation, not business logic extensibility
- **The trigger**: Evolved requirement — customers needed async audit trails alongside sync validation on the same operation. Boolean hooks can't express "both, in this order."

Links: [[hooks/four-hook-signatures]] [[projects/v1-user-management]]
