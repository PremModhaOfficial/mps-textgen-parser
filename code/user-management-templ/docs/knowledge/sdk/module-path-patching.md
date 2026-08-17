---
title: "SDK Module Path Patching"
category: sdk
project: v3
created: 2026-03-20
---

The Motadata Go SDK uses `motadatagosdk` as its Go module name internally, but the generated project needs `dev.azure.com/Motadata/NextGen/motadata-go-sdk`.

## The Problem
- SDK source lives at `~/projects/nextgen/gosdk/motadata-go-sdk/src/motadatagosdk/`
- Its `go.mod` says `module motadatagosdk`
- Generated code imports `dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core`
- These don't match — Go refuses to build

## The Fix (in sync.sh)
```bash
# Patch go.mod
sed -i 's|^module motadatagosdk|module dev.azure.com/Motadata/NextGen/motadata-go-sdk|' go.mod

# Fix ALL Go imports
find . -name '*.go' -exec sed -i 's|"motadatagosdk/|"dev.azure.com/Motadata/NextGen/motadata-go-sdk/|g' {} +
```

## Why Vendored (Not Go Modules)
The SDK is **not hosted on the internet yet** — it's an internal Motadata package without a public or private Go module registry. Vendoring + patching creates an environment where the generated files land and "everything works like magic."

This is a pragmatic solution: copy the SDK source, patch the module path to match what the generated Go imports expect, and include it in the Docker build context.

## Open Questions (lower priority)
- Has the sed approach ever corrupted a file? (Not reported so far)
- When the SDK gets published to a registry, sync.sh can drop the vendoring + patching steps entirely — just use normal `go mod` imports

Links: [[tooling/sync-pipeline]] [[projects/v3-user-managment-mps]]
