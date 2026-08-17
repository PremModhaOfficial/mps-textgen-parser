---
title: "sync.sh — The Generation Pipeline"
category: tooling
project: v3
created: 2026-03-20
---

`sync.sh` is the bridge between MPS generation and a running Docker service.

## Pipeline Steps
1. Copy `.go` and `.sql` files from MPS `source_gen/` into `src/`
2. Run `merge_hooks.py` to extract and merge hook stubs
3. Extract project name from `main.go` first-line comment (`// MotaDSLUM`)
4. Sync vendored SDK from `~/projects/nextgen/gosdk/` via rsync
5. Patch SDK module path: `motadatagosdk` → `dev.azure.com/Motadata/NextGen/motadata-go-sdk`
6. Fix SDK internal imports to match patched module path
7. Kill any processes on ports 4229/8222 (`fuser -k 4229/tcp 8222/tcp`)
8. Docker compose build + up

## The SDK Patching Problem
The Motadata Go SDK source uses `motadatagosdk` as its module name, but the project's `go.mod` expects `dev.azure.com/Motadata/NextGen/motadata-go-sdk`. sync.sh does:
```bash
sed -i 's|^module motadatagosdk|module dev.azure.com/...|' go.mod
find . -name '*.go' -exec sed -i 's|"motadatagosdk/|"dev.azure.com/...|g' {} +
```

## Answers (from interview)
- **Why vendor the SDK?** — SDK is not hosted on the internet yet. Vendoring + patching is the only option (see [[sdk/module-path-patching]])
- **Project name from comment** — `head -1 src/main.go | sed 's|^// *||'` extracts the name (e.g., `motadata_user_management`). This name gets baked into the Dockerfile (`--build-arg APP_NAME`) and docker-compose.yml. Using the first-line comment is a convention — the DSL's `Main` concept generates it. A config file would be cleaner but this works and is generated automatically.
- **Service name propagation**: The extracted name flows through the entire stack: main.go comment → sync.sh → Docker build arg → container name. Single source of truth is the DSL's `Main` concept.

Links: [[tooling/merge-hooks-story]] [[sdk/module-path-patching]] [[deployment/docker-nats-setup]]
