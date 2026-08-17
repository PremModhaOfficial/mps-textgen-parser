---
title: "The Holder Concept Pattern — 0..1 Wrapper Around 0..n Children"
category: mps-gotchas
project: general
created: 2026-03-20
---

A design pattern used throughout the UserManagement DSL: instead of `0..n` children directly, use a `0..1 HolderConcept` that contains `1..n actualConcept`.

## The Pattern
```
Entity
  └── HookTypeHooksHolder (0..1)   ← holder
       └── Hook (1..n)              ← actual children
```

Instead of:
```
Entity
  └── Hook (0..n)                   ← direct children
```

## Why This Was Chosen
1. **Clean editor for two cases**: When there are 0 hooks vs. 1+ hooks, the editor can show completely different UI depending on whether the holder exists
2. **Different behavior for 0 vs non-0**: The holder's presence/absence acts as a boolean flag — "this entity has hooks" vs "this entity doesn't"
3. **Workaround for limited styling knowledge**: Didn't know the in-depth MPS editor styling options to make `0..n` look good in both empty and populated states

## The Downside
In some cases this made the editor **less friendly** than direct `0..n` would have been — extra nesting, extra clicks to create the holder before adding children.

## Research Result: The Holder Was Unnecessary
Vanilla MPS editor **can** handle cardinality-dependent display without a HolderConcept:

1. **Visibility queries** (`show if`) on cells — e.g., `show if: node.returns.size >= 2` on parenthesis cells
2. **Empty cell property** on collection cells — renders a placeholder when the child list is empty vs populated
3. **Conditional Editor extension** from JetBrains/MPS-extensions — cleaner syntax for conditional rendering

For the Go function return example:
- 0 returns: collection empty, parens hidden (visibility query `size >= 2` = false)
- 1 return: collection shows one item, parens still hidden
- 2+ returns: collection shows items, parens visible

**Zero new concepts needed.** The HolderConcept was a knowledge gap workaround.

## The Lesson
This is a pragmatic workaround born from incomplete knowledge of MPS editor styling. Direct `0..n` with visibility queries is cleaner and avoids the extra nesting that sometimes made the editor less friendly.

Links: [[mps-gotchas/scoping-still-unclear]] [[post-success/future-improvements]] [[mps-gotchas/one-concept-one-file-constraint]]
