# Onboarding Guide: MotaDSLUM — NATS User Management DSL

Welcome to the MotaDSLUM project. This guide will get you up to speed on the architecture, the code generation pipeline, and the critical gotchas that experienced developers on this project have discovered.

---

## Before You Start

### Prerequisites
- **JetBrains MPS** (3.2 or later) — installed and functional
- **Go toolchain** (1.18+) — for compiling generated code
- **Docker & Docker Compose** — for running the NATS+app stack
- **Python 3.9+** — for the `parse_textgen.py` and `merge_hooks.py` scripts
- **NATS CLI** (optional but recommended) — for testing via `nats req`/`nats reply`

### Key Directories
| Path | Purpose |
|------|---------|
| `~/UserManagmentMps` | The DSL source (MPS project) |
| `~/projects/dsl/src/code/user-management-templ/` | Generated code lives here (this repo) |
| `~/projects/nextgen/gosdk/` | Vendored Motadata Go SDK |
| `src/` | Generated `.go` files (overwritten on sync) |
| `src/userDefinedHooks.go` | **ONLY safe file to edit** — hook implementations |
| `sync.sh` | The generation pipeline entry point |
| `merge_hooks.py` | Extracts hooks from generated code |
| `demo.sh` | 21 e2e tests (mock DAL pattern) |

---

## The 5-Minute Overview

### What the DSL Does
This DSL transforms **Domain-Specific Language definitions** (written in MPS) into a complete **NATS microservice** with:
- Entity handlers (User, Roles, Permissions) → CRUD operations on NATS subjects
- Relation handlers (UserRoles, RolesPermissions) → Assignment/removal operations
- SQL schema generation (PostgreSQL)
- Sophisticated hook system for customer extensibility
- OpenTelemetry tracing integration
- Motadata Go SDK integration

### The Generation Pipeline
```
MPS DSL (~/UserManagmentMps/)
    ↓
TextGen generates .go files
    ↓
sync.sh:
  1. Copy files from MPS into src/
  2. Run merge_hooks.py (preserve user implementations)
  3. Patch SDK module paths
  4. Docker compose build + up
    ↓
Running NATS microservice on port 4229 (external) / 4222 (internal)
```

**Key rule**: Never hand-edit generated `.go` files except `userDefinedHooks.go`. All other changes must be made in the MPS DSL, then regenerated via `sync.sh`.

---

## Critical Things You MUST Know

### 1. One Concept = One Generated File (MPS Constraint)
In MPS, there is a **strict 1:1 mapping** between a root concept instance and a TextGen output file.

**Why this matters:**
- If all entities were defined in a single model file (Prisma-style), TextGen would generate ONE file with all handlers
- To get `user.go`, `roles.go`, `permissions.go` as **separate files**, each Entity must be a **separate root definition** in the sandbox
- This constraint **forced the entire architecture**: each Entity is its own top-level node in the sandbox

**The workaround:**
```
Sandbox structure:
  Main (root)        → main.go
  User (root)        → user.go (includes UserRoles handlers)
  Roles (root)       → roles.go (includes RolesPermissions handlers)
  Permissions (root) → permissions.go
  SqlSchem (root)    → sqlPrem_init_sql.sql
```

**Lesson:** Understand this constraint before wondering why the DSL structure looks the way it does.

### 2. TextGen Is a Walled Garden
MPS's projectional editor **cannot exchange content with external editors**:
- **No copy-paste INTO TextGen** — can't paste code from external editors (must type character-by-character)
- **No copy-paste FROM TextGen** — can't copy templates to edit in VS Code
- **No external tooling** — can't use AI or scripts to edit templates directly

**The breakthrough:** The `parse_textgen.py` script (1 day of development) **reverse-engineered** MPS's internal TextGen format, enabling:
- Write templates in readable `.md` format
- Parse to `.mps` (MPS XML format)
- Load back into MPS
- Result: templates are now editable in any text editor, AI-assisted, and diffable in git

**Lesson:** If a tool's UX is a blocker, build a bridge rather than fighting the tool.

