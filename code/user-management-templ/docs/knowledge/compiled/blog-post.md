# Building a Production DSL with JetBrains MPS: What I Learned from 21 Failed Experiments

Every entity handler looked identical. Same validation pattern. Same NATS reply structure. Same hook pipeline. Same error handling. 1,790 lines of Go spread across three files, but the cognitive structure was the same repeated three times.

That's when I knew: there had to be a better way.

This is the story of building a production domain-specific language (DSL) with JetBrains MPS—from the initial frustration to the breakthrough to production. Along the way, I'll share what worked, what didn't, and the surprising lesson that the real power of MPS isn't fighting its constraints but *extending them* with tooling.

## The Hook: Too Much Repetition

I was building a NATS microservice that bridges a User Management layer with a Data Access Layer. Each entity (User, Role, Permission) needed handlers for create, read, update, delete, and list operations. Each operation required:

1. A struct definition (the event type)
2. Validation logic (check required fields, password strength)
3. Hook execution (pre-create validation, post-create notification)
4. NATS subject routing (map to the DAL)
5. OTEL span creation (observability)

By the third entity, I realized: I'm copy-pasting handler structure. This isn't engineering; it's bureaucracy. And that's when you need a DSL.

## The Learning Journey: 21 MPS Experiments

I didn't start by building UserManagement. I started by learning MPS from first principles.

The `/home/prem-modha/MPSProjects/` directory contains 21 experimental projects, each teaching me something about language design in MPS:

- **Calculator DSLs** (from tutorials) — learned editor definitions and basic syntax
- **Shape DSLs** — practiced with MPS's expression language
- **gtog** (Go-to-Go transpiler) — deep dive into syntax checking, references, and scoping
- **pcrd** (CRUD DSL) — first working end-to-end example (Gin server + SQL schema + Docker)
- **config languages** (with typos like `congifLang`) — rapid iteration, failed experiments, learning by trying

The typos in project names tell the story: `claculator`, `congifLang`, `testtttt`. Speed over polish. The goal wasn't to ship these; it was to understand MPS's mental model.

By the time I reached UserManagement, I'd internalized:
- How concepts relate to each other
- How MPS's projectional editor works (why copy-paste doesn't work)
- The power of constraints in language design
- The fundamental limitation: one concept = one generated file

That last one would haunt me through the entire project.

## First Real Attempt: UMAN (The Prisma-Style Mistake)

UMAN felt natural. I'd define a single `Schema` file that contained all entities—like Prisma's schema.prisma or GraphQL's schema. One file to rule them all.

```
Models {
  User {
    id: int64 [primaryKey, auto]
    username: string
    email: email [unique]
    password: password [hidden]
    created_at: timestamp [auto]
  }
  Role {
    id: int64 [primaryKey, auto]
    name: string
    permissions: Permission[] [manyToMany]
  }
  Permission {
    id: int64 [primaryKey, auto]
    name: string
    resource: string
  }
}
```

TextGen produced exactly what I asked for: a single generated file with all entity handlers. Problem: that violates the S in SOLID. One file shouldn't have three responsibilities.

I needed separate files: `user.go`, `roles.go`, `permissions.go`.

But MPS has a constraint: **one root concept = one generated file**. If all entities live in a single `Models` root, TextGen generates one file. To get three files, I need three root concepts.

That's the architectural lesson from UMAN. The mistake wasn't the feature set; it was the schema-centric thinking (thinking in tables) instead of entity-centric thinking (thinking in domain objects).

## The Shift: Entity-Centric Architecture

UserManagement/UserManagmentMps (note the typo in the original name—it stuck) flipped the architecture:

- **Top-level concept**: `Entity` (a domain object like User, Role, Permission)
- **Each Entity is a separate root in the sandbox**
- **Relations live as children of entities** (UserRole is a child of User, not a separate root)
- **Handlers are generated WITH their entity** (not in a separate file)

Now TextGen produced the right structure:
- `main.go` (from a Main root) — NATS service setup
- `user.go` (from User root) — User handlers + UserRoles relation handlers
- `roles.go` (from Roles root) — Roles handlers + RolesPermissions relation handlers
- `permissions.go` (from Permissions root) — Permissions handlers
- `sqlPrem_init_sql.sql` (from SqlSchema root) — DDL schema

