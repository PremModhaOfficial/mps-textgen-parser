---
title: "parse_textgen.py — The Reverse-Engineered TextGen Bridge"
category: textgen
project: v3
created: 2026-03-20
---

The breakthrough solution: a Python script that converts human-readable markdown into MPS TextGen `.mps` format.

**Path**: `/home/prem-modha/projects/dsl/src/parse_textgen.py`

## How It Works
1. Author writes a `.md` file that mirrors how TextGen templates look (readable text with placeholders)
2. `parse_textgen.py` parses this markdown format
3. Converts it to MPS's internal XML format (`.mps` files)
4. The `.mps` file is loaded into MPS as a TextGen template
5. MPS treats it as if it was authored natively

## The Reverse Engineering
The `.mps` format for TextGen was **reverse-engineered** from existing TextGen templates:
- Examined existing TextGen `.mps` XML structure
- Mapped each TextGen node type to its XML representation
- Used AI assistance extensively to decode the format ("a lot of tokens")
- Built the parser to produce valid `.mps` files

## Why This Matters
- Breaks the MPS walled garden — can now edit templates in any text editor
- Enables AI-assisted template authoring (write markdown → parse → import)
- Makes TextGen changes diffable in git (markdown is text, not binary)
- Dramatically speeds up template iteration

## Git Evidence
Recent commits show active development of this bridge:
- `parse_textgen.py` modified in latest commits
- PR #1 and #2: "fix grouping append parts" — fixing parser edge cases
- "format expressions and control flows as strings" — improving parser capabilities

## Answers (from interview)
- **Time**: One full day — including bug fixing and handling all permutations and combinations of TextGen node types
- **Hardest parts**: The permutations and combinations of how TextGen nodes map to XML structure — each node type has different child arrangements
- **Current limitations**: Not specified — the parser handles the Entity and Relation TextGen templates fully. Edge cases may exist for future concept types.

Links: [[textgen/textgen-walled-garden]] [[projects/v1-user-management]] [[mps-gotchas/one-concept-one-file-constraint]]
