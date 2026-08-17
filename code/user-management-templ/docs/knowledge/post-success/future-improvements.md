---
title: "Post-Success: What to Do Differently Next Time"
category: post-success
project: general
created: 2026-03-20
---

After the UserManagmentMps DSL succeeded, these are the identified improvements for future projects.

## 1. Use DSL Foundry TextGen for New Projects
Despite the failed M2 attempt, DSL Foundry TextGen is the right direction for future work. The M2 failure was due to lack of guides and not understanding groupings — NOT because the plugin was bad. Now with more MPS experience, the groupings should be tractable.

**Why**: Escapes the char-by-char editing pain of native TextGen without needing the parse_textgen.py reverse-engineering bridge.

## 2. Use External Helper Scripts Instead of New Root Concepts
Key lesson from merge_hooks.py and sync.sh: when MPS's 1:1 concept-to-file constraint gets in the way, don't fight it by creating cluttering root concepts. Instead, generate a combined output and use external scripts to split/transform it.

**Pattern**: Generate → External script transforms → Final output
**Already proven**: merge_hooks.py (extract hooks), sync.sh (patch SDK paths)

## 3. Use MPS Extensions (mbeddr)
The JetBrains/MPS-extensions project provides tools that automate tedious MPS tasks:
- **Grammar-extended editor languages** — easier editor definitions than vanilla MPS
- **Helper utilities** — things that take significant time normally but become trivial with these extensions

**Why**: Reduces the "MPS boilerplate" — the overhead of defining editors, constraints, and typesystem rules for each concept.

## The Meta-Insight
All three improvements share a theme: **don't fight MPS, extend it**. The successful pattern is:
- Use MPS for what it's good at (structure, constraints, concept modeling)
- Use external tools for what MPS is bad at (text editing, file manipulation, build pipelines)
- Use community extensions to reduce MPS-specific boilerplate

Links: [[textgen/textgen-walled-garden]] [[projects/m2-dslfoundry-attempt]] [[tooling/merge-hooks-story]] [[mps-gotchas/one-concept-one-file-constraint]]
