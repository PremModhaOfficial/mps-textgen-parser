---
title: "UserManageMentM2 — The DSL Foundry TextGen Attempt (Failed)"
category: textgen
project: experiments
created: 2026-03-20
---

A desperate attempt to escape the pain of MPS native TextGen by switching to the DSL Foundry TextGen plugin.

## The Pain It Tried to Solve
MPS native TextGen requires **char-by-char editing** — laborious manual work to construct output templates. The DSL Foundry TextGen plugin promised a better authoring experience.

## Why It Failed
1. **No good guides** — couldn't find documentation for the DSL Foundry TextGen plugin
2. **Groupings were confusing** — didn't understand how DSL Foundry's groupings and manipulation worked
3. **Hard to reason about** — without understanding groupings, the plugin's model was opaque

## The Insight
The "char-by-char editing" pain of native TextGen was eventually solved differently — not by switching plugins, but by building a **Python bridge** (`parse_textgen.py`) that could work with the native TextGen output format.

## Resolution
Not worth deep-diving — the M2 attempt was unsuccessful and didn't produce lasting insights. The DSL Foundry plugin itself is valid (now in JetBrains/MPS-extensions), but the attempt was abandoned without significant learnings beyond "couldn't find guides, couldn't reason about groupings."

Links: [[projects/v1-user-management]] [[textgen/textgen-walled-garden]] [[textgen/python-bridge-solution]]