**The insight**: Don't fight MPS's 1:1 constraint. Work with it. Decompose your domain into separate entities, each as a root concept, each generating its own file.

This wasn't just a technical fix. It changed how I thought about the domain. Instead of "what does the database look like?" (schema-centric), I asked "what are the business objects and their operations?" (entity-centric). That shift unlocked everything else: hooks, NATS subjects, prioritized validation.

## Production v3: The Four-Hook System

By March 2025, UserManagmentMps was production-ready with a sophisticated hook system I'd evolved through iteration:

```go
// Pre-sync: can abort the operation by returning an error
func (s *UserHandler) preCreateValidate(ctx context.Context, span tracer.Span, event *CreateUserEvent) error {
    if event.Password == "" {
        return errors.New("password required")
    }
    return nil
}

// Pre-async: fire-and-forget, runs via workerpool
func (s *UserHandler) preCreateNotifyStakeholders(ctx context.Context, span tracer.Span, event *CreateUserEvent) {
    // send notification without blocking create
}

// Post-sync: transforms response data after DAL reply
func (s *UserHandler) postGetFilterPII(ctx context.Context, span tracer.Span, event *GetUserEvent, data []byte) []byte {
    // strip sensitive fields before returning to client
    return sanitized
}

// Post-async: fire-and-forget after DAL reply
func (s *UserHandler) postCreateNotifyAdmin(ctx context.Context, span tracer.Span, event *CreateUserEvent, data []byte) {
    // audit log, metrics, alerts
}
```

Four signatures emerged naturally from the two axes: **pre/post × sync/async**. Each served a distinct purpose:

- **Pre-sync** validators run first (highest priority), can reject the operation
- **Pre-async** notifications run in parallel, can't block
- **Post-sync** transformers modify the response before it goes to the client
- **Post-async** sidecars (audit, metrics, notifications) run after the response is sent

The hook system was the biggest evolution from v1's boolean flags. Priority ordering meant multiple hooks per operation could coexist without collision.

Generating 30+ hook stubs across four entity handlers seemed impossible—until we solved it with external Python.

## The Hardest Problems: TextGen's Walled Garden

### Problem 1: Char-by-Char Editing

The Entity TextGen template is ~300 lines of `append` statements. Every single line was typed character-by-character in MPS's projectional editor.

```
text gen component for concept Entity {
  file path: "/src/"
  extension: "go"
  file name: node.name.toLowerCase()
  body:
    append {package main\n} ;
    append {import (\n} ;
    append {"github.com/nats-io/nats.go"\n} ;
    // ... 290 more lines like this
}
```

This is painful:
- Can't copy-paste working Go code into the template
- Can't use AI to help write templates (can't export them)
- Can't diff changes in git (the `.mps` files are binary XML)
- Every `{`, `}`, `\n` requires an editor action

I tried everything:
- Copy code snippets, paste into MPS → doesn't work (walled garden)
- Export template to external file, edit in VS Code, import back → MPS doesn't support import
- Use DSL Foundry (a third-party TextGen plugin) → failed; couldn't understand groupings

Then I realized: **the constraint isn't MPS, it's the UX**. And UX can be worked around with tooling.

### Problem 2: Reverse-Engineering the Bridge

I built `parse_textgen.py`—a Python script that:

1. Accepts a human-readable `.md` file that looks like TextGen output
2. Reverse-engineers the MPS internal `.mps` XML format
3. Produces valid `.mps` files that MPS loads as TextGen templates

This took one full day and "a lot of AI tokens" to decode the permutations of how TextGen nodes map to XML. But once working, the workflow became:

```bash
1. Write template.md in any text editor
2. run parse_textgen.py template.md
3. Load the generated .mps file into MPS
4. MPS treats it as if you typed it natively
```

Suddenly, templates became version-controllable (markdown is text, not binary). AI could help write them. Iteration speed 10x'd.

This is the breakthrough insight: **when MPS's editor UX is a blocker, don't find a better plugin—build a bridge that lets you work outside MPS and import results back in**.

