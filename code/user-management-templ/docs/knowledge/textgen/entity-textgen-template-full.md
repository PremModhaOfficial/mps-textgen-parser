---
title: "The Full Entity TextGen Template — 300+ Lines of append Statements"
category: textgen
project: v3
depth: deep
screenshot: images/entity-user-definition-1.png
created: 2026-03-20
---

This is the complete TextGen template for the Entity concept. Every line of the generated `user.go`, `roles.go`, and `permissions.go` is produced by this template. Each `append {...}` statement was typed character-by-character in MPS's projectional editor.

## Template Structure

```
text gen component for concept Entity {
  file path: "/src/"
  extension: "go"
  file name: node.name.toLowerCase()

  body:
    1. Package declaration + imports
    2. Entity struct (foreach field → Go struct field with json/db tags)
    3. Event types per operation (CreatedEvent, UpdatedEvent, etc.)
    4. Handler struct + constructor
    5. Handler methods per operation:
       a. OTEL span + context injection
       b. JSON unmarshal into typed event
       c. Pre-hooks (sorted by priority, sync vs async)
       d. Field validation (per operation type)
       e. Span attributes
       f. DAL request (subject + headers + timeout)
       g. Post-hooks (sync transforms response, async fire-and-forget)
       h. Respond
    6. #HOOKS_START / #HOOKS_END — hook stubs for merge_hooks.py
    7. Relation handlers (delegated to Relation TextGen via forEach)
}
```

## Key Patterns in the Template

### Conditional imports
```
append ${node.needWorkerPool() ? "\"...workerpool\"" : ""} \n ;
```
Only imports workerpool if any hook is async.

### Dynamic validation
```
foreach field in node.fields {
  if (!(field.hasAnotation("primaryKey")) && !(field.hasAnotation("auto"))
      && !(field.hasAnotation("hidden")) && !(field.hasAnotation("nullable"))) {
    append { || event.Name.Field == ""} ;
  }
}
```
Builds a validation expression by iterating non-auto, non-hidden, non-nullable fields.

### Hook sorting and dispatch
```
foreach hook in prehook.getHooksSorted() {
  if (hook.isAsync) {
    // workerpool.AsyncWithCtx pattern
  } else {
    // if err := s.hookName(); err != nil { return } pattern
  }
}
```
Hooks are sorted by priority (via `getHooksSorted()` behavior), then dispatched based on `isAsync`.

### The #HOOKS markers
```
append {// #HOOKS_START} \n ;
// ... generate all hook stubs ...
append {// #HOOKS_END} \n ;
```
These markers are what `merge_hooks.py` uses to extract stubs.

## Why This Template Illustrates the Pain
- ~300 lines of `append` statements, each typed char-by-char in MPS
- Every `{` and `}` and `\n` is a deliberate editor action
- String interpolation via `${...}` requires navigating MPS completion menus
- No copy-paste from a working Go file — every pattern must be re-typed
- Debugging: if line 47 of `user.go` is wrong, you must mentally trace which `append` produced it

## The Full Template
<details>
<summary>Complete Entity TextGen (click to expand)</summary>

See the full template in the user's interview response — too large to inline here but preserved in the interview transcript at `.omc/specs/deep-interview-dsl-learnings.md`.

The template is approximately 300 lines of TextGen code producing ~400 lines of Go per entity.

</details>

Links: [[textgen/textgen-walled-garden]] [[textgen/textgen-debugging-pain]] [[hooks/four-hook-signatures]] [[tooling/merge-hooks-story]] [[mps-gotchas/one-concept-one-file-constraint]]
