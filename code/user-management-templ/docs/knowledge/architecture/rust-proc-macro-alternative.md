---
title: "Alternative DSL Approach: Rust Declarative Macros (Working Prototype)"
category: architecture
project: experiments
depth: deep
created: 2026-03-20
---

A working alternative to the MPS DSL — Rust **declarative macros** (`macro_rules!`) using `schema!` and `many_to_many!` patterns (with the `paste` crate for identifier concatenation) that generate entity structs, CRUD stores, query builders, and cascade deletes. Note: these are `macro_rules!` (declarative), not `#[proc_macro]` (procedural) — the title was corrected.

**Path**: `/home/prem-modha/projects/rust/userManMacros`

## What It Generates (from ~10 lines of declaration)

```rust
schema! {
    User {
        name: String,
        email: Email,
        password: Password,
        age: u32,
    }
    indexes { name: String, email: Email }
    uniques { email: Email }
}

many_to_many!(UserRole, User, Role);
```

This single declaration generates:
- `User` struct with all fields + auto-incrementing `id`
- `UserStore` with `create`, `read`, `update`, `delete`, `list` methods
- `UserQuery` query builder with `filter`, `order_by`, `limit`, `skip`, `first`, `count`, `exists`, `pluck`, `map`, `print_table`
- `find_by_*` methods for indexed fields
- `find_by_*_id` methods for foreign keys
- `describe()` schema introspection
- Unique constraint enforcement on create/update
- `UserRoleStore` with `link`, `unlink`, `*_ids_for_*` traversal, cascade unlink

## Features
- **Domain types**: `Email` (with validation), `Password` (hashed, hidden in debug)
- **Query builder**: chainable `.filter().order_by().limit().skip().first()`
- **Relations**: `belongs_to { User, Post }` generates FK fields and find methods
- **M:M joins**: `many_to_many!` generates join table with link/unlink/traverse
- **Cascade deletes**: `Database.delete_user_cascade()` removes user + posts + comments + role links
- **Schema introspection**: `describe()` prints table structure with PKs, FKs, indexes, uniques
- **Seed macro**: `seed!(db.users => (...), (...))` for bulk insert

## Comparison with MPS DSL

| Aspect | Rust Macros | MPS DSL |
|--------|------------|---------|
| Declaration syntax | Rust macro invocation | MPS projectional editor |
| Generated output | In-memory Rust structs | Go NATS microservice files |
| Validation | Compile-time (Rust type system) | MPS constraints + typesystem |
| Editing UX | Normal text editor + cargo | MPS projectional (char-by-char) |
| Hook system | None | Named, prioritized, async hooks |
| Distribution | cargo dependency | MPS plugin |
| NATS/messaging | Not present | Core feature |
| SQL generation | Schema introspection only | Full DDL |

## The Insight
Rust macros solve the TextGen UX problem completely — templates are normal Rust code, editable in any editor, with compiler error messages. But they lack the rich DSL editor experience (autocomplete, validation, visual structure) that MPS provides.

The ideal might be a hybrid: MPS for the entity definition UX → generates a Rust macro invocation → Rust compiles to the final output. Or the JSON/Jinja approach: MPS → JSON AST → external templates.

Links: [[textgen/textgen-walled-garden]] [[architecture/json-ast-jinja-idea]] [[projects/v3-user-managment-mps]]
