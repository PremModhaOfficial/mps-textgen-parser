---
title: "Hooks as Customer Extensibility — Not Just Developer Convenience"
category: hooks
project: v3
created: 2026-03-20
---

The hook system isn't internal tooling — it's a **customer-facing extensibility mechanism**.

## The Real Requirement
Customers need the ability to:
- Add **custom behavior** on arbitrary operations (create, update, delete, get, list)
- Control **ordering** of that behavior (priority system)
- Choose between **blocking** (sync) and **fire-and-forget** (async) execution
- Have behavior that serves customer needs AND internal platform needs

## Concrete Examples
- **Async audit trail**: fire-and-forget logging after every operation (customer compliance)
- **Sync data enrichment**: enrich requested data per customer requirements before returning (blocking, must complete)
- **Platform needs**: internal hooks for cache invalidation, PII filtering, admin notifications

## Why This Drove the Evolution
- **Boolean hooks** (v1): "has hook: yes/no" — one hook per operation, no ordering, no async. Useless for multiple customers with different needs.
- **Named hooks** (intermediate): can have multiple hooks, but no priority or async control.
- **Named + prioritized + async** (v3): full customer extensibility — any number of hooks, ordered by priority, each independently sync or async.

## The 4 Signatures Serve 4 Use Cases
| Signature | Use Case | Example |
|-----------|----------|---------|
| Pre-sync (returns error) | Validation, authorization | Reject if customer policy violated |
| Pre-async (fire-and-forget) | Logging, notifications | Audit trail before operation |
| Post-sync (transforms response) | Data enrichment, PII filtering | Strip fields per customer config |
| Post-async (fire-and-forget) | Cache invalidation, alerts | Notify admin after user created |

## The Insight
A hook system designed only for developers would be simpler. The complexity (priority, async, 4 signatures) exists because **customers** need to compose arbitrary behavior on shared operations without modifying generated code.

Links: [[hooks/four-hook-signatures]] [[hooks/boolean-hooks-v1]] [[tooling/merge-hooks-story]]
