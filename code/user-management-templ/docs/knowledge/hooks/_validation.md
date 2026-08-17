---
title: "Validation Report — Hooks Category"
validator: validator-hooks
date: 2026-03-20
nodes: [boolean-hooks-v1.md, four-hook-signatures.md, hooks-as-customer-extensibility.md]
---

# Validation Report: Hooks Category

## Summary

| Claim | Status | Evidence |
|-------|--------|----------|
| 4 hook signatures (pre-sync, pre-async, post-sync, post-async) | ✅ VERIFIED | `userDefinedHooks.go` + `user.go` |
| Async hooks use workerpool (fire-and-forget) | ✅ VERIFIED | `user.go` uses `workerpool.AsyncWithCtx()` |
| Pre-sync aborts on error | ✅ VERIFIED | `user.go:86-91` — error returned → `req.RespondError` |
| OTEL `RecordError`, `SetAttributes`, `End` on span | ✅ VERIFIED | `tracer/types.go` Span interface; SDK wraps `go.opentelemetry.io/otel/trace` |
| 4 signatures are idiomatic Go patterns | ✅ VERIFIED | context-first, error return, method receivers — all standard |
| Priority "lower number = earlier" example | ⚠️ DISCREPANCY | Doc's sandbox example contradicts actual generated code order |
| boolean-hooks v1 historical facts | ✅ PLAUSIBLE | Cannot verify against MPS sources, but architecturally consistent |

---

## Detailed Findings

### 1. Four Hook Signatures — VERIFIED

Source: `src/userDefinedHooks.go`, `src/user.go`

All four signatures exist exactly as documented:

```go
// Pre-sync — returns error (aborts on non-nil)
func (s *UserHandler) preCreateD(ctx context.Context, span tracer.Span, event *UserCreatedEvent) error

// Pre-async — no return (fire-and-forget)
func (s *UserHandler) preCreateM(ctx context.Context, span tracer.Span, event *UserCreatedEvent)

// Post-sync — transforms response (returns []byte)
func (s *UserHandler) postGetFilterPII(ctx context.Context, span tracer.Span, event *UserGetRequest, data []byte) []byte

// Post-async — no return (fire-and-forget, data read-only)
func (s *UserHandler) postCreateNotifyAdmin(ctx context.Context, span tracer.Span, event *UserCreatedEvent, data []byte)
```

The concrete examples in the doc (preCreateD, preCreateM, postGetFilterPII, postCreateNotifyAdmin) all appear in the actual code.

---

### 2. Async Hooks Use Workerpool — VERIFIED (with precision note)

The doc says "fire-and-forget via workerpool." The actual implementation uses:

```go
// user.go:80-82
_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
    s.preCreateM(c, span, &event)
})
```

This is **NOT a raw `go func()`** — it uses `workerpool.AsyncWithCtx()` from `motadata-go-sdk/core/pool/workerpool`, which wraps the `ants/v2` goroutine pool library with:
- Bounded queue (default MaxQueueSize: 10000)
- Auto-scaling (min 2 / max 6 workers by default)
- Panic recovery per task
- Context-aware cancellation

**Precision note:** The doc is correct to call it "workerpool" rather than just "goroutine." The distinction matters: tasks may queue under load, not launch immediately.

---

### 3. Pre-Sync Aborts the Operation — VERIFIED

From `user.go:84-91`:
```go
if err := s.preCreateD(ctx, span, &event); err != nil {
    logger.Error(ctx, "hook failed", ...)
    span.RecordError(err)
    _ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
    return
}
```
Pre-sync failure short-circuits the handler: no DAL call is made, error response sent. ✅

Async hooks that were already dispatched before the sync hook failure **are NOT cancelled** — they continue in the workerpool. This edge case is noted as "tacit knowledge needed" in the doc, and the code confirms it.

---

### 4. OTEL Tracer.Span API — VERIFIED (SDK wrapper, not raw OTEL)

Source: `motadata-go-sdk/otel/tracer/types.go`

The project uses a **custom `tracer.Span` interface** that wraps `go.opentelemetry.io/otel/trace.Span`. The methods called in hooks are:

