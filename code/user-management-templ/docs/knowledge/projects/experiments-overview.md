---
title: "MPS Experiments — 21 Projects in MPSProjects/"
category: experiments
project: experiments
created: 2026-03-20
---

Beyond the 3 main projects, 21 experimental MPS projects live in `/home/prem-modha/MPSProjects/`. These represent the learning journey with JetBrains MPS itself.

## Project List
| Name | Likely Purpose |
|------|---------------|
| AlgebraL2 | Algebra/math DSL |
| DEC | Unknown — likely early experiment |
| TaskMaster | Task management DSL |
| UserManageMentM2 | User management iteration (between v1 and UMAN?) |
| claculator | Calculator DSL (typo in name) |
| config | Configuration language |
| congifLang | Config language (another typo — pattern of fast iteration) |
| gfunc | Go function generation? |
| gg | Unknown — very short name, likely throwaway |
| gltogl | Unknown |
| golangTgolang | Go-to-Go transpiler? |
| gtog | Go code generator (has README) |
| newShape | Shape DSL (likely from MPS tutorial) |
| pcrd | Unknown |
| shape | Shape DSL (earlier version of newShape?) |
| t, test, testtttt | Test projects (increasing desperation in naming) |
| throwed | Unknown — possibly "thrown away" experiments |

## Patterns Observed
- Typos in project names (`claculator`, `congifLang`) — suggests rapid creation, not polished work
- Multiple attempts at same concept (shape/newShape, config/congifLang, test/testtttt)
- `gtog` is notable — a Go code generator language, possibly the precursor to TextGen templates
- `UserManageMentM2` sits between v1 and UMAN — may hold transitional learnings

## Answers (from interview)
- **Useful experiments**: `gtog` (Go syntax recreation) taught scoping, references, and MPS language design. `pcrd` proved end-to-end generation works. Calculator/shapes taught MPS basics.
- **gtog → TextGen influence**: gtog taught MPS fundamentals (syntax checks, refs, scoping) but didn't directly influence the TextGen template design — TextGen was learned separately through iteration.
- **UserManageMentM2**: Desperate attempt to use DSL Foundry TextGen plugin instead of native TextGen. Failed — couldn't find guides or understand groupings. See [[projects/m2-dslfoundry-attempt]].

Links: [[projects/v1-user-management]] [[projects/v2-uman]]
