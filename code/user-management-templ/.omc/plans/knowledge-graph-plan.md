# Implementation Plan: DSL Learnings Knowledge Graph

> Source spec: `.omc/specs/deep-interview-dsl-learnings.md`
> Checklist: `.omc/knowledge-extraction-checklist.md`

## RALPLAN-DR Summary

### Principles
1. **Graph-first, document-second** — nodes and links are the source of truth; documents are compiled outputs
2. **Tacit before automated** — interview the user first (they forget details over time), then fill gaps from code
3. **One insight per node** — small, linkable units beat monolithic pages
4. **Screenshot at the boundary** — prompt for MPS editor screenshots where .mps XML diverges from visual representation
5. **AI-consumable index** — a machine-readable map so AI doesn't need `!ls`/`!tree`

### Decision Drivers
1. **Resumability** — context compaction can happen mid-interview; checklist + index must allow seamless resume
2. **Compilation flexibility** — the same graph must compile into meetings, blogs, docs, onboarding without restructuring
3. **Cross-project traceability** — connections between v1/UMAN/v3 must be first-class, not buried in prose

### Viable Options

**Option A: Flat linked markdown (Obsidian-style)** — CHOSEN
- Each node is a `.md` file in `docs/knowledge/`
- Frontmatter for metadata (category, project, tags)
- `[[wiki-links]]` for connections
- `_index.json` for AI consumption
- Pros: Simple, familiar, works with Obsidian directly, easy to read/write
- Cons: No enforced schema, links can break silently
- Trade-off: Simplicity over rigidity — appropriate for a knowledge capture that will evolve

**Option B: Structured JSON graph + rendered markdown**
- `graph.json` as canonical source with nodes/edges
- Markdown files rendered from JSON
- Pros: Strict schema, queryable, no broken links
- Cons: Editing JSON is painful, can't use Obsidian natively, adds build step
- Rejected: Over-engineered for a capture that needs to grow organically during interviews

## 1. Directory Structure

```
docs/knowledge/
├── _index.json              # Machine-readable graph index for AI
├── _template.md             # Node template (copy for new nodes)
├── README.md                # How to use this knowledge graph
│
├── projects/                # Project context nodes
│   ├── v1-user-management.md
│   ├── v2-uman.md
│   ├── v3-user-managment-mps.md
│   ├── experiments-overview.md
│   └── dls-project.md
│
├── textgen/                 # Topic: TextGen pitfalls
│   ├── textgen-iteration-hell.md
│   ├── python-bridge-solution.md
│   ├── append-parts-gotcha.md
│   └── ...
│
├── hooks/                   # Topic: Hook system evolution
│   ├── boolean-hooks-v1.md
│   ├── named-hooks-v2.md
│   ├── priority-async-hooks-v3.md
│   ├── four-hook-signatures.md
│   └── ...
│
├── mps-gotchas/             # Topic: MPS editor/constraint surprises
│   └── ...
│
├── architecture/            # Topic: Architecture pivots
│   ├── v1-to-uman-pivot.md
│   ├── rest-to-nats-pivot.md
│   ├── schema-vs-entity-centric.md
│   └── ...
│
├── sdk/                     # Topic: SDK integration
│   ├── module-path-patching.md
│   ├── adapter-pattern.md
│   └── ...
│
├── tooling/                 # Topics: merge_hooks.py, sync.sh
│   ├── merge-hooks-story.md
│   ├── sync-pipeline.md
│   └── ...
│
├── deployment/              # Topic: Docker/NATS
│   └── ...
│
├── testing/                 # Topic: demo.sh, e2e approach
│   └── ...
│
├── dsl-evolution/           # Topic: DSL concept maturation
│   └── ...
│
├── post-success/            # Topic: Post-success lessons
│   └── ...
│
├── experiments/             # Topic: ~22 experimental projects
│   └── ...
│
├── cross-cutting/           # Topic: OTEL, multi-tenancy, SQL
│   └── ...
│
└── images/                  # Screenshots (MPS editor views)
    └── ...
```

