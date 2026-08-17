---
title: "SQL Schema Generation from DSL"
category: cross-cutting
project: v3
created: 2026-03-20
---

The DSL generates SQL DDL for PostgreSQL from Entity and Relation definitions.

## Generated Schema (`sqlPrem_init_sql.sql`)
```sql
CREATE SCHEMA IF NOT EXISTS iam;

CREATE TABLE IF NOT EXISTS iam.user (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  password VARCHAR(255) NOT NULL,
  user_name VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  user_mail VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS iam.roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  desc TEXT,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS iam.user_roless (
  user_id UUID REFERENCES iam.user(id) ON DELETE CASCADE,
  roles_id UUID REFERENCES iam.roles(id) ON DELETE CASCADE,
  assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, roles_id)
);
```

## What the DSL Maps
| DSL Concept | SQL Output |
|-------------|-----------|
| Entity | TABLE with columns |
| Field(uuid, primaryKey) | `UUID PRIMARY KEY DEFAULT gen_random_uuid()` (PostgreSQL 13+; older versions need `pgcrypto` extension) |
| Field(string, required) | `VARCHAR(255) NOT NULL` |
| Field(time, auto) | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` |
| Field(email, unique) | `VARCHAR(255) NOT NULL UNIQUE` |
| Field(text, nullable) | `TEXT` (no NOT NULL) |
| Relation | Junction table with composite PK + CASCADE |

## Answers (from interview)
- **`iam` schema** — fully configurable in the DSL via the `dbSchema` property on the `SqlSchem` root concept. `iam` was just the chosen value.
- **`user_roless` (double s)** — user mistake, not a TextGen bug. The hardcoded `s` suffix for pluralization wasn't removed. Should use variables instead of hardcoded strings in TextGen.
- **No Permissions table** — forgot to add a reference to the Permissions entity in the `SqlSchem` root's `entityrefs`. The SQL generation itself works fine — just a missing reference in the sandbox configuration.

Links: [[projects/v3-user-managment-mps]] [[dsl-evolution/field-types-evolution]]