### 3. TextGen Has No Debugger
When generated code is wrong, you have no automated debugging tools:
- **No line numbers** — can't trace which TextGen node produced line 47 of output
- **No breakpoints** — can't pause execution mid-generation
- **No step-through** — can't watch the template expand

**Workarounds:**
1. Add comments in TextGen output to trace which section generated code
2. Use `System.out.println()` inside `${ }` blocks (output appears in MPS Messages console)
3. Use "Preview Generated Text" (right-click root node) — fastest feedback loop, no file I/O
4. Binary search the template — comment out sections to isolate bad output

**Lesson:** Expect 4+ iterations when modifying TextGen templates. Plan accordingly.

### 4. Holder Pattern Was Unnecessary
The DSL uses a "Holder" concept pattern (e.g., `HookTypeHooksHolder` containing 1..n `Hook` children) to show different UI for "no hooks" vs "has hooks" states.

**Research finding:** Vanilla MPS can achieve the same result with **visibility queries** on collection cells (no extra concept needed). The Holder was a knowledge gap workaround.

**Lesson:** Don't assume complex patterns are necessary — MPS may already provide what you need through visibility queries or conditional editors.

### 5. Hook Priority Must Be Ordered Correctly
Hooks execute **sorted by priority** (highest priority number first). Critical design rule:

**Always place VALIDATING (sync) hooks at high priority, before async hooks.**

**Why:** Async hooks are fire-and-forget — they cannot be cancelled. If a pre-async hook fires then a later pre-sync validator fails, the async hook **keeps running**. Bad situation.

**Example:**
```go
preCreateValidate (sync, priority 10)  ← Runs first, can abort
postCreateAudit   (async, priority 5)   ← Runs after validation passes
```

---

## How to Add a New Entity

1. **Open the MPS project** (`~/UserManagmentMps`)

2. **Create a new Entity root** in the sandbox:
   - Right-click sandbox → New Child → Entity
   - Name: `YourEntity`
   - Set `name` property to the Go type name (e.g., `User`)

3. **Add fields** to the Entity:
   - Add Field children with:
     - `fieldName`: the Go field name
     - `fieldType`: enum (UUID, String, Time, Text, Email, etc.)
     - Modifiers: `primaryKey`, `auto`, `hidden`, `nullable`, `unique`

4. **Add hooks** (optional):
   - Add Hook children to HookTypeHooksHolder
   - Set `hookName`, `priority`, `isAsync`
   - Naming convention: `pre/post` + OperationName + Description

5. **Add relations** (if needed):
   - Add Relation children to the Entity
   - Link to other entities via `toEntity` reference
   - Define operations (assign, remove, list)

6. **Update the Main root**:
   - Add an `EntityRef` to the Main concept pointing to your new Entity
   - This registers the entity with the microservice

7. **Update SqlSchem**:
   - Add an `entityref` to SqlSchem pointing to your new Entity
   - This includes it in the generated SQL schema

8. **Regenerate**:
   ```bash
   cd ~/projects/dsl/src/code/user-management-templ/
   ./sync.sh
   ```

9. **Verify**:
   ```bash
   ./demo.sh
   ```

---

## How to Modify a Hook

### Understanding Hook Preservation
When you regenerate via `sync.sh`, the generated `*.go` files are **completely overwritten**, but your hook **implementations are preserved** via `merge_hooks.py`.

### The Process

1. **Edit in userDefinedHooks.go** (NOT in generated files):
   ```go
   func (s *UserHandler) preCreateValidate(ctx context.Context, span tracer.Span, event *UserCreateEvent) error {
       // Your implementation here
       return nil
   }
   ```

2. **When changing from sync to async** (or vice versa):
   - Edit the MPS Hook concept: change `isAsync` property
   - Regenerate via `sync.sh`
   - `merge_hooks.py` **detects the signature change**
   - The old implementation is **replaced** with a new stub, and the old code is printed as a warning
   - Port your logic to the new signature

3. **When adding a new hook**:
   - Add a Hook concept in the Entity (set `hookName`, `priority`, `isAsync`)
   - Regenerate via `sync.sh`
   - `merge_hooks.py` **appends the new stub** to `userDefinedHooks.go`
   - Implement the hook body

