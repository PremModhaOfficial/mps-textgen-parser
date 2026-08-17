---
title: "Architecture Validation Report"
category: architecture
created: 2026-03-20
validator: validator-architecture
---

# Architecture Knowledge Graph — Validation Report

Validated against internet sources on 2026-03-20. Covers 7 nodes:
`json-ast-jinja-idea.md`, `plugin-vs-ide-distribution.md`, `rust-proc-macro-alternative.md`,
`schema-vs-entity-centric.md`, `v1-user-management.md`, `v2-uman.md`, `v3-user-managment-mps.md`

---

## plugin-vs-ide-distribution.md

### Claim 1: Plugin .zip distribution
> "Build a plugin `.zip` that users install into any MPS instance"

**VERIFIED**
MPS official documentation confirms: *"Languages come packaged as ordinary zip files, which you unzip into the MPS plugin directory and which MPS will load upon restart."* Build scripts produce a zip file containing all necessary jar files following MPS plugin layout conventions.

Source: https://www.jetbrains.com/help/mps/building-mps-language-plugins.html

---

### Claim 2: mbeddr as example of standalone IDE distribution
> "Build a complete IDE distribution that bundles MPS + the DSL language (like mbeddr)"

**VERIFIED (with nuance)**
mbeddr DOES offer both distribution modes:
- **Plugin distribution**: a ZIP of plugins installed into an existing MPS instance
- **Standalone IDE**: a bundled installer that includes MPS

The claim that mbeddr is an example of the standalone IDE approach is accurate. However, the node omits that mbeddr also provides a plugin-only ZIP option, so the two approaches are not mutually exclusive for a real product. The standalone IDE claim is correct.

Source: https://mbeddr.com/download.html
Source: https://dslfoundry.com/how-i-make-an-open-source-mbeddr-developer-release/

---

### Claim 3: Plugin approach has "no version tracking burden"
> "Runs on top of any compatible MPS instance — no version tracking burden"

**PARTIALLY INACCURATE**
MPS plugin build scripts `depend on mps`, which means the plugin IS compiled against a specific MPS version. Plugins are not version-agnostic — they must be built against the target MPS version. The burden is *reduced* compared to standalone IDE (no cross-platform build matrix), but version compatibility must still be managed between the plugin and MPS.

The relative claim (plugin is simpler than standalone IDE) is correct. The absolute claim of "no version tracking burden" is an overstatement.

Source: https://www.jetbrains.com/help/mps/building-mps-language-plugins.html

---

## json-ast-jinja-idea.md

### Claim 4: TextGen has a 1:1 concept-to-file constraint
> "dodges the 1:1 concept-to-file constraint completely"

**VERIFIED**
MPS official documentation and support forums confirm: *"TextGen only allows one file per root node, and you can only have one TextGen component per concept."* This is a known, documented architectural limitation of MPS TextGen. Generating multiple files from one concept requires the model-to-model Generator instead.

Source: https://mps-support.jetbrains.com/hc/en-us/community/posts/205830349
Source: https://www.jetbrains.com/help/mps/textgen.html

---

### Claim 5: Jinja templates can consume JSON and produce Go files
> "External Jinja templates consume the JSON and produce Go files"

**VERIFIED**
Jinja2 is a well-established Python templating engine with support for loops, conditionals, filters, and custom extensions. It is widely used for code generation from structured data (including JSON-loaded Python dicts). The proposed MPS → JSON AST → Jinja → Go pipeline is technically sound and commonly used in similar code generation pipelines (e.g., OpenAPI generators, Ansible, Salt).

Source: https://jinja.palletsprojects.com/en/3.1.x/

---

## rust-proc-macro-alternative.md

### Claim 6: Title says "proc_macros" but body says "declarative macros"

**INTERNAL INCONSISTENCY — NEEDS CLARIFICATION**

The node title is: *"Alternative DSL Approach: Rust proc_macros (Working Prototype)"*
The node body says: *"Rust declarative macros (`schema!` and `many_to_many!`)"*

These are different Rust constructs:

| Type | Syntax | Crate type |
|------|--------|-----------|
| Declarative (`macro_rules!`) | Pattern-matching, simpler | Regular crate |
| Procedural (`proc_macro`) | Full Rust code, operates on token streams | Must be `proc-macro = true` crate |

**Both** can generate complex structs and impl blocks. `macro_rules!` macros CAN handle struct + impl generation with repetition syntax (`$($field:ident: $ty:ty),*`), which is plausible for the described use case. For highly complex AST manipulation (like external resource access at compile time), proc_macros are required.

The distinction matters: if `/home/prem-modha/projects/rust/userManMacros` uses `macro_rules!`, the title ("proc_macros") is misleading. If it uses a `#[proc_macro]` function, the body ("declarative macros") is wrong. One or the other needs correction once the actual implementation is inspected.

**Recommendation**: Check `userManMacros/src/lib.rs` — if it contains `macro_rules! schema`, the body is correct and the title should say "declarative macros". If it has `#[proc_macro]`, the title is correct and the body should say "procedural macros".

