---
title: "Validation Report — MPS Gotchas Category"
validator: validator-mps-gotchas
date: 2026-03-20
nodes_checked: 3
---

## Summary

| Node | Overall Verdict |
|------|----------------|
| `one-concept-one-file-constraint.md` | ✅ CONFIRMED (with nuance) |
| `holder-concept-pattern.md` | ⚠️ PARTIALLY INACCURATE — two sub-claims are wrong |
| `scoping-still-unclear.md` | ⚠️ MINOR INACCURACY — wrong method name |

---

## Node 1: `one-concept-one-file-constraint.md`

### Claim: 1:1 mapping between root concept instance and TextGen output file

**Verdict: ✅ CONFIRMED**

Official MPS TextGen documentation confirms: each rootable concept with TextGen defined produces one output file. The "target file can also be specified" per TextGen component. You cannot get multiple files out of one root concept instance via TextGen alone — the support forum questions asking "TextGen — multiple files output per concept instance?" confirm this is a known limitation, and the standard answer is to use the Generator (model-to-model) aspect instead of TextGen.

**Source:** https://www.jetbrains.com/help/mps/textgen.html
**Community confirmation:** https://mps-support.jetbrains.com/hc/en-us/community/posts/206613545

### Nuance worth noting

The documentation states "MPS does not create files for root concept automatically" — meaning TextGen must be explicitly defined per concept. Sub-concepts do not automatically get a file. The claim in the node is accurate in practice: one root concept instance with TextGen = one output file.

### Architectural conclusion (that each entity must be a separate root)

**Verdict: ✅ CONFIRMED** — This follows directly from the 1:1 mapping. To get `user.go`, `roles.go`, `permissions.go` as separate files, each must be a separate root concept instance. Confirmed valid reasoning.

---

## Node 2: `holder-concept-pattern.md`

### Claim A: MPS `show if` visibility queries exist on cells

**Verdict: ✅ CONFIRMED**

Official MPS editor documentation states: "Elements that should only be visible under certain condition should have its `show if` property set." This is a core editor feature accessible via the Inspector (Alt+2). The example `show if: node.returns.size >= 2` is a plausible predicate expression.

**Source:** https://www.jetbrains.com/help/mps/editor.html
**Source:** http://mbeddr.com/mps-platform-docs/aspects/editor/ — also confirms: "The meta-model allows specifying the `show if` property... the generator will generate the correct code for it."

**Caveat:** The mbeddr platform docs note the inspector doesn't visually expose `show if` for custom cells — you may need the reflective editor to set it. But the feature exists.

---

### Claim B: `[%ifEmpty%]` in collection cells renders differently when list is empty vs populated

**Verdict: ❌ NOT CONFIRMED — likely inaccurate label**

No official MPS documentation, MPS Extensions documentation, or community posts use the syntax `[%ifEmpty%]` as an editor feature. This notation resembles TextGen template syntax (`$list{...}`, `${...}`) not editor cell syntax. The actual MPS editor mechanism for collection empty-state rendering is:

- **"empty cell"** property on collection cells in the editor — controls what is displayed when a collection has no children
- **`text*` / `empty-text*`** flags on cells — custom placeholder text when empty
- **`show if`** on wrapper cells — hide surrounding punctuation (brackets, keywords) when collection is empty

The label `[%ifEmpty%]` does not appear in any official source. It may be a misremembering or informal shorthand for the "empty cell" property. The underlying behavior described (render differently when empty vs populated) is real — but the specific syntax `[%ifEmpty%]` is not a recognized MPS API or editor DSL construct.

**Sources:**
https://www.jetbrains.com/help/mps/editor-cookbook.html
https://specificlanguages.com/posts/basic-editors/use-empty-text-for-empty-cells/

---

### Claim C: mbeddr `ConditionalCellLayout` as a cleaner syntax for conditional display

**Verdict: ❌ INACCURATE — wrong feature name**

`ConditionalCellLayout` does not exist in the MPS Extensions or mbeddr platform. The `de.itemis.mps.editor.celllayout` language (`Cell Layout`) handles **spatial layout** (margins, borders, grid, grow/push) — not conditional display logic.