## 2. Node Template

Two modes: **lightweight** (during interview) and **structured** (post-interview pass).

### Lightweight mode (capture fast, structure later)
```markdown
---
title: "Short descriptive title"
category: textgen | hooks | mps-gotchas | architecture | sdk | tooling | deployment | testing | dsl-evolution | post-success | experiments | cross-cutting
project: v1 | v2-uman | v3 | experiments | dls | general
created: YYYY-MM-DD
---

{Free-text capture — whatever the user says, get it down}

Links: [[related-node]]
```

### Structured mode (post-interview structure pass)
```markdown
---
title: "Short descriptive title"
category: textgen | hooks | mps-gotchas | architecture | sdk | tooling | deployment | testing | dsl-evolution | post-success | experiments | cross-cutting
project: v1 | v2-uman | v3 | experiments | dls | general
tags: [tag1, tag2]
depth: shallow | deep
screenshot: images/slug.png
created: YYYY-MM-DD
---

# {Title}

## The Situation
What was happening / what was the goal at that point.

## What Happened
The specific experience, decision, or surprise.

## The Insight
The generalizable lesson — what would you tell someone facing this.

## Links
- [[related-node-1]] — how it connects (type: caused|evolved|replaced|failed-then-succeeded|depends-on|contradicts|related)
- See also: [[project-node]] for broader context
```

> **Architect feedback applied**: `screenshot_needed: bool` replaced with `screenshot: path` (doubles as tracker and reference). Added `depth: shallow|deep` for compiler prioritization. Lightweight mode preserves interview momentum.

## 3. AI-Consumable Index (`_index.json`) — GENERATED ONLY

`_index.json` is **never hand-maintained**. A generation script (`docs/knowledge/build_index.py`) walks all `.md` files, parses frontmatter + `[[wiki-links]]`, and produces the index. Run it on-demand after interview sessions.

```python
# docs/knowledge/build_index.py (Phase 1 deliverable)
# Walks docs/knowledge/**/*.md, extracts frontmatter + [[links]]
# Produces _index.json with nodes[], edges[], categories[], projects{}
```

Output format:
```json
{
  "version": 1,
  "generated": "2026-03-19",
  "description": "DSL development learnings knowledge graph",
  "node_count": 0,
  "categories": ["textgen", "hooks", "mps-gotchas", ...],
  "projects": {
    "v1": { "path": "/home/prem-modha/MPSProjects/UserManagement", "status": "failed-prototype" },
    "v2-uman": { "path": "/home/prem-modha/MPSProjects/UMAN", "status": "working-abandoned" },
    "v3": { "path": "/home/prem-modha/UserManagmentMps", "status": "production" }
  },
  "nodes": [
    { "file": "textgen/textgen-iteration-hell.md", "title": "...", "category": "textgen", "project": "v1", "depth": "deep", "links": ["textgen/python-bridge-solution.md"] }
  ],
  "edges": [
    { "from": "textgen/textgen-iteration-hell.md", "to": "textgen/python-bridge-solution.md", "type": "caused" }
  ]
}
```

> **Architect feedback applied**: Single source of truth is frontmatter. `_index.json` is a derived artifact, never edited directly. Eliminates staleness risk during active interviewing.

## 3a. Relationship to Existing `docs/` AsciiDoc

The project already has `docs/sections/*.adoc` (7 Feynman-style guide sections) and `docs/lessons.adoc`. There is direct overlap with the knowledge graph — particularly `07-lessons-and-red-flags.adoc` covers hooks, MPS gotchas, and tooling.

**Decision**: The knowledge graph is the **upstream source**. The AsciiDoc guide is a **compiled downstream output**.

- `docs/knowledge/` = raw capture (the graph, edited during interviews)
- `docs/sections/*.adoc` = polished compilation (written by AI from the graph)
- `docs/lessons.adoc` = compiled from knowledge nodes tagged with war-story patterns

