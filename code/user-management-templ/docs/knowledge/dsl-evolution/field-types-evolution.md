---
title: "Field Types — From SQL Primitives to Domain Types"
category: dsl-evolution
project: general
created: 2026-03-20
---

How the DSL's type system evolved across projects.

## v2-UMAN: SQL Primitives
`dataType` enum: `int64`, `string`, `float`, `boolean`, `timestamp`
These map 1:1 to database column types. No domain semantics.

## v3: Domain Types
`FieldType` enum: `uuid`, `string`, `text`, `int`, `bool`, `time`, `json`, `email`, `password`, `entityRef`
These carry business meaning — `email` implies validation and UNIQUE constraint, `password` implies hidden from JSON, `uuid` implies auto-generated primary key.

## What Changed
| v2 Type | v3 Type | What it gained |
|---------|---------|---------------|
| `string` | `email` | UNIQUE constraint, validation |
| `string` | `password` | `hidden` annotation, policy checking |
| `int64` | `uuid` | Auto-generation via `gen_random_uuid()` |
| (none) | `json` | JSONB storage |
| (none) | `entityRef` | Foreign key reference |

## Design Notes
- **Domain types emerged from requirements**: The company needed email validation, password hiding, and UUID auto-generation as standard behaviors. Making them types (not annotations) ensures they're always applied correctly — you can't forget to hide a password if the type enforces it.
- **`entityRef` as a type**: Pragmatic choice — keeps field definitions uniform. An entityRef field is just another column in the struct, but the DSL knows it's a foreign key reference for SQL generation and relation handling.

Links: [[cross-cutting/sql-schema-generation]] [[projects/v2-uman]] [[projects/v3-user-managment-mps]]
