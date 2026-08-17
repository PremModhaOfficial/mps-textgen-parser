---
title: "SDK Category Validation Report"
category: sdk
type: validation
created: 2026-03-20
validated_by: validator-sdk
---

# Validation Report: `sdk/module-path-patching.md`

## Summary

All core facts in the document are **VALID**. One platform caveat and one nuance about the standard Go alternative are noted below.

---

## Fact-by-Fact Verification

### 1. Go module path must match import paths
**Status: CONFIRMED**

Go requires that the module path declared in `go.mod` matches the import paths used in source files. If `go.mod` declares `module motadatagosdk` but the generated code imports `dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core`, the build fails. This is a fundamental Go modules invariant.

**Source:** [Go Modules Reference — go.dev/ref/mod](https://go.dev/ref/mod)

---

### 2. `replace` directive is the standard way to use local modules
**Status: CONFIRMED — but with an important nuance**

The `replace` directive in `go.mod` is the documented standard mechanism for redirecting a module path to a local directory. However, the `replace` directive only works at the consumer level. It does **not** fix import paths *inside* the SDK itself. If the SDK's own Go files (`events/core`, `events/transport/nats`, etc.) import each other using `motadatagosdk/...`, a `replace` in the consumer's `go.mod` would not fix those internal cross-package imports — they would still resolve to the old name.

The `sed` approach in `sync.sh` patches both the `go.mod` declaration AND all internal SDK import strings, which is why it works where a plain `replace` directive would not. This is a recognized (if non-standard) workaround for vendoring unhosted SDK code with mismatched internal imports.

**Sources:**
- [Go Modules Reference — go.dev/ref/mod](https://go.dev/ref/mod)
- [go.mod file reference — go.dev/doc/modules/gomod-ref](https://go.dev/doc/modules/gomod-ref)

---

### 3. Vendoring is a recognized Go pattern for non-published modules
**Status: CONFIRMED**

`go mod vendor` is an official Go command that copies dependencies into a `vendor/` directory. Since Go 1.14, if a `vendor/` directory exists and `go.mod` specifies `go 1.14+`, the toolchain defaults to `-mod=vendor` automatically. Using vendoring for internal/unhosted packages (not yet on a module registry) is a well-documented pattern.

**Sources:**
- [Go Modules Reference — Vendoring](https://go.dev/ref/mod#vendoring)
- [Go Wiki: Modules](https://go.dev/wiki/Modules)

---

### 4. `sed -i` behavior for in-place editing of Go files
**Status: CONFIRMED for Linux — CAVEAT for macOS**

On **Linux** (GNU sed), `sed -i 's|...|...|'` performs in-place editing without creating a backup file. This is correct and safe for patching `.go` files and `go.mod`.

**macOS caveat (not applicable here):** BSD sed (macOS) requires `sed -i ''` (empty string suffix). Running `sed -i 's|...|...|'` on macOS would fail with a "invalid command code" error. Since this project targets Linux (confirmed by environment info), the current syntax is correct.

**Sources:**
- [sed in-place editing Linux vs macOS — thoughtbot](https://thoughtbot.com/blog/sed-102-replace-in-place)
- [Portable sed -i — johndcook.com](https://www.johndcook.com/blog/2023/10/18/portable-sed-i/)

---

### 5. `rsync --delete` removes extra files from destination
**Status: CONFIRMED**

The `--delete` flag tells rsync to delete files in the destination that do not exist in the source, making the destination an exact mirror. By default (without `--delete`), rsync only adds/updates files and never removes anything from the destination.

**Sources:**
- [DigitalOcean: How To Use Rsync](https://www.digitalocean.com/community/tutorials/how-to-use-rsync-to-sync-local-and-remote-directories)
- [explainshell.com: rsync --delete](https://www.explainshell.com/explain?cmd=rsync+--delete)

---

## Overall Assessment

| Claim | Status |
|-------|--------|
| Go module path must match import paths | ✅ Confirmed |
| `replace` directive is standard for local modules | ✅ Confirmed (with nuance: sed approach is necessary here because it also fixes *internal* SDK imports) |
| Vendoring is recognized for non-published modules | ✅ Confirmed |
| `sed -i` works for in-place Go file editing | ✅ Confirmed (Linux only; macOS needs `sed -i ''`) |
| `rsync --delete` removes extra destination files | ✅ Confirmed |

**No false facts detected.** The sed-based patching approach is pragmatic and correct given the SDK's internal import structure. The standard `replace` directive alone would be insufficient here.