### Problem 3: The Hook Stub Extraction

Generated code needed a `userDefinedHooks.go` file with hook implementations. But MPS's 1:1 concept-to-file constraint meant creating a root concept just for hook stubs would clutter the sandbox.

Solution: **generate hooks INSIDE entity files, then extract them with Python**.

`merge_hooks.py` does three things:

1. **Strips hooks from entity files** — looks for `#HOOKS_START` / `#HOOKS_END` markers
2. **Preserves existing implementations** — new hooks are appended, old hooks are kept unless their signature changed
3. **Detects signature changes** — if a hook changed from sync to async (different signature), replace it and warn the user

This solved the preservation problem: when TextGen regenerates entity files, user code isn't lost. It's extracted, merged, and preserved.

```python
def merge_hooks(entity_files, userDefinedHooks):
    for entity_file in entity_files:
        hooks = extract_between_markers(entity_file, "#HOOKS_START", "#HOOKS_END")
        for hook in hooks:
            key = f"{hook.receiver_type}.{hook.name}"
            if key not in userDefinedHooks:
                userDefinedHooks[key] = hook_stub
            elif signature_changed(hook, userDefinedHooks[key]):
                warn(f"Signature changed: {key}")
                userDefinedHooks[key] = hook_stub
            else:
                # Keep user's implementation
                pass
```

## What Actually Worked: "Don't Fight MPS, Extend It"

Three improvements emerged as the philosophy:

### 1. Separate Concepts for Separate Root Files

Each entity = separate root concept. Don't fight the 1:1 constraint; embrace it.

```
Sandbox:
  Main → main.go
  User → user.go
  Roles → roles.go
  Permissions → permissions.go
  SqlSchema → sqlPrem_init_sql.sql
```

Each one generates independently, cleanly.

### 2. External Scripts for Transformation

When MPS's 1:1 constraint gets in the way (like with hooks), generate a combined output and use external scripts to transform it:

- `merge_hooks.py` — extracts hooks from generated files and consolidates them
- `sync.sh` — copies `.go` files from MPS, patches SDK imports, builds Docker

The pattern: **MPS generates → External script transforms → Final output**.

### 3. Use Extensions, Not Workarounds

Future projects should use:
- **DSL Foundry TextGen** — escapes the char-by-char pain (the M2 attempt failed due to lack of guides, not plugin quality)
- **MPS extensions** (mbeddr) — reduces boilerplate for editors, constraints, typesystems

The meta-principle: **use MPS for structure and constraints, use external tools for text manipulation and build pipelines**.

## The Alternatives I Explored (And Why They Didn't Win)

### Rust Declarative Macros

I built a working prototype in `/home/prem-modha/projects/rust/userManMacros` that generates User/Role/Permission structs, CRUD stores, and query builders:

```rust
schema! {
    User {
        name: String,
        email: Email,
        password: Password,
        age: u32,
    }
    indexes { name: String, email: Email }
}

many_to_many!(UserRole, User, Role);
```

This generates:
- Struct with auto-incrementing ID
- `UserStore` with create/read/update/delete/list
- Query builder with filter/order_by/limit/skip
- Cascade delete support
- Domain types with validation

**Why it lost**: No hook system, no NATS integration, no OTEL. The macro system excels at generating Rust code, but the broader DSL features lived in MPS's constraint system (validating field types, enforcing naming conventions, generating NATS subject patterns). You can't do that in macros without significant complexity.

**The insight**: Rust macros solve the *editing* problem (templates are normal Rust code) but not the *structure* problem (defining valid domains, constraints, concepts). MPS excels at structure but fails at editing.

### JSON AST → Jinja Templates

Another idea: MPS generates a JSON AST of the domain, then external Jinja templates produce the Go code.

**Why it didn't happen**: It's a partial solution. It solves the TextGen walled garden (external Jinja templates in any text editor) but requires building a JSON serializer in MPS, then learning Jinja well enough to handle all the permutations of entity fields, hook signatures, and relation types. The `parse_textgen.py` bridge was faster.

## The 21 Experiments in Perspective

