# DSL Code Generation Journey: From Boilerplate to Production

---

## Slide 1: Title Slide

**Building a Code Generation DSL**

*From 1,790 Lines of Boilerplate to Automated Generation*

**Project**: UserManagmentMps — NATS Microservice Code Generator
**Timeline**: Feb 2 - Mar 19, 2025 (6 weeks)
**Status**: Production ✅ (21 e2e tests passing)

**Speaker notes**: Today we're walking through a 6-week journey building a Domain-Specific Language (DSL) using JetBrains MPS. We started with the problem of repetitive boilerplate code in NATS microservices and ended with a fully functional code generator. This talk covers the architecture decisions, technical challenges, and lessons learned.

---

## Slide 2: Agenda

**What We'll Cover**

- **The Problem**: Why we built a DSL (1,790 lines of hand-written code = pain)
- **The Journey**: Three iterations across two projects — what changed each time
- **Key Decisions**: The architectural pivots that shaped the final design
- **What We Built**: Artifacts, stats, and validation
- **Lessons Learned**: Top 5 insights from the trenches
- **Future Directions**: Where this goes next (DSL Foundry, Rust alternatives, etc.)
- **Q&A**: Your turn

**Speaker notes**: This is a technical deep-dive, not a product launch. If you're interested in DSL design, MPS quirks, or code generation patterns, this is for you. We'll focus on the decisions that mattered and the mistakes that taught us something.

---

## Slide 3: The Problem — Why Boilerplate?

**The Context**

A NATS microservice for User Management needs:
- ✅ Request-reply pattern (publish → wait for response)
- ✅ Entity handlers (User, Roles, Permissions)
- ✅ Hook system (pre/post validations, side effects)
- ✅ OTEL tracing on every operation
- ✅ SDK integration (Motadata Go SDK)
- ✅ SQL schema generation
- ✅ Docker deployment

**The Problem**: Each new entity meant writing all of this by hand.

**Speaker notes**: Imagine adding a fourth entity. You copy-paste the User handler, rename 50+ occurrences, adjust the schema, register hooks. It's mechanical and error-prone. Why aren't we automating this?

---

## Slide 4: The Boilerplate Reality

**What One Entity Requires**

For the `User` entity alone:

| Artifact | Lines | Content |
|----------|-------|---------|
| Handler code | 595 | CRUD operations, hook stubs, validation |
| SQL schema | 8 | Table definition with constraints |
| Hook stubs | 30+ | 4 signatures × multiple operations |
| Integration tests | 21 | End-to-end validation |

**Total per entity**: ~600 lines of hand-written code

**For 3 entities (User, Roles, Permissions)**: **1,790 lines**

**Speaker notes**: Most of this is boilerplate. The actual business logic (validation, hooks) might be 100 lines. The rest is plumbing — subjects, handlers, response formatting. The case for automation is clear.

---

## Slide 5: The Vision

**A DSL That Generates All of This**

```
Entity Definition (in DSL)
     ↓
MPS TextGen Templates
     ↓
Generates (automatically):
  - NATS handler code (1,790 lines)
  - SQL schema (31 lines)
  - Hook stubs (preserved across regeneration)
  - Integration tests (scaffolding)
```

**Why This Matters**: Define once in the DSL, regenerate infinitely. Refactoring a field type updates the schema, handlers, and tests at once.

**Speaker notes**: The goal isn't to write clever code — it's to eliminate repetition. A developer should define "User has username, email, password" once, and the system generates everything else.

---

## Slide 6: The Journey — v1: UMAN (Feb 2-3, 2025)

**First Serious Attempt — The Prisma-Style Approach**

**Path**: `/home/prem-modha/MPSProjects/UMAN`

**The Idea**: Define all entities and relations in a single model file (like Prisma).

```
Models (root concept)
  ├─ User
  ├─ Role
  └─ Permission
```

**What Happened**: MPS's core constraint struck — **one root concept = one generated file**. This created a single `models.go` with all entity handlers jumbled together. It was a **Single Responsibility Principle violation**.

**Lesson**: Schema-centric architecture doesn't work with MPS's 1:1 constraint. You need entity-centric design.

**Status**: Working generation but architecturally flawed → Abandoned

**Speaker notes**: This was a good learning moment early. We discovered a fundamental constraint that would shape everything that came next. The constraint sounds simple, but its implications are deep.

---

