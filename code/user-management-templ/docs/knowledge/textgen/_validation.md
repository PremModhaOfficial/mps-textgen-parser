---
title: "TextGen Category — Validation Report"
category: textgen
generated: 2026-03-20
validator: validator-textgen
---

# TextGen Validation Report

Validation of factual claims across 6 TextGen knowledge nodes:
- `textgen-walled-garden.md`
- `textgen-debugging-pain.md`
- `python-bridge-solution.md`
- `hardcoded-strings-bug.md`
- `entity-textgen-template-full.md`
- `projects/m2-dslfoundry-attempt.md`

---

## Claim 1: MPS TextGen has no copy-paste (must type char-by-char)

**Status: ✅ CONFIRMED**

The MPS projectional editor does not support arbitrary text copy-paste. You can copy/paste AST nodes (whole structural nodes), but not arbitrary text fragments. This is a fundamental constraint of the projectional editor design — the editor manipulates the AST directly and requires whole-node selection for copy/paste. Pasting free-form text into TextGen buffers is not possible through normal clipboard operations.

**Sources:**
- [MPS Community: copy/paste text in MPS editor](https://mps-support.jetbrains.com/hc/en-us/community/posts/206610375-copy-paste-text-in-MPS-editor) — confirms no textual copy/paste in MPS projectional editor
- [MPS Documentation: Cut, copy, paste](https://www.jetbrains.com/help/mps/cutting-copying-and-pasting.html) — documents node-based copy-paste only
- [DSL Foundry plaintextgen](https://github.com/DSLFoundry/mps-plaintextgen) — explicitly designed to address this gap (paste text from buffer into MPS)

**Nuance:** The DSL Foundry plaintextgen plugin was *built specifically to solve this problem*, which confirms the limitation in native MPS/TextGen. The claim in the nodes is accurate.

---

## Claim 2: No breakpoints or line tracing in TextGen

**Status: ✅ CONFIRMED (with nuance)**

There is no in-place TextGen step-debugger. The MPS debugger exists for BaseLanguage behaviors, but to use breakpoints during generation you must create a separate run configuration and launch a second MPS instance — making it impractical for normal TextGen iteration. "Printf debugging" is acknowledged by JetBrains themselves as the common fallback.

The node's claim that "no plugin or extension (including mbeddr/MPS-extensions) adds a TextGen step-debugger" is consistent with what search results show — no such plugin was found.

**Sources:**
- [MPS Support: Debug behavior and generator with breakpoints](https://mps-support.jetbrains.com/hc/en-us/community/posts/206612255-Debug-behavior-and-generator-with-breakpoints) — confirms you must run a second MPS instance for breakpoints; printf debugging is a known pain point
- [MPS Documentation: Using the MPS Debugger](https://www.jetbrains.com/help/mps/using-mps-debugger.html)
- [MPS Documentation: TextGen](https://www.jetbrains.com/help/mps/textgen.html)

---

## Claim 3: System.out.println inside `${ }` blocks works, output in MPS Messages console

**Status: ⚠️ PLAUSIBLE BUT NOT DIRECTLY CONFIRMED**

Web sources confirm printf-style debugging is the de-facto community approach for MPS TextGen. However, no official source explicitly states that `System.out.println` inside TextGen `${ }` expression blocks specifically routes output to the MPS Messages console. This is a widely-used community pattern, but the specific mechanism (Messages console vs. stdout/run window) was not confirmed by official docs.

**Recommendation:** This claim should be treated as community knowledge. It is highly plausible (TextGen `${ }` blocks execute as Java/BaseLanguage expressions on the JVM inside MPS), but tag the node as "community-observed, not officially documented."

**Sources:**
- [MPS Support: Debug behavior and generator with breakpoints](https://mps-support.jetbrains.com/hc/en-us/community/posts/206612255-Debug-behavior-and-generator-with-breakpoints) — printf debugging referenced but not TextGen-specific

---

## Claim 4: "Preview Generated Text" right-click action exists, available in MPS 2021.x+

**Status: ⚠️ PARTIALLY CONFIRMED — availability pre-dates 2021, exact version unknown**

The "Preview Generated Text" right-click action on a root node is a real feature confirmed in multiple MPS tutorials and the plaintextgen tutorial. However, MPS 2021.1 release notes do NOT mention this as a new feature, suggesting it existed before 2021.x. The claim that it is "available in MPS 2021.x+" is likely imprecise — it was probably available in earlier MPS versions.

**Recommendation:** Update the node to say "available in MPS (version uncertain, at least through 2021.x)" rather than implying it was introduced in 2021.

**Sources:**
- [Plaintextgen tutorial — DSLFoundry](https://dslfoundry.com/plaintextgen-tutorial/) — uses "Preview Generated Text" right-click action
- [MPS 2021.1 Release Notes](https://blog.jetbrains.com/mps/2021/05/mps-2021-1-has-been-released/) — does not mention this feature as new, suggesting it predates 2021.1
- [Generator User Guide Demo5](https://www.jetbrains.com/help/mps/generator-user-guide-demo5.html) — references preview functionality

---

## Claim 5: DSL Foundry TextGen plugin exists and offers a better authoring experience

**Status: ✅ CONFIRMED**

The DSL Foundry `mps-plaintextgen` plugin is real, exists on GitHub, and does offer copy-paste of text blocks with macro parameterization. It was designed precisely to solve the walled-garden problem of native TextGen. The plugin is now archived/deprecated and has been moved into the official [JetBrains MPS-extensions](https://jetbrains.github.io/MPS-extensions/extensions/generator/plaintext-gen/) project.

**Additional finding:** The claim that documentation was hard to find (m2-dslfoundry-attempt.md) was plausible at the time — the main docs are a GitHub README and a tutorial on dslfoundry.com, with no official JetBrains docs page. The project has since become a maintained JetBrains extension.

**Sources:**
- [GitHub: DSLFoundry/mps-plaintextgen](https://github.com/DSLFoundry/mps-plaintextgen)
- [DSLFoundry plaintextgen tutorial](https://dslfoundry.com/plaintextgen-tutorial/)
- [JetBrains MPS-extensions: Plaintext Generator](https://jetbrains.github.io/MPS-extensions/extensions/generator/plaintext-gen/)

---

## Claim 6: .mps files are XML format

**Status: ✅ CONFIRMED**

MPS saves models in an XML-based format. The `.mps` files are XML documents representing the AST. Multiple official sources confirm this.

**Sources:**
- [MPS Documentation: Custom Persistence Cookbook](https://www.jetbrains.com/help/mps/custom-persistence-cookbook.html)
- [DSLFoundry: How to write a dump of an MPS model to XML](https://dslfoundry.com/how-to-write-a-dump-of-an-mps-model-to-xml-based-on-its-structure/)
- Web search results consistently describe `.mps` as "XML-based format" representing the AST

---

## Claim 7: No good guides existed for DSL Foundry TextGen plugin (m2-dslfoundry-attempt.md)

**Status: ⚠️ PARTIALLY CONFIRMED — sparse docs, but not zero**

At the time of the M2 attempt, DSL Foundry documentation consisted of a GitHub README and a single tutorial blog post. It was not well-documented. The claim of "no good guides" is defensible — there was minimal documentation. The plugin has since been better integrated into MPS-extensions with more coverage.

**Sources:**
- [DSLFoundry GitHub](https://github.com/DSLFoundry)
- [Plaintextgen tutorial](https://dslfoundry.com/plaintextgen-tutorial/)

---

## Summary Table

| Claim | Node | Status | Confidence |
|-------|------|--------|------------|
| No copy-paste into TextGen (char-by-char) | textgen-walled-garden | ✅ Confirmed | High |
| No copy-paste from TextGen | textgen-walled-garden | ✅ Confirmed | High |
| No breakpoints / step-debugger in TextGen | textgen-debugging-pain | ✅ Confirmed | High |
| System.out.println in `${ }` → Messages console | textgen-debugging-pain | ⚠️ Plausible, unconfirmed | Medium |
| "Preview Generated Text" right-click action exists | textgen-debugging-pain | ✅ Confirmed (feature real) | High |
| "Preview Generated Text" added in MPS 2021.x+ | textgen-debugging-pain | ⚠️ Likely imprecise — predates 2021 | Low |
| DSL Foundry TextGen plugin exists | m2-dslfoundry-attempt | ✅ Confirmed | High |
| DSL Foundry hard to find docs for | m2-dslfoundry-attempt | ⚠️ Partially confirmed | Medium |
| .mps files are XML format | textgen-walled-garden, python-bridge-solution | ✅ Confirmed | High |
| Hardcoded string bug (double 's') — pattern claim | hardcoded-strings-bug | ✅ Plausible, domain-specific | N/A |

---

## Recommended Node Updates

1. **textgen-debugging-pain.md**: Change "available in MPS 2021.x+" to "available in MPS (version uncertain, at least 2021.x — likely earlier)" and mark `System.out.println` claim as community-observed.
2. **m2-dslfoundry-attempt.md**: Note that DSL Foundry plaintextgen is now part of official JetBrains MPS-extensions (not just a third-party plugin).

---

*Validated: 2026-03-20 | Sources: JetBrains MPS docs, JetBrains MPS Blog, DSLFoundry GitHub, MPS Community forums*
