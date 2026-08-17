# Deep Interview Spec: DSL Development Learnings Knowledge Graph

## Metadata
- Interview ID: dsl-learnings-001
- Rounds: 6
- Final Ambiguity Score: 19.5%
- Type: brownfield
- Generated: 2026-03-19
- Threshold: 20%
- Status: PASSED

## Clarity Breakdown
| Dimension | Score | Weight | Weighted |
|-----------|-------|--------|----------|
| Goal Clarity | 0.85 | 35% | 0.2975 |
| Constraint Clarity | 0.75 | 25% | 0.1875 |
| Success Criteria | 0.80 | 25% | 0.2000 |
| Context Clarity | 0.80 | 15% | 0.1200 |
| **Total Clarity** | | | **0.805** |
| **Ambiguity** | | | **19.5%** |

## Goal
Build a **persistent knowledge graph** (Obsidian-style linked markdown files) that captures ALL learnings from developing the MPS DSL across multiple projects. This graph serves as a **reusable AI context layer** — instead of `!ls` and `!tree` every time, AI can read the graph to compile targeted outputs (meeting presentations, blog posts, internal docs, onboarding material) without re-exploring the codebase each time.

## Constraints
- **Format**: Linked markdown files with `[[wiki-links]]` creating graph connections
- **Storage**: `docs/knowledge/` directory (or similar) within the project
- **Screenshots**: Prompted at key points — MPS editor views are different from .mps XML files; readers need to see what the DSL author sees
- **Labels**: Clear enough for any person to understand, not just MPS experts
- **Non-linear**: Graph structure supports jumping between topics; not forced chronological
- **Extraction method**: Deep interview to pull tacit knowledge, combined with automated codebase analysis