## Slide 7: The Journey — v2: UserManagement (Feb 9-Mar 19, 2025)

**The Main Project — Entity-Centric Architecture**

**Paths**:
- `/home/prem-modha/MPSProjects/UserManagement` (original)
- `/home/prem-modha/UserManagmentMps` (reclone, production version)

**The Pivot**: Each entity is now a **separate root concept** in the sandbox.

```
Sandbox
  ├─ Main (root) → main.go
  ├─ User (root) → user.go (+ UserRoles handlers)
  ├─ Roles (root) → roles.go (+ RolesPermissions handlers)
  ├─ Permissions (root) → permissions.go
  └─ SqlSchema (root) → sqlPrem_init_sql.sql
```

**Result**: Separate files! Proper separation of concerns.

**Why Reclone**: MPS's IDE state got corrupted. Faster to reclone than debug MPS's internal state — a common pattern in MPS development.

**Status**: Production-ready ✅

**Speaker notes**: This is the design that stuck. By respecting MPS's constraint instead of fighting it, we got clean separation. Every entity is independently defined and independently generated.

---

## Slide 8: Timeline — Key Milestones

**6-Week Journey**

```
Week 1 (Feb 2-8)
  └─ v1 (UMAN) discovers the single-file problem

Week 2-4 (Feb 9-Mar 8)
  ├─ v2 (UserManagement) pivots to entity-centric
  ├─ Hook system evolves: boolean → named → prioritized + async
  └─ TextGen templates heavily iterated (4+ versions in git)

Week 5-6 (Mar 12-19)
  ├─ Reclone to fix MPS state issues
  ├─ Python TextGen bridge (parse_textgen.py) built
  ├─ SDK + OTEL integration added
  ├─ 21 e2e tests pass
  └─ Production launch ✅
```

**Speaker notes**: The pace accelerated in weeks 5-6 because we had learned the right architecture. The earlier versions weren't wasted — they taught us the constraints.

---

## Slide 9: Key Decision #1 — The 1:1 Constraint

**MPS Core Rule: One Concept = One Generated File**

**What This Means**

When MPS generates a file from a TextGen template, there is a **strict 1:1 mapping** between a root concept instance and an output file.

- 1 `User` concept → 1 `user.go` file
- 1 `Roles` concept → 1 `roles.go` file
- All in a single file? → 1 file with everything (bad design)

**How We Worked Around It**

Instead of fighting the constraint, we **embraced it** as a design rule:

- **Each entity is its own root** → separate, manageable files
- **Relations are children of entities** → generate inside entity files
- **Helper files need external scripts** → `merge_hooks.py` and `sync.sh` extract and transform

**The Insight**: "Don't fight MPS, extend it" — use external tools for what MPS can't do cleanly.

**Speaker notes**: This constraint caused us to pivot twice. But once we stopped seeing it as a limitation and more as a design discipline, everything clicked into place.

---

## Slide 10: Key Decision #2 — The TextGen Walled Garden

**The Problem: Copy-Paste Doesn't Work Inside MPS**

MPS's projectional editor is a **walled garden**:
- ❌ Can't paste code snippets INTO TextGen
- ❌ Can't copy TextGen OUTPUT to an external editor
- ❌ Can't use AI tools to edit templates directly
- ❌ Changes aren't diffable (binary `.mps` format)

Writing a 400-line handler template meant typing every character through MPS completion menus. This is a **real blocker for large templates**.

**The Solution: parse_textgen.py**

**Built**: A Python script that:
1. Parses human-readable markdown (looks like TextGen output)
2. Reverse-engineers MPS's internal TextGen `.mps` XML format
3. Outputs valid `.mps` files loadable into MPS

**Impact**: Write templates in any text editor → parse to `.mps` → import into MPS

**Cost**: One day of work + "a lot of AI tokens" to figure out the `.mps` format

**Speaker notes**: This was a breakthrough moment. Instead of finding a better MPS plugin (DSL Foundry attempt failed), we built a bridge that let us work outside the tool and import results back in. It's an example of pragmatic problem-solving — the tool has a constraint, so automate the bridge.

---

## Slide 11: Key Decision #3 — The Hook System

**Evolution Across Versions**

**v1 (UMAN)**: No hooks at all
- ❌ No extensibility

**v2a (UserManagement)**: Boolean hooks
```go
type User struct {
  PreCreateD bool
  PostCreateN bool
  // ...
}
```
- ✅ Simple
- ❌ Inflexible (only one hook per operation)

