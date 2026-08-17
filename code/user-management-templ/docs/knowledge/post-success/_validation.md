---
title: "Validation: post-success/future-improvements.md"
validator: validator-post-success
date: 2026-03-20
task: "9"
---

# Validation Report — post-success/future-improvements.md

## Claims Validated

### 1. DSL Foundry TextGen Plugin

**Claim**: "DSL Foundry TextGen is the right direction for future work."

**Verdict**: ✅ VALID (with important status note)

**Findings**:
- The plugin EXISTS and was originally at `github.com/DSLFoundry/mps-plaintextgen`
- It is listed on the JetBrains Plugin Marketplace (ID: 8444)
- The original DSLFoundry repo is **ARCHIVED** (last commit April 2019)
- The plugin has been **migrated** to `github.com/JetBrains/MPS-extensions` where it is actively maintained as "Plaintext Generator"
- Documentation: https://jetbrains.github.io/MPS-extensions/extensions/generator/plaintext-gen/

**Action for author**: Update references to point to the new home in JetBrains/MPS-extensions, not the original DSLFoundry repo.

---

### 2. mbeddr/MPS-extensions GitHub Project

**Claim**: "The mbeddr/MPS-extensions project provides tools that automate tedious MPS tasks"

**Verdict**: ✅ VALID (with repo URL correction)

**Findings**:
- The project EXISTS and is actively maintained
- Correct GitHub URL: `https://github.com/JetBrains/MPS-extensions` (NOT `mbeddr/MPS-extensions`)
- The mbeddr team originally developed these extensions; they were migrated to JetBrains
- The mbeddr GitHub org (`github.com/mbeddr`) still exists but `mbeddr/MPS-extensions` is not the right repo
- Active: 953 releases, latest November 2025; 5,454 commits; maintained by itemis, JetBrains, and open source community
- Documentation: https://jetbrains.github.io/MPS-extensions/

**Action for author**: Correct the attribution from "mbeddr/MPS-extensions" to "JetBrains/MPS-extensions". The mbeddr origin is historical context, not the current home.

---

### 3. Grammar-Extended Editor Languages in mbeddr/MPS-extensions

**Claim**: "Grammar-extended editor languages — easier editor definitions than vanilla MPS"

**Verdict**: ✅ VALID

**Findings**:
- This feature exists as **Grammar Cells** (`com.mbeddr.mpsutil.grammarcells`) in MPS-extensions
- Documented at: https://jetbrains.github.io/MPS-extensions/extensions/editor/grammar-cells/
- Description: "an alternative approach to the MPS transformation language for declaratively specifying textual notations and their interactions for the MPS editor"
- Provides: `grammar.flag`, `grammar.constant`, `grammar.optional`, `grammar.wrap`, `grammar.substitute`, `grammar.brackets`, `grammar.rules`
- Reduces boilerplate in editor definitions, consistent with the claim

---

### 4. MPS-extensions Provides Helper Utilities

**Claim**: "Helper utilities — things that take significant time normally but become trivial with these extensions"

**Verdict**: ✅ VALID

**Findings**:
- MPS-extensions provides a broad set of utilities beyond Grammar Cells
- Official description: "The MPS extensions aim to ease language development within MPS"
- Maintained by itemis, JetBrains, and open source community
- Covers editor, generator (plaintextgen), and other MPS development areas

---

## Summary

| Claim | Verdict | Notes |
|-------|---------|-------|
| DSL Foundry TextGen exists | ✅ Valid | Migrated to JetBrains/MPS-extensions; original repo archived |
| mbeddr/MPS-extensions exists on GitHub | ⚠️ Partially correct | Repo is JetBrains/MPS-extensions, not mbeddr/MPS-extensions |
| Grammar-extended editor languages exist | ✅ Valid | Called "Grammar Cells" in MPS-extensions |
| MPS-extensions provides helper utilities | ✅ Valid | Actively maintained, broad coverage |

## Corrections Needed in Source Document

1. **DSL Foundry TextGen**: Note that the canonical current location is JetBrains/MPS-extensions, not the original DSLFoundry GitHub org.
2. **Repo attribution**: Change "mbeddr/MPS-extensions" to "JetBrains/MPS-extensions". Mbeddr is the historical origin, JetBrains is the current home.

## Sources

- https://github.com/DSLFoundry/mps-plaintextgen (archived, original DSL Foundry repo)
- https://github.com/JetBrains/MPS-extensions (current, active home)
- https://plugins.jetbrains.com/plugin/8444-com-dslfoundry-plaintextgen (JetBrains marketplace)
- https://jetbrains.github.io/MPS-extensions/extensions/editor/grammar-cells/ (Grammar Cells docs)
- https://jetbrains.github.io/MPS-extensions/extensions/generator/plaintext-gen/ (Plaintext Generator docs)
- https://mbeddr.com/platform.html (mbeddr platform overview)
