---
title: "Docker Compose + NATS Setup"
category: deployment
project: v3
created: 2026-03-20
---

The generated microservice runs via Docker Compose with a NATS server.

## docker-compose.yml
Two services:
- **nats**: `nats:2.10-alpine` with JetStream enabled, healthcheck on `:8222/healthz`
- **app**: Built from `Dockerfile`, connects to `nats://nats:4222` internally

## Port Mapping
- `4229:4222` — NATS client (external 4229, internal 4222)
- `8222:8222` — NATS HTTP monitoring

The app uses `NATS_URL=nats://nats:4222` (internal Docker network).
`demo.sh` connects to `nats://localhost:4229` (host-side port).

## Answers (from interview)
- **Port 4229**: Local conflict on 4222 — trivial detail, just picked an available port
- **Healthcheck timing**: No significant issues — Docker Compose `service_healthy` condition handles the startup order
- **Port conflicts**: `sync.sh` runs `fuser -k 4229/tcp 8222/tcp 2>/dev/null || true` to kill any lingering processes before starting fresh

Links: [[tooling/sync-pipeline]] [[testing/demo-sh-e2e]] [[projects/v3-user-managment-mps]]