**v2b-v3 (Final)**: Named, prioritized, async hooks
```go
// 4 signatures, each serving a distinct purpose
PreSync:  func(...) error                    // Can abort
PreAsync: func(...)                          // Fire-and-forget
PostSync: func(...) ([]byte) []byte          // Transform response
PostAsync: func(..., data []byte)            // Side effects
```

**Why 4 Signatures?**

Emerged naturally from two axes: **pre/post × sync/async**

| Signature | Use Case |
|-----------|----------|
| Pre-sync | Validate before operation (can reject) |
| Pre-async | Send notification without blocking |
| Post-sync | Transform response (strip PII) |
| Post-async | Log event after operation |

**Design Rule**: Always place validating sync hooks at HIGHER priority (run first), so they reject before async hooks fire.

**Speaker notes**: The hook system is the biggest evolution in the DSL. It emerged from real use cases — customers needed extensibility, and that drove the design. This is a pattern: constraints and real requirements shape better designs than grand planning.

---

## Slide 12: What We Built — The Artifacts

**Generated Output (1,790 LOC)**

```
src/
├─ main.go (155 lines)              NATS server setup, endpoint wiring
├─ user.go (595 lines)              User entity + UserRoles relation
├─ roles.go (481 lines)             Roles entity + RolesPermissions
├─ permissions.go (307 lines)       Permissions entity
├─ userDefinedHooks.go (221 lines)  Hook stubs (preserved across regen)
└─ sqlPrem_init_sql.sql (31 lines)  IAM schema (3 tables)
```

**Key Stats**

- **30+ hook stubs** across 4 handler types
- **1,790 lines** of generated code
- **4 hook signatures** (pre-sync, pre-async, post-sync, post-async)
- **21 e2e tests** validating the full pipeline
- **Docker Compose** deployment with NATS + app

**Speaker notes**: The code is production-ready. It includes OTEL tracing on every operation, proper error handling, and full SDK integration. This is not toy code.

---

## Slide 13: What We Built — Generated Code Example

**Entity Definition (in DSL)**

```
Entity User {
  id: uuid (primaryKey)
  username: string
  email: string (nullable)
  password: string (hidden)
  created_at: timestamp
}

Relation UserRoles {
  from: User
  to: Role
  operations: [assign, remove, list]
}
```

**Generated Handler (excerpt from user.go)**

```go
func (h *UserHandler) HandleCreate(req *User) error {
  // Pre-sync hooks (validate)
  if err := h.preCreateD(ctx, span, event); err != nil {
    return err
  }

  // Pre-async hooks (fire-and-forget)
  h.workerpool.Submit(func() {
    h.preCreateM(ctx, span, event)
  })

  // Publish to DAL
  resp, _ := h.publisher.Request(ctx, "motadata.user.db.create", ...)

  // Post-sync hooks (transform response)
  resp = h.postCreateFilter(ctx, span, event, resp)

  // Post-async hooks (side effects)
  h.workerpool.Submit(func() {
    h.postCreateNotify(ctx, span, event, resp)
  })

  return h.Reply(req, resp)
}
```

**Speaker notes**: Every operation follows this pattern: pre-sync validation → pre-async side effects → DAL call → post-sync transform → post-async cleanup. The hooks are preserved across regeneration, so developers can implement business logic in `userDefinedHooks.go` without losing it on the next generation.

---

## Slide 14: Testing — 21 E2E Tests

**demo.sh — Integration Test Suite**

**Architecture**: Runs against Docker Compose stack with mock DAL

```bash
┌─ NATS (port 4229)
├─ UserManagementService (generated)
├─ Mock DAL (5 responders on motadata.*.db.*)
└─ Traffic observer (capture all messages)
```

**Coverage**: 21 tests

| Category | Tests | What's Tested |
|----------|-------|---------------|
| User CRUD | 4 | Create (valid, weak password, missing fields), get, update, delete |
| User ops | 3 | List, validation pipeline, hook execution |
| Roles CRUD | 5 | Create, get, update, delete, list |
| UserRoles | 3 | Assign, remove, list |
| Permissions CRUD | 4 | All CRUD ops |
| Permissions relations | 2 | Assign, remove |
| Service metadata | 1 | NATS `$SRV.INFO` discovery |

**Why No Unit Tests?**