The correct MPS Extensions feature for conditional editor display is the **Conditional Editor** extension (listed separately in https://jetbrains.github.io/MPS-extensions/). This provides condition-based overriding of editor definitions.

Additionally, the claim associates `ConditionalCellLayout` with `de.itemis.mps.editor.celllayout`. The actual namespace for Grammar Cells (the textual notation helper) is `com.mbeddr.mpsutil.grammarcells`, not `de.itemis`.

**Corrected statement:** Use the **Conditional Editor** extension from MPS-extensions (https://jetbrains.github.io/MPS-extensions/) for condition-based rendering, not a `ConditionalCellLayout` which does not exist.

**Sources:**
https://jetbrains.github.io/MPS-extensions/
https://jetbrains.github.io/MPS-extensions/extensions/editor/celllayout/
https://jetbrains.github.io/MPS-extensions/extensions/editor/grammar-cells/

---

### Overall assessment of holder-concept-pattern.md

The **core insight is valid**: a HolderConcept was an unnecessary workaround; direct `0..n` with editor visibility features could have handled empty vs populated states cleanly. However, two of the three supporting mechanisms are described inaccurately:

- `show if` — ✅ real feature, correctly named
- `[%ifEmpty%]` — ❌ not a recognized MPS editor syntax
- `ConditionalCellLayout` — ❌ does not exist; correct feature is "Conditional Editor" extension

---

## Node 3: `scoping-still-unclear.md`

### Claim A: `ScopeProvider` is the MPS mechanism for controlling variable visibility

**Verdict: ✅ CONFIRMED**

Official MPS Scopes documentation confirms: "MPS starts looking for the closest ancestor to the reference node that implements `ScopeProvider` and who can provide scope for the current kind." Nodes implement the `ScopeProvider` concept interface to control what is visible for references.

**Source:** https://www.jetbrains.com/help/mps/scopes.html

---

### Claim B: "Override `ScopeProvider.scope()` to return a list of valid variables"

**Verdict: ⚠️ MINOR INACCURACY — wrong method name**

The method is **`getScope()`**, not `scope()`. The official documentation states the method signature is `getScope(kind, child)` where:
- `kind` — the concept of the possible target for the reference
- `child` — the child node from which the scope request came

The file uses `ScopeProvider.scope()` which is not the correct method name. This is a minor naming error that would confuse someone trying to implement it.

**Corrected statement:** Override `ScopeProvider.getScope(concept<> kind, SContainmentLink link, int index)` (or the two-parameter variant `getScope(kind, child)`) to return valid variables.

**Source:** https://www.jetbrains.com/help/mps/scopes.html

---

### Claim C: The UserManagement DSL sidesteps scoping by using entity references by name

**Verdict: ✅ PLAUSIBLE / OUT OF SCOPE FOR EXTERNAL VALIDATION**

This is an architectural decision about this specific project. Not externally verifiable, but consistent with the observed DSL design (EntityRef nodes rather than variable binding). No contradiction found.

---

---

## Claim 6 (from task, not explicit in any single file): MPS projectional editor doesn't support traditional copy-paste for TextGen

**Verdict: ✅ VERIFIED (with clarification on scope)**

**Source:** [copy/paste text in MPS editor — JetBrains Support](https://mps-support.jetbrains.com/hc/en-us/community/posts/206610375-copy-paste-text-in-MPS-editor)

**Evidence:**
> "There is no textual copy/paste support in the MPS editors."

Copy-paste in MPS operates at the **AST node level**, not the text level. You cannot select an arbitrary character range and paste it as raw text — the editor always operates on the structured tree.

**Clarification:** The limitation is not specific to TextGen; it applies to the MPS projectional editor in general. TextGen only runs at generation time, not during editing. The practical consequence is: you cannot paste snippets of Go code (or any other text) directly into an MPS editor — you must construct nodes through the editor's menus and completions. This is a well-known usability friction point of projectional editors.

---

## Corrections Needed

| Node | Item | Correction |
|------|------|------------|
| `holder-concept-pattern.md` | `[%ifEmpty%]` claim | Replace with: MPS collection cells have an **"empty cell"** property that controls display when empty; no `[%ifEmpty%]` editor syntax exists |
| `holder-concept-pattern.md` | `mbeddr ConditionalCellLayout` | Replace with: **Conditional Editor** extension from MPS-extensions; `de.itemis.mps.editor.celllayout` is Cell Layout (spatial layout only), not conditional display |
| `scoping-still-unclear.md` | `ScopeProvider.scope()` | Replace with: `ScopeProvider.getScope()` — correct method name per official docs |

---

## Final Summary

| # | Claim | Status |
|---|-------|--------|
| 1 | 1:1 root concept → TextGen output file | **VERIFIED** |
| 2 | `show if` visibility predicate on editor cells | **VERIFIED** |
| 3 | `ScopeProvider` API controls variable scoping | **VERIFIED** (method is `getScope()`, not `scope()`) |
| 4 | mbeddr→JetBrains/MPS-extensions has grammar cells + `de.itemis.mps.editor.celllayout` | **VERIFIED** (package name correct for Cell Layout; `ConditionalCellLayout` name **CORRECTED** → "Conditional Editor") |
| 5 | `[%ifEmpty%]` in MPS collection cells | **UNVERIFIED** — no documentation found; real mechanism is `empty cell` property |
| 6 | MPS projectional editor lacks traditional text copy-paste | **VERIFIED** (general editor limitation, not TextGen-specific) |

---

## Sources

- [TextGen | MPS Documentation](https://www.jetbrains.com/help/mps/textgen.html)
- [Editor | MPS Documentation](https://www.jetbrains.com/help/mps/editor.html)
- [Editor Cookbook | MPS Documentation](https://www.jetbrains.com/help/mps/editor-cookbook.html)
- [Scopes | MPS Documentation](https://www.jetbrains.com/help/mps/scopes.html)
- [MPS Extensions](https://jetbrains.github.io/MPS-extensions/)
- [Cell Layout — MPS Extensions](https://jetbrains.github.io/MPS-extensions/extensions/editor/celllayout/)
- [Grammar Cells — MPS Extensions](https://jetbrains.github.io/MPS-extensions/extensions/editor/grammar-cells/)
- [Editor Aspect — MPS Platform Docs (mbeddr)](http://mbeddr.com/mps-platform-docs/aspects/editor/)
- [Use empty text for empty cells — specificlanguages.com](https://specificlanguages.com/posts/basic-editors/use-empty-text-for-empty-cells/)
- [copy/paste text in MPS editor — JetBrains Support](https://mps-support.jetbrains.com/hc/en-us/community/posts/206610375-copy-paste-text-in-MPS-editor)
- [TextGen: multiple files output per concept instance? — JetBrains Support](https://mps-support.jetbrains.com/hc/en-us/community/posts/206613545)
