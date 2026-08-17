---
title: "pcrd — Prem's CRUD DSL (Working End-to-End, Same File Problem)"
category: experiments
project: experiments
created: 2026-03-20
---

A fully working DSL for generating CRUD servers — but hit the same 1:1 file constraint as UMAN.

## What pcrd Generated (End-to-End)
- Entity and Relation model declarations
- Go server for CRUD operations (using Gin framework)
- Complete SQL schema generation
- Database connection setup
- Docker Compose file + Dockerfile
- Full end-to-end: DSL definition → running containerized server

## Why It Matters
pcrd proves the DSL concept **works** — you CAN generate a complete server from declarations. The problem wasn't the generation capability, it was the **file organization** forced by MPS's 1:1 concept-to-file mapping.

## The Same Problem
Like UMAN, pcrd defined all entities in a single model → generated one giant file → violated single responsibility. The shift to UserManagement's "one Entity root per file" solved this.

## Answers (from interview)
- **Timeline**: pcrd came **before** UMAN. Order was: pcrd → UMAN → UserManagement/UserManagmentMps
- **Docker influence**: Yes — UserManagement reuses the Docker pattern from pcrd, but handled via external scripts (sync.sh) rather than generated from the DSL
- **Gin influence**: pcrd generated Gin (HTTP) servers; UserManagement switched to NATS microservices. The handler pattern is different but the entity/CRUD generation concept carried over

Links: [[mps-gotchas/one-concept-one-file-constraint]] [[projects/v2-uman]] [[projects/v1-user-management]]