Generating unit tests would require a TextGen template ~10x more complex (all field combinations, all type permutations). Mock DAL is sufficient for validating business logic in isolation.

**Status**: ✅ All 21 passing

**Speaker notes**: The test suite is a quick-and-dirty check that the system works end-to-end. It's not comprehensive, but it's good enough to catch breaking changes. Proper unit tests would be generated later if needed.

---

## Slide 15: Lessons Learned — Top 5 Insights

**#1: Respect MPS Constraints Instead of Fighting Them**

The 1:1 concept-to-file mapping seemed limiting. But when we embraced it as a design rule (each entity is its own root), we got cleaner architecture. The lesson: understand the tool's model and design within it, not against it.

**#2: When Tool UX is Blocking, Build a Bridge**

TextGen's walled garden made iteration painful. Instead of hunting for a better plugin (DSL Foundry attempt failed), we built `parse_textgen.py` to let us work outside MPS and import results. Pragmatism over perfectionism.

**#3: External Scripts Beat New Root Concepts**

The 1:1 constraint meant we couldn't generate `userDefinedHooks.go` directly from a root concept (it would clutter the sandbox). Solution: generate hooks inside entity files, then use `merge_hooks.py` to extract. Pattern: Generate → Script transforms → Final output.

**#4: Hook Signatures Emerge From Real Use Cases**

Started with booleans, evolved to named hooks, then prioritized + async. Each step was driven by a real constraint (extensibility, ordering, non-blocking side effects). Design follows requirements, not the reverse.

**#5: Documentation During Development Saves Months Later**

Captured learnings in a linked knowledge graph immediately. When it came time to compile presentations and guides, we had a structured artifact ready. Future projects will reference these same nodes.

**Speaker notes**: These are patterns, not rules. Each project might need different solutions. But the meta-pattern holds: respect constraints, be pragmatic, iterate on feedback, document as you go.

---

## Slide 16: What Didn't Work

**Failed Experiments (Valuable Learning)**

**DSL Foundry TextGen (M2 attempt)**
- Goal: Escape MPS's text editing pain
- Problem: Lack of documentation, complexity of groupings not well-understood
- Lesson: The tool is good, but you need to understand it deeply first
- Revisit: With current MPS knowledge, this is now tractable

**Over-engineering Entity Definitions**
- Goal: Make the DSL maximally flexible
- Problem: Complexity exploded (10+ field attributes, conditional generation, edge cases)
- Lesson: Start simple, add complexity when requirements demand it
- Final design: Minimal set of attributes, sufficient for 90% of use cases

**Trying to Hide MPS Complexity**
- Goal: Make the IDE "user-friendly" for non-experts
- Problem: MPS's model is too rich to hide; trying to simplify it broke things
- Lesson: Expose the model clearly, train users on it, don't pretend it's simple

**Speaker notes**: These failures aren't wasted. Each taught us what matters and what doesn't. The final design is simple because we stripped away the unnecessary complexity.

---

## Slide 17: Future Directions — Phase 2

**What We'd Do Differently Next Time**

**#1: Use DSL Foundry TextGen**

Despite the M2 failure, DSL Foundry is the right direction. With more MPS experience, the grouping complexity should be tractable. This escapes the char-by-char editing pain without needing `parse_textgen.py`.

**#2: Use MPS Extensions (mbeddr)**

The JetBrains MPS-extensions library provides:
- Grammar-extended editor languages (reduce boilerplate)
- Helper utilities (automate tedious tasks)

This reduces the "MPS overhead" that ate weeks of development.

**#3: Hybrid: MPS → JSON AST → Jinja Templates**

An alternative architecture:
```
MPS Entity Definition
  ↓
JSON AST (machine-readable)
  ↓
External Jinja2 Templates (easy to edit, version control friendly)
  ↓
Final Code
```

