---
title: "The Four Hook Signatures"
category: hooks
project: v3
created: 2026-03-20
---

v3 generates 4 distinct hook signatures, each serving a different lifecycle role.

## The Signatures (from `userDefinedHooks.go`)

**Pre-sync** — can abort the operation by returning an error:
```go
func (s *Handler) hookName(ctx context.Context, span tracer.Span, event *Event) error
```

**Pre-async** — fire-and-forget via workerpool, cannot abort:
```go
func (s *Handler) hookName(ctx context.Context, span tracer.Span, event *Event)
```

**Post-sync** — transforms the response data, runs after DAL reply:
```go
func (s *Handler) hookName(ctx context.Context, span tracer.Span, event *Event, data []byte) []byte
```

**Post-async** — fire-and-forget after DAL reply, receives response data read-only:
```go
func (s *Handler) hookName(ctx context.Context, span tracer.Span, event *Event, data []byte)
```

## Examples from User Entity
- `preCreateD` (pre-sync): validates before create
- `preCreateM` (pre-async): fires notification without blocking
- `postGetFilterPII` (post-sync): strips sensitive fields from response
- `postCreateNotifyAdmin` (post-async): sends admin notification after create

## How Priority Works
Hooks execute sorted by priority value. From the sandbox:
- `filterPII` priority 3, `cacheResult` priority 2, `enrichData` priority 1

**Note:** The generated code in `user.go HandleGet` runs filterPII → cacheResult → enrichData (highest priority number first). The exact sort order (ascending vs descending) depends on the `getHooksSorted()` behavior implementation in the DSL — verify against your MPS behavior definition.

## Answers (from interview)
- **Why 4 signatures?** They emerged naturally from the two axes: pre/post × sync/async. Each combination serves a distinct use case. No single-hook-with-options pattern was considered — the separation is intentional.
- **Async hooks and failure:** By design, async hooks are fire-and-forget. If a pre-async hook fires and then a later pre-sync hook fails, the async hook **keeps running** — it cannot be cancelled. **Design rule:** always place validating sync hooks at higher priority (lower number) so they run BEFORE async hooks. If you put an async hook before a sync validator, the async hook fires even if validation will fail.
- **Priority system:** Was not in v1 (v1 had boolean hooks only). Emerged when the customer extensibility requirement demanded multiple hooks per operation with controlled ordering.

Links: [[hooks/boolean-hooks-v1]] [[projects/v3-user-managment-mps]] [[tooling/merge-hooks-story]]
