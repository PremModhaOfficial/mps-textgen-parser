---
title: "TextGen's Walled Garden — No Copy-Paste, No External Editing"
category: textgen
project: general
created: 2026-03-20
---

The core TextGen pain: MPS's projectional editor is a **walled garden**. The inside and outside of MPS cannot interact with each other.

## The Three Walls
1. **No copy-paste INTO TextGen** — can't paste code snippets from external editors into TextGen buffers. Must type character-by-character and select completions.
2. **No copy-paste FROM TextGen** — can't copy TextGen content to an external `.md` file for editing or AI assistance.
3. **No external tooling** — can't use AI to edit TextGen directly, can't parse scripts to modify templates, can't use any tool outside MPS.

## Why This Hurts
- Writing a 400-line handler template means typing every character through the projectional editor
- Can't use AI coding assistants to help write templates
- Can't iterate on templates in a normal text editor and import
- Can't diff TextGen changes in git (binary `.mps` format)
- Every change requires navigating MPS's completion menus

## The Workaround: parse_textgen.py
Built a Python script (`/home/prem-modha/projects/dsl/src/parse_textgen.py`) that:
1. Takes `.md` files written in a format that LOOKS like TextGen output
2. **Reverse-engineers** the MPS internal TextGen `.mps` XML format
3. Produces `.mps` files that can be loaded into MPS as TextGen templates
4. Used "a lot of AI tokens" to figure out the reverse-engineered format

This breaks through the walled garden: write templates externally → parse to .mps → load into MPS.

## The Insight
When a tool's editing UX is a blocker, sometimes the right answer isn't to find a better plugin (DSL Foundry attempt failed) but to **build a bridge** that lets you work outside the tool and import results back in.

Links: [[textgen/python-bridge-solution]] [[projects/m2-dslfoundry-attempt]] [[mps-gotchas/one-concept-one-file-constraint]]