Benefits: Best of both worlds (MPS's rich editor + external template flexibility)

**#4: Rust Proc Macros as an Alternative**

We built a working prototype in Rust declarative macros:
```rust
schema! {
  User {
    name: String,
    email: Email,
    password: Password,
  }
}
many_to_many!(UserRole, User, Role);
```

Generates everything the MPS DSL does, but:
- ✅ Editable in any text editor
- ✅ Full IDE support (rust-analyzer, error messages)
- ✅ Distributable as a cargo dependency
- ❌ No visual DSL editor
- ❌ No NATS/messaging integration (yet)

**Speaker notes**: Phase 2 isn't decided yet. The MPS DSL works and is production-ready. But if we were starting fresh, these alternatives would be strong contenders. The key lesson is: there are many ways to build a DSL. Pick the one that matches your constraints (editing UX, distribution, team expertise).

---

## Slide 18: The Meta-Pattern

**Don't Fight Your Tools — Extend Them**

All the successful patterns share a theme:

```
MPS is good at:
  ✅ Structure (entities, fields, concepts)
  ✅ Constraints (validation, type checking)
  ✅ Visual editor (autocomplete, refactoring)

MPS is bad at:
  ❌ Text editing (projectional editor is limiting)
  ❌ File manipulation (1:1 concept-to-file rule)
  ❌ External integration (isolated from CLI tools)

Solution:
  ✓ Use MPS for structure and constraints
  ✓ Use external scripts for file transformation
  ✓ Use external editors (parse_textgen.py) for large templates
  ✓ Use MPS extensions for boilerplate reduction
```

**The Insight**: A well-designed DSL isn't one monolithic tool. It's a **pipeline** where each tool does what it's best at.

**Speaker notes**: This might be the most important slide. Too many teams try to make MPS do everything (because they've invested in it), and they end up with overcomplicated projects. The successful pattern is to respect each tool's strengths and use external tools to fill the gaps.

---

## Slide 19: Q&A

**What's Your Question?**

**Common Topics (if you're thinking about this):**

- **DSL design**: How do you choose between MPS, Rust macros, and other approaches?
- **Code generation**: How do you keep generated code maintainable?
- **Testing**: How much effort to test a DSL?
- **Scaling**: How many entities before this breaks?
- **Distribution**: How do you ship a DSL to users?
- **Experimentation**: How do you prototype before committing?

**Open Discussion**: The knowledge graph is complete and linked. Any specific node you want to dive deeper into?

**Speaker notes**: This is the interactive part. Encourage questions. People often have skepticism about DSLs, and that's worth addressing head-on. "Is this over-engineering?" is a fair question, and the answer is context-dependent.

---

## Slide 20: References & Next Steps

**Artifacts**

- **DSL Repository**: `/home/prem-modha/UserManagmentMps` (production)
- **Knowledge Graph**: `/home/prem-modha/projects/dsl/src/code/user-management-templ/docs/knowledge/` (41 linked nodes)
- **Generated Code**: `/home/prem-modha/projects/dsl/src/code/user-management-templ/src/` (1,790 LOC)
- **Test Suite**: `demo.sh` (21 e2e tests)
- **Parse TextGen Bridge**: `/home/prem-modha/projects/dsl/src/parse_textgen.py`

**Key Files to Read**

1. **For Architecture**: `mps-gotchas/one-concept-one-file-constraint.md`
2. **For Technical Details**: `textgen/entity-textgen-template-full.md`
3. **For Lessons**: `post-success/future-improvements.md`
4. **For Implementation**: `tooling/merge-hooks-story.md` + `tooling/sync-pipeline.md`

**Next Steps**

- [ ] Phase 2: Evaluate DSL Foundry vs. Rust macros
- [ ] Phase 3: Onboarding guide for new DSL users
- [ ] Phase 4: Expand to other microservices (billing, auth, etc.)

**Speaker notes**: The knowledge graph is your reference layer. It's organized, linked, and indexed. Future presentations and onboarding will pull from it. The investment in documentation now saves months of context reconstruction later.

---

## Appendix: Key Metrics

**Development Timeline**
- Total: 6 weeks (Feb 2 - Mar 19, 2025)
- Versions: 3 (v1 UMAN, v2 UserManagement, v3 UserManagmentMps)
- Commits: 30+ (across all projects and reclones)
- Experiments: 21 MPS learning projects

**Code Output**
- Generated lines of code: 1,790
- Hook stubs managed: 30+
- E2E tests: 21 (all passing)
- Tables in schema: 3 (User, Roles, UserRoles, RolesPermissions, Permissions)

**Knowledge Captured**
- Knowledge nodes: 41
- Categories: 12
- Wiki links: 80+
- Screenshots: 5 (MPS editor examples)

**Time Investment**
- TextGen bridge (`parse_textgen.py`): 1 day
- MPS constraint adaptation: 2 weeks (across v1→v2 pivot)
- Testing harness: 1 day
- Documentation: 3 days (concurrent with development)
