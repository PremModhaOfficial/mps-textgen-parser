---
title: "Alternative Approach: MPS TextGen → JSON AST → Jinja Templates"
category: architecture
project: general
created: 2026-03-20
---

An unrealized but promising idea for solving the TextGen pain: instead of generating Go code directly from TextGen, generate a **validated JSON AST** and use **Jinja templates** for the actual code generation.

## The Idea
```
MPS DSL → TextGen generates JSON (single file) → Jinja templates → Go files
          (validated against DSL)                   (external tooling)
```

1. The MPS editor stays exactly as it is (same UX for the human)
2. TextGen generates ONE big JSON object — the DSL's AST as structured data
3. The JSON is guaranteed valid because it's generated from the validated DSL properties
4. External Jinja templates consume the JSON and produce Go files
5. File multiplexing (one file per entity) is handled by external scripts

## Why This Is Elegant
- **TextGen becomes trivial** — just serialize the AST to JSON, no char-by-char Go code authoring
- **Jinja templates are normal text files** — editable in any editor, AI-assistable, diffable in git
- **Escapes the walled garden** entirely — all code generation logic lives outside MPS
- **File splitting is external** — dodges the 1:1 concept-to-file constraint completely
- **Debugging is easy** — inspect the JSON to verify the AST, inspect Jinja output for code issues

## Trade-offs
- Adds a Jinja dependency to the build pipeline
- Two-stage generation means two places where bugs can hide
- JSON schema must exactly mirror the DSL structure (coupling)

## Status
Idea only — not implemented yet. But the parse_textgen.py bridge is a step toward this: it already works outside MPS to manipulate generated output.

Links: [[textgen/textgen-walled-garden]] [[textgen/python-bridge-solution]] [[mps-gotchas/one-concept-one-file-constraint]] [[post-success/future-improvements]]
