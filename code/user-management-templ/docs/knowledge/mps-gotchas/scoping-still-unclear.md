---
title: "MPS Scoping — Still Not Fully Understood"
category: mps-gotchas
project: general
created: 2026-03-20
---

Variable scoping in MPS remains a knowledge gap even after building the full DSL.

## What Scoping Means in MPS
When your DSL has variables (like Go's `var x = 5`), MPS needs to know:
- Which variables are **valid** in a given context (visible in scope)
- Which variables are **invalid** (not in scope, should show errors)
- How to **export** variables from one scope to another

## What Was Learned (from gtog + calculator)
The Go-to-Go project (`gtog`) was the most complex learning exercise for scoping. Calculator also practiced this. Key concepts:
- `ScopeProvider` — the MPS mechanism for controlling variable visibility
- Override `ScopeProvider.getScope()` to return a list of valid variables
- Must explicitly manage what's in scope and what's not

## What's Still Unclear
- How `ScopeProvider` works in detail
- How to make a list of variables in scope and export them
- How to override `ScopeProvider.getScope()` correctly
- The interaction between scoping and the MPS constraint system

## Why This Matters for the DSL
The current UserManagement DSL sidesteps complex scoping — entities reference each other by name, not by variable binding. But future DSL features (like computed fields or conditional logic in templates) would need proper scoping.

## The Insight
You can build a production DSL without fully mastering every MPS feature. Scoping was learned enough for the learning projects but the UserManagement DSL was designed to avoid needing deep scoping knowledge.

Links: [[experiments/mps-learning-projects]] [[mps-gotchas/one-concept-one-file-constraint]]