Looking back, the 21 MPS experiments weren't failures. They were investment in understanding MPS's mental model before betting the farm on a production project.

Each taught something:
- **Calculator/Shape DSLs** — MPS editor basics
- **gtog** — language design (syntax, references, scoping)
- **pcrd** — proof that end-to-end generation works
- **UMAN** — what NOT to do (single-file schema-centric architecture)

By the time I built v3 (UserManagmentMps), I understood MPS's strengths (structure, constraints, concepts) and weaknesses (projectional editor UX, walled garden). That's why v3 worked: it leveraged the strengths and worked *around* the weaknesses with external tools.

## The Stats (What Production Looks Like)

**Generated Code**:
- 1,790 lines of Go across 4 files
- 30+ hook stubs (pre/post × sync/async)
- 21 e2e integration tests
- Docker Compose deployment (NATS on 4229)

**Tooling**:
- `parse_textgen.py` — 1 day to reverse-engineer and implement
- `merge_hooks.py` — preserves user implementations across regeneration
- `sync.sh` — copies output from MPS, patches SDK paths, builds Docker
- `demo.sh` — 21 tests via `nats req` against mock DAL

**The TextGen Template**:
- ~300 lines of append statements
- ~400 lines of Go generated per entity
- Each line typed char-by-char in MPS's projectional editor

## Five Things I Wish I Knew

### 1. The 1:1 Constraint Is a Feature, Not a Bug

I fought it at first. Then I realized: decomposing your domain into separate root concepts forces clean architecture. Each entity is independent. No God file with three responsibilities.

**For next time**: Design your domain around the constraint. Think in separate concepts, not unified schemas.

### 2. The Walled Garden Can Be Bypassed

MPS's projectional editor prevents copy-paste. But you can work around it with external tooling (parse_textgen.py). Once you realize the constraint is UX, not technical, solutions appear.

**For next time**: Build a text → MPS bridge. Makes template authoring and iteration 10x faster.

### 3. Hooks Are Your Extensibility Layer

Pre/post × sync/async × priority = a powerful customer extension mechanism. Customers can validate, transform, and side-effect on every operation without touching generated code.

**For next time**: Design hooks early. They're not just developer convenience; they're the extensibility contract.

### 4. External Scripts Solve MPS's Structural Limits

When MPS's 1:1 constraint prevents what you need (like consolidating 30 hook stubs), generate combined output and use external Python/Bash to transform it. Don't fight MPS; extend it.

**For next time**: Use scripts for file merging, path patching, dependency management. Let MPS focus on generating correct Go code.

### 5. Learn MPS Before Betting On It

The 21 experiments weren't wasted time. They were essential. By the time I reached UserManagement v3, MPS's mental model was intuitive. Concepts, relations, hooks, constraints—they weren't alien. They were natural.

**For next time**: Budget time for learning. A working DSL without understanding MPS's design philosophy is brittle and unmaintainable.

## Closing: What This Project Proved

You can build a production DSL with MPS. It works. The generated NATS microservice is running, tests pass, hooks are executing, customers can extend it.

But success requires accepting MPS's constraints and extending them with external tooling. The naive approach—fight the projectional editor, demand a text editor, refuse the 1:1 concept-to-file mapping—leads to frustration.

The winning approach: **use MPS for what it's good at (structure, constraints, concepts), and use external tools (Python, Bash, templates) for what it's bad at (text editing, file manipulation, build pipelines)**.

That's the philosophy that took v3 from prototype to production.

The NATS microservice now handles User, Role, and Permission operations. More entities can be added by following the pattern. Hooks are the extensibility layer. The codebase is maintainable because each entity is separate, each responsibility is clear, and the tooling (merge_hooks.py, sync.sh) preserves user code across regeneration.

If you're considering MPS for your DSL: learn the fundamentals, understand the constraints, extend don't fight, and build the bridges that make the UX tolerable.

The payoff is worth it.

---

**Where to go from here**: Read the full knowledge graph at `/docs/knowledge/` for deep dives into TextGen, hooks, the merge strategy, alternatives explored, and lessons learned from each of the 21 experiments.