Source: https://doc.rust-lang.org/reference/procedural-macros.html
Source: https://doc.rust-lang.org/reference/macros-by-example.html

---

### Claim 7: `schema!` generates structs, stores, query builders from ~10 lines
> "This single declaration generates: User struct, UserStore with CRUD, UserQuery builder..."

**PLAUSIBLE / UNVERIFIED AGAINST SOURCE**
Generating this scope of code from a single macro invocation is technically feasible with either declarative or procedural Rust macros. This is similar to what Diesel, SeaORM, and other Rust ORMs do. Cannot be independently verified without reading the actual source at `/home/prem-modha/projects/rust/userManMacros` — but the claim is architecturally credible.

---

## schema-vs-entity-centric.md

### Claim 8: Schema-centric → Entity-centric is a meaningful architectural shift
> "v2 → v3 moved from 'what does the database look like?' to 'what does the business domain look like?'"

**VERIFIED (design principle)**
This maps precisely to Domain-Driven Design (DDD) terminology. "Schema-centric" mirrors an anemic/data-driven model; "Entity-centric" mirrors a domain model. The claim that entity-centric enables business logic hooks and NATS subjects mapped to domain operations (not SQL operations) is consistent with DDD literature and microservice design best practices.

This is a design philosophy claim, not a verifiable technical fact — but it is well-grounded in established patterns.

---

## v3-user-managment-mps.md

### Claim 9: NATS micro.Service API
> "main.go: 155 lines — NATS micro.Service setup, SDK connection, endpoint wiring"

**VERIFIED**
The `github.com/nats-io/nats.go/micro` package provides exactly a `micro.Service` interface with `AddService(nc, config)` for creating microservices, `AddEndpoint()`, `AddGroup()`, `Info()`, `Stats()`, and `Stop()`. The package supports `HandlerFunc`, `Request` (with `Respond()`, `Error()`, `Data()`), and `EndpointConfig` with named subjects — consistent with the described usage.

Source: https://pkg.go.dev/github.com/nats-io/nats.go/micro

---

### Claim 10: NATS subject pattern `motadata.{entity}.{operation}`
> "Subject pattern: `motadata.{entity}.{operation}` → handler → `motadata.{entity}.db.{operation}` (DAL)"

**VERIFIED**
NATS official documentation endorses dot-separated hierarchical subject naming encoding "business or physical entities" and "business intent." The pattern `{domain}.{entity}.{operation}` is a canonical NATS microservice subject design. The two-tier pattern (application → DAL via subject suffix `.db.`) is a valid architectural use of NATS subject hierarchies.

Source: https://docs.nats.io/nats-concepts/subjects

---

## v1-user-management.md / v2-uman.md

### Claim 11: UMAN used Prisma-style single-file schema approach
> "All entities and relations were defined in a single file, similar to how Prisma schema files work"

**VERIFIED (design comparison)**
Prisma does use a single-schema-file approach (`schema.prisma`) for all entity and relation definitions. The comparison is accurate and commonly understood by developers.

Source: https://www.prisma.io/docs/concepts/components/prisma-schema

---

## Summary

| # | Claim | Status | Node |
|---|-------|--------|------|
| 1 | MPS plugin distributed as .zip | ✅ VERIFIED | plugin-vs-ide-distribution |
| 2 | mbeddr as standalone IDE example | ✅ VERIFIED (with nuance) | plugin-vs-ide-distribution |
| 3 | Plugin has "no version tracking burden" | ⚠️ PARTIALLY INACCURATE | plugin-vs-ide-distribution |
| 4 | TextGen 1:1 concept-to-file constraint | ✅ VERIFIED | json-ast-jinja-idea |
| 5 | Jinja2 feasible for JSON→code generation | ✅ VERIFIED | json-ast-jinja-idea |
| 6 | "proc_macros" title vs "declarative macros" body | ❌ INTERNAL INCONSISTENCY | rust-proc-macro-alternative |
| 7 | schema! macro generates structs+stores+queries | ✅ PLAUSIBLE | rust-proc-macro-alternative |
| 8 | Schema→Entity shift is meaningful architectural move | ✅ VERIFIED (DDD) | schema-vs-entity-centric |
| 9 | NATS micro.Service API exists as described | ✅ VERIFIED | v3-user-managment-mps |
| 10 | NATS subject pattern `{domain}.{entity}.{operation}` | ✅ VERIFIED | v3-user-managment-mps |
| 11 | UMAN ≈ Prisma single-file schema style | ✅ VERIFIED | v2-uman |

### Issues Requiring Action
1. **plugin-vs-ide-distribution.md** — soften "no version tracking burden" to "reduced version management complexity"
2. **rust-proc-macro-alternative.md** — title/body inconsistency: inspect `/home/prem-modha/projects/rust/userManMacros/src/lib.rs` to determine if `macro_rules!` (→ fix title) or `#[proc_macro]` (→ fix body)