| Method | On SDK Span | Maps to OTEL |
|--------|-------------|--------------|
| `span.End()` | ✅ exists | `trace.Span.End()` |
| `span.RecordError(err)` | ✅ exists | `trace.Span.RecordError(err error, options ...EventOption)` |
| `span.SetAttributes(...)` | ✅ exists | `trace.Span.SetAttributes(kv ...attribute.KeyValue)` |
| `hookSpan.SetOK()` | ✅ exists (SDK-only) | calls `SetStatus(codes.Ok, "")` internally |
| `span.SetError(err)` | ✅ exists (SDK-only) | calls `SetStatus(codes.Error, ...)` + `RecordError(err)` |

**Key distinction:** `SetOK()` and `SetError()` are SDK convenience methods, not standard OTEL Span interface methods. The docs don't claim they are OTEL standard — this is fine.

Official OTEL Go source: https://pkg.go.dev/go.opentelemetry.io/otel/trace#Span

---

### 5. Hook Signatures Are Idiomatic Go — VERIFIED

All four signatures follow established Go conventions:
- `context.Context` as first parameter ✅
- `error` return for cancellable operations ✅
- Method receivers on handler struct (`*UserHandler`, `*RolesHandler`, etc.) ✅
- `[]byte` transformation pattern (post-sync) is idiomatic for middleware pipelines ✅
- Fire-and-forget via managed pool (not naked goroutines) is production-safe ✅

The pattern is consistent with how Go middleware and hook systems are typically designed (e.g., `net/http` middleware chains).

---

### 6. Priority System Example — ⚠️ DISCREPANCY

**What the doc claims** (`four-hook-signatures.md`):
> "filterPII priority 3, cacheResult priority 2, enrichData priority 1 → Execution order: enrichData → cacheResult → filterPII"

This implies **lower priority number = runs earlier**.

**What the generated code shows** (`src/user.go HandleGet`):
```go
responseData = s.postGetFilterPII(...)           // runs 1st
workerpool.AsyncWithCtx(...postGetCacheResult...)  // dispatched 2nd
responseData = s.postGetEnrichData(...)           // runs 3rd
```

Execution order in code: **filterPII → cacheResult → enrichData** — the OPPOSITE of the doc's example.

**Assessment:** The doc says its example is "from the sandbox," suggesting it may describe a different configuration than the production-generated code. However, this creates confusion: either (a) the sandbox has different priorities than the generated project, or (b) the priority numbering description ("lower = earlier") is inverted.

**Recommendation:** Clarify whether the priority example in the doc is:
1. From a separate sandbox with different priority values, OR
2. Incorrectly ordered relative to the actual generated code

This requires checking the MPS DSL definitions where priorities are set per hook.

---

### 7. Boolean Hooks v1 Historical Claims — PLAUSIBLE (unverifiable from code)

The `boolean-hooks-v1.md` describes the v1 system as:
- Binary preHooks/postHooks per operation (yes/no)
- No naming, no priority, no async

These are historical/architectural claims about an earlier DSL version. The MPS source files (`.mps`) are in the working tree but not human-readable XML. The claims are **architecturally consistent** with what you'd expect before adding named, prioritized, async hooks. Cannot definitively verify without examining the MPS model files.

---

### 8. core.Request Interface (NATS micro) — VERIFIED (custom SDK, not nats.go/micro directly)

The project uses `core.Request` from the SDK (`events/core/request.go`), NOT `nats.go/micro.Request` directly. The SDK's `core.Request` interface includes:
- `Context() context.Context`
- `Data() []byte`, `Headers() nats.Header`, `Header(key string) string`
- `Respond(data []byte, opts ...RespondOption) error`
- `RespondError(code, description string, data []byte, opts ...RespondOption) error`

This is a superset of the standard NATS micro pattern (adds JetStream ack/nak, options pattern). The hooks docs don't directly claim to describe `nats.go/micro.Request`, so no contradiction.

---

## Issues Requiring Action

| Issue | Severity | Action |
|-------|----------|--------|
| Priority example (enrichData=1 → first) contradicts generated code (filterPII → first) | MEDIUM | Clarify with DSL author; check MPS priority values for get hooks |
| v1 boolean hooks claims unverifiable from code | LOW | Accept as author's tacit knowledge; tag as needs-MPS-verification |
| Async hooks may queue under pool load (not instant dispatch) | LOW | Add nuance to doc if precision matters for SLA reasoning |
