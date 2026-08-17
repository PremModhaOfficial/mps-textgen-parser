---
title: "Validation Report: Testing Category"
validated: 2026-03-20
---

# Testing Category Validation

## Node: demo-sh-e2e.md

| # | Claim | Verdict | Evidence |
|---|-------|---------|----------|
| 1 | NATS CLI has `nats req` command | VERIFIED | Official NATS CLI tool (github.com/nats-io/natscli). `nats req <subject> <payload>` sends request-reply messages. |
| 2 | NATS CLI has `nats reply` command | VERIFIED | `nats reply <subject> <response>` creates a responder that replies to requests on the given subject. Supports wildcards (`>`). |
| 3 | NATS CLI has `nats sub` command | VERIFIED | `nats sub <subject>` subscribes and prints received messages. Used for traffic observation. |
| 4 | NATS request-reply is a real messaging pattern | VERIFIED | Core NATS pattern documented at docs.nats.io. One-to-one request with timeout, built into the protocol. |
| 5 | `$SRV.INFO` is a valid NATS micro service discovery subject | VERIFIED | NATS micro framework (nats.go/micro) exposes service metadata on `$SRV.INFO`, `$SRV.PING`, and `$SRV.STATS` subjects. Part of the NATS services API. |
| 6 | Mock responders via `nats reply` with wildcards is a valid testing pattern | VERIFIED | `nats reply "motadata.*.db.>" '{"status":"ok"}'` is valid — the `>` wildcard matches any number of tokens. Common pattern for mocking downstream services in NATS-based architectures. |
| 7 | Docker Compose port mapping 4229→4222 | VERIFIED | Standard `HOST:CONTAINER` port syntax. External port 4229 maps to internal NATS port 4222. |

## Summary
- **7/7 claims VERIFIED**
- No corrections needed
- The mock DAL pattern using `nats reply` with wildcards is a pragmatic and valid approach for integration testing NATS microservices without a real database
