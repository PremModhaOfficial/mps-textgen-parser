# Knowledge Extraction — Master Checklist

> Resume guide: Read this file + `.omc/specs/deep-interview-dsl-learnings.md` to restore context.
> Storage: Knowledge nodes go in `docs/knowledge/` as linked markdown files.

## Pipeline Status
- [x] Deep interview — spec crystallized (ambiguity 19.5%)
- [x] Ralplan consensus — Planner done, Architect ITERATE (applied), Critic ACCEPT (fixes applied)
- [x] Knowledge graph scaffold — 14 dirs, template, README, build_index.py
- [x] Automated extraction — skeleton nodes from codebase
- [x] Interview rounds — 8 rounds of tacit knowledge capture
- [x] Screenshots — 9 MPS editor screenshots captured and organized
- [ ] Cross-linking pass — `[[wiki-links]]` verified between nodes
- [ ] Validation — AI can read graph and produce summary without `!ls`/`!tree`

## Final Stats (2026-03-20)
- **29 nodes**, **7,036 words**, **87 edges**, **12 categories**, **9 screenshots**, **0 broken links**

## Topic Tracker (12 categories)

Mark: `[ ]` pending, `[~]` interviewing, `[x]` captured, `[!]` needs screenshot

| # | Topic | Status | Nodes | Notes |
|---|-------|--------|-------|-------|
| 1 | TextGen pitfalls | [x] | 5 | Walled garden, python bridge, M2 DSL Foundry, debugging pain, full Entity template, hardcoded strings bug |
| 2 | Hook system evolution | [x] | 3 | Boolean→named→prioritized+async, 4 signatures, customer extensibility |
| 3 | MPS editor/constraint gotchas | [x] | 3 | 1:1 concept-to-file constraint, scoping unclear, holder concept pattern |
| 4 | Architecture pivots | [x] | 4 | Schema vs entity, JSON/Jinja idea, plugin vs IDE, project timeline corrected |
| 5 | SDK integration | [x] | 1 | Vendored because not hosted online, module path patching |
| 6 | merge_hooks.py story | [x] | 1 | Born from 1:1 constraint, 3-stage evolution, sig change detection |
| 7 | Deployment/Docker | [x] | 1 | Port mapping, compose setup |
| 8 | Testing approach | [x] | 1 | demo.sh, mock DAL, 21 e2e tests |
| 9 | DSL concept evolution | [x] | 1 | Field types from SQL primitives to domain types |
| 10 | Post-success lessons | [x] | 1 | DSL Foundry, external scripts, MPS extensions |
| 11 | Experimental projects | [x] | 3 | pcrd (working CRUD DSL!), learning projects, DLS/toys first workspace |
| 12 | Cross-cutting concerns | [x] | 1 | SQL gen, hardcoded string bug |

## Convergence Tracker
| Round | New insights? | Topics touched |
|-------|--------------|----------------|
| 1 | Yes — UMAN came first, reclone story, M2 DSL Foundry | Architecture, TextGen |
| 2 | Yes — experiments by name, gtog = MPS learning | Experiments |
| 3 | Yes — TextGen walled garden, no copy-paste, parse_textgen.py | TextGen |
| 4 | Yes — hooks as customer extensibility, not just dev convenience | Hooks |
| 5 | Yes — SDK vendored (not online), merge_hooks from 1:1 constraint | SDK, Tooling |
| 6 | Yes — port 4229 trivial, user_roless bug, OTEL no issues | Deploy, Cross-cutting |
| 7 | Yes — DLS/toys = first workspace, scoping gaps | Experiments, MPS gotchas |
| 8 | Yes — Rust proc_macro attempt, JSON/Jinja idea, holder pattern, plugin vs IDE | Architecture, MPS gotchas |

## Screenshots Captured
| File | What it shows |
|------|--------------|
| entity-user-definition-1.png | User entity: fields, types, annotations, hooks |
| entity-user-definition-2.png | User entity continued: more hooks, relations |
| roles-entity-with-relation.png | Roles entity with Permission_pivot relation |
| permissions-entity-hooks.png | Permissions entity with pre/post hooks |
| sandbox-overview.png | Logical View: all 5 root concepts |
| behavior-tree.png | All behavior definitions |
| structure-tree-all-concepts.png | All 19 structure concepts |
| entity-concept-definition.png | Entity concept: properties, children, references |
| main-concept-definition.png | Main concept: config properties |

## Key Decisions
- Graph format: Obsidian-style `[[wiki-links]]` in markdown
- Purpose: AI context layer → compile to meetings/blogs/docs/onboarding
- Done gate: Topic checklist + convergence + user says done
- Git history is NOT reliable (messy saves) — graph captures the WHY
- Knowledge graph is UPSTREAM, AsciiDoc docs are DOWNSTREAM compilations

## Pending Research (background agents)
- MPS TextGen debugging methods
- MPS editor cardinality-dependent styling (can 0..n replace holder pattern?)
- Rust proc_macro DSL attempt at ~/projects/rust/userManMacros/
