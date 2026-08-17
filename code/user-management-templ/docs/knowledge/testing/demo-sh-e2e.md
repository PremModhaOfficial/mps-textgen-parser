---
title: "demo.sh — 21 E2E Tests with Mock DAL"
category: testing
project: v3
created: 2026-03-20
---

`demo.sh` is the integration test suite. It tests the generated NATS microservice against a Docker Compose stack.

## Architecture
- Connects to NATS on port 4229 (docker compose maps 4229→4222)
- Spawns 5 mock DAL responders: `nats reply "motadata.{entity}.db.>" '{"status":"ok"}'`
- Sends test messages via `nats req` and checks responses
- Traffic observer captures all messages on `motadata.>`

## Test Coverage (21 tests)
| # | Test | What it validates |
|---|------|-------------------|
| 1-3 | User create (valid, weak password, missing fields) | Validation + hook pipeline |
| 4-6 | User get, list, update | CRUD operations |
| 7 | Role create | Entity handler |
| 8-10 | User-Role assign, list, remove | Relation operations |
| 11 | User delete | Delete pipeline |
| 12-17 | Permission CRUD + validation | Full entity lifecycle |
| 18-20 | Role-Permission assign, list, remove | Second relation |
| 21 | Service discovery (`$SRV.INFO`) | NATS micro metadata |

## The Mock DAL Pattern
Instead of a real database, mock DAL responders return `{"status":"ok"}` for any `*.db.*` subject. This tests the business layer in isolation.

## Answers (from interview)
- **demo.sh origin**: AI-generated quick-and-dirty check to verify the service works with validation
- **Why no unit tests?** Generating proper test cases would require a whole new TextGen template that's ~10x more complex than the current Entity TextGen — it would need to generate tests for all combinations and permutations of entity fields and their types. Not worth the effort at this stage.
- **Mock DAL is sufficient** for now — it validates the business logic layer (validation, hooks, routing) without needing a real database. The DAL is a separate downstream service.

Links: [[projects/v3-user-managment-mps]] [[deployment/docker-nats-setup]]
