# Validation Report — Experiments Category

**Validator**: validator-experiments
**Date**: 2026-03-20
**Nodes checked**: mps-learning-projects.md, pcrd-working-crud-dsl.md
**Note**: experiments-overview.md was listed in the task but does **not exist** in the directory.

---

## Claim 1: MPS has official tutorials for Calculator and Shapes DSLs

**Status**: ✅ CONFIRMED

Both are real, official JetBrains MPS tutorials:

- **Shapes**: An introductory tutorial — [Shapes — An Introductory MPS Tutorial](https://www.jetbrains.com/help/mps/shapes-an-introductory-mps-tutorial.html)
- **Calculator**: An in-depth tutorial covering code generation, type-system, and scoping — [MPS Calculator Language Tutorial](https://www.jetbrains.com/help/mps/mps-calculator-language-tutorial.html)

The `mps-learning-projects.md` node refers to `claculator` (note typo) and `shape`/`newShape` as "tutorial-based projects." This is accurate — these are standard MPS tutorial projects that new users build when learning the tool.

---

## Claim 2: MPS Sandbox Concept — how root nodes work

**Status**: ✅ CONFIRMED

MPS uses a "sandbox" solution within tutorial projects for testing the DSL being built. The official MPS GitHub repo confirms this pattern:

- Path: `samples/calculator-tutorial/solutions/jetbrains.mps.calculator.sandbox/`
- File: `sandbox.mps` — a root node file in the sandbox solution

A **sandbox** in MPS is a separate module/solution where you write programs in your new language to test it. Root nodes are top-level model elements (the unit of persistence in MPS). The sandbox solution contains root nodes written in the tutorial DSL. This is distinct from concepts, which are the schema/structure definitions. The `mps-learning-projects.md` note is consistent with this — projects like `claculator` and `shape` are for building these tutorial DSLs and experimenting in their sandboxes.

Source: [MPS GitHub — calculator sandbox.mps](https://github.com/JetBrains/MPS/blob/master/samples/calculator-tutorial/solutions/jetbrains.mps.calculator.sandbox/jetbrains/mps/calculator/sandbox/sandbox.mps)

---

## Claim 3: Gin is a Go web framework

**Status**: ✅ CONFIRMED

Gin is a well-established, high-performance HTTP web framework for Go:

- Official site: [gin-gonic.com](https://gin-gonic.com/)
- GitHub: [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- Described as "up to 40x faster than Martini" due to httprouter
- Used for building REST APIs, web apps, and microservices

The `pcrd-working-crud-dsl.md` node states that pcrd generated "a Go server for CRUD operations (using Gin framework)." This is plausible and consistent — Gin is the most commonly used Go framework for CRUD REST APIs. The claim is internally consistent with the DSL generating a complete server stack.

---

## Claim 4: Prisma uses a single-file schema definition (schema.prisma)

**Status**: ✅ CONFIRMED (with nuance)

Prisma's primary and default schema approach is a **single file** named `schema.prisma`, written in Prisma Schema Language (PSL). The file contains:
- `datasource` block (DB connection)
- `generator` block (client config)
- `model` blocks (data models and relations)

**Nuance**: Since Prisma ORM v5.15.0, multi-file schema support was added. However, the single-file `schema.prisma` pattern remains the standard default and what virtually all Prisma documentation/tutorials use.

Source: [Prisma Schema Overview](https://www.prisma.io/docs/orm/prisma-schema/overview)

The `pcrd-working-crud-dsl.md` node describes pcrd as hitting MPS's 1:1 concept-to-file constraint as a problem (one giant generated file). This makes sense in contrast to patterns like Prisma where a single schema file is intentional design — in MPS, a single generated file was an accidental anti-pattern caused by putting all entities in one model.

---

## Missing File

**experiments-overview.md**: This file was listed in the task as one of 3 nodes but **does not exist**. No validation possible. Either it was not yet created or was deleted. The directory only contains:
- `mps-learning-projects.md`
- `pcrd-working-crud-dsl.md`

---

## Summary

| Claim | Status | Notes |
|-------|--------|-------|
| MPS official calculator tutorial | ✅ Confirmed | jetbrains.com/help/mps/mps-calculator-language-tutorial.html |
| MPS official shapes tutorial | ✅ Confirmed | jetbrains.com/help/mps/shapes-an-introductory-mps-tutorial.html |
| MPS sandbox concept + root nodes | ✅ Confirmed | Sandbox = test solution for your DSL; root nodes are file-level model elements |
| Gin is a Go web framework | ✅ Confirmed | gin-gonic.com, widely used for REST APIs |
| Prisma single-file schema.prisma | ✅ Confirmed | Default pattern; multi-file added in v5.15.0 but not standard |
| experiments-overview.md exists | ❌ File missing | Not found in directory |

All verifiable claims are **accurate**. The `experiments-overview.md` file is absent and should be created or the task reference updated.