The existing AsciiDoc skeleton in `docs/sections/` stays as the output target. The Feynman docs skill (`/feynman-docs write <section>`) will read knowledge graph nodes as source material instead of (or in addition to) raw codebase exploration.

> **Architect feedback applied**: Explicitly defines knowledge graph as upstream, AsciiDoc as downstream. No two sources of truth.

## 4. Phase 1: Automated Extraction (~30 min)

Pull facts from the codebase to create **skeleton nodes** that the interview will enrich:

| What to extract | Source | Becomes node in |
|----------------|--------|-----------------|
| Project structures (file counts, language models) | `ls`, `find` across 3 projects | `projects/` |
| Concept list (Entity, Relation, Field, Hook, etc.) | `.structure.mps` files | `dsl-evolution/` |
| TextGen guide content | `*_guide.md` files in v1 and v3 | `textgen/` |
| Generated code stats (LOC per entity, handler count) | `src/*.go` | `architecture/` |
| Hook signature catalogue | `userDefinedHooks.go` | `hooks/` |
| SQL schema details | `sqlPrem_init_sql.sql` | `cross-cutting/` |
| Experimental project names + sizes | `MPSProjects/*/` | `experiments/` |
| sync.sh + merge_hooks.py logic | Script source code | `tooling/` |
| Docker/compose config | `docker-compose.yml`, `Dockerfile` | `deployment/` |
| demo.sh test catalogue | `demo.sh` | `testing/` |

| Experimental project scan (names, sizes, languages) | `MPSProjects/*/` directory listing | `experiments/` + interview prompts |

> **Architect feedback applied**: Experiment scan moved to Phase 1 (was topic 10 in interviews). Surfaces project names early so interviewer can ask "did you hit this in experiment X too?" during topics 2-5.

**Estimated output**: ~15-20 skeleton nodes with codebase facts, no tacit knowledge yet.

## 5. Phase 2: Interview Workflow

Systematic deep interview across 12 topics. One topic at a time, but user can jump between topics (graph supports non-linear).

### Interview order (recommended, not enforced)
1. **Architecture pivots** — sets the narrative frame (WHY did things change?)
2. **TextGen pitfalls** — the hardest part, richest lessons
3. **MPS editor/constraint gotchas** — screenshot-heavy, do early while fresh
4. **Hook system evolution** — key design story
5. **DSL concept evolution** — builds on architecture + hooks
6. **merge_hooks.py story** — concrete tooling decision
7. **SDK integration** — technical but well-documented in code
8. **Deployment/Docker** — smaller topic
9. **Testing approach** — smaller topic
10. **Experimental projects** — broad survey, may reveal forgotten lessons
11. **Cross-cutting concerns** — OTEL, multi-tenancy (fills gaps)
12. **Post-success lessons** — reflective, best done last

### Per-topic interview protocol
1. Show user the skeleton node (automated facts)
2. Ask: "What's missing? What surprised you? What would you warn someone about?"
3. Probe for the WHY behind decisions
4. Ask for `[[connections]]` to other topics
5. Prompt for screenshots where MPS editor view matters
6. Update node, run `build_index.py`, mark topic in checklist

### Convergence detection
After each round, check: did this round produce new nodes or substantive edits?
Track in `.omc/knowledge-extraction-checklist.md` convergence tracker.
3 consecutive rounds with no new insights → declare topic complete.

## 6. Phase 3: Screenshot Integration

### Key screenshot points (prompt user)
| When | What to capture | Why |
|------|----------------|-----|
| MPS editor showing Entity definition | How User/Roles/Permissions look in the DSL | .mps XML is unreadable; editor view is the real UX |
| MPS editor showing Relation with hooks | The hook priority/async UI | Shows the declarative power |
| TextGen template editor | What writing a TextGen template looks like | Core development experience |
| MPS constraint/typesystem editor | How validation rules are defined | Often surprising to non-MPS users |
| Sandbox model instance | A concrete DSL instance before generation | Shows the "user experience" of the DSL |
| Generated code diff (v1 vs v3) | How output evolved | Visual impact of DSL improvements |

