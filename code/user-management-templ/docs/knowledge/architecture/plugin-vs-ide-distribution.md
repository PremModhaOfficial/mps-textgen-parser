---
title: "Distributing the DSL: Plugin Approach vs Standalone IDE"
category: architecture
project: general
depth: deep
created: 2026-03-20
---

To make the DSL usable by others in the organization, it needs to be packaged as a software artifact. Two approaches exist in MPS.

## The Problem
The DSL works on the developer's machine inside MPS. But to distribute it to others:
- It needs to be a **software entity** that others can install
- The MPS build system is complex and poorly understood — "full of files I don't know about and what they do"
- Packaging DSL Foundry as a dependency was a struggle (for the generation plan to work on other machines)

## Option A: Standalone MPS IDE (like mbeddr)
Build a complete IDE distribution that bundles MPS + the DSL language.

**Downsides:**
- Must track MPS version AND language version compatibility
- Must handle cross-platform build system conflicts (OS-specific MPS libs)
- Much more complex build configuration
- Larger distribution artifact

## Option B: MPS Plugin (Chosen)
Build a plugin `.zip` that users install into any MPS instance.

**Advantages:**
- Only manage the language and its dependencies
- Simple build → plugin solution
- Users install the `.zip` into their existing MPS, import the language in the sandbox, and use it
- Runs on top of any compatible MPS instance — reduced version tracking burden (still need MPS version compatibility, but avoid cross-platform OS-specific build issues)

## The Struggle
Even the plugin approach has friction:
- The MPS build system has many files whose purpose isn't clear
- Packaging DSL Foundry as a plugin dependency (for generation) was a specific pain point
- Need to research: what's the current best practice for packaging MPS plugins with third-party dependencies?

## The Insight
Plugin > IDE for distribution complexity. The trade-off: plugins depend on users having MPS installed, but avoid the entire cross-platform build matrix problem.

Links: [[post-success/future-improvements]] [[projects/m2-dslfoundry-attempt]] [[mps-gotchas/one-concept-one-file-constraint]]
