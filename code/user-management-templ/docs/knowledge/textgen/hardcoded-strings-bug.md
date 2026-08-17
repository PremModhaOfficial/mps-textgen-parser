---
title: "TextGen Lesson: Use Variables, Not Hardcoded Strings"
category: textgen
project: v3
depth: shallow
created: 2026-03-20
---

The SQL schema generates `user_roless` (double 's') because the plural suffix was hardcoded as a string literal instead of using a variable or computed value.

## The Bug
Junction table name: `iam.user_roless` — the `s` for plural was appended as a hardcoded string, doubling up when the entity name already ended in `s` (Roles → Roless).

## The Lesson
In TextGen templates, prefer **variables and computed expressions** over hardcoded string fragments. When you hardcode `"s"` for pluralization, you lose context about what you're appending to.

## Broader Pattern
This is a category of TextGen bugs: anytime you concatenate hardcoded strings in templates, you risk:
- Double suffixes (this bug)
- Missing separators
- Case mismatches

Links: [[cross-cutting/sql-schema-generation]] [[textgen/textgen-walled-garden]]