### The Four Hook Signatures

| Signature | When It Runs | Return Type | Use Case |
|-----------|--------------|-------------|----------|
| Pre-sync | Before operation, blocking | `error` | Validation, authorization — can abort |
| Pre-async | Before operation, non-blocking | `void` | Fire-and-forget logging, audit trail |
| Post-sync | After DAL reply, transforming | `[]byte` | Data enrichment, PII filtering |
| Post-async | After DAL reply, non-blocking | `void` | Cache invalidation, admin notifications |

---

## How to Run and Test

### Prerequisites
- Docker & Docker Compose running
- NATS CLI installed (optional)
- Generated code already compiled via `./sync.sh`

### Running the Service
```bash
cd ~/projects/dsl/src/code/user-management-templ/
./sync.sh        # Generate, build, start Docker Compose
```

The NATS server is now on `localhost:4229` (external) / `nats://nats:4222` (internal).

### Running the Tests
```bash
./demo.sh
```

This runs 21 e2e tests using a **mock DAL pattern**:
- 5 mock responders on `motadata.{entity}.db.>` subjects return `{"status":"ok"}`
- Tests send requests to `motadata.{entity}.{operation}` and verify responses
- Validates: validation, hooks, CRUD routing, relations, service discovery

### Manual Testing
```bash
# Send a request (NATS CLI)
nats req motadata.user.create '{"username":"alice","email":"alice@example.com","password":"secret123"}'

# Listen for responses
nats sub 'motadata.>'
```

---

## Common Mistakes and How to Fix Them

### Mistake 1: Double-S Bug (`user_roless` instead of `user_roles`)
**Symptom:** Junction table name has double suffix (e.g., `user_roless`).

**Root cause:** TextGen hardcoded the `s` suffix for pluralization as a string literal, not a variable.

**Fix:**
- Open the Relation TextGen template in MPS
- Find the hardcoded `"s"` concatenation
- Replace with a variable or computed expression that handles the entity name correctly
- Use `parse_textgen.py` to edit the template as markdown, then parse back to `.mps`

**Lesson:** Always use variables for dynamic values, never hardcoded strings in TextGen.

### Mistake 2: Missing Entity in SQL Schema
**Symptom:** One entity CRUD works in the microservice, but no SQL table is created.

**Root cause:** Forgot to add an `entityref` to the `SqlSchem` root pointing to that entity.

**Fix:**
- Open the `SqlSchem` root in the sandbox
- Add a new `entityref` node pointing to the missing entity
- Regenerate via `sync.sh`

### Mistake 3: SDK Import Errors After `sync.sh`
**Symptom:** Compile errors like `unknown import "motadatagosdk/events/core"`

**Root cause:** `sync.sh` patches the SDK module path from `motadatagosdk` to `dev.azure.com/Motadata/NextGen/motadata-go-sdk`, but the patching may have failed.

**Fix:**
- Check `sync.sh` log output for sed errors
- Manually verify `go.mod` and all `.go` files have the correct module path
- Run `sync.sh` again

### Mistake 4: Hook Signature Mismatch After Changes
**Symptom:** Code compiles but hook doesn't run, or wrong number of parameters.

**Root cause:** Changed hook from sync to async (or vice versa) in MPS, but didn't port the implementation to the new signature.

**Fix:**
- Check `merge_hooks.py` output for "signature mismatch" warnings
- The old implementation is printed — copy it and port to the new signature in `userDefinedHooks.go`
- Example: `error` return type → `void`, add the `[]byte data` parameter for post-hooks

### Mistake 5: Async Hook Fires Even After Sync Validator Fails
**Symptom:** Notification sent even though validation failed.

**Root cause:** Async hook had higher priority than sync validator. Async hooks are fire-and-forget — they can't be cancelled.

**Fix:**
- In the Entity's Hook definitions, **lower the priority number** of validating sync hooks below async hooks
- Rule: Sync validators always run FIRST (higher priority number)

---

## Where to Find Things

### DSL Source Code
**Location:** `~/UserManagmentMps/languages/UserManagement/`

