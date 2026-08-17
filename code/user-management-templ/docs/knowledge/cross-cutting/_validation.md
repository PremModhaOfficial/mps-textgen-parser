---
title: "Validation Report: SQL Schema Generation"
validator: validator-cross-cutting
validated: 2026-03-20
node: sql-schema-generation.md
---

# Validation Report: `sql-schema-generation.md`

## Summary

All core PostgreSQL facts in this node are **CORRECT**. One version caveat and one naming anomaly flagged.

---

## Fact-by-Fact Results

### 1. `gen_random_uuid()` — UUID PRIMARY KEY

**Status: CONFIRMED** ✓

`gen_random_uuid()` is a **built-in core function** in PostgreSQL 13+. It generates a version 4 (random) UUID and requires **no extension**.

- Before PostgreSQL 13: required the `pgcrypto` extension
- PostgreSQL 13+: built-in, pgcrypto's version internally delegates to core
- The pattern `UUID PRIMARY KEY DEFAULT gen_random_uuid()` is valid and idiomatic

**Version caveat (not noted in the document):** If the target PostgreSQL version is < 13, the `pgcrypto` extension must be enabled first (`CREATE EXTENSION IF NOT EXISTS pgcrypto`). The document should note a minimum version requirement.

**Sources:**
- https://www.postgresql.org/docs/current/functions-uuid.html
- https://www.postgresql.org/docs/13/functions-uuid.html

---

### 2. `TIMESTAMPTZ NOT NULL DEFAULT NOW()`

**Status: CONFIRMED** ✓

`TIMESTAMPTZ` (alias: `TIMESTAMP WITH TIME ZONE`) is a standard PostgreSQL type. `DEFAULT NOW()` sets the column to the current timestamp at insert time. The PostgreSQL wiki explicitly recommends `timestamptz` over plain `timestamp` for most use cases.

The DSL mapping `Field(time, auto) → TIMESTAMPTZ NOT NULL DEFAULT NOW()` is correct.

**Sources:**
- https://www.postgresql.org/docs/current/ddl-constraints.html
- PostgreSQL wiki: Don't Do This (https://wiki.postgresql.org/wiki/Don't_Do_This)

---

### 3. `ON DELETE CASCADE` behavior

**Status: CONFIRMED** ✓

`ON DELETE CASCADE` on a foreign key constraint causes PostgreSQL to automatically delete referencing rows when the referenced row is deleted. This is standard, well-documented DDL behavior.

The junction table `iam.user_roless` uses this correctly on both FK columns:
```sql
user_id UUID REFERENCES iam.user(id) ON DELETE CASCADE,
roles_id UUID REFERENCES iam.roles(id) ON DELETE CASCADE,
```
Deleting a user or a role will cascade-delete the corresponding join rows.

**Sources:**
- https://www.postgresql.org/docs/current/ddl-constraints.html (Section: Foreign Keys)

---

### 4. Junction table pattern with composite primary keys

**Status: CONFIRMED** ✓

The pattern used in `iam.user_roless`:
```sql
PRIMARY KEY (user_id, roles_id)
```
is the canonical PostgreSQL pattern for many-to-many junction tables. Composite PKs span multiple columns; PostgreSQL enforces uniqueness across the combination and automatically creates a unique B-tree index. This overlaps with the FK columns to prevent duplicate assignments.

**Sources:**
- https://www.postgresql.org/docs/current/ddl-constraints.html (Section: Primary Keys, Section: Foreign Keys)

---

### 5. `VARCHAR(255)` as string column type

**Status: CONFIRMED with note** ✓⚠

`VARCHAR(255)` (`CHARACTER VARYING(255)`) is valid PostgreSQL. Strings shorter than 255 chars use only the space needed. Attempts to store longer strings raise an error.

**Note:** PostgreSQL documentation and community best practice (PostgreSQL wiki "Don't Do This") advise preferring `TEXT` over `VARCHAR(n)` unless a length constraint is intentionally enforced. The 255 limit is a MySQL/legacy convention with no special meaning in PostgreSQL. The DSL using `VARCHAR(255)` for all string fields is functional but somewhat arbitrary — if the intent is just to hold typical string data, `TEXT` would be more idiomatic PostgreSQL.

**Sources:**
- https://www.postgresql.org/docs/current/datatype-character.html

---

## Issues / Open Questions (from the document itself)

These were flagged as "Tacit Knowledge Needed" in the source node and remain unresolved by external doc lookup — they require inspection of the DSL/TextGen templates:

| Issue | Status |
|-------|--------|
| `iam` schema — configurable in DSL? | Not resolvable from PG docs; requires DSL source inspection |
| `user_roless` double-s naming | Likely a TextGen bug (string concatenation `user` + `roles` + `s`); not a PG concern |
| No Permissions table in SQL despite Go entity | Requires TextGen source inspection — possibly SQL gen not wired for Permissions entity |

---

## Overall Verdict

**All five PostgreSQL facts are correct.** The schema pattern used is valid, idiomatic (with the caveat on `VARCHAR(255)` vs `TEXT`), and reflects real PostgreSQL behavior.

**One recommended addition to the source doc:**
> Note that `gen_random_uuid()` requires PostgreSQL 13+. For PG 12 and below, add `CREATE EXTENSION IF NOT EXISTS pgcrypto;` before use.