### Protocol
- When a node has a non-empty `screenshot:` field but the referenced file doesn't exist, prompt: "Can you take a screenshot of [X] in MPS and save it to `docs/knowledge/images/{slug}.png`?"
- Reference in node as `![{caption}](images/{slug}.png)`

## 7. Phase 4: Cross-Linking Pass

After all topics are covered:
1. Read all nodes and identify implicit connections (same project, same concept, cause-effect)
2. Add `[[wiki-links]]` where missing
3. Verify no orphan nodes (every node has at least one link)
4. Run `build_index.py` to regenerate edges
5. Generate a text-based graph summary for AI context

### Connection types
| Type | Meaning | Example |
|------|---------|---------|
| `caused` | A led to B | TextGen failure → Python bridge |
| `evolved` | A became B | Bool hooks → named hooks |
| `replaced` | A was abandoned for B | REST (UMAN) → NATS (v3) |
| `failed-then-succeeded` | Failed in project X, worked in Y | v1 TextGen → v3 TextGen |
| `depends-on` | B requires understanding A | Hook signatures → merge_hooks.py |
| `contradicts` | A was proven wrong by B | v1 assumption about TextGen → v3 reality |
| `related` | Topically connected | OTEL tracing ↔ SDK integration |

## 8. Phase 5: Validation

| Criterion | How to verify |
|-----------|--------------|
| Topic checklist complete | All 12 categories marked `[x]` in checklist |
| Each node is standalone | Node reads coherently without other nodes |
| Cross-project connections exist | At least 5 `failed-then-succeeded` or `evolved` edges |
| Tacit knowledge captured | Nodes contain insights not derivable from code |
| Screenshots present | All nodes with `screenshot:` field have corresponding image files at the referenced path |
| AI-consumable | Feed `_index.json` + 3 random nodes to AI, ask for meeting summary — does it work? |
| Convergence reached | 3 rounds no new insights in convergence tracker |

## Risk Assessment

| Phase | Risk | Mitigation |
|-------|------|------------|
| Automated extraction | MPS XML files are hard to parse | Use guide.md files + file structure, not XML parsing |
| Interview | Context compaction mid-interview | Checklist + `_index.json` allow resume |
| Interview | User forgets tacit details | Start with architecture pivots (narrative triggers memory) |
| Screenshots | User may not have MPS open | Batch screenshot requests, don't block on them |
| Cross-linking | Over-linking makes graph noisy | Limit to 3-5 links per node, use typed connections |
| Validation | "AI-consumable" is subjective | Concrete test: give AI the index + 3 nodes, ask for summary |

## Estimated Node Count

| Category | Estimated nodes |
|----------|----------------|
| Projects | 5 |
| TextGen | 4-6 |
| Hooks | 4-5 |
| MPS gotchas | 3-5 |
| Architecture | 3-4 |
| SDK | 2-3 |
| Tooling | 2-3 |
| Deployment | 1-2 |
| Testing | 1-2 |
| DSL evolution | 3-5 |
| Post-success | 2-4 |
| Experiments | 3-5 |
| Cross-cutting | 2-3 |
| **Total** | **35-52 nodes** |

## ADR: Knowledge Graph Format

**Decision**: Use flat linked markdown files with frontmatter and `[[wiki-links]]`

**Drivers**: Need for organic growth during interviews, Obsidian compatibility, AI readability

**Alternatives considered**:
- JSON graph database → rejected: painful to edit during interviews
- Single monolithic document → rejected: can't compile to different formats
- AsciiDoc (match existing docs/) → rejected: no wiki-link convention, harder for graph

**Why chosen**: Lowest friction for interview-driven content creation while maintaining machine-readable structure via `_index.json`

**Consequences**: Links can break if files are renamed; mitigated by `_index.json` regeneration script

**Follow-ups**: After knowledge capture is complete, build a compilation script that reads `_index.json` and produces targeted outputs (meeting deck, blog post, etc.)