- **Entity definitions:** `models/UserManagement.structure.mps` (structure) + sandbox concepts
- **Hook concepts:** `models/UserManagement.structure.mps`
- **TextGen templates:** `models/*.textGen.mps` files
  - `Entity.textGen_guide.md` — Entity handler template
  - `Relation.textGen_guide.md` — Relation handler template
  - `SqlSchema.textGen_guide.md` — SQL schema template

### Generated Code
**Location:** `~/projects/dsl/src/code/user-management-templ/src/`

- `main.go` (155 lines) — NATS micro.Service setup, endpoints
- `user.go`, `roles.go`, `permissions.go` — Entity handlers
- `userDefinedHooks.go` — **SAFE TO EDIT** — all hook implementations
- `sqlPrem_init_sql.sql` — Generated PostgreSQL schema

### Pipeline Scripts
**Location:** `~/projects/dsl/src/code/user-management-templ/`

- `sync.sh` — Full pipeline: generate → merge → patch → docker compose
- `merge_hooks.py` — Extracts hooks from generated code, merges with user implementations
- `demo.sh` — 21 e2e tests
- `docker-compose.yml` — NATS server + app definition
- `Dockerfile` — App container image

### Vendored SDK
**Location:** `~/projects/nextgen/gosdk/motadata-go-sdk/`

- `events/core/` — Publisher, Subscriber, Message interfaces
- `events/transport/nats/` — NATS implementation
- `otel/tracer/` — OpenTelemetry integration
- `workerpool/` — Async hook execution

### External Documentation
**Knowledge graph:** `~/projects/dsl/src/code/user-management-templ/docs/knowledge/`
- Deep dives on TextGen, hooks, SQL generation, gotchas, and alternatives

---

## Further Reading

### For Understanding MPS Gotchas
- **One Concept = One File:** `docs/knowledge/mps-gotchas/one-concept-one-file-constraint.md`
- **Holder Pattern:** `docs/knowledge/mps-gotchas/holder-concept-pattern.md`
- **TextGen Debugging:** `docs/knowledge/textgen/textgen-debugging-pain.md`

### For Understanding TextGen
- **Walled Garden Problem:** `docs/knowledge/textgen/textgen-walled-garden.md`
- **Python Bridge Solution:** `docs/knowledge/textgen/python-bridge-solution.md`
- **Full Entity Template:** `docs/knowledge/textgen/entity-textgen-template-full.md`

### For Understanding Hooks
- **Four Hook Signatures:** `docs/knowledge/hooks/four-hook-signatures.md`
- **Customer Extensibility:** `docs/knowledge/hooks/hooks-as-customer-extensibility.md`

### For Understanding the Pipeline
- **sync.sh:** `docs/knowledge/tooling/sync-pipeline.md`
- **merge_hooks.py:** `docs/knowledge/tooling/merge-hooks-story.md`
- **Docker Setup:** `docs/knowledge/deployment/docker-nats-setup.md`

### For Testing
- **demo.sh E2E Tests:** `docs/knowledge/testing/demo-sh-e2e.md`

### For Architecture Context
- **v3 (Production):** `docs/knowledge/projects/v3-user-managment-mps.md`
- **v2 (What NOT to do):** `docs/knowledge/projects/v2-uman.md` — Schema-centric, single-file mistake

### For SQL Generation
- **SQL Schema from DSL:** `docs/knowledge/cross-cutting/sql-schema-generation.md`

---

## Quick Checklist: Your First Change

- [ ] Clone/open `~/UserManagmentMps` in MPS
- [ ] Navigate to the sandbox and locate the Entity you want to modify
- [ ] Make your change (new field, new hook, etc.)
- [ ] Run `sync.sh` in `~/projects/dsl/src/code/user-management-templ/`
- [ ] Check `merge_hooks.py` output for warnings
- [ ] If you added/modified a hook, edit `src/userDefinedHooks.go`
- [ ] Run `./demo.sh` to verify
- [ ] Review the generated `src/` files (read-only, for reference)

---

**Welcome aboard.** This project is a sophisticated example of DSL engineering with real constraints. Understand the constraints, and the architecture becomes clear.