## Non-Goals
- NOT a polished final document (that's the AI "compilation" step later)
- NOT limited to the 3 main projects — includes experiments and post-success lessons
- NOT a git history reconstruction (git is "desperate saves", not semantic history)

## Acceptance Criteria
- [ ] Topic checklist covers all major areas (see categories below)
- [ ] Each knowledge node is a standalone markdown file with clear title, category, and `[[links]]` to related nodes
- [ ] Cross-project connections are explicit (e.g., "this failed in v1" → "this succeeded in v3")
- [ ] Tacit knowledge is extracted through interview (not derivable from code alone)
- [ ] Screenshot prompts placed at key points where MPS editor view matters
- [ ] Convergence reached: 3 consecutive interview rounds produce no new insights
- [ ] Graph can be read by AI to produce a meeting summary without additional codebase exploration

## Topic Checklist (Coverage Gates)
- [ ] **TextGen pitfalls** — MPS TextGen limitations, the 4+ iteration attempts in v1, the Python bridge solution
- [ ] **Hook system evolution** — boolean hooks → named hooks → prioritized + async hooks, the 4 hook signatures
- [ ] **MPS editor/constraint gotchas** — things about JetBrains MPS itself that surprised you
- [ ] **Architecture pivots** — v1 prototype → UMAN REST → v3 NATS, why each shift happened
- [ ] **SDK integration** — Motadata Go SDK challenges, module path patching, adapter pattern
- [ ] **merge_hooks.py story** — why it exists, key-by-receiver-type design, signature change detection
- [ ] **Deployment/Docker** — containerization learnings, NATS compose setup, port mapping
- [ ] **Testing approach** — demo.sh mock DAL pattern, the 21 e2e tests, what's tested and what's not
- [ ] **DSL concept evolution** — Entity/Relation/Field matured across projects, what was added and why
- [ ] **Post-success lessons** — improvements identified after v3 worked, future directions
- [ ] **Experimental projects** — what was tried in MPSProjects (gtog, AlgebraL2, TaskMaster, etc.) and what was learned
- [ ] **Cross-cutting concerns** — OTEL tracing, multi-tenancy, SQL schema generation

## Assumptions Exposed & Resolved
| Assumption | Challenge | Resolution |
|------------|-----------|------------|
| Git history captures the journey | Contrarian: git shows WHAT, not WHY | User confirmed: git is messy saves, not semantic history. Graph needed. |
| A flat document is sufficient | User wanted all 4 output formats | Resolved: graph as source → AI compiles to specific formats |
| Chronological structure is best | User wanted non-linear exploration | Resolved: graph structure supports cycling back between topics |
| Code artifacts capture everything | Tacit knowledge question | User confirmed: MPS editor experience, dead-end knowledge, cross-project connections are only in their head |

## Technical Context
### Projects to cover
1. **MPSProjects/UserManagement** (v1) — Prototype, 8 commits, heavy TextGen iteration, no working output
2. **MPSProjects/UMAN** (v2) — Working REST API + SQL generation, Schema-centric architecture
3. **UserManagmentMps** (v3) — Production NATS microservice, Entity-centric, named/prioritized/async hooks, SDK+OTEL, Python TextGen bridge
4. **~22 experimental projects** in MPSProjects/ — various MPS learning experiments
5. **Post-success work** — improvements and ideas after v3 achieved the solution
6. **projects/DLS** — additional DSL-related work

### Key codebase paths
- `/home/prem-modha/MPSProjects/` — all MPS experiments
- `/home/prem-modha/UserManagmentMps/` — production DSL
- `/home/prem-modha/projects/DLS/` — DSL work
- `/home/prem-modha/projects/dsl/src/code/user-management-templ/` — generated output + tooling

## Ontology (Key Entities)
| Entity | Type | Fields | Relationships |
|--------|------|--------|---------------|
| Knowledge Node | core domain | title, insight, category, source_project, links[] | connects to other Knowledge Nodes, belongs to Learning Category |
| MPS Project | supporting | name, path, status (prototype/working/production), timeline | contains Knowledge Nodes, evolves DSL Concepts |
| DSL Concept | core domain | name, version_history, current_form | defined in MPS Projects, generates code |
| Output Format | supporting | type (meeting/blog/doc/onboarding), audience | compiled from Knowledge Graph |
| Knowledge Graph | structure | nodes[], edges[], metadata | composed of Knowledge Nodes and Connections |
| Connection | relation | from_node, to_node, relationship_type (caused/evolved/replaced/failed-then-succeeded) | links Knowledge Nodes |
| Learning Category | taxonomy | name (TextGen/hooks/MPS-gotchas/SDK/deploy/test/etc.) | groups Knowledge Nodes |
| AI Context Layer | meta | purpose (reduce !ls and !tree), compilation_targets | implemented by Knowledge Graph |
| Screenshot | evidence | path, caption, what_it_shows, mps_editor_vs_xml | attached to Knowledge Nodes |

## Ontology Convergence
| Round | Entity Count | New | Changed | Stable | Stability Ratio |
|-------|-------------|-----|---------|--------|----------------|
| 1 | 4 | 4 | - | - | N/A |
| 2 | 6 | 2 | 0 | 4 | 67% |
| 3 | 7 | 1 | 0 | 6 | 86% |
| 4 | 8 | 1 | 0 | 7 | 88% |
| 5 | 9 | 1 | 0 | 8 | 89% |
| 6 | 9 | 0 | 0 | 9 | 100% |

## Interview Transcript
<details>
<summary>Full Q&A (6 rounds)</summary>

### Round 1
**Q:** What is the PURPOSE of capturing these learnings? Who will read them and what should they be able to do after reading?
**A:** All four purposes (internal docs, Feynman guide, conference/blog, onboarding) — but clarified: "I want to capture the whole context and then use AI to 'compile' to specific needs like for presenting to meetings vs blogs vs internal documentation etc."
**Ambiguity:** 59% (Goal: 0.6, Constraints: 0.2, Criteria: 0.2, Context: 0.7)

### Round 2
**Q:** What structure should this raw capture take? Think of it as the 'source format' before AI compiles it.
**A:** "I want you to organize it in related brain-like structure similar to the Obsidian.md graph view so it can be compiled into a nice document at once with all the right connections."
**Ambiguity:** 46% (Goal: 0.7, Constraints: 0.5, Criteria: 0.25, Context: 0.7)

### Round 3
**Q:** What categories of learning MUST be in the knowledge graph for it to feel complete?
**A:** All of the above + more (technical pivots, MPS gotchas, DSL evolution, SDK integration, merge_hooks.py, deployment, testing)
**Ambiguity:** 38% (Goal: 0.75, Constraints: 0.5, Criteria: 0.5, Context: 0.75)

### Round 4 (Contrarian Mode)
**Q:** You want to capture EVERYTHING into a graph. But what if that's overkill? Git history, code, and TextGen guides already capture a lot. What does the graph add?
**A:** Three key reasons: (1) Git is "desperate mess of saves" — not semantic history, (2) The graph reduces AI context cost — no more !ls and !tree every time, (3) Cross-project connections and tacit knowledge aren't in any artifact. Also: "can we use document links to make the graphs?"
**Ambiguity:** 30% (Goal: 0.8, Constraints: 0.65, Criteria: 0.55, Context: 0.8)

### Round 5
**Q:** The tacit knowledge is only in YOUR head. What are the top tacit things you learned that no artifact captures?
**A:** "This is where you come in — interview me deep. The doc will need screenshots too. MPS editor view is different from .mps files. There are more lessons than the files and projects mentioned here, including post-success lessons and improvements."
**Ambiguity:** 26% (Goal: 0.85, Constraints: 0.7, Criteria: 0.6, Context: 0.8)

### Round 6
**Q:** How will we know the knowledge capture is 'complete enough'?
**A:** "Have some basic checkbox but also convergence too, and I may like to cycle back as we have it structured as a graph so it can be not chronological or out of the current talking points."
**Ambiguity:** 19.5% (Goal: 0.85, Constraints: 0.75, Criteria: 0.8, Context: 0.8)

</details>
</content>
</invoke>