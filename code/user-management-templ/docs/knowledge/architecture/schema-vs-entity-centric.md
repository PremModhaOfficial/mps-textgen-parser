---
title: "Schema-Centric vs Entity-Centric Architecture"
category: architecture
project: general
created: 2026-03-20
---

The biggest architectural pivot between v2-UMAN and v3: how the DSL models the domain.

## v2-UMAN: Schema-Centric
- Top-level concept: `Schema` (a database table)
- Fields defined as columns with SQL types
- Thinking in tables, not domain objects
- `FieldRefrence` for joins (database concept)
- Generated REST APIs and SQL schemas

## v3: Entity-Centric
- Top-level concept: `Entity` (a domain object)
- Fields with domain types (`email`, `password`, `uuid`)
- Thinking in business logic, not storage
- `Relation` as a first-class concept (not a foreign key trick)
- Generated NATS microservice handlers with hooks

## The Shift
v2 → v3 moved from "what does the database look like?" to "what does the business domain look like?"

This is significant because:
- Entity-centric enables hooks (business logic sits on entities, not tables)
- Entity-centric enables NATS subjects that map to domain operations (`motadata.user.create`)
- Schema-centric made SQL generation natural but handler generation awkward

## Answers (from interview)
- **UMAN was a learning POC** following Prisma-like definition structure. It was never meant for production — just learning MPS by mimicking a known pattern.
- **UserManagement is a company task** — an internal productized version with a clear scope: a business logic layer on top of the Data Access Layer (DAL), providing custom fields, validation layers, and hooks.
- **The shift wasn't schema→entity as a design realization** — it was POC→product. The company requirement defined the entity-centric architecture: each entity is a NATS microservice endpoint with its own validation and hook pipeline.

Links: [[projects/v2-uman]] [[projects/v3-user-managment-mps]] [[projects/v1-user-management]]
